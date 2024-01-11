package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransferTx(t *testing.T) {

	store := NewStore(testDb)

	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)

	//fmt.Println(">> before:", fromAccount.Balance, toAccount.Balance)
	chanErr := make(chan error)
	transferRes := make(chan TransferRes)
	exited := make(map[int]bool)

	amount := int64(10)
	n := 5
	//wg := sync.WaitGroup{}
	for i := 0; i < n; i++ {
		//wg.Add(1)
		go func() {
			//defer func() {
			//	wg.Done()
			//}()
			res, err := store.TransferTx(context.Background(), TransferParams{
				FromAccountId: fromAccount.ID,
				ToAccountId:   toAccount.ID,
				Amount:        amount,
			})
			chanErr <- err
			transferRes <- res
		}()

	}
	//wg.Wait()

	for i := 0; i < n; i++ {
		err := <-chanErr
		res := <-transferRes
		require.NoError(t, err)
		require.NotEmpty(t, res)

		transfer := res.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, transfer.FromAccountID, fromAccount.ID)
		require.Equal(t, transfer.ToAccountID, toAccount.ID)
		require.Equal(t, transfer.Amount, amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)
		_, err = testQueries.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		fromEntry := res.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromEntry.Amount, -amount)
		require.Equal(t, fromEntry.AccountID, fromAccount.ID)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)
		_, err = testQueries.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := res.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, toEntry.Amount, amount)
		require.Equal(t, toEntry.AccountID, toAccount.ID)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)
		_, err = testQueries.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		resFromAccount := res.FromAccount
		require.NotEmpty(t, resFromAccount)
		require.Equal(t, resFromAccount.ID, fromAccount.ID)

		resToAccount := res.ToAccount
		require.NotEmpty(t, resToAccount)
		require.Equal(t, resToAccount.ID, toAccount.ID)

		//fmt.Println(">> tx:", resFromAccount.Balance, resToAccount.Balance)
		diff1 := fromAccount.Balance - resFromAccount.Balance
		diff2 := resToAccount.Balance - toAccount.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := diff1 / amount
		require.True(t, k >= 1 && int(k) <= n)
		require.NotContains(t, exited, k)
		exited[i] = true
	}

	updateFromAccount, err := testQueries.GetAccount(context.Background(), fromAccount.ID)
	updateToAccount, err := testQueries.GetAccount(context.Background(), toAccount.ID)
	//fmt.Println(">> after:", updateFromAccount.Balance, updateToAccount.Balance)
	require.NoError(t, err)

	require.Equal(t, fromAccount.Balance-int64(n)*amount, updateFromAccount.Balance)
	require.Equal(t, toAccount.Balance+int64(n)*amount, updateToAccount.Balance)

}

func TestTransferTxDeadlock(t *testing.T) {

	store := NewStore(testDb)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	//fmt.Println(">> before:", account1.Balance, account2.Balance)
	chanErr := make(chan error)

	amount := int64(10)
	n := 10
	for i := 0; i < n; i++ {
		fromAccountId := account1.ID
		toAccountId := account2.ID
		if i%2 == 1 {
			fromAccountId = account2.ID
			toAccountId = account1.ID
		}

		go func() {
			_, err := store.TransferTx(context.Background(), TransferParams{
				FromAccountId: fromAccountId,
				ToAccountId:   toAccountId,
				Amount:        amount,
			})
			chanErr <- err
		}()

	}

	for i := 0; i < n; i++ {
		err := <-chanErr
		require.NoError(t, err)
	}

	updateFromAccount, err := testQueries.GetAccount(context.Background(), account1.ID)
	updateToAccount, err := testQueries.GetAccount(context.Background(), account2.ID)
	//fmt.Println(">> after:", updateFromAccount.Balance, updateToAccount.Balance)
	require.NoError(t, err)

	require.Equal(t, account1.Balance, updateFromAccount.Balance)
	require.Equal(t, account2.Balance, updateToAccount.Balance)

}
