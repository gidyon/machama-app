package chamamember

import (
	"context"

	"github.com/gidyon/machama-app/pkg/api/chama"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("ListChamaMembers", func() {
	var (
		listReq *chama.ListChamaMembersRequest
		ctx     context.Context
	)

	BeforeEach(func() {
		listReq = &chama.ListChamaMembersRequest{
			PageToken: "",
			PageSize:  defaultPageSize,
			Filter:    &chama.ChamaMemberFilter{},
		}
		ctx = context.TODO()
	})

	Describe("ListChams with malformed request", func() {
		It("should fail when the request is nil", func() {
			listReq = nil
			listRes, err := ChamaMemberAPI.ListChamaMembers(ctx, listReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(listRes).Should(BeNil())
		})
		It("should fail when the page token is wrong", func() {
			listReq.PageToken = "weird"
			listRes, err := ChamaMemberAPI.ListChamaMembers(ctx, listReq)
			Expect(err).Should(HaveOccurred())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
			Expect(listRes).Should(BeNil())
		})
	})

	Describe("ListChamaMembers with well formed request", func() {
		Context("Lets create so many members first", func() {
			It("should succeed", func() {
				Expect(laodMockData(100)).ShouldNot(HaveOccurred())
			})
		})

		It("should succeed even when filter is nil", func() {
			listRes, err := ChamaMemberAPI.ListChamaMembers(ctx, &chama.ListChamaMembersRequest{
				PageToken: "",
				PageSize:  defaultPageSize,
				Filter:    nil,
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(listRes.ChamaMembers).ShouldNot(BeNil())
		})

		Describe("Looping through all members", func() {
			It("should succeed", func() {
				var pageToken, nextPageToken = "", "test"
				for nextPageToken != "" {
					listRes, err := ChamaMemberAPI.ListChamaMembers(ctx, &chama.ListChamaMembersRequest{
						PageToken: pageToken,
						PageSize:  defaultPageSize,
					})
					Expect(err).ShouldNot(HaveOccurred())
					Expect(listRes.ChamaMembers).ShouldNot(BeNil())

					nextPageToken = listRes.NextPageToken
					pageToken = nextPageToken
				}
			})
			It("should succeed when page sixe is 0", func() {
				var pageToken, nextPageToken = "", "test"
				for nextPageToken != "" {
					listRes, err := ChamaMemberAPI.ListChamaMembers(ctx, &chama.ListChamaMembersRequest{
						PageToken: pageToken,
						PageSize:  0,
					})
					Expect(err).ShouldNot(HaveOccurred())
					Expect(listRes.ChamaMembers).ShouldNot(BeNil())

					nextPageToken = listRes.NextPageToken
					pageToken = nextPageToken
				}
			})
			It("should succeed is larger than default", func() {
				var pageToken, nextPageToken = "", "test"
				for nextPageToken != "" {
					listRes, err := ChamaMemberAPI.ListChamaMembers(ctx, &chama.ListChamaMembersRequest{
						PageToken: pageToken,
						PageSize:  1000,
					})
					Expect(err).ShouldNot(HaveOccurred())
					Expect(listRes.ChamaMembers).ShouldNot(BeNil())

					nextPageToken = listRes.NextPageToken
					pageToken = nextPageToken
				}
			})
		})

		Describe("ListChamaMembers with filters", func() {
			chamaIDs := []string{"1"}

			It("should succeed when CHAMA_ID filter is on", func() {
				var pageToken, nextPageToken = "", "test"
				for nextPageToken != "" {
					listRes, err := ChamaMemberAPI.ListChamaMembers(ctx, &chama.ListChamaMembersRequest{
						PageToken: pageToken,
						PageSize:  defaultPageSize,
						Filter: &chama.ChamaMemberFilter{
							ChamaIds: chamaIDs,
						},
					})
					Expect(err).ShouldNot(HaveOccurred())
					Expect(listRes.ChamaMembers).ShouldNot(BeNil())

					nextPageToken = listRes.NextPageToken
					pageToken = nextPageToken

					for _, chamaPB := range listRes.ChamaMembers {
						Expect(chamaPB.ChamaId).Should(BeElementOf(chamaIDs))
					}
				}
			})
		})
	})
})
