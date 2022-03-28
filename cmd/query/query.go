package main

import (
	"fmt"
	"projects/xinhuatool/pkg/login"
	"projects/xinhuatool/pkg/office"
	_ "projects/xinhuatool/pkg/office"
	"projects/xinhuatool/pkg/page"
	"projects/xinhuatool/util"
	"strconv"
	"sync"
	"time"
)

type Task struct {
	id      string
	keyword string
	num     int
}

type TaskSet []Task

func (ts *TaskSet) load(p string) {
	sli, err := office.ReadExcel(p, "")
	if err != nil {
		panic("读取关键词表失败.")
	}

	for _, item := range sli[1:] {
		var tk Task
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
		*ts = append(*ts, tk)
	}
}

func (ts *TaskSet) toChan() chan Task {
	l := len(*ts)
	c := make(chan Task, l)
	defer close(c)

	for _, item := range *ts {
		c <- item
	}

	return c
}

// 保存的结果字段数目
var itemNum = 7

// 单页面显示结果总是显示50条
var NumPerPage = 50

type QueryParam struct {
	participle   string
	pg           int
	tp           int
	showDataFlag int
}

type Worker struct {
	sess       *login.Session
	queryParam *QueryParam
	taskChan   chan Task
	resultChan chan []string
	wg         *sync.WaitGroup
}

func NewWorker(queryParam *QueryParam, taskChan chan Task, resultChan chan []string, wg *sync.WaitGroup) *Worker {
	ua := util.UA.Get()
	header := map[string]string{"User-Agent": ua}
	sess := login.NewSess(header, "", 20)
	sess.LoginGuest()
	page.HistoryFrame(sess)
	worker := &Worker{
		sess:       sess,
		queryParam: queryParam,
		taskChan:   taskChan,
		resultChan: resultChan,
		wg:         wg,
	}
	return worker
}

func (worker *Worker) Run() {
	for {
		task, ok := <-worker.taskChan
		if !ok {
			break
		}
		worker.singleTask(task)
	}

	worker.wg.Done()
}

func (worker *Worker) singleTask(task Task) error {
	kw := task.keyword
	queryNum := task.num
	worker.queryParam.pg = 1 // 每个任务都是从第一个页面开始
	var al *page.ArticleList
	var err error
	for {
		al, err = worker.singleQuery(kw, worker.queryParam)
		if err != nil {
			fmt.Println("重新登录Guest用户中...")
			worker.sess.LoginGuest()
			continue
		}
		break
	}

	articleTotal := al.ArticleTotal

	// 检索结果为0
	if articleTotal == 0 {
		item := make([]string, itemNum)
		item[0] = kw
		worker.resultChan <- item
		return nil
	}

	// 检索结果不为0
	metas := al.MetaSli()
	for _, meta := range metas {
		item := []string{kw}
		item = append(item, meta...)
		worker.resultChan <- item
	}

	// 后续页面
	pageNum := util.GetPageNum(articleTotal, queryNum, NumPerPage)
	if pageNum <= 1 {
		return nil
	}
	for i := 2; i <= pageNum; i++ {
		worker.queryParam.pg = i
		var al *page.ArticleList
		var err error
		for {
			al, err = worker.singleQuery(kw, worker.queryParam)
			if err != nil {
				fmt.Println("重新登录Guest用户中...")
				worker.sess.LoginGuest()
				continue
			}
			break
		}

		metas := al.MetaSli()
		for _, meta := range metas {
			item := []string{kw}
			item = append(item, meta...)
			worker.resultChan <- item
		}

	}

	return nil
}

func (worker *Worker) singleQuery(kw string, queryParam *QueryParam) (*page.ArticleList, error) {
	participle := queryParam.participle
	pg := queryParam.pg
	tp := queryParam.tp
	showDataFlag := queryParam.showDataFlag

	resp, err := page.SearchKeyword(worker.sess, kw, participle, pg, tp, showDataFlag)
	if err != nil {
		return nil, err
	}

	al := page.NewArticleListFromResp(resp)
	return al, nil
}

func printProc(taskChan chan Task, taskTot int) {
	tot := strconv.Itoa(taskTot)
	for {
		l := strconv.Itoa(len(taskChan))
		fmt.Println("完成进度：", l+"/"+tot)
		time.Sleep(time.Second * 3)
		if l == "0" {
			break
		}
	}
}

func saveSlice(p string, res [][]string) error {
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

func saveResult(p string, resultChan chan []string) {
	var temp [][]string
	for {
		item, ok := <-resultChan
		if ok {
			temp = append(temp, item)
		} else {
			break
		}
	}
	saveSlice(p, temp)
}

// func worker(tasks *TaskSet, result *[][]string) {
// 	ua := util.UA.Get()
// 	header := map[string]string{"User-Agent": ua}
// 	sess := login.NewSess(header, "", 20)
// 	sess.LoginGuest()

// 	page.HistoryFrame(sess)

// 	// 进行检索
// 	participle := "on"
// 	pg := 1
// 	tp := page.IMG
// 	showDataFlag := page.SHOWALL

// 	itemLen := 7
// 	taskNum := len(*tasks)
// 	var tempResult [][]string
// 	for i, tk := range *tasks {
// 		fmt.Println("进行任务：", strconv.Itoa(i+1)+"/"+strconv.Itoa(taskNum))
// 		kw := tk.keyword
// 		resp, err := page.SearchKeyword(sess, kw, participle, pg, tp, showDataFlag)
// 		if err != nil {
// 			item := make([]string, itemLen)
// 			item[0] = kw
// 			tempResult = append(tempResult, item)
// 			fmt.Println(err)
// 			continue
// 		}

// 		al := page.NewArticleListFromResp(resp)
// 		metas := al.MetaSli()
// 		for _, meta := range metas {
// 			item := []string{kw}
// 			item = append(item, meta...)
// 			tempResult = append(tempResult, item)
// 		}
// 	}

// 	*result = tempResult
// }
