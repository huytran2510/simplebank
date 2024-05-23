package db

import (
	"context"
	// "database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDb)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// Run n current transfer transactions
	n := 5
	amount := int64(10)
	errs := make(chan error, n)
	results := make(chan TransferTxResult, n)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountId: int64(account1.ID),
				ToAccountId:   int64(account2.ID),
				Amount:        int64(amount),
			})

			errs <- err
			results <- result
		}()
	}

	// Check results
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		fmt.Printf("results: %+v\n", results)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.Equal(t, account1.ID, int32(result.Transfer.FromAccountID))
		require.Equal(t, account2.ID, int32(result.Transfer.ToAccountID))
		require.Equal(t, amount, result.Transfer.Amount)
		require.NotZero(t, result.Transfer.ID)
		require.NotZero(t, result.Transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)
		fmt.Printf("fromTransfer: %+v\n", result.Transfer)

		//check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, int64(account1.ID), fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)
		fmt.Printf("fromEntry: %+v\n", result.FromEntry)

		_, err = testQueries.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, int64(account2.ID), toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = testQueries.GetEntry(context.Background(), toEntry.AccountID)
		require.NoError(t, err)
		fmt.Printf("toEntry: %+v\n", result.ToEntry)

		//check account
		fmt.Printf("fromAccount: %+v\n", result.FromAccount)
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// check balances
		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, int64(diff1)%amount == 0) // 1 * amount, 2 * amount, 3 * amount, ..., n * amount

		k := int(int64(diff1) / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}
	// check the final updated balance
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, int64(account1.Balance)-int64(n)*amount, int64(updatedAccount1.Balance))
	require.Equal(t, int64(account2.Balance)+int64(n)*amount, int64(updatedAccount2.Balance))
}

func TestTransferTx2(t *testing.T) {
	store := NewStore(testDb)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// Run n current transfer transactions
	n := 5
	amount := int64(10)
	errs := make(chan error, n)
	results := make(chan TransferTxResult, n)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx2(context.Background(), TransferTxParams{
				FromAccountId: int64(account1.ID),
				ToAccountId:   int64(account2.ID),
				Amount:        int64(amount),
			})

			errs <- err
			results <- result
		}()
	}

	// Check results
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		fmt.Printf("results: %+v\n", results)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.Equal(t, account1.ID, int32(result.Transfer.FromAccountID))
		require.Equal(t, account2.ID, int32(result.Transfer.ToAccountID))
		require.Equal(t, amount, result.Transfer.Amount)
		require.NotZero(t, result.Transfer.ID)
		require.NotZero(t, result.Transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)
		fmt.Printf("fromTransfer: %+v\n", result.Transfer)

		//check account
		fmt.Printf("fromAccount: %+v\n", result.FromAccount)
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// check balances
		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		fmt.Println(">> diff1:", diff1)
		fmt.Println(">> diff2:", diff2)

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, int64(diff1)%amount == 0) // 1 * amount, 2 * amount, 3 * amount, ..., n * amount

		k := int(int64(diff1) / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}
	// check the final updated balance
	// updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	// require.NoError(t, err)

	// updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	// require.NoError(t, err)

	// fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)

	// require.Equal(t, int64(account1.Balance)-int64(n)*amount, int64(updatedAccount1.Balance))
	// require.Equal(t, int64(account2.Balance)+int64(n)*amount, int64(updatedAccount2.Balance))
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDb)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// Run n current transfer transactions
	n := 10
	amount := int64(10)
	errs := make(chan error, n)
	// results := make(chan TransferTxResult, n)

	for i := 0; i < n; i++ {
		fromAccount := account1.ID
        toAccount := account2.ID

        if i % 2 == 0 {
            fromAccount, toAccount = toAccount, fromAccount
        }

		go func(fromAccount, toAccount int64) {
            _, err := store.TransferTxDeadlock(context.Background(), TransferTxParams{
                FromAccountId: int64(fromAccount),
                ToAccountId:   int64(toAccount),
                Amount:        amount,
            })
    
            errs <- err
        }(int64(fromAccount), int64(toAccount))
	}

	// Check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}
	// check the final updated balance
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)

	require.Equal(t, int64(account1.Balance), int64(updatedAccount1.Balance))
	require.Equal(t, int64(account2.Balance), int64(updatedAccount2.Balance))
}
