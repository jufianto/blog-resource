package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/jufianto/blog-resource/big-files-processing/store"
)

func ReadWithConcurrent(ctx context.Context, file io.Reader, limit int, dbstore store.StoreInterface) (int, error) {
	recordChan := make(chan []string, 100000)
	wg := new(sync.WaitGroup)
	totalWorker := worker
	totalChan := make(chan int, worker)

	go readCSVRAW(ctx, file, recordChan, limit)

	go func() {
		wg.Wait()
		close(totalChan)
	}()

	for i := 0; i < totalWorker; i++ {
		log.Println("starting worker ", i)
		wg.Add(1)
		go parseCSV(ctx, recordChan, totalChan, wg, i, dbstore)
	}

	var total = 0
	for t := range totalChan {
		total += t
	}

	return total, nil
}

func readCSVRAW(ctx context.Context, file io.Reader, recordChan chan []string, limit int) {
	reader := csv.NewReader(file)
	_, _ = reader.Read() // read head
	defer close(recordChan)
	var i = 1
	sep := 50000
	for {
		select {
		case <-ctx.Done():
			fmt.Println("context cancelled, stopping read csv")
			return
		default:
			// will continue the for
		}

		if i%sep == 0 {
			log.Printf("%d has read from csv \n", i)
		}

		if i > limit && limit > -1 {
			log.Printf("limited read count on loop %d \n", i-1)
			break
		}
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				fmt.Println("done loaded data")
				break
			}
			fmt.Println("error read", err)
			break
		}
		// log.Println("send data", i ,record[6])
		recordChan <- record
		i++
	}
}

func parseCSV(ctx context.Context, recordChan chan []string, totalChan chan int, wg *sync.WaitGroup, wid int, dbstore store.StoreInterface) {
	defer wg.Done()

	var i = 1
	sepLocal := sep

	for record := range recordChan {
		// log.Println("record to save left: ", len(recordChan))
		sale, err := getResult(record)
		if err != nil {
			fmt.Printf("failed on read record: %v \n", err)
			continue
		}

		id := uuid.New()
		if err := dbstore.InsertSales(ctx, store.DataSales{
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
			continue
		}

		// log.Println("get result from worker", wid, sale.OrderID)
		if i%sepLocal == 0 {
			log.Printf("%d record has been saved on worker %d \n", i, wid)
		}
		i++
	}

	totalChan <- i - 1
	fmt.Println("done with total data", i)

}
