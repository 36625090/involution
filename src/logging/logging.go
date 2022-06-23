package logging

import (
	"github.com/36625090/involution/option"
	"github.com/36625090/involution/utils"
	"github.com/hashicorp/go-hclog"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type rotationPolicy string

const (
	rotationPolicyDay  rotationPolicy = "day"
	rotationPolicyHour rotationPolicy = "hour"
)

type rotatedLogging struct {
	app     string
	option  option.Log
	files   []*os.File
	opts    *hclog.LoggerOptions
	logger  hclog.InterceptLogger
	sigChan chan struct{}
}

func (l *rotatedLogging) Flush() error {
	l.close()
	if err := l.rename(); err != nil {
		return err
	}
	writer, err := l.openWriter()
	l.opts.Output = writer
	return err
}

func (l *rotatedLogging) start() {
	resettable := l.logger.(hclog.OutputResettable)
	timer := time.NewTimer(l.nextRoundOfMilliDuration())
	defer timer.Stop()
	log.Println("next round", l.nextRoundOfMilliDuration(), "at", time.Now().Add(l.nextRoundOfMilliDuration()).Format(time.RFC3339))
	for {
		select {
		case <-timer.C:
			resettable.ResetOutputWithFlush(l.opts, l)
			log.Println("next round", l.nextRoundOfMilliDuration(), "at", time.Now().Add(l.nextRoundOfMilliDuration()).Format(time.RFC3339))
			timer.Reset(l.nextRoundOfMilliDuration())
		case <-l.sigChan:
			l.close()
			return
		}
	}

}

func (l *rotatedLogging) close() error {
	for _, file := range l.files {
		if err := file.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (l *rotatedLogging) rename() error {
	for _, file := range l.files {
		if err := l.rotateRename(file); err != nil {
			return err
		}
	}
	return nil
}

func (l *rotatedLogging) nextRoundOfMilliDuration() time.Duration {
	policy := rotationPolicy(l.option.Rotate)
	if policy == rotationPolicyHour {
		return nextHourOfMilliDuration()
	}
	return nextDayOfMilliDuration()
}

//NewLogger 日志文件初始化方法，如有需要请自己实现日志轮转
func NewLogger(app string, option option.Log) (hclog.InterceptLogger, error) {
	logging := &rotatedLogging{
		app:    app,
		option: option,
	}

	if err := os.Mkdir(option.Path, os.ModePerm); err != nil && !os.IsExist(err) {
		return nil, err
	}

	leveledWriter, err := logging.openWriter()
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

func (l *rotatedLogging) openWriter() (*hclog.LeveledWriter, error) {

	standard, err := open(l.option.Path, l.app, hclog.NoLevel)
	if err != nil {
		return nil, err
	}
	trace, err := open(l.option.Path, l.app, hclog.Trace)
	if err != nil {
		return nil, err
	}
	l.files = []*os.File{standard.(*os.File), trace.(*os.File)}

	if l.option.Console {
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

func (l rotatedLogging) rotateRename(file *os.File) error {
	//time.RFC3339 2006-01-02T15:04:05Z07:00
	var past string
	if rotationPolicy(l.option.Rotate) == rotationPolicyDay {
		past = time.Now().Add(-time.Minute * 5).Format("2006-01-02")
	} else {
		past = time.Now().Add(-time.Minute * 5).Format("2006-01-02-15")
	}

	newName := strings.Replace(file.Name(), ".log", "_"+past+".log", 4)
	return os.Rename(file.Name(), newName)
}

func nextDayOfMilliDuration() time.Duration {
	_, offset := time.Now().Local().Zone()
	intervalDayOfMilli := int64(3600000 * 24)
	return time.Duration(intervalDayOfMilli-time.Now().UnixMilli()%intervalDayOfMilli-int64(offset*1000)) * time.Millisecond
}

func nextHourOfMilliDuration() time.Duration {
	intervalHourOfMilli := int64(3600000)
	return time.Duration(intervalHourOfMilli-time.Now().UnixMilli()%intervalHourOfMilli) * time.Millisecond
}
