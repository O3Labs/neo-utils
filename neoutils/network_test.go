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
		"http://seed2.neo.org:10332",
		"https://seed1.neo.org:10331",
		"http://seed3.neo.org:10332",
		"http://seed4.neo.org:10332",
		"http://seed5.neo.org:10332",
		"http://seed1.cityofzion.io:8080",
		"http://seed2.cityofzion.io:8080",
		"http://seed3.cityofzion.io:8080",
		"http://seed4.cityofzion.io:8080",
		"http://seed2.aphelion-neo.com:10332",
		"http://seed2.aphelion-neo.com:10332",
		"https://seed3.switcheo.network:10331",
		"https://seed2.switcheo.network:10331",

		"http://node2.ams2.bridgeprotocol.io:10332",
		"http://seed1.o3node.org:10332",
		"http://seed2.o3node.org:10332",
		"http://seed3.o3node.org:10332",
	}
	commaSeparated := strings.Join(urls, ",")
	best := SelectBestSeedNode(commaSeparated)
	if best != nil {
		log.Printf("best node %+v %v %vms", best.URL, best.BlockCount, best.ResponseTime)
	}
}

func TestGetBestO3Node(t *testing.T) {
	urls := []string{
		"http://seed1.o3node.org:10332",
		"http://seed2.o3node.org:10332",
		"http://seed3.o3node.org:10332",
	}
	commaSeparated := strings.Join(urls, ",")
	best := SelectBestSeedNode(commaSeparated)
	if best != nil {
		log.Printf("best node %+v %v %vms", best.URL, best.BlockCount, best.ResponseTime)
	}
}
