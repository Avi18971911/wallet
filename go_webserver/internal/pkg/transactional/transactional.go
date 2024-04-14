package transactional

import "context"

// Transactional defines the operations for starting, committing, and rolling back a transaction.
type Transactional interface {
	BeginTransaction(ctx context.Context) (TransactionContext, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// TransactionContext represents a context for a set of operations that should be executed within the same transaction.
// This could be as simple as a context.Context interface, or it could provide additional methods for transaction-specific operations.
type TransactionContext interface {
	context.Context
}
