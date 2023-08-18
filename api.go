package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// TODO - change to environement config
var databaseURL string = "postgresql://myuser:mypassword@postgres:5432/mydb"

func main() {
	actx := AppContext{}

	// probably not the preferred pattern for depencency injection but this is fine for now
	if os.Args[0] == "test" {
		ms := NewMemoryStore(0) // add a count to lengthen process time by seconds
		actx.storer = &ms
	} else {
		ps := NewPostgresStore(databaseURL)
		actx.storer = &ps
		close := ps.startDatabase()
		defer close()
	}

	handleRequests(actx)
	fmt.Println("-------------------- main is over")
}

type responseError struct {
	Error string `json:"error"`
}

type AppContext struct {
	storer Storer
}

func handleRequests(actx AppContext) {
	sm := http.NewServeMux()
	wm := middleWare(sm)

	sm.Handle("/people", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { handlePeople(&actx, w, r) }))
	sm.Handle("/people/", http.StripPrefix("/people/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { handlePerson(&actx, w, r) })))

	s := http.Server{
		Addr:         ":8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      wm,
	}

	log.Fatal(s.ListenAndServe())
}

func middleWare(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		r = r.WithContext(ctx)
		w.Header().Add("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	})
}

func handlePeople(actx *AppContext, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handlePeopleGET(actx, w, r)
	case "POST":
		handlePeoplePOST(actx, w, r)
	}
}

func handlePerson(actx *AppContext, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handlePersonGET(actx, w, r)
	case "PATCH":
		handlePersonPATCH(actx, w, r)
	case "PUT":
		handlePersonPUT(actx, w, r)
	case "DELETE":
		handlePersonDELETE(actx, w, r)
	}
}

func handlePeopleGET(actx *AppContext, w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	people, err := actx.storer.allPeople(ctx)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, responseError{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, people)
}

func handlePersonGET(actx *AppContext, w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, responseError{Error: "invalid request"})
		return
	}

	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	person, err := actx.storer.personForID(ctx, id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, responseError{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, *person)
}

func handlePeoplePOST(actx *AppContext, w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var p Person
	if err := decoder.Decode(&p); err != nil {
		writeJSON(w, http.StatusBadRequest, responseError{Error: err.Error()})
		return
	}

	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	pu, err := actx.storer.addPerson(ctx, p)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, responseError{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, pu)
}

func handlePersonPUT(actx *AppContext, w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, responseError{Error: "invalid request"})
		return
	}

	decoder := json.NewDecoder(r.Body)
	var p Person
	if err := decoder.Decode(&p); err != nil {
		writeJSON(w, http.StatusBadRequest, responseError{Error: err.Error()})
		return
	}

	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	pu, err := actx.storer.updatePerson(ctx, id, p)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, responseError{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, pu)
}

func handlePersonDELETE(actx *AppContext, w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, responseError{Error: "invalid request"})
		return
	}

	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := actx.storer.deletePerson(ctx, id); err != nil {
		writeJSON(w, http.StatusBadRequest, responseError{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

// TODO - patch is difficult because it's hard to determine which values
// are actually passed - you can encode into a map and then convert based on the
// keys, but is that the most efficient way to handle this?
func handlePersonPATCH(actx *AppContext, w http.ResponseWriter, r *http.Request) {
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(value)
}
