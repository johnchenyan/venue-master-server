package cmd

import (
	"fmt"
	"github.com/deatil/lakego-doak-admin/admin/controller"
	"github.com/deatil/lakego-doak-admin/admin/model"
	"github.com/deatil/lakego-doak/lakego/command"
	"github.com/xuri/excelize/v2"
	"golang.org/x/xerrors"
	"strconv"
	"time"
)

var ImportDataToSqlCmd = &command.Command{
	Use:          "lakego-admin:import-data-to-sql",
	Short:        "lakego-admin import-data-to-sql.",
	Example:      "{execfile} lakego-admin:import-data-to-sql",
	SilenceUsage: true,
	PreRun: func(cmd *command.Command, args []string) {

	},
	Run: func(cmd *command.Command, args []string) {
		err := ImportDataToSql()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	},
}

var excelPath string

func init() {
	pf := ImportDataToSqlCmd.Flags()
	pf.StringVarP(&excelPath, "excelPath", "e", "", "数据文件")

	command.MarkFlagRequired(pf, "excelPath")
}

func ImportDataToSql() error {
	println("start import data to sql...")
	if excelPath == "" {
		return xerrors.Errorf("文件路径不能为空")
	}

	f, err := excelize.OpenFile(excelPath)
	if err != nil {
		return err
	}

	// 假设数据在第一个工作表中
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return err
	}

	layout := "01/02/2006" // 日期格式，Go 使用固定的模板表示法

	// 创建 map 来存储 SettlementPointName 和 SettlementPointType 的集合
	settlementMap := make(map[string]map[string]struct{})

	for i, row := range rows {
		if i == 0 {
			continue
		}

		// 解析字符串为 time.Time 类型
		date, err := time.Parse(layout, row[0])
		if err != nil {
			fmt.Println("解析日期错误:", err)
			continue
		}

		// 转换字符串为 uint8
		deliveryHour, err := strconv.ParseUint(row[1], 10, 8) // 基数10，位数8
		if err != nil {
			fmt.Println("转换错误:", err)
			continue
		}

		deliveryInterval, err := strconv.ParseUint(row[2], 10, 8) // 基数10，位数8
		if err != nil {
			fmt.Println("转换错误:", err)
			continue
		}

		settlementPointPrice, err := strconv.ParseFloat(row[7], 64) // 基数10，位数8
		if err != nil {
			fmt.Println("转换错误:", err)
			continue
		}

		data := model.SettlementData{
			SettlementPointName:  row[4],
			SettlementPointType:  row[5],
			DeliveryDate:         date,
			DeliveryHour:         uint8(deliveryHour),
			DeliveryInterval:     uint8(deliveryInterval),
			SettlementPointPrice: settlementPointPrice,
		}

		// 在集合中插入 SettlementPointName 和 SettlementPointType
		if _, exists := settlementMap[data.SettlementPointName]; !exists {
			settlementMap[data.SettlementPointName] = make(map[string]struct{})
		}
		settlementMap[data.SettlementPointName][data.SettlementPointType] = struct{}{}

		err = controller.CreateSettlementData(data)
		if err != nil {
			fmt.Println("插入错误:", err)
			continue
		}
	}

	for pointName, pointTypes := range settlementMap {
		for pointType, _ := range pointTypes {
			err := controller.CreateSettlementPoint(model.SettlementPoint{
				SettlementPointName: pointName,
				SettlementPointType: pointType,
			})
			if err != nil {
				println(err.Error())
				continue
			}
		}
	}

	println("import data to sql success.")

	return nil
}
