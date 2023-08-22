package main

import "context"

type StorerStub struct {
	allPeopleStub    func(ctx context.Context) ([]Person, error)
	personForIDStub  func(ctx context.Context, id int) (*Person, error)
	addPersonStub    func(ctx context.Context, p Person) (Person, error)
	deletePersonStub func(ctx context.Context, id int) error
	updatePersonStub func(ctx context.Context, id int, p Person) (Person, error)
}

func (ss StorerStub) allPeople(ctx context.Context) ([]Person, error) {
	return ss.allPeopleStub(ctx)
}

func (ss StorerStub) personForID(ctx context.Context, id int) (*Person, error) {
	return ss.personForIDStub(ctx, id)
}

func (ss StorerStub) addPerson(ctx context.Context, p Person) (Person, error) {
	return ss.addPersonStub(ctx, p)
}

func (ss StorerStub) deletePerson(ctx context.Context, id int) error {
	return ss.deletePersonStub(ctx, id)
}

func (ss StorerStub) updatePerson(ctx context.Context, id int, p Person) (Person, error) {
	return ss.updatePersonStub(ctx, id, p)
}
