package login

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"projects/xinhuatool/util"
	"regexp"
	"strings"
	"time"
)

// 请求相关的几个url
var LoginURL = "http://home.xinhua-news.com/loginxhs"
var HostURL = "home.xinhua-news.com"
var OriginURL = "http://home.xinhua-news.com"
var QueryURL = "http://home.xinhua-news.com/history/page"
var QueryReferer = "http://home.xinhua-news.com/work/search"

var GuestLoginURL = "http://home.xinhua-news.com/loginxhs/visitorSSOLogin"

type Session struct {
	Header  map[string]string
	Client  *http.Client
	Birth   time.Time
	Account string
	Passwd  string
}

func NewSess(header map[string]string, proxyURL string, timeout int, auth ...string) *Session {
	sess := &Session{}

	// Header
	sess.Header = header

	// Client
	client := &http.Client{
		Timeout: time.Second * time.Duration(timeout),
	}
	if proxyURL != "" {
		proxy := func(_ *http.Request) (*url.URL, error) {
			return url.Parse(proxyURL)
		}
		transport := &http.Transport{Proxy: proxy}
		client.Transport = transport
	}
	sess.Client = client

	// Birth
	sess.Birth = time.Now()

	// auth
	if len(auth) == 2 {
		sess.Account = auth[0]
		sess.Passwd = auth[1]
	}

	return sess
}

// cookieKV从Set-Cookie字段取出key和value
func cookieKV(s string) string {
	segs := strings.Split(s, ";")
	for i, seg := range segs {
		segs[i] = strings.TrimSpace(seg)
	}

	segs = strings.Split(segs[0], "=")
	if len(segs) >= 2 {
		return segs[0] + "=" + segs[1]
	}
	return ""
}

// SetCookie从response中获取Set-Cookie字段，设置Session
func (sess *Session) SetCookie(resp *http.Response) {
	segs, ok := resp.Header["Set-Cookie"]
	if !ok {
		return
	}

	var cookies []string
	for _, seg := range segs {
		cookies = append(cookies, cookieKV(seg))
	}
	cookie := strings.Join(cookies, "; ")
	if sess.Header["Cookie"] == "" {
		sess.Header["Cookie"] = cookie
	} else {
		sess.Header["Cookie"] += "; " + cookie
	}
}

func RawGuestPasswd(s string) string {
	pat := regexp.MustCompile(`<input type="hidden" id="anonymous_p" value="(.*?)"`)
	res := pat.FindStringSubmatch(s)
	if len(res) >= 2 {
		return res[1]
	}

	return ""
}

// 获取guest原始密码，初步设置cookie。所有cookie不以url区分，赌不会出现碰撞情况
func (sess *Session) LoginFrame() (string, error) {
	resp, err := util.GetWithHead(LoginURL, sess.Header, nil, sess.Client)
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()

	// 设置Cookie
	sess.SetCookie(resp)

	// 获取anonymous_p字段
	byteData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	rawPasswd := RawGuestPasswd(string(byteData))

	return rawPasswd, nil
}

// 获取完整password
func Passwd(raw string) string {
	s := raw + "_xhs"
	byteData := make([]byte, 1024)
	encoding := base64.StdEncoding
	encoding.Encode(byteData, []byte(s))

	idx := bytes.IndexByte(byteData, 0)
	byteData = byteData[:idx]

	seg := string(byteData)
	passwd := seg + "_pwdtag"
	return passwd
}

// 登录Guest
func (sess *Session) LoginGuest() error {
	// 准备工作
	raw, _ := sess.LoginFrame()
	sess.Account = "xhqmgonggong"
	sess.Passwd = Passwd(raw)

	data := map[string][]string{
		"loginReferer": {""},
		"userName":     {sess.Account},
		"password":     {sess.Passwd},
		"ssoPageType":  {"SSO_PRIVATE_PAGE"},
		"timezone":     {"Asia/Shanghai"},
	}

	sess.Header["Host"] = HostURL
	sess.Header["Referer"] = LoginURL
	sess.Header["Origin"] = OriginURL
	sess.Header["Content-Type"] = "application/x-www-form-urlencoded"

	fmt.Println("即将登录Guest用户：")
	resp, err := PostHandleRedict(GuestLoginURL, sess, data)
	if err != nil {
		fmt.Println("Guest用户登录失败：", err)
		return err
	}

	sess.SetCookie(resp)

	fmt.Println("Guest用户登录成功.")
	return nil
}

// 登录DBK
func (sess *Session) LoginDBK() error {

	return nil
}

func PostHandleRedict(url string, sess *Session, data url.Values) (*http.Response, error) {
	var client *http.Client
	if sess.Client == nil {
		client = &http.Client{}
	} else {
		client = sess.Client
	}

	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	// jar, _ := cookiejar.New(nil)
	// client.Jar = jar

	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	for k, v := range sess.Header {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

Redirect:
	if resp.StatusCode != 200 && resp.StatusCode != 302 {
		err := fmt.Errorf("重定向跳转失败")
		return nil, err
	}

	if resp.StatusCode == 200 {
		return resp, nil
	}

	if resp.StatusCode == 302 {
		location := resp.Header.Get("location")
		if location == "" {
			err := fmt.Errorf("重定向location获取为空")
			return nil, err
		}
		// fmt.Println("正在跳转：", location)

		setCookie := resp.Header.Get("Set-Cookie")
		if setCookie != "" {
			sess.SetCookie(resp)
		}

		// 请求重定向
		req, err := http.NewRequest("GET", location, nil)
		if err != nil {
			return nil, err
		}
		for k, v := range sess.Header {
			req.Header.Set(k, v)
		}
		resp, err = client.Do(req)
		if err != nil {
			return nil, err
		}

		goto Redirect
	}

	err = fmt.Errorf("重定向跳转失败")
	return nil, err
}
