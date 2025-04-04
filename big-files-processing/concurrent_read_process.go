package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"sync"
)

func ReadWithConcurrent(file io.Reader) ([]DataSales, error) {
	recordChan := make(chan []string, 10000)
	resultChan := make(chan DataSales, 10000)
	result := []DataSales{}
	wg := new(sync.WaitGroup)
	go readCSVRAW(file, recordChan, 100)

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	worker := 4
	for i := 0; i < worker; i++ {
		log.Println("starting worker ", i)
		wg.Add(1)
		go parseCSV(recordChan, resultChan, wg, i)
	}

	for res := range resultChan {
		result = append(result, res)
	}

	return result, nil
}

func readCSVRAW(file io.Reader, recordChan chan []string, limit int) {
	reader := csv.NewReader(file)
	_, _ = reader.Read() // read head
	defer close(recordChan)
	var i = 1
	for {
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

func parseCSV(recordChan chan []string, resultChan chan DataSales, wg *sync.WaitGroup, wid int) {
	defer wg.Done()

	for record := range recordChan {
		sale, err := getResult(record)
		if err != nil {
			fmt.Println("failed on read record: %v", err)
		}
		log.Println("get result from worker", wid, sale.OrderID)
		resultChan <- sale
	}
}
