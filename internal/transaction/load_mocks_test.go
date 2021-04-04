package transaction

import (
	"fmt"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/machama-app/internal/models"
	"github.com/gidyon/machama-app/pkg/api/transaction"
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

func randomTxType() transaction.TransactionType {
	if time.Now().Unix()%2 == 0 {
		return transaction.TransactionType_DEPOSIT
	}
	return transaction.TransactionType_WITHDRAWAL
}

func mockTransaction() *transaction.Transaction {
	return &transaction.Transaction{
		ActorId:           randomID(),
		AccountId:         randomID(),
		Description:       randomDescription(),
		TransactionType:   randomTxType(),
		TransactionAmount: randomdata.Decimal(1000, 10000),
	}
}

func laodMockData(count int) error {
	dbs := make([]*models.Transaction, 0, count)
	for i := 0; i < count; i++ {
		db, err := models.TransactionModel(mockTransaction())
		if err != nil {
			return err
		}
		dbs = append(dbs, db)
	}
	return TransactionAPIServer.SQLDB.CreateInBatches(dbs, count).Error
}
