package loanproduct

import (
	"fmt"

	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/machama-app/internal/models"
	"github.com/gidyon/machama-app/pkg/api/loan"
)

func randomID() string {
	return fmt.Sprint(randomdata.Number(1, 10))
}

func randomDescription() string {
	desc := randomdata.Paragraph()
	if len(desc) > 100 {
		return desc[:100]
	}
	return desc
}

func mockLoanProduct() *loan.LoanProduct {
	return &loan.LoanProduct{
		ChamaId:           randomID(),
		ProductId:         randomID(),
		Name:              randomdata.SillyName(),
		Description:       randomDescription(),
		LoanDurationDays:  int32(randomdata.Number(10, 100)),
		LoanMinimumAmount: 100.00,
		LoanMaximumAmount: 1000000.00,
		InterestRate:      float32(randomdata.Decimal(3, 30)),
	}
}

func laodMockData(count int) error {
	dbs := make([]*models.LoanProduct, 0, count)
	for i := 0; i < count; i++ {
		db, err := models.LoanProductModel(mockLoanProduct())
		if err != nil {
			return err
		}
		dbs = append(dbs, db)
	}
	return LoanProductAPIServer.SQLDB.CreateInBatches(dbs, count).Error
}
