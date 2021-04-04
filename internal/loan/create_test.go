package loan

import (
	"context"

	"github.com/gidyon/machama-app/pkg/api/loan"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("CreateLoan", func() {
	var (
		createReq *loan.CreateLoanRequest
		ctx       context.Context
	)

	BeforeEach(func() {
		createReq = &loan.CreateLoanRequest{
			Loan: mockLoan(),
		}
		ctx = context.TODO()
	})

	Describe("CreateLoan with malformed request", func() {
		It("should fail when the request is nil", func() {
			createReq = nil
			createRes, err := LoanAPI.CreateLoan(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when loan is nil", func() {
			createReq.Loan = nil
			createRes, err := LoanAPI.CreateLoan(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when chama id is missing", func() {
			createReq.Loan.ChamaId = ""
			createRes, err := LoanAPI.CreateLoan(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when plan id is missing", func() {
			createReq.Loan.PlanId = ""
			createRes, err := LoanAPI.CreateLoan(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when member id is missing", func() {
			createReq.Loan.MemberId = ""
			createRes, err := LoanAPI.CreateLoan(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when loanee names is missing", func() {
			createReq.Loan.LoaneeNames = ""
			createRes, err := LoanAPI.CreateLoan(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when loanee email is missing", func() {
			createReq.Loan.LoaneePhone = ""
			createRes, err := LoanAPI.CreateLoan(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when loanee national id is missing", func() {
			createReq.Loan.NationalId = ""
			createRes, err := LoanAPI.CreateLoan(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
	})

	Describe("CreateLoan with wellformed request", func() {
		It("should succeed", func() {
			createRes, err := LoanAPI.CreateLoan(ctx, createReq)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(createRes).ShouldNot(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.OK))
		})
	})
})
