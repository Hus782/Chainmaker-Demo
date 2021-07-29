/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main
import (
	//"chainmaker.org/chainmaker-go/common/log"
	sdk "chainmaker.org/chainmaker-sdk-go"
	"chainmaker.org/chainmaker-sdk-go/pb/protogo/common"
	//"chainmaker.org/chainmaker-sdk-go/pb/protogo/config"
	"errors"
	"fmt"
	//"go.uber.org/zap"
	"io/ioutil"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"strings"
	"math/big"
)

const (
	createContractTimeout = 5
	chainId        = "chain1"
	orgId1         = "wx-org1.chainmaker.org"
	orgId2         = "wx-org2.chainmaker.org"
	orgId3         = "wx-org3.chainmaker.org"
	orgId4         = "wx-org4.chainmaker.org"

	certPathPrefix = "./testdata"
	tlsHostName    = "chainmaker.org"

	nodeAddr1 = "127.0.0.1:12301"
	connCnt1  = 5

	nodeAddr2 = "127.0.0.1:12301"
	connCnt2  = 5

	certPathFormat = "/crypto-config/%s/ca"

	voteContractName = "Vote006"
	voteVersion      = "1.0.0"
	voteByteCodePath = "./testdata/vote/Vote.bin"
	voteABIPath      = "./testdata/vote/Vote.abi"
)
var (
	caPaths = []string{
		certPathPrefix + fmt.Sprintf(certPathFormat, orgId1),
		certPathPrefix + fmt.Sprintf(certPathFormat, orgId2),
		certPathPrefix + fmt.Sprintf(certPathFormat, orgId3),
		certPathPrefix + fmt.Sprintf(certPathFormat, orgId4),
	}


	userKeyPath = certPathPrefix + "/crypto-config/%s/user/client1/client1.tls.key"
	userCrtPath = certPathPrefix + "/crypto-config/%s/user/client1/client1.tls.crt"
	userSignKeyPath = certPathPrefix + "/crypto-config/%s/user/client1/client1.sign.key"
	userSignCrtPath = certPathPrefix + "/crypto-config/%s/user/client1/client1.sign.crt"

	adminKeyPath = certPathPrefix + "/crypto-config/%s/user/admin1/admin1.tls.key"
	adminCrtPath = certPathPrefix + "/crypto-config/%s/user/admin1/admin1.tls.crt"
)
var (
	node1 *sdk.NodeConfig
	node2 *sdk.NodeConfig
	cnt int64 = 0
)

func main()  {
	client, err := createClientWithCertBytes()
	admin1, err := createAdmin(orgId1)
	admin2, err := createAdmin(orgId2)
	admin3, err := createAdmin(orgId3)
	admin4, err := createAdmin(orgId4)
	//var cnt int = 0
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("====================== CreateUserContract ======================")
	CreateUserContract(client, admin1, admin2, admin3, admin4, voteContractName,voteVersion,voteByteCodePath,true)
	
	fmt.Println("====================== Get Items Count ======================")
	GetVotingItemsCount(client, voteContractName, "getItemsCount")

	//fmt.Println("====================== Getting User ======================")
	//GetUser(client,"78329178923178912hdiojd21i9d2j19id21ji")

	fmt.Println("====================== Getting Items ======================")
	GetAllItems(client)

	fmt.Println("====================== Vote for Item ======================")
	VoteForItem(client, voteContractName, "vote", 1)
	
	fmt.Println("====================== Getting Items ======================")
	GetAllItems(client)
}




func GetVotingItemsCount(client *sdk.ChainClient, contractName, 
	_method string)(*common.TxResponse, error) {
	
	abiJson, err := ioutil.ReadFile(voteABIPath)
	myAbi, err := abi.JSON(strings.NewReader(string(abiJson)))
	dataByte, err := myAbi.Pack(_method)
	dataString := hex.EncodeToString(dataByte)
	method := dataString[0:8]
	pairs := map[string]string{
		"data": dataString,
	}

	resp, err := client.QueryContract(contractName, method, pairs, -1)
	//fmt.Printf("resp: %s\n", resp)
	//fmt.Printf("\n GetContractResult: %s\n", resp.ContractResult)
	//fmt.Printf(fmt.Sprintf("%T", resp.GetContractResult()))

	val, err := myAbi.Unpack(_method,  resp.ContractResult.Result)
	fmt.Printf("Items count = %d\n", val[0])
	num := val[0].(*big.Int)
	cnt = num.Int64()
	if err=checkProposalRequestResp(resp, true);err!=nil{
		return nil,err
	}

	return resp,nil
}
func GetAllItems(client *sdk.ChainClient)() {
		var i int64
		for i = 0; i < cnt; i++{
			GetItem(client, i)
		  }
		  return
	}

func GetItem(client *sdk.ChainClient, 
	ID int64)(*common.TxResponse, error) {
	
	abiJson, err := ioutil.ReadFile(voteABIPath)
	myAbi, err := abi.JSON(strings.NewReader(string(abiJson)))
	dataByte, err := myAbi.Pack("items",  big.NewInt(ID))
	dataString := hex.EncodeToString(dataByte)
	method := dataString[0:8]
	pairs := map[string]string{
		"data": dataString,
	}

	resp, err := client.QueryContract(voteContractName, method, pairs, -1)
	//fmt.Printf("resp: %s\n", resp)
	//fmt.Printf("\n GetContractResult: %s\n", resp.ContractResult)
	//fmt.Printf(fmt.Sprintf("%T", resp.GetContractResult()))

	val, err := myAbi.Unpack("items",  resp.ContractResult.Result)
	fmt.Printf("Name: %s, ID: %s, Number of votes: %s\n", val[0], val[1], val[2])

	
	if err=checkProposalRequestResp(resp, true);err!=nil{
		return nil,err
	}

	return resp,nil
}

func GetUser(client *sdk.ChainClient, 
	 addr string)(*common.TxResponse, error) {
	
	abiJson, err := ioutil.ReadFile(voteABIPath)
	myAbi, err := abi.JSON(strings.NewReader(string(abiJson)))
	dataByte, err := myAbi.Pack("voters",  "addr")
	dataString := hex.EncodeToString(dataByte)
	method := dataString[0:8]
	pairs := map[string]string{
		"data": dataString,
	}
	
	resp, err := client.QueryContract(voteContractName, method, pairs, -1)
	//fmt.Printf("resp: %s\n", resp)
	//fmt.Printf("\n GetContractResult: %s\n", resp.ContractResult)
	//fmt.Printf(fmt.Sprintf("%T", resp.GetContractResult()))

	val, err := myAbi.Unpack("voters",  resp.ContractResult.Result)
	fmt.Printf("voted: %s, vote: %s \n", val[0], val[1])

	
	if err=checkProposalRequestResp(resp, true);err!=nil{
		return nil,err
	}

	return resp,nil
}
func VoteForItem(client *sdk.ChainClient, contractName, 
	_method string, ID int64)(*common.TxResponse, error) {
	
	abiJson, err := ioutil.ReadFile(voteABIPath)
	myAbi, err := abi.JSON(strings.NewReader(string(abiJson)))
	dataByte, err := myAbi.Pack(_method, big.NewInt(ID))

	dataString := hex.EncodeToString(dataByte)
	method := dataString[0:8]
	pairs := map[string]string{
		"data": dataString,
	}

	resp, err := client.InvokeContract(contractName, method, "", pairs, -1, true)
	if err != nil {
		return nil,err
	}

     if err=checkProposalRequestResp(resp,true);err!=nil{
		//fmt.Printf("resp: %s\n", resp)

     	return nil,err
	 }
	//fmt.Printf("resp: %s\n", resp)

	val, err := myAbi.Unpack(_method,  resp.ContractResult.Result)
	fmt.Printf("%s\n", val[0])

	if err=checkProposalRequestResp(resp, true);err!=nil{
		return nil,err
	}

	return resp,nil
}


func createClientWithCertBytes() (*sdk.ChainClient, error) {

	userCrtBytes, err := ioutil.ReadFile(fmt.Sprintf(userCrtPath, orgId1))
	if err != nil {
		return nil, err
	}

	userKeyBytes, err := ioutil.ReadFile(fmt.Sprintf(userKeyPath, orgId1))
	if err != nil {
		return nil, err
	}

	userSignCrtBytes, err := ioutil.ReadFile(fmt.Sprintf(userSignCrtPath, orgId1))
	if err != nil {
		return nil, err
	}

	userSignKeyBytes, err := ioutil.ReadFile(fmt.Sprintf(userSignKeyPath, orgId1))
	if err != nil {
		return nil, err
	}

	chainClient, err := sdk.NewChainClient(
		sdk.WithConfPath("./testdata/sdk_config.yml"),
		sdk.WithUserCrtBytes(userCrtBytes),
		sdk.WithUserKeyBytes(userKeyBytes),
		sdk.WithUserSignKeyBytes(userSignKeyBytes),
		sdk.WithUserSignCrtBytes(userSignCrtBytes),
	)

	if err != nil {
		return nil, err
	}

	err = chainClient.EnableCertHash()
	if err != nil {
		return nil, err
	}

	return chainClient, nil
}

func createAdmin(orgId string) (*sdk.ChainClient, error) {
	if node1 == nil {
		node1 = createNode(nodeAddr1, connCnt1)
	}

	if node2 == nil {
		node2 = createNode(nodeAddr2, connCnt2)
	}

	adminClient, err := sdk.NewChainClient(
		sdk.WithChainClientOrgId(orgId),
		sdk.WithChainClientChainId(chainId),
		sdk.WithUserKeyFilePath(fmt.Sprintf(adminKeyPath, orgId)),
		sdk.WithUserCrtFilePath(fmt.Sprintf(adminCrtPath, orgId)),
		sdk.AddChainClientNodeConfig(node1),
		sdk.AddChainClientNodeConfig(node2),
	)
	if err != nil {
		return nil, err
	}

	err = adminClient.EnableCertHash()
	if err != nil {
		return nil, err
	}

	return adminClient, nil
}


func createNode(nodeAddr string, connCnt int) *sdk.NodeConfig {
	node := sdk.NewNodeConfig(
		// 节点地址，格式：127.0.0.1:12301
		sdk.WithNodeAddr(nodeAddr),
		// 节点连接数
		sdk.WithNodeConnCnt(connCnt),
		// 节点是否启用TLS认证
		sdk.WithNodeUseTLS(true),
		// 根证书路径，支持多个
		sdk.WithNodeCAPaths(caPaths),
		// TLS Hostname
		sdk.WithNodeTLSHostName(tlsHostName),
	)

	return node
}

func CreateUserContract(client *sdk.ChainClient,
	admin1, admin2, admin3, admin4 *sdk.ChainClient,contractName, version, byteCodePath string, withSyncResult bool) (*common.TxResponse, error) {
	codeBytes, err := ioutil.ReadFile(voteByteCodePath)

	resp, err := createUserContract(client, admin1, admin2, admin3, admin4,
		contractName, version, string(codeBytes), common.RuntimeType_EVM, nil, withSyncResult)
		
	fmt.Printf("CREATE EVM Vote contract resp: %s\n", resp)

	return resp,err
}

func createUserContract(client *sdk.ChainClient, admin1, admin2, admin3, admin4 *sdk.ChainClient,
	contractName, version, byteCodePath string, runtime common.RuntimeType, kvs []*common.KeyValuePair, withSyncResult bool) (*common.TxResponse, error) {

	payloadBytes, err := client.CreateContractCreatePayload(contractName, version, byteCodePath, runtime, kvs)
	if err != nil {
		return nil, err
	}

	// 各组织Admin权限用户签名
	signedPayloadBytes1, err := admin1.SignContractManagePayload(payloadBytes)
	if err != nil {
		return nil, err
	}

	signedPayloadBytes2, err := admin2.SignContractManagePayload(payloadBytes)
	if err != nil {
		return nil, err
	}

	signedPayloadBytes3, err := admin3.SignContractManagePayload(payloadBytes)
	if err != nil {
		return nil, err
	}

	signedPayloadBytes4, err := admin4.SignContractManagePayload(payloadBytes)
	if err != nil {
		return nil, err
	}

	// 收集并合并签名
	mergeSignedPayloadBytes, err := client.MergeContractManageSignedPayload([][]byte{signedPayloadBytes1,
		signedPayloadBytes2, signedPayloadBytes3, signedPayloadBytes4})
	if err != nil {
		return nil, err
	}

	// 发送创建合约请求
	resp, err := client.SendContractManageRequest(mergeSignedPayloadBytes, createContractTimeout, withSyncResult)
	if err != nil {
		return nil, err
	}

	err = checkProposalRequestResp(resp, true)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func checkProposalRequestResp(resp *common.TxResponse, needContractResult bool) error {
	if resp.Code != common.TxStatusCode_SUCCESS {
		return errors.New(resp.Message)
	}

	if needContractResult && resp.ContractResult == nil {
		return fmt.Errorf("contract result is nil")
	}

	if resp.ContractResult != nil && resp.ContractResult.Code != common.ContractResultCode_OK {
		return errors.New(resp.ContractResult.Message)
	}
	return nil
}

