package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lniklison/bill-splitter-go/internal/bills"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func NewPostgres(pool *pgxpool.Pool) *Postgres {
	return &Postgres{pool: pool}
}

func (r *Postgres) GetBill(ctx context.Context, billID int64) (bills.Bill, error) {
	var bill bills.Bill
	err := r.pool.QueryRow(ctx, `
		select id, description, total_cents
		from bills
		where id = $1
	`, billID).Scan(&bill.ID, &bill.Description, &bill.TotalAmount)
	if errors.Is(err, pgx.ErrNoRows) {
		return bills.Bill{}, bills.ErrBillNotFound
	}
	if err != nil {
		return bills.Bill{}, fmt.Errorf("load bill: %w", err)
	}

	shares, err := r.loadShares(ctx, billID)
	if err != nil {
		return bills.Bill{}, err
	}
	bill.Shares = shares
	return bill, nil
}

func (r *Postgres) ReplaceShares(ctx context.Context, billID int64, shares []bills.Share) (bill bills.Bill, returnErr error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return bills.Bill{}, fmt.Errorf("begin replacing bill shares: %w", err)
	}
	defer func() {
		if returnErr != nil {
			_ = tx.Rollback(context.WithoutCancel(ctx))
		}
	}()

	err = tx.QueryRow(ctx, `
		select id, description, total_cents
		from bills
		where id = $1
		for update
	`, billID).Scan(&bill.ID, &bill.Description, &bill.TotalAmount)
	if errors.Is(err, pgx.ErrNoRows) {
		return bills.Bill{}, bills.ErrBillNotFound
	}
	if err != nil {
		return bills.Bill{}, fmt.Errorf("lock bill: %w", err)
	}

	if _, err := tx.Exec(ctx, "delete from bill_shares where bill_id = $1", billID); err != nil {
		return bills.Bill{}, fmt.Errorf("delete bill shares: %w", err)
	}
	for _, share := range shares {
		if _, err := tx.Exec(ctx, `
			insert into bill_shares (
				bill_id, position, person_name, normalized_name, percentage_basis_points
			) values ($1, $2, $3, $4, $5)
		`, billID, share.Position, share.PersonName, share.NormalizedName, share.Percentage); err != nil {
			return bills.Bill{}, fmt.Errorf("insert bill share: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return bills.Bill{}, fmt.Errorf("commit bill shares: %w", err)
	}
	bill.Shares = shares
	return bill, nil
}

func (r *Postgres) loadShares(ctx context.Context, billID int64) ([]bills.Share, error) {
	rows, err := r.pool.Query(ctx, `
		select position, person_name, normalized_name, percentage_basis_points
		from bill_shares
		where bill_id = $1
		order by position
	`, billID)
	if err != nil {
		return nil, fmt.Errorf("load bill shares: %w", err)
	}
	defer rows.Close()

	var shares []bills.Share
	for rows.Next() {
		var share bills.Share
		if err := rows.Scan(
			&share.Position,
			&share.PersonName,
			&share.NormalizedName,
			&share.Percentage,
		); err != nil {
			return nil, fmt.Errorf("scan bill share: %w", err)
		}
		shares = append(shares, share)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read bill shares: %w", err)
	}
	return shares, nil
}
