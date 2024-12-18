package shopify

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
)


type Repository interface {
	Close() 
	PutOrder(ctx context.Context, order Order) error
	GetOrdersForShopAndAccount(ctx context.Context, shopName string, accountId string) ([]Order, error)
	SyncOrders(ctx context.Context, shopName string, sinceId string, limit int, token string) ([]Order, error)
	UpdateOrder(ctx context.Context, order Order, accountId string, shopName string) error	
	StoreToken(ctx context.Context, shopName string, accountId string, token string) error

}

type postgresRepository struct{
	db *sql.DB
}


func NewPostgresRepository(url string) (*postgresRepository, error) {
	db, err:= sql.Open("postgres", url)
	if err!= nil{
		return nil, err
	}

	err = db.Ping()
	if err!= nil{
		return nil, err
	}
	return &postgresRepository{db: db}, nil
}

func (r *postgresRepository) PutOrder(ctx context.Context, order Order) error{

	tx, err:= r.db.BeginTx(ctx, nil)
	if err!= nil{
		return err
	}
	defer func(){
		if err!= nil{
			tx.Rollback()
			return 
		}
		err = tx.Commit()

	}()
	tx.ExecContext(
		ctx,
		`INSERT INTO orders (id, created_at, updated_at, shop_name, account_id, order_id, total_price) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		order.ID, order.CreatedAt, order.UpdatedAt, order.ShopName, order.AccountId, order.OrderId, order.TotalPrice,
	)
	if err!= nil{
		return err
	}

	return nil


}


func (r *postgresRepository) GetOrdersForShopAndAccount(ctx context.Context, shopName string, accountId string) ([]Order, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, created_at, updated_at, shop_name, account_id, order_id, total_price 
		FROM orders WHERE shop_name = $1 AND account_id = $2`,
		shopName, accountId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		if err := rows.Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt, &order.ShopName, &order.AccountId, &order.OrderId, &order.TotalPrice); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}


func (r *postgresRepository) SyncOrders(ctx context.Context, shopName string, sinceId string, limit int, token string) ([]Order, error) {
	// Placeholder for the actual syncing logic with Shopify or other external systems.
	// Here we query orders based on the `sinceId` and `limit` to fetch a subset of orders.

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, created_at, updated_at, shop_name, account_id, order_id, total_price 
		FROM orders WHERE shop_name = $1 AND id > $2 LIMIT $3`,
		shopName, sinceId, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		if err := rows.Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt, &order.ShopName, &order.AccountId, &order.OrderId, &order.TotalPrice); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}


func (r *postgresRepository) UpdateOrder(ctx context.Context, order Order, accountId string, shopName string) error {
	// Update order details like `updated_at` and `total_price`.
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE orders SET updated_at = $1, total_price = $2 WHERE order_id = $3 AND account_id = $4 AND shop_name = $5`,
		order.UpdatedAt, order.TotalPrice, order.OrderId, accountId, shopName,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *postgresRepository) StoreToken(ctx context.Context, shopName string, accountId string, token string) error {
	// Store or update the token for the shop and account.
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO tokens (shop_name, account_id, token) 
		VALUES ($1, $2, $3) 
		ON CONFLICT (shop_name, account_id) 
		DO UPDATE SET token = $3`,
		shopName, accountId, token,
	)
	if err != nil {
		return err
	}
	return nil
}


func (r *postgresRepository) Close() {
	r.db.Close()
}
