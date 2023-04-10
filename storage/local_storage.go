package storage

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/defi-pool-share/dps-webapi/blockchain/events"
	"github.com/dgraph-io/badger/v3"
)

var db *badger.DB

func InitLocalStorage() {
	var err error
	db, err = badger.Open(badger.DefaultOptions(os.Getenv("LOCAL_STORAGE_PATH")))
	if err != nil {
		log.Fatalf("Failed to open BadgerDB: %v", err)
	}
	defer db.Close()

}

func SaveLoan(loan *events.LoanCreatedEvent) {
	loanJSON, err := json.Marshal(loan)
	if err != nil {
		log.Fatalf("Failed to marshal loan: %v", err)
	}

	err = db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte("loan:"+strconv.FormatInt(loan.LoanIndex, 10)), loanJSON)
	})

	if err != nil {
		log.Fatalf("Failed to store loan in BadgerDB: %v", err)
	}
}

func FetchAllLoans() ([]*events.LoanCreatedEvent, error) {
	var loans []*events.LoanCreatedEvent

	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = []byte("loan:")
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			value, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}

			var loan events.LoanCreatedEvent
			err = json.Unmarshal(value, &loan)
			if err != nil {
				return err
			}
			loans = append(loans, &loan)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return loans, nil
}
