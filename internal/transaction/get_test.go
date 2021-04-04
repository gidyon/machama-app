package transaction

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

var _ = Describe("GetTransaction", func() {
	var (
		getReq *transaction.GetTransactionRequest
		ctx    context.Context
	)

	BeforeEach(func() {
		getReq = &transaction.GetTransactionRequest{
			TransactionId: "1",
		}
		ctx = context.TODO()
	})

	Describe("GetTransaction with malformed request", func() {
		It("should fail when the request is nil", func() {
			getReq = nil
			getRes, err := TransactionAPI.GetTransaction(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(getRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when transaction id is missing", func() {
			getReq.TransactionId = ""
			getRes, err := TransactionAPI.GetTransaction(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(getRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when transaction does not exist", func() {
			getReq.TransactionId = "oops"
			getRes, err := TransactionAPI.GetTransaction(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(getRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.NotFound))
		})
	})

	Describe("GetTransaction with well formed request", func() {
		var transactionID string

		Context("Lets create transaction first", func() {
			It("should succeed", func() {
				chamaMemberDB, err := models.TransactionModel(mockTransaction())
				Expect(err).ShouldNot(HaveOccurred())

				Expect(TransactionAPIServer.SQLDB.Create(chamaMemberDB).Error).ShouldNot(HaveOccurred())

				transactionID = fmt.Sprint(chamaMemberDB.ID)
			})
		})

		Describe("Getting the transaction", func() {
			It("should succeed", func() {
				getRes, err := TransactionAPI.GetTransaction(ctx, &transaction.GetTransactionRequest{
					TransactionId: transactionID,
				})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(getRes).ShouldNot(BeNil())
				Expect(status.Code(err)).Should(Equal(codes.OK))
			})
		})
	})
})
