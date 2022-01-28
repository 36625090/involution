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
		"--http.trace",
		"--log.path",
		"../logs",
		"--app","example",
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

	err = inv.Start()
	if err != nil {
		return err
	}
	return nil
}

func runCalls(b *testing.B) {


	n := runtime.NumCPU() -4
	since := time.Now()
	count := 10
	var wg sync.WaitGroup
	wg.Add(n)

	call := func() {
		cli := http.Client{
			Transport:     &http.Transport{
			},
			Timeout:       time.Second * 30,
		}
		body := `{"method":"account.user.login","data":"{\"mobile\":\"12312312123\",\"source\":\"122\",\"open_id\":null,\"union_id\":null,\"verify_code\":\"222\"}","timestamp":161239123,"version":"1.0","sign":"D41D8CD98F00B204E9800998ECF8427E","sign_type":"md5"}`
		defer wg.Done()


		for i := 0; i < count; i++ {
			var reader = bytes.NewReader([]byte(body))
			req, err := http.NewRequest("POST", "http://localhost:8080/example/api", reader)
			if err != nil {
				return
			}
			req.Header.Set("Content-Type", "application/json")
			if resp, err := cli.Do(req); err != nil {
				b.Error(err)
			}else {
				data, _  := ioutil.ReadAll(resp.Body)
				b.Log(string(data))
				resp.Body.Close()
			}
		}
	}
	for i:=0; i < n;i++{
		go call()
		log.Println("start call", i)
	}

	wg.Wait()
	b.Log("calls ",count * n,"finished time latency:", time.Now().Sub(since))
}

func BenchmarkServer(b *testing.B) {
	//go func() {
	//	if err := mockServer(); err != nil {
	//		b.Fatal(err)
	//		return
	//	}
	//}()
	time.Sleep(time.Second * 3)
	runCalls(b)
	cli := http.DefaultClient
	resp, err := cli.Get("http://localhost:8080/example/health")
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