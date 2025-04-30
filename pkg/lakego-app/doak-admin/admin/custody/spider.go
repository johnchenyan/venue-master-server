package custody

import (
	"fmt"
	"strings"
)

func spider(url string) ([]*HashRateEntry, error) {
	if strings.Contains(url, "antpool") {
		rUrl, err := convertURL(url)
		if err != nil {
			return nil, err
		}

		hss, err := fetchAntPoolRecv(rUrl)
		if err != nil {
			return nil, err
		}

		return hss, nil
	} else if strings.Contains(url, "f2pool") {
		hs, err := fetchF2PoolRecv(url)
		if err != nil {
			return nil, err
		}
		return hs, nil
	}

	return nil, fmt.Errorf("no router path")
}
