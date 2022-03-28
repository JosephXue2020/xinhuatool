package main

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	path := "../var/图片页面.html"
	info := make(map[string]string)

	parseFrame(path, &info)
}

// 解析图片frame页面
func parseFrame(p string, res *map[string]string) error {
	reader, _ := os.Open(p)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return err
	}

	// 图片地址
	imgURL, ok := doc.Find("ul#links").Find("img").Attr("src")
	if ok {
		(*res)["imgURL"] = imgURL
	}

	// gdsDesc
	// gdsDesc := doc.Find("ul#links").Find("gdsDesc his-detail-pic search-tag-02240201 list-detail-page").Text()
	gdsDesc := doc.Find("ul#links").Find(".gdsDesc").Text()
	(*res)["gdsDesc"] = gdsDesc

	// 富标签
	tags := doc.Find("#pro-dt").Find(".key-wrap").Find("div").Text()
	tags = strings.TrimSpace(tags)
	tags = strings.ReplaceAll(tags, "\n", " ")
	pat := regexp.MustCompile(`[\s]{1,}`)
	tags = pat.ReplaceAllString(tags, ",")
	(*res)["tags"] = tags

	fmt.Println(res)
	return nil
}

// 给定总数，需要数，单数，获取页面数
func GetPageNum(tot, queryNum, numPerPage int) int {
	if tot <= numPerPage {
		return 1
	}

	actualNum := tot
	if queryNum > 0 {
		actualNum = MinInTwo(tot, queryNum)
	}
	pageNum := int(math.Ceil(float64(actualNum) / float64(numPerPage)))
	return pageNum
}

// 两个数值中最小的
func MinInTwo(i, j int) int {
	if i <= j {
		return i
	}
	return j
}
