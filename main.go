package main

import (
	"sync"
	"github.com/hyetpang/taipingyangdata/crawl"
)



var cookie = "name=value; USERNAME=CHD11941; BRANCH_ORG_NAME=%E5%9B%9B%E5%B7%9D%E5%88%86%E5%85%AC%E5%8F%B8; PRE_NAME=CHD; JSESSIONID=35FB1A4334160D6AEC564E9308962CF7; XSESSION=261_1_xingxiaozhichiwangzhan_agent-ebw.ba37ac2f-9688-11e9-a9be-0242b7149d61"

// 爬虫获取数据
func main() {
	// 初始化文件
	// cookies := getLoginCookies()
	// 获取信息
	var wg sync.WaitGroup
	crawl.StartCrawl(cookie, &wg)
	wg.Wait()
}
