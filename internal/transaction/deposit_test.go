package transaction

import (
	"context"

	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/machama-app/pkg/api/transaction"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Deposit", func() {
	var (
		depositReq *transaction.DepositRequest
		ctx        context.Context
	)

	BeforeEach(func() {
		depositReq = &transaction.DepositRequest{
			ActorId:     randomID(),
			AccountId:   randomID(),
			Description: randomDescription(),
			Amount:      randomdata.Decimal(1000, 10000),
		}
		ctx = context.TODO()
	})

	Describe("Deposit with malformed request", func() {
		It("should fail when the request is nil", func() {
			depositReq = nil
			depRes, err := TransactionAPI.Deposit(ctx, depositReq)
			Expect(err).Should(HaveOccurred())
			Expect(depRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when actor id is missing", func() {
			depositReq.ActorId = ""
			depRes, err := TransactionAPI.Deposit(ctx, depositReq)
			Expect(err).Should(HaveOccurred())
			Expect(depRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when account id is missing", func() {
			depositReq.AccountId = ""
			depRes, err := TransactionAPI.Deposit(ctx, depositReq)
			Expect(err).Should(HaveOccurred())
			Expect(depRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when description is missing", func() {
			depositReq.Description = ""
			depRes, err := TransactionAPI.Deposit(ctx, depositReq)
			Expect(err).Should(HaveOccurred())
			Expect(depRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when amount is missing", func() {
			depositReq.Amount = 0
			depRes, err := TransactionAPI.Deposit(ctx, depositReq)
			Expect(err).Should(HaveOccurred())
			Expect(depRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when amount is less than zero", func() {
			depositReq.Amount = -10
			depRes, err := TransactionAPI.Deposit(ctx, depositReq)
			Expect(err).Should(HaveOccurred())
			Expect(depRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
	})

	Describe("Deposit with well formed request", func() {
		It("should succees", func() {
			depRes, err := TransactionAPI.Deposit(ctx, depositReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(depRes).ShouldNot(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.OK))
		})
	})
})
