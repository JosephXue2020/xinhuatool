package main

import (
	"fmt"
	"projects/xinhuatool/pkg/login"
	"projects/xinhuatool/pkg/office"
	_ "projects/xinhuatool/pkg/office"
	"projects/xinhuatool/pkg/page"
	"projects/xinhuatool/util"
	"strconv"
	"time"
)

func main() {
	// 读取数据
	path := "./关键词表.xlsx"
	tasks := readExcel(path)
	fmt.Println("读取关键词表成功.")

	// 干活
	var result = new([][]string)
	worker(tasks, result)
	// fmt.Println(result)
	// for _, v := range *result {
	// 	fmt.Println(v)
	// }

	// 任务完成
	outPath := "./data/检索结果.xlsx"
	saveResult(outPath, *result)
	fmt.Println("Complete!")

	// 退出
	fmt.Println("按任意键退出：")
	var s string
	fmt.Scanln(&s)

	// test()

}

type task struct {
	id      string
	keyword string
	num     int
}

func readExcel(p string) []task {
	sli, err := office.ReadExcel(p, "")
	if err != nil {
		panic("读取关键词表失败.")
	}

	var tasks []task
	for _, item := range sli[1:] {
		var tk task
		tk.id = item[0]
		tk.keyword = item[1]
		if item[2] != "" {
			tk.num, err = strconv.Atoi(item[2])
			if err != nil {
				panic("关键词表limit列错误.")
			}
		} else {
			tk.num = -1
		}
		tasks = append(tasks, tk)
	}

	return tasks
}

func saveResult(p string, res [][]string) error {
	col := []string{
		"keyword",
		"title",
		"gdsid",
		"abstraction",
		"introduction",
		"copyright",
		"tags",
	}
	return office.WriteExcel(p, res, col)
}

func worker(tasks []task, result *[][]string) {
	ua := util.UA.Get()
	header := map[string]string{"User-Agent": ua}
	sess := login.NewSess(header, "", 20)
	sess.LoginGuest()

	page.HistoryFrame(sess)

	// 进行检索
	participle := "on"
	pg := 1
	tp := page.IMG
	showDataFlag := page.SHOWALL

	itemLen := 7
	taskNum := len(tasks)
	var tempResult [][]string
	for i, tk := range tasks {
		fmt.Println("进行任务：", strconv.Itoa(i+1)+"/"+strconv.Itoa(taskNum))
		kw := tk.keyword
		resp, err := page.SearchKeyword(sess, kw, participle, pg, tp, showDataFlag)
		if err != nil {
			item := make([]string, itemLen)
			item[0] = kw
			tempResult = append(tempResult, item)
			fmt.Println(err)
			continue
		}

		al := page.NewArticleListFromResp(resp)
		metas := al.MetaSli()
		for _, meta := range metas {
			item := []string{kw}
			item = append(item, meta...)
			tempResult = append(tempResult, item)
		}
	}

	*result = tempResult
}

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
