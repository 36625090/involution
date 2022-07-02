package option

import (
	"github.com/go-various/consul"
	"github.com/jessevdk/go-flags"
	"os"
)

type Consul struct {
	DC           string `long:"consul.dc" default:"dc1" description:"Datacenter for consul cluster"`
	Address      string `long:"consul.address" default:"127.0.0.1:8500" description:"Sets the consul address"`
	AclToken     string `long:"consul.acl_token" description:"Token for consul config read"`
	ConfigKey    string `long:"consul.config_key" description:"Key for consul config read"`
	ConfigFormat string `long:"consul.config_format" default:"hcl" choice:"hcl" choice:"yaml" choice:"json" choice:"properties" description:"Format for config content"`
}

type Http struct {
	Path         string `long:"http.path" default:"" description:"Path for http server context"`
	Address      string `long:"http.address" default:"0.0.0.0" description:"Address for http server listening"`
	Port         int    `long:"http.port" default:"8080" description:"Port for http server listening"`
	Cors         bool   `long:"http.cors" description:"Support Cors access"`
	Trace        bool   `long:"http.trace" description:"Trace http requests"`
	Sign         bool   `long:"http.sign" description:"Sign verification requests"`
	IdleTimeout  int    `long:"http.idle" default:"30"  description:"Timeout(seconds) for idle connection"`
	ReadTimeout  int    `long:"http.read" default:"5" description:"Timeout(seconds) for read  client request"`
	WriteTimeout int    `long:"http.write" default:"10" description:"Timeout(seconds) for write to client request"`
	KeepAlive    bool   `long:"http.keepalive" description:"Keep-Alive"`
}

//Log logging settings
type Log struct {
	Console bool   `long:"log.console" description:"Set log output to console"`
	Path    string `long:"log.path" default:"logs" required:"true" description:"Sets the path to log file"`
	Level   string `long:"log.level" default:"info" description:"Sets the log level" choice:"info" choice:"warn" choice:"error" choice:"debug" choice:"trace" `
	Format  string `long:"log.format" default:"text" description:"Sets the log format" choice:"text" choice:"json"`
	Rotate  string `long:"log.rotate" default:"day" description:"Rotates the log" choice:"day" choice:"hour" `
}

//Options 服务参数选项（OOPs 英语不好，注释描述凑合看，写中文怕终端乱码 ^ - ^）
type Options struct {
	App           string `long:"app" required:"true" description:"App name for service"`
	Profile       string `long:"profile" description:" Profile for runtime"`
	ConfigFile    string `long:"config" description:"Config file for runtime"`
	Log           Log    `group:"log"`
	UseConsul     bool   `long:"consul" description:"Enable consul"`
	Consul        Consul `group:"consul"`
	Http          Http   `group:"http"`
	Pprof         bool   `long:"pprof" description:"Enable profiling"`
	PprofAddr     string `long:"pprof.port" default:"127.0.0.1:32768" description:"Listen port on Pprof server"`
	Ui            bool   `long:"ui" description:"Enable document ui support"`
	Newrelic      bool   `long:"newrelic" description:"Enable newrelic support"`
	NewrelicKey   string `long:"newrelic.key" description:"Key for newrelic access"`
	NewrelicTrace bool   `long:"newrelic.trace" description:"Trace on newrelic access"`
}

var opts Options
var parser *flags.Parser

func init() {
	parser = flags.NewParser(&opts, flags.Default)
}

func NewOptions() (*Options, error) {
	if _, err := parser.ParseArgs(os.Args[1:]); err != nil {
		//os.Stdout.WriteString(err.Error())
		return nil, err
	}
	return &opts, nil
}

func (o *Options) ConsulConfig() *consul.Config {

	return &consul.Config{
		Datacenter:  o.Consul.DC,
		ZoneAddress: o.Consul.Address,
		Token:       o.Consul.AclToken,
		Application: struct {
			Name    string
			Profile string
		}{
			Name:    o.App,
			Profile: o.Profile,
		},
		Config: struct {
			DataKey string
			Format  string
		}{
			DataKey: o.Consul.ConfigKey,
			Format:  o.Consul.ConfigFormat,
		},
	}
}
