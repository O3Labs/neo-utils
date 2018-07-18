package neorpc_test

import (
	"log"
	"testing"

	"github.com/o3labs/neo-utils/neoutils/neorpc"
)

func TestEndpoint(t *testing.T) {
	client := neorpc.NewClient("http://localhost:30333")
	if client == nil {
		t.Fail()
	}
	log.Printf("%v", client)
}

func TestGetContractState(t *testing.T) {
	client := neorpc.NewClient("http://localhost:30333")
	if client == nil {
		t.Fail()
	}

	result := client.GetContractState("ce575ae1bb6153330d20c560acb434dc5755241b")
	log.Printf("%+v", result)
}

func TestSendRawTransaction(t *testing.T) {
	client := neorpc.NewClient("http://localhost:30333")
	if client == nil {
		t.Fail()
	}
	raw := "d1004208e8030000000000001423ba2703c53263e8d6e522dc32203339dcd8eee952c10c6d696e74546f6b656e73546f671b245557dc34b4ac60c5200d335361bbe15a57ce01f11e74686973697361756e69717565746f6b656e5f66726f6d5f73747269706501e216181b1f9a773f93064af30be44679f34ec878788afa1727aa60057eb39a96000001e72d286979ee6cb1b7e65dfddfb2e384100b8d148e7758de42e4168b71792c60010000000000000023ba2703c53263e8d6e522dc32203339dcd8eee9014140f55e2b2914c409396904b8c5a1e8ec0ffc0b62f8b1b996beae7c65ceca7e11a3dbab011038b948ec380c5b22ba474f013ca6de61051dda487a5bec17196115412321031a6c6fbbdf02ca351745fa86b9ba5a9452d785ac4f7fc2b7548ca2a46c4fcf4aacce575ae1bb6153330d20c560acb434dc5755241b"
	result := client.SendRawTransaction(raw)
	log.Printf("%+v", result)
}

func TestGetRawTransaction(t *testing.T) {

	client := neorpc.NewClient("http://localhost:30333")
	if client == nil {
		t.Fail()
	}
	txID := "bde02f8c6482e23d5b465259e3e438f0acacaba2a7a938d5eecd90bba0e9d1ad"
	result := client.GetRawTransaction(txID)
	log.Printf("%+v", result)
}

func TestGetBlock(t *testing.T) {
	client := neorpc.NewClient("http://seed2.o3node.org:10332")
	if client == nil {
		t.Fail()
	}
	txID := "5ba40a700fbdd72344d2903629fac10b55e7a957d17d38e475a20ab18766fa7b"
	result := client.GetBlock(txID)
	log.Printf("%+v", len(result.Result.Tx))
}

func TestGetBlockByIndex(t *testing.T) {
	client := neorpc.NewClient("http://seed2.o3node.org:10332")
	if client == nil {
		t.Fail()
	}
	index := 2188171
	result := client.GetBlockByIndex(index)
	log.Printf("%+v", result)
}

func TestGetBlockCount(t *testing.T) {
	client := neorpc.NewClient("http://seed2.o3node.org:10332")
	if client == nil {
		t.Fail()
	}

	result := client.GetBlockCount()
	log.Printf("%+v", result.Result)
}

func TestGetAccountState(t *testing.T) {
	client := neorpc.NewClient("http://seed2.o3node.org:10332")
	if client == nil {
		t.Fail()
	}

	result := client.GetAccountState("AdSBfV9kMmN2Q3xMYSbU33HWQA1dCc9CV3")
	log.Printf("%+v", result.Result)
}

func TestGetTokenBalance(t *testing.T) {
	client := neorpc.NewClient("http://seed3.aphelion-neo.com:10332")
	if client == nil {
		t.Fail()
	}

	result := client.GetTokenBalance("fc732edee1efdf968c23c20a9628eaa5a6ccb934", "AcydXy1MvrzaT8qD3Qe4B8mqEoinTvRy8U")
	log.Printf("%+v", result.Result)
}

func TestInvokeScript(t *testing.T) {
	client := neorpc.NewClient("http://seed2.neo.org:20332")
	script := "00c1046e616d6567f8e679d19048360e414c82d82fdb33486438d37c00c10673796d626f6c67f8e679d19048360e414c82d82fdb33486438d37c00c10b746f74616c537570706c7967f8e679d19048360e414c82d82fdb33486438d37c"
	if client == nil {
		t.Fail()
	}

	result := client.InvokeScript(script)
	log.Printf("%+v", result.Result)
}
