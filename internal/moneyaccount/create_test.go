package moneyaccount

import (
	"context"

	"github.com/gidyon/machama-app/pkg/api/transaction"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("CreateChamaAccount", func() {
	var (
		createReq *transaction.CreateChamaAccountRequest
		ctx       context.Context
	)

	BeforeEach(func() {
		createReq = &transaction.CreateChamaAccountRequest{
			ChamaAccount: mockChamaAccount(),
		}
		ctx = context.TODO()
	})

	Describe("CreateChamaAccount with malformed request", func() {
		It("should fail when the request is nil", func() {
			createReq = nil
			createRes, err := ChamaAccountAPI.CreateChamaAccount(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when account is nil", func() {
			createReq.ChamaAccount = nil
			createRes, err := ChamaAccountAPI.CreateChamaAccount(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when owner id is missing", func() {
			createReq.ChamaAccount.OwnerId = ""
			createRes, err := ChamaAccountAPI.CreateChamaAccount(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when account name is missing", func() {
			createReq.ChamaAccount.AccountName = ""
			createRes, err := ChamaAccountAPI.CreateChamaAccount(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when account type is missing", func() {
			createReq.ChamaAccount.AccountType = transaction.AccountType_ACCOUNT_TYPE_UNSPECIFIED
			createRes, err := ChamaAccountAPI.CreateChamaAccount(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
	})

	Describe("CreateChamaAccount with wellformed request", func() {
		It("should succeed", func() {
			createRes, err := ChamaAccountAPI.CreateChamaAccount(ctx, createReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(createRes).ShouldNot(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.OK))
		})
	})
})
