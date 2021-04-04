package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gidyon/machama-app/pkg/api/chama"
	"github.com/gidyon/micro/v2/utils/errs"
)

// MemberId      string            `protobuf:"bytes,1,opt,name=member_id,json=memberId,proto3" json:"member_id,omitempty"`
// ChamaId       string            `protobuf:"bytes,2,opt,name=chama_id,json=chamaId,proto3" json:"chama_id,omitempty"`
// FirstName     string            `protobuf:"bytes,3,opt,name=first_name,json=firstName,proto3" json:"first_name,omitempty"`
// LastName      string            `protobuf:"bytes,4,opt,name=last_name,json=lastName,proto3" json:"last_name,omitempty"`
// Email         string            `protobuf:"bytes,5,opt,name=email,proto3" json:"email,omitempty"`
// Phone         string            `protobuf:"bytes,6,opt,name=phone,proto3" json:"phone,omitempty"`
// IdNumber      string            `protobuf:"bytes,7,opt,name=id_number,json=idNumber,proto3" json:"id_number,omitempty"`
// Residence     string            `protobuf:"bytes,8,opt,name=residence,proto3" json:"residence,omitempty"`
// JobDetails    map[string]string `protobuf:"bytes,9,rep,name=job_details,json=jobDetails,proto3" json:"job_details,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
// Kyc           map[string]string `protobuf:"bytes,10,rep,name=kyc,proto3" json:"kyc,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
// Beneficiaries []*TrustPerson    `protobuf:"bytes,11,rep,name=beneficiaries,proto3" json:"beneficiaries,omitempty"`
// Guarantees    []*TrustPerson    `protobuf:"bytes,12,rep,name=guarantees,proto3" json:"guarantees,omitempty"`
// Active        bool              `protobuf:"varint,13,opt,name=active,proto3" json:"active,omitempty"`
// Status        string            `protobuf:"bytes,14,opt,name=status,proto3" json:"status,omitempty"`
// RegisterDate  string            `protobuf:"bytes,15,opt,name=register_date,json=registerDate,proto3" json:"register_date,omitempty"`

type ChamaMember struct {
	ID            uint      `gorm:"primaryKey;autoIncrement"`
	ChamaID       string    `gorm:"type:varchar(15);not null"`
	FirstName     string    `gorm:"type:varchar(30);not null"`
	LastName      string    `gorm:"type:varchar(30);not null"`
	Phone         string    `gorm:"type:varchar(15);not null"`
	IDNumber      string    `gorm:"type:varchar(10);not null"`
	Email         string    `gorm:"type:varchar(50)"`
	Residence     string    `gorm:"type:varchar(100)"`
	Status        string    `gorm:"type:varchar(100)"`
	Beneficiaries []byte    `gorm:"type:json"`
	Guarantees    []byte    `gorm:"type:json"`
	JobDetails    []byte    `gorm:"type:json"`
	KYC           []byte    `gorm:"type:json"`
	Active        bool      `gorm:"type:tinyint(1)"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
}

func (*ChamaMember) TableName() string {
	return "chama_members"
}

func ChamaMemberModel(pb *chama.ChamaMember) (*ChamaMember, error) {
	if pb == nil {
		return nil, errs.NilObject("chama member")
	}

	db := &ChamaMember{
		ChamaID:   pb.ChamaId,
		FirstName: pb.FirstName,
		LastName:  pb.LastName,
		Phone:     pb.Phone,
		IDNumber:  pb.IdNumber,
		Email:     pb.Email,
		Residence: pb.Residence,
		Status:    pb.Status,
		Active:    pb.Active,
	}

	if len(pb.JobDetails) != 0 {
		bs, err := json.Marshal(pb.JobDetails)
		if err != nil {
			return nil, errs.FromJSONMarshal(err, "job details")
		}
		db.JobDetails = bs
	}

	if len(pb.Guarantees) != 0 {
		bs, err := json.Marshal(pb.Guarantees)
		if err != nil {
			return nil, errs.FromJSONMarshal(err, "guarantees")
		}
		db.Guarantees = bs
	}

	if len(pb.Beneficiaries) != 0 {
		bs, err := json.Marshal(pb.Beneficiaries)
		if err != nil {
			return nil, errs.FromJSONMarshal(err, "beneficiaries")
		}
		db.Beneficiaries = bs
	}

	if len(pb.Kyc) != 0 {
		bs, err := json.Marshal(pb.Kyc)
		if err != nil {
			return nil, errs.FromJSONMarshal(err, "kyc")
		}
		db.KYC = bs
	}

	return db, nil
}

func ChamaMemberProto(db *ChamaMember) (*chama.ChamaMember, error) {
	if db == nil {
		return nil, errs.NilObject("chama member")
	}

	pb := &chama.ChamaMember{
		MemberId:     fmt.Sprint(db.ID),
		ChamaId:      db.ChamaID,
		FirstName:    db.FirstName,
		LastName:     db.LastName,
		Email:        db.Email,
		Phone:        db.Phone,
		IdNumber:     db.IDNumber,
		Residence:    db.Residence,
		Status:       db.Status,
		Active:       db.Active,
		UpdatedDate:  db.UpdatedAt.String(),
		RegisterDate: db.CreatedAt.String(),
	}

	if len(db.Beneficiaries) != 0 {
		err := json.Unmarshal(db.Beneficiaries, &pb.Beneficiaries)
		if err != nil {
			return nil, errs.FromJSONUnMarshal(err, "beneficiaries")
		}
	}

	if len(db.Guarantees) != 0 {
		err := json.Unmarshal(db.Guarantees, &pb.Guarantees)
		if err != nil {
			return nil, errs.FromJSONUnMarshal(err, "guarantees")
		}
	}

	if len(db.JobDetails) != 0 {
		err := json.Unmarshal(db.JobDetails, &pb.JobDetails)
		if err != nil {
			return nil, errs.FromJSONUnMarshal(err, "job details")
		}
	}

	if len(db.KYC) != 0 {
		err := json.Unmarshal(db.KYC, &pb.Kyc)
		if err != nil {
			return nil, errs.FromJSONUnMarshal(err, "Kyc")
		}
	}

	return pb, nil
}
