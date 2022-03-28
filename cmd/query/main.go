package main

import (
	"flag"
	"fmt"
	_ "projects/xinhuatool/pkg/office"
	"projects/xinhuatool/pkg/page"
	"projects/xinhuatool/util"
	"sync"
)

func main() {
	// 读取数据
	path := "./关键词表.xlsx"
	var tasks TaskSet
	tasks.load(path)
	fmt.Println("读取关键词表成功.")

	// 检索任务数
	taskChan := tasks.toChan()
	taskTot := len(taskChan)
	fmt.Println("检索任务总数：", taskTot)

	// 检索参数
	pcFlag := flag.String("participle", "on", "Help: on means accurate; off mean fuzzy.")
	tpFlag := flag.Int("type", 2, "Help: 2 is image; 3 is audio; 4 is vedio.")
	flag.Parse()
	queryParam := &QueryParam{
		participle:   *pcFlag,
		pg:           1,
		tp:           *tpFlag,
		showDataFlag: page.SHOWALL,
	}

	// 结果chan
	resultChan := make(chan []string, 1000000)

	// worker数量
	workerNum := util.WorkerNum
	var wg sync.WaitGroup
	wg.Add(workerNum)

	// 创建worker干活
	for i := 0; i < workerNum; i++ {
		worker := NewWorker(queryParam, taskChan, resultChan, &wg)
		go worker.Run()
	}

	// 打印进度信息
	go printProc(taskChan, taskTot)

	wg.Wait()

	close(resultChan)

	// 保存结果
	saveResult("./data/检索结果.xlsx", resultChan)

	fmt.Println("Complete!")

	// test()

}
