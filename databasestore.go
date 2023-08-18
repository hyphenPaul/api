package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	dbURL string
	pool  *pgxpool.Pool
}

func NewPostgresStore(dbURL string) PostgresStore {
	return PostgresStore{
		dbURL: dbURL,
	}
}

func (ps *PostgresStore) startDatabase() func() {
	fmt.Println("Starting the database")
	dbpool, err := pgxpool.New(context.Background(), ps.dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool %v\n", err)
		os.Exit(1)
	}
	ps.pool = dbpool
	return func() {
		dbpool.Close()
	}
}

func (ps PostgresStore) allPeople(ctx context.Context) ([]Person, error) {
	q := `
  SELECT id, lastname, firstname, age
  FROM people
  ORDER BY id DESC
  `
	rows, err := ps.pool.Query(ctx, q)
	if err != nil {
		return []Person{}, err
	}
	defer rows.Close()

	res := []Person{}
	for rows.Next() {
		var p Person
		if err = rows.Scan(&p.ID, &p.LastName, &p.FirstName, &p.Age); err != nil {
			return []Person{}, err
		}

		res = append(res, p)
	}

	return res, nil
}

func (ps PostgresStore) personForID(ctx context.Context, id int) (*Person, error) {
	q := `
  SELECT id, lastname, firstname, age
  FROM people
  WHERE id = $1
  `
	var p Person
	row := ps.pool.QueryRow(ctx, q, id)

	if err := row.Scan(&p.ID, &p.LastName, &p.FirstName, &p.Age); err != nil {
		return nil, err
	}

	return &p, nil
}

func (ps PostgresStore) addPerson(ctx context.Context, p Person) (Person, error) {
	q := `
  INSERT INTO people (lastname, firstname, age)
  VALUES ($1, $2, $3)
  RETURNING id
  `
	var id int
	row := ps.pool.QueryRow(ctx, q, p.FirstName, p.LastName, p.Age)
	if err := row.Scan(&id); err != nil {
		return p, err
	}

	p.ID = id
	return p, nil
}

func (ps PostgresStore) deletePerson(ctx context.Context, id int) error {
	q := `
  DELETE FROM people
  WHERE id = $1
  `
	_, err := ps.pool.Exec(ctx, q, id)
	if err != nil {
		return err
	}

	return nil
}

func (ps PostgresStore) updatePerson(ctx context.Context, id int, p Person) (Person, error) {
	fmt.Printf("%+v\n", p)
	q := `
  UPDATE people
  SET firstname=$1, lastname=$2, age=$3
  WHERE id = $4
  RETURNING id, firstname, lastname, age
  `
	row := ps.pool.QueryRow(ctx, q, p.FirstName, p.LastName, p.Age, id)
	var up Person
	if err := row.Scan(&up.ID, &up.LastName, &up.FirstName, &up.Age); err != nil {
		return p, err
	}

	return up, nil
}
