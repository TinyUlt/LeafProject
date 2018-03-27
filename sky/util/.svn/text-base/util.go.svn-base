package util

import (
    "os"
    "flag"
    "strings"
    "io/ioutil"
    "path/filepath"
    "runtime/debug"
    "sky/log"
    "sky/jsonhelper"
)

//defaults
var (
    prefix = flag.String("prefix", ".", "" )
    name   = flag.String("name", selfName(), "")
    conf   = flag.String("conf", "conf/service.conf", "")
)

var fd   *os.File

func Initial() jsonhelper.JSONObject {
    flag.Parse()
    os.Chdir(*prefix)

    b, err := ioutil.ReadFile(*conf)
    if err != nil {
        panic(err)
    }

    node := jsonhelper.NewJSONObjectFromBuf(b)
    if node.Len() == 0 {
        b = decrypt([]byte("application/x-www-form-urlencode"), b)
        node = jsonhelper.NewJSONObjectFromBuf(b)
    }

    node = node.GetAsObject(*name)
    level := node.GetAsObject("u").GetAsString("log_level")  
    if path := node.GetAsObject("u").GetAsString("log_file"); path != "" {
        if path == "auto" {
            path = "logs/" + *name + ".log" 
        }        
        if fd, _ = os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660); fd == nil { 
            panic("can't open file:" + path)
        } 
        log.NewGlobal(fd, level, 19)
    } else {
        log.SetGlobalLevel(level)
    }

    return node
}

func Final() {
   if e := recover(); e != nil {
        log.Critical("panic, recovered: '%v'", e) 
        log.Critical("\n\n--------------Stack-----------------\n\n")
        log.Critical(debug.Stack())
        log.Critical("\n\n--------------Stack-----------------\n\n")
    }   
    if fd != nil {       
        fd.Close()
        fd = nil
    }
}

func Recover(cb ...func()) {
    if e := recover(); e != nil {
        log.Critical("panic, recovered: '%v'", e) 
        log.Critical("\n\n--------------Stack-----------------\n\n")
        log.Critical(debug.Stack())
        log.Critical("\n\n--------------Stack-----------------\n\n")

        for _, f := range cb {
            f()
        }
    }   
}

func selfName() string {
    return strings.TrimRight(filepath.Base(os.Args[0]), filepath.Ext(os.Args[0]))  
}

