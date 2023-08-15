package main

import "context"

type Storer interface {
	allPeople(ctx context.Context) ([]Person, error)
	personForID(ctx context.Context, id int) (*Person, error)
	addPerson(ctx context.Context, p Person) (Person, error)
	deletePerson(ctx context.Context, id int) error
	updatePerson(ctx context.Context, id int, p Person) (Person, error)
}
