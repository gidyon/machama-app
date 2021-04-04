package models

import (
	"fmt"
	"time"

	"github.com/gidyon/machama-app/pkg/api/loan"
	"github.com/gidyon/micro/v2/utils/errs"
)

// ProductId              string  `protobuf:"bytes,1,opt,name=plan_id,json=planId,proto3" json:"plan_id,omitempty"`
// ChamaId             string  `protobuf:"bytes,2,opt,name=chama_id,json=chamaId,proto3" json:"chama_id,omitempty"`
// Name                string  `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
// Description         string  `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
// InterestRate        string  `protobuf:"bytes,5,opt,name=interest_rate,json=interestRate,proto3" json:"interest_rate,omitempty"`
// LoanDurationDays    int32   `protobuf:"varint,6,opt,name=loan_duration_days,json=loanDurationDays,proto3" json:"loan_duration_days,omitempty"`
// LoanMinimumAmount   float32 `protobuf:"fixed32,7,opt,name=loan_minimum_amount,json=loanMinimumAmount,proto3" json:"loan_minimum_amount,omitempty"`
// LoanMaximumAmount   float32 `protobuf:"fixed32,8,opt,name=loan_maximun_amount,json=loanMaximunAmount,proto3" json:"loan_maximun_amount,omitempty"`
// LoanAccountBalance  float32 `protobuf:"fixed32,9,opt,name=loan_account_balance,json=loanAccountBalance,proto3" json:"loan_account_balance,omitempty"`
// LoanInterestBalance float32 `protobuf:"fixed32,10,opt,name=loan_interest_balance,json=loanInterestBalance,proto3" json:"loan_interest_balance,omitempty"`
// LoanSettledBalance  float32 `protobuf:"fixed32,11,opt,name=loan_settled_balance,json=loanSettledBalance,proto3" json:"loan_settled_balance,omitempty"`
// SettledLoan         int32   `protobuf:"varint,12,opt,name=settled_loan,json=settledLoan,proto3" json:"settled_loan,omitempty"`
// ActiveLoans         int32   `protobuf:"varint,13,opt,name=active_loans,json=activeLoans,proto3" json:"active_loans,omitempty"`
// TotalLoans          int32   `protobuf:"varint,14,opt,name=total_loans,json=totalLoans,proto3" json:"total_loans,omitempty"`
// CreateDate          string  `protobuf:"bytes,15,opt,name=create_date,json=createDate,proto3" json:"create_date,omitempty"`

type LoanProduct struct {
	ID                  uint      `gorm:"primaryKey;autoIncrement"`
	ChamaID             string    `gorm:"type:varchar(15);not null"`
	Name                string    `gorm:"type:varchar(50);not null"`
	Description         string    `gorm:"type:varchar(200)"`
	InterestRate        float32   `gorm:"type:float(3)"`
	LoanDurationDays    int32     `gorm:"type:int(10)"`
	LoanMinimumAmount   float64   `gorm:"type:float(15)"`
	LoanMaximumAmount   float64   `gorm:"type:float(15)"`
	LoanAccountBalance  float64   `gorm:"type:float(15)"`
	LoanInterestBalance float64   `gorm:"type:float(15)"`
	LoanSettledBalance  float64   `gorm:"type:float(15)"`
	SettledLoans        int32     `gorm:"type:int(10)"`
	ActiveLoans         int32     `gorm:"type:int(10)"`
	TotalLoans          int32     `gorm:"type:int(10)"`
	UpdatedAt           time.Time `gorm:"autoUpdateTime"`
	CreatedAt           time.Time `gorm:"autoCreateTime"`
}

func LoanProductModel(pb *loan.LoanProduct) (*LoanProduct, error) {
	if pb == nil {
		return nil, errs.NilObject("loan plan")
	}
	db := &LoanProduct{
		ChamaID:             pb.ChamaId,
		Name:                pb.Name,
		Description:         pb.Description,
		InterestRate:        pb.InterestRate,
		LoanDurationDays:    pb.LoanDurationDays,
		LoanMinimumAmount:   pb.LoanMinimumAmount,
		LoanMaximumAmount:   pb.LoanMaximumAmount,
		LoanAccountBalance:  pb.LoanAccountBalance,
		LoanInterestBalance: pb.LoanInterestBalance,
		LoanSettledBalance:  pb.LoanSettledBalance,
		SettledLoans:        pb.SettledLoans,
		ActiveLoans:         pb.ActiveLoans,
		TotalLoans:          pb.TotalLoans,
	}
	return db, nil
}

func LoanProductProto(db *LoanProduct) (*loan.LoanProduct, error) {
	if db == nil {
		return nil, errs.NilObject("loan plan")
	}
	pb := &loan.LoanProduct{
		ProductId:           fmt.Sprint(db.ID),
		ChamaId:             db.ChamaID,
		Name:                db.Name,
		Description:         db.Description,
		LoanDurationDays:    db.LoanDurationDays,
		InterestRate:        db.InterestRate,
		LoanMinimumAmount:   db.LoanMinimumAmount,
		LoanMaximumAmount:   db.LoanMaximumAmount,
		LoanAccountBalance:  db.LoanAccountBalance,
		LoanInterestBalance: db.LoanInterestBalance,
		LoanSettledBalance:  db.LoanSettledBalance,
		SettledLoans:        db.SettledLoans,
		ActiveLoans:         db.ActiveLoans,
		TotalLoans:          db.TotalLoans,
		UpdatedDate:         db.UpdatedAt.String(),
		CreatedDate:         db.CreatedAt.String(),
	}
	return pb, nil
}
