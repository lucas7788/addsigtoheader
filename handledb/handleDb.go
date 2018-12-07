package handledb

import (
	"encoding/json"
	"fmt"
	"github.com/ontio/ontology-crypto/keypair"
	"github.com/ontio/ontology/account"
	"github.com/ontio/ontology/cmd/common"
	"github.com/ontio/ontology/cmd/utils"
	ccommon "github.com/ontio/ontology/common"
	. "github.com/ontio/ontology/consensus/vbft/config"
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

func isChangePubKey(accs []*account.Account, peerConfigs []*PeerConfig) bool {
	var samePubNum int = 0
	for m:=0; m < len(accs);m++ {
		pks := keypair.SerializePublicKey(accs[m].PublicKey)
		pkstr := ccommon.ToHexString(pks)
		for i:=0;i<len(peerConfigs);i++ {
			if pkstr == peerConfigs[i].ID {
				samePubNum++
				break
			}
		}
	}
	if samePubNum == 6 {
		return true
	}
	return false
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

	var lastConfigBlockNum uint32 //记录上一个周期的值

	var totalChangeNum int //用来统计变化次数

	var start int = 0  //用来记录从第几个钱包开始追加签名
	var end int = 7

	var peerConfigs []*PeerConfig

	for i := 0; uint32(i) <= blockCurrHeight; i++ {
		blockHash, err := blockStore.GetBlockHash(uint32(i))
		if err != nil {
			return fmt.Errorf("GetBlockHash error %s", err)
		}
		block, err := blockStore.GetBlock(blockHash)

		blkInfo := &VbftBlockInfo{}
		if err := json.Unmarshal(block.Header.ConsensusPayload, blkInfo); err != nil {
			return fmt.Errorf("unmarshal blockInfo: %s", err)
		}
        if i == 0 {
			lastConfigBlockNum = blkInfo.LastConfigBlockNum
			peerConfigs = blkInfo.NewChainConfig.Peers
		}else {
			if lastConfigBlockNum != blkInfo.LastConfigBlockNum {
				lastConfigBlockNum = blkInfo.LastConfigBlockNum
				//比较公钥是否改变
				b := isChangePubKey(accs[start:end], peerConfigs)
				if b {
					totalChangeNum = totalChangeNum + 1
					start = 7
					end = 14
				}
			}
		}
		var accSig [][]byte
		for k := start; k < end; k++ {
			sigData, err := utils.Sign(blockHash.ToArray(), accs[k])
			if err != nil {
				return fmt.Errorf("GetBlock error %s", err)
			}
			accSig = append(accSig, sigData)
		}
		block.Header.SigData = accSig
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
