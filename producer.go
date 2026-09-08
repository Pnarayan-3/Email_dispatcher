package main

import (
	"encoding/csv"
	"os"
)

func loadRecipient(filepath string, ch chan Recipient) error {
	defer close(ch)
	
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return err
	}

	defer f.Close()

	for _, record := range records[1:] {
		//fmt.Println(record)
		ch <- Recipient{
			Name:  record[0],
			Email: record[1],
		}
		//send ->consumer->channels
	}
	return nil
}
