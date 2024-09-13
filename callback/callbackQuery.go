package callback

import (
	"time"
)

type CallbackInfo struct {
	url  string
	data any
	cb   func(any)
}

var query chan *CallbackInfo

var stop chan any

type LogerFun func(string, ...any)

type Logger struct {
	Info  LogerFun
	Error LogerFun
}

var logger Logger

func AddCallback(url string, data any, cb func(any)) {
	go func() {
		query <- &CallbackInfo{
			url:  url,
			data: data,
			cb:   cb,
		}
	}()
}

func StartConsume(infoLog LogerFun, errlog LogerFun) {

	logger = Logger{
		Info:  infoLog,
		Error: errlog,
	}
	query = make(chan *CallbackInfo)
	stop = make(chan any)
	go func() {
		for {
			select {
			case <-stop:
				close(query)
				return
			case data := <-query:
				logger.Info("begin dispose %+v", data)
				err := Instance().Callback(data.url, data.data)
				if err != nil {
					logger.Error("dispose %+v error  %s", data, err)
					// 创建一个定时器，3秒后向retryCh通道发送当前时间
					timer := time.After(5 * time.Second)
					<-timer
					logger.Error("retry %+v", data)
					AddCallback(data.url, data.data, data.cb)
				} else {
					logger.Info("dispose %+v sucess", data)
					if data.cb != nil {
						logger.Info("dispose %+v sucess callback", data)
						data.cb(data.data)
					}
				}
			}
		}
	}()
}

func Stop() {
	close(stop)
}
