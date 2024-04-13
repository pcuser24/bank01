package db

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	_ "github.com/stretchr/testify/require"
	"github.com/user2410/simplebank/util"
	"testing"
	"time"
)

func createRandomEntry(t *testing.T, acc *Account) Entry {
	args := CreateEntryParams{
		AccountID: acc.ID,
		Amount:    util.RandomMoney(),
	}
	entry, err := testQueries.CreateEntry(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, entry.AccountID, args.AccountID)
	require.Equal(t, entry.Amount, args.Amount)
	require.NotZero(t, entry.ID, entry.CreatedAt)
	return entry
}

func TestQueries_CreateEntry(t *testing.T) {
	acc := createRandomAccount(t)
	entry := createRandomEntry(t, &acc)

	// cleanup
	require.NoError(t, testQueries.DeleteEntry(context.Background(), entry.ID))
	require.NoError(t, testQueries.DeleteAccount(context.Background(), acc.ID))
}

func TestQueries_GetEntry(t *testing.T) {
	acc := createRandomAccount(t)
	entry1 := createRandomEntry(t, &acc)

	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)

	// cleanup
	require.NoError(t, testQueries.DeleteEntry(context.Background(), entry1.ID))
	require.NoError(t, testQueries.DeleteAccount(context.Background(), acc.ID))
}

func TestQueries_DeleteEntry(t *testing.T) {
	acc := createRandomAccount(t)
	entry := createRandomEntry(t, &acc)

	err := testQueries.DeleteEntry(context.Background(), entry.ID)
	require.NoError(t, err)

	entry, err = testQueries.GetEntry(context.Background(), entry.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, entry)

	// cleanup
	require.NoError(t, testQueries.DeleteAccount(context.Background(), acc.ID))
}
