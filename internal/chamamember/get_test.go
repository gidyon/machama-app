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

var _ = Describe("GetChamaMember", func() {
	var (
		getReq *chama.GetChamaMemberRequest
		ctx    context.Context
	)

	BeforeEach(func() {
		getReq = &chama.GetChamaMemberRequest{
			MemberId: "1",
		}
		ctx = context.TODO()
	})

	Describe("GetChamaMember with malformed request", func() {
		It("should fail when the request is nil", func() {
			getReq = nil
			getRes, err := ChamaMemberAPI.GetChamaMember(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(getRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when chama member id is missing", func() {
			getReq.MemberId = ""
			getRes, err := ChamaMemberAPI.GetChamaMember(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(getRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when chama member does not exist", func() {
			getReq.MemberId = "oops"
			getRes, err := ChamaMemberAPI.GetChamaMember(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(getRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.NotFound))
		})
	})

	Describe("GetChamaMember with well formed request", func() {
		var chamaMemberID string

		Context("Lets create chama member first", func() {
			It("should succeed", func() {
				chamaMemberDB, err := models.ChamaMemberModel(mockChamaMember())
				Expect(err).ShouldNot(HaveOccurred())

				Expect(ChamaMemberAPIServer.SQLDB.Create(chamaMemberDB).Error).ShouldNot(HaveOccurred())

				chamaMemberID = fmt.Sprint(chamaMemberDB.ID)
			})
		})

		Describe("Getting the chama member", func() {
			It("should succeed", func() {
				getRes, err := ChamaMemberAPI.GetChamaMember(ctx, &chama.GetChamaMemberRequest{
					MemberId: chamaMemberID,
				})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(getRes).ShouldNot(BeNil())
				Expect(status.Code(err)).Should(Equal(codes.OK))
			})
		})
	})
})
