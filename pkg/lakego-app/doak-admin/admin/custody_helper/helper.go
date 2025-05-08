package custody_helper

import (
	"github.com/deatil/lakego-doak-admin/admin/model"
	"strconv"
)

func TotalEnergy(lastDayHashRate, lastDayHashUnit string, custodyInfo model.CustodyInfo) (float64, error) {
	var value float64
	hashRate, err := strconv.ParseFloat(lastDayHashRate, 64)
	if err != nil {
		return 0, err
	}
	switch lastDayHashUnit {
	case "TH/s":
		value = hashRate
	case "PH/s":
		value = hashRate * 1000
	case "EH/s":
		value = hashRate * 1000 * 1000
	}

	energyRatio, err := strconv.ParseFloat(custodyInfo.EnergyRatio, 64)
	if err != nil {
		return 0, err
	}

	return value * energyRatio * 3600 * 24, nil
}

func TotalHostingFee(energy float64, custodyInfo model.CustodyInfo) (float64, error) {
	energyKwh := joulesToKWh(energy)
	basicHostingFee, err := strconv.ParseFloat(custodyInfo.BasicHostingFee, 64)
	if err != nil {
		return 0, err
	}
	return basicHostingFee * energyKwh, nil
}

func joulesToKWh(joules float64) float64 {
	return joules / 3600000
}

func TotalIncomeUSD(incomeBtc, averagePrice string) (float64, error) {
	incomeBtcFloat, err := strconv.ParseFloat(incomeBtc, 64)
	if err != nil {
		return 0, err
	}

	averagePriceFloat, err := strconv.ParseFloat(averagePrice, 64)
	if err != nil {
		return 0, err
	}
	return averagePriceFloat * incomeBtcFloat, nil
}
