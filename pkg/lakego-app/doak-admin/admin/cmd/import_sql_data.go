package cmd

import (
	"fmt"
	"github.com/deatil/lakego-doak-admin/admin/controller"
	"github.com/deatil/lakego-doak-admin/admin/model"
	"github.com/deatil/lakego-doak/lakego/command"
	"github.com/xuri/excelize/v2"
	"golang.org/x/xerrors"
	"strconv"
	"strings"
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
var excelType string

func init() {
	pf := ImportDataToSqlCmd.Flags()
	pf.StringVarP(&excelPath, "excelPath", "e", "", "数据文件")
	pf.StringVarP(&excelType, "excelType", "t", "", "数据类型")

	command.MarkFlagRequired(pf, "excelPath")
	command.MarkFlagRequired(pf, "excelType")
}

func ImportDataToSql() error {
	println("start import data to sql...")
	if excelPath == "" {
		return xerrors.Errorf("文件路径不能为空")
	}

	if excelType == "" {
		return xerrors.Errorf("文件类型不能为空")
	}

	if excelType == "normal" {
		err := saveSettlementData()
		if err != nil {
			return err
		}
	} else if excelType == "special" {
		err := saveSettlementDataT()
		if err != nil {
			return err
		}
	} else {
		return xerrors.Errorf("错误的文件类型")
	}

	println("import data to sql success.")

	return nil
}

func saveSettlementData() error {
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
	return nil
}

func saveSettlementDataT() error {
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
	settlementMap := make(map[string]struct{})

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

		fmt.Println(row[0], row[1])

		// 转换字符串为 uint8
		deliveryHour, err := convertTimeStringToNumber(row[1]) // 基数10，位数8
		if err != nil {
			fmt.Println("转换错误:", err)
			continue
		}

		settlementPointPrice, err := strconv.ParseFloat(strings.TrimSpace(row[5]), 64) // 基数10，位数8
		if err != nil {
			fmt.Println("转换错误:", err)
			continue
		}

		data := model.SettlementDataT{
			SettlementPointName:  row[3],
			DeliveryDate:         date,
			DeliveryHour:         uint8(deliveryHour),
			SettlementPointPrice: settlementPointPrice,
		}

		// 在集合中插入 SettlementPointName
		if _, exists := settlementMap[data.SettlementPointName]; !exists {
			settlementMap[data.SettlementPointName] = struct{}{}
		}

		err = controller.CreateSettlementDataT(data)
		if err != nil {
			fmt.Println("插入错误:", err)
			continue
		}
	}

	for pointName, _ := range settlementMap {
		err := controller.CreateSettlementPointT(model.SettlementPointT{
			SettlementPointName: pointName,
		})
		if err != nil {
			println(err.Error())
			continue
		}
	}

	return nil
}

func convertTimeStringToNumber(timeStr string) (int, error) {
	// 分割字符串，获取小时部分
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid time format")
	}

	// 将小时部分转换为整数
	hourStr := parts[0]
	hour, err := strconv.Atoi(hourStr)
	if err != nil {
		return 0, err
	}

	// 检查小时范围
	if hour < 0 || hour > 24 {
		return 0, fmt.Errorf("hour must be between 0 and 24")
	}

	return hour - 1, nil
}
