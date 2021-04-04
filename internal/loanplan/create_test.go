package loanproduct

import (
	"context"

	"github.com/gidyon/machama-app/pkg/api/loan"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("CreateLoanProduct", func() {
	var (
		createReq *loan.CreateLoanProductRequest
		ctx       context.Context
	)

	BeforeEach(func() {
		createReq = &loan.CreateLoanProductRequest{
			LoanProduct: mockLoanProduct(),
		}
		ctx = context.TODO()
	})

	Describe("CreateLoanProduct with malformed request", func() {
		It("should fail when the request is nil", func() {
			createReq = nil
			createRes, err := LoanProductAPI.CreateLoanProduct(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when loan product is nil", func() {
			createReq.LoanProduct = nil
			createRes, err := LoanProductAPI.CreateLoanProduct(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when chama id is missing", func() {
			createReq.LoanProduct.ChamaId = ""
			createRes, err := LoanProductAPI.CreateLoanProduct(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when product name is missing", func() {
			createReq.LoanProduct.Name = ""
			createRes, err := LoanProductAPI.CreateLoanProduct(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when product insterest rate is missing", func() {
			createReq.LoanProduct.InterestRate = 0
			createRes, err := LoanProductAPI.CreateLoanProduct(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
	})

	Describe("CreateLoanProduct with wellformed request", func() {
		It("should succeed", func() {
			createRes, err := LoanProductAPI.CreateLoanProduct(ctx, createReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(createRes).ShouldNot(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.OK))
		})
	})
})
