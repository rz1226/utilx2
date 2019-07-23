package inix

import (
	"github.com/c4pt0r/ini"
)

var conf = ini.NewConf("test.ini")

var (
	V1 = conf.String("section_1", "field1", "v1")
	V2 = conf.String("section_1", "field2", "df")
)

func init() {
	conf.Parse()
}
