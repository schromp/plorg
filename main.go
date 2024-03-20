package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"

	"plorg/.gen/defaultdb/public/model"
	. "plorg/.gen/defaultdb/public/table"

	crdbpgx "github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	// "plorg/.gen/defaultdb/public/model"

	. "github.com/go-jet/jet/v2/postgres"
)

type User struct {
	Name string
	Age  int
}

func initTable(ctx context.Context, tx pgx.Tx) error {
	// Dropping existing table if it exists
	log.Println("Drop existing accounts table if necessary.")
	if _, err := tx.Exec(ctx, "DROP TABLE IF EXISTS accounts"); err != nil {
		return err
	}

	// Create the accounts table
	log.Println("Creating accounts table.")
	if _, err := tx.Exec(ctx,
		"CREATE TABLE accounts (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), balance INT8)"); err != nil {
		return err
	}
	return nil
}

func main() {

	os.Setenv("DATABASE_URL", "postgresql://root@127.0.0.1:26257/defaultdb?sslmode=disable")

	// Database
	// Setup pgx config
	config, err := pgx.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	config.RuntimeParams["application_name"] = "$ docs_simplecrud_gopgx"

	// Connect to db
	conn, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatal(err)
	}
	// Close the db connection when app exits
	defer conn.Close(context.Background())

	// Create a table using a transaction
	err = crdbpgx.ExecuteTx(context.Background(), conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		return initTable(context.Background(), tx)
	})

	// Jet things

	// var accounts []struct {
	// 	model.Accounts
	// }

	meh := int64(100)
	accStr := model.Accounts{
		ID:      uuid.New(),
		Balance: &meh,
	}

	newAccount, newArgs := Accounts.INSERT(Accounts.ID, Accounts.Balance).MODEL(accStr).Sql()

	err = crdbpgx.ExecuteTx(context.Background(), conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		log.Println("Creating new rows...")
		ret, err := tx.Exec(context.Background(), newAccount, newArgs...)
		log.Println(ret)
		return err
	})
	if err != nil {
		log.Fatal(err)
	}

	account := SELECT(Accounts.AllColumns).FROM(Accounts)
	sql, args := account.Sql()

	rows, err := conn.Query(context.Background(), sql, args...)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(rows)

	// Webserver
	tmpl, err := template.New("index.html").ParseFiles("./templates/index.html")
	if err != nil {
		log.Fatal(err)
	}

	innertmpl, ierr := template.New("inner").Parse(`<h3> Inner template: {{.Age}}</h3>`)
	if ierr != nil {
		log.Fatal(ierr)
	}

	user := User{Name: "John", Age: 31}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err = tmpl.Execute(w, user)
		if err != nil {
			log.Fatal(err)
		}
	})

	http.HandleFunc("/inner", func(w http.ResponseWriter, r *http.Request) {
		err = innertmpl.Execute(w, user)
	})

	log.Println("Server started on: http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
