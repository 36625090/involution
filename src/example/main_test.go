package main

import (
	"bytes"
	"github.com/36625090/involution"
	"github.com/36625090/involution/example/services/account/controller"
	"github.com/36625090/involution/logical"
	"github.com/36625090/involution/option"
	"github.com/36625090/involution/utils"
	"github.com/gin-gonic/gin"
	"github.com/jessevdk/go-flags"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"sync"
	"testing"
	"time"
)

func mockServer() error {
	var opts option.Options
	var args = []string{
		"--log.level",
		"info",
		"--http.path","/account",
		"--http.port","8081",
		"--http.cors",
		"--http.logging",
		"--log.path",
		"../logs",
		"--app","account",
	}
	gin.SetMode("release")
	parser := flags.NewParser(&opts, flags.Default)
	if _, err := parser.ParseArgs(args); err != nil {
		return err
	}
	factories := map[string]logical.Factory{
		"account": controller.Factory,
	}

	inv, err := involution.DefaultInvolution(&opts, factories)
	if err != nil {
		return err
	}

	if err := inv.RegisterService(); err != nil {
		log.Fatal("register service:", err)
		return err
	}

	return inv.Start()
}

func runCalls(logger *log.Logger) {


	n := runtime.NumCPU() -2
	since := time.Now()
	count := 10000
	var wg sync.WaitGroup
	wg.Add(n)

	call := func() {
		cli := http.DefaultClient
		body := `{"method":"account.user.logout","data":"{}","timestamp":"20201029150000","version":"1.0","sign":"D41D8CD98F00B204E9800998ECF8427E","sign_type":"md5"}`
		defer wg.Done()
		for i := 0; i < count; i++ {

			if _, err := cli.Post("http://localhost:8081/account/api", "application/json", bytes.NewReader([]byte(body))); err != nil {
				logger.Fatal(err)
				return
			}
			time.Sleep(time.Millisecond * 10)
		}
	}
	for i:=0; i < n;i++{
		go call()
		log.Println("start call", i)
	}

	wg.Wait()
	logger.Println("calls ",count * n,"finished time latency:", time.Now().Sub(since))
}

func BenchmarkServer(b *testing.B) {
	go func() {
		if err := mockServer(); err != nil {
			b.Fatal(err)
			return
		}
	}()
	time.Sleep(time.Second * 3)
	runCalls(log.Default())
	cli := http.DefaultClient
	resp, err := cli.Get("http://localhost:8081/account/health")
	if err != nil {
		b.Fatal(err)
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		b.Fatal(err)
		return
	}
	b.Log(utils.JSONDump(string(bs)))
}