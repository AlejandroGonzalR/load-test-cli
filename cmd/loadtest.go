package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/tcnksm/go-httpstat"
)

func ExecRequest(url string) {
	rep := ClientOpts.Concurrency
	results := make(chan map[string]int)

	// Use goroutine to send multiple time-consuming jobs to the channel.
	for i := 0; i < rep; i++ {
		go func(num int) {
			results <- mockHTTPRequest(url)
		}(i)
	}

	// Receive results from the channel and use them.
	for i := 0; i < rep; i++ {
		for key, value := range <-results {
			log.Printf("%s: %dms\n", key, value)
		}
	}
}

func mockHTTPRequest(url string) map[string]int {
	// Create a new HTTP request
	req, err := http.NewRequest(ClientOpts.Method, url, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Create a httpstat powered context
	var result httpstat.Result
	ctx := httpstat.WithHTTPStat(req.Context(), &result)
	req = req.WithContext(ctx)

	// Send request by default HTTP client
	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if _, err := io.Copy(ioutil.Discard, res.Body); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	res.Body.Close()

	results := make(map[string]int)
	results["DNS lookup"] = int(result.TCPConnection / time.Millisecond)
	results["TCP connection"] = int(result.TCPConnection / time.Millisecond)
	results["TLS handshake"] = int(result.TLSHandshake / time.Millisecond)
	results["Server processing"] = int(result.ServerProcessing / time.Millisecond)
	results["Content transfer"] = int(result.ContentTransfer(time.Now()) / time.Millisecond)

	return results
}
