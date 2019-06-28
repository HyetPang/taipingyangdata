package crawl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/hyetpang/taipingyangdata/excel"
	"github.com/hyetpang/taipingyangdata/model"
)

// CrawlUsers 获取用户信息
func CrawlUsers(cookie string) {
	userDataCap := 100
	users := make([]*model.UserData, 0, userDataCap)
	// TODO: 从参数设置cookie
	totalPage := 1
	curPage := 1
	for curPage <= totalPage {
		req, err := http.NewRequest(http.MethodPost, "https://www.cpic.com.cn/agent-ebw/view/fcRegion/fcRegionCust.jsp", strings.NewReader(fmt.Sprintf("currentPageIndex=%d&pageSize=10&autoload=true&apname=&apid=&minBirthMonth=&maxBirthMonth=&minCreateDate=&maxCreateDate=&isTelVist=&isPolicyReview=&isVist=&signInsPolicyFlg=&sortOrder=&matchStartDate=&matchEndDate=&custFlag=&regionNumber=CDSP000353&_action=findByParams", curPage)))
		if err != nil {
			panic(err)
		}
		req.Header.Set("Cookie", cookie)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}
		// 每一页有10条数据
		data := make(map[string]interface{})
		var buf bytes.Buffer
		io.Copy(&buf, resp.Body)
		resp.Body.Close()
		newstring := strings.Replace(buf.String(), "'", "\"", -1)
		err = json.Unmarshal([]byte(newstring), &data)
		// err = json.Unmarshal(buf.Bytes(), data)
		if err != nil {
			panic(err)
		}
		curPage++
		result := data["result"].(map[string]interface{})
		// 获取详情
		totalPage, _ = strconv.Atoi(result["totalPage"].(string))
		// 获取数据
		items := result["items"].([]interface{})
		for _, v := range items {
			var u model.UserData
			a := v.(map[string]interface{})
			u.Name = a["apname"].(string)
			u.IDCard = a["apid"].(string)
			// 获取手机号，家庭地址等信息
			phoneAndAddrReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://www.cpic.com.cn/agent-ebw/view/fccustomer/orphan/myOrphanInfo.jsp?policyHorderNo=%s=&signInsPolicyFlg=0&custType=&isSprerec=", a["clsaaname"].(string)), nil)
			if err != nil {
				panic(err)
			}
			phoneAndAddrReq.Header.Set("Cookie", cookie)
			phoneAndAddrReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			phoneAndAddrResp, err := http.DefaultClient.Do(phoneAndAddrReq)
			if err != nil {
				panic(err)
			}
			doc, err := goquery.NewDocumentFromReader(phoneAndAddrResp.Body)
			if err != nil {
				panic(err)
			}
			phone := doc.Find("body > table.open_table > tbody > tr:nth-child(3) > td:nth-child(6)").Text()
			address := doc.Find("body > table.open_table > tbody > tr:nth-child(5) > td:nth-child(2)").Text()
			phoneAndAddrResp.Body.Close()
			u.PhoneNum = phone
			u.Address = address
			// 获取保单信息
			policyCurPage := 1
			policyTotalPage := 1
			for policyCurPage <= policyTotalPage {
				policyReq, err := http.NewRequest(http.MethodPost, "https://www.cpic.com.cn/agent-ebw/view/fccustomer/orphan/myOrphanPolicys.jsp", strings.NewReader(fmt.Sprintf(fmt.Sprintf("currentPageIndex=%d&pageSize=5&autoload=true&asysc=sync&policyHorderNo=%s&_action=findByParams", policyCurPage, u.IDCard))))
				if err != nil {
					panic(err)
				}
				policyReq.Header.Set("Cookie", cookie)
				policyReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				policyResp, err := http.DefaultClient.Do(policyReq)
				if err != nil {
					panic(err)
				}
				var buf bytes.Buffer
				io.Copy(&buf, policyResp.Body)
				policyResp.Body.Close()
				newstring := strings.Replace(buf.String(), "'", "\"", -1)
				policyDataWrap := make(map[string]interface{})
				err = json.Unmarshal([]byte(newstring), &policyDataWrap)
				if err != nil {
					panic(err)
				}
				policyData := policyDataWrap["result"].(map[string]interface{})
				policyTotalPage, _ = strconv.Atoi(policyData["totalPage"].(string))
				policyCurPage++
				// 处理数据
				policyItems := policyData["items"].([]interface{})
				for _, v := range policyItems {
					var realU model.UserData
					realU.PhoneNum = u.PhoneNum
					realU.IDCard = u.IDCard
					realU.Address = u.Address
					realU.Name = u.Name
					item := v.(map[string]interface{})
					// 获取保单详情
					realU.PolicyNum = item["policyNo"].(string)
					realU.PolicyName = item["insuranceName"].(string)
					policyDetailsReq, err := http.NewRequest(http.MethodPost, "https://www.cpic.com.cn/agent-ebw/view/fccustomer/orphan/myOrphanPolicys.jsp", strings.NewReader(fmt.Sprintf("policyNo=%s&insuranceCode=%s&_action=findPolicyDetail", realU.PolicyNum, item["insuranceCode"].(string))))
					if err != nil {
						panic(err)
					}
					policyDetailsReq.Header.Set("Cookie", cookie)
					policyDetailsReq.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
					policyDetailsResp, err := http.DefaultClient.Do(policyDetailsReq)
					if err != nil {
						panic(err)
					}
					var policyDetailsBuf bytes.Buffer
					io.Copy(&policyDetailsBuf, policyDetailsResp.Body)
					policyDetailsResp.Body.Close()
					policyDetailsData := make(map[string]interface{})
					newstring := strings.Replace(policyDetailsBuf.String(), "'", "\"", -1)
					fmt.Println(newstring)
					err = json.Unmarshal([]byte(newstring), &policyDetailsData)
					if err != nil {
						panic(err)
					}
					policyDetailsResult := policyDetailsData["result"].(map[string]interface{})
					fmt.Printf("%+v\n", policyDetailsResult)
					// 判断上面是否获取到手机号，地址等信息，没有获取到就这里重新获取
					if realU.PhoneNum == "" {
						phoneNum := policyDetailsResult["insuredmobile"]
						if phoneNum == nil {
							realU.PhoneNum = "无"
						} else {
							realU.PhoneNum = phoneNum.(string)
						}
					}
					if realU.Address == "" {
						add := policyDetailsResult["payaddress"]
						if add == nil {
							realU.Address = "无"
						} else {
							realU.Address = "(缴费地址)" + add.(string)
						}
					}
					// 投保日期
					realU.BuyedYear = policyDetailsResult["respStartDate"].(string)
					realU.ExpireDate = policyDetailsResult["respendDate"].(string)
					payYeaynum := policyDetailsResult["payYearnum"]
					if payYeaynum  == nil {
						realU.BuyedYear = "无"
					} else {
						realU.BuyedYear = payYeaynum.(string)
					}
					realU.HowMuch = policyDetailsResult["premium"].(string)
					realU.PolicyType = "孤儿单客户"
					// 保存
					users = append(users, &realU)
				}
				if len(users) >= userDataCap {
					// 大于等于限制的容量，开始写入数据
					excel.SaveData(users, userDataCap)
				}
			}
			if len(users) >= userDataCap {
				// 大于等于限制的容量，开始写入数据
				excel.SaveData(users, userDataCap)
			}
		}
		if len(users) >= userDataCap {
			// 大于等于限制的容量，开始写入数据
			excel.SaveData(users, userDataCap)
		}
	}
	excel.SaveData(users, userDataCap)
}
