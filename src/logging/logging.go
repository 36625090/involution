package logging

import (
	"fmt"
	"github.com/36625090/involution/option"
	"github.com/36625090/involution/utils"
	"github.com/hashicorp/go-hclog"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type flush struct {

}

func (f flush) Flush() error {
	panic("implement me")
}

type _logging struct {
	app string
	option option.Log
	files []*os.File
	opts *hclog.LoggerOptions
	logger hclog.InterceptLogger
	sigChan chan struct{}
}
func (l _logging) Flush() error {
	l.close()
	l.rename()
	writer, err := openWriter(l.app, l.option)
	l.opts.Output = writer
	return err
}

func (l _logging) start()  {
	resettable := l.logger.(hclog.OutputResettable)

	timer := time.NewTimer( left() )
	for{
		select {
		case <-timer.C:
			log.Println("resettable")
			resettable.ResetOutputWithFlush(l.opts, &l)
			timer.Reset( left()  )
		case <-l.sigChan:
			timer.Stop()
			l.close()
			return
		case <-time.NewTicker(time.Millisecond * 800).C:
			l.logger.Info("test rotation")
		}
	}
}

func (l _logging) close()  {
	for _, file := range l.files {
		file.Close()
	}
}

func (l _logging) rename()  {
	for _, file := range l.files {
		rename(file)
	}
}

var logging _logging


//NewLogger 日志文件初始化方法，如有需要请自己实现日志轮转
func NewLogger(app string, option option.Log) (hclog.InterceptLogger, error) {
	os.Mkdir(option.Path, os.ModePerm)

	leveledWriter, err := openWriter(app, option)
	if err != nil {
		return nil, err
	}
	fmt.Println(left())
	opts := &hclog.LoggerOptions{
		IncludeLocation: true,
		Output:          leveledWriter,
		Level:           hclog.LevelFromString(option.Level),
		JSONFormat:      option.Format == "json",
	}

	logger := hclog.NewInterceptLogger(opts)
	{
		logging.app = app
		logging.option = option
		logging.opts = opts
		logging.logger = logger
		logging.sigChan = utils.MakeShutdownCh()
		go logging.start()
	}

	reset := logger.(hclog.OutputResettable)
	fmt.Println(reset)
	return logger, nil
}

func  openWriter(app string, option option.Log)(*hclog.LeveledWriter,error)  {

	standard, err := open(option.Path, app, hclog.NoLevel)
	if err != nil {
		return nil, err
	}
	trace, err := open(option.Path, app, hclog.Trace)
	if err != nil {
		return nil, err
	}
	logging.files = []*os.File{standard.(*os.File), trace.(*os.File)}

	if option.Console {
		trace = io.MultiWriter(trace, os.Stdout)
		standard = io.MultiWriter(standard, os.Stdout)
	}

	leveledWriter := hclog.NewLeveledWriter(standard, map[hclog.Level]io.Writer{
		hclog.Trace:   trace,
		hclog.NoLevel: standard,
	})



	return leveledWriter, nil
}

func open(path, app string, level hclog.Level) (io.Writer, error) {
	var name string
	if level == hclog.NoLevel {
		name = strings.Join([]string{path, app + ".log"}, "/")
	} else {
		name = strings.Join([]string{path, app + "_" + level.String() + ".log"}, "/")
	}

	return os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_RDWR, os.ModePerm)

}

func rename(file *os.File)error{
	//time.RFC3339 2006-01-02T15:04:05Z07:00
	lastDay := time.Now().Add(-time.Hour).Format("2006-01-02_15-04")
	newName := strings.Replace(file.Name(), ".log", "_"+lastDay+".log",4)
	return os.Rename(file.Name(), newName)
}

func  left() time.Duration {
	var inv int64 = 600
	return time.Duration(inv - time.Now().Unix() % inv) * time.Second
}