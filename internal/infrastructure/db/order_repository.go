package db

import (
	"context"
	"database/sql"

	"github.com/dprio/redis-limiter/internal/domain"
)

type (
	OrderRepository interface {
		Save(ctx context.Context, order *domain.Order) (domain.Order, error)
		GetAll(ctx context.Context) ([]domain.Order, error)
	}
)

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (repo *orderRepository) Save(ctx context.Context, order *domain.Order) (domain.Order, error) {

	stmt, err := repo.db.PrepareContext(ctx, "INSERT INTO orders (id, price, tax, final_price) VALUES (?,?,?,?)")
	if err != nil {
		return domain.Order{}, err
	}

	_, err = stmt.ExecContext(ctx, order.ID, order.Price, order.Tax, order.FinalPrice)
	if err != nil {
		return domain.Order{}, err
	}

	return *order, nil
}

func (repo *orderRepository) GetAll(ctx context.Context) ([]domain.Order, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT * FROM orders")
	if err != nil {
		return []domain.Order{}, err
	}
	defer rows.Close()

	orders := make([]domain.Order, 0)
	for rows.Next() {
		var order domain.Order
		if err := rows.Scan(&order.ID, &order.Price, &order.Tax, &order.FinalPrice); err != nil {
			return []domain.Order{}, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}
