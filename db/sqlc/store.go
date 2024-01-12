package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	Store interface {
		Querier
		TransferTx(ctx context.Context, arg TransferParams) (res TransferRes, err error)
	}

	StorePostgresSql struct {
		db *pgxpool.Pool
		*Queries
	}

	//

	//StoreMysqlSql struct {
	//	db *pgxpool.Pool
	//	*Queries
	//}

	TransferParams struct {
		FromAccountId int64 `json:"from_account_id"`
		ToAccountId   int64 `json:"to_account_id"`
		Amount        int64 `json:"amount"`
	}

	TransferRes struct {
		Transfer    Transfer `json:"transfer"`
		FromAccount Account  `json:"from_account"`
		ToAccount   Account  `json:"to_account"`
		FromEntry   Entry    `json:"from_entry"`
		ToEntry     Entry    `json:"to_entry"`
	}
)

func NewStore(db *pgxpool.Pool) Store {
	return &StorePostgresSql{
		Queries: New(db),
		db:      db,
	}
}

func (store *StorePostgresSql) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.Begin(ctx)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit(ctx)
}

func (store *StorePostgresSql) TransferTx(ctx context.Context, arg TransferParams) (res TransferRes, err error) {
	err = store.execTx(ctx, func(queries *Queries) error {
		res.Transfer, err = queries.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountId,
			ToAccountID:   arg.ToAccountId,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		res.FromEntry, err = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountId,
			Amount:    -arg.Amount,
		})

		if err != nil {
			return err
		}

		res.ToEntry, err = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountId,
			Amount:    arg.Amount,
		})

		if err != nil {
			return err
		}

		if arg.FromAccountId < arg.ToAccountId {
			res.FromAccount, res.ToAccount, err = store.addMoney(ctx, queries, arg.FromAccountId, -arg.Amount, arg.ToAccountId, arg.Amount)
		} else {
			res.ToAccount, res.FromAccount, err = store.addMoney(ctx, queries, arg.ToAccountId, arg.Amount, arg.FromAccountId, -arg.Amount)
		}

		return nil
	})
	return
}

func (store *StorePostgresSql) addMoney(
	ctx context.Context,
	q *Queries,
	accountID1, accountAmount1, accountID2, accountAmount2 int64,
) (account1, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: accountAmount1,
		ID:     accountID1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: accountAmount2,
		ID:     accountID2,
	})

	if err != nil {
		return
	}

	return
}
