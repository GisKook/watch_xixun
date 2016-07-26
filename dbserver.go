package watch_xixun

import (
	"database/sql"
	"fmt"
	"github.com/giskook/watch_xixun/protocol"
	_ "github.com/jbarham/gopgsqldriver"
	"log"
	"strconv"
	"strings"
)

type DbServer struct {
	Db         *sql.DB
	Config     *DatabaseConfiguration
	SqlChan    chan string
	ActionChan chan string
}

func NewDbServer(config *DatabaseConfiguration) (*DbServer, error) {
	connstring := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", config.User, config.Passwd, config.Host, config.Port, config.Dbname)
	db, err := sql.Open("postgres", connstring)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &DbServer{
		Db:         db,
		Config:     config,
		SqlChan:    make(chan string),
		ActionChan: make(chan string),
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

func (db *DbServer) Action(action_string string) {
	value := strings.Split(action_string, ",")
	imei, _ := strconv.ParseUint(value[0], 10, 64)
	action_pkg := protocol.Parse_Action("181437EQ>;EPAXEM", value[0], "0001", value[1])
	c := NewConns().GetConn(imei)
	if c != nil {
		c.SendToClient(action_pkg)
	}
}

func (db *DbServer) Do() {
	for {
		select {
		case sql := <-db.SqlChan:
			db.Exec(sql)
			log.Println("insert to db")
		case action := <-db.ActionChan:
			db.Action(action)
		}

	}
}
