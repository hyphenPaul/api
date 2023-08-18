package main

import (
	"context"
	"fmt"
	"time"
)

type MemoryStore struct {
	people       []Person
	sleepSeconds int
}

func NewMemoryStore(sleepSeconds int) MemoryStore {
	people := []Person{
		{ID: 1, FirstName: "Bob", LastName: "Barker", Age: 53},
		{ID: 2, FirstName: "Fred", LastName: "Flintstone", Age: 44},
		{ID: 3, FirstName: "Joan", LastName: "Jet", Age: 49},
	}

	return MemoryStore{people: people, sleepSeconds: sleepSeconds}
}

func (m MemoryStore) allPeople(ctx context.Context) ([]Person, error) {
	ch := make(chan []Person, 0)

	go func() {
		time.Sleep(time.Duration(m.sleepSeconds) * time.Second)
		ch <- m.people
	}()

	select {
	case people := <-ch:
		return people, nil
	case <-ctx.Done():
		return []Person{}, ctx.Err()
	}
}

func (m MemoryStore) personForID(ctx context.Context, id int) (*Person, error) {
	type ret struct {
		person *Person
		error  error
	}

	ch := make(chan ret, 1)

	go func() {
		time.Sleep(time.Duration(m.sleepSeconds) * time.Second)
		for _, person := range m.people {
			if person.ID == id {
				ch <- ret{&person, nil}
				return
			}
		}

		ch <- ret{nil, fmt.Errorf("Person not found for id: %d", id)}
	}()

	select {
	case ret := <-ch:
		return ret.person, ret.error
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (m *MemoryStore) addPerson(ctx context.Context, p Person) (Person, error) {
	type ret struct {
		person Person
		error  error
	}

	ch := make(chan ret, 1)

	go func() {
		time.Sleep(time.Duration(m.sleepSeconds) * time.Second)
		for _, i := range m.people {
			if i.ID == p.ID {
				ch <- ret{p, fmt.Errorf("The person ID is alread taken: %v", p)}
				return
			}
		}

		m.people = append(m.people, p)
		ch <- ret{p, nil}
	}()

	select {
	case ret := <-ch:
		return ret.person, ret.error
	case <-ctx.Done():
		return p, ctx.Err()
	}
}

func (m *MemoryStore) deletePerson(ctx context.Context, id int) error {
	ch := make(chan error, 1)

	go func() {
		time.Sleep(time.Duration(m.sleepSeconds) * time.Second)
		for i, p := range m.people {
			if p.ID == id {
				m.people = delete_at_index(m.people, i)
				ch <- nil
				return
			}
		}

		ch <- fmt.Errorf("No person exists for ID: %d", id)
	}()

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (m *MemoryStore) updatePerson(ctx context.Context, id int, p Person) (Person, error) {
	type ret struct {
		person Person
		error  error
	}

	ch := make(chan ret, 1)

	go func() {
		time.Sleep(time.Duration(m.sleepSeconds) * time.Second)
		for i, ep := range m.people {
			if ep.ID == id {
				p.ID = id
				m.people[i] = p
				ch <- ret{p, nil}
				return
			}
		}

		ch <- ret{p, fmt.Errorf("No person exists for ID: %d", id)}
	}()

	select {
	case ret := <-ch:
		return ret.person, ret.error
	case <-ctx.Done():
		return p, ctx.Err()
	}
}

func delete_at_index(people []Person, index int) []Person {
	return append(people[:index], people[(index+1):]...)
}
