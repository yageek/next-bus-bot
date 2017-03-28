package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type Stop struct {
	Name string `json:"stopName"`
	Code string `json:"stopCode"`
}

type stopRecord struct {
	Stops []Stop `json:"stops"`
}

type StopDB struct {
	NameMatching map[string]string
	NameList     []string
}

func NewStopDB() (*StopDB, error) {

	resp, err := http.Get("http://prod.ivtr-od.tpg.ch/v1/GetStops.json?key=" + os.Getenv("TPG_KEY"))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var record stopRecord
	err = json.NewDecoder(resp.Body).Decode(&record)
	if err != nil {
		return nil, err
	}
	log.Println("Donwloading DATA OK")

	nameMap := make(map[string]string)
	nameList := make([]string, len(record.Stops))

	for i, stop := range record.Stops {
		nameMap[stop.Code] = stop.Name
		nameList[i] = stop.Name
	}
	log.Println("Database OK")
	return &StopDB{
		NameMatching: nameMap,
		NameList:     nameList,
	}, nil
}
