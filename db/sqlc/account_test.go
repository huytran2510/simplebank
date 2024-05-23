package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"simplebank/db/util"
)

func createRandomAccount(t *testing.T) Account {
	// user := createRandomUser(t)

	arg := CreateAccountParams{
		Owner:   util.RandomString(6) ,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)

	return account
}
