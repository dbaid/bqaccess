package main

import (
	"cloud.google.com/go/bigquery"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/api/iterator"
	"time"
)

type Input struct {
	Input_script_bytes        string `json:"input_script_bytes"`
	Input_script_string       string `json:"input_script_string"`
	Input_script_string_error string `json:"input_script_string_error"`
	Input_sequence_number     uint64 `json:"input_sequence_number"`
	Input_pubkey_base58       string `json:"input_pubkey_base58"`
	Input_pubkey_base58_error string `json:"input_pubkey_base58_error"`
}

func main() {
	projectID := `bqaccess-346102`
	ctx := context.Background()
	client, _ := bigquery.NewClient(ctx, projectID)
	stmt := `SELECT to_json_string(inputs) FROM ` +
		"`bigquery-public-data.bitcoin_blockchain.transactions`" +
		` where array_length(inputs) = 2 LIMIT 1`
	q := client.Query(stmt)
	job, _ := q.Run(ctx)
	status, _ := job.Wait(ctx)
	if err := status.Err(); err != nil {
		msg := fmt.Errorf("bigquery.Value error : %v", err)
		var tz, _ = time.LoadLocation("Asia/Taipei")
		fmt.Println(msg.Error() + ` at ` + time.Now().In(tz).Format("2009-03-03 20:20:20"))
		return 
	}
	it, _ := job.Read(ctx)
	for {
		var row []bigquery.Value
		var Inputs []Input
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			msg := fmt.Errorf("bigquery.Value error : %v", err)
			var tz, _ = time.LoadLocation("Asia/Taipei")
			fmt.Println(msg.Error() + ` at ` + time.Now().In(tz).Format("2009-03-03 20:20:20"))
			return 
		}
        fmt.Println("============big query row array data====================")
		fmt.Println(row[0])
		input_str := fmt.Sprintf("%v", row[0])
		err = json.Unmarshal([]byte(input_str), &Inputs)
		fmt.Println("=============go slice element data====================")
		for _, input := range Inputs {
			fmt.Println("input_script_bytes :", input.Input_script_bytes)
			fmt.Println("input_script_string :", input.Input_script_string)
			fmt.Println("input_script_string_error :", input.Input_script_string_error)
			fmt.Println("input_sequence_number :", input.Input_sequence_number)
			fmt.Println("input_pubkey_base58 :", input.Input_pubkey_base58)
			fmt.Println("input_pubkey_base58_error :", input.Input_pubkey_base58_error)
		}
	}

}
