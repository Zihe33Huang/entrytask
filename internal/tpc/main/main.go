package main

import (
	"database/sql"
	ziherpc "entrytask/internal/rpc"
	tcp "entrytask/internal/tpc"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net"
	"sync"
)

func startServer() {
	var userService tcp.UserService
	if err := ziherpc.Register(&userService); err != nil {
		log.Fatal("register error:", err)
	}
	var authService tcp.AuthService
	if err := ziherpc.Register(&authService); err != nil {
		log.Fatal("register error:", err)
	}
	// pick a free port
	l, err := net.Listen("tcp", "localhost:8888")
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Println("start rpc server on", l.Addr())
	ziherpc.Accept(l)
}

func main() {
	//var err error
	//db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/entrytask")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//tcp.SetDB(db)
	//startServer()
	//defer db.Close()

	createTable()
	generateMillionUsers()
}

func createTable() {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/entrytask")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("DROP TABLE IF EXISTS users")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
    	id       BIGINT  PRIMARY KEY,
		username VARCHAR(255),
		password VARCHAR(255),
		nickname VARCHAR(255),
		profile_image VARCHAR(255)
	)`)

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE UNIQUE INDEX idx_username ON users (username)`)
	if err != nil {
		log.Fatal(err)
	}
}

func generateMillionUsers() {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/entrytask")
	db.SetMaxIdleConns(10)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	wg := &sync.WaitGroup{}
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go insertUsers(db, wg, i)
	}

	wg.Wait()
}

func insertUsers(db *sql.DB, wg *sync.WaitGroup, idx int) {
	defer wg.Done()

	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	batchSize := 1000
	totalRecords := 1000000
	starting := totalRecords * idx
	values := make([]interface{}, batchSize*5) // Assuming 2 columns

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	insertQuery := "INSERT INTO users ( id, username, password, nickname, profile_image) VALUES "
	valueStr := "(?, ?, ?, ?, ?)"

	for i := 0; i < totalRecords/batchSize; i++ {
		query := insertQuery
		for j := 0; j < batchSize; j++ {
			if j > 0 {
				query += ", "
			}
			query += valueStr
			userId := starting + i*batchSize + j
			values[j*5] = userId
			values[j*5+1] = fmt.Sprintf("user_%d", userId)
			values[j*5+2] = fmt.Sprintf("psw_%d", userId)
			values[j*5+3] = fmt.Sprintf("nickname_%d", userId)
			values[j*5+4] = fmt.Sprintf("")
		}

		stmt, err := tx.Prepare(query)
		if err != nil {
			log.Fatal(err)
		}

		_, err = stmt.Exec(values...)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
		stmt.Close()
	}
	err = tx.Commit()

	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	//stmt, err := db.Prepare("INSERT INTO users (id, username, password, nickname, profile_image) VALUES (?, ?, ?, ?, ?)")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer stmt.Close()
	//batchSize := 100000
	//starting := idx * batchSize
	//for i := 0; i < 100000; i++ {
	//	userId := starting + i
	//	username := fmt.Sprintf("user-%d", userId)
	//	var password string
	//	password, err = tcp.Encode(fmt.Sprintf("password-%d", userId))
	//	if err != nil {
	//		log.Panic("encode error", err.Error())
	//	}
	//	nickname := fmt.Sprintf("nickname-%d", userId)
	//	profile_image := ""
	//
	//	_, err = stmt.Exec(userId, username, password, nickname, profile_image)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	// Print progress bar
	//}

}
