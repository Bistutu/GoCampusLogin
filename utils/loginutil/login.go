package loginutil

import (
	"crypto/tls"
	"errors"
	"log"
	"net/http"
	"net/url"

	"github.com/gocolly/colly"

	"GoCampusLogin/utils"
	"GoCampusLogin/utils/captureutil"
	"GoCampusLogin/utils/encryutil"
	"GoCampusLogin/utils/httputil"
)

const (
	baseUrl  = "https://wxjw.bistu.edu.cn/authserver/login" // BISTU 教务网登录页面
	loginUrl = "https://wxjw.bistu.edu.cn/authserver/login?service=https://jwxt.bistu.edu.cn:443/jwapp/sys/emaphome/portal/index.do"
	maxRetry = 3 // 最大重试次数
)

// Login 登录模块
func Login(username, password string) ([]*http.Cookie, error) {
	// 登录请求，最多重试 3 次
	for i := 0; i < maxRetry; i++ {
		cookies, err := login(username, password)
		// 如果请求成功则直接返回
		if err == nil {
			return cookies, nil
		}
		log.Fatalf("username: %v, login fail and current count: %d, err: %v", username, i, err)
	}
	return nil, errors.New("达到最大重试次数，登录失败！")
}

// 内部登录模块
func login(username, password string) ([]*http.Cookie, error) {
	params, password, err := getParams(username, password)
	if err != nil {
		return nil, err
	}
	// 发起登录请求
	header := http.Header{}
	header.Add("Origin", "http://wxjw.bistu.edu.cn")
	header.Add("Referer", "http://wxjw.bistu.edu.cn/authserver/login?service=http://jwxt.bistu.edu.cn/jwapp/sys/emaphome/portal/index.do")
	resp, err := httputil.PostForm(loginUrl, header, params)
	if err != nil {
		log.Fatalf("登录失败: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	return resp.Cookies(), nil
}

// 构建登录前的请求参数
func getParams(username string, password string) (url.Values, string, error) {
	var encryptKey, execution string
	// params 登录请求参数
	params := url.Values{}

	// 创建 goColly
	c := colly.NewCollector()
	c.AllowURLRevisit = true
	// 忽略证书错误
	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})
	// 设置回调器，获取AES加密密钥 pwdEncryptSalt 和认证标识 execution
	c.OnHTML("#execution,#pwdEncryptSalt", func(e *colly.HTMLElement) {
		switch e.Attr("id") {
		case "execution":
			execution = e.Attr("value")
		case "pwdEncryptSalt":
			encryptKey = e.Attr("value")
		}
	})
	// 请求登录页
	err := c.Visit(loginUrl)
	if err != nil {
		log.Fatalf("fail to visit login web: %v", err)
		return nil, "", err
	}

	// 对密码进行 AES-CBC-PKCS7Padding 加密
	password, err = encryutil.CBCEncrypt([]byte(utils.RandomNString(64)+password), []byte(encryptKey))
	if err != nil {
		return nil, "", err
	}

	// 将 colly 的 cookie 转移到 golang http
	httputil.AddCookie(baseUrl, c.Cookies(baseUrl))
	// 准备登录请求参数
	params["username"] = []string{username}
	params["password"] = []string{password}
	params["_eventId"] = []string{"submit"}
	params["cllt"] = []string{"userNameLogin"}
	params["dllt"] = []string{"generalLogin"}
	params["lt"] = []string{""}
	params["execution"] = []string{execution}
	// 判断是否需要验证码
	var code string
	if captureutil.IsNeedCaptcha(username) {
		retryCount := 0
		for len(code) != 4 && retryCount < 8 { // 验证码最大识别次数为 8
			code = captureutil.GetCaptureCode()
			retryCount++
		}
		if len(code) != 4 {
			return nil, "", errors.New("验证码识别失败")
		}
	}
	params["captcha"] = []string{code}
	return params, password, nil
}
