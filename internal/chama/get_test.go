package chama

import (
	"context"

	"github.com/gidyon/machama-app/internal/models"
	"github.com/gidyon/machama-app/pkg/api/chama"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("GetChama", func() {
	var (
		getReq *chama.GetChamaRequest
		ctx    context.Context
	)

	BeforeEach(func() {
		getReq = &chama.GetChamaRequest{
			ChamaId: "1",
		}
		ctx = context.TODO()
	})

	Describe("GetChama with malformed request", func() {
		It("should fail when the request is nil", func() {
			getReq = nil
			getRes, err := ChamaAPI.GetChama(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(getRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when chama id is missing", func() {
			getReq.ChamaId = ""
			getRes, err := ChamaAPI.GetChama(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(getRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when chama does not exist", func() {
			getReq.ChamaId = "oops"
			getRes, err := ChamaAPI.GetChama(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(getRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.NotFound))
		})
	})

	Describe("GetChama with well formed request", func() {
		var chamaID string

		Context("Lets create chama first", func() {
			It("should succeed", func() {
				chamaDB := mockChama()
				Expect(ChamaAPIServer.SQLDB.Create(chamaDB).Error).ShouldNot(HaveOccurred())

				v, err := models.ChamaProto(chamaDB)
				Expect(err).ShouldNot(HaveOccurred())

				chamaID = v.ChamaId
			})
		})

		Describe("Getting the chama", func() {
			It("should succeed", func() {
				getRes, err := ChamaAPI.GetChama(ctx, &chama.GetChamaRequest{
					ChamaId: chamaID,
				})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(getRes).ShouldNot(BeNil())
				Expect(status.Code(err)).Should(Equal(codes.OK))
			})
		})
	})
})
