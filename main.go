package main

import (
	"fmt"
	"log"

	"github.com/ashwanthkumar/wasp-cli/client"
	"github.com/ashwanthkumar/wasp-cli/util"
	"github.com/buger/jsonparser"

	"github.com/indix/abelwatch/abel"
)

func main() {
	wasp := client.WASP{
		Url: "http://wasp.indix.tv:9000",
	}
	data, err := wasp.List("production")
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	var keys []string
	util.JsonDecode(data, &keys)
	for _, key := range keys {
		fmt.Println(key)
	}

	fmt.Println("==================================================")

	abelClient := abel.Abel{
		URL: "http://abel.prod.indix.tv:3330",
	}
	count, datatype, err := abelClient.Get("completed", []string{"rmn_p1_variant_20180918", "www.finishline.com"})
	if err != nil || datatype == jsonparser.NotExist {
		log.Fatalf("%v\n", err)
	}
	fmt.Printf("Count: %d\n", count)
}
