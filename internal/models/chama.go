package models

import (
	"fmt"
	"time"

	"github.com/gidyon/machama-app/pkg/api/chama"
	"github.com/gidyon/micro/v2/utils/errs"
)

// ChamaId        string  `protobuf:"bytes,1,opt,name=chama_id,json=chamaId,proto3" json:"chama_id,omitempty"`
// Name           string  `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
// Description    string  `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
// Status         string  `protobuf:"bytes,4,opt,name=status,proto3" json:"status,omitempty"`
// AccountBalance float32 `protobuf:"fixed32,5,opt,name=account_balance,json=accountBalance,proto3" json:"account_balance,omitempty"`
// Active         bool    `protobuf:"varint,6,opt,name=active,proto3" json:"active,omitempty"`
// CreateDate     string  `protobuf:"bytes,7,opt,name=create_date,json=createDate,proto3" json:"create_date,omitempty"`

type Chama struct {
	ID             uint      `gorm:"primaryKey;autoIncrement"`
	CreatorID      string    `gorm:"type:varchar(50);not null"`
	Name           string    `gorm:"type:varchar(50);not null"`
	Description    string    `gorm:"type:varchar(200)"`
	Status         string    `gorm:"type:varchar(100)"`
	AccountBalance float64   `gorm:"type:float(15)"`
	Active         bool      `gorm:"type:tinyint(1)"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
}

func (*Chama) TableName() string {
	return "chamas"
}

func ChamaModel(pb *chama.Chama) (*Chama, error) {
	if pb == nil {
		return nil, errs.NilObject("chama")
	}
	return &Chama{
		Name:           pb.Name,
		CreatorID:      pb.CreatorId,
		Description:    pb.Description,
		Status:         pb.Status,
		AccountBalance: pb.AccountBalance,
		Active:         pb.Active,
	}, nil
}

func ChamaProto(db *Chama) (*chama.Chama, error) {
	if db == nil {
		return nil, errs.NilObject("chama")
	}
	return &chama.Chama{
		ChamaId:        fmt.Sprint(db.ID),
		CreatorId:      db.CreatorID,
		Name:           db.Name,
		Description:    db.Description,
		Status:         db.Status,
		AccountBalance: db.AccountBalance,
		Active:         db.Active,
		UpdatedDate:    db.UpdatedAt.String(),
		CreatedDate:    db.CreatedAt.String(),
	}, nil
}
