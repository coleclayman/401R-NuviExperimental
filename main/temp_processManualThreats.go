package main

import (
	"database/sql"
	"fmt"
	"os"
	"../pgdatabase"
	"../filter"
	"../models"
	"encoding/json"
	"io/ioutil"
	"strings"
)

func main() {
	//nuvi, err := sql.Open("postgres", "postgres://go:gogo2017@localhost/nuvisocialthreat?sslmode=disable")
	//SqlError(err, "Error connecting to socialthreat")
	dao := pgdatabase.NewDAO()
	dao.Test()
	twitterFilter := new(filter.Filter)
	twitterFilter.InitFilter("danger.csv")
	manual, err := sql.Open("postgres", "postgres://go:gogo2017@localhost/manualthreats?sslmode=disable")
	SqlError(err, "Error connecting to socialthreat")

	rows, err := manual.Query("select uuid, threat, post from post")
	var id string
	var threat bool
	var storedPost []byte
	var post models.Post
	var body = "@relation deepTwitterData\n\n@attribute content string\n@attribute Class {false, true}\n\n@data\n"
	for rows.Next() {
		err := rows.Scan(&id, &threat, &storedPost)
		SqlError(err, "Error Scanning result set")

		json.Unmarshal(storedPost, &post)
		var threatString string
		if threat {
			threatString = "'true'"
		} else {
			threatString = "'false'"
		}
		var line = "'" + strings.Replace(strings.Replace(post.Raw_body_text, "\n", " ", -1), "'", "", -1) + "'," + threatString + "\n"
		fmt.Print(line)
		body = body + line
		//dao.AddRawPost(id, storedPost)
		//processed := twitterFilter.Preprocess(&post)
		//dao.AddProcessedPost(id, processed)
		//if threat {
		//	dao.LabelThreat(id, true)
		//} else {
		//	dao.LabelThreat(id, false)
		//}
	}
	ioutil.WriteFile("deepData.arff", []byte(body), 0644)
	fmt.Println("Complete")
}

func SqlError(err error, message string) {
	if err != nil {
		fmt.Println("SQL Error: ", message, "-> ", err.Error())
		os.Exit(1)
	}
}