package postgres

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/smugglerv1/pkg/orders"
)

type Repo struct {
	db *sql.DB
}

func New(repo string) (*Repo, error) {
	const op = "repo.postgres.new"
	db, err := sql.Open("postgres", repo)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	return &Repo{db: db}, nil
}

func (r *Repo) Stop() error {
	return r.db.Close()
}

func (r *Repo) SaveOrder(ctx context.Context, customerId int32, productId int32, quantity int32) (string, error) {

	const op = "repository.postgres.SaveOrder"
	sqlQuery := `INSERT INTO orders (customer_id, product_id, quantity) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, sqlQuery, customerId, productId, quantity)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	status := "Order has been created"
	return status, nil
}

func (r *Repo) GetOrder(ctx context.Context, customerId int32) ([]*orders.Order, error) {

	const op = "repository.postgres.GetOrder"
	queryRow := `SELECT  orderid ,customerid, productid, quantity FROM orders WHERE customerid =$1`
	rows, err := r.db.QueryContext(ctx, queryRow, customerId)
	if err != nil {
		return nil, err
	}
	var ords []*orders.Order

	for rows.Next() {
		var o *orders.Order
		err = rows.Scan(&o.OrderId, &o.CustomerId, &o.ProductId, &o.Quantity)
		ords = append(ords, o)
		if err != nil {
			return nil, err
		}
	}
	return ords, nil
}
