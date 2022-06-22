package logging

import (
	"github.com/36625090/involution/option"
	"github.com/36625090/involution/utils"
	"github.com/hashicorp/go-hclog"
	"io"
	"os"
	"strings"
	"time"
)

const rotationInterval =  time.Hour * 24

type _logging struct {
	app     string
	option  option.Log
	files   []*os.File
	opts    *hclog.LoggerOptions
	logger  hclog.InterceptLogger
	sigChan chan struct{}
}

func (l _logging) Flush() error {
	l.close()
	l.rename()
	writer, err := openWriter(l.app, l.option)
	l.opts.Output = writer
	return err
}

func (l _logging) start() {
	resettable := l.logger.(hclog.OutputResettable)
	timer := time.NewTimer(nextDayLeft())
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			resettable.ResetOutputWithFlush(l.opts, &l)
			time.Sleep(time.Second * 10)
			timer.Reset(nextDayLeft())
		case <-l.sigChan:
			l.close()
			return
		}
	}

}

func (l _logging) close() {
	for _, file := range l.files {
		file.Close()
	}
}

func (l _logging) rename() {
	for _, file := range l.files {
		rename(file)
	}
}

var logging _logging

//NewLogger 日志文件初始化方法，如有需要请自己实现日志轮转
func NewLogger(app string, option option.Log) (hclog.InterceptLogger, error) {
	if err := os.Mkdir(option.Path, os.ModePerm); err != nil && !os.IsExist(err) {
		return nil, err
	}

	leveledWriter, err := openWriter(app, option)
	if err != nil {
		return nil, err
	}

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
	return logger, nil
}

func openWriter(app string, option option.Log) (*hclog.LeveledWriter, error) {

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

func rename(file *os.File) error {
	//time.RFC3339 2006-01-02T15:04:05Z07:00
	lastDay := time.Now().Add(-rotationInterval).Format("2006-01-02")
	newName := strings.Replace(file.Name(), ".log", "_"+lastDay+".log", 4)
	return os.Rename(file.Name(), newName)
}

func nextDayLeft() time.Duration {
	_, offset := time.Now().Zone()
	var inv = int64(rotationInterval / time.Millisecond)
	return time.Duration(inv - time.Now().UnixMilli() % inv - int64(offset * 1000)) * time.Millisecond
}
