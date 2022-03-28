package main

import (
	"fmt"
	"projects/xinhuatool/pkg/media"
	"projects/xinhuatool/util"
	"sync"
)

func main() {
	// 读取数据
	path := "./图片id列表.xlsx"
	ids := readID(path)
	fmt.Println(ids)

	// id通道
	idNum := len(ids)
	idChan := make(chan string, idNum)
	for _, v := range ids {
		idChan <- v
	}
	close(idChan)

	// 结果通道
	resultChan := make(chan media.Metas, idNum)

	// 图片保存目录
	timePrefix := util.TimePrefix()
	imgDir := "./data/" + timePrefix + "image"
	util.ExistOrCreate(imgDir)

	// 下载器数量
	downloaderNum := util.DownloaderNum
	fmt.Println(downloaderNum)
	var wg sync.WaitGroup
	wg.Add(downloaderNum)

	// 创建线程干活
	for i := 0; i < downloaderNum; i++ {
		downloader := NewDownloader(idChan, resultChan, imgDir, &wg)
		go downloader.Run()
	}

	// 打印进度信息
	go printProc(idChan, idNum)

	wg.Wait()

	close(resultChan)

	// 保存结果
	saveResult("./data/"+timePrefix+"图片下载信息表.xlsx", resultChan)

	fmt.Println("Complete!")
}
