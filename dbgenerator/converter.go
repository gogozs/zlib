package dbgenerator

import (
	"os"
	"strings"
	"text/template"

	"github.com/gogozs/zlib/storage/xsql"
)

type DBConverter struct {
	config *ConverterConfig
	parser *DBParser
}

type ConverterConfig struct {
	DBConfig      DBConfig
	ConvertConfig ConvertConfig
}

type DBConfig struct {
	Host     string
	Username string
	Password string
	Port     int
	Dbname   string
}

type ConvertConfig struct {
	TemplateDir        string
	TablePrefix        string
	IncludeTables      []string
	ExcludeTables      []string
	ModelPackagePath   string
	DaoPackagePath     string
	ServicePackagePath string
	ServerPackagePath  string
}

func NewDBConverter(config *ConverterConfig) (*DBConverter, error) {
	db, err := xsql.NewLogDB(&xsql.SQLConfig{
		Host:     config.DBConfig.Host,
		Username: config.DBConfig.Username,
		Password: config.DBConfig.Password,
		Port:     config.DBConfig.Port,
		Dbname:   config.DBConfig.Dbname,
	})
	if err != nil {
		return nil, err
	}
	parser := NewDBParser(db, config.DBConfig.Dbname)
	return &DBConverter{
		config: config,
		parser: parser,
	}, nil
}

func (c *DBConverter) Convert() error {
	tables, err := c.getTargetTables()
	if err != nil {
		return err
	}
	for _, table := range tables {
		if err = c.convert(table); err != nil {
			return err
		}
	}

	return nil
}

func (c *DBConverter) convert(table string) error {
	tableData, err := c.parser.GetTableData(table, nil)
	if err != nil {
		return err
	}
	templates := template.New("root")
	if templates, err = templates.ParseGlob(c.config.ConvertConfig.TemplateDir + "/*"); err != nil {
		return err
	}
	for _, tmpl := range templates.Templates() {
		tmplName := tmpl.Name()
		dstFile := strings.ReplaceAll(tmplName, "[model]", strings.ToLower(tableData.TableName))
		dstFile = strings.ReplaceAll(dstFile, "[Model]", tableData.TableName)
		pkgpath := strings.ReplaceAll(c.config.ConvertConfig.ModelPackagePath, ".", "/")
		dstFile = strings.ReplaceAll(dstFile, "[pkgpath]", pkgpath)
		f, err := os.OpenFile(dstFile, os.O_WRONLY|os.O_CREATE, 0766)
		if err != nil {
			return err
		}
		if err := tmpl.ExecuteTemplate(f, tmplName, tableData); err != nil {
			return err
		}
	}
	return nil
}

func (c *DBConverter) getTargetTables() ([]string, error) {
	if len(c.config.ConvertConfig.IncludeTables) > 0 {
		return c.config.ConvertConfig.IncludeTables, nil
	}
	tables, err := c.parser.GetTables()
	if err != nil {
		return nil, err
	}
	targetTables := make([]string, 0, len(tables))
	m := make(map[string]interface{}, len(c.config.ConvertConfig.ExcludeTables))
	for _, t := range c.config.ConvertConfig.ExcludeTables {
		m[t] = struct{}{}
	}
	for _, table := range tables {
		_, ok := m[table]
		if !ok {
			targetTables = append(targetTables, table)
		}
	}
	return targetTables, nil
}
