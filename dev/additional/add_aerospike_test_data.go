package main

import (
	as "github.com/aerospike/aerospike-client-go"
)

func main() {
	// Initialize policy.
	policy := as.NewWritePolicy(0, 0)
	client, err := as.NewClient("127.0.0.1", 3000)
	if err != nil {
		panic(err)
	}

	// ---------------------
	// ---------------------
	// ---------------------
	key, _ := as.NewKey("test", "campaigns", 3)
	//
	bin := make(map[string]string)
	bin["icon"] = "http://logiq.one/img/blog2.jpg"
	bin["image"] = "http://logiq.one/img/blog2.jpg"
	bin["click-action"] = "http://gurfv.pro/?target=-7EBNQCgQAAAOSUwNXQQAFAQEREQoRCQoRDUIRDRIAAX9hZGNvbWJvATE"
	bin["country"] = "THN"
	bin["node"] = "test_node"
	var bins []*as.Bin
	binU := as.NewBin("user-id", 7651)

	for k, v := range bin {
		bins = append(bins, as.NewBin(k, v))
	}
	binP := as.NewBin("price", 10000.1)
	err = client.PutBins(policy, key, binU, bins[0], bins[1], bins[2], bins[3], bins[4], binP)
	if err != nil {
		panic(err)
	}

	// ---------------------
	//existed, err := client.Delete(policy, key)
	//if  err != nil {
	//	panic(err)
	//}
	//log.Println(existed)
	// ---------------------
	// ---------------------
	// ---------------------

	//key2, _ := as.NewKey("test", "dsp_proxy_tracker", "push:1")
	//binIconUrl := as.NewBin("click-action", "http://localhost:8090/click?campaign-id=2")
	//err = client.PutBins(policy, key2, binIconUrl)
	//if err != nil {
	//	panic(err)
	//}
}
