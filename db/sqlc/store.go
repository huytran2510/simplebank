package db

import (
    "database/sql"
    "context"
	"fmt"
)

// Stores provives all functions to execute db queries
type Store struct {
	*Queries 
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
        Queries: New(db),
        db: db,
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
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v , tx rollback error: %v", err, rbErr)
        }
		return err
	}
	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountId int64 `json:"from_account_id`
	ToAccountId int64 `json:"to_account_id`
    Amount int64 `json:"amount`
}


type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// Will TransferTx function cause a deadlock?
    // Yes, it will cause a deadlock.
    // The reason is that we are using the same transaction for both queries.
    // The first query is creating a transfer, and the second query is creating an entry.
    // If the first query is executed first, then the second query will wait for the first query to finish.
    // If the second query is executed first, then the first query will wait for the second query to finish.
    // The problem is that both queries are using the same transaction.
func (store *Store) TransferTx(ctx context.Context,arg TransferTxParams) (TransferTxResult,error) {
	var result TransferTxResult
	err := store.execTx(ctx,func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountId,
			ToAccountID:   arg.ToAccountId,
			Amount:        arg.Amount,
		})

		if err!= nil {
            return err
        }

		fromEntry, err := q.CreateEntry(ctx, CreateEntryParams{
            AccountID: arg.FromAccountId,
            Amount:    -arg.Amount, // Negative amount for the FromEntry
        })
        if err != nil {
            return err
        }
        result.FromEntry = fromEntry

		
		toEntry, err := q.CreateEntry(ctx, CreateEntryParams{
            AccountID: arg.ToAccountId,
            Amount:    arg.Amount, // Positive amount for the ToEntry
        })
        if err != nil {
            return err
        }
        result.ToEntry = toEntry

		if arg.FromAccountId < arg.ToAccountId {
			result.FromAccount,result.ToAccount, err = addMoney(ctx,q,arg.FromAccountId,-arg.Amount,arg.ToAccountId,arg.Amount) 
		} else {
			result.ToAccount,result.ToAccount, err = addMoney(ctx,q,arg.FromAccountId,arg.Amount,arg.ToAccountId,-arg.Amount) 
		}
		return err
	})

	return result,err
}

func (store *Store) TransferTx2(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// Tạo bản ghi chuyển khoản
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountId,
			ToAccountID:   arg.ToAccountId,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		// Khóa tài khoản theo thứ tự ID để tránh deadlock
		fmt.Println("arg.FromAccountId", arg.FromAccountId)
		fmt.Println("arg.ToAccountId", arg.ToAccountId)
		if arg.FromAccountId < arg.ToAccountId {
			err = lockAccounts(ctx, q, arg.FromAccountId, arg.ToAccountId)
		} else {
			err = lockAccounts(ctx, q, arg.ToAccountId, arg.FromAccountId)
		}
		fmt.Print("err",err)
		if err != nil {
			return err
		}

		// Lấy thông tin tài khoản gốc và cập nhật số dư
		account1, err := q.GetAccount(ctx, int32(arg.FromAccountId))
		if err != nil {
			return err
		}

		result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      int32(arg.FromAccountId),
			Balance: account1.Balance - int32(arg.Amount),
		})
		if err != nil {
			return err
		}

		// Lấy thông tin tài khoản đích và cập nhật số dư
		account2, err := q.GetAccount(ctx, int32(arg.ToAccountId))
		if err != nil {
			return err
		}

		result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      int32(arg.ToAccountId),
			Balance: account2.Balance + int32(arg.Amount),
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}

func lockAccounts(ctx context.Context, q *Queries, accountID1, accountID2 int64) error {
	_, err := q.db.ExecContext(ctx, `SELECT id FROM ACCOUNT WHERE id IN ($1, $2) FOR UPDATE`, accountID1, accountID2)
	return err
}



func addMoney(ctx context.Context, q *Queries, accountId1 int64, amount1 int64, accountId2 int64, amount2 int64) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     int32(accountId1),
		Amount: int32(amount1),
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     int32(accountId2),
		Amount: int32(amount2),
	})
	if err != nil {
		return
	}
	return
}

func (store *Store) TransferTxDeadlock(ctx context.Context,arg TransferTxParams) (TransferTxResult,error) {
	var result TransferTxResult
	err := store.execTx(ctx,func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountId,
			ToAccountID:   arg.ToAccountId,
			Amount:        arg.Amount,
		})

		if err!= nil {
            return err
        }

		fromEntry, err := q.CreateEntry(ctx, CreateEntryParams{
            AccountID: arg.FromAccountId,
            Amount:    -arg.Amount, // Negative amount for the FromEntry
        })
        if err != nil {
            return err
        }
        result.FromEntry = fromEntry

		
		toEntry, err := q.CreateEntry(ctx, CreateEntryParams{
            AccountID: arg.ToAccountId,
            Amount:    arg.Amount, // Positive amount for the ToEntry
        })
        if err != nil {
            return err
        }
        result.ToEntry = toEntry

		if arg.FromAccountId < arg.ToAccountId {
			result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:      int32(arg.FromAccountId),
                Amount: -int32(arg.Amount),
			})

			if err != nil {
				return err
			}

			result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
                ID:      int32(arg.ToAccountId),
                Amount: int32(arg.Amount),
            })

			if err != nil {
				return err
			}

		} else {
			result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
                ID:      int32(arg.ToAccountId),
                Amount: int32(arg.Amount),
            })

			if err != nil {
				return err
			}

			result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:      int32(arg.FromAccountId),
                Amount: -int32(arg.Amount),
			})

			if err != nil {
				return err
			}
		}
		return nil
	})

	return result,err
}