package user

import (
	"database/sql"
	"dbConnection/acc"
	"dbConnection/helper"
	"fmt"
	"strconv"
	"time"
)

func Main(db *sql.DB) {
	for {
		fmt.Println("\t\t\t User")
		fmt.Println("\t\t (1) Open account")
		fmt.Println("\t\t (2) Login")
		fmt.Println("\t\t (3) Back")

		input := helper.InputText("\t\t\t Input")
		switch input {
		case "1":
			newAccount(db)
		case "2":
			ok, id := login(db)
			if ok {
				profile(db, id)
			} else {
				fmt.Println("Login is unsuccessful")
			}
		case "3":
			return
		default:
			fmt.Println("Wrong input")
		}
	}
}

func profile(db *sql.DB, id int) {
	for {
		fmt.Println("\t\t\t User's Profile")
		fmt.Println("\t\t (1) Personal Info")
		fmt.Println("\t\t (2) Update personal Info")
		fmt.Println("\t\t (3) Cash Refill")
		fmt.Println("\t\t (4) Money transfer")
		fmt.Println("\t\t (5) Contact Admin")
		fmt.Println("\t\t (6) Log out")

		input := helper.InputText("\t\t\t Input")
		switch input {
		case "1":
			viewAccInfo(db, id)
		case "2":
			updateAccInfo(db, id)
		case "3":
			cashRefill(db, id)
		case "4":
			transfer(db, id)
		case "5":
			contactAdmin(db, id)
		case "6":
			return
		default:
			fmt.Println("Wrong input")
		}
	}
}

func newAccount(db *sql.DB) {
	//takes user input for firstname, lastname, date of birth, gender,
	//				username, password, PIN
	//sets default values for registration time = current time,
	//		balance = 0, account status = active, id
	//adds all above mentioned data to DB: bank,
	//								Tables: users

	fmt.Println("\tOpen a new account")
	//firstname
	ok, firstname := acc.FirstOrLastName("First name")
	if !ok {
		fmt.Println("Firstname input is not valid")
		return
	}

	//lastname
	ok, lastname := acc.FirstOrLastName("Last name")
	if !ok {
		fmt.Println("Lastname is not valid")
		return
	}

	//date of birth
	ok, date := acc.DateOfBirth()
	if !ok {
		fmt.Println("Date of Birth is not valid")
		return
	}

	//gender
	ok, gender := acc.Gender()
	if !ok {
		fmt.Println("Input is not valid")
		return
	}

	//username
	ok, username := acc.Username()
	if !ok {
		fmt.Println("Username is not valid")
		return
	}
	ok, _, _ = isUsernameExist(username, db)
	if ok {
		fmt.Println("This username already exists in DB")
		return
	}

	//password
	ok, password := acc.Password()
	if !ok {
		fmt.Println("Password is not valid")
		return
	}

	//PIN
	ok, pin := pin()
	if !ok {
		fmt.Println("PIN is not valid")
		return
	}

	//move data to DB, Table: users
	sqlStatement := `
		INSERT INTO users (firstname, lastname, dateofbirth, gender, pin,
					username, password)
		VALUES ($1, $2, $3, $4, $5, $6, $7);`
	_, err := db.Exec(sqlStatement, firstname, lastname, date, gender,
		pin, username, password)
	helper.CheckError(err)

	fmt.Println("New account is successfully opened!")
}

func login(db *sql.DB) (bool, int) {
	//takes input for username and password,
	//checks if account with 'username' exists in DB
	//returns (TRUE, ID) of account if username exists and
	//									input password matches
	//returns (FALSE, 0) if username does not exist in DB and
	//				input password does not match with password in DB

	fmt.Println("\t Login")
	//username
	ok, username := acc.Username()
	if !ok {
		fmt.Println("Username input is not valid")
		return false, 0
	}
	ok, id, passWord := isUsernameExist(username, db)
	if !ok {
		fmt.Println("This username does not exist in DB")
		return false, 0
	}

	//password
	ok, password := acc.Password()
	if !ok {
		fmt.Println("Password is not valid")
		return false, 0
	}

	//validate input password with passWord from DB
	if password != passWord {
		fmt.Println("Wrong password")
		return false, 0
	}

	fmt.Println("Login is successful!")
	return true, id
}

func cashRefill(db *sql.DB, id int) {
	//this function imitates the process of refilling balance
	//			with cash
	//takes input for PIN,
	//validates input PIN with PIN from DB,
	//takes input for cash, performs transaction to balance of account
	//updates 'users' table after transaction,
	//adds data about transaction to 'transfers' table,
	//prints "Success" if PIN is correct and transaction is successful
	//prints "Error" if PIN is wrong and transaction is unsuccessful

	//user input for PIN
	ok, pin := pin()
	if !ok {
		fmt.Println("PIN is wrong!")
		return
	}

	//query DB to get PIN and balance
	pinDB := ""
	var balance float64
	sqlStatement := `
			SELECT pin, balance FROM users
			WHERE userid = $1;`
	err := db.QueryRow(sqlStatement, id).Scan(&pinDB, &balance)
	helper.CheckError(err)

	//validate pin with PIN from DB
	if pin != pinDB {
		fmt.Println("Wrong PIN code")
		return
	}

	//user input for cash to make transaction
	cash := helper.InputText("Cash")
	cashFloat := helper.ConvStrToFloat64(cash)
	if cashFloat == -1. {
		fmt.Println("Wrong value for cash")
		return
	}

	//update balance
	sqlStatement = `
			UPDATE users
			SET balance = $1
			WHERE userid = $2;`
	_, err = db.Exec(sqlStatement, cashFloat+balance, id)
	helper.CheckError(err)

	//update transfers table
	sqlStatement = `
			INSERT INTO transfers (senderid, recieverid, amount)
			VALUES ($1, $2, $3);`
	_, err = db.Exec(sqlStatement, id, id, cashFloat)
	helper.CheckError(err)

	fmt.Println("Transaction is successful!")
	fmt.Println("Balance is updated")
}

func transfer(db *sql.DB, idSender int) {
	//this function imitates the process of transferring money from
	//	one account to another
	//(1) takes user input for PIN of sender,
	//(2) validates input PIN with PIN from DB,
	//(3) takes user input for money and validates it
	//(4) check if transferred money <= balance of sender,
	//(5) updates balances of Sender and Reciever
	//(6) updates transfers table and addes new transfer data to it
	//prints "Success" if transfer is successful
	//prints "Error" if one of above mentioned validations fails

	//user input for sender's PIN
	ok, pin := pin()
	if !ok {
		fmt.Println("Wrong PIN code")
		return
	}

	//query DB and retrieve data for PIN and balance of Sender
	var pinSender string
	var balanceSender float64
	sqlStatement := `
			SELECT pin, balance FROM users
			WHERE userid = $1;`
	err := db.QueryRow(sqlStatement, idSender).Scan(&pinSender, &balanceSender)
	helper.CheckError(err)

	//input for Receiver's id
	receiver := helper.InputText("Receiver Id")
	if !helper.IsCleanNum(receiver) {
		fmt.Println("Wrong ID")
		return
	}
	idReciever, err := strconv.Atoi(receiver)
	helper.CheckError(err)

	//transfer cannot be done when idSender and idReciever are the same
	if idSender == idReciever {
		fmt.Println("Transfer cannot be performed")
		return
	}

	//query DB, retrieve data for Reciever's balance if idReceiver exists
	var balanceReciever float64
	sqlStatement = `
			SELECT balance FROM users
			WHERE userid = $1;`
	err = db.QueryRow(sqlStatement, idReciever).Scan(&balanceReciever)
	if err == sql.ErrNoRows {
		fmt.Printf("ID = %v does not exist in DB\n", idReciever)
		return
	}
	helper.CheckError(err)

	//validate input PIN with PIN from DB
	if pin != pinSender {
		fmt.Println("Wrong PIN code")
		return
	}

	//user input for money
	money := helper.InputText("Money")
	moneyFloat := helper.ConvStrToFloat64(money)
	if moneyFloat == -1 {
		fmt.Println("Invalid input for money")
		return
	}

	//transferred money should be less than or equal to Sender's balance
	if moneyFloat > balanceSender {
		fmt.Println("There is not enough money in balance")
		return
	}

	//perform transfer, update Sender's and Recipiever's balances
	sqlStatementSender := `
			UPDATE users
			SET balance = $1
			WHERE userid = $2;`
	_, err = db.Exec(sqlStatementSender, balanceSender-moneyFloat,
		idSender)
	helper.CheckError(err)

	sqlStatementReciever := `
			UPDATE users
			SET balance = $1
			WHERE userid = $2;`
	_, err = db.Exec(sqlStatementReciever, balanceReciever+moneyFloat,
		idReciever)
	helper.CheckError(err)

	//update transfers table
	sqlStatement = `
			INSERT INTO transfers (senderid, recieverid, amount)
			VALUES ($1, $2, $3);`
	_, err = db.Exec(sqlStatement, idSender, idReciever, moneyFloat)
	helper.CheckError(err)

	fmt.Println("Transfer is successful")
}

func updateAccInfo(db *sql.DB, id int) {
	//displays menu to show what data can be updated,
	//takes user input for Update Choice,
	//calls required update function
	//prints Success if update is done
	//prints Unsuccess if update is not done

	for {
		fmt.Println("\t\t\t Update:")
		fmt.Println("\t\t (1) Firstname")
		fmt.Println("\t\t (2) Lastname")
		fmt.Println("\t\t (3) Date of Birth")
		fmt.Println("\t\t (4) Gender")
		fmt.Println("\t\t (5) Username")
		fmt.Println("\t\t (6) Password")
		fmt.Println("\t\t (7) Back")

		var ok bool
		input := helper.InputText("\tInput")
		switch input {
		case "1":
			ok = updateFirstname(db, id)
		case "2":
			ok = updateLastname(db, id)
		case "3":
			ok = updateDOB(db, id)
		case "4":
			ok = updateGender(db, id)
		case "5":
			ok = updateUsername(db, id)
		case "6":
			ok = updatePassword(db, id)
		case "7":
			return
		default:
			fmt.Println("Wrong input")
		}

		if !ok {
			fmt.Println("Update is unsuccessful")
		}
		fmt.Println("Update is successful")
	}
}

func viewAccInfo(db *sql.DB, id int) {
	//displays menu to show what data can be viewed,
	//takes user input for View Choice,
	//queries DB for current value of View Choice,
	//prints required Data if view is successful,

	for {
		fmt.Println("\t\t View:")
		fmt.Println("(1) Firstname")
		fmt.Println("(2) Lastname")
		fmt.Println("(3) Date of Birth")
		fmt.Println("(4) Gender")
		fmt.Println("(5) Username")
		fmt.Println("(6) Password")
		fmt.Println("(7) PIN")
		fmt.Println("(8) Status")
		fmt.Println("(9) Balance")
		fmt.Println("(10) Transfer History")
		fmt.Println("(11) Registration Time")
		fmt.Println("(12) No change")

		view := helper.InputText("\tInput")
		switch view {
		case "1":
			firstname := getFirstname(db, id)
			fmt.Println("Firstname: ", firstname)
		case "2":
			lastname := getLastname(db, id)
			fmt.Println("Lastname: ", lastname)
		case "3":
			dateOfBirth := getDateOfBirth(db, id)
			fmt.Println("Date of Birth: ", dateOfBirth)
		case "4":
			gender := getGender(db, id)
			fmt.Println("Gender: ", gender)
		case "5":
			username := getUsername(db, id)
			fmt.Println("Username: ", username)
		case "6":
			password := getPassword(db, id)
			fmt.Println("Password: ", password)
		case "7":
			pin := getPIN(db, id)
			fmt.Println("PIN: ", pin)
		case "8":
			status := getStatus(db, id)
			fmt.Println("Status: ", status)
		case "9":
			balance := getBalance(db, id)
			fmt.Println("Balance: ", balance)
		case "10":
			fmt.Println("Transfer History: ")
			getTransferHistory(db, id)
		case "11":
			regTime := getRegTime(db, id)
			fmt.Println("Registration Time: ", regTime)
		case "12":
			fmt.Println("No change")
			return
		default:
			fmt.Println("Wrong input")
		}
	}
}

func contactAdmin(db *sql.DB, id int) {
	//takes input for message to be sent to admins,
	//message can be an issue regarding system like:
	//	deletion, disactivation of account, problems with transction, etc
	//message is moved to Table: messages
	//prints Success if query is performed
	//causes ERR if query is unsuccessful

	//user input for message
	message := helper.InputText("Message: ")

	//Insert data to DB
	sqlStatement := `
			INSERT INTO messages (userid, adminid, message)
			VALUES ($1, $2, $3);`
	_, err := db.Exec(sqlStatement, id, 1, message)
	helper.CheckError(err)

	fmt.Println("Message is sent successfully")
}

func pin() (bool, string) {
	//takes user input for PIN code of account,
	//PIN code should be 4 digit number only,
	//returns (TRUE, PIN) if it passes validations,
	//returns (FALSE, "") empty string if it does not pass validations

	pin := helper.InputText("PIN")
	if helper.IsCleanNum(pin) && helper.IsLengthInLimit(pin, 4) {
		return true, pin
	}

	return false, ""
}

func isUsernameExist(username string, db *sql.DB) (bool, int, string) {
	//HELPER:
	//VALIDATOR:
	//makes query to DB to check for username,
	//returns (TRUE, id, password) if username exists in DB
	//returns (FALSE, 0, "" empty string) if username does exist in DB
	var id int
	var password string

	sqlStatement := `
			SELECT userid, password
			FROM users
			WHERE username = $1;`
	row := db.QueryRow(sqlStatement, username)

	switch err := row.Scan(&id, &password); err {
	case sql.ErrNoRows:
		return false, 0, ""
	case nil:
		return true, id, password
	default:
		panic(err)
	}
}

func updateFirstname(db *sql.DB, id int) bool {
	//retrieves current Firstname from DB,
	//takes user input for new Firstname and validates it
	//updates firstname column of DB

	//query DB and retrieve current firstname
	current := getFirstname(db, id)
	fmt.Println("Current firstname: ", current)

	//user input for new firstname
	ok, new := acc.FirstOrLastName("Firstname")
	if !ok {
		fmt.Println("Invalid input for Firstname")
		return false
	}

	//update DB
	sqlStatement := `
			UPDATE users 
			SET firstname = $1
			WHERE userid = $2;`
	_, err := db.Exec(sqlStatement, new, id)
	helper.CheckError(err)

	fmt.Println("Firstname is updated successfully!")
	return true
}

func updateLastname(db *sql.DB, id int) bool {
	//retrieves current Lastname from DB,
	//takes user input for new Lastname and validates it
	//updates lastname column of DB

	//query DB and retrieve current lastname
	current := getLastname(db, id)
	fmt.Println("Current lastname: ", current)

	//user input for new lastname
	ok, new := acc.FirstOrLastName("Lastname")
	if !ok {
		fmt.Println("Invalid input for Lastname")
		return false
	}

	//update DB
	sqlStatement := `
			UPDATE users
			SET lastname = $1
			WHERE userid = $2;`
	_, err := db.Exec(sqlStatement, new, id)
	helper.CheckError(err)

	fmt.Println("Lastname is updated successfully!")
	return true
}

func updateDOB(db *sql.DB, id int) bool {
	//retrieves current Date of birth from DB,
	//takes user input for new Date of birth and validates it
	//updates dateofbirth column of DB

	//query DB and retrieve current Date of birth
	current := getDateOfBirth(db, id)
	fmt.Println("Current Date of Birth: ", current)

	//user input for new dateofbirth
	ok, new := acc.DateOfBirth()
	if !ok {
		fmt.Println("Invalid input for Date of Birth")
		return false
	}

	//update DB
	sqlStatement := `
			UPDATE users
			SET dateofbirth = $1
			WHERE userid = $2;`
	_, err := db.Exec(sqlStatement, new, id)
	helper.CheckError(err)

	fmt.Println("Date of Birth is updated successfully!")
	return true
}

func updateGender(db *sql.DB, id int) bool {
	//retrieves current Gender from DB,
	//takes user input for new Gender and validates it
	//updates gender column of DB

	//query DB and retrieve current gender
	current := getGender(db, id)
	fmt.Println("Current gender: ", current)

	//user input for new gender
	ok, new := acc.Gender()
	if !ok {
		fmt.Println("Invalid input for Gender")
		return false
	}

	//update DB
	sqlStatement := `
			UPDATE users
			SET gender = $1
			WHERE userid = $2;`
	_, err := db.Exec(sqlStatement, new, id)
	helper.CheckError(err)

	fmt.Println("Gender is updated successfully!")
	return true
}

func updateUsername(db *sql.DB, id int) bool {
	//retrieves current Username from DB,
	//takes user input for new Username and validates it
	//check whether or not new username is already used by another account
	//updates username column of DB

	//query DB and retrieve current username
	current := getUsername(db, id)
	fmt.Println("Current username: ", current)

	//user input for new firstname
	ok, new := acc.Username()
	if !ok {
		fmt.Println("Invalid input for Username")
		return false
	}
	ok, _, _ = isUsernameExist(new, db)
	if ok {
		fmt.Println("This username already exists")
		return false
	}

	//update DB
	sqlStatement := `
			UPDATE users
			SET username = $1
			WHERE userid = $2;`
	_, err := db.Exec(sqlStatement, new, id)
	helper.CheckError(err)

	fmt.Println("Username is updated successfully!")
	return true
}

func updatePassword(db *sql.DB, id int) bool {
	//retrieves current password from DB,
	//takes user input for new password and validates it
	//updates password column of DB

	//query DB and retrieve current password
	current := getPassword(db, id)
	fmt.Println("Current password: ", current)

	//user input for new password
	ok, new := acc.Password()
	if !ok {
		fmt.Println("Invalid input for Password")
		return false
	}

	//update DB
	sqlStatement := `
			UPDATE users
			SET password = $1
			WHERE userid = $2;`
	_, err := db.Exec(sqlStatement, new, id)
	helper.CheckError(err)

	fmt.Println("Password is updated successfully!")
	return true
}

func getFirstname(db *sql.DB, id int) string {
	//makes query to DB and gets firstname
	//returns Firstname

	var firstname string
	sqlStatement := `
			SELECT firstname FROM users
			WHERE userid = $1;`

	err := db.QueryRow(sqlStatement, id).Scan(&firstname)
	helper.CheckError(err)

	return firstname
}

func getLastname(db *sql.DB, id int) string {
	//makes query to DB and gets lastname
	//returns Lastname
	var lastname string
	sqlStatement := `
			SELECT lastname FROM users
			WHERE userid = $1;`

	err := db.QueryRow(sqlStatement, id).Scan(&lastname)
	helper.CheckError(err)

	return lastname
}

func getDateOfBirth(db *sql.DB, id int) time.Time {
	//makes query to DB and gets date of birth
	//returns Date of Birth
	var dateOfBirth time.Time
	sqlStatement := `
			SELECT dateofbirth FROM users
			WHERE userid = $1;`

	err := db.QueryRow(sqlStatement, id).Scan(&dateOfBirth)
	helper.CheckError(err)

	return dateOfBirth
}

func getGender(db *sql.DB, id int) string {
	//makes query to DB and gets gender
	//returns Gender
	var gender string
	sqlStatement := `
			SELECT gender FROM users
			WHERE userid = $1;`

	err := db.QueryRow(sqlStatement, id).Scan(&gender)
	helper.CheckError(err)

	return gender
}

func getUsername(db *sql.DB, id int) string {
	//gets username from DB
	//returns Username
	var username string
	sqlStatement := `
			SELECT username FROM users
			WHERE userid = $1;`

	err := db.QueryRow(sqlStatement, id).Scan(&username)
	helper.CheckError(err)

	return username
}

func getPassword(db *sql.DB, id int) string {
	//gets password from DB
	//returns Password
	var password string
	sqlStatement := `
			SELECT password FROM users
			WHERE userid = $1;`

	err := db.QueryRow(sqlStatement, id).Scan(&password)
	helper.CheckError(err)

	return password
}

func getPIN(db *sql.DB, id int) string {
	//makes query to DB and gets PIN
	//returns PIN
	var pin string
	sqlStatement := `
			SELECT pin FROM users
			WHERE userid = $1;`

	err := db.QueryRow(sqlStatement, id).Scan(&pin)
	helper.CheckError(err)

	return pin
}

func getStatus(db *sql.DB, id int) bool {
	//makes query to DB and gets status
	//returns Status
	var status bool
	sqlStatement := `
			SELECT status FROM users
			WHERE userid = $1;`

	err := db.QueryRow(sqlStatement, id).Scan(&status)
	helper.CheckError(err)

	return status
}

func getBalance(db *sql.DB, id int) float64 {
	//makes query to DB and gets balance
	//returns Balance
	var balance float64
	sqlStatement := `
			SELECT balance FROM users
			WHERE userid = $1;`

	err := db.QueryRow(sqlStatement, id).Scan(&balance)
	helper.CheckError(err)

	return balance
}

func getTransferHistory(db *sql.DB, id int) {
	//makes query to DB and gets all transfer data about id
	//returns Transfer History
	var transferid, senderid int
	var amount float64
	var date time.Time
	sqlStatement := `
			SELECT transferid, senderid, amount, date
			FROM transfers
			WHERE senderid = $1 OR recieverid = $2;`

	rows, err := db.Query(sqlStatement, id, id)
	helper.CheckError(err)
	defer rows.Close()

	fmt.Println("Transfer ID\t\tSender ID\t\tAmount\t\tDate")
	for rows.Next() {
		err = rows.Scan(&transferid, &senderid, &amount, &date)
		helper.CheckError(err)
		fmt.Printf("%v\t\t\t%v\t\t\t%v\t\t\t%v\n", transferid, senderid,
			amount, date)
	}

	err = rows.Err()
	helper.CheckError(err)
}

func getRegTime(db *sql.DB, id int) time.Time {
	//makes query to DB and gets Registration Time
	//returns Registration Time
	var regtime time.Time
	sqlStatement := `
			SELECT regtime FROM users
			WHERE userid = $1;`

	err := db.QueryRow(sqlStatement, id).Scan(&regtime)
	helper.CheckError(err)
	return regtime
}
