package media

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"projects/xinhuatool/pkg/login"
	"projects/xinhuatool/util"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// 图片frame页面基础url
var FrameURLBase = "http://home.xinhua-news.com//gdsdetailxhs/media/3/"

// 图片下载页面基础url
var ImgURLBase = "http://img.xinhua-news.com/archive/image"

type Metas struct {
	ID         string
	ImgURL     string
	GdsDesc    string
	Tag        string
	ImgOutPath string
}

// 解析图片frame页面
func parseFrame(resp *http.Response, res *Metas) error {
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	// 图片地址
	imgURL, ok := doc.Find("ul#links").Find("img").Attr("src")
	if ok {
		res.ImgURL = imgURL
	} else {
		res.ImgURL = ""
	}

	// gdsDesc
	// gdsDesc := doc.Find("ul#links").Find("gdsDesc his-detail-pic search-tag-02240201 list-detail-page").Text()
	gdsDesc := doc.Find("ul#links").Find(".gdsDesc").Text()
	res.GdsDesc = gdsDesc

	// 富标签
	tags := doc.Find("#pro-dt").Find(".key-wrap").Find("div").Text()
	tags = strings.TrimSpace(tags)
	tags = strings.ReplaceAll(tags, "\n", " ")
	pat := regexp.MustCompile(`[\s]{1,}`)
	tags = pat.ReplaceAllString(tags, ",")
	res.Tag = tags

	return nil
}

// 下载图像
func DownloadIMG(sess *login.Session, id string, imgDir string) (Metas, error) {
	info := Metas{}
	info.ID = id

	frameURL := FrameURLBase + id
	resp, err := util.GetWithHead(frameURL, sess.Header, nil, sess.Client)
	if err != nil {
		return info, err
	}
	defer resp.Body.Close()

	err = parseFrame(resp, &info)
	if err != nil {
		return info, err
	}

	imgURL := info.ImgURL
	if len(imgURL) == 0 {
		err = fmt.Errorf("没有解析到图片URL")
		return info, err
	}

	resp, err = util.GetWithHead(imgURL, sess.Header, nil, sess.Client)
	if err != nil {
		return info, err
	}
	defer resp.Body.Close()

	tpStr := resp.Header.Get("Content-Type")
	segs := strings.Split(tpStr, "/")
	suffix := segs[len(segs)-1]
	imgOutPath := filepath.Join(imgDir, id+"."+suffix)

	byteData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return info, err
	}

	err = util.SaveBytes(imgOutPath, byteData)
	if err != nil {
		return info, err
	}
	info.ImgOutPath = imgOutPath

	return info, nil
}
