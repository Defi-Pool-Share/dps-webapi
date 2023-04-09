package blockchain

import (
	"log"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	loanCreatedEventSignature = "LoanCreated (index_topic_1 address _from, index_topic_2 uint256 _loanIndex)"
)

func InitBlockchainListener() {
	client, err := ethclient.Dial(os.Getenv(("ETH_NETWORK_URL")))
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum client: %v", err)
	}
	_ = client
	log.Printf("Connected to the Ethereum network with success")
}
