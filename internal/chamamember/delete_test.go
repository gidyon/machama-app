package chamamember

import (
	"context"
	"fmt"

	"github.com/gidyon/machama-app/internal/models"
	"github.com/gidyon/machama-app/pkg/api/chama"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("DeleteChamaMember", func() {
	var (
		getReq *chama.DeleteChamaMemberRequest
		ctx    context.Context
	)

	BeforeEach(func() {
		getReq = &chama.DeleteChamaMemberRequest{
			MemberId: "1",
		}
		ctx = context.TODO()
	})

	Describe("DeleteChamaMember with malformed request", func() {
		It("should fail when the request is nil", func() {
			getReq = nil
			getRes, err := ChamaMemberAPI.DeleteChamaMember(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(getRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when chama member id is missing", func() {
			getReq.MemberId = ""
			getRes, err := ChamaMemberAPI.DeleteChamaMember(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(getRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
	})

	Describe("DeleteChamaMember with well formed request", func() {
		var chamaMemberID string

		Context("Lets create chama member first", func() {
			It("should succeed", func() {
				chamaMemberDB, err := models.ChamaMemberModel(mockChamaMember())
				Expect(err).ShouldNot(HaveOccurred())

				Expect(ChamaMemberAPIServer.SQLDB.Create(chamaMemberDB).Error).ShouldNot(HaveOccurred())

				chamaMemberID = fmt.Sprint(chamaMemberDB.ID)
			})
		})

		Describe("Deleting the chama member", func() {
			It("should succeed", func() {
				getRes, err := ChamaMemberAPI.DeleteChamaMember(ctx, &chama.DeleteChamaMemberRequest{
					MemberId: chamaMemberID,
				})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(getRes).ShouldNot(BeNil())
				Expect(status.Code(err)).Should(Equal(codes.OK))
			})
		})

		Describe("Getting the chama member", func() {
			It("should fail because we deleted them", func() {
				getRes, err := ChamaMemberAPI.GetChamaMember(ctx, &chama.GetChamaMemberRequest{
					MemberId: chamaMemberID,
				})
				Expect(err).Should(HaveOccurred())
				Expect(getRes).Should(BeNil())
				Expect(status.Code(err)).Should(Equal(codes.NotFound))
			})
		})
	})
})
