package moneyaccount

import (
	"context"

	"github.com/gidyon/machama-app/pkg/api/transaction"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("ListChamaAccounts", func() {
	var (
		listReq *transaction.ListChamaAccountsRequest
		ctx     context.Context
	)

	BeforeEach(func() {
		listReq = &transaction.ListChamaAccountsRequest{
			PageToken: "",
			PageSize:  defaultPageSize,
			Filter:    &transaction.ChamaAccountFilter{},
		}
		ctx = context.TODO()
	})

	Describe("ListChamaAccounts with malformed request", func() {
		It("should fail when the request is nil", func() {
			listReq = nil
			listRes, err := ChamaAccountAPI.ListChamaAccounts(ctx, listReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(listRes).Should(BeNil())
		})
		It("should fail when the page token is wrong", func() {
			listReq.PageToken = "weird"
			listRes, err := ChamaAccountAPI.ListChamaAccounts(ctx, listReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(listRes).Should(BeNil())
		})
	})

	Describe("ListChamaAccounts with well formed request", func() {
		Context("Lets create so many members first", func() {
			It("should succeed", func() {
				Expect(laodMockData(100)).ShouldNot(HaveOccurred())
			})
		})

		It("should succeed even when filter is nil", func() {
			listRes, err := ChamaAccountAPI.ListChamaAccounts(ctx, &transaction.ListChamaAccountsRequest{
				PageToken: "",
				PageSize:  defaultPageSize,
				Filter:    nil,
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(listRes.ChamaAccounts).ShouldNot(BeNil())
		})

		Describe("Looping through all members", func() {
			It("should succeed", func() {
				var pageToken, nextPageToken = "", "test"
				for nextPageToken != "" {
					listRes, err := ChamaAccountAPI.ListChamaAccounts(ctx, &transaction.ListChamaAccountsRequest{
						PageToken: pageToken,
						PageSize:  defaultPageSize,
					})
					Expect(err).ShouldNot(HaveOccurred())
					Expect(listRes.ChamaAccounts).ShouldNot(BeNil())

					nextPageToken = listRes.NextPageToken
					pageToken = nextPageToken
				}
			})
			It("should succeed when page sixe is 0", func() {
				var pageToken, nextPageToken = "", "test"
				for nextPageToken != "" {
					listRes, err := ChamaAccountAPI.ListChamaAccounts(ctx, &transaction.ListChamaAccountsRequest{
						PageToken: pageToken,
						PageSize:  0,
					})
					Expect(err).ShouldNot(HaveOccurred())
					Expect(listRes.ChamaAccounts).ShouldNot(BeNil())

					nextPageToken = listRes.NextPageToken
					pageToken = nextPageToken
				}
			})
			It("should succeed is larger than default", func() {
				var pageToken, nextPageToken = "", "test"
				for nextPageToken != "" {
					listRes, err := ChamaAccountAPI.ListChamaAccounts(ctx, &transaction.ListChamaAccountsRequest{
						PageToken: pageToken,
						PageSize:  1000,
					})
					Expect(err).ShouldNot(HaveOccurred())
					Expect(listRes.ChamaAccounts).ShouldNot(BeNil())

					nextPageToken = listRes.NextPageToken
					pageToken = nextPageToken
				}
			})
		})

		Describe("ListChamaAccounts with filters", func() {
			accountIDs := []string{"1"}

			It("should succeed when ACCOUNT_ID filter is on", func() {
				var pageToken, nextPageToken = "", "test"
				for nextPageToken != "" {
					listRes, err := ChamaAccountAPI.ListChamaAccounts(ctx, &transaction.ListChamaAccountsRequest{
						PageToken: pageToken,
						PageSize:  defaultPageSize,
						Filter: &transaction.ChamaAccountFilter{
							AccountIds: accountIDs,
						},
					})
					Expect(err).ShouldNot(HaveOccurred())
					Expect(listRes.ChamaAccounts).ShouldNot(BeNil())

					nextPageToken = listRes.NextPageToken
					pageToken = nextPageToken

					for _, chamaPB := range listRes.ChamaAccounts {
						Expect(chamaPB.AccountId).Should(BeElementOf(accountIDs))
					}
				}
			})

			ownerIds := []string{"3", "4"}

			It("should succeed when OWNER_ID filter is on", func() {
				var pageToken, nextPageToken = "", "test"
				for nextPageToken != "" {
					listRes, err := ChamaAccountAPI.ListChamaAccounts(ctx, &transaction.ListChamaAccountsRequest{
						PageToken: pageToken,
						PageSize:  defaultPageSize,
						Filter: &transaction.ChamaAccountFilter{
							OwnerIds: ownerIds,
						},
					})
					Expect(err).ShouldNot(HaveOccurred())
					Expect(listRes.ChamaAccounts).ShouldNot(BeNil())

					nextPageToken = listRes.NextPageToken
					pageToken = nextPageToken

					for _, chamaPB := range listRes.ChamaAccounts {
						Expect(chamaPB.OwnerId).Should(BeElementOf(ownerIds))
					}
				}
			})

			accountType := transaction.AccountType_SAVINGS_ACCOUNT

			It("should succeed when ACCOUNT_TYPE filter is on", func() {
				var pageToken, nextPageToken = "", "test"
				for nextPageToken != "" {
					listRes, err := ChamaAccountAPI.ListChamaAccounts(ctx, &transaction.ListChamaAccountsRequest{
						PageToken: pageToken,
						PageSize:  defaultPageSize,
						Filter: &transaction.ChamaAccountFilter{
							AccountType: accountType,
						},
					})
					Expect(err).ShouldNot(HaveOccurred())
					Expect(listRes.ChamaAccounts).ShouldNot(BeNil())

					nextPageToken = listRes.NextPageToken
					pageToken = nextPageToken

					for _, chamaPB := range listRes.ChamaAccounts {
						Expect(chamaPB.AccountType).Should(Equal(accountType))
					}
				}
			})

			It("should succeed when WITHDRAWABLE filter is on", func() {
				var pageToken, nextPageToken = "", "test"
				for nextPageToken != "" {
					listRes, err := ChamaAccountAPI.ListChamaAccounts(ctx, &transaction.ListChamaAccountsRequest{
						PageToken: pageToken,
						PageSize:  defaultPageSize,
						Filter: &transaction.ChamaAccountFilter{
							AccountType:  accountType,
							Withdrawable: true,
						},
					})
					Expect(err).ShouldNot(HaveOccurred())
					Expect(listRes.ChamaAccounts).ShouldNot(BeNil())

					nextPageToken = listRes.NextPageToken
					pageToken = nextPageToken

					for _, chamaPB := range listRes.ChamaAccounts {
						Expect(chamaPB.Withdrawable).Should(BeTrue())
					}
				}
			})
		})
	})
})
