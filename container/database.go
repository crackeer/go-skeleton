package container

import (
	"fmt"
	"strings"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	DriverMySQL  = "mysql"
	DriverSQLite = "sqlite"
)

var (
	readDB  map[string]*gorm.DB = make(map[string]*gorm.DB)
	writeDB map[string]*gorm.DB = make(map[string]*gorm.DB)
)

func initDatabase(list []DatabaseConfig) error {
	for _, conf := range list {
		writeConn, readConn, err := connectDB(conf)
		if err != nil {
			return err
		}
		readDB[strings.ToLower(conf.Name)] = readConn
		writeDB[strings.ToLower(conf.Name)] = writeConn
	}
	return nil
}

func connectDB(conf DatabaseConfig) (*gorm.DB, *gorm.DB, error) {
	var (
		readDB, writeDB *gorm.DB
		err             error
	)
	if conf.Driver == DriverSQLite {
		writeDB, err = connectSQLite(conf.File)
		if err != nil {
			return nil, nil, fmt.Errorf("connect sqlite write db `%s` error: %s", conf.File, err.Error())
		}
		return writeDB, writeDB, nil
	}

	if conf.Driver == DriverMySQL {
		writeDB, err = connectMySQL(conf.WriteUser, conf.WritePassword, conf.WriteHost, conf.Database, conf.Charset)
		if err != nil {
			return nil, nil, fmt.Errorf("connect mysql write db `%s` error: %s", conf.WriteHost, err.Error())
		}
		readDB, err = connectMySQL(conf.ReadUser, conf.ReadPassword, conf.ReadHost, conf.Database, conf.Charset)
		if err != nil {
			return nil, nil, fmt.Errorf("connect mysql read db `%s` error: %s", conf.ReadHost, err.Error())
		}
		return writeDB, readDB, nil
	}
	return nil, nil, fmt.Errorf("unsupported database driver: %s", conf.Driver)
}

// GetWriteDB ...
//
//	@param name
//	@return *gorm.DB
func GetWriteDB(name string) *gorm.DB {
	if db, exists := writeDB[name]; exists {
		return db
	}
	return nil
}

// GetReadDB ...
//
//	@param name
//	@return *gorm.DB
func GetReadDB(name string) *gorm.DB {
	if db, exists := readDB[name]; exists {
		return db
	}
	return nil
}

func connectMySQL(user, password, host, database, charset string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local", user, password, host, database, charset)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func connectSQLite(file string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(file), &gorm.Config{})
}
