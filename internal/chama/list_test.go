package chama

import (
	"context"

	"github.com/gidyon/machama-app/pkg/api/chama"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("ListChamas", func() {
	var (
		listReq *chama.ListChamasRequest
		ctx     context.Context
	)

	BeforeEach(func() {
		listReq = &chama.ListChamasRequest{
			PageToken: "",
			PageSize:  defaultPageSize,
			Filter:    &chama.ChamaFilter{},
		}
		ctx = context.TODO()
	})

	Describe("ListChams with malformed request", func() {
		It("should fail when the request is nil", func() {
			listReq = nil
			listRes, err := ChamaAPI.ListChamas(ctx, listReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(listRes).Should(BeNil())
		})
		It("should fail when the page token is wrong", func() {
			listReq.PageToken = "weird"
			listRes, err := ChamaAPI.ListChamas(ctx, listReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(listRes).Should(BeNil())
		})
	})

	Describe("ListChamas with well formed request", func() {
		Context("Lets create so many chamas first", func() {
			It("should succeed", func() {
				Expect(laodMockData(100)).ShouldNot(HaveOccurred())
			})
		})

		It("should succeed even when filter is nil", func() {
			listRes, err := ChamaAPI.ListChamas(ctx, &chama.ListChamasRequest{
				PageToken: "",
				PageSize:  defaultPageSize,
				Filter:    nil,
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(listRes.Chamas).ShouldNot(BeNil())
		})

		Describe("Looping through all chamas", func() {
			It("should succeed", func() {
				var pageToken, nextPageToken = "", "test"
				for nextPageToken != "" {
					listRes, err := ChamaAPI.ListChamas(ctx, &chama.ListChamasRequest{
						PageToken: pageToken,
						PageSize:  defaultPageSize,
					})
					Expect(err).ShouldNot(HaveOccurred())
					Expect(listRes.Chamas).ShouldNot(BeNil())

					nextPageToken = listRes.NextPageToken
					pageToken = nextPageToken
				}
			})
			It("should succeed when page sixe is 0", func() {
				var pageToken, nextPageToken = "", "test"
				for nextPageToken != "" {
					listRes, err := ChamaAPI.ListChamas(ctx, &chama.ListChamasRequest{
						PageToken: pageToken,
						PageSize:  0,
					})
					Expect(err).ShouldNot(HaveOccurred())
					Expect(listRes.Chamas).ShouldNot(BeNil())

					nextPageToken = listRes.NextPageToken
					pageToken = nextPageToken
				}
			})
			It("should succeed is larger than default", func() {
				var pageToken, nextPageToken = "", "test"
				for nextPageToken != "" {
					listRes, err := ChamaAPI.ListChamas(ctx, &chama.ListChamasRequest{
						PageToken: pageToken,
						PageSize:  1000,
					})
					Expect(err).ShouldNot(HaveOccurred())
					Expect(listRes.Chamas).ShouldNot(BeNil())

					nextPageToken = listRes.NextPageToken
					pageToken = nextPageToken
				}
			})
		})

		Describe("ListChamas with filters", func() {
			creatorIDs := []string{"1"}

			It("should succeed when CREATOR_ID filter is on", func() {
				var pageToken, nextPageToken = "", "test"
				for nextPageToken != "" {
					listRes, err := ChamaAPI.ListChamas(ctx, &chama.ListChamasRequest{
						PageToken: pageToken,
						PageSize:  defaultPageSize,
						Filter: &chama.ChamaFilter{
							CreatorIds: creatorIDs,
						},
					})
					Expect(err).ShouldNot(HaveOccurred())
					Expect(listRes.Chamas).ShouldNot(BeNil())

					nextPageToken = listRes.NextPageToken
					pageToken = nextPageToken

					for _, chamaPB := range listRes.Chamas {
						Expect(chamaPB.CreatorId).Should(BeElementOf(creatorIDs))
					}
				}
			})
		})
	})
})
