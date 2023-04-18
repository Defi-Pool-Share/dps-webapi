package main

import (
	"log"

	"github.com/defi-pool-share/dps-webapi/blockchain"
	"github.com/defi-pool-share/dps-webapi/net"
	"github.com/defi-pool-share/dps-webapi/storage"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	storage.InitLocalStorage()
	go blockchain.InitBlockchainListener()
	net.InitAPI()
}
