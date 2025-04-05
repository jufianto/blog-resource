package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jufianto/blog-resource/big-files-processing/store"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile("env/env.yaml")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("viper config not found: %v", err)
	}

}

func main() {

	start := time.Now()

	file, err := os.Open("../resource/sales_5000000.csv")
	if err != nil {
		log.Fatal("failed to open ", err)
	}
	defer file.Close()

	var result []DataSales

	// concurrentWork := true
	// if !concurrentWork {
	// 	result, err = nonConcurrentMethod(file)
	// 	if err != nil {
	// 		log.Fatal("error", err)
	// 	}
	// } else {
	// 	result, err = concurrentMethod(file)
	// 	if err != nil {
	// 		log.Fatal("error", err)
	// 	}
	// }

	// test databases
	sqlConnStr := fmt.Sprintf(`postgres://%s:%s@%s:%s/%s?sslmode=disable`,
		viper.GetString("database.user"),
		viper.GetString("database.password"),
		viper.GetString("database.host"),
		viper.GetString("database.port"),
		viper.GetString("database.dbname"),
	)

	ctx := context.Background()

	configPool, err := pgxpool.ParseConfig(sqlConnStr)
	if err != nil {
		log.Fatal(err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, configPool)
	if err != nil {
		log.Fatalf("failed to fetch pool connection: %v \n", err)
	}
	defer pool.Close()

	storedb := store.NewStore(pool)

	id := uuid.New()

	if err := storedb.InsertSales(ctx, store.DataSales{
		Region:        "Asia",
		ID:            id,
		Country:       "ID",
		ItemType:      "ONL",
		SalesChannel:  "ONL",
		OrderPriority: "H",
		OrderDate:     time.Now(),
		OrderID:       fmt.Sprintf("%d", time.Now().Unix()),
		ShipDate:      time.Now().Add(24 * time.Hour),
		UnitSold:      100,
		UnitPrice:     100000,
		UnitCost:      120000,
		TotalRevenue:  900000,
		TotalCost:     800000,
		TotalProfit:   700000,
	}); err != nil {
		log.Fatalf("err insert: %v", err)
	}

	fmt.Println("total data", len(result))

	fmt.Printf("end in %d ms \n", time.Since(start).Milliseconds())
}

type DataSales struct {
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

func nonConcurrentMethod(file io.Reader) ([]DataSales, error) {
	log.Println("starting with non-concurrent method")
	data, err := SingleReadCSV(file, 100)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func concurrentMethod(file io.Reader) ([]DataSales, error) {
	log.Println("starting with concurrent method")

	data, err := ReadWithConcurrent(file)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func getResult(record []string) (DataSales, error) {
	unitSold, err := strconv.ParseFloat(record[8], 64)
	if err != nil {
		return DataSales{}, err
	}

	unitPrice, err := strconv.ParseFloat(record[9], 64)
	if err != nil {
		return DataSales{}, err
	}

	unitCost, err := strconv.ParseFloat(record[10], 64)
	if err != nil {
		return DataSales{}, err
	}

	totalRevenue, err := strconv.ParseFloat(record[11], 64)
	if err != nil {
		return DataSales{}, err
	}
	totalCost, err := strconv.ParseFloat(record[12], 64)
	if err != nil {
		return DataSales{}, err
	}
	totalProfit, err := strconv.ParseFloat(record[13], 64)
	if err != nil {
		return DataSales{}, err
	}

	orderDate, err := time.Parse("2006-01-02", record[5])
	if err != nil {
		return DataSales{}, err
	}

	shipDate, err := time.Parse("2006-01-02", record[7])
	if err != nil {
		return DataSales{}, err
	}

	time.Sleep(5 * time.Millisecond)

	sales := DataSales{
		Region:        record[0],
		Country:       record[1],
		ItemType:      record[2],
		SalesChannel:  record[3],
		OrderPriority: record[4],
		OrderDate:     orderDate,
		OrderID:       record[6],
		ShipDate:      shipDate,
		UnitSold:      unitSold,
		UnitPrice:     unitPrice,
		UnitCost:      unitCost,
		TotalRevenue:  totalRevenue,
		TotalCost:     totalCost,
		TotalProfit:   totalProfit,
	}
	return sales, nil
}
