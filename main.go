package main

import (
	// "fmt"
	"html/template"
	"log"
	"net/http"
)

type User struct {
	Name string
	Age  int
}

func main() {

	tmpl, err := template.New("index.html").ParseFiles("./templates/index.html")
	if err != nil {
		log.Fatal(err)
	}

	innertmpl, ierr := template.New("inner").Parse(`<h3> Inner template: {{.Age}}</h3>`)
	if ierr != nil {
		log.Fatal(ierr)
	}

	user := User{Name: "John", Age: 30}

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
