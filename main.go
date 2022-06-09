package main

import (
	"database/sql"
	"dbConnection/admin"
	"dbConnection/helper"
	"dbConnection/user"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	HOST     = "localhost"
	PORT     = 5432
	USER     = "_____" //default: "postgres"
	PASSWORD = "_____"
	DBNAME   = "bank"
)

func main() {
	//database connection
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		HOST, PORT, USER, PASSWORD, DBNAME)
	db, err := sql.Open("postgres", psqlInfo)
	helper.CheckError(err)

	defer db.Close()

	//Main menu
	for {
		fmt.Println("\t\t\t Main Menu")
		fmt.Println("\t\t (1) Admin")
		fmt.Println("\t\t (2) User")
		fmt.Println("\t\t (3) Exit")

		input := helper.InputText("\t\t\t Input")
		switch input {
		case "1":
			admin.Main(db)
		case "2":
			user.Main(db)
		case "3":
			return
		default:
			fmt.Println("Wrong input")
		}
	}
}
