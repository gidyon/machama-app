package loanproduct

import (
	"context"
	"errors"

	"github.com/gidyon/machama-app/internal/models"
	"github.com/gidyon/machama-app/pkg/api/loan"
	"github.com/gidyon/micro/v2/pkg/middleware/grpc/auth"
	"github.com/gidyon/micro/v2/utils/errs"

	"github.com/speps/go-hashids"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type Options struct {
	SQLDB         *gorm.DB
	PageHasher    *hashids.HashID
	Logger        grpclog.LoggerV2
	Auth          auth.API
	AllowedGroups []string
}

type loanProductAPIServer struct {
	loan.UnimplementedLoanProductAPIServer
	*Options
}

func NewLoanProductAPI(ctx context.Context, opt *Options) (loan.LoanProductAPIServer, error) {
	// Validation
	switch {
	case ctx == nil:
		return nil, errors.New("missing context")
	case opt == nil:
		return nil, errors.New("missing options")
	case opt.SQLDB == nil:
		return nil, errors.New("missing sql db")
	case opt.Logger == nil:
		return nil, errors.New("missing logger")
	case opt.Auth == nil:
		return nil, errors.New("missing auth API")
	case opt.PageHasher == nil:
		return nil, errors.New("missing pagination hasher")
	default:
		if len(opt.AllowedGroups) == 0 {
			opt.AllowedGroups = opt.Auth.AdminGroups()
		}
	}

	LoanProductAPI := &loanProductAPIServer{
		Options: opt,
	}

	return LoanProductAPI, nil
}

func ValidateLoanProduct(pb *loan.LoanProduct) error {
	switch {
	case pb == nil:
	case pb.ChamaId == "":
		return errs.MissingField("chama id")
	case pb.Name == "":
		return errs.MissingField("plan name")
	case pb.InterestRate == 0:
		return errs.MissingField("interest rate")
	}
	return nil
}

func (LoanProductAPI *loanProductAPIServer) CreateLoanProduct(
	ctx context.Context, req *loan.CreateLoanProductRequest,
) (*emptypb.Empty, error) {
	// Authorization
	_, err := LoanProductAPI.Auth.AuthorizeGroup(ctx, LoanProductAPI.AllowedGroups...)
	if err != nil {
		return nil, err
	}

	// Validate
	switch {
	case req == nil:
		return nil, errs.NilObject("request body")
	default:
		err = ValidateLoanProduct(req.LoanProduct)
		if err != nil {
			return nil, err
		}
	}

	db, err := models.LoanProductModel(req.LoanProduct)
	if err != nil {
		return nil, err
	}

	err = LoanProductAPI.SQLDB.Create(db).Error
	if err != nil {
		return nil, errs.FailedToSave("LoanProduct", err)
	}

	return &emptypb.Empty{}, nil
}

func (LoanProductAPI *loanProductAPIServer) UpdateLoanProduct(
	ctx context.Context, req *loan.UpdateLoanProductRequest,
) (*emptypb.Empty, error) {
	// Authorization
	_, err := LoanProductAPI.Auth.AuthorizeGroup(ctx, LoanProductAPI.AllowedGroups...)
	if err != nil {
		return nil, err
	}

	// Validate
	switch {
	case req == nil:
		return nil, errs.MissingField("update request")
	case req.LoanProduct == nil:
		return nil, errs.MissingField("loan product")
	case req.LoanProduct.ProductId == "":
		return nil, errs.MissingField("product id")
	}

	db, err := models.LoanProductModel(req.LoanProduct)
	if err != nil {
		return nil, err
	}

	err = LoanProductAPI.SQLDB.Where("id = ?", req.LoanProduct.ProductId).Updates(db).Error
	if err != nil {
		return nil, errs.FailedToUpdate("LoanProduct", err)
	}

	return &emptypb.Empty{}, nil
}

const defaultPageSize = 50

func (LoanProductAPI *loanProductAPIServer) ListLoanProducts(
	ctx context.Context, req *loan.ListLoanProductsRequest,
) (*loan.ListLoanProductsResponse, error) {
	// Authorization
	actor, err := LoanProductAPI.Auth.AuthorizeGroup(ctx, LoanProductAPI.AllowedGroups...)
	if err != nil {
		return nil, err
	}

	// Validation
	switch {
	case req == nil:
		return nil, errs.NilObject("list request")
	}

	pageSize := req.GetPageSize()
	switch {
	case pageSize <= 0:
		pageSize = defaultPageSize
	case pageSize > defaultPageSize:
		if !LoanProductAPI.Auth.IsAdmin(actor.Group) {
			pageSize = defaultPageSize
		}
	}

	var ID uint
	pageToken := req.GetPageToken()
	if pageToken != "" {
		ids, err := LoanProductAPI.PageHasher.DecodeInt64WithError(req.GetPageToken())
		if err != nil {
			return nil, errs.WrapErrorWithCodeAndMsg(codes.InvalidArgument, err, "failed to parse page token")
		}
		ID = uint(ids[0])
	}

	db := LoanProductAPI.SQLDB.Unscoped().Limit(int(pageSize + 1)).Order("id DESC")
	if ID != 0 {
		db = db.Where("id<?", ID)
	}

	// Apply tv filters
	if req.Filter != nil {
		if len(req.Filter.ChamaIds) != 0 {
			db = db.Where("chama_id IN (?)", req.Filter.ChamaIds)
		}
	}

	dbs := make([]*models.LoanProduct, 0, pageSize+1)
	err = db.Find(&dbs).Error
	switch {
	case err == nil:
	default:
		return nil, errs.SQLQueryFailed(err, "LIST")
	}

	pbs := make([]*loan.LoanProduct, 0, len(dbs))
	for i, db := range dbs {
		if i == int(pageSize) {
			break
		}

		pb, err := models.LoanProductProto(db)
		if err != nil {
			return nil, err
		}

		pbs = append(pbs, pb)

		ID = db.ID
	}

	var token string
	if len(dbs) > int(pageSize) {
		// Next page token
		token, err = LoanProductAPI.PageHasher.EncodeInt64([]int64{int64(ID)})
		if err != nil {
			return nil, errs.WrapErrorWithCodeAndMsg(codes.InvalidArgument, err, "failed to generate next page token")
		}
	}

	return &loan.ListLoanProductsResponse{
		NextPageToken: token,
		LoanProducts:  pbs,
	}, nil
}

func (LoanProductAPI *loanProductAPIServer) GetLoanProduct(
	ctx context.Context, req *loan.GetLoanProductRequest,
) (*loan.LoanProduct, error) {
	// Authorization
	_, err := LoanProductAPI.Auth.AuthorizeGroup(ctx, LoanProductAPI.AllowedGroups...)
	if err != nil {
		return nil, err
	}

	// Validation
	switch {
	case req == nil:
		return nil, errs.MissingField("request body")
	case req.ProductId == "":
		return nil, errs.MissingField("plan id")
	}

	db := &models.LoanProduct{}

	err = LoanProductAPI.SQLDB.First(db, "id = ?", req.ProductId).Error
	switch {
	case err == nil:
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, errs.DoesNotExist("LoanProduct", req.ProductId)
	default:
		return nil, errs.FailedToFind("LoanProduct", err)
	}

	return models.LoanProductProto(db)
}

func (LoanProductAPI *loanProductAPIServer) DeleteLoanProduct(
	ctx context.Context, req *loan.DeleteLoanProductRequest,
) (*emptypb.Empty, error) {
	// Authorization
	_, err := LoanProductAPI.Auth.AuthorizeGroup(ctx, LoanProductAPI.AllowedGroups...)
	if err != nil {
		return nil, err
	}

	// Validation
	switch {
	case req == nil:
		return nil, errs.MissingField("request body")
	case req.ProductId == "":
		return nil, errs.MissingField("plan id")
	}

	err = LoanProductAPI.SQLDB.Delete(&models.LoanProduct{}, "id = ?", req.ProductId).Error
	switch {
	case err == nil:
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, errs.DoesNotExist("LoanProduct", req.ProductId)
	default:
		return nil, errs.FailedToFind("LoanProduct", err)
	}

	return &emptypb.Empty{}, nil
}
