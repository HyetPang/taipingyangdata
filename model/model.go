package model

// UserData 客户姓名，保单号，购买的产品名字，哪年购买的，哪年到期，缴费多少，缴费多少年，手机号，家庭地址
type UserData struct {
	Name       string // 客户姓名
	IDCard     string //  身份证号
	PolicyNum  string // 保单号
	PolicyName string // 购买的产品名字
	BuyDate    string // 购买的日期
	ExpireDate string // 过期日期
	HowMuch    string // 多少钱,
	BuyedYear  string // 缴费多少年
	Address    string // 家庭地址
	PhoneNum   string // 手机号码
	PolicyType string // 保单类型 孤儿单客户，银保客户，车险客户
}