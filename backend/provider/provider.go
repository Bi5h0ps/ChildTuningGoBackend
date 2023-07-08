package provider

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"strings"
)

var HttpClientProvider HttpClient
var DatabaseEngine *gorm.DB

//const defaultDB = "root:Nmdhj2e2d@tcp(127.0.0.1:3306)/childTuningDB?parseTime=true"

const defaultDB = "root:Password2023!@tcp(127.0.0.1:3306)/childTuningDB?parseTime=true"

func init() {
	HttpClientProvider = NewHttpClient()
	var err error
	DatabaseEngine, err = gorm.Open(mysql.Open(defaultDB), &gorm.Config{})
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func FixTimeFormat(input string) string {
	s := input[:len(input)-1]
	return strings.ReplaceAll(s, "T", " ")
}
