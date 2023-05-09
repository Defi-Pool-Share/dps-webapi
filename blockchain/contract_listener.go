package blockchain

import (
	"context"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"reflect"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/defi-pool-share/dps-webapi/blockchain/contractEntity"
	"github.com/defi-pool-share/dps-webapi/blockchain/events"
	"github.com/defi-pool-share/dps-webapi/storage"
)

const (
	loanCreatedEventSignature  = "LoanCreated(address,uint256)"
	loanUpdatedEventSignature  = "LoanUpdated(uint256)"
	loanWithdrawEventSignature = "LoanWithdraw(uint256)"
)

var dpsABI *abi.ABI
var client *ethclient.Client

func InitBlockchainListener() {
	var err error
	client, err = ethclient.Dial(os.Getenv("ETH_NETWORK_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum client: %v", err)
	}
	log.Printf("Connected to the Ethereum network with success")

	// Parse the JSON ABI
	abiFile, err := os.Open("./abi/DPSLendingUniswapLiquidity.json")
	if err != nil {
		log.Fatalf("Failed to open ABI file: %v", err)
	}
	defer abiFile.Close()

	abiBytes, err := ioutil.ReadAll(abiFile)
	if err != nil {
		log.Fatalf("Failed to read ABI file: %v", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(string(abiBytes)))
	if err != nil {
		log.Fatalf("Failed to to parse ABI file: %v", err)
	}
	dpsABI = &parsedABI

	// Setup listener
	contractAddress := common.HexToAddress(os.Getenv("DPS_CONTRACT_ADDR"))
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatalf("Failed to subscribe to logs: %v", err)
	}
	defer sub.Unsubscribe()

	for {
		select {
		case err := <-sub.Err():
			log.Fatalf("Error from logs subscription: %v", err)
		case vLog := <-logs:
			// created
			if vLog.Topics[0] == crypto.Keccak256Hash([]byte(loanCreatedEventSignature)) {
				handleLoanCreatedEvent(vLog)
			}
			// update
			if vLog.Topics[0] == crypto.Keccak256Hash([]byte(loanUpdatedEventSignature)) {
				handleLoanUpdatedEvent(vLog)
			}
			// widthdraw
			if vLog.Topics[0] == crypto.Keccak256Hash([]byte(loanWithdrawEventSignature)) {
				handleLoanWithdrawEvent(vLog)
			}
		}
	}
}

func handleLoanCreatedEvent(vLog types.Log) {
	_from := common.BytesToAddress(vLog.Topics[1][12:])
	_loanIndex := new(big.Int).SetBytes(vLog.Topics[2].Bytes()).Int64()
	loanEvent := &events.LoanCreatedEvent{
		Addr:      _from.Hex(),
		LoanIndex: _loanIndex,
	}

	log.Printf("New LoanCreated event received (from: %s, loanIndex: %d)", loanEvent.Addr, loanEvent.LoanIndex)

	loan, err := getLoanInfo(client, loanEvent.LoanIndex)
	if err != nil {
		log.Fatalf("%v", err)
	}

	storage.SaveLoan(loan)
}

func handleLoanUpdatedEvent(vLog types.Log) {
	_loanIndex := new(big.Int).SetBytes(vLog.Topics[1].Bytes()).Int64()

	log.Printf("New LoanUpdated event received (loanIndex: %d)", _loanIndex)

	loan, err := getLoanInfo(client, _loanIndex)
	if err != nil {
		log.Fatalf("%v", err)
	}

	storage.SaveLoan(loan)
}

func handleLoanWithdrawEvent(vLog types.Log) {
	_loanIndex := new(big.Int).SetBytes(vLog.Topics[1].Bytes()).Int64()

	log.Printf("New LoanWithdraw event received (loanIndex: %d)", _loanIndex)
	storage.DeleteLoanByIndex(_loanIndex)
}

func getLoanInfo(client *ethclient.Client, index int64) (*contractEntity.Loan, error) {
	loanInfo, err := dpsABI.Pack("getLoanInfo", big.NewInt(index))
	if err != nil {
		return nil, err
	}

	contractAddress := common.HexToAddress(os.Getenv("DPS_CONTRACT_ADDR"))
	callMsg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: loanInfo,
	}

	res, err := client.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		return nil, err
	}

	unpackedOutputs, err := dpsABI.Unpack("getLoanInfo", res)
	if err != nil {
		return nil, err
	}
	var loan contractEntity.Loan
	loanStruct := unpackedOutputs[0]
	loanValue := reflect.ValueOf(loanStruct)

	loan.Lender = loanValue.FieldByName("Lender").Interface().(common.Address)
	loan.Borrower = loanValue.FieldByName("Borrower").Interface().(common.Address)
	loan.TokenId = loanValue.FieldByName("TokenId").Interface().(*big.Int).Int64()
	loan.LoanAmount = loanValue.FieldByName("LoanAmount").Interface().(*big.Int).String()
	loan.CreationTime = loanValue.FieldByName("CreationTime").Interface().(*big.Int).Int64()
	loan.StartTime = loanValue.FieldByName("StartTime").Interface().(*big.Int).Int64()
	loan.EndTime = loanValue.FieldByName("EndTime").Interface().(*big.Int).Int64()
	loan.AcceptedToken = loanValue.FieldByName("AcceptedToken").Interface().(common.Address)
	loan.IsActive = loanValue.FieldByName("IsActive").Interface().(bool)
	loan.LoanIndex = loanValue.FieldByName("LoanIndex").Interface().(*big.Int).Int64()

	log.Printf("%v", loan)

	return &loan, nil
}
