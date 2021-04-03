package chama

import (
	"fmt"

	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/machama-app/internal/models"
)

func randomCreatorID() string {
	return fmt.Sprint(randomdata.Number(1, 10))
}

func randomDescription() string {
	desc := randomdata.Paragraph()
	if len(desc) > 100 {
		return desc[:100]
	}
	return desc
}

func randomStatus() string {
	st := randomdata.Paragraph()
	if len(st) > 50 {
		return st[:50]
	}
	return st
}

func mockChama() *models.Chama {
	return &models.Chama{
		CreatorID:      randomCreatorID(),
		Name:           randomdata.SillyName() + " SACCO",
		Description:    randomDescription(),
		Status:         randomStatus(),
		AccountBalance: 0,
		Active:         false,
	}
}

func laodMockData(count int) error {
	dbs := make([]*models.Chama, 0, count)
	for i := 0; i < count; i++ {
		dbs = append(dbs, mockChama())
	}
	return ChamaAPIServer.SQLDB.CreateInBatches(dbs, count).Error
}
