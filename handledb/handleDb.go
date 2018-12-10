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

func GetAccounts(ctx *cli.Context, walletDirs []string) (map[string]*account.Account, error) {

	var accsMap = make(map[string]*account.Account, 0)
	for i := 0; i < len(walletDirs); i++ {
		wallet, err := account.Open(walletDirs[i])

		passwd, err := common.GetPasswd(ctx)
		if err != nil {
			return nil, err
		}
		acc, err := wallet.GetDefaultAccount(passwd)
		if err != nil {
			return nil, err
		}
		pubkey := keypair.SerializePublicKey(acc.PublicKey)
		pubkeyStr := ccommon.ToHexString(pubkey)
		accsMap[pubkeyStr] = acc
	}
	return accsMap, nil
}

func getSigAccs(accsMap map[string]*account.Account, peerConfigs []*PeerConfig) ([]*account.Account, error) {

	var sigAccs = make([]*account.Account, 0)
	for i:=0;i<len(peerConfigs);i++ {
		acc := accsMap[peerConfigs[i].ID]
		if acc == nil {
			return nil,fmt.Errorf("no pubkey  error %s",peerConfigs[i].ID)
		}
		sigAccs = append(sigAccs, acc)
	}
	return sigAccs,nil
}
func AddSigToHeader(dataDir, saveToDir string, accsMap map[string]*account.Account) error {

	blockStore, err := ledgerstore.NewBlockStore(fmt.Sprintf("%s%s%s", dataDir, string(os.PathSeparator), ledgerstore.DBDirBlock), true)

	if err != nil {
		return fmt.Errorf("NewBlockStore error %s", err)
	}
	blockStore2, err := ledgerstore.NewBlockStore(fmt.Sprintf("%s%s%s", saveToDir, string(os.PathSeparator), ledgerstore.DBDirBlock), true)

	if err != nil {
		return fmt.Errorf("NewBlockStore error %s", err)
	}

	_, blockCurrHeight, err := blockStore.GetCurrentBlock()
	if err != nil {
		return err
	}

	var lastConfigBlockNum uint32 //记录上一个周期的值

	var sigAccount = make([]*account.Account, 0)

	var lastConfigBlockNumChange = 0

	fmt.Println("GetCurrentBlock:", blockCurrHeight)


	for i := 0; uint32(i) <= blockCurrHeight; i++ {
		blockHash, err := blockStore.GetBlockHash(uint32(i))
		if err != nil {
			return fmt.Errorf("GetBlockHash error %s", err)
		}
		block, err := blockStore.GetBlock(blockHash)
		if err != nil {
			return fmt.Errorf("GetBlock error %s, blockHash:%s", err, blockHash)
		}

		blkInfo := &VbftBlockInfo{}
		if err := json.Unmarshal(block.Header.ConsensusPayload, blkInfo); err != nil {
			return fmt.Errorf("unmarshal blockInfo: %s", err)
		}
		if i == 0 {
			lastConfigBlockNum = blkInfo.LastConfigBlockNum
			sigAccount,err = getSigAccs(accsMap, blkInfo.NewChainConfig.Peers)
			if err != nil {
				return fmt.Errorf("getSigAccs: %s", err)
			}
			continue
		}
        if i == 1 {
			lastConfigBlockNum = blkInfo.LastConfigBlockNum
		}else {
			if lastConfigBlockNum != blkInfo.LastConfigBlockNum {
				lastConfigBlockNumChange++
				lastConfigBlockNum = blkInfo.LastConfigBlockNum
				//获得需要签名的account
				sigAccount,err = getSigAccs(accsMap, blkInfo.NewChainConfig.Peers)
				if err != nil {
					return fmt.Errorf("i>0, getSigAccs: %s", err)
				}
			}
		}
		if len(sigAccount) != 7 {
			return fmt.Errorf("sigAccount length is not 7 error")
		}
		var accSig = make([][]byte, 0)
		var bookKeepers = make([]keypair.PublicKey, 0)
		for k := 0; k < len(sigAccount); k++ {
			sigData, err := utils.Sign(blockHash.ToArray(), sigAccount[k])
			if err != nil {
				return fmt.Errorf("Sign error %s", err)
			}
			accSig = append(accSig, sigData)
			bookKeepers = append(bookKeepers, sigAccount[k].PublicKey)
		}
		fmt.Println(fmt.Sprintf("blockHeight: %d, before: %d, after: %d", i, len(block.Header.SigData), len(accSig)))
		block.Header.Bookkeepers = bookKeepers
		block.Header.SigData = accSig
		blockStore2.NewBatch()
		err = blockStore2.SaveBlock(block)
		if err != nil {
			fmt.Println("SaveBlock, error:", err)
			return err
		}
		blockStore2.CommitTo()
	}
	fmt.Println("lastConfigBlockNumChange:", lastConfigBlockNumChange)
	blockStore.Close()
	blockStore2.Close()
	return nil
}
