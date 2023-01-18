package btc

import (
	"fmt"
	"github.com/reinaldocomputer/BTC_Billionaire/internal/platform/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	DatetimeZeroError    = "DateTime is invalid"
	IncorrectAmountError = "Amount is invalid"
)

func (t *Transaction) Valid() error {
	t.DateTime = t.DateTime.UTC()
	if t.DateTime.IsZero() {
		return fmt.Errorf(DatetimeZeroError)
	}
	if t.Amount <= 0 {
		return fmt.Errorf(IncorrectAmountError)
	}
	return nil
}

func NewTransaction(data SendBTCRequest) *Transaction {
	return &Transaction{
		DateTime: data.DateTime,
		Amount:   data.Amount,
		DataBase: &mongodb.MongoDB{},
	}
}
func (t *Transaction) SendBTC() error {
	return t.DataBase.Store(*t)
}

func NewHistory(data HistoryRequest) *History {
	return &History{
		StartDateTime: data.StartDateTime,
		EndDateTime:   data.EndDateTime,
		DataBase:      &mongodb.MongoDB{},
	}
}
func (t *History) Valid() error {
	t.StartDateTime = t.StartDateTime.UTC()
	t.EndDateTime = t.EndDateTime.UTC()

	if t.StartDateTime.IsZero() || t.EndDateTime.IsZero() {
		return fmt.Errorf(DatetimeZeroError)
	}
	return nil
}
func (t *History) GetHistory() ([]SendBTCRequest, error) {
	result, err := t.DataBase.Find(t.StartDateTime, t.EndDateTime)
	if err != nil {
		return nil, err
	}
	var transaction SendBTCRequest
	var transactions []SendBTCRequest
	for _, element := range result {
		var resultBytes []byte
		resultBytes, err = bson.Marshal(element)
		if err != nil {
			return nil, err
		}
		err = bson.Unmarshal(resultBytes, &transaction)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}
