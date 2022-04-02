package main

import (
	"flag"
	"fmt"
	"projects/xinhuatool/pkg/office"
	"projects/xinhuatool/util"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func main() {
	mod := flag.Int("m", 0, "0，转移检索结果，给下载图片做准备；1，将下载好的结果粘贴进excel")
	flag.Parse()
	switch *mod {
	case MoveQueryResult:
		moveQueryResult()
	case InputImg:
		fmt.Println("将要执行插入图片操作")
		inputImg()
	default:
		panic("使用了模块不具备的功能")
	}
}

// 标记本glue程序不同的功能
const (
	MoveQueryResult = iota
	InputImg
)

func moveQueryResult() {
	src := "./data/检索结果.xlsx"
	items, err := office.ReadExcel(src, "Sheet1")
	if err != nil {
		panic("转移检索结果失败")
	}

	var gdsids [][]string
	for i, item := range items {
		gdsids = append(gdsids, []string{strconv.Itoa(i), item[2]})
	}
	gdsids[0][0] = "idx"

	outPath := "./图片id列表.xlsx"
	office.WriteExcel(outPath, gdsids[1:], gdsids[0])
}

func idx2ColMap() map[int]string {
	var letters []string
	var single []string
	for ch := 'A'; ch <= 'Z'; ch++ {
		single = append(single, string(ch))
	}

	letters = append(letters, single...)
	for _, ch1 := range single {
		for _, ch2 := range single {
			letters = append(letters, ch1+ch2)
		}
	}

	res := make(map[int]string, len(letters))
	for i, v := range letters {
		res[i] = string(v)
	}

	return res
}

// 从0开始
func idx2Col(i int) string {
	m := idx2ColMap()
	return m[i]
}

func inputImg() {
	src1 := "./data/检索结果.xlsx"
	items1, err := office.ReadExcel(src1, "Sheet1")
	if err != nil {
		panic("读取检索结果表格失败")
	}

	src2 := "./data/图片下载信息表.xlsx"
	items2, err := office.ReadExcel(src2, "Sheet1")
	if err != nil {
		panic("读取图片下载信息表失败")
	}

	if len(items1) != len(items2) {
		panic("两个表格不匹配")
	}
	l := len(items1)

	var itemTot [][]string
	for i := 0; i < l; i++ {
		var temp []string
		temp = append(items1[i], items2[i][1:]...)
		itemTot = append(itemTot, temp)
	}

	util.RemoveFileIfExist(src1)
	util.RemoveFileIfExist(src2)
	office.WriteExcel(src1, itemTot[1:], itemTot[0])

	// 再次读取
	items, err := office.ReadExcel(src1, "Sheet1")
	if err != nil {
		panic("结果表格失败")
	}

	l = len(items)
	colnum := len(items[0])
	colName := idx2Col(colnum)

	xlsx, err := excelize.OpenFile(src1)
	if err != nil {
		panic("读取图片下载信息表失败")
	}

	for i := 1; i < l; i++ {
		item := items[i]
		imgPath := item[colnum-1]
		if len(imgPath) == 0 {
			continue
		}
		imgPath = "./" + imgPath
		cell := colName + strconv.Itoa(i+1)
		office.AddImgToExcel(xlsx, "Sheet1", cell, 80, 80, imgPath)
	}

	err = xlsx.SaveAs(src1)
	if err != nil {
		panic("保存最终结果失败")
	}
}
