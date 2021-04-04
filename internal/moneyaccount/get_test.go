package moneyaccount

import (
	"context"
	"fmt"

	"github.com/gidyon/machama-app/internal/models"
	"github.com/gidyon/machama-app/pkg/api/transaction"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("GetChamaAccount", func() {
	var (
		getReq *transaction.GetChamaAccountRequest
		ctx    context.Context
	)

	BeforeEach(func() {
		getReq = &transaction.GetChamaAccountRequest{
			AccountId: "1",
		}
		ctx = context.TODO()
	})

	Describe("GetChamaAccount with malformed request", func() {
		It("should fail when the request is nil", func() {
			getReq = nil
			getRes, err := ChamaAccountAPI.GetChamaAccount(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(getRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when chama account id is missing", func() {
			getReq.AccountId = ""
			getRes, err := ChamaAccountAPI.GetChamaAccount(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(getRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when chama account does not exist", func() {
			getReq.AccountId = "oops"
			getRes, err := ChamaAccountAPI.GetChamaAccount(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(getRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.NotFound))
		})
	})

	Describe("GetChamaAccount with well formed request", func() {
		var chamaMemberID, ownerID, accountName string

		Context("Lets create chama account first", func() {
			It("should succeed", func() {
				chamaMemberDB, err := models.ChamaAccountModel(mockChamaAccount())
				Expect(err).ShouldNot(HaveOccurred())

				Expect(ChamaAccountAPIServer.SQLDB.Create(chamaMemberDB).Error).ShouldNot(HaveOccurred())

				chamaMemberID = fmt.Sprint(chamaMemberDB.ID)
				ownerID = chamaMemberDB.OwnerID
				accountName = chamaMemberDB.AccountName
			})
		})

		Describe("Getting the chama account", func() {
			It("should succeed when account id is present", func() {
				getRes, err := ChamaAccountAPI.GetChamaAccount(ctx, &transaction.GetChamaAccountRequest{
					AccountId: chamaMemberID,
				})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(getRes).ShouldNot(BeNil())
				Expect(status.Code(err)).Should(Equal(codes.OK))
			})

			It("should succeed when owner id and account_name is provided", func() {
				getRes, err := ChamaAccountAPI.GetChamaAccount(ctx, &transaction.GetChamaAccountRequest{
					OwnerId:     ownerID,
					AccountName: accountName,
				})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(getRes).ShouldNot(BeNil())
				Expect(status.Code(err)).Should(Equal(codes.OK))
			})
		})
	})
})
