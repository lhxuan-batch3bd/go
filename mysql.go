package mysql

import (
	"database/sql"
	"github.com/user/stories/config"
)
import _ "github.com/go-sql-driver/mysql"

func DbConnect() (db *sql.DB) {
	configDB := config.GetConfigDB()
	db, err := sql.Open(configDB.Driver, configDB.UserName+":"+configDB.Password+"@/"+configDB.DBName)
	if err != nil {
		panic(err.Error())
	}
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	return db
}
