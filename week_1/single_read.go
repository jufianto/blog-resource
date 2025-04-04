package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
)

func SingleReadCSV(file io.Reader, limit int) (sales []DataSales, err error) {
	reader := csv.NewReader(file)
	var i = 1

	// skip header
	_, _ = reader.Read()

	for {
		if i > limit && limit >= -1 {
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
			return nil, err
		}

		sale, err := getResult(record)
		if err != nil {
			log.Printf("failed to get result on loop %d: %v \n", i+1, err)
			return nil, err
		}

		sales = append(sales, sale)
		i++
	}
	return sales, nil
}
