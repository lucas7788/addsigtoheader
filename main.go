package main

import (
	"fmt"
	"github.com/ontio/addsigtoheader/handledb"
	"github.com/ontio/ontology/cmd"
	"github.com/urfave/cli"
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
	//原数据库目录
	rawDbDir := "/Users/sss/gopath/src/github.com/ontio/ontology/Chain/ontology"
	//保存到哪个目录
	toDbDir := "/Users/sss/gopath/src/github.com/ontio/ontology/Chain/ontology2"
	//钱包文件目录
	walletDir := []string{
		"/Users/sss/gopath/src/github.com/ontio/ontology/wallet.dat",
		"/Users/sss/gopath/src/github.com/ontio/ontology/wallet.dat",
	}
	accs, err := handledb.GetAccounts(ctx, walletDir)
	if err != nil {
		return err
	}
	err = handledb.AddSigToHeader(rawDbDir, toDbDir, accs)
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
