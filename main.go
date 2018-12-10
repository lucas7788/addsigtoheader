package main

import (
	"encoding/json"
	"fmt"
	"github.com/ontio/addsigtoheader/handledb"
	"github.com/ontio/ontology/cmd"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
	"runtime"
)

func setupAPP() *cli.App {
	app := cli.NewApp()
	app.Usage = "Ontology Tool"
	app.Action = startHandleDb
	app.Version = "1.0"
	app.Copyright = "Copyright in 2018 The Ontology Authors"
	app.Before = func(context *cli.Context) error {
		runtime.GOMAXPROCS(runtime.NumCPU())
		return nil
	}
	return app
}

func startHandleDb(ctx *cli.Context) error {

	data, err := ioutil.ReadFile("./config.json")
	if err != nil {
		return err
	}
	config := &handledb.Config{}
	err = json.Unmarshal(data, config)
	if err != nil {
		return err
	}
	//if len(config.WalletDir) != 14 {
	//	return fmt.Errorf("wallet file num is wrong")
	//}
	accsMap, err := handledb.GetAccounts(ctx, config.WalletDir)
	if err != nil {
		return err
	}
	err = handledb.AddSigToHeader(config.RawDbDir, config.ToDbDir, accsMap)
	if err != nil {
		fmt.Println("AddSigToHeader error:", err)
	}
	return nil
}

func main() {
	if err := setupAPP().Run(os.Args); err != nil {
		cmd.PrintErrorMsg(err.Error())
		os.Exit(1)
	}
}
