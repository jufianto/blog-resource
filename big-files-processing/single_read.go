package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"

	"github.com/google/uuid"
	"github.com/jufianto/blog-resource/big-files-processing/store"
)

func SingleReadCSV(file io.Reader, limit int, str store.StoreInterface) (totalInserted int, err error) {
	reader := csv.NewReader(file)
	var i = 1

	// skip header
	_, _ = reader.Read()

	for {
		if i > limit && limit > -1 {
			log.Printf("limited read count: %d \n", i-1)
			break
		}

		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				fmt.Println("done loaded the data csv")
				break
			}
			fmt.Printf("error read data on loop %d: %v", i+1, err)
			return i - 1, err
		}

		sale, err := getResult(record)
		if err != nil {
			log.Printf("failed to get result on loop %d: %v \n", i+1, err)
			return i - 1, err
		}

		ctx := context.Background()
		id := uuid.New()

		if err := str.InsertSales(ctx, store.DataSales{
			ID:            id,
			Region:        sale.Region,
			Country:       sale.Country,
			ItemType:      sale.ItemType,
			SalesChannel:  sale.SalesChannel,
			OrderPriority: sale.OrderPriority,
			OrderDate:     sale.OrderDate,
			OrderID:       sale.OrderID,
			ShipDate:      sale.ShipDate,
			UnitSold:      sale.UnitSold,
			UnitPrice:     sale.UnitPrice,
			UnitCost:      sale.UnitCost,
			TotalRevenue:  sale.TotalRevenue,
			TotalCost:     sale.TotalCost,
			TotalProfit:   sale.TotalProfit,
		}); err != nil {
			log.Printf("[ERROR] failed to insert data: %v \n", err)
			return i - 1, err
		}

		if i%sep == 0 {
			log.Printf("%d has been inserted \n", i)
		}

		i++
	}
	return i - 1, nil
}
