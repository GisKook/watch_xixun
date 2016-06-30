package watch_xixun

import (
	"database/sql"
	"fmt"
	"log"
)

type DbServer struct {
	Db      *sql.DB
	Config  *DatabaseConfiguration
	SqlChan chan string
}

func NewExecDatabase(config *DatabaseConfiguration) (*DbServer, error) {
	connstring := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", config.User, config.Passwd, config.Host, config.Port, config.Dbname)
	db, err := sql.Open("postgres", connstring)
	if err != nil {
		return nil, err
	}

	return &DbServer{
		Db:      db,
		Config:  config,
		SqlChan: make(chan string),
	}, nil
}

func (db *DbServer) Stop() {
	db.Db.Close()
}

func (db *DbServer) Exec(sql string) error {
	_, err := db.Db.Exec(sql)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (db *DbServer) Insert(sql string) {
	db.SqlChan <- sql
}

func (db *DbServer) Do() {
	for {
		select {
		case sql := <-db.SqlChan:
			db.Exec(sql)
			log.Println("insert to db")
		}
	}
}
