package mgo

import (
	// "time"
	"fmt"
	"sky/jsonhelper"
	_log "sky/log"
)

func DialEx(node jsonhelper.JSONObject) (*Collection, error) {
	//_log.Debug(">>>>>>>>>> DialEx %v", node)
	var url string
	if len(node.GetAsString("user")) > 0 {
		url = fmt.Sprintf("%s:%s@%s", node.GetAsString("user"), node.GetAsString("pass"), node.GetAsString("addr"))
	} else {
		url = fmt.Sprintf("%s", node.GetAsString("addr"))
	}
	session, _err := Dial(url)
	//session.SetMode(Eventual, true)
	if session == nil || _err != nil {
		_log.Debug("数据库联接失败 %v %v ", session, _err)
		return nil, _err
	}
	return session.DB(node.GetAsString("db")).C(node.GetAsString("collection")), nil
}

func DialDb(node jsonhelper.JSONObject) (*Session, error) {
	//_log.Debug(">>>>>>>>>> DialEx %v", node)
	var url string
	if len(node.GetAsString("user")) > 0 {
		url = fmt.Sprintf("%s:%s@%s", node.GetAsString("user"), node.GetAsString("pass"), node.GetAsString("addr"))
	} else {
		url = fmt.Sprintf("%s", node.GetAsString("addr"))
	}
	session, _err := Dial(url)
	//session.SetMode(Eventual, true)
	if session == nil || _err != nil {
		_log.Debug("数据库联接失败 %v %v ", session, _err)

	}
	return session, _err
}

func Close(coll *Collection) {
	coll.Database.Session.Close()
}
