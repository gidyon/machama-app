package moneyaccount

import (
	"fmt"

	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/machama-app/internal/models"
	"github.com/gidyon/machama-app/pkg/api/transaction"
)

func randomID() string {
	return fmt.Sprint(randomdata.Number(1, 10))
}

func mockChamaAccount() *transaction.ChamaAccount {
	return &transaction.ChamaAccount{
		OwnerId:              randomID(),
		AccountName:          randomdata.SillyName(),
		AccountType:          transaction.AccountType_SAVINGS_ACCOUNT,
		Withdrawable:         true,
		AvailableAmount:      randomdata.Decimal(10000, 100000),
		TotalDepositedAmount: randomdata.Decimal(10000, 1000000),
		TotalWithdrawnAmount: randomdata.Decimal(10000, 1000000),
		LastDepositedAmount:  randomdata.Decimal(1000, 10000),
		LastWithdrawnAmount:  randomdata.Decimal(1000, 10000),
		Active:               false,
	}
}

func laodMockData(count int) error {
	dbs := make([]*models.ChamaAccount, 0, count)
	for i := 0; i < count; i++ {
		db, err := models.ChamaAccountModel(mockChamaAccount())
		if err != nil {
			return err
		}
		dbs = append(dbs, db)
	}
	return ChamaAccountAPIServer.SQLDB.CreateInBatches(dbs, count).Error
}
