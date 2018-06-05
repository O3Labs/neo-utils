package neoutils

import (
	"encoding/json"
	"net"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

type customTransport struct {
	rtp       http.RoundTripper
	dialer    *net.Dialer
	connStart time.Time
	connEnd   time.Time
	reqStart  time.Time
	reqEnd    time.Time
}

func newTransport() *customTransport {

	tr := &customTransport{
		dialer: &net.Dialer{
			Timeout:   1 * time.Second, //keep timeout low
			KeepAlive: 1 * time.Second,
		},
	}
	tr.rtp = &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		Dial:                  tr.dial,
		TLSHandshakeTimeout:   1 * time.Second,
		ResponseHeaderTimeout: 1 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	return tr
}

func (tr *customTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	tr.reqStart = time.Now()
	resp, err := tr.rtp.RoundTrip(r)
	tr.reqEnd = time.Now()
	return resp, err
}

func (tr *customTransport) dial(network, addr string) (net.Conn, error) {
	tr.connStart = time.Now()
	cn, err := tr.dialer.Dial(network, addr)
	tr.connEnd = time.Now()
	return cn, err
}

func (tr *customTransport) ReqDuration() time.Duration {
	return tr.Duration() - tr.ConnDuration()
}

func (tr *customTransport) ConnDuration() time.Duration {
	return tr.connEnd.Sub(tr.connStart)
}

func (tr *customTransport) Duration() time.Duration {
	return tr.reqEnd.Sub(tr.reqStart)
}

type BlockCountResponse struct {
	Jsonrpc      string `json:"jsonrpc"`
	ID           int    `json:"id"`
	Result       int    `json:"result"`
	ResponseTime int64  `json:"-"`
}

func fetchSeedNode(url string) *BlockCountResponse {
	//instead of using default http client. we use a transport one here.
	//because we need to mearure the time. Request, Response and total duration.
	//to select the best node among the nodes that has the highest blockcount by picking the least latency node.
	transport := newTransport()
	client := http.Client{Transport: transport}
	payload := strings.NewReader(" {\"jsonrpc\": \"2.0\", \"method\": \"getblockcount\", \"params\": [], \"id\": 3}")
	res, err := client.Post(url, "application/json", payload)
	if err != nil || res == nil {
		return nil
	}
	defer res.Body.Close()
	blockResponse := BlockCountResponse{}
	err = json.NewDecoder(res.Body).Decode(&blockResponse)
	if err != nil {
		return nil
	}
	blockResponse.ResponseTime = transport.ReqDuration().Nanoseconds()
	return &blockResponse
}

type FetchSeedRequest struct {
	Response *BlockCountResponse
	URL      string
}

type SeedNodeResponse struct {
	URL          string
	BlockCount   int
	ResponseTime int64 //milliseconds
}

type NodeList struct {
	URL []string
}

//go mobile bind does not support slice parameters...yet
//https://github.com/golang/go/issues/12113

func SelectBestSeedNode(commaSeparatedURLs string) *SeedNodeResponse {
	urls := strings.Split(commaSeparatedURLs, ",")
	ch := make(chan *FetchSeedRequest, len(urls))
	fetchedList := []string{}
	wg := sync.WaitGroup{}
	listHighestNodes := []SeedNodeResponse{}
	for _, url := range urls {
		go func(url string) {
			res := fetchSeedNode(url)
			ch <- &FetchSeedRequest{res, url}
		}(url)
	}
	wg.Add(1)

loop:
	for {
		select {
		case request := <-ch:
			if request.Response != nil {
				listHighestNodes = append(listHighestNodes, SeedNodeResponse{
					URL:          request.URL,
					BlockCount:   request.Response.Result,
					ResponseTime: request.Response.ResponseTime / int64(time.Millisecond),
				})
			}

			fetchedList = append(fetchedList, request.URL)
			if len(fetchedList) == len(urls) {

				// if len(listHighestNodes) == 0 {
				// 	continue
				// }
				wg.Done()
				break loop
			}
		}
	}
	//wait for the operation
	wg.Wait()
	//using sort.SliceStable to sort min response time first
	sort.SliceStable(listHighestNodes, func(i, j int) bool {
		return listHighestNodes[i].ResponseTime < listHighestNodes[j].ResponseTime
	})
	//using sort.SliceStable to sort block count and preserve the sorted position
	sort.SliceStable(listHighestNodes, func(i, j int) bool {
		return listHighestNodes[i].BlockCount > listHighestNodes[j].BlockCount
	})
	if len(listHighestNodes) == 0 {
		return nil
	}
	return &listHighestNodes[0]
}
