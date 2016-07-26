package watch_xixun

import (
	"fmt"
	"log"
	"time"
)

type WatchData struct {
	dbserver *DbServer
}

func NewWatchData(dbserver *DbServer) *WatchData {
	return &WatchData{
		dbserver: dbserver,
	}
}
func (w *WatchData) DoRead() string {
	var imme, action string
	err := w.dbserver.Db.QueryRow("select imme,action from t_watchaction").Scan(&imme, &action)
	if err == nil {

		delsql := fmt.Sprintf("delete from t_watchaction where imme='%s'", imme)
		_, err := w.dbserver.Db.Exec(delsql)
		if err != nil {
			log.Println(err)
			return ""
		}
		imme += ","
		imme += action
		return imme
	}
	return ""
}
func (w *WatchData) Do() {
	t1 := time.NewTimer(time.Second * 1)

	for {
		select {

		case <-t1.C:
			//println("5s timer")
			msg := w.DoRead()
			if msg != "" {
				w.dbserver.ActionChan <- msg
			}
			t1.Reset(time.Second * 5)

		}
	}
}
