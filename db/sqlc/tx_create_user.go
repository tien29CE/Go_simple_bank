package db

import "context"

// CreateUserTxParams contains the input paramters of the transfer transaction
type CreateUserTxParams struct {
	CreateUserParams
	AfterCreate func(user User) error
}

// CreateUserTxResult is the result of the transfer transaction
type CreateUserTxResult struct {
	User User
}

// CreateUserTx performs a money from one account to the other.
// It creates a transfer record, add acount entries, and update acounts' balance withen a single database transaction
func (store *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := store.execTX(ctx, func(q *Queries) error {
		var err error

		result.User, err = q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}

		return arg.AfterCreate(result.User)
	})

	return result, err
}
