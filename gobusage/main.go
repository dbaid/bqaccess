package main

import (
	"database/sql"
	"encoding/gob"
	"fmt"
	"log"
	"os"

	go_ora "github.com/sijms/go-ora/v2"
)

type test_struc struct {
	Id     int
	Name   string
	Sup    string
	Crdate uint64
}

func main() {
	var test_val test_struc
	var tests []test_struc
	urlOptions := map[string]string{
		"CONNECTION TIMEOUT": "3",
	}
	databaseUrl := go_ora.BuildUrl("localhost", 1521, "XE", "system", "xxxxx", urlOptions)
	conn, err := sql.Open("oracle", databaseUrl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = conn.Ping()
	if err != nil {
		fmt.Println(err)
		return
	}
	query_string := "select id, name ,sup, cast((create_date - TO_DATE('01011970000000' , 'ddmmyyyyhh24miss')) * 86400 as number(19)) from test"
	rows, err := conn.Query(query_string)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&test_val.Id, &test_val.Name, &test_val.Sup, &test_val.Crdate)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%d %q %q %d \n", test_val.Id, test_val.Name, test_val.Sup, test_val.Crdate)
		tests = append(tests, test_val)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("%q", tests)

	dataFile, err := os.Create("integerdata.gob")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// serialize the data
	dataEncoder := gob.NewEncoder(dataFile)
	dataEncoder.Encode(tests)

	dataFile.Close()

	var tests2 []test_struc
	dataFile2, err := os.Open("integerdata.gob")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("1")
	dataDecoder := gob.NewDecoder(dataFile2)
	fmt.Println("2")
	err = dataDecoder.Decode(&tests2)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dataFile2.Close()

	fmt.Println("tests2: ", tests2)

}


