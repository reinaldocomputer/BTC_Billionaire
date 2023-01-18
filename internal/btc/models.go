package btc

import (
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type Storer interface {
	Store(data interface{}) error
}

type Finder interface {
	Find(start time.Time, end time.Time) ([]bson.M, error)
}

type SendBTCRequest struct {
	DateTime time.Time `json:"datetime" mongodb:"datetime"`
	Amount   float64   `json:"amount" mongodb:"amount"`
}

type Transaction struct {
	DateTime time.Time `json:"datetime" mongodb:"datetime"`
	Amount   float64   `json:"amount" mongodb:"amount"`
	DataBase Storer
}

type HistoryRequest struct {
	StartDateTime time.Time `json:"startDatetime"`
	EndDateTime   time.Time `json:"endDatetime"`
	DataBase      Finder
}

type History struct {
	StartDateTime time.Time `json:"startDatetime"`
	EndDateTime   time.Time `json:"endDatetime"`
	DataBase      Finder
}
