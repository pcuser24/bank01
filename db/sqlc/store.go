package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store : All functions of query struct + 1 function to execute money transfer
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// SQLStore provides all functionalities to run databases individually
// as well as their combination within transactions
type SQLStore struct {
	*Queries
	db *sql.DB
}

// NewStore create a new Store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(queries *Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	query := New(tx)
	err = fn(query)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			// rollback error
			return fmt.Errorf("tx err: %v\nrbErr:%v\n", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// Exported transactions

// TransferTxParams contains input parameters of the transaction account
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from-account"`
	ToAccount   Account  `json:"to-account"`
	FromEntry   Entry    `json:"from-entry"`
	ToEntry     Entry    `json:"to-entry"`
}

//var txKey = struct{}{}

// TransferTx perform a money transfer from one account to another
// Create a transfer record, add account entries, update accounts' balance within a single transaction
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		//txName := ctx.Value(txKey)

		// transfer record
		//fmt.Println(txName, "create transfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		// account entries
		//fmt.Println(txName, "create entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		//fmt.Println(txName, "create entry 2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// update accounts' balance
		// To avoid deadlock, update account with smaller id first
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addAmount(ctx, q, arg.FromAccountID, arg.ToAccountID, -arg.Amount, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addAmount(ctx, q, arg.ToAccountID, arg.FromAccountID, arg.Amount, -arg.Amount)
		}

		return err
	})

	return result, err
}

func addAmount(
	ctx context.Context,
	q *Queries,
	acc1ID, acc2ID int64,
	amount1, amount2 int64,
) (acc1, acc2 Account, err error) {
	acc1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: amount1,
		ID:     acc1ID,
	})
	if err != nil {
		return
	}
	acc2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: amount2,
		ID:     acc2ID,
	})
	return
}
