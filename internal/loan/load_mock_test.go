package loan

import (
	"fmt"

	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/machama-app/internal/models"
	"github.com/gidyon/machama-app/pkg/api/loan"
)

func randomID() string {
	return fmt.Sprint(randomdata.Number(1, 10))
}

func randomPhone() string {
	return fmt.Sprintf("2547%d", randomdata.Number(10000000, 99999999))
}

func mockLoan() *loan.Loan {
	return &loan.Loan{
		ChamaId:      randomID(),
		ProductId:    randomID(),
		MemberId:     randomID(),
		LoaneeNames:  randomdata.SillyName(),
		LoaneePhone:  randomPhone(),
		LoaneeEmail:  randomdata.Email(),
		NationalId:   fmt.Sprint(randomdata.Number(22222222, 44444444)),
		DurationDays: int32(randomdata.Number(10, 100)),
		InterestRate: float32(randomdata.Decimal(3, 30)),
	}
}

func laodMockData(count int) error {
	dbs := make([]*models.Loan, 0, count)
	for i := 0; i < count; i++ {
		db, err := models.LoanModel(mockLoan())
		if err != nil {
			return err
		}
		dbs = append(dbs, db)
	}
	return LoanAPIServer.SQLDB.CreateInBatches(dbs, count).Error
}
