package transactions

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
)

type NoTransactionMgr struct{}

func NoTransactions() *NoTransactionMgr {
	return &NoTransactionMgr{}
}

func (tm *NoTransactionMgr) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

type TransactionMgr struct {
	client *mongo.Client
}

func NewTransactionMgr(client *mongo.Client) *TransactionMgr {
	return &TransactionMgr{client: client}
}

func (tm *TransactionMgr) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return tm.client.UseSession(ctx, func(sc mongo.SessionContext) error {
		defer sc.EndSession(ctx)

		if err := sc.StartTransaction(); err != nil {
			return fmt.Errorf("failed to start transaction: %w", err)
		}

		if err := fn(sc); err != nil {
			return fmt.Errorf("failed in transation: %w", err)
		}

		if err := sc.CommitTransaction(ctx); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}

		return nil
	})
}
