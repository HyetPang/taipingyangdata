package excel

import (
	"os"
	"github.com/hyetpang/taipingyangdata/model"
	"github.com/tealeg/xlsx"
)

var excelFile *xlsx.File

func init() {
	excelFile = xlsx.NewFile()
	sheet1, err := excelFile.AddSheet("sheet1")
	if err != nil {
		panic(err)
	}
	row := sheet1.AddRow()
	row.AddCell().SetString("客户姓名")
	row.AddCell().SetString("身份证")
	row.AddCell().SetString("保单号")
	row.AddCell().SetString("产品名字")
	row.AddCell().SetString("购买日期")
	row.AddCell().SetString("到期日期")
	row.AddCell().SetString("缴费")
	row.AddCell().SetString("缴费年限")
	row.AddCell().SetString("地址")
	row.AddCell().SetString("手机号")
	row.AddCell().SetString("保单类型")
}

// SaveData 保存数据
func SaveData(users []*model.UserData, userDataCap int, ) {
	sheet := excelFile.Sheet["sheet1"]
	for _, v := range users {
		row := sheet.AddRow()
		row.AddCell().SetString(v.Name)
		row.AddCell().SetString(v.IDCard)
		row.AddCell().SetString(v.PolicyNum)
		row.AddCell().SetString(v.PolicyName)
		row.AddCell().SetString(v.BuyDate)
		row.AddCell().SetString(v.ExpireDate)
		row.AddCell().SetString(v.HowMuch)
		row.AddCell().SetString(v.BuyedYear)
		row.AddCell().SetString(v.Address)
		row.AddCell().SetString(v.PhoneNum)
		row.AddCell().SetString(v.PolicyType)
	}
	users = make([]*model.UserData, 0, userDataCap)
}

// WriteToDisk 写到磁盘
func WriteToDisk(file *os.File) {
	excelFile.Write(file)
}