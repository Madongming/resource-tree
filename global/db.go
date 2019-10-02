package global

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	log "github.com/cihub/seelog"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type DBLog struct{}

func (dblog *DBLog) Print(v ...interface{}) {
	log.Trace(v)
}

var db *gorm.DB

func initDb() {
	if db != nil {
		return
	}

	dataSourceName := fmt.Sprintf(
		"%s:%s@%s(%s)/%s?%s",
		Configs.Mysql.User,
		Configs.Mysql.Pass,
		Configs.Mysql.Protocol,
		Configs.Mysql.Address,
		Configs.Mysql.Db,
		Configs.Mysql.Params,
	)

	if err := retry("init db", 10, func() (e error) {
		db, e = gorm.Open("mysql", dataSourceName)
		return e
	}); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	if os.Getenv("GIN_MODE") == "debug" || os.Getenv("GIN_MODE") == "test" {
		db.LogMode(true)
		db.SetLogger(&DBLog{})
	}

	db.LogMode(true)
	db.SetLogger(&DBLog{})
	//db.SingularTable(DBSingularTable)
	db.DB().SetMaxIdleConns(int(Configs.Mysql.MaxOpenConns))
	db.DB().SetMaxOpenConns(int(Configs.Mysql.MaxIdleConns))
}

func DB() *gorm.DB {
	return db.New()
}

func Begin() *gorm.DB {
	return db.Begin()
}

func retry(desc string, retryTimes int, method func() error) (err error) {
	for i := 1; i < retryTimes+1; i++ {
		err = method()
		if err == nil {
			return
		} else if err == sql.ErrNoRows {
			log.Warnf("%s no rows found: %s", desc, err)
			return
		}
		log.Warnf("%s failed: %s, will retry %d times", desc, err, i)
		time.Sleep(time.Duration(i*i) * time.Second)
	}
	log.Warnf("%s failed after retry %d times, last error: %s", desc, retryTimes, err)
	return err
}
