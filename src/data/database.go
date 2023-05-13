package data

import (
	"crypto/tls"
	"fmt"
	"os"

	"bitbucket.org/marcoboschetti/catalogscraper/src/entities"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

var pgConnection *pg.DB

func SetDbConnection() {
	if pgConnection != nil {
		return
	}

	address := os.Getenv("POSTGRES_URL")
	port := os.Getenv("POSTGRES_PORT")
	username := os.Getenv("POSTGRES_USERNAME")
	password := os.Getenv("POSTGRES_PASSWORD")
	database := os.Getenv("POSTGRES_DATABASE")

	pgConnection = pg.Connect(&pg.Options{
		Addr:      fmt.Sprintf("%s:%s", address, port),
		Database:  database,
		User:      username,
		Password:  password,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
	})

	err := createSchema(pgConnection)
	if err != nil {
		panic(err.Error())
	}
}

func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{
		(*entities.Category)(nil),
		(*entities.Subcategory)(nil),
		(*entities.Product)(nil),
	} {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp:          false,
			IfNotExists:   true,
			FKConstraints: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
