package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	ErrNoNextDepartures = errors.New("No planned departures")
)

type Departure struct {
	WaitingTime          int `json:"waitingTimeMillis"`
	ConectionWaitingTime int `json:"conectionWaitingTime"`
}
type departureRecord struct {
	Departures []Departure `json:"departures"`
}
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
		nameMap[stop.Name] = stop.Code
		nameList[i] = stop.Name
	}
	log.Println("Database OK")
	return &StopDB{
		NameMatching: nameMap,
		NameList:     nameList,
	}, nil
}

func (s *StopDB) getStopCode(q string) (string, error) {
	for _, value := range s.NameList {
		if value == q {
			return s.NameMatching[value], nil
		}
	}
	return "", errors.New("Name not found")
}
func (s *StopDB) GetNextBus(stopCode string) (time.Duration, error) {

	resp, err := http.Get("http://prod.ivtr-od.tpg.ch/v1/GetNextDepartures.json?key=" + os.Getenv("TPG_KEY") + "&stopCode=" + stopCode)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()
	var record departureRecord
	err = json.NewDecoder(resp.Body).Decode(&record)
	if err != nil {
		return -1, err
	}

	if len(record.Departures) < 1 {
		return -1, ErrNoNextDepartures
	}

	waitingTime := time.Duration(record.Departures[0].WaitingTime) * time.Millisecond
	return waitingTime, nil
}
