package neoutils

import (
	"log"
	"strings"
	"testing"
)

func TestFetchBlockCount(t *testing.T) {
	url := "http://seed2.neo.org:10332/"
	res := fetchSeedNode(url)
	log.Printf("%v", res)
}

func TestFetchDownNodeBlockCount(t *testing.T) {
	url := "http://seed1.cityofzion.io:8080"
	res := fetchSeedNode(url)
	log.Printf("%v", res)
}

func TestBestNode(t *testing.T) {
	urls := []string{
		"http://seed1.neo.org:10332",
		"http://seed2.neo.org:10332",
		"http://seed3.neo.org:10332",
		"http://seed4.neo.org:10332",
		"http://seed5.neo.org:10332",
		"http://seed1.cityofzion.io:8080",
		"http://seed2.cityofzion.io:8080",
		"http://seed3.cityofzion.io:8080",
		"http://seed4.cityofzion.io:8080",
		"http://seed5.cityofzion.io:8080",
		"http://node1.o3.network:10332",
		"http://node2.o3.network:10332",
	}
	commaSeparated := strings.Join(urls, ",")
	best := SelectBestSeedNode(commaSeparated)
	if best != nil {
		log.Printf("best node %+v %v %vms", best.URL, best.BlockCount, best.ResponseTime)
	}

}
