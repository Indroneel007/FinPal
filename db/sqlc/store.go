package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

//var txKey = struct{}{}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("err : %v and rbErr : %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(*Queries) error {
		var err error

		//txName := ctx.Value(txKey)

		//fmt.Println(txName, "create Transfer")
		result.Transfer, err = store.CreateTransfer(ctx, CreateTransferParams(arg))
		if err != nil {
			return fmt.Errorf("create transfer: %w", err)
		}

		//fmt.Println(txName, "create FromEntry")
		result.FromEntry, err = store.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return fmt.Errorf("create from entry: %w", err)
		}

		//fmt.Println(txName, "ToEntry")
		result.ToEntry, err = store.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return fmt.Errorf("create to entry: %w", err)
		}

		//fmt.Println(txName, "update account1")
		account1, err := store.GetAccountForUpdate(ctx, arg.FromAccountID)
		if err != nil {
			return fmt.Errorf("get account for update1: %w", err)
		}

		//fmt.Println(txName, "update account2")
		account2, err := store.GetAccountForUpdate(ctx, arg.ToAccountID)
		if err != nil {
			return fmt.Errorf("get account for update2: %w", err)
		}

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, err = store.UpdateAcount(ctx, UpdateAcountParams{
				ID:      arg.FromAccountID,
				Balance: account1.Balance - arg.Amount,
			})
			if err != nil {
				return err
			}

			result.ToAccount, err = store.UpdateAcount(ctx, UpdateAcountParams{
				ID:      arg.ToAccountID,
				Balance: account2.Balance + arg.Amount,
			})

		} else {
			result.ToAccount, err = store.UpdateAcount(ctx, UpdateAcountParams{
				ID:      arg.ToAccountID,
				Balance: account2.Balance + arg.Amount,
			})
			if err != nil {
				return err
			}

			result.FromAccount, err = store.UpdateAcount(ctx, UpdateAcountParams{
				ID:      arg.FromAccountID,
				Balance: account1.Balance - arg.Amount,
			})
			if err != nil {
				return err
			}
		}
		return err
	})
	return result, err
}

type CreateGroupTxParams struct {
	Username  string `json:"username"`
	GroupName string `json:"group_name"`
	Currency  string `json:"currency"`
	Type      string `json:"type"`
}

type CreateGroupTxResult struct {
	Group   Group   `json:"group"`
	Account Account `json:"account"`
}

func (store *Store) CreateGroupTx(ctx context.Context, arg CreateGroupTxParams) (CreateGroupTxResult, error) {
	var result CreateGroupTxResult
	err := store.execTx(ctx, func(*Queries) error {
		var err error

		result.Group, err = store.CreateGroup(ctx, CreateGroupParams{
			GroupName: arg.GroupName,
			Currency:  arg.Currency,
			Type:      arg.Type,
		})
		if err != nil {
			return fmt.Errorf("create group error: %w", err)
		}

		result.Account, err = store.CreateAccountWithGroup(ctx, CreateAccountWithGroupParams{
			Owner:    arg.Username,
			Balance:  0,
			Currency: arg.Currency,
			Type:     arg.Type,
			GroupID:  sql.NullInt64{Int64: result.Group.ID, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("create account with group error: %w", err)
		}

		return nil
	})
	return result, err
}
