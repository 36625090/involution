package logging

import (
	"github.com/36625090/involution/option"
	"github.com/hashicorp/go-hclog"
	"io"
	"os"
	"strings"
)

//NewLogger 日志文件初始化方法，如有需要请自己实现日志轮转
func NewLogger(app string, option option.Log) (hclog.InterceptLogger, error) {
	os.Mkdir(option.Path, os.ModePerm)

	standard, err := open(option.Path, app, hclog.NoLevel)
	if err != nil {
		return nil, err
	}

	trace, err := open(option.Path, app, hclog.Trace)
	if err != nil {
		return nil, err
	}

	if option.Console {
		trace = io.MultiWriter(trace, os.Stdout)
		standard = io.MultiWriter(standard, os.Stdout)
	}

	leveledWriter := hclog.NewLeveledWriter(standard, map[hclog.Level]io.Writer{
		hclog.Trace:   trace,
		hclog.NoLevel: standard,
	})

	logger := hclog.NewInterceptLogger(&hclog.LoggerOptions{
		IncludeLocation: true,
		Output:          leveledWriter,
		Level:           hclog.LevelFromString(option.Level),
		JSONFormat:      option.Format == "json",
	})

	return logger, nil
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
