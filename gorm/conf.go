package gorm

type MysqlConf struct {
	Dialect  string
	User     string
	Password string
	Host     string
	Port     int
	Database string
	Charset  string
	ShowSql  bool
	MaxOpen  int
	MaxIdle  int
}
