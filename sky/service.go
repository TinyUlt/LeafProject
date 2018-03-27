package sky

import (
	"sync"
	"time"
	"syscall"
	"os"
	"os/signal"
	"runtime/debug"
	"sky/log"
	"sky/net"
	"sky/jsonhelper"
)

var (
	//version string
	MajorVersion 	= "0"
	MinorVersion 	= "0"
	RevisionVersion = "0"
	BuildVersion 	= "2"
	Version 		= MajorVersion+"."+MinorVersion+"."+RevisionVersion+"."+ BuildVersion
	BuildTime 		= "1979-01-01 00:00:00.000"
	)


type Delegate interface {
	OnStart(s *Service)
	OnStop(s *Service)
	OnTick(delta time.Duration)
}

type Service struct {
	delegate  Delegate
	done      sync.WaitGroup
	ticker    *time.Ticker
    ns        *net.Server
    conf 	  jsonhelper.JSONObject
}

func CreateService(d Delegate, conf jsonhelper.JSONObject) *Service {
	log.Trace("service version:%s", Version)
	log.Trace("build time:%s", BuildTime)
	s := &Service {delegate:d, conf:conf}
	s.ns = net.NewServer(conf)
    t := conf.GetAsObject("u").GetAsInt("tick")
	s.ticker = time.NewTicker(time.Duration(t)*time.Millisecond)
	s.register()
	return s
}

func (self *Service) Start() *sync.WaitGroup {	
	self.ns.Start()

    self.delegate.OnStart(self)
    
    go self.mux()
    
	self.done.Add(1)
	
	return &(self.done)
}

func (self *Service) Stop() { 
	self.delegate.OnStop(self)

	self.ns.Stop()

	self.done.Done()
}

func (self *Service) GetModule(key string) *net.Module {
    return self.ns.GetModule(key)
}

func (self *Service) mux() {
	defer func () {
	    if e := recover(); e != nil {
	        log.Critical("---------recover---------")
	        log.Critical("panic: %v", e)
	        log.Critical(debug.Stack())
	        log.Critical("---------recover---------")
	        self.mux()
	    }   
	}()

    var call *net.Call
    prev := time.Now()
    ch   := make(chan os.Signal)
    signal.Notify(ch, syscall.SIGINT, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGSEGV, syscall.SIGTERM)   
    for {
        select {
        case <-self.ticker.C:
            self.delegate.OnTick(time.Since(prev))
            prev = time.Now()
        case call = <-self.ns.Input:
           self.ns.Process(call)
        case s := <-ch:
            log.Trace("Signal:%q", s.(syscall.Signal))
            self.Stop()
            log.Trace("Stopd:%q", s.(syscall.Signal))
            return
        }
    }
}
