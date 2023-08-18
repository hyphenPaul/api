package main

// import (
// 	"context"
// 	"fmt"
// 	"os"
//
// 	"github.com/jackc/pgx/v5"
// 	"github.com/jackc/pgx/v5/pgxpool"
// )
//
// var databaseURL string = "postgresql://myuser:mypassword@postgres:5432/mydb"
//
// func startDatabase() {
// 	fmt.Println("Starting the database")
// 	dbpool, err := pgxpool.New(context.Background(), databaseURL)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Unable to create connection pool %v\n", err)
// 		os.Exit(1)
// 	}
// 	defer dbpool.Close()
//
// 	var p Person
// 	q := "select id, lastname, firstname, age from people where id = 122"
// 	err = dbpool.QueryRow(context.Background(), q).Scan(&p.ID, &p.LastName, &p.FirstName, &p.Age)
// 	if err != nil {
// 		if err == pgx.ErrNoRows {
// 			fmt.Println("There are no rows for that id: 122")
// 			return
// 		}
// 		fmt.Fprintf(os.Stderr, "Unable to populate p %v\n", err)
// 	}
// 	fmt.Printf("Person: %+v\n", p)
// }
