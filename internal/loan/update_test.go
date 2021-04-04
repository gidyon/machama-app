package loan

import (
	"context"
	"fmt"

	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/machama-app/internal/models"
	"github.com/gidyon/machama-app/pkg/api/loan"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("UpdateLoan", func() {
	var (
		createReq *loan.UpdateLoanRequest
		ctx       context.Context
	)

	BeforeEach(func() {
		createReq = &loan.UpdateLoanRequest{
			Loan: mockLoan(),
		}
		ctx = context.TODO()
	})

	Describe("UpdateLoan with malformed request", func() {
		It("should fail when the request is nil", func() {
			createReq = nil
			createRes, err := LoanAPI.UpdateLoan(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when loan is nil", func() {
			createReq.Loan = nil
			createRes, err := LoanAPI.UpdateLoan(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when loan id is missing", func() {
			createReq.Loan.LoanId = ""
			createRes, err := LoanAPI.UpdateLoan(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
	})

	Describe("UpdateLoan with wellformed request", func() {
		var initialLoan, newLoan *loan.Loan

		Context("Lets create loan", func() {
			It("should succeed", func() {
				loandPB := mockLoan()
				v, err := models.LoanModel(loandPB)
				Expect(err).ShouldNot(HaveOccurred())

				Expect(LoanAPIServer.SQLDB.Create(v).Error).ShouldNot(HaveOccurred())

				initialLoan = loandPB

				newLoan = &loan.Loan{
					LoanId: fmt.Sprint(v.ID),
				}
			})
		})

		Describe("Updating initial loan", func() {
			It("should update non zero fields", func() {
				newLoan.LoaneeEmail = randomdata.Email()
				newLoan.LoaneePhone = randomPhone()
				newLoan.LoaneeNames = "just updated"

				updateRes, err := LoanAPI.UpdateLoan(ctx, &loan.UpdateLoanRequest{
					Loan: newLoan,
				})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(updateRes).ShouldNot(BeNil())
				Expect(status.Code(err)).Should(Equal(codes.OK))
			})
		})

		Context("Lets get loan and check whether its updated", func() {
			It("should succeed in getting loan", func() {
				loanPB, err := LoanAPI.GetLoan(ctx, &loan.GetLoanRequest{
					LoanId: newLoan.LoanId,
				})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(loanPB).ShouldNot(BeNil())
				Expect(status.Code(err)).Should(Equal(codes.OK))

				// Assertions
				Expect(loanPB.LoaneeEmail).Should(Equal(newLoan.LoaneeEmail))
				Expect(loanPB.LoaneeEmail).ShouldNot(Equal(initialLoan.LoaneeEmail))

				Expect(loanPB.LoaneePhone).Should(Equal(newLoan.LoaneePhone))
				Expect(loanPB.LoaneePhone).ShouldNot(Equal(initialLoan.LoaneePhone))

				Expect(loanPB.LoaneeNames).Should(Equal(newLoan.LoaneeNames))
				Expect(loanPB.LoaneeNames).ShouldNot(Equal(initialLoan.LoaneeNames))
			})
		})
	})
})
