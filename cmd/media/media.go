package main

import (
	"fmt"
	"projects/xinhuatool/pkg/login"
	"projects/xinhuatool/pkg/media"
	"projects/xinhuatool/pkg/office"
	"projects/xinhuatool/pkg/page"
	"projects/xinhuatool/util"
	"strconv"
	"sync"
	"time"
)

func readID(p string) []string {
	sheet, err := office.ReadExcel(p, "")
	if err != nil {
		panic("读取图片id列表失败.")
	}

	var ids []string
	for _, v := range sheet[1:] {
		ids = append(ids, v[1])
	}
	return ids
}

type Downloader struct {
	sess       *login.Session
	idChan     chan string
	resultChan chan media.Metas
	imgDir     string
	wg         *sync.WaitGroup
}

func NewDownloader(idChan chan string, resultChan chan media.Metas, imgDir string, wg *sync.WaitGroup) *Downloader {
	ua := util.UA.Get()
	header := map[string]string{"User-Agent": ua}
	sess := login.NewSess(header, "", 20)
	sess.LoginGuest()
	page.HistoryFrame(sess)
	worker := &Downloader{
		sess:       sess,
		idChan:     idChan,
		resultChan: resultChan,
		imgDir:     imgDir,
		wg:         wg,
	}
	return worker
}

func (dl *Downloader) Run() {
	for {
		id, ok := <-dl.idChan
		if !ok {
			break
		}
		success := dl.singleRun(id)
		if !success {
			dl.sess.LoginGuest()
		}

		time.Sleep(time.Millisecond * time.Duration(util.RandMilliSec()))
	}

	dl.wg.Done()
}

func (dl *Downloader) singleRun(id string) bool {
	metas, err := media.DownloadIMG(dl.sess, id, dl.imgDir)
	if err != nil {
		// fmt.Println(err)
		return false
	}

	if len(metas.ImgURL) == 0 {
		return false
	}

	dl.resultChan <- metas
	return true
}

func printProc(idChan chan string, idTot int) {
	tot := strconv.Itoa(idTot)
	for {
		l := strconv.Itoa(idTot - len(idChan))
		fmt.Println("完成进度：", l+"/"+tot)
		time.Sleep(time.Second * 3)
		if l == tot {
			break
		}
	}
}

func saveSlice(p string, res []media.Metas) error {
	keys := []string{
		"ID",
		"ImgURL",
		"GdsDesc",
		"Tag",
		"ImgOutPath",
	}

	var items [][]string
	for _, meta := range res {
		var temp []string
		temp = append(temp, meta.ID, meta.ImgURL, meta.GdsDesc, meta.Tag, meta.ImgOutPath)
		items = append(items, temp)
	}

	return office.WriteExcel(p, items, keys)
}

func saveResult(p string, resultChan chan media.Metas) {
	var temp []media.Metas
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
