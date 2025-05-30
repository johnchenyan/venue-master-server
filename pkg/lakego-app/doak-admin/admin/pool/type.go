package pool

type HashRateEntry struct {
	LastDayHashRate string
	LastDayHashUnit string
	LastDayRecv     string
	LastDayTime     string
}

type FBProfitEntry struct {
	LastDayRecv string
	LastDayTime string
}

type RealTimeStatus struct {
	RealTimeHash string
	HSLast1D     string
	HashUnit     string
	Online       int
	Offline      int
}
