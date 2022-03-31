package office

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/nfnt/resize"
)

// ReadExcel function reads the excel file and return 2 dimension slice
func ReadExcel(path string, sheetName string) ([][]string, error) {
	if sheetName == "" {
		sheetName = "Sheet1"
	}
	var r [][]string
	fd, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Println(err.Error())
		return r, err
	}
	rows := fd.GetRows(sheetName)
	r = rows
	return r, err
}

// ReadExcel function reads the excel file and return 2 dimension slice
func ReadExcelFromReader(reader io.Reader, sheetName string) ([][]string, error) {
	if sheetName == "" {
		sheetName = "Sheet1"
	}
	var r [][]string
	fd, err := excelize.OpenReader(reader)
	if err != nil {
		fmt.Println(err.Error())
		return r, err
	}
	rows := fd.GetRows(sheetName)
	r = rows
	return r, err
}

// writeData writes the data to an *excel.File type variable
func writeData(dataIn interface{}, col []string) (*excelize.File, error) {
	f := excelize.NewFile()
	sheetName := "Sheet1"
	f.NewSheet(sheetName)

	// Get data
	var data [][]interface{}
	v := reflect.ValueOf(dataIn)
	if v.Kind() != reflect.Slice {
		err := errors.New("Data to write should be 2 dim slice.")
		return f, err
	}
	for i := 0; i < v.Len(); i++ {
		itemV := v.Index(i)
		if itemV.Kind() != reflect.Slice {
			err := errors.New("Data to write should be 2 dim slice.")
			return f, err
		}
		var innerSli []interface{}
		for j := 0; j < itemV.Len(); j++ {
			innerSli = append(innerSli, itemV.Index(j))
		}
		data = append(data, innerSli)
	}

	// Write data to excelize.File
	f.SetSheetRow(sheetName, "A1", &col)
	for i, item := range data {
		f.SetSheetRow(sheetName, "A"+strconv.Itoa(i+2), &item)
	}

	return f, nil
}

// WriteExcel writes the data to an xlsx file
func WriteExcel(p string, dataIn interface{}, col []string) error {
	f, err := writeData(dataIn, col)
	if err != nil {
		return err
	}

	// Write to file
	err = f.SaveAs(p)
	if err != nil {
		return err
	}

	return nil
}

// WriteExcelToWriter writes the data to an io.Writer
func WriteExcelToWriter(w io.Writer, dataIn interface{}, col []string) error {
	f, err := writeData(dataIn, col)
	if err != nil {
		return err
	}

	// Write to an io.Writer variable
	err = f.Write(w)
	if err != nil {
		return err
	}

	return nil
}

func AddImgToExcel(xlsx *excelize.File, sheet, cell string, width, height float64, imgPath string) error {
	f, err := os.Open(imgPath)
	if err != nil {
		fmt.Println("os.Open err:", err)
		return err
	}
	defer f.Close()

	var buffer bytes.Buffer
	buffer.ReadFrom(f)
	var img image.Image
	ext := strings.ToLower(filepath.Ext(imgPath))
	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(bytes.NewReader(buffer.Bytes()))
	case ".png":
		img, err = png.Decode(bytes.NewReader(buffer.Bytes()))
	default:
		err = errors.New("仅支持jpg,jpeg,png格式图片")
	}
	if err != nil {
		return err
	}

	var m image.Image
	if img.Bounds().Dx() > img.Bounds().Dy() {
		m = resize.Resize(120, 0, img, resize.Lanczos3)
	} else {
		m = resize.Resize(0, 120, img, resize.Lanczos3)
	}

	// write new image to file
	encodebuffer := bytes.NewBuffer(nil)
	jpeg.Encode(encodebuffer, m, nil)

	// 原代码存在问题，单元格不会随之改变，稍作调整
	// xlsx.SetColWidth("Sheet1", "E", "E", 30)
	rowNumStr := cell[1:]
	rowNum, err := strconv.Atoi(rowNumStr)
	if err != nil {
		return err
	}
	xlsx.SetRowHeight("Sheet1", rowNum, 90)

	// format := `{"lock_aspect_ratio": true, "locked": true, "positioning": "oneCell"}`
	format := `{"x_scale": 0.95, "y_scale": 0.95, "lock_aspect_ratio": true, "locked": true, "positioning": "absolute"}`
	// err = xlsx.AddPicture(sheet, location, outname, `{"lock_aspect_ratio": true, "locked": true, "positioning": "absolute"}`)//oneCell
	err = xlsx.AddPictureFromBytes(sheet, cell, format, "xx", ".jpg", encodebuffer.Bytes())
	if err != nil {
		return err
	}

	return nil
}
