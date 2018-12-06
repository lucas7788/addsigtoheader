package handledb

import (
	"bytes"
	"fmt"
	"github.com/ontio/ontology/account"
	"github.com/ontio/ontology/cmd/common"
	"github.com/ontio/ontology/cmd/utils"
	"github.com/ontio/ontology/core/store/ledgerstore"
	"github.com/urfave/cli"
	"os"
)

func GetAccounts(ctx *cli.Context, walletDirs []string) ([]*account.Account, error) {
	var accs []*account.Account
	for i := 0; i < len(walletDirs); i++ {
		wallet, err := account.Open(walletDirs[i])

		passwd, err := common.GetPasswd(ctx)
		if err != nil {
			return nil, err
		}
		defer common.ClearPasswd(passwd)
		acc, err := wallet.GetAccountByIndex(0, passwd)
		if err != nil {
			return nil, err
		}
		accs = append(accs, acc)
	}
	return accs, nil
}
func AddSigToHeader(dataDir, saveToDir string, accs []*account.Account) error {

	blockStore, err := ledgerstore.NewBlockStore(fmt.Sprintf("%s%s%s", dataDir, string(os.PathSeparator), ledgerstore.DBDirBlock), true)

	if err != nil {
		return fmt.Errorf("NewBlockStore error %s", err)
	}
	blockStore2, err := ledgerstore.NewBlockStore(fmt.Sprintf("%s%s%s", saveToDir, string(os.PathSeparator), ledgerstore.DBDirBlock), true)

	if err != nil {
		return fmt.Errorf("NewBlockStore error %s", err)
	}

	_, blockCurrHeight, err := blockStore.GetCurrentBlock()
	for i := 0; uint32(i) <= blockCurrHeight; i++ {
		blockHash, err := blockStore.GetBlockHash(uint32(i))
		if err != nil {
			return fmt.Errorf("GetBlockHash error %s", err)
		}
		block, err := blockStore.GetBlock(blockHash)
		for i := 0; i < len(accs); i++ {
			sigdata, err := utils.Sign(blockHash.ToArray(), accs[i])
			if err != nil {
				return fmt.Errorf("GetBlock error %s", err)
			}
			hasSig := block.Header.SigData[:]
			for j := 0; j < len(hasSig); j++ {
				if bytes.Contains(hasSig[j], sigdata) {
					continue
				}
			}
			block.Header.SigData = append(block.Header.SigData, sigdata)
		}
		blockStore2.NewBatch()
		err = blockStore2.SaveBlock(block)
		if err != nil {
			fmt.Println("SaveBlock, error:", err)
			return err
		}
		blockStore2.CommitTo()
	}
	blockStore2.Close()
	return nil
}
