package transaction

import (
	"context"

	"github.com/gidyon/machama-app/pkg/api/transaction"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("ListTransactions", func() {
	var (
		listReq *transaction.ListTransactionsRequest
		ctx     context.Context
	)

	BeforeEach(func() {
		listReq = &transaction.ListTransactionsRequest{
			PageToken: "",
			PageSize:  defaultPageSize,
			Filter:    &transaction.TransactionFilter{},
		}
		ctx = context.TODO()
	})

	Describe("ListTransactions with malformed request", func() {
		It("should fail when the request is nil", func() {
			listReq = nil
			listRes, err := TransactionAPI.ListTransactions(ctx, listReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(listRes).Should(BeNil())
		})
		It("should fail when the page token is wrong", func() {
			listReq.PageToken = "weird"
			listRes, err := TransactionAPI.ListTransactions(ctx, listReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(listRes).Should(BeNil())
		})
	})

	Describe("ListTransactions with well formed request", func() {
		Context("Lets create so many members first", func() {
			It("should succeed", func() {
				Expect(laodMockData(100)).ShouldNot(HaveOccurred())
			})
		})

		It("should succeed even when filter is nil", func() {
			listRes, err := TransactionAPI.ListTransactions(ctx, &transaction.ListTransactionsRequest{
				PageToken: "",
				PageSize:  defaultPageSize,
				Filter:    nil,
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(listRes.Transactions).ShouldNot(BeNil())
		})

		Describe("Looping through all members", func() {
			It("should succeed", func() {
				var pageToken, nextPageToken = "", "test"
				for nextPageToken != "" {
					listRes, err := TransactionAPI.ListTransactions(ctx, &transaction.ListTransactionsRequest{
						PageToken: pageToken,
						PageSize:  defaultPageSize,
					})
					Expect(err).ShouldNot(HaveOccurred())
					Expect(listRes.Transactions).ShouldNot(BeNil())

					nextPageToken = listRes.NextPageToken
					pageToken = nextPageToken
				}
			})
			It("should succeed when page sixe is 0", func() {
				var pageToken, nextPageToken = "", "test"
				for nextPageToken != "" {
					listRes, err := TransactionAPI.ListTransactions(ctx, &transaction.ListTransactionsRequest{
						PageToken: pageToken,
						PageSize:  0,
					})
					Expect(err).ShouldNot(HaveOccurred())
					Expect(listRes.Transactions).ShouldNot(BeNil())

					nextPageToken = listRes.NextPageToken
					pageToken = nextPageToken
				}
			})
			It("should succeed is larger than default", func() {
				var pageToken, nextPageToken = "", "test"
				for nextPageToken != "" {
					listRes, err := TransactionAPI.ListTransactions(ctx, &transaction.ListTransactionsRequest{
						PageToken: pageToken,
						PageSize:  1000,
					})
					Expect(err).ShouldNot(HaveOccurred())
					Expect(listRes.Transactions).ShouldNot(BeNil())

					nextPageToken = listRes.NextPageToken
					pageToken = nextPageToken
				}
			})
		})

		Describe("ListTransactions with filters", func() {
			accountIDs := []string{"1"}

			It("should succeed when ACCOUNT_ID filter is on", func() {
				var pageToken, nextPageToken = "", "test"
				for nextPageToken != "" {
					listRes, err := TransactionAPI.ListTransactions(ctx, &transaction.ListTransactionsRequest{
						PageToken: pageToken,
						PageSize:  defaultPageSize,
						Filter: &transaction.TransactionFilter{
							AccountIds: accountIDs,
						},
					})
					Expect(err).ShouldNot(HaveOccurred())
					Expect(listRes.Transactions).ShouldNot(BeNil())

					nextPageToken = listRes.NextPageToken
					pageToken = nextPageToken

					for _, chamaPB := range listRes.Transactions {
						Expect(chamaPB.AccountId).Should(BeElementOf(accountIDs))
					}
				}
			})

			actorIDs := []string{"3", "4"}

			It("should succeed when ACTOR_ID filter is on", func() {
				var pageToken, nextPageToken = "", "test"
				for nextPageToken != "" {
					listRes, err := TransactionAPI.ListTransactions(ctx, &transaction.ListTransactionsRequest{
						PageToken: pageToken,
						PageSize:  defaultPageSize,
						Filter: &transaction.TransactionFilter{
							ActorIds: actorIDs,
						},
					})
					Expect(err).ShouldNot(HaveOccurred())
					Expect(listRes.Transactions).ShouldNot(BeNil())

					nextPageToken = listRes.NextPageToken
					pageToken = nextPageToken

					for _, chamaPB := range listRes.Transactions {
						Expect(chamaPB.ActorId).Should(BeElementOf(actorIDs))
					}
				}
			})

			txIDs := []string{"3", "4", "5", "6", "7"}

			It("should succeed when TRANSACTION_ID filter is on", func() {
				var pageToken, nextPageToken = "", "test"
				for nextPageToken != "" {
					listRes, err := TransactionAPI.ListTransactions(ctx, &transaction.ListTransactionsRequest{
						PageToken: pageToken,
						PageSize:  defaultPageSize,
						Filter: &transaction.TransactionFilter{
							TransactionIds: txIDs,
						},
					})
					Expect(err).ShouldNot(HaveOccurred())
					Expect(listRes.Transactions).ShouldNot(BeNil())

					nextPageToken = listRes.NextPageToken
					pageToken = nextPageToken

					for _, chamaPB := range listRes.Transactions {
						Expect(chamaPB.TransactionId).Should(BeElementOf(txIDs))
					}
				}
			})

			txType := transaction.TransactionType_DEPOSIT

			It("should succeed when TRANSACTION_TYPE filter is on", func() {
				var pageToken, nextPageToken = "", "test"
				for nextPageToken != "" {
					listRes, err := TransactionAPI.ListTransactions(ctx, &transaction.ListTransactionsRequest{
						PageToken: pageToken,
						PageSize:  defaultPageSize,
						Filter: &transaction.TransactionFilter{
							TransactionType: txType,
						},
					})
					Expect(err).ShouldNot(HaveOccurred())
					Expect(listRes.Transactions).ShouldNot(BeNil())

					nextPageToken = listRes.NextPageToken
					pageToken = nextPageToken

					for _, chamaPB := range listRes.Transactions {
						Expect(chamaPB.TransactionType).Should(Equal(txType))
					}
				}
			})
		})
	})
})
