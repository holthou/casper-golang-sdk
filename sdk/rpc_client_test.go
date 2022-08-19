package sdk

import (
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math/big"
	"testing"
)

//var client = NewRpcClient("http://3.136.227.9:7777/rpc")
var client = NewRpcClient("https://node-clarity-testnet.make.services/rpc")

func TestRpcClient_GetLatestBlock(t *testing.T) {
	b, err := client.GetLatestBlock()

	if err != nil {
		t.Errorf("can't get latest block")
	}
	fmt.Println(b)
}

func TestRpcClient_GetDeploy(t *testing.T) {
	//hash := "1dfdf144eb0422eae3076cd8a17e55089010a133c6555c881ac0b9e2714a1605"
	hash := "107b9d1405b1eb2778da6cbb0ecf09cd1026c9f7e13c6e92906ab4d6082cfcd1"
	result, err := client.GetDeploy(hash)

	if err != nil {
		t.Errorf("can't get deploy info")
	}

	for _, test := range result.Deploy.Session.Transfer.Args {
		if len(test) != 2 {
			continue
		}

		b, ok := test[1].(map[string]interface{})
		if !ok {
			continue
		}

		switch test[0] {
		case "amount":
			a, ok := b["parsed"]
			if ok {
				fmt.Println(a)
			}
		case "target":
			if b["cl_type"] != "PublicKey" {
				continue
			}
			a, ok := b["parsed"]
			if ok {
				fmt.Println(a)
			}
		}

	}
}

func TestRpcClient_GetBlockState(t *testing.T) {
	stateRootHash := "c0eb76e0c3c7a928a0cb43e82eb4fad683d9ad626bcd3b7835a466c0587b0fff"
	key := "account-hash-a9efd010c7cee2245b5bad77e70d9beb73c8776cbe4698b2d8fdf6c8433d5ba0"
	path := []string{"special_value"}
	_, err := client.GetStateItem(stateRootHash, key, path)

	if err != nil {
		t.Errorf("can't get block state")
	}
}

func TestRpcClient_GetAccountBalance(t *testing.T) {
	stateRootHash := "c0eb76e0c3c7a928a0cb43e82eb4fad683d9ad626bcd3b7835a466c0587b0fff"
	key := "account-hash-a9efd010c7cee2245b5bad77e70d9beb73c8776cbe4698b2d8fdf6c8433d5ba0"

	balanceUref := client.GetAccountMainPurseURef(key)

	_, err := client.GetAccountBalance(stateRootHash, balanceUref)

	if err != nil {
		t.Errorf("can't get account balance")
	}
}

func TestRpcClient_GetAccountBalanceByKeypair(t *testing.T) {
	stateRootHashNew, err := client.GetStateRootHash()
	if err != nil {
		return
	}
	path := []string{"special_value"}
	_, err = client.GetStateItem(stateRootHashNew.StateRootHash, sourceKeyPair.AccountHash(), path)
	if err != nil {
		_, err := client.GetAccountBalanceByKeypair(stateRootHashNew.StateRootHash, sourceKeyPair)

		if err != nil {
			t.Errorf("can't get account balance")
		}
	}
}

func TestRpcClient_GetBlockByHeight(t *testing.T) {
	_, err := client.GetBlockByHeight(1034)

	if err != nil {
		t.Errorf("can't get block by height")
	}
}

func TestRpcClient_GetBlockTransfersByHeight(t *testing.T) {
	_, err := client.GetBlockByHeight(1034)

	if err != nil {
		t.Errorf("can't get block transfers by height")
	}
}

func TestRpcClient_GetBlockByHash(t *testing.T) {
	_, err := client.GetBlockByHash("")

	if err != nil {
		t.Errorf("can't get block by hash")
	}
}

func TestRpcClient_GetBlockTransfersByHash(t *testing.T) {
	_, err := client.GetBlockTransfersByHash("")

	if err != nil {
		t.Errorf("can't get block transfers by hash")
	}
}

func TestRpcClient_GetLatestBlockTransfers(t *testing.T) {
	_, err := client.GetLatestBlockTransfers()

	if err != nil {
		t.Errorf("can't get latest block transfers")
	}
}

func TestRpcClient_GetValidator(t *testing.T) {
	_, err := client.GetValidator()

	if err != nil {
		t.Errorf("can't get validator")
	}
}

func TestRpcClient_GetStatus(t *testing.T) {
	_, err := client.GetStatus()

	if err != nil {
		t.Errorf("can't get status")
	}
}

func TestRpcClient_GetPeers(t *testing.T) {
	_, err := client.GetPeers()

	if err != nil {
		t.Errorf("can't get peers")
	}
}

//make sure your account has balance
func TestRpcClient_PutDeploy(t *testing.T) {
	deploy := NewTransferToUniqAddress(*source, UniqAddress{
		PublicKey:  dest,
		TransferId: 10,
	}, big.NewInt(3000000000), big.NewInt(10000), "casper-test", "")

	assert.True(t, deploy.ValidateDeploy())
	deploy.SignDeploy(sourceKeyPair)

	result, err := client.PutDeploy(*deploy)

	if !assert.NoError(t, err) {
		t.Errorf("error : %v", err)
		return
	}

	assert.Equal(t, hex.EncodeToString(deploy.Hash), result.Hash)
}

//make sure your account has balance
func TestRpcClient_delegate(t *testing.T) {
	//delegate file
	modulePath := "../keypair/test_account_keys/contract/delegate.wasm"
	module, _ := ioutil.ReadFile(modulePath)
	validator := "0109b48a169e6163078a07b6248f330133236c6e390fe915813c187c3f268c213e"

	deploy := NewDelegate(*source, validator, big.NewInt(500000000000), big.NewInt(3000000000), "casper-test", module)

	assert.True(t, deploy.ValidateDeploy())
	deploy.SignDeploy(sourceKeyPair)

	result, err := client.PutDeploy(*deploy)

	if !assert.NoError(t, err) {
		t.Errorf("error : %v", err)
		return
	}

	assert.Equal(t, hex.EncodeToString(deploy.Hash), result.Hash)
}

func TestRpcClient_undelegate(t *testing.T) {
	//undelegate file
	modulePath := "../keypair/test_account_keys/contract/undelegate.wasm"
	module, _ := ioutil.ReadFile(modulePath)
	validator := "0109b48a169e6163078a07b6248f330133236c6e390fe915813c187c3f268c213e"

	deploy := NewDelegate(*source, validator, big.NewInt(500000000000), big.NewInt(1000000000), "casper-test", module)

	assert.True(t, deploy.ValidateDeploy())
	deploy.SignDeploy(sourceKeyPair)

	result, err := client.PutDeploy(*deploy)

	if !assert.NoError(t, err) {
		t.Errorf("error : %v", err)
		return
	}

	fmt.Println(result.Hash)
	assert.Equal(t, hex.EncodeToString(deploy.Hash), result.Hash)
}

func TestRpcClient_GetBalance(t *testing.T) {
	publicKey := "01c68d9b5eed7c7e8e8901fd1c354ab044859f45a579c632c48137b2a7d45fe3c0"
	result, err := client.GetLiquidBalance(publicKey)
	assert.NoError(t, err)
	fmt.Printf("Liquid balance:%s \n", result.String())

	result2, err := client.GetStackingBalance(publicKey)
	assert.NoError(t, err)
	fmt.Printf("Staked balance:%s \n", result2.String())
}
