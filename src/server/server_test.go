package server

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestServer_Initialize(t *testing.T) {
	t.Log(os.Getenv("GOOS"), os.Getenv("GOARCH"))
}

func BenchmarkServer(b *testing.B) {

	runCalls(log.Default())
}

func runCalls(logger *log.Logger) {

	cli := http.DefaultClient

	body := `{"method":"account.user.logout","data":"{}","timestamp":"20201029150000","version":"1.0","sign":"D41D8CD98F00B204E9800998ECF8427E","sign_type":"md5"}`
	req, err := http.NewRequest("POST", "http://localhost:8080/account/api", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		logger.Fatal(err)
		return
	}
	n := runtime.NumCPU() -2
	since := time.Now()
	count := 10000 * n
	var wg sync.WaitGroup
	wg.Add(n)

	call := func() {
		defer wg.Done()
		for i := 0; i < count; i++ {
			if resp, err := cli.Post("http://localhost:8080/account/api", "application/json", bytes.NewReader([]byte(body))); err != nil {
				logger.Fatal(err)
				return
			}else {
				resp.Body.Close()
			}

			cli.CloseIdleConnections()
			time.Sleep(time.Millisecond * 10)
		}
	}
	for i:=0; i<n; i++{
		go call()
	}
	wg.Wait()
	logger.Println("latency:", time.Now().Sub(since))
}
