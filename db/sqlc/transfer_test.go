package db

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"github.com/user2410/simplebank/util"
	"testing"
	"time"
)

func createRandomTransfer(t *testing.T, acc1, acc2 *Account) Transfer {
	args := CreateTransferParams{
		FromAccountID: acc1.ID,
		ToAccountID:   acc2.ID,
		Amount:        util.RandomMoney(),
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.Equal(t, transfer.FromAccountID, acc1.ID)
	require.Equal(t, transfer.ToAccountID, acc2.ID)
	require.Equal(t, transfer.Amount, args.Amount)
	require.NotZero(t, transfer.ID, transfer.CreatedAt)

	return transfer
}

func TestQueries_CreateTransfer(t *testing.T) {
	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)
	transfer := createRandomTransfer(t, &acc1, &acc2)

	// cleanup
	require.NoError(t, testQueries.DeleteTransfer(context.Background(), transfer.ID))
	require.NoError(t, testQueries.DeleteAccount(context.Background(), acc1.ID))
	require.NoError(t, testQueries.DeleteAccount(context.Background(), acc2.ID))

}

func TestQueries_GetTransfer(t *testing.T) {
	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)
	transfer1 := createRandomTransfer(t, &acc1, &acc2)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)
	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)

	// cleanup
	require.NoError(t, testQueries.DeleteTransfer(context.Background(), transfer1.ID))
	require.NoError(t, testQueries.DeleteAccount(context.Background(), acc1.ID))
	require.NoError(t, testQueries.DeleteAccount(context.Background(), acc2.ID))
}

func TestQueries_DeleteTransfer(t *testing.T) {
	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)
	transfer := createRandomTransfer(t, &acc1, &acc2)

	err := testQueries.DeleteTransfer(context.Background(), transfer.ID)
	require.NoError(t, err)

	transfer, err = testQueries.GetTransfer(context.Background(), transfer.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, transfer)

	// cleanup
	require.NoError(t, testQueries.DeleteAccount(context.Background(), acc1.ID))
	require.NoError(t, testQueries.DeleteAccount(context.Background(), acc2.ID))
}
