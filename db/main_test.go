package db_test

import (
	"context"
	"database/sql"
	"go-nuxt-blogs/db"
	"log"
	"os"
	"testing"
	"time"
)

var testDB *sql.DB
var testRepo db.Repository

func TestMain(m *testing.M) {
	testDB = openDB()
	testRepo = db.NewRepository(testDB)
	os.Exit(m.Run())
}

func openDB() *sql.DB {
	db, err := sql.Open("sqlite3", "../db.sqlite")
	if err != nil {
		log.Fatalln("openSQLite() failed: ", err)
	}

	// 连接到SQLite3,最多尝试连接5次
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatalln("db.PingContext() failed: ", err)
	}

	return db
}
