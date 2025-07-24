package db

import (
	"context"
	//"database/sql"
	"examples/SimpleBankProject/util"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func CreateRandomAccount(t *testing.T, username string) Account {
	user := username

	args := CreateAccountParams{
		Owner:    user,
		Balance:  util.RandomAmount(),
		Currency: util.RandomCurrency(),
		Type:     util.RandomType(),
	}
	account, err := testQueries.CreateAccount(context.Background(), args)

	require.NoError(t, err)

	require.NotEmpty(t, account)
	require.Equal(t, args.Balance, account.Balance)
	require.Equal(t, args.Currency, account.Currency)
	require.Equal(t, args.Owner, account.Owner)
	require.Equal(t, args.Type, account.Type)
	require.NotEmpty(t, account.CreatedAt)

	require.NotZero(t, account.CreatedAt)
	require.NotZero(t, account.ID)
	fmt.Printf("Created account: %v\n", account.ID)

	return account
}

func CreateRandomAccountWithCurrencyType(t *testing.T, username, currency, accType string) Account {
	args := CreateAccountParams{
		Owner:    username,
		Balance:  util.RandomAmount(),
		Currency: currency,
		Type:     accType,
	}
	account, err := testQueries.CreateAccount(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	return account
}

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	user1 := CreateRandomUser(t)
	user2 := CreateRandomUser(t)

	accountA := CreateRandomAccount(t, user1.Username)
	accountB := CreateRandomAccountWithCurrencyType(t, user2.Username, accountA.Currency, accountA.Type)

	account1, err := testQueries.GetAccountByOwnerCurrencyType(context.Background(), GetAccountByOwnerCurrencyTypeParams{
		Owner:    user1.Username,
		Currency: accountA.Currency,
		Type:     accountA.Type,
	})
	require.NoError(t, err)
	require.NotEmpty(t, account1)
	require.Equal(t, accountA.ID, account1.ID)

	account2, err := testQueries.GetAccountByOwnerCurrencyType(context.Background(), GetAccountByOwnerCurrencyTypeParams{
		Owner:    user2.Username,
		Currency: accountB.Currency,
		Type:     accountB.Type,
	})
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, accountB.ID, account2.ID)

	amount := int64(10)
	n := 2
	//existed := make(map[int]bool)
	//performing n concurrent transfers
	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		//txName := fmt.Sprintf("tx %d", i+1)
		go func() {
			//ctx := context.WithValue(context.Background())
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromUsername: account1.Owner,
				ToUsername:   account2.Owner,
				Currency:     account1.Currency,
				Type:         account1.Type,
				Amount:       amount,
			})

			fmt.Printf("Transfer %d: %v\n", i+1, result.FromEntry.AccountID)

			errs <- err
			results <- result
		}()
	}

	//check results
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		fmt.Printf("Transfer ID: %d, From: %d\n", transfer.ID, transfer.FromAccountID)

		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		fmt.Printf("From Entry ID: %d, Account ID: %d\n", fromEntry.ID, fromEntry.AccountID)

		_, err = store.GetEntry(context.Background(), fromEntry.AccountID)
		require.NoError(t, err)

		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		//require.NotContains(t, existed, k)
		//existed[k+1] = true
	}

	updatedAccount1, err := testQueries.GetAccountForUpdate(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccountForUpdate(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)

}

/*func CreateRandomAccountWithGroup(t *testing.T) Account {
	user := CreateRandomUser(t)

	group := CreateRandomGroup(t)

	args := CreateAccountWithGroupParams{
		Owner:    user.Username,
		Balance:  util.RandomAmount(),
		Currency: util.RandomCurrency(),
		Type:     util.RandomType(),
		GroupID:  sql.NullInt64{Int64: group.ID, Valid: true},
	}
	account, err := testQueries.CreateAccountWithGroup(context.Background(), args)

	require.NoError(t, err)

	require.NotEmpty(t, account)
	require.Equal(t, args.Balance, account.Balance)
	require.Equal(t, args.Currency, account.Currency)
	require.Equal(t, args.Owner, account.Owner)
	require.Equal(t, args.Type, account.Type)
	require.Equal(t, args.GroupID, account.GroupID)
	require.NotEmpty(t, account.CreatedAt)

	require.NotZero(t, account.CreatedAt)
	require.NotZero(t, account.ID)
	fmt.Printf("Created account with group: %v\n", account.ID)

	return account
}

func CreateRandomGroup(t *testing.T) Group {
	store := NewStore(testDB)

	groupName := util.RandomGroupName()
	currency := util.RandomCurrency()
	groupType := util.RandomType()

	group, err := store.CreateGroup(context.Background(), CreateGroupParams{
		GroupName: groupName,
		Currency:  currency,
		Type:      groupType,
	})

	require.NoError(t, err)
	require.NotEmpty(t, group)
	require.Equal(t, groupName, group.GroupName)
	require.Equal(t, currency, group.Currency)
	require.Equal(t, groupType, group.Type)
	require.NotZero(t, group.ID)
	require.NotZero(t, group.CreatedAt)

	return group
}

func CreateGroupTx(t *testing.T) {
	store := NewStore(testDB)

	username := CreateRandomUser(t).Username
	groupName := util.RandomGroupName()
	currency := util.RandomCurrency()
	groupType := util.RandomType()

	result, err := store.CreateGroupTx(context.Background(), CreateGroupTxParams{
		Username:  username,
		GroupName: groupName,
		Currency:  currency,
		Type:      groupType,
	})

	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, groupName, result.Group.GroupName)
	require.Equal(t, currency, result.Group.Currency)
	require.Equal(t, groupType, result.Group.Type)
	require.NotZero(t, result.Group.ID)
	require.NotZero(t, result.Group.CreatedAt)

	require.Equal(t, username, result.Account.Owner)
	require.Equal(t, currency, result.Account.Currency)
	require.Equal(t, groupType, result.Account.Type)
	require.NotZero(t, result.Account.ID)
	require.NotZero(t, result.Account.CreatedAt)
}*/
