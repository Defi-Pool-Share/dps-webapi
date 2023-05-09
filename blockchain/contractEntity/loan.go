package contractEntity

import (
	"github.com/ethereum/go-ethereum/common"
)

type Loan struct {
	Lender        common.Address
	Borrower      common.Address
	TokenId       int64
	LoanAmount    strin
	CreationTime  int64
	StartTime     int64
	EndTime       int64
	AcceptedToken common.Address
	IsActive      bool
	LoanIndex     int64
}
