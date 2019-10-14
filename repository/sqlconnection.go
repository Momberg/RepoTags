package repository

import (
	"log"

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
	db, err = sqlx.Open("mysql", "root@tcp(localhost:3306)/gittags?parseTime=true")
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
		log.Println("[GetDBConnection] Erro no ping na conexao: ", err.Error())
		return
	}
	localdb = db
	return
}
