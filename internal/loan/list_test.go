package loan

import (
	"context"

	"github.com/gidyon/machama-app/pkg/api/loan"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("ListLoans", func() {
	var (
		listReq *loan.ListLoansRequest
		ctx     context.Context
	)

	BeforeEach(func() {
		listReq = &loan.ListLoansRequest{
			PageToken: "",
			PageSize:  defaultPageSize,
			Filter:    &loan.LoanFilter{},
		}
		ctx = context.TODO()
	})

	Describe("ListLoans with malformed request", func() {
		It("should fail when the request is nil", func() {
			listReq = nil
			listRes, err := LoanAPI.ListLoans(ctx, listReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(listRes).Should(BeNil())
		})
		It("should fail when the page token is wrong", func() {
			listReq.PageToken = "weird"
			listRes, err := LoanAPI.ListLoans(ctx, listReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(listRes).Should(BeNil())
		})
	})

	Describe("ListLoans with well formed request", func() {
		Context("Lets create so many loans first", func() {
			It("should succeed", func() {
				Expect(laodMockData(100)).ShouldNot(HaveOccurred())
			})
		})

		It("should succeed even when filter is nil", func() {
			listRes, err := LoanAPI.ListLoans(ctx, &loan.ListLoansRequest{
				PageToken: "",
				PageSize:  defaultPageSize,
				Filter:    nil,
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(listRes.Loans).ShouldNot(BeNil())
		})

		Describe("Looping through all loans", func() {
			It("should succeed", func() {
				var pageToken, nextPageToken = "", "test"
				for nextPageToken != "" {
					listRes, err := LoanAPI.ListLoans(ctx, &loan.ListLoansRequest{
						PageToken: pageToken,
						PageSize:  defaultPageSize,
					})
					Expect(err).ShouldNot(HaveOccurred())
					Expect(listRes.Loans).ShouldNot(BeNil())

					nextPageToken = listRes.NextPageToken
					pageToken = nextPageToken
				}
			})
			It("should succeed when page sixe is 0", func() {
				var pageToken, nextPageToken = "", "test"
				for nextPageToken != "" {
					listRes, err := LoanAPI.ListLoans(ctx, &loan.ListLoansRequest{
						PageToken: pageToken,
						PageSize:  0,
					})
					Expect(err).ShouldNot(HaveOccurred())
					Expect(listRes.Loans).ShouldNot(BeNil())

					nextPageToken = listRes.NextPageToken
					pageToken = nextPageToken
				}
			})
			It("should succeed is larger than default", func() {
				var pageToken, nextPageToken = "", "test"
				for nextPageToken != "" {
					listRes, err := LoanAPI.ListLoans(ctx, &loan.ListLoansRequest{
						PageToken: pageToken,
						PageSize:  1000,
					})
					Expect(err).ShouldNot(HaveOccurred())
					Expect(listRes.Loans).ShouldNot(BeNil())

					nextPageToken = listRes.NextPageToken
					pageToken = nextPageToken
				}
			})
		})

		Describe("ListLoans with filters", func() {
			chamaIDs := []string{"1"}
			productIds := []string{"1", "2"}

			It("should succeed when CHAMA_ID filter is on", func() {
				var pageToken, nextPageToken = "", "test"
				for nextPageToken != "" {
					listRes, err := LoanAPI.ListLoans(ctx, &loan.ListLoansRequest{
						PageToken: pageToken,
						PageSize:  defaultPageSize,
						Filter: &loan.LoanFilter{
							ChamaIds: chamaIDs,
						},
					})
					Expect(err).ShouldNot(HaveOccurred())
					Expect(listRes.Loans).ShouldNot(BeNil())

					nextPageToken = listRes.NextPageToken
					pageToken = nextPageToken

					for _, chamaPB := range listRes.Loans {
						Expect(chamaPB.ChamaId).Should(BeElementOf(chamaIDs))
					}
				}
			})

			It("should succeed when PRODUCT_ID filter is on", func() {
				var pageToken, nextPageToken = "", "test"
				for nextPageToken != "" {
					listRes, err := LoanAPI.ListLoans(ctx, &loan.ListLoansRequest{
						PageToken: pageToken,
						PageSize:  defaultPageSize,
						Filter: &loan.LoanFilter{
							ProductIds: productIds,
						},
					})
					Expect(err).ShouldNot(HaveOccurred())
					Expect(listRes.Loans).ShouldNot(BeNil())

					nextPageToken = listRes.NextPageToken
					pageToken = nextPageToken

					for _, chamaPB := range listRes.Loans {
						Expect(chamaPB.ProductId).Should(BeElementOf(productIds))
					}
				}
			})
		})
	})
})
