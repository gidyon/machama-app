package loanproduct

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

var _ = Describe("UpdateLoanProduct", func() {
	var (
		updateReq *loan.UpdateLoanProductRequest
		ctx       context.Context
	)

	BeforeEach(func() {
		updateReq = &loan.UpdateLoanProductRequest{
			LoanProduct: mockLoanProduct(),
		}
		ctx = context.TODO()
	})

	Describe("UpdateLoanProduct with malformed request", func() {
		It("should fail when the request is nil", func() {
			updateReq = nil
			updateRes, err := LoanProductAPI.UpdateLoanProduct(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(updateRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when loan product is nil", func() {
			updateReq.LoanProduct = nil
			updateRes, err := LoanProductAPI.UpdateLoanProduct(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(updateRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when loan product id is missing", func() {
			updateReq.LoanProduct.ProductId = ""
			updateRes, err := LoanProductAPI.UpdateLoanProduct(ctx, updateReq)
			Expect(err).Should(HaveOccurred())
			Expect(updateRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
	})

	Describe("UpdateLoanProduct with wellformed request", func() {
		var initialProduct, newProduct *loan.LoanProduct

		Context("Lets create loan product", func() {
			It("should succeed", func() {
				productPB := mockLoanProduct()
				v, err := models.LoanProductModel(productPB)
				Expect(err).ShouldNot(HaveOccurred())

				Expect(LoanProductAPIServer.SQLDB.Create(v).Error).ShouldNot(HaveOccurred())

				initialProduct = productPB

				newProduct = &loan.LoanProduct{
					ProductId: fmt.Sprint(v.ID),
				}
			})
		})

		Describe("Updating initial loan", func() {
			It("should update non zero fields", func() {
				newProduct.Name = randomdata.SillyName()
				newProduct.Description = randomDescription()
				newProduct.InterestRate = float32(randomdata.Decimal(3, 20))

				updateRes, err := LoanProductAPI.UpdateLoanProduct(ctx, &loan.UpdateLoanProductRequest{
					LoanProduct: newProduct,
				})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(updateRes).ShouldNot(BeNil())
				Expect(status.Code(err)).Should(Equal(codes.OK))
			})
		})

		Context("Lets get product and check whether its updated", func() {
			It("should succeed in getting loan", func() {
				productPB, err := LoanProductAPI.GetLoanProduct(ctx, &loan.GetLoanProductRequest{
					ProductId: newProduct.ProductId,
				})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(productPB).ShouldNot(BeNil())
				Expect(status.Code(err)).Should(Equal(codes.OK))

				// Assertions
				Expect(productPB.Name).Should(Equal(newProduct.Name))
				Expect(productPB.Name).ShouldNot(Equal(initialProduct.Name))

				Expect(productPB.Description).Should(Equal(newProduct.Description))
				Expect(productPB.Description).ShouldNot(Equal(initialProduct.Description))

				Expect(productPB.InterestRate).Should(Equal(newProduct.InterestRate))
				Expect(productPB.InterestRate).ShouldNot(Equal(initialProduct.InterestRate))
			})
		})
	})
})
