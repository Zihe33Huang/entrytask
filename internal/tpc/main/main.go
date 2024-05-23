package main

import (
	"database/sql"
	ziherpc "entrytask/internal/rpc"
	tcp "entrytask/internal/tpc"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"math/rand"
	"net"
	"time"
)

func startServer() {
	var userService tcp.UserService
	if err := ziherpc.Register(&userService); err != nil {
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
	var err error
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/entrytask")
	if err != nil {
		log.Fatal(err)
	}
	tcp.SetDB(db)
	startServer()
	//insertUsers()

	//defer db.Close()
}

func createTable() {
	_, err := tcp.GetDB().Exec(`CREATE TABLE IF NOT EXISTS users (
		username VARCHAR(255) PRIMARY KEY,
		password VARCHAR(255),
		nickname VARCHAR(255),
		profile_image VARCHAR(255)
	)`)
	if err != nil {
		log.Fatal(err)
	}
}

func insertUsers() {
	rand.Seed(time.Now().UnixNano())
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/entrytask")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO users (username, password, nickname, profile_image) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for i := 0; i < 1000; i++ {
		username := fmt.Sprintf("user-%d", i)
		encode, err := tcp.Encode(fmt.Sprintf("password-%d", i))
		password := encode
		nickname := fmt.Sprintf("nickname-%d", i)
		profile_image := fmt.Sprintf("profile_image-%d.jpg", i)

		_, err = stmt.Exec(username, password, nickname, profile_image)
		if err != nil {
			log.Fatal(err)
		}
	}
}
