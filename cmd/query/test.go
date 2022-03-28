package main

import (
	"fmt"
	"projects/xinhuatool/pkg/login"
	"projects/xinhuatool/pkg/page"
	"projects/xinhuatool/util"
	"time"
)

func test() {
	// 测试user-agent
	ua := util.UA.Get()
	fmt.Println(ua)

	// 测试获取时间点
	EndDay := util.Datetime
	fmt.Println("截止日期：", EndDay)

	// 测试设置值，设置成昨天
	now := time.Now()
	tm := now.Add(-(time.Hour * 24))
	util.SetDatetime(tm)

	// 测试获取账户
	fmt.Println(util.Account, ": ", util.Passwd)

	// 测试登录功能
	header := map[string]string{"User-Agent": ua}
	sess := login.NewSess(header, "", 20)
	sess.LoginGuest()
	// fmt.Println(sess)

	// 获取检索页面
	page.HistoryFrame(sess)

	// 进行检索
	kw := "春节"
	participle := "on"
	pg := 1
	tp := page.IMG
	showDataFlag := page.SHOWALL
	resp, err := page.SearchKeyword(sess, kw, participle, pg, tp, showDataFlag)
	if err != nil {
		panic("检索关键词未得到正确响应.")
	}
	al := page.NewArticleListFromResp(resp)
	al.GetMeta()
	// fmt.Println(al.Metas())
}
