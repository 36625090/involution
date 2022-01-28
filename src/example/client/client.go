package main

import (
	"bytes"
	"context"
	"github.com/36625090/involution/utils"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"sync"
	"time"
)

func main() {
	debug.SetPanicOnFault(true)
	if len(os.Args) < 3{
		os.Args = append(os.Args, "1000")
	}
	url := os.Args[2]
	log.SetFlags(log.LstdFlags|log.Lshortfile)
	count, _ := strconv.Atoi(os.Args[1])
	//time.Sleep(time.Second * 3)
	log.Println(count, url)
	runCalls(count, url)
	cli := http.DefaultClient
	resp, err := cli.Get(url + "/example/health")
	if err != nil {
		log.Fatal(err)
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println(utils.JSONDump(string(bs)))
}

func runCalls(count int, url string) {

	n := runtime.NumCPU()
	since := time.Now()
	var wg sync.WaitGroup
	wg.Add(n)

	call := func() {
		body := `{"method":"account.user.home","data":"{\"mobile\":\"12312312123\",\"source\":\"122\",\"open_id\":null,\"union_id\":null,\"verify_code\":\"222\"}","timestamp":161239123,"version":"1.0","sign":"D41D8CD98F00B204E9800998ECF8427E","sign_type":"md5"}`
		defer wg.Done()
		cli := http.DefaultClient

		for i := 0; i < count; i++ {
			var reader = bytes.NewReader([]byte(body))
			req, err := http.NewRequest("POST", url+"/example/api", reader)
			if err != nil {
				return
			}

			req.Header.Set("Content-Type", "application/json")
			if _, err := cli.Do(req); err != nil {
				log.Println("ERROR:",err, req.Response)
				debug.PrintStack()
			}
			req.Clone(context.Background())
		}
	}

	for i:=0; i < n;i++{
		go call()
		log.Printf("start call: %d, %d/pre routine", i,count)
	}

	wg.Wait()
	diff := time.Now().Sub(since)
	log.Println("calls ",count * n,"finished time latency:", diff," ",
		float64(diff.Milliseconds()) / float64(count * n),"ms/pre request")
}