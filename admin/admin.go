package admin

import (
	"database/sql"
	"dbConnection/acc"
	"dbConnection/helper"
	"fmt"
	"strconv"
)

const KEY = "123"

//admin with ID = 1 is a root admin that can control other admins:
//(1) open a new admin account
//(2) delete an admin account

func Main(db *sql.DB) {
	for {
		fmt.Println("\t\t\t Admin")
		fmt.Println("\t\t (1) Login")
		fmt.Println("\t\t (2) Back")

		input := helper.InputText("\t\t\t Input")
		switch input {
		case "1":
			ok, id := login(db)
			if ok {
				profile(db, id)
			} else {
				fmt.Println("Login is unsuccessful")
			}
		case "2":
			return
		default:
			fmt.Println("Wrong input")
		}
	}
}

func profile(db *sql.DB, id int) {
	for {
		fmt.Println("\t\t\t Admin's Profile")
		fmt.Println("\t\t (1) New Admin")
		fmt.Println("\t\t (2) Delete Admin")
		fmt.Println("\t\t (3) Delete User")
		fmt.Println("\t\t (4) View unhandled messages")
		fmt.Println("\t\t (5) Handle messages")
		fmt.Println("\t\t (6) Log out")

		input := helper.InputText("\t\t\t Input")
		switch input {
		case "1":
			newAdmin(db, id)
		case "2":
			deleteAdmin(db, id)
		case "3":
			deleteUser(db, id)
		case "4":
			unhandledMessages(db)
		case "5":
			handleMessage(db, id)
		case "6":
			return
		default:
			fmt.Println("Wrong input")
		}
	}
}

func login(db *sql.DB) (bool, int) {
	//takes input for id and password
	//validates input id/password with ID/PASSWORD from DB
	//returns (TRUE, id) if id/password matches
	//returns (FALSE, 0) if id/password do not match

	//user input for id and check it from DB
	id := helper.InputText("ID")
	if !helper.IsCleanNum(id) {
		fmt.Println("Invalid ID")
		return false, 0
	}
	//convert id from String to Int
	idInt, err := strconv.Atoi(id)
	helper.CheckError(err)

	exist, passwordDB := isIdExist(db, idInt)
	if !exist {
		fmt.Println("This ID does not exist in DB")
		return false, 0
	}

	//user input for password and its validation with DB
	ok, password := acc.Password()
	if !ok {
		fmt.Println("Invalid password")
		return false, 0
	}
	if password != passwordDB {
		fmt.Println("Wrong password")
		return false, 0
	}

	fmt.Println("Login is successful")
	return true, idInt
}

func newAdmin(db *sql.DB, id int) {
	//ONLY id = 1 can use this function
	//takes user input for special Key
	//id for admin is assigned by system
	if id != 1 {
		fmt.Println("Admin does not have access to this functionality")
		return
	}

	//user input for Key and its validation
	key := helper.InputText("Key")
	if key != KEY {
		fmt.Println("Wrong key")
		return
	}

	//user input for password
	ok, password := acc.Password()
	if !ok {
		fmt.Println("Wrong input for password")
		return
	}

	//open a new admin
	sqlStatement := `
			INSERT INTO admins (password)
			VALUES ($1);`
	_, err := db.Exec(sqlStatement, password)
	helper.CheckError(err)

	fmt.Println("New Admin account is successfully opened")
}

func deleteAdmin(db *sql.DB, id int) {
	//only Root Admin with id = 1 can access to this method
	//takes input for Admin Id to be deleted
	//deletes Admin from DB, table: Admins
	if id != 1 {
		fmt.Println("Admin does not have access to this functionality")
		return
	}

	//user input for Key and its validation
	key := helper.InputText("Key")
	if key != KEY {
		fmt.Println("Wrong key")
		return
	}

	//input for Admin ID to be deleted
	deleteID := helper.InputText("Admin ID")
	if !helper.IsCleanNum(deleteID) {
		fmt.Println("Invalid ID")
		return
	}
	if deleteID == "1" {
		fmt.Println("Root Admin cannot be deleted")
		return
	}
	//check DB if Admin Id exists before its deletion
	deleteIDInt, _ := strconv.Atoi(deleteID)
	exist, _ := isIdExist(db, deleteIDInt)
	if !exist {
		fmt.Printf("Admin with ID = %v does not exist in Database\n",
			deleteIDInt)
		return
	}

	//make query to DB and delete Admin
	sqlStatement := `
		DELETE FROM admins
		WHERE adminid = $1;`
	_, err := db.Exec(sqlStatement, deleteIDInt)
	helper.CheckError(err)

	fmt.Println("Deleted successfully")
}

func deleteUser(db *sql.DB, id int) {
	//all admins can delete user accounts
	//takes user input for User Id to be deleted
	//validates its existance from table: users
	//deletes user from DB

	//user input for User ID
	user := helper.InputText("User ID")
	if !helper.IsCleanNum(user) {
		fmt.Println("Invalid id")
		return
	}
	userID, err := strconv.Atoi(user)
	helper.CheckError(err)

	exist, _ := isUserExist(db, userID)
	if !exist {
		fmt.Println("User does not exist")
		return
	}

	//query DB and delete user
	sqlStatement := `
			DELETE FROM users
			WHERE userid = $1;`
	_, err = db.Exec(sqlStatement, userID)
	helper.CheckError(err)

	fmt.Println("User deleted successfully")
}

func unhandledMessages(db *sql.DB) {
	//displays the list of all messages that are unhandled:
	//		column: "ischecked" = false
	//using this function, admin can see the messages that need to be
	//solved
	//messages are ordered by date
	sqlStatement := `
			SELECT messageid, message
			FROM messages
			WHERE ischecked = false
			ORDER BY date;`
	rows, err := db.Query(sqlStatement)
	helper.CheckError(err)
	defer rows.Close()

	fmt.Println("Message ID\t\tMessage")
	for rows.Next() {
		var id int
		var message string
		err = rows.Scan(&id, &message)
		helper.CheckError(err)
		fmt.Printf("%v\t\t%v\n", id, message)
	}
	err = rows.Err()
	helper.CheckError(err)
}

func handleMessage(db *sql.DB, id int) {
	//This function is called when the problem mentioned in a message is
	//	sorted out and the 'ischecked' column is given TRUE value
	//take input for message id,
	//updates DB, table: messages

	//input for message id
	message := helper.InputText("Message ID")
	if !helper.IsCleanNum(message) {
		fmt.Println("Invalid id")
		return
	}
	messageId, err := strconv.Atoi(message)
	helper.CheckError(err)

	//check if Message ID exists in DB, table: messages
	exist, _ := isMessageExist(db, messageId)
	if !exist {
		fmt.Println("Message does not exist")
		return
	}

	//make query to DB and update table: messages
	sqlStatement := `
			UPDATE messages
			SET ischecked = true,
				adminid = $1
			WHERE messageid = $2;`
	_, err = db.Exec(sqlStatement, id, messageId)
	helper.CheckError(err)

	fmt.Println("Message is handled successfully")
}

func isIdExist(db *sql.DB, id int) (bool, string) {
	//HELPER:
	//VALIDATOR:
	//makes query to DB to check for id,
	//returns (TRUE, password) if id exists in DB
	//returns (FALSE, "" empty string) if id does exist in DB
	var password string

	sqlStatement := `
			SELECT password
			FROM admins
			WHERE adminid = $1;`
	row := db.QueryRow(sqlStatement, id)

	switch err := row.Scan(&password); err {
	case sql.ErrNoRows:
		return false, ""
	case nil:
		return true, password
	default:
		panic(err)
	}
}

func isUserExist(db *sql.DB, id int) (bool, string) {
	//HELPER:
	//VALIDATOR:
	//makes query to DB to check for id,
	//returns TRUE if id exists in DB
	//returns FALSE if id does exist in DB
	var username string

	sqlStatement := `
			SELECT username
			FROM users
			WHERE userid = $1;`
	row := db.QueryRow(sqlStatement, id)

	switch err := row.Scan(&username); err {
	case sql.ErrNoRows:
		return false, ""
	case nil:
		return true, username
	default:
		panic(err)
	}
}

func isMessageExist(db *sql.DB, id int) (bool, int) {
	//HELPER:
	//VALIDATOR:
	//makes query to DB to check for id,
	//returns TRUE if id exists in DB
	//returns FALSE if id does exist in DB
	var user int

	sqlStatement := `
			SELECT userid
			FROM messages
			WHERE messageid = $1;`
	row := db.QueryRow(sqlStatement, id)

	switch err := row.Scan(&user); err {
	case sql.ErrNoRows:
		return false, -1
	case nil:
		return true, user
	default:
		panic(err)
	}
}
