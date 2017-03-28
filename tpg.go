package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type StopDB struct {
	NameMatching map[string]string
}

func NewStopDB() (*StopDB, error) {

	resp, err := http.Get("http://prod.ivtr-od.tpg.ch/v1/GetStops.json?key=" + os.Getenv("TPG_KEY"))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	value, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println("JSON:", string(value))
	return &StopDB{}, nil
}
