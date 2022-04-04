package main

import (
	"cloud.google.com/go/bigquery"
	"context"
	"fmt"
	"google.golang.org/api/iterator"
	"strings"
	"time"
)

func main() {
	projectID := `bqaccess-346103`
	ctx := context.Background()
	client, _ := bigquery.NewClient(ctx, projectID)
	stmt := `SELECT 	to_json_string(addresses)   FROM ` +
		"`bigquery-public-data.crypto_bitcoin.inputs`" +
		` where array_length(addresses) = 2 LIMIT 1`
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
		address_str := fmt.Sprintf("%v", row[0])
		str1 := strings.Split(address_str, ",")
		fmt.Println("=============go slice element data====================")
		for i, str := range str1 {
			addr := strings.Trim(str, "\"[]")
			fmt.Println("Item :", i, ", data  :", addr)
		}
	}

}
