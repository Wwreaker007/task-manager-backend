package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"
	"gopkg.in/gomail.v2"
)

const (
	PATH = ""
	PORT = ":9000"
)

func GetSMTPDialer() (*gomail.Dialer, error) {
	intPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		log.Println("Error in parsing SMTP port : ", err)
		return nil, err
	}
	return gomail.NewDialer(os.Getenv("SMTP_HOST"), intPort, os.Getenv("SMTP_EMAIL"), os.Getenv("SMTP_PASSWORD")), nil
}

func GetPostgressDBConnector() (*sql.DB, error) {
	url := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PSWD"), os.Getenv("DB_NAME"))

	// Get Postgres DB connection
	db, err := sql.Open("postgres", url)
	if(err != nil) {
		fmt.Println("unable to connect to postgresDB : ", err.Error(), url)
		return nil, err
	}

	// Ping to check the connection with the DB
	err = db.Ping()
	if err != nil {
		fmt.Println("Error connecting to the database: " + err.Error())
		return nil, err
	}
	fmt.Println("SUCCESFULLY CONNECTED TO DB")
	return db, nil
}

func main(){
	// Establish Postgres DB connection
	client, err := GetPostgressDBConnector()
	if err != nil {
		log.Fatalln("Error connecting to the database: ", err)
	}

	// Initialize SFTP dialer
	dialer, err := GetSMTPDialer()
	if err != nil {
		log.Fatalln("Error in connecting to SFTP server : ", err)
	}

	// Secret to be used in the JWT based authorization
	secretKey := os.Getenv("JWT_SECRET")

	// Create a new server
	server := NewServer(PATH, PORT, secretKey, client, dialer)

	// Service startup
	err = server.ServerStartup()
	if err != nil {
		log.Fatalln("Error starting the server: ", err)
	}
}