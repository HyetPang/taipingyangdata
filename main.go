package main

import (
	"os"
	"github.com/google/uuid"
	"github.com/hyetpang/taipingyangdata/excel"
	"github.com/hyetpang/taipingyangdata/crawl"
)



var cookie = "name=value; USERNAME=CHD11941; BRANCH_ORG_NAME=%E5%9B%9B%E5%B7%9D%E5%88%86%E5%85%AC%E5%8F%B8; PRE_NAME=CHD; JSESSIONID=0D3433FAE4AB39F7E17EF95C23B0A289; XSESSION=261_1_xingxiaozhichiwangzhan_agent-ebw.ba218c1b-9688-11e9-a9be-0242b7149d61"

// 爬虫获取数据
func main() {
	// 初始化文件
	// cookies := getLoginCookies()
	// 获取信息
	crawl.CrawlUsers(cookie)
	file, _ := os.Create(uuid.New().String() + ".xlsx")
	excel.WriteToDisk(file)
}
