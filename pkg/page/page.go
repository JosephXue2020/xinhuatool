package page

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"projects/xinhuatool/pkg/login"
	"projects/xinhuatool/util"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// 历史资源库url
var HistoryURL = "http://home.xinhua-news.com/history/page"

// 历史资源库检索url
var HistoryQueryURL = "http://home.xinhua-news.com/history/page"

// 数据列表样例
var SampleData = map[string]string{
	"pageNumber":   "1",
	"pageSize":     "50",
	"multimedia":   "true",
	"pageType":     "list", // 代表返回列表
	"keyword":      "春节",
	"gdsTypeIds":   "type_2", // 1代表文本，2代表图片，3代表音频, 4代表视频...
	"languageId":   "",
	"showDataFlag": "1",  // 1代表所有，2代表付费用户的已订购
	"participle":   "on", // on代表精确，off代表非精确（模糊）
	"searchType":   "all",
}

// 检索的类型
const (
	TEXT = 1 + iota
	IMG
	AUDIO
	VIDEO
)

const (
	SHOWALL = 1 + iota
	SHOWPAYED
)

// 获取history主页面
func HistoryFrame(sess *login.Session) (*http.Response, error) {
	resp, err := util.GetWithHead(HistoryURL, sess.Header, nil, sess.Client)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// 发送检索请求
func SearchKeyword(sess *login.Session, kw, participle string, pg, tp, showDataFlag int) (*http.Response, error) {
	m := util.CopyMap(SampleData)
	m["keyword"] = kw
	m["gdsTypeIds"] = "type_" + strconv.Itoa(tp)
	m["showDataFlag"] = strconv.Itoa(showDataFlag)
	m["participle"] = participle
	data := util.WrapMap(m)

	resp, err := util.PostForm(HistoryQueryURL, sess.Header, data, sess.Client)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// 解析检索返回的列表
type Meta struct {
	title        string
	gdsid        string
	abstraction  string
	introduction string
	copyright    string
	tags         []string
}

type ArticleList struct {
	text           string
	metas          []Meta
	pageNum        int
	pageTotal      int
	articlePerPage int
	articleTotal   int

	normFlag bool // true代表是正常的列表页面
}

// 构造函数
func NewArticleList(text string) *ArticleList {
	al := &ArticleList{text: text}

	al.GetNum()
	al.GetMeta()

	return al
}

// 构造函数从resp
func NewArticleListFromResp(resp *http.Response) *ArticleList {
	byteData, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	text := string(byteData)

	return NewArticleList(text)
}

func (al *ArticleList) GetNum() {
	pat := regexp.MustCompile(`<div id="listPageControlbar" .*?></div>`)
	finds := pat.FindAllString(al.text, -1)
	if len(finds) != 1 {
		al.normFlag = false
	} else {
		al.normFlag = true
	}

	if al.normFlag {
		seg := finds[0]
		// pageNum
		pat := regexp.MustCompile(`data-page-number="(.*?)"`)
		pageNum := pat.FindStringSubmatch(seg)[1]
		al.pageNum, _ = strconv.Atoi(pageNum)

		// pageTotal
		pat = regexp.MustCompile(`data-total-page="(.*?)"`)
		pageTotal := pat.FindStringSubmatch(seg)[1]
		al.pageTotal, _ = strconv.Atoi(pageTotal)

		// articlePerPage
		pat = regexp.MustCompile(`data-page-size="(.*?)"`)
		articlePerPage := pat.FindStringSubmatch(seg)[1]
		al.articlePerPage, _ = strconv.Atoi(articlePerPage)

		// articleTotal
		pat = regexp.MustCompile(`data-total-row="(.*?)"`)
		articleTotal := pat.FindStringSubmatch(seg)[1]
		al.articleTotal, _ = strconv.Atoi(articleTotal)
	}

}

func (al *ArticleList) GetMeta() error {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(al.text))
	if err != nil {
		return nil
	}
	doc.Find("div.date-list-item-warp").Each(func(i int, s *goquery.Selection) {
		var meta Meta

		gdsid, _ := s.Find("article").Attr("gdsid")
		meta.gdsid = strings.TrimSpace(gdsid)

		title := s.Find(".his-list-preview").Find("h4.tit").Text()
		meta.title = strings.TrimSpace(title)

		abstr := s.Find(".item-com").Find("span").Text()
		meta.abstraction = strings.TrimSpace(abstr)

		cr := s.Find("div.copyright").Text()
		meta.copyright = strings.TrimSpace(cr)

		al.metas = append(al.metas, meta)
	})

	// fmt.Println("数量：", len(al.metas))
	return nil
}

func (al *ArticleList) Metas() []Meta {
	return al.metas
}

func (al *ArticleList) MetaSli() [][]string {
	var res [][]string
	for _, meta := range al.metas {
		var item []string
		v := reflect.ValueOf(meta)
		for i := 0; i < v.NumField()-1; i++ {
			item = append(item, v.Field(i).String())
		}
		tags := meta.tags
		if tags != nil {
			jsonData, _ := json.Marshal(tags)
			item = append(item, string(jsonData))
		}
		item = append(item, "")

		res = append(res, item)
	}
	return res
}
