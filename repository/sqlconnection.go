package repository

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	/*
		github.com/go-sql-driver/mysql não é usado diretamente pela aplicação
	*/
	_ "github.com/go-sql-driver/mysql"
)

//variavel singleton que armazena a conexao
var db *sqlx.DB

//OpenConnection open MYSQL connection
func OpenConnection() (db *sqlx.DB, err error) {
	err = nil
	db, err = sqlx.Open("mysql", createURLConnection())
	if err != nil {
		log.Println("[OpenConnection] Connection error: ", err.Error())
		return
	}
	err = db.Ping()
	if err != nil {
		log.Println("[OpenConnection] Connection ping error: ", err.Error())
		return
	}
	return
}

//GetDBConnection get DB connection
func GetDBConnection() (localdb *sqlx.DB, err error) {
	if db == nil {
		db, err = OpenConnection()
		if err != nil {
			log.Println("[GetDBConnection] Connection error: ", err.Error())
			return
		}
	}
	err = db.Ping()
	if err != nil {
		log.Println("[GetDBConnection] Connection ping error: ", err.Error())
		return
	}
	localdb = db
	return
}

//createURLConnection create the mysql url to connect
func createURLConnection() string {
	user := os.Getenv("MYSQLUSER")
	pass := os.Getenv("MYSQLPASS")
	host := os.Getenv("MYSQLHOST")
	url := ""
	if host == "" {
		log.Println("[createUrlConnection] Configure the mysql host.")
	} else if user == "" {
		log.Println("[createUrlConnection] Configure the mysql user.")
	} else if pass != "" {
		url = user + ":" + pass + "@tcp(" + host + ")/gittags?parseTime=true"
	} else {
		url = user + "@tcp(" + host + ")/gittags?parseTime=true"
	}
	return url
}
