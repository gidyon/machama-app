package chamamember

import (
	"context"

	"github.com/gidyon/machama-app/pkg/api/chama"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("CreateChamaMember", func() {
	var (
		createReq *chama.CreateChamaMemberRequest
		ctx       context.Context
	)

	BeforeEach(func() {
		createReq = &chama.CreateChamaMemberRequest{
			ChamaMember: mockChamaMember(),
		}
		ctx = context.TODO()
	})

	Describe("CreateChamaMember with malformed request", func() {
		It("should fail when the request is nil", func() {
			createReq = nil
			createRes, err := ChamaMemberAPI.CreateChamaMember(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when member is nil", func() {
			createReq.ChamaMember = nil
			createRes, err := ChamaMemberAPI.CreateChamaMember(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when member id is missing", func() {
			createReq.ChamaMember.ChamaId = ""
			createRes, err := ChamaMemberAPI.CreateChamaMember(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when member first name is missing", func() {
			createReq.ChamaMember.FirstName = ""
			createRes, err := ChamaMemberAPI.CreateChamaMember(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when member last name is missing", func() {
			createReq.ChamaMember.LastName = ""
			createRes, err := ChamaMemberAPI.CreateChamaMember(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when member phone is missing", func() {
			createReq.ChamaMember.Phone = ""
			createRes, err := ChamaMemberAPI.CreateChamaMember(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when member id number is missing", func() {
			createReq.ChamaMember.IdNumber = ""
			createRes, err := ChamaMemberAPI.CreateChamaMember(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
	})

	Describe("CreateChamaMember with wellformed request", func() {
		It("should succeed", func() {
			createRes, err := ChamaMemberAPI.CreateChamaMember(ctx, createReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(createRes).ShouldNot(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.OK))
		})
	})
})
