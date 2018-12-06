package handledb

import (
	"fmt"
	"github.com/ontio/ontology/account"
	cutils "github.com/ontio/ontology/cmd/utils"
	"github.com/ontio/ontology/core/store/ledgerstore"
	"os"
)

func AddSigToHeader(dataDir string, walletDir string) error {
	wallet, err := account.Open(walletDir)
	acc,err := wallet.GetAccountByAddress("AV9STaqJVr1rvibQQB3KM8mc6qYsvgiv8B",[]byte("111111"))
	if err != nil {
		return err
	}
	blockStore, err := ledgerstore.NewBlockStore(fmt.Sprintf("%s%s%s", dataDir, string(os.PathSeparator), ledgerstore.DBDirBlock), true)
	if err != nil {
		return  fmt.Errorf("NewBlockStore error %s", err)
	}

	_, blockCurrHeight, err := blockStore.GetCurrentBlock()
	for i:=0; uint32(i) < blockCurrHeight;i++ {
		blockHash,err := blockStore.GetBlockHash(uint32(i))
		if err !=nil {
			return  fmt.Errorf("GetBlockHash error %s", err)
		}
		block,err := blockStore.GetBlock(blockHash)
		sigdata,err := cutils.Sign(blockHash.ToArray(), acc)
		if err !=nil {
			return  fmt.Errorf("GetBlock error %s", err)
		}
		block.Header.SigData = append(block.Header.SigData, sigdata)
		blockStore.NewBatch()
		err = blockStore.SaveBlock(block)
		if err != nil {
			fmt.Println("SaveBlock, error:", err)
			return err
		}
		blockStore.CommitTo()
	}
	blockStore.Close()
	return nil
}
