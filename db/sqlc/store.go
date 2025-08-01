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
	FromUsername string `json:"from_username"`
	ToUsername   string `json:"to_username"`
	Currency     string `json:"currency"`
	Type         string `json:"type"`
	Amount       int64  `json:"amount"`
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
		args := GetAccountByOwnerCurrencyTypeParams{
			Owner:    arg.FromUsername,
			Currency: arg.Currency,
			Type:     arg.Type,
		}
		//fmt.Println(txName, "create Transfer")
		FromAccount, err := store.GetAccountByOwnerCurrencyType(ctx, args)
		if err == sql.ErrNoRows {
			// Only create if not found
			FromAccount, err = store.CreateAccount(ctx, CreateAccountParams{
				Owner:    arg.FromUsername,
				Currency: arg.Currency,
				Type:     args.Type,
			})
			if err != nil {
				return fmt.Errorf("create from account: %w", err)
			}
		} else if err != nil {
			return fmt.Errorf("get account: %w", err)
		}

		args = GetAccountByOwnerCurrencyTypeParams{
			Owner:    arg.ToUsername,
			Currency: arg.Currency,
			Type:     arg.Type,
		}
		ToAccount, err := store.GetAccountByOwnerCurrencyType(ctx, args)
		if err == sql.ErrNoRows {
			// Only create if not found
			ToAccount, err = store.CreateAccount(ctx, CreateAccountParams{
				Owner:    arg.ToUsername,
				Currency: arg.Currency,
				Type:     args.Type,
			})
			if err != nil {
				return fmt.Errorf("create to account: %w", err)
			}
		} else if err != nil {
			return fmt.Errorf("get account: %w", err)
		}

		FromAccountID := FromAccount.ID
		ToAccountID := ToAccount.ID

		result.Transfer, err = store.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: FromAccountID,
			ToAccountID:   ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return fmt.Errorf("create transfer: %w", err)
		}

		//fmt.Println(txName, "create FromEntry")
		result.FromEntry, err = store.CreateEntry(ctx, CreateEntryParams{
			AccountID: FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return fmt.Errorf("create from entry: %w", err)
		}

		//fmt.Println(txName, "ToEntry")
		result.ToEntry, err = store.CreateEntry(ctx, CreateEntryParams{
			AccountID: ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return fmt.Errorf("create to entry: %w", err)
		}

		//fmt.Println(txName, "update account1")
		account1, err := store.GetAccountForUpdate(ctx, FromAccountID)
		if err != nil {
			return fmt.Errorf("get account for update1: %w", err)
		}

		//fmt.Println(txName, "update account2")
		account2, err := store.GetAccountForUpdate(ctx, ToAccountID)
		if err != nil {
			return fmt.Errorf("get account for update2: %w", err)
		}

		if FromAccountID < ToAccountID {
			result.FromAccount, err = store.UpdateAcount(ctx, UpdateAcountParams{
				ID:      FromAccountID,
				Balance: account1.Balance - arg.Amount,
			})
			if err != nil {
				return err
			}

			result.ToAccount, err = store.UpdateAcount(ctx, UpdateAcountParams{
				ID:      ToAccountID,
				Balance: account2.Balance + arg.Amount,
			})

		} else {
			result.ToAccount, err = store.UpdateAcount(ctx, UpdateAcountParams{
				ID:      ToAccountID,
				Balance: account2.Balance + arg.Amount,
			})
			if err != nil {
				return err
			}

			result.FromAccount, err = store.UpdateAcount(ctx, UpdateAcountParams{
				ID:      FromAccountID,
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

		args := GetAccountByOwnerCurrencyTypeGroupIDParams{
			Owner:    arg.Username,
			Currency: arg.Currency,
			Type:     arg.Type,
			GroupID:  sql.NullInt64{Int64: result.Group.ID, Valid: true},
		}

		ToAccount, err := store.GetAccountByOwnerCurrencyTypeGroupID(ctx, args)
		if err == sql.ErrNoRows {
			// Only create if not found
			ToAccount, err = store.CreateAccountWithGroup(ctx, CreateAccountWithGroupParams{
				Owner:    arg.Username,
				Balance:  0,
				Currency: arg.Currency,
				Type:     args.Type,
				GroupID:  sql.NullInt64{Int64: result.Group.ID, Valid: true},
			})
			if err != nil {
				return fmt.Errorf("create to account: %w", err)
			}
		} else if err != nil {
			return fmt.Errorf("get account: %w", err)
		}

		result.Account = ToAccount

		return nil
	})
	return result, err
}

type GroupTransferTxParams struct {
	FromUsername string `json:"from_username"`
	ToUsername   string `json:"to_username"`
	Currency     string `json:"currency"`
	Type         string `json:"type"`
	Amount       int64  `json:"amount"`
	GroupID      int64  `json:"group_id"`
}

func (store *Store) GroupTransactionTx(ctx context.Context, arg GroupTransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(*Queries) error {
		var err error

		//txName := ctx.Value(txKey)
		args := GetAccountByOwnerCurrencyTypeGroupIDParams{
			Owner:    arg.FromUsername,
			Currency: arg.Currency,
			Type:     arg.Type,
			GroupID:  sql.NullInt64{Int64: arg.GroupID, Valid: true},
		}
		//fmt.Println(txName, "create Transfer")
		FromAccount, err := store.GetAccountByOwnerCurrencyTypeGroupID(ctx, args)
		if err == sql.ErrNoRows {
			// Only create if not found
			FromAccount, err = store.CreateAccountWithGroup(ctx, CreateAccountWithGroupParams{
				Owner:    arg.FromUsername,
				Currency: arg.Currency,
				Type:     args.Type,
			})
			if err != nil {
				return fmt.Errorf("create from account: %w", err)
			}
		} else if err != nil {
			return fmt.Errorf("get account: %w", err)
		}

		args = GetAccountByOwnerCurrencyTypeGroupIDParams{
			Owner:    arg.ToUsername,
			Currency: arg.Currency,
			Type:     arg.Type,
			GroupID:  sql.NullInt64{Int64: arg.GroupID, Valid: true},
		}
		ToAccount, err := store.GetAccountByOwnerCurrencyTypeGroupID(ctx, args)
		if err == sql.ErrNoRows {
			// Only create if not found
			ToAccount, err = store.CreateAccountWithGroup(ctx, CreateAccountWithGroupParams{
				Owner:    arg.ToUsername,
				Currency: arg.Currency,
				Type:     args.Type,
			})
			if err != nil {
				return fmt.Errorf("create to account: %w", err)
			}
		} else if err != nil {
			return fmt.Errorf("get account: %w", err)
		}

		FromAccountID := FromAccount.ID
		ToAccountID := ToAccount.ID

		result.Transfer, err = store.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: FromAccountID,
			ToAccountID:   ToAccountID,
			Amount:        arg.Amount,
			GroupID:       sql.NullInt64{Int64: arg.GroupID, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("create transfer: %w", err)
		}

		//fmt.Println(txName, "create FromEntry")
		result.FromEntry, err = store.CreateEntry(ctx, CreateEntryParams{
			AccountID: FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return fmt.Errorf("create from entry: %w", err)
		}

		//fmt.Println(txName, "ToEntry")
		result.ToEntry, err = store.CreateEntry(ctx, CreateEntryParams{
			AccountID: ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return fmt.Errorf("create to entry: %w", err)
		}

		//fmt.Println(txName, "update account1")
		account1, err := store.GetAccountForUpdate(ctx, FromAccountID)
		if err != nil {
			return fmt.Errorf("get account for update1: %w", err)
		}

		//fmt.Println(txName, "update account2")
		account2, err := store.GetAccountForUpdate(ctx, ToAccountID)
		if err != nil {
			return fmt.Errorf("get account for update2: %w", err)
		}

		if FromAccountID < ToAccountID {
			result.FromAccount, err = store.UpdateAcount(ctx, UpdateAcountParams{
				ID:      FromAccountID,
				Balance: account1.Balance - arg.Amount,
			})
			if err != nil {
				return err
			}

			result.ToAccount, err = store.UpdateAcount(ctx, UpdateAcountParams{
				ID:      ToAccountID,
				Balance: account2.Balance + arg.Amount,
			})

		} else {
			result.ToAccount, err = store.UpdateAcount(ctx, UpdateAcountParams{
				ID:      ToAccountID,
				Balance: account2.Balance + arg.Amount,
			})
			if err != nil {
				return err
			}

			result.FromAccount, err = store.UpdateAcount(ctx, UpdateAcountParams{
				ID:      FromAccountID,
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
