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

var _ = Describe("Withdraw", func() {
	var (
		withdrawReq *transaction.WithdrawRequest
		ctx         context.Context
	)

	BeforeEach(func() {
		withdrawReq = &transaction.WithdrawRequest{
			ActorId:     randomID(),
			AccountId:   randomID(),
			Description: randomDescription(),
			Amount:      randomdata.Decimal(1000, 10000),
		}
		ctx = context.TODO()
	})

	Describe("Withdraw with malformed request", func() {
		It("should fail when the request is nil", func() {
			withdrawReq = nil
			withdrawRes, err := TransactionAPI.Withdraw(ctx, withdrawReq)
			Expect(err).Should(HaveOccurred())
			Expect(withdrawRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when actor id is missing", func() {
			withdrawReq.ActorId = ""
			withdrawRes, err := TransactionAPI.Withdraw(ctx, withdrawReq)
			Expect(err).Should(HaveOccurred())
			Expect(withdrawRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when account id is missing", func() {
			withdrawReq.AccountId = ""
			withdrawRes, err := TransactionAPI.Withdraw(ctx, withdrawReq)
			Expect(err).Should(HaveOccurred())
			Expect(withdrawRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when description is missing", func() {
			withdrawReq.Description = ""
			withdrawRes, err := TransactionAPI.Withdraw(ctx, withdrawReq)
			Expect(err).Should(HaveOccurred())
			Expect(withdrawRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when amount is missing", func() {
			withdrawReq.Amount = 0
			withdrawRes, err := TransactionAPI.Withdraw(ctx, withdrawReq)
			Expect(err).Should(HaveOccurred())
			Expect(withdrawRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when amount is less than zero", func() {
			withdrawReq.Amount = -10
			withdrawRes, err := TransactionAPI.Withdraw(ctx, withdrawReq)
			Expect(err).Should(HaveOccurred())
			Expect(withdrawRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
	})

	Describe("Withdraw with well formed request", func() {
		It("should succees", func() {
			withdrawRes, err := TransactionAPI.Withdraw(ctx, withdrawReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(withdrawRes).ShouldNot(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.OK))
		})
	})
})
