package utils

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func Insert_history(startTime time.Time, endTime time.Time, putCount int, crawlCount int) {
	db, err := sql.Open("mysql", "root:5623130@tcp(127.0.0.1:3306)/crawler")
	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	result, err := db.Exec("Insert Into crawl_history (crawl_start_time, crawl_end_time, elastic_put_count, crawl_count) Values (? ? ? ?)",
		startTime, endTime, putCount, crawlCount)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(result)
	// Connect and check the server version
}
