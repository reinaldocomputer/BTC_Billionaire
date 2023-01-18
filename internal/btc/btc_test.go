package btc

import (
	"fmt"
	"github.com/reinaldocomputer/BTC_Billionaire/internal/platform/mongodb"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
	"testing"
	"time"
)

type MockFinderStorer struct {
	CallbackStore func(data interface{}) error
	CallbackFind  func(start time.Time, end time.Time) ([]bson.M, error)
}

func (m *MockFinderStorer) Store(data interface{}) error {
	if m.CallbackStore != nil {
		return m.CallbackStore(data)
	}
	return nil

}

func (m *MockFinderStorer) Find(start time.Time, end time.Time) ([]bson.M, error) {
	if m.CallbackFind != nil {
		return m.CallbackFind(start, end)
	}
	return []bson.M{}, nil
}

func TestHistory_GetHistory(t1 *testing.T) {
	type fields struct {
		StartDateTime time.Time
		EndDateTime   time.Time
		DataBase      Finder
	}
	tests := []struct {
		name    string
		fields  fields
		want    []SendBTCRequest
		wantErr error
	}{
		{
			name: "Success",
			fields: fields{
				StartDateTime: time.Date(2022, 1, 22, 8, 0, 0, 0, time.UTC),
				EndDateTime:   time.Date(2022, 1, 22, 9, 0, 0, 0, time.UTC),
				DataBase: &MockFinderStorer{
					CallbackFind: func(start time.Time, end time.Time) ([]bson.M, error) {
						return []bson.M{
							{
								"datetime": time.Date(2022, 1, 22, 8, 30, 0, 0,
									time.UTC),
								"amount": 1000.00,
							},
						}, nil
					},
				},
			},
			want: []SendBTCRequest{
				{
					DateTime: time.Date(2022, 1, 22, 8, 30, 0, 0, time.UTC),
					Amount:   1000.00,
				},
			},
			wantErr: nil,
		},
		{
			name: "Find Error",
			fields: fields{
				StartDateTime: time.Date(2022, 1, 22, 8, 0, 0, 0, time.UTC),
				EndDateTime:   time.Date(2022, 1, 22, 9, 0, 0, 0, time.UTC),
				DataBase: &MockFinderStorer{
					CallbackFind: func(start time.Time, end time.Time) ([]bson.M, error) {
						return []bson.M{}, fmt.Errorf("Finder Error")
					},
				},
			},
			want:    nil,
			wantErr: fmt.Errorf("Finder Error"),
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &History{
				StartDateTime: tt.fields.StartDateTime,
				EndDateTime:   tt.fields.EndDateTime,
				DataBase:      tt.fields.DataBase,
			}
			got, err := t.GetHistory()
			assert.EqualValues(t1, tt.want, got)
			assert.EqualValues(t1, tt.wantErr, err)
		})
	}
}

func TestHistory_Valid(t1 *testing.T) {
	type fields struct {
		StartDateTime time.Time
		EndDateTime   time.Time
		DataBase      Finder
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				StartDateTime: time.Date(2022, 1, 22, 8, 0, 0, 0, time.UTC),
				EndDateTime:   time.Date(2022, 1, 22, 9, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name: "Zeror Error",
			fields: fields{
				EndDateTime: time.Date(2022, 1, 22, 9, 0, 0, 0, time.UTC),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &History{
				StartDateTime: tt.fields.StartDateTime,
				EndDateTime:   tt.fields.EndDateTime,
				DataBase:      tt.fields.DataBase,
			}
			if err := t.Valid(); (err != nil) != tt.wantErr {
				t1.Errorf("Valid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewHistory(t *testing.T) {
	type args struct {
		data HistoryRequest
	}
	tests := []struct {
		name string
		args args
		want *History
	}{
		{
			name: "Success NewHistory",
			args: args{
				data: HistoryRequest{
					StartDateTime: time.Date(2022, 1, 22, 8, 0, 0, 0, time.UTC),
					EndDateTime:   time.Date(2022, 1, 22, 9, 0, 0, 0, time.UTC),
				},
			},
			want: &History{
				StartDateTime: time.Date(2022, 1, 22, 8, 0, 0, 0, time.UTC),
				EndDateTime:   time.Date(2022, 1, 22, 9, 0, 0, 0, time.UTC),
				DataBase:      &mongodb.MongoDB{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHistory(tt.args.data); !assert.EqualValues(t, got, tt.want) {
				t.Errorf("NewHistory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewTransaction(t *testing.T) {
	type args struct {
		data SendBTCRequest
	}
	tests := []struct {
		name string
		args args
		want *Transaction
	}{
		{
			name: "New Transaction Success",
			args: args{
				data: SendBTCRequest{
					DateTime: time.Date(2022, 1, 22, 8, 0, 0, 0, time.UTC),
					Amount:   500,
				},
			},
			want: &Transaction{
				DateTime: time.Date(2022, 1, 22, 8, 0, 0, 0, time.UTC),
				Amount:   500,
				DataBase: &mongodb.MongoDB{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTransaction(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransaction_SendBTC(t1 *testing.T) {
	type fields struct {
		DateTime time.Time
		Amount   float64
		DataBase Storer
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "SendBTC Success",
			fields: fields{
				DateTime: time.Date(2022, 1, 22, 8, 0, 0, 0, time.UTC),
				Amount:   100,
				DataBase: &MockFinderStorer{
					CallbackStore: func(data interface{}) error {
						return nil
					},
				},
			},
			wantErr: false,
		},
		{
			name: "SendBTC Error",
			fields: fields{
				DateTime: time.Date(2022, 1, 22, 8, 0, 0, 0, time.UTC),
				Amount:   100,
				DataBase: &MockFinderStorer{
					CallbackStore: func(data interface{}) error {
						return fmt.Errorf("Error")
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Transaction{
				DateTime: tt.fields.DateTime,
				Amount:   tt.fields.Amount,
				DataBase: tt.fields.DataBase,
			}
			if err := t.SendBTC(); (err != nil) != tt.wantErr {
				t1.Errorf("SendBTC() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTransaction_Valid(t1 *testing.T) {
	type fields struct {
		DateTime time.Time
		Amount   float64
		DataBase Storer
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Transaction Valid Success",
			fields: fields{
				DateTime: time.Date(2022, 1, 22, 8, 0, 0, 0, time.UTC),
				Amount:   100,
			},
			wantErr: false,
		},
		{
			name: "Transaction Valid Zero Error",
			fields: fields{
				Amount: 100,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Transaction{
				DateTime: tt.fields.DateTime,
				Amount:   tt.fields.Amount,
				DataBase: tt.fields.DataBase,
			}
			if err := t.Valid(); (err != nil) != tt.wantErr {
				t1.Errorf("Valid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
