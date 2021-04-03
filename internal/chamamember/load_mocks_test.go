package chamamember

import (
	"fmt"

	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/machama-app/internal/models"
	"github.com/gidyon/machama-app/pkg/api/chama"
)

func randomChamaID() string {
	return fmt.Sprint(randomdata.Number(1, 10))
}

func randomPhone() string {
	return fmt.Sprintf("2547%d", randomdata.Number(10000000, 99999999))
}

func randomStatus() string {
	st := randomdata.Paragraph()
	if len(st) > 50 {
		return st[:50]
	}
	return st
}

func randomTrustPersion() *chama.TrustPerson {
	return &chama.TrustPerson{
		Name:  randomdata.SillyName(),
		Email: randomdata.Email(),
		Phone: randomPhone(),
	}
}

func mockChamaMember() *chama.ChamaMember {
	return &chama.ChamaMember{
		ChamaId:   randomChamaID(),
		FirstName: randomdata.FirstName(1),
		LastName:  randomdata.LastName(),
		Phone:     randomPhone(),
		IdNumber:  fmt.Sprint(randomdata.Number(22222222, 44444444)),
		Email:     randomdata.Email(),
		Residence: randomdata.Address(),
		Status:    randomStatus(),
		Beneficiaries: []*chama.TrustPerson{
			randomTrustPersion(),
			randomTrustPersion(),
			randomTrustPersion(),
		},
		Guarantees: []*chama.TrustPerson{
			randomTrustPersion(),
			randomTrustPersion(),
			randomTrustPersion(),
		},
		Kyc: map[string]string{
			"id_card_url": randomdata.IpV4Address(),
		},
		Active: false,
	}
}

func laodMockData(count int) error {
	dbs := make([]*models.ChamaMember, 0, count)
	for i := 0; i < count; i++ {
		db, err := models.ChamaMemberModel(mockChamaMember())
		if err != nil {
			return err
		}
		dbs = append(dbs, db)
	}
	return ChamaMemberAPIServer.SQLDB.CreateInBatches(dbs, count).Error
}
