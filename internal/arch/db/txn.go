package db

// Txn defines the abstract interface for transactions.
type Txn interface {

	// Validate returns an error if the transaction is not valid.
	Validate() error

	// End ends the transaction.
	End()
}
