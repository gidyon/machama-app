package models

import (
	"fmt"
	"time"

	"github.com/gidyon/machama-app/pkg/api/loan"
	"github.com/gidyon/micro/v2/utils/errs"
)

// LoanId        string  `protobuf:"bytes,1,opt,name=loan_id,json=loanId,proto3" json:"loan_id,omitempty"`
// ChamaId       string  `protobuf:"bytes,2,opt,name=chama_id,json=chamaId,proto3" json:"chama_id,omitempty"`
// ProductId        string  `protobuf:"bytes,3,opt,name=plan_id,json=planId,proto3" json:"plan_id,omitempty"`
// MemberId      string  `protobuf:"bytes,4,opt,name=member_id,json=memberId,proto3" json:"member_id,omitempty"`
// LoaneeNames   string  `protobuf:"bytes,5,opt,name=loanee_names,json=loaneeNames,proto3" json:"loanee_names,omitempty"`
// LoaneePhone   string  `protobuf:"bytes,6,opt,name=loanee_phone,json=loaneePhone,proto3" json:"loanee_phone,omitempty"`
// LoaneeEmail   string  `protobuf:"bytes,7,opt,name=loanee_email,json=loaneeEmail,proto3" json:"loanee_email,omitempty"`
// NationalId    string  `protobuf:"bytes,8,opt,name=national_id,json=nationalId,proto3" json:"national_id,omitempty"`
// Approved      bool    `protobuf:"varint,9,opt,name=approved,proto3" json:"approved,omitempty"`
// Duration      int32   `protobuf:"varint,10,opt,name=duration,proto3" json:"duration,omitempty"`
// LoanAmount    float32 `protobuf:"fixed32,11,opt,name=loan_amount,json=loanAmount,proto3" json:"loan_amount,omitempty"`
// InterestRate  float32 `protobuf:"fixed32,12,opt,name=interest_rate,json=interestRate,proto3" json:"interest_rate,omitempty"`
// SettledAmount float32 `protobuf:"fixed32,13,opt,name=settled_amount,json=settledAmount,proto3" json:"settled_amount,omitempty"`
// PenaltyAmount float32 `protobuf:"fixed32,14,opt,name=penalty_amount,json=penaltyAmount,proto3" json:"penalty_amount,omitempty"`
// BorrowedDate  string  `protobuf:"bytes,15,opt,name=borrowed_date,json=borrowedDate,proto3" json:"borrowed_date,omitempty"`

type Loan struct {
	ID            uint      `gorm:"primaryKey;autoIncrement"`
	ChamaID       string    `gorm:"type:varchar(15);not null"`
	ProductID     string    `gorm:"type:varchar(15);not null"`
	MemberID      string    `gorm:"type:varchar(15);not null"`
	LoaneeNames   string    `gorm:"type:varchar(30);not null"`
	LoaneePhone   string    `gorm:"type:varchar(15);not null"`
	NationalID    string    `gorm:"type:varchar(10);not null"`
	LoaneeEmail   string    `gorm:"type:varchar(50)"`
	Approved      bool      `gorm:"type:tinyint(1)"`
	DurationDays  int32     `gorm:"type:int(10)"`
	InterestRate  float32   `gorm:"type:float(3)"`
	LoanAmount    float64   `gorm:"type:float(15)"`
	SettledAmount float64   `gorm:"type:float(15)"`
	PenaltyAmount float64   `gorm:"type:float(15)"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
}

func LoanModel(pb *loan.Loan) (*Loan, error) {
	if pb == nil {
		return nil, errs.NilObject("loan")
	}
	db := &Loan{
		ChamaID:       pb.ChamaId,
		ProductID:     pb.ProductId,
		MemberID:      pb.MemberId,
		LoaneeNames:   pb.LoaneeNames,
		LoaneePhone:   pb.LoaneePhone,
		LoaneeEmail:   pb.LoaneeEmail,
		NationalID:    pb.NationalId,
		Approved:      pb.Approved,
		DurationDays:  pb.DurationDays,
		LoanAmount:    pb.LoanAmount,
		InterestRate:  pb.InterestRate,
		SettledAmount: pb.SettledAmount,
		PenaltyAmount: pb.PenaltyAmount,
	}
	return db, nil
}

func LoanProto(db *Loan) (*loan.Loan, error) {
	if db == nil {
		return nil, errs.NilObject("loan")
	}
	pb := &loan.Loan{
		LoanId:        fmt.Sprint(db.ID),
		ChamaId:       db.ChamaID,
		ProductId:     db.ProductID,
		MemberId:      db.MemberID,
		LoaneeNames:   db.LoaneeNames,
		LoaneePhone:   db.LoaneePhone,
		LoaneeEmail:   db.LoaneeEmail,
		NationalId:    db.NationalID,
		Approved:      db.Approved,
		DurationDays:  db.DurationDays,
		LoanAmount:    db.LoanAmount,
		InterestRate:  db.InterestRate,
		SettledAmount: db.SettledAmount,
		PenaltyAmount: db.PenaltyAmount,
		UpdatedDate:   db.UpdatedAt.String(),
		BorrowedDate:  db.CreatedAt.String(),
	}
	return pb, nil
}
