package moneyaccount

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/machama-app/internal/models"
	"github.com/gidyon/machama-app/pkg/api/transaction"
	"github.com/gidyon/micro/v2"
	"github.com/gidyon/micro/v2/pkg/conn"
	"github.com/gidyon/micro/v2/pkg/mocks"
	"github.com/gidyon/micro/v2/utils/encryption"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

func TestChama(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Money Account Suite")
}

var (
	ChamaAccountAPIServer *moneyAccountAPIServer
	ChamaAccountAPI       transaction.ChamaAccountAPIServer
	modelsStructs         = []interface{}{
		&models.ChamaAccount{},
	}
	schema = "machama"
)

func startDB() (*gorm.DB, error) {
	return conn.OpenGormConn(&conn.DBOptions{
		Dialect:  "mysql",
		Address:  "localhost:3306",
		User:     "root",
		Password: "hakty11",
		Schema:   schema,
	})
}

var _ = BeforeSuite(func() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rand.Seed(time.Now().UnixNano())

	// Start real database
	db, err := startDB()
	Expect(err).ShouldNot(HaveOccurred())

	db = db.Debug()

	// db = db.Debug()
	err = db.Migrator().DropTable(modelsStructs...)
	Expect(err).ShouldNot(HaveOccurred())

	err = db.Migrator().AutoMigrate(modelsStructs...)
	Expect(err).ShouldNot(HaveOccurred())

	hasher, err := encryption.NewHasher(string([]byte(randomdata.RandStringRunes(32))))
	Expect(err).ShouldNot(HaveOccurred())

	logger := micro.NewLogger("ussdlog", 0)

	authAPI := mocks.AuthAPI

	opt := &Options{
		SQLDB:      db,
		Logger:     logger,
		PageHasher: hasher,
		Auth:       authAPI,
	}

	// Create ussdlog API
	ChamaAccountAPI, err = NewChamaAccountAPI(ctx, opt)
	Expect(err).ShouldNot(HaveOccurred())

	var ok bool
	ChamaAccountAPIServer, ok = ChamaAccountAPI.(*moneyAccountAPIServer)
	Expect(ok).Should(BeTrue())

	// Load data
	Expect(laodMockData(100)).ShouldNot(HaveOccurred())

	_, err = NewChamaAccountAPI(ctx, nil)
	Expect(err).Should(HaveOccurred())

	opt.SQLDB = nil
	_, err = NewChamaAccountAPI(ctx, opt)
	Expect(err).Should(HaveOccurred())

	opt.SQLDB = db
	opt.Logger = nil
	_, err = NewChamaAccountAPI(ctx, opt)
	Expect(err).Should(HaveOccurred())

	opt.Logger = logger
	opt.PageHasher = nil
	_, err = NewChamaAccountAPI(ctx, opt)
	Expect(err).Should(HaveOccurred())

	opt.PageHasher = hasher
	opt.Auth = nil
	_, err = NewChamaAccountAPI(ctx, opt)
	Expect(err).Should(HaveOccurred())

	opt.Auth = authAPI
	_, err = NewChamaAccountAPI(ctx, opt)
	Expect(err).ShouldNot(HaveOccurred())
})
