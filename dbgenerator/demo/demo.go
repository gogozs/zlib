package main

import (
	"github.com/gogozs/zlib/dbgenerator"
	"github.com/gogozs/zlib/xlog"
)

func main() {
	config := &dbgenerator.ConverterConfig{
		DBConfig: dbgenerator.DBConfig{
			Host:     "127.0.0.1",
			Username: "root",
			Password: "test",
			Port:     3306,
			Dbname:   "admin",
		},
		ConvertConfig: dbgenerator.ConvertConfig{
			TemplateDir:   "templates",
			TablePrefix:   "t_",
			IncludeTables: nil,
			ExcludeTables: nil,
			PackagePath:   "",
		},
	}
	c, err := dbgenerator.NewDBConverter(config)
	if err != nil {
		xlog.Fatal(err.Error())
	}
	if err = c.Convert(); err != nil {
		xlog.Fatal(err.Error())
	}
}
