package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	Store struct {
		*Queries
		db *pgxpool.Pool
	}

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

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
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

func (store *Store) TransferTx(ctx context.Context, arg TransferParams) (res TransferRes, err error) {
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

		res.FromAccount, err = queries.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.FromAccountId,
			Amount: -arg.Amount,
		})

		if err != nil {
			return err
		}

		res.ToAccount, err = queries.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.ToAccountId,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}
		return nil
	})
	return
}
