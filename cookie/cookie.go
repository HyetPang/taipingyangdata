package cookie

import (
	"os"
	"errors"
	"regexp"
	"io"
	"strings"
	"time"
	"net/http"
	"fmt"
)

var reqClient = http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

func getLoginCookies() map[string]string {
	cookies := make(map[string]string)
	// reqClient := http.Client{
	// 	CheckRedirect: func(req *http.Request, via []*http.Request) error {
	// 		return http.ErrUseLastResponse
	// 	},
	// }
	code, err := login(cookies)
	if err != nil {
		panic(err)
	}
	d := fmt.Sprintf("j_username=67-72-68-49-49-57-52-49&j_password=90-121-121-49-50-51-52-53-54-64&j_captcha=%s&codeTime=%d", code, time.Now().UnixNano()/1e6)
	data := strings.NewReader(d)
	request, err := http.NewRequest(http.MethodPost, "https://www.cpic.com.cn/agent-ebw/j_spring_security_check", data)
	if err != nil {
		panic(err)
	}
	for k, v := range cookies {
		request.AddCookie(&http.Cookie{
			Name:  k,
			Value: v,
		})
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r, err := reqClient.Do(request)
	if err != nil {
		panic(err)
	}
	if strings.Contains(r.Request.URL.String(), "https://www.cpic.com.cn/agent-ebw/login.html") {
		url := r.Request.URL.String()
		err := url[len(url)-1:]
		if err == "9" {
			fmt.Println("验证码错误！")
		} else if err == "1" {
			fmt.Println("用户名或密码错误！")
		} else if err == "10" {
			fmt.Println("用户不是营销或区拓类型！")
		} else if err == "2" {
			fmt.Println("验证码已经失效，请重新输入！")
		} else if err == "3" {
			fmt.Println("连续登录失败5次数后锁定账号，请30分钟后尝试！")
		}
		return nil
	}
	for r.StatusCode == http.StatusFound {
		for _, v := range r.Cookies() {
			cookies[v.Name] = v.Value
		}
		// 如果是重定向,开始获取数据
		url, err := r.Location()
		if err == http.ErrNoLocation {
			break
		}
		if err != nil && err != http.ErrNoLocation {
			panic(err)
		}
		request, err := http.NewRequest(http.MethodGet, url.String(), nil)
		if err != nil {
			panic(err)
		}
		var cookie string
		for k, c := range cookies {
			cookie += fmt.Sprintf("%s=%s; ", k, c)
		}
		request.Header.Set("Cookie", cookie[:len(cookie)-2])
		resp, err := reqClient.Do(request)
		if err != nil {
			panic(err)
		}
		r = resp
	}
	req, err := http.NewRequest(http.MethodPost, "https://www.cpic.com.cn/agent-ebw/view/fcRegion/innerAddressDetails.jsp", strings.NewReader("fcRegionCodeId=10848745&_action=fcRegionCodeInfo"))
	if err != nil {
		panic(err)
	}
	var cookie string
	for k, c := range cookies {
		cookie += fmt.Sprintf("%s=%s; ", k, c)
	}
	req.Header.Set("Cookie", cookie[:len(cookie)-2])
	resp, err := reqClient.Do(req)
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, resp.Body)
	resp.Body.Close()
	return cookies
}

func login(cookies map[string]string) (string, error) {
	resp, err := http.Get("https://www.cpic.com.cn/agent-ebw/login.html")
	if err != nil {
		panic(err)
	}
	for _, v := range resp.Cookies() {
		cookies[v.Name] = v.Value
	}
	// 验证码试三次
	i := 0
	var text string
	for matched, _ := regexp.MatchString(`\d{4}`, text); !matched; {
		if i == 3 {
			return "", errors.New("验证码出错")
		}
		text = getImageCode(cookies)
		i++
	}
	return text, nil
}

func getImageCode(cookies map[string]string) string {
	// imgRequ, err := http.NewRequest(http.MethodGet, "https://www.cpic.com.cn/agent-ebw/jcaptcha.jpg", nil)
	// if err != nil {
	// 	panic(err)
	// }
	// for k,v := range cookies {
	// 	imgRequ.AddCookie(&http.Cookie{
	// 		Name:k,
	// 		Value:v,
	// 	})
	// }
	// imgResp, err := http.DefaultClient.Do(imgRequ)
	// if err != nil {
	// 	panic(err)
	// }
	// defer imgResp.Body.Close()
	// f, err := os.Create("img/img.jpg")
	// if err != nil {
	// 	panic(err)
	// }
	// _, err = io.Copy(f, imgResp.Body)
	// defer os.Remove("img/img.jpg")
	// if err != nil {
	// 	panic(err)
	// }
	// client := gosseract.NewClient()
	// client.SetImage("img/img.jpg")
	// text, _ := client.Text()
	// client.Close()
	// f.Close()
	// os.Remove("img/img.jpg")
	// return text
	return ""
}