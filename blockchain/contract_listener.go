package blockchain

import (
	"context"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/defi-pool-share/dps-webapi/blockchain/events"
	"github.com/defi-pool-share/dps-webapi/storage"
)

const (
	loanCreatedEventSignature = "LoanCreated(address,uint256)"
)

func InitBlockchainListener() {
	client, err := ethclient.Dial(os.Getenv("ETH_NETWORK_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum client: %v", err)
	}
	_ = client
	log.Printf("Connected to the Ethereum network with success")

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
			if vLog.Topics[0] == crypto.Keccak256Hash([]byte(loanCreatedEventSignature)) {
				handleLoanCreatedEvent(vLog)
			}
		}
	}
}

func handleLoanCreatedEvent(vLog types.Log) {
	_from := common.BytesToAddress(vLog.Topics[1][12:])
	_loanIndex := new(big.Int).SetBytes(vLog.Topics[2].Bytes()).Int64()
	loan := &events.LoanCreatedEvent{
		Addr:      _from.Hex(),
		LoanIndex: _loanIndex,
	}

	log.Printf("New LoanCreated event received (from: %s, loanIndex: %d)", loan.Addr, loan.LoanIndex)

	storage.SaveLoan(loan)
}
