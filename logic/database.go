package logic

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func initConnection(
	host,
	port,
	login,
	password,
	dbName string,
) (*sql.DB, error) {
	db, err := sql.Open(
		"mysql",
		fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s",
			login,
			password,
			host,
			port,
			dbName,
		),
	)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
