package util

import (
	"math"
	"math/rand"
	"os"
	"time"
)

// 复制map
func CopyMap(m map[string]string) map[string]string {
	newMap := make(map[string]string)
	for k, v := range m {
		newMap[k] = v
	}
	return newMap
}

// map多加一层
func WrapMap(m map[string]string) map[string][]string {
	nm := make(map[string][]string)
	for k, v := range m {
		nm[k] = []string{v}
	}
	return nm
}

// 两个数值中最小的
func MinInTwo(i, j int) int {
	if i <= j {
		return i
	}
	return j
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

func SaveBytes(p string, data []byte) error {
	fd, err := os.OpenFile(p, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer fd.Close()

	_, err = fd.Write(data)
	return err
}

func RandMilliSec() int {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(1000)
	return i
}

func TimePrefix() string {
	now := time.Now()
	str := now.Format("2006.01.02.15.04.05")
	return str
}

func ExistOrCreate(path string) {
	info, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		os.Mkdir(path, 0777)
		return
	}

	if !info.IsDir() {
		os.Mkdir(path, 0777)
		return
	}
}
