package order

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type Repository interface {
	Close()
	PutOrder(ctx context.Context, order Order) error
	GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (Repository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return &postgresRepository{db: db}, nil
}

func (r *postgresRepository) Close() {
	r.db.Close()
}

func (r *postgresRepository) PutOrder(ctx context.Context, order Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Insert order
	_, err = tx.ExecContext(
		ctx,
		"INSERT INTO orders (id, created_at, account_id, total_price) VALUES ($1, $2, $3, $4)",
		order.ID,
		order.CreatedAt,
		order.AccountID,
		order.TotalPrice,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Prepare COPY for order_products
	stmt, err := tx.Prepare(pq.CopyIn("order_products", "order_id", "product_id", "quantity"))
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, product := range order.Products {
		_, err = stmt.Exec(order.ID, product.ID, product.Quantity)
		if err != nil {
			stmt.Close()
			tx.Rollback()
			return err
		}
	}

	_, err = stmt.Exec() // finalize COPY
	if err != nil {
		stmt.Close()
		tx.Rollback()
		return err
	}

	stmt.Close()
	return tx.Commit()
}

func (r *postgresRepository) GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT o.id, o.created_at, o.account_id, o.total_price, op.product_id, op.quantity
		FROM orders o
		JOIN order_products op ON o.id = op.order_id
		WHERE o.account_id = $1
		ORDER BY o.created_at
		OFFSET $2 LIMIT $3`,
		accountID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ordersMap := make(map[string]*Order)

	for rows.Next() {
		var (
			orderID      string
			createdAt    time.Time
			accountIDRow string
			totalPrice   float64
			productID    string
			quantity     uint64
		)
		if err := rows.Scan(&orderID, &createdAt, &accountIDRow, &totalPrice, &productID, &quantity); err != nil {
			return nil, err
		}

		order, exists := ordersMap[orderID]
		if !exists {
			order = &Order{
				ID:         orderID,
				CreatedAt:  createdAt,
				AccountID:  accountIDRow,
				TotalPrice: totalPrice,
				Products:   []OrderProduct{},
			}
			ordersMap[orderID] = order
		}

		order.Products = append(order.Products, OrderProduct{
			ID: productID,
			Quantity:  quantity,
		})
	}

	orders := make([]Order, 0, len(ordersMap))
	for _, o := range ordersMap {
		orders = append(orders, *o)
	}

	return orders, nil
}
