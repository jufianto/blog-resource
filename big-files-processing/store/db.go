package store

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ StoreInterface = (*Store)(nil)

type DataSales struct {
	ID            uuid.UUID
	Region        string
	Country       string
	ItemType      string
	SalesChannel  string
	OrderPriority string
	OrderDate     time.Time
	OrderID       string
	ShipDate      time.Time
	UnitSold      float64
	UnitPrice     float64
	UnitCost      float64
	TotalRevenue  float64
	TotalCost     float64
	TotalProfit   float64
}

type StoreInterface interface {
	InsertSales(ctx context.Context, sales DataSales) error
}

type Store struct {
	pool *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

func (s *Store) InsertSales(ctx context.Context, sales DataSales) error {

	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		log.Println("error get pool", err)
		return err
	}
	defer conn.Release()

	query := `INSERT into transactions 
		(id, region, country, item_type, sales_channel, 
		order_priority, order_date, order_id, ship_date, units_sold, 
		unit_price, unit_cost, total_revenue, total_cost, total_profit) 
			values ($1, $2, $3, $4, $5, $6, $7, $8, $9,$10, $11, $12, $13,$14,$15)`

	_, err = conn.Exec(
		ctx, query,
		sales.ID,
		sales.Region,
		sales.Country,
		sales.ItemType,
		sales.SalesChannel,
		sales.OrderPriority,
		sales.OrderDate,
		sales.OrderID,
		sales.ShipDate,
		sales.UnitSold,
		sales.UnitPrice,
		sales.UnitCost,
		sales.TotalRevenue,
		sales.TotalCost,
		sales.TotalProfit,
	)
	if err != nil {
		log.Println("error from db", err)
		return err
	}
	return nil
}
