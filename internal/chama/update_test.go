package chama

import (
	"context"

	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/machama-app/internal/models"
	"github.com/gidyon/machama-app/pkg/api/chama"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("UpdateChama", func() {
	var (
		createReq *chama.UpdateChamaRequest
		ctx       context.Context
	)

	BeforeEach(func() {
		chamaPB, err := models.ChamaProto(mockChama())
		Expect(err).ShouldNot(HaveOccurred())
		createReq = &chama.UpdateChamaRequest{
			Chama: chamaPB,
		}
		ctx = context.TODO()
	})

	Describe("UpdateChama with malformed request", func() {
		It("should fail when the request is nil", func() {
			createReq = nil
			createRes, err := ChamaAPI.UpdateChama(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when chama is nil", func() {
			createReq.Chama = nil
			createRes, err := ChamaAPI.UpdateChama(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when chama id is missing", func() {
			createReq.Chama.ChamaId = ""
			createRes, err := ChamaAPI.UpdateChama(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
	})

	Describe("UpdateChama with wellformed request", func() {
		var initialChama, newChama *chama.Chama

		Context("Lets create chama", func() {
			It("should succeed", func() {
				chamaDB := mockChama()
				Expect(ChamaAPIServer.SQLDB.Create(chamaDB).Error).ShouldNot(HaveOccurred())

				v, err := models.ChamaProto(chamaDB)
				Expect(err).ShouldNot(HaveOccurred())
				initialChama = v

				newChama = &chama.Chama{
					ChamaId: initialChama.ChamaId,
				}
			})
		})

		Describe("Updating initial chama", func() {
			It("should update non zero fields", func() {
				newChama.AccountBalance = 200
				newChama.Name = randomdata.SillyName()
				newChama.Status = "just updated"

				updateRes, err := ChamaAPI.UpdateChama(ctx, &chama.UpdateChamaRequest{
					Chama: newChama,
				})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(updateRes).ShouldNot(BeNil())
				Expect(status.Code(err)).Should(Equal(codes.OK))
			})
		})

		Context("Lets get chama and check whether its updated", func() {
			It("should succeed in getting cham", func() {
				chamaPB, err := ChamaAPI.GetChama(ctx, &chama.GetChamaRequest{
					ChamaId: initialChama.ChamaId,
				})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(chamaPB).ShouldNot(BeNil())
				Expect(status.Code(err)).Should(Equal(codes.OK))

				// Assertions
				Expect(chamaPB.Name).Should(Equal(newChama.Name))
				Expect(chamaPB.Name).ShouldNot(Equal(initialChama.Name))

				Expect(chamaPB.Status).Should(Equal(newChama.Status))
				Expect(chamaPB.Status).ShouldNot(Equal(initialChama.Status))

				Expect(chamaPB.AccountBalance).Should(Equal(newChama.AccountBalance))
				Expect(chamaPB.AccountBalance).ShouldNot(Equal(initialChama.AccountBalance))
			})
		})
	})
})
