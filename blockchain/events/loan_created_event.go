package events

type LoanCreatedEvent struct {
	Addr      string
	LoanIndex int64
}
