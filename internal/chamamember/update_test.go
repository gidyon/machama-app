package chamamember

import (
	"context"
	"fmt"

	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/machama-app/internal/models"
	"github.com/gidyon/machama-app/pkg/api/chama"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("UpdateChamaMember", func() {
	var (
		createReq *chama.UpdateChamaMemberRequest
		ctx       context.Context
	)

	BeforeEach(func() {
		createReq = &chama.UpdateChamaMemberRequest{
			ChamaMember: mockChamaMember(),
		}
		ctx = context.TODO()
	})

	Describe("UpdateChamaMember with malformed request", func() {
		It("should fail when the request is nil", func() {
			createReq = nil
			createRes, err := ChamaMemberAPI.UpdateChamaMember(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when member is nil", func() {
			createReq.ChamaMember = nil
			createRes, err := ChamaMemberAPI.UpdateChamaMember(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
		It("should fail when member id is missing", func() {
			createReq.ChamaMember.ChamaId = ""
			createRes, err := ChamaMemberAPI.UpdateChamaMember(ctx, createReq)
			Expect(err).Should(HaveOccurred())
			Expect(createRes).Should(BeNil())
			Expect(status.Code(err)).Should(Equal(codes.InvalidArgument))
		})
	})

	Describe("UpdateChamaMember with wellformed request", func() {
		var initialChamaMember, newChamaMember *chama.ChamaMember

		Context("Lets create chama member", func() {
			It("should succeed", func() {
				chamaMemberPB := mockChamaMember()
				v, err := models.ChamaMemberModel(chamaMemberPB)
				Expect(err).ShouldNot(HaveOccurred())

				Expect(ChamaMemberAPIServer.SQLDB.Create(v).Error).ShouldNot(HaveOccurred())

				initialChamaMember = chamaMemberPB

				newChamaMember = &chama.ChamaMember{
					MemberId: fmt.Sprint(v.ID),
				}
			})
		})

		Describe("Updating initial chama", func() {
			It("should update non zero fields", func() {
				newChamaMember.Email = randomdata.Email()
				newChamaMember.Phone = randomPhone()
				newChamaMember.Status = "just updated"

				updateRes, err := ChamaMemberAPI.UpdateChamaMember(ctx, &chama.UpdateChamaMemberRequest{
					ChamaMember: newChamaMember,
				})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(updateRes).ShouldNot(BeNil())
				Expect(status.Code(err)).Should(Equal(codes.OK))
			})
		})

		Context("Lets get chama and check whether its updated", func() {
			It("should succeed in getting chama member", func() {
				chamaPB, err := ChamaMemberAPI.GetChamaMember(ctx, &chama.GetChamaMemberRequest{
					MemberId: newChamaMember.MemberId,
				})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(chamaPB).ShouldNot(BeNil())
				Expect(status.Code(err)).Should(Equal(codes.OK))

				// Assertions
				Expect(chamaPB.Email).Should(Equal(newChamaMember.Email))
				Expect(chamaPB.Email).ShouldNot(Equal(initialChamaMember.Email))

				Expect(chamaPB.Status).Should(Equal(newChamaMember.Status))
				Expect(chamaPB.Status).ShouldNot(Equal(initialChamaMember.Status))

				Expect(chamaPB.Phone).Should(Equal(newChamaMember.Phone))
				Expect(chamaPB.Phone).ShouldNot(Equal(initialChamaMember.Phone))
			})
		})
	})
})
