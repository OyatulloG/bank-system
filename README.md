# Bank System
# Go, PostgreSQL

# This is a project imitating the banking system. It provides 2 types of accounts: admin and user.
  
  Admin has features:
    - Open new admin
    - Delete admin
    - Delete user
    - View unchecked user queries
    - Check user query
  
  User has features:
    - View account information
    - Update account information
    - Cash refill
    - Transfer 
    - Contact Admin, send query

# Content:
    (1.0) Run the program
    (2.0) Specific details
    (3.0) Database
    (4.0) Main functions

## 1.0 Run the program
1.1 Set up the database first. SQL commands are given in "SQLcommand.sql"
    file. Run this file in psql terminal. 

1.2 Change database details in "main.go" file: USER, PASSWORD.
    PORT value will not be changed, if it has a default value in your computer
    Guide for connecting Go to Database can be found here:
    https://www.calhoun.io/connecting-to-a-postgresql-database-with-gos-database-sql-package/

1.3 Run the "main.go" file in terminal.

## 2.0 Specific details
2.1 Admin with id = 1 is a root admin. Root Admin can Open and Delete
    admin accounts. Root Admin has a special KEY = "123".

2.2 Username for user account is valid if there is no white space and
    the length is equal or more than 4.

2.3 Password for user/admin account is valid if there is at least 1 digit,
    at least 1 symbol and the length is equal or more than 6.

## 3.0 Database
3.1 Database and its tables are created when the "SQLcommand.sql" 
    file runs.

3.2 Database: bank
    4 Tables: users, admins, messages, transfers
    
