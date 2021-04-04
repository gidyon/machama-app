package models

import (
	"fmt"
	"time"

	"github.com/gidyon/machama-app/pkg/api/transaction"
	"github.com/gidyon/micro/v2/utils/errs"
)

// AccountId            string      `protobuf:"bytes,1,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"`
// OwnerId              string      `protobuf:"bytes,2,opt,name=owner_id,json=ownerId,proto3" json:"owner_id,omitempty"`
// AccountName          string      `protobuf:"bytes,3,opt,name=account_name,json=accountName,proto3" json:"account_name,omitempty"`
// AccountType          AccountType `protobuf:"varint,4,opt,name=account_type,json=accountType,proto3,enum=gidyon.transaction.AccountType" json:"account_type,omitempty"`
// Withdrawable         bool        `protobuf:"varint,5,opt,name=withdrawable,proto3" json:"withdrawable,omitempty"`
// TotalAmountDeposited float64     `protobuf:"fixed64,6,opt,name=total_amount_deposited,json=totalAmountDeposited,proto3" json:"total_amount_deposited,omitempty"`
// AvailableAmount      float64     `protobuf:"fixed64,7,opt,name=available_amount,json=availableAmount,proto3" json:"available_amount,omitempty"`
// TotalWithdrawnAmount      float64     `protobuf:"fixed64,8,opt,name=withdrawn_amount,json=withdrawnAmount,proto3" json:"withdrawn_amount,omitempty"`
// LastDepositedAmount  float64     `protobuf:"fixed64,9,opt,name=last_deposited_amount,json=lastDepositedAmount,proto3" json:"last_deposited_amount,omitempty"`
// LastTotalWithdrawnAmount  float64     `protobuf:"fixed64,10,opt,name=last_withdrawn_amount,json=lastTotalWithdrawnAmount,proto3" json:"last_withdrawn_amount,omitempty"`
// CreatedDate          string      `protobuf:"bytes,11,opt,name=created_date,json=createdDate,proto3" json:"created_date,omitempty"`
// UpdatedDate          string      `protobuf:"bytes,12,opt,name=updated_date,json=updatedDate,proto3" json:"updated_date,omitempty"`

type ChamaAccount struct {
	ID                   uint      `gorm:"primaryKey;autoIncrement"`
	OwnerID              string    `gorm:"type:varchar(50);not null"`
	AccountName          string    `gorm:"type:varchar(50);not null"`
	AccountType          string    `gorm:"type:varchar(30);not null"`
	Withdrawable         bool      `gorm:"type:tinyint(1)"`
	TotalDepositedAmount float64   `gorm:"type:float(15)"`
	AvailableAmount      float64   `gorm:"type:float(15)"`
	TotalWithdrawnAmount float64   `gorm:"type:float(15)"`
	LastDepositedAmount  float64   `gorm:"type:float(15)"`
	LastWithdrawnAmount  float64   `gorm:"type:float(15)"`
	Active               bool      `gorm:"type:tinyint(1)"`
	UpdatedAt            time.Time `gorm:"autoUpdateTime"`
	CreatedAt            time.Time `gorm:"autoCreateTime"`
}

func (*ChamaAccount) TableName() string {
	return "chama_accounts"
}

func ChamaAccountModel(pb *transaction.ChamaAccount) (*ChamaAccount, error) {
	if pb == nil {
		return nil, errs.NilObject("chama account")
	}
	db := &ChamaAccount{
		OwnerID:              pb.OwnerId,
		AccountName:          pb.AccountName,
		AccountType:          pb.AccountType.String(),
		Withdrawable:         pb.Withdrawable,
		TotalDepositedAmount: pb.TotalDepositedAmount,
		AvailableAmount:      pb.AvailableAmount,
		TotalWithdrawnAmount: pb.TotalWithdrawnAmount,
		LastDepositedAmount:  pb.LastDepositedAmount,
		LastWithdrawnAmount:  pb.LastWithdrawnAmount,
		Active:               pb.Active,
	}
	return db, nil
}

func ChamaAccountProto(db *ChamaAccount) (*transaction.ChamaAccount, error) {
	if db == nil {
		return nil, errs.NilObject("chama account")
	}
	pb := &transaction.ChamaAccount{
		AccountId:            fmt.Sprint(db.ID),
		OwnerId:              db.OwnerID,
		AccountName:          db.AccountName,
		AccountType:          transaction.AccountType(transaction.AccountType_value[db.AccountType]),
		Withdrawable:         db.Withdrawable,
		AvailableAmount:      db.AvailableAmount,
		TotalDepositedAmount: db.TotalDepositedAmount,
		TotalWithdrawnAmount: db.TotalWithdrawnAmount,
		LastDepositedAmount:  db.LastDepositedAmount,
		LastWithdrawnAmount:  db.LastWithdrawnAmount,
		Active:               db.Active,
		CreatedDate:          db.CreatedAt.String(),
		UpdatedDate:          db.CreatedAt.String(),
	}
	return pb, nil
}
