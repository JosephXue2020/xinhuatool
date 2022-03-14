package util

import (
	"log"
	"time"

	"github.com/go-ini/ini"
)

var Cfg *ini.File

var Datetime *time.Time
var Account string
var Passwd string

func init() {
	var err error
	Cfg, err = ini.Load("./config/config.ini")
	if err != nil {
		log.Fatal(err)
	}

	LoadDatetime()

	LoadAuth()
}

func LoadDatetime() {
	sec, err := Cfg.GetSection("datetime")
	if err != nil {
		log.Println("不存在datetime章节，将取值当前日期")
		now := time.Now()
		Datetime = &now
		return
	}

	v, err := sec.GetKey("EndDay")
	if err != nil {
		log.Println("不存在EndDay字段，将取值当前日期")
		now := time.Now()
		Datetime = &now
		return
	}

	s := v.String()
	tm, err := time.Parse("2006-01-02", s)
	if err != nil {
		log.Println("EndDay字段解析失败，将取值当前日期")
		now := time.Now()
		Datetime = &now
		return
	}

	Datetime = &tm
}

func SetDatetime(tm time.Time) {
	sec, err := Cfg.GetSection("datetime")
	if err != nil {
		sec, _ = Cfg.NewSection("datetime")
	}

	s := tm.Format("2006-01-02")

	k, err := sec.GetKey("EndDay")
	if err != nil {
		k, _ = sec.NewKey("EndDay", s)
	}

	k.SetValue(s)

	// 要保存
	Cfg.SaveTo("./config/config.ini")
}

func LoadAuth() {
	sec, err := Cfg.GetSection("authority")
	if err != nil {
		log.Fatal(err)
	}

	v, err := sec.GetKey("account")
	if err != nil {
		log.Fatal(err)
	}
	Account = v.String()

	v, err = sec.GetKey("password")
	if err != nil {
		log.Fatal(err)
	}
	Passwd = v.String()
}
