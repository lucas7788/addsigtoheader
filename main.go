package main

import (
	"addsigtoheader/handledb"
	"fmt"
)

func main(){
	err := handledb.AddSigToHeader("/Users/sss/gopath/src/github.com/ontio/ontology/Chain/ontology",
		"/Users/sss/gopath/src/github.com/ontio/ontology/wallet.dat")
	if err != nil {
		fmt.Println("AddSigToHeader error:", err)
	}
}
