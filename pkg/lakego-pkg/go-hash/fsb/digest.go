package fsb

import (
    "errors"

    "github.com/deatil/go-hash/whirlpool"
)

// digest represents the partial evaluation of a checksum.
type digest struct {
    /* the parameters of FSB */
    n, w, r, p uint32 /* n multiple of w, r mulitple of 32/64 */
    /* other useful parameters to avoid recomputing */
    b uint32/* number of QC blocks */
    inputsize uint32 /* number of input bits from the message (s-r) per round *MUST* be multiple of 8 */
    bpc, bfiv, bfm uint32 /* for each column : number of input bits, number from the iv, number from the message */
    /* the "first line" of matrix H */
    firstLine [][][]byte
    /* hash length */
    hashbitlen int
    /* current syndrome */
    syndrome []byte
    /* space to store new syndrome */
    newSyndrome []uint32
    /* input buffer */
    buffer []byte
    /* number of bits in the buffer */
    count uint32
    /* number of bits hashed */
    databitlen uint64
}

// New returns a new *digest computing the fsb checksum
func newDigest(hashbitlen int) (*digest, error) {
    if hashbitlen == 0 {
        return nil, errors.New("go-hash/fsb: hash size can't be zero")
    }

    switch hashbitlen {
        case 48, 160, 224, 256, 384, 512:
            break
        default:
            return nil, errors.New("go-hash/fsb: non-byte hash sizes are not supported")
    }

    d := new(digest)
    d.hashbitlen = hashbitlen
    d.Reset()

    return d, nil
}

func (d *digest) Reset() {
    var i,j,k,l int
    var Pi_line []byte

    for i=0; i<NUMBER_OF_PARAMETERS; i++ {
        if d.hashbitlen == parameters[i][0] {
            d.n = uint32(parameters[i][1])
            d.w = uint32(parameters[i][2])
            d.r = uint32(parameters[i][3])
            d.p = uint32(parameters[i][4])
            d.b = d.n/d.r
            d.bpc = uint32(logarithm(d.n/d.w))
            d.inputsize = d.w*d.bpc-d.r
            d.bfiv = d.r/d.w
            d.bfm = d.inputsize/d.w
            d.databitlen = 0

            /* compute the first QC matrix line */
            d.firstLine = make([][][]byte, d.b)
            for k=0; k<int(d.b); k++ {
                d.firstLine[k] = make([][]byte, 8)
                d.firstLine[k][0] = make([]byte, int((d.p+d.r)>>3)+1)

                Pi_line = Pi[k*int((d.p>>3)+1):]
                for j=0; j<int(d.p>>3); j++ {
                    d.firstLine[k][0][int(d.r>>3)+j] = Pi_line[j]
                }

                d.firstLine[k][0][(d.p+d.r)>>3] = Pi_line[d.p>>3]&(((1<<(d.p&7))-1)<<(8-(d.p&7)))
                for j=0; j<int(d.r>>3); j++ {
                    d.firstLine[k][0][j] = d.firstLine[k][0][int(d.p>>3)+j]<<(d.p&7)
                    d.firstLine[k][0][j] ^= d.firstLine[k][0][int(d.p>>3)+j+1]>>(8-(d.p&7))
                }

                for j=1; j<8; j++ {
                    d.firstLine[k][j] = make([]byte, int((d.p+d.r)>>3)+1)
                    for l=0; l<int(d.p+d.r)>>3; l++ {
                        d.firstLine[k][j][l] ^= d.firstLine[k][0][l] >> j
                        d.firstLine[k][j][l+1] ^= d.firstLine[k][0][l] << (8-j)
                    }
                }
            }

            d.syndrome = make([]byte, LUI(d.r) * 4)
            d.newSyndrome = make([]uint32, LUI(d.r))
            d.buffer = make([]byte, d.inputsize>>3)
            d.count = 0
        }
    }
}

func (d *digest) Size() int {
    return d.hashbitlen >> 3
}

func (d *digest) BlockSize() int {
    return int(d.inputsize >> 3)
}

func (d *digest) Write(data []byte) (nn int, err error) {
    nn = len(data)

    databitlen := uint32(8 * len(data))
    d.write(data, databitlen)

    return
}

func (d *digest) write(data []byte, databitlen uint32) (err error) {
    var tmp, i uint32
    var remaining byte

    /* we check if this Update will fill one buffer */
    if databitlen + d.count < d.inputsize {
        /* we simply need to copy data to the buffer. Either it is aligned or not. */
        if (d.count & 7) == 0 {
            if databitlen > 0 {
                copy(d.buffer[d.count>>3:], data[:((databitlen-1)>>3) + 1])
            }

            d.databitlen += uint64(databitlen)
            d.count += databitlen
            return nil
        } else {
            d.buffer[d.count>>3] ^= d.buffer[d.count>>3] & ((1<<(8-(d.count&7)))-1)
            for i = 0; i <= databitlen>>3; i++ {
                d.buffer[d.count>>3+i] ^= data[i]>>(d.count&7)
                d.buffer[d.count>>3+i+1] = data[i]<<(8-(d.count&7))
            }

            d.databitlen += uint64(databitlen)
            d.count += databitlen
            return nil
        }
    } else {
        /* we fill up the buffer, perform a hash and recursively call Update */
        if (d.count & 7) == 0 {
            tmp = d.inputsize - d.count
            copy(d.buffer[d.count>>3:], data[:tmp>>3])

            d.databitlen += uint64(tmp)
            d.count += tmp
            d.performHash()

            return d.write(data[tmp>>3:], databitlen - tmp)
        } else {
            /* tmp contains the number of bits we have to read to fill up the buffer */
            tmp = d.inputsize - d.count

            d.buffer[d.count>>3] ^= d.buffer[d.count>>3] & ((1<<(8-(d.count&7)))-1)
            for i=0; i<(tmp>>3); i++ {
                d.buffer[d.count>>3+i] ^= data[i]>>(d.count&7)
                d.buffer[d.count>>3+i+1] = data[i]<<(8-(d.count&7))
            }

            d.buffer[(d.inputsize>>3)-1] ^= data[tmp>>3]>>(d.count&7)

            /* perform this round's hash */
            d.databitlen += uint64(tmp)
            d.count += tmp
            d.performHash()

            /* we check if there are still some bits to input */
            if databitlen > tmp {
                /* we check if these bits are stored in more than the end of the byte data[tmp>>3] already read */
                if databitlen > (((tmp>>3)+1)<<3) {
                    /* we first re-input the remaining bits in data[tmp>>3] then perform the recursive call */
                    remaining = byte(uint32(data[tmp>>3]) << (tmp&7))
                    d.write([]byte{remaining}, uint32(8-(tmp&7)))

                    return d.write(data[tmp>>3:], databitlen - (((tmp>>3)+1)<<3))
                } else {
                    /* we simply input the remaining bits of data[tmp>>3] */
                    remaining = byte(uint32(data[tmp>>3]) << (tmp&7))

                    return d.write([]byte{remaining}, databitlen - tmp)
                }
            } else {
                return nil
            }
        }
    }

    return
}

func (d *digest) Sum(in []byte) []byte {
    // Make a copy of d so that caller can keep writing and summing.
    d0 := d.copy()
    hash := d0.checkSum()
    return append(in, hash...)
}

func (d *digest) checkSum() []byte {
    var padding, whirlOutput []byte
    var i int

    databitlen := d.databitlen
    if d.count+65 > d.inputsize {
        padding = make([]byte, int(d.inputsize>>3))
        padding[0] = 1<<7
        d.write(padding, d.inputsize-d.count)

        padding[0] = 0
        for i = 0; i < 8; i++ {
            padding[int(d.inputsize>>3)-1-i] = byte(databitlen>>(8*i))
        }

        d.write(padding, d.inputsize)
    } else {
        padding = make([]byte, int((d.inputsize-d.count)>>3)+1)
        padding[0] = 1<<7;
        d.write(padding, d.inputsize-d.count-64)

        for i = 0; i < 8; i++ {
            padding[7-i] = byte(databitlen>>(8*i))
        }

        d.write(padding, 64)
    }

    /* The final round of FSB is finished, now we simply apply the final transform: Whirlpool */
    whirlpoolState := whirlpool.New()
    whirlpoolState.Write(d.syndrome[:d.r >> 3])
    whirlOutput = whirlpoolState.Sum(nil)

    hashval := make([]byte, d.hashbitlen>>3)
    for i = 0; i < (d.hashbitlen>>3); i++ {
        hashval[i] = whirlOutput[i]
    }

    return hashval
}

func (d *digest) performHash() {
    var i,j,index,bidx,tmp int
    var temp []uint32

    for i := range d.newSyndrome {
        d.newSyndrome[i] = 0
    }

    for i=0; i<int(d.w); i++ {
        index = i<<d.bpc
        switch d.bfiv {
            case 2:
                index ^= int(d.syndrome[i>>2]>>(6-((i&3)<<1)))&3
                break;
            case 4:
                index ^= int(d.syndrome[i>>1]>>(4-((i&1)<<2)))&15
                break;
            case 8:
                index ^= int(d.syndrome[i])
                break;
            default:
                tmp = (i+1)*int(d.bfiv)
                for j = i*int(d.bfiv); j < tmp; j++ {
                    index ^= ((int(d.syndrome[j>>3])>>(7-(j&7)))&1)<<(tmp-j-1)
                }
        }

        switch (d.bfm) {
            case 2:
                index ^= int((uint32((d.buffer[i>>2]>>(6-((i&3)<<1)))&3)) << d.bfiv)
                break;
            case 4:
                index ^= int((uint32((d.buffer[i>>1]>>(4-((i&1)<<2)))&15)) << d.bfiv)
                break;
            case 8:
                index ^= int(uint32(d.buffer[i]) << d.bfiv)
                break;
            default:
                tmp = (i+1)*int(d.bfm)
                for j = i*int(d.bfm); j < tmp; j++ {
                    index ^= int((d.buffer[j>>3]>>(7-(j&7)))&1)<<(tmp-j-1+int(d.bfiv))
                }
        }

        bidx = index/int(d.r) /* index of the vector */
        index = index - bidx*int(d.r) /* shift to perform on the vector */

        /* we have finished computing the vector index and shift, now we XOR it! */
        temp = bytesToUints(d.firstLine[bidx][index&7][int(d.r>>3)-(index>>3):])
        for j = int(d.r)/(4<<3)-1; j >= 0; j-- {
            d.newSyndrome[j] ^= temp[j]
        }
    }

    temp = d.newSyndrome;
    d.newSyndrome = bytesToUints(d.syndrome)
    d.syndrome = uintsToBytes(temp)
    d.count = 0
}

func (d *digest) copy() *digest {
    d0 := &digest{}

    d0.n, d0.w, d0.r, d0.p = d.n, d.w, d.r, d.p
    d0.b = d.b

    d0.inputsize = d.inputsize
    d0.bpc, d0.bfiv, d0.bfm = d.bpc, d.bfiv, d.bfm

    d0.firstLine = make([][][]byte, len(d.firstLine))
    copy(d0.firstLine, d.firstLine)

    d0.hashbitlen = d.hashbitlen

    d0.syndrome = make([]byte, len(d.syndrome))
    copy(d0.syndrome, d.syndrome)

    d0.newSyndrome = make([]uint32, len(d.newSyndrome))
    copy(d0.newSyndrome, d.newSyndrome)

    d0.buffer = make([]byte, len(d.buffer))
    copy(d0.buffer, d.buffer)

    d0.count = d.count
    d0.databitlen = d.databitlen

    return d0
}
