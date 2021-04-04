package models

import (
	"fmt"
	"time"

	"github.com/gidyon/machama-app/pkg/api/transaction"
	"github.com/gidyon/micro/v2/utils/errs"
)

// TransactionId          string          `protobuf:"bytes,1,opt,name=transaction_id,json=transactionId,proto3" json:"transaction_id,omitempty"`
// ActorId                string          `protobuf:"bytes,2,opt,name=actor_id,json=actorId,proto3" json:"actor_id,omitempty"`
// AccountId              string          `protobuf:"bytes,3,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"`
// Description            string          `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
// TransactionType        TransactionType `protobuf:"varint,5,opt,name=transaction_type,json=transactionType,proto3,enum=gidyon.transaction.TransactionType" json:"transaction_type,omitempty"`
// TransactionAmount      float64         `protobuf:"fixed64,6,opt,name=transaction_amount,json=transactionAmount,proto3" json:"transaction_amount,omitempty"`
// TransactionTimeSeconds int64           `protobuf:"varint,7,opt,name=transaction_time_seconds,json=transactionTimeSeconds,proto3" json:"transaction_time_seconds,omitempty"`

type Transaction struct {
	ID                uint      `gorm:"primaryKey;autoIncrement"`
	ActorID           string    `gorm:"type:varchar(50);not null"`
	AccountID         string    `gorm:"type:varchar(50);not null"`
	Description       string    `gorm:"type:varchar(200);not null"`
	TransactionType   string    `gorm:"type:varchar(30);not null"`
	TransactionAmount float64   `gorm:"type:float(15)"`
	CreatedAt         time.Time `gorm:"autoCreateTime"`
}

func (*Transaction) TableName() string {
	return "transactions"
}

func TransactionModel(pb *transaction.Transaction) (*Transaction, error) {
	if pb == nil {
		return nil, errs.NilObject("transaction")
	}
	db := &Transaction{
		ActorID:           pb.ActorId,
		AccountID:         pb.AccountId,
		Description:       pb.Description,
		TransactionType:   pb.TransactionType.String(),
		TransactionAmount: pb.TransactionAmount,
	}
	return db, nil
}

func TransactionProto(db *Transaction) (*transaction.Transaction, error) {
	if db == nil {
		return nil, errs.NilObject("transaction")
	}
	pb := &transaction.Transaction{
		TransactionId:          fmt.Sprint(db.ID),
		ActorId:                db.ActorID,
		AccountId:              db.AccountID,
		Description:            db.Description,
		TransactionType:        transaction.TransactionType(transaction.TransactionType_value[db.TransactionType]),
		TransactionAmount:      db.TransactionAmount,
		TransactionTimeSeconds: db.CreatedAt.Unix(),
	}
	return pb, nil
}
