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

var _ = Describe("CreateChama", func() {
	var (
		createReq *chama.CreateChamaRequest
		ctx       context.Context
	)

	BeforeEach(func() {
		chamaPB, err := models.ChamaProto(mockChama())
		Expect(err).ShouldNot(HaveOccurred())
		createReq = &chama.CreateChamaRequest{
			Chama: chamaPB,
		}
		ctx = context.TODO()
	})

	Describe("CreateChama with malformed request", func() {
		It("should fail when the request is nil", func() {
			createReq = nil
			createRes, err := ChamaAPI.CreateChama(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when chama is nil", func() {
			createReq.Chama = nil
			createRes, err := ChamaAPI.CreateChama(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when chama name is missing", func() {
			createReq.Chama.Name = ""
			createRes, err := ChamaAPI.CreateChama(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when chama creator is missing", func() {
			createReq.Chama.CreatorId = ""
			createRes, err := ChamaAPI.CreateChama(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
	})

	Describe("CreateChama with wellformed request", func() {
		It("should succeed", func() {
			createRes, err := ChamaAPI.CreateChama(ctx, createReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(createRes).ShouldNot(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.OK))
		})
	})
})
