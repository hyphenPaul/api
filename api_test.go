package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

func Test_handlePersonGET(t *testing.T) {
	expperson := Person{FirstName: "Foo", LastName: "Bar", Age: 22}
	ss := StorerStub{
		personForIDStub: func(ctx context.Context, id int) (*Person, error) {
			return &expperson, nil
		},
	}
	h := newTestHandler(ss)

	server := httptest.NewServer(h)
	defer server.Close()

	res, err := http.Get(server.URL + "/people/1")
	if err != nil {
		t.Errorf("error during http.Get: %s", err.Error())
		return
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("got status %d but expected %d", res.StatusCode, http.StatusOK)
		return
	}

	decoder := json.NewDecoder(res.Body)
	defer res.Body.Close()
	var p Person
	err = decoder.Decode(&p)
	if err != nil {
		t.Errorf("error during decode: %s", err.Error())
		return
	}

	if p != expperson {
		t.Errorf("got response %v but got %v", p, expperson)
		return
	}
}

func Test_handlePersonGETBadRequest(t *testing.T) {
	ss := StorerStub{
		personForIDStub: func(ctx context.Context, id int) (*Person, error) {
			return &Person{}, nil
		},
	}
	h := newTestHandler(ss)

	server := httptest.NewServer(h)
	defer server.Close()

	res, err := http.Get(server.URL + "/people/foo")
	if err != nil {
		t.Errorf("error during http.Get: %s", err.Error())
		return
	}

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("got status %d but expected %d", res.StatusCode, http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(res.Body)
	defer res.Body.Close()
	var r map[string]string
	err = decoder.Decode(&r)
	if err != nil {
		t.Errorf("error during decode: %s", err.Error())
		return
	}

	expmap := map[string]string{"error": "invalid request"}

	if !reflect.DeepEqual(r, expmap) {
		t.Errorf("got response %v but got %v", r, expmap)
		return
	}
}

func Test_handlePersonGETBadID(t *testing.T) {
	ss := StorerStub{
		personForIDStub: func(ctx context.Context, id int) (*Person, error) {
			return &Person{}, errors.New("Foo Bar")
		},
	}
	h := newTestHandler(ss)

	server := httptest.NewServer(h)
	defer server.Close()

	res, err := http.Get(server.URL + "/people/123")
	if err != nil {
		t.Errorf("error during http.Get: %s", err.Error())
		return
	}

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("got status %d but expected %d", res.StatusCode, http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(res.Body)
	defer res.Body.Close()
	var r map[string]string
	err = decoder.Decode(&r)
	if err != nil {
		t.Errorf("error during decode: %s", err.Error())
		return
	}

	expmap := map[string]string{"error": "Foo Bar"}

	if !reflect.DeepEqual(r, expmap) {
		t.Errorf("got response %v but got %v", r, expmap)
		return
	}
}

func Test_handlePersonGETTimeout(t *testing.T) {
	ss := StorerStub{
		personForIDStub: func(ctx context.Context, id int) (*Person, error) {
			type ret struct {
				person *Person
				error  error
			}

			ch := make(chan ret, 1)

			go func() {
				time.Sleep(time.Second)
				ch <- ret{person: &Person{}, error: nil}
			}()

			time.Sleep(time.Second)

			select {
			case ret := <-ch:
				return ret.person, ret.error
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		},
	}
	h := newTestHandler(ss)

	server := httptest.NewServer(h)
	defer server.Close()

	res, err := http.Get(server.URL + "/people/123")
	if err != nil {
		t.Errorf("error during http.Get: %s", err.Error())
		return
	}

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("got status %d but expected %d", res.StatusCode, http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(res.Body)
	defer res.Body.Close()
	var r map[string]string
	err = decoder.Decode(&r)
	if err != nil {
		t.Errorf("error during decode: %s", err.Error())
		return
	}

	expmap := map[string]string{"error": "context deadline exceeded"}

	if !reflect.DeepEqual(r, expmap) {
		t.Errorf("got response %v but got %v", r, expmap)
		return
	}
}

func Test_handlePeopleGET(t *testing.T) {
	exppeople := []Person{
		{FirstName: "Foo", LastName: "Bar", Age: 22},
		{FirstName: "Bin", LastName: "Baz", Age: 24},
	}
	ss := StorerStub{
		allPeopleStub: func(ctx context.Context) ([]Person, error) {
			return exppeople, nil
		},
	}
	h := newTestHandler(ss)

	server := httptest.NewServer(h)
	defer server.Close()

	res, err := http.Get(server.URL + "/people")
	if err != nil {
		t.Errorf("error during http.Get: %s", err.Error())
		return
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("got status %d but expected %d", res.StatusCode, http.StatusOK)
		return
	}

	decoder := json.NewDecoder(res.Body)
	defer res.Body.Close()
	var pl []Person
	err = decoder.Decode(&pl)
	if err != nil {
		t.Errorf("error during decode: %s", err.Error())
		return
	}

	if !reflect.DeepEqual(pl, exppeople) {
		t.Errorf("got response %v but got %v", pl, exppeople)
		return
	}
}

func Test_handlePeopleGETBadRequest(t *testing.T) {
	ss := StorerStub{
		allPeopleStub: func(ctx context.Context) ([]Person, error) {
			return []Person{}, errors.New("Something went wrong")
		},
	}
	h := newTestHandler(ss)

	server := httptest.NewServer(h)
	defer server.Close()

	res, err := http.Get(server.URL + "/people")
	if err != nil {
		t.Errorf("error during http.Get: %s", err.Error())
		return
	}

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("got status %d but expected %d", res.StatusCode, http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(res.Body)
	defer res.Body.Close()
	var r map[string]string
	err = decoder.Decode(&r)
	if err != nil {
		t.Errorf("error during decode: %s", err.Error())
		return
	}

	expmap := map[string]string{"error": "Something went wrong"}

	if !reflect.DeepEqual(r, expmap) {
		t.Errorf("got response %v but got %v", r, expmap)
		return
	}
}

func Test_handlePeoplePOST(t *testing.T) {
	expperson := Person{FirstName: "Foo", LastName: "Bar", Age: 22}
	ss := StorerStub{
		addPersonStub: func(ctx context.Context, p Person) (Person, error) {
			return p, nil
		},
	}
	h := newTestHandler(ss)

	server := httptest.NewServer(h)
	defer server.Close()

	b, err := json.Marshal(&expperson)
	if err != nil {
		t.Errorf("Error marshaling json: %s", err.Error())
		return
	}
	params := strings.NewReader(string(b))
	res, err := http.Post(server.URL+"/people", "application/json", params)
	if err != nil {
		t.Errorf("error during http.Get: %s", err.Error())
		return
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("got status %d but expected %d", res.StatusCode, http.StatusOK)
		return
	}

	var p Person
	decoder := json.NewDecoder(res.Body)
	defer res.Body.Close()
	err = decoder.Decode(&p)
	if err != nil {
		t.Errorf("error during decode: %s", err.Error())
		return
	}

	if !reflect.DeepEqual(p, expperson) {
		t.Errorf("got response %v but got %v", p, expperson)
		return
	}
}

func Test_handlePeoplePOSTBadRequest(t *testing.T) {
	ss := StorerStub{
		addPersonStub: func(ctx context.Context, p Person) (Person, error) {
			return Person{}, errors.New("Something went wrong")
		},
	}
	h := newTestHandler(ss)

	server := httptest.NewServer(h)
	defer server.Close()

	m := map[string]string{"foo": "bar"}
	b, err := json.Marshal(m)
	if err != nil {
		t.Errorf("Error marshaling json: %s", err.Error())
		return
	}
	params := strings.NewReader(string(b))
	res, err := http.Post(server.URL+"/people", "application/json", params)
	if err != nil {
		t.Errorf("error during http.Get: %s", err.Error())
		return
	}

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("got status %d but expected %d", res.StatusCode, http.StatusBadRequest)
		return
	}

	var resm map[string]string
	decoder := json.NewDecoder(res.Body)
	defer res.Body.Close()
	err = decoder.Decode(&resm)
	if err != nil {
		t.Errorf("error during decode: %s", err.Error())
		return
	}

	expm := map[string]string{"error": "Something went wrong"}

	if !reflect.DeepEqual(expm, resm) {
		t.Errorf("got response %v but got %v", expm, resm)
		return
	}
}

func Test_handlePersonsPUT(t *testing.T) {
	expperson := Person{FirstName: "Foo", LastName: "Bar", Age: 22}
	ss := StorerStub{
		updatePersonStub: func(ctx context.Context, id int, p Person) (Person, error) {
			return expperson, nil
		},
	}
	h := newTestHandler(ss)

	server := httptest.NewServer(h)
	defer server.Close()

	b, err := json.Marshal(&expperson)
	if err != nil {
		t.Errorf("Error marshaling json: %s", err.Error())
		return
	}
	params := strings.NewReader(string(b))
	req, err := http.NewRequest(http.MethodPut, server.URL+"/people/1", params)
	if err != nil {
		t.Errorf("error during http.Get: %s", err.Error())
		return
	}

	req.Header.Set("Content-Type", "application/json")
	cli := &http.Client{}

	res, err := cli.Do(req)
	if res.StatusCode != http.StatusOK {
		t.Errorf("got status %d but expected %d", res.StatusCode, http.StatusOK)
		return
	}

	var p Person
	decoder := json.NewDecoder(res.Body)
	defer res.Body.Close()
	err = decoder.Decode(&p)
	if err != nil {
		t.Errorf("error during decode: %s", err.Error())
		return
	}

	if !reflect.DeepEqual(p, expperson) {
		t.Errorf("got response %v but got %v", p, expperson)
		return
	}
}

func Test_handlePersonsPUTBadRequest(t *testing.T) {
	expperson := Person{FirstName: "Foo", LastName: "Bar", Age: 22}
	ss := StorerStub{
		updatePersonStub: func(ctx context.Context, id int, p Person) (Person, error) {
			return expperson, nil
		},
	}
	h := newTestHandler(ss)

	server := httptest.NewServer(h)
	defer server.Close()

	b, err := json.Marshal(&expperson)
	if err != nil {
		t.Errorf("Error marshaling json: %s", err.Error())
		return
	}
	params := strings.NewReader(string(b))
	req, err := http.NewRequest(http.MethodPut, server.URL+"/people/foo", params)
	if err != nil {
		t.Errorf("error during http.Get: %s", err.Error())
		return
	}

	req.Header.Set("Content-Type", "application/json")
	cli := &http.Client{}

	res, err := cli.Do(req)
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("got status %d but expected %d", res.StatusCode, http.StatusBadRequest)
		return
	}

	var m map[string]string
	decoder := json.NewDecoder(res.Body)
	defer res.Body.Close()
	err = decoder.Decode(&m)
	if err != nil {
		t.Errorf("error during decode: %s", err.Error())
		return
	}

	expm := map[string]string{"error": "invalid request"}

	if !reflect.DeepEqual(m, expm) {
		t.Errorf("got response %v but got %v", m, expm)
		return
	}
}

func Test_handlePersonDELETE(t *testing.T) {
	ss := StorerStub{
		deletePersonStub: func(ctx context.Context, id int) error {
			return nil
		},
	}
	h := newTestHandler(ss)

	server := httptest.NewServer(h)
	defer server.Close()

	req, err := http.NewRequest(http.MethodDelete, server.URL+"/people/1", nil)
	if err != nil {
		t.Errorf("New request error: %s", err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json")

	cli := &http.Client{}
	res, err := cli.Do(req)
	if err != nil {
		t.Errorf("Error during cli.Do: %s", err.Error())
		return
	}

	var m map[string]bool
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&m)
	if err != nil {
		t.Errorf("Error during decode: %s", err.Error())
		return
	}

	expm := map[string]bool{"success": true}
	if !reflect.DeepEqual(expm, m) {
		t.Errorf("expected response %v got %v", m, expm)
	}
}

func newTestHandler(ss Storer) http.Handler {
	actx := AppContext{
		storer:  ss,
		timeout: 30 * time.Millisecond,
		logger:  noopLogger{},
	}

	return NewHandler(actx)
}
