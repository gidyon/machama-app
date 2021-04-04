package loan

import (
	"context"
	"fmt"

	"github.com/gidyon/machama-app/internal/models"
	"github.com/gidyon/machama-app/pkg/api/loan"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("GetLoan", func() {
	var (
		getReq *loan.GetLoanRequest
		ctx    context.Context
	)

	BeforeEach(func() {
		getReq = &loan.GetLoanRequest{
			LoanId: "1",
		}
		ctx = context.TODO()
	})

	Describe("GetLoan with malformed request", func() {
		It("should fail when the request is nil", func() {
			getReq = nil
			getRes, err := LoanAPI.GetLoan(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(getRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when loan id is missing", func() {
			getReq.LoanId = ""
			getRes, err := LoanAPI.GetLoan(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(getRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when loan does not exist", func() {
			getReq.LoanId = "oops"
			getRes, err := LoanAPI.GetLoan(ctx, getReq)
			Expect(err).Should(HaveOccurred())
			Expect(getRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.NotFound))
		})
	})

	Describe("GetLoan with well formed request", func() {
		var loanID string

		Context("Lets create loan first", func() {
			It("should succeed", func() {
				loanDB, err := models.LoanModel(mockLoan())
				Expect(err).ShouldNot(HaveOccurred())

				Expect(LoanAPIServer.SQLDB.Create(loanDB).Error).ShouldNot(HaveOccurred())

				loanID = fmt.Sprint(loanDB.ID)
			})
		})

		Describe("Getting the loan", func() {
			It("should succeed", func() {
				getRes, err := LoanAPI.GetLoan(ctx, &loan.GetLoanRequest{
					LoanId: loanID,
				})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(getRes).ShouldNot(BeNil())
				Expect(status.Code(err)).Should(Equal(codes.OK))
			})
		})
	})
})
