package loan

import (
	"context"
	"errors"
	"fmt"

	"github.com/gidyon/machama-app/internal/models"
	"github.com/gidyon/machama-app/pkg/api/loan"
	"github.com/gidyon/machama-app/pkg/api/transaction"
	"github.com/gidyon/micro/v2/pkg/middleware/grpc/auth"
	"github.com/gidyon/micro/v2/utils/errs"
	"github.com/gidyon/services/pkg/utils/mdutil"

	"github.com/speps/go-hashids"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type Options struct {
	MoneyAccountAPI transaction.ChamaAccountAPIServer
	TransactionAPI  transaction.TransactionAPIServer
	SQLDB           *gorm.DB
	PageHasher      *hashids.HashID
	Logger          grpclog.LoggerV2
	Auth            auth.API
	AllowedGroups   []string
}

type loanAPIServer struct {
	loan.UnimplementedLoanAPIServer
	*Options
}

func NewLoanAPI(ctx context.Context, opt *Options) (loan.LoanAPIServer, error) {
	// Validation
	switch {
	case ctx == nil:
		return nil, errors.New("missing context")
	case opt == nil:
		return nil, errors.New("missing options")
	case opt.TransactionAPI == nil:
		return nil, errors.New("missing transaction API")
	case opt.MoneyAccountAPI == nil:
		return nil, errors.New("missing chama account API")
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

	loanAPI := &loanAPIServer{
		Options: opt,
	}

	return loanAPI, nil
}

func ValidateLoan(pb *loan.Loan) error {
	switch {
	case pb == nil:
		return errs.MissingField("loan")
	case pb.ChamaId == "":
		return errs.MissingField("loan id")
	case pb.ProductId == "":
		return errs.MissingField("plan id")
	case pb.MemberId == "":
		return errs.MissingField("member id")
	case pb.LoaneeNames == "":
		return errs.MissingField("loanee names")
	case pb.LoaneePhone == "":
		return errs.MissingField("loanee phone")
	case pb.NationalId == "":
		return errs.MissingField("loanee national id")
	}
	return nil
}

func (loanAPI *loanAPIServer) CreateLoan(
	ctx context.Context, req *loan.CreateLoanRequest,
) (*emptypb.Empty, error) {
	// Authorization
	_, err := loanAPI.Auth.AuthorizeGroup(ctx, loanAPI.AllowedGroups...)
	if err != nil {
		return nil, err
	}

	// Validate
	switch {
	case req == nil:
		return nil, errs.NilObject("request body")
	default:
		err = ValidateLoan(req.Loan)
		if err != nil {
			return nil, err
		}
	}

	db, err := models.LoanModel(req.Loan)
	if err != nil {
		return nil, err
	}

	err = loanAPI.SQLDB.Create(db).Error
	if err != nil {
		return nil, errs.FailedToSave("loan", err)
	}

	return &emptypb.Empty{}, nil
}

func (loanAPI *loanAPIServer) UpdateLoan(
	ctx context.Context, req *loan.UpdateLoanRequest,
) (*emptypb.Empty, error) {
	// Authorization
	_, err := loanAPI.Auth.AuthorizeGroup(ctx, loanAPI.AllowedGroups...)
	if err != nil {
		return nil, err
	}

	// Validate
	switch {
	case req == nil:
		return nil, errs.MissingField("update request")
	case req.Loan == nil:
		return nil, errs.MissingField("loan")
	case req.Loan.LoanId == "":
		return nil, errs.MissingField("loan id")
	}

	db, err := models.LoanModel(req.Loan)
	if err != nil {
		return nil, err
	}

	err = loanAPI.SQLDB.Where("id = ?", req.Loan.LoanId).Updates(db).Error
	if err != nil {
		return nil, errs.FailedToUpdate("loan", err)
	}

	return &emptypb.Empty{}, nil
}

const defaultPageSize = 50

func (loanAPI *loanAPIServer) ListLoans(
	ctx context.Context, req *loan.ListLoansRequest,
) (*loan.ListLoansResponse, error) {
	// Authorization
	actor, err := loanAPI.Auth.AuthorizeGroup(ctx, loanAPI.AllowedGroups...)
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
		if !loanAPI.Auth.IsAdmin(actor.Group) {
			pageSize = defaultPageSize
		}
	}

	var ID uint
	pageToken := req.GetPageToken()
	if pageToken != "" {
		ids, err := loanAPI.PageHasher.DecodeInt64WithError(req.GetPageToken())
		if err != nil {
			return nil, errs.WrapErrorWithCodeAndMsg(codes.InvalidArgument, err, "failed to parse page token")
		}
		ID = uint(ids[0])
	}

	db := loanAPI.SQLDB.Unscoped().Limit(int(pageSize + 1)).Order("id DESC")
	if ID != 0 {
		db = db.Where("id<?", ID)
	}

	// Apply tv filters
	if req.Filter != nil {
		if len(req.Filter.ChamaIds) != 0 {
			db = db.Where("chama_id IN (?)", req.Filter.ChamaIds)
		}
		if len(req.Filter.ProductIds) != 0 {
			db = db.Where("product_id IN (?)", req.Filter.ProductIds)
		}
	}

	dbs := make([]*models.Loan, 0, pageSize+1)
	err = db.Find(&dbs).Error
	switch {
	case err == nil:
	default:
		return nil, errs.SQLQueryFailed(err, "LIST")
	}

	pbs := make([]*loan.Loan, 0, len(dbs))
	for i, db := range dbs {
		if i == int(pageSize) {
			break
		}

		pb, err := models.LoanProto(db)
		if err != nil {
			return nil, err
		}

		pbs = append(pbs, pb)

		ID = db.ID
	}

	var token string
	if len(dbs) > int(pageSize) {
		// Next page token
		token, err = loanAPI.PageHasher.EncodeInt64([]int64{int64(ID)})
		if err != nil {
			return nil, errs.WrapErrorWithCodeAndMsg(codes.InvalidArgument, err, "failed to generate next page token")
		}
	}

	return &loan.ListLoansResponse{
		NextPageToken: token,
		Loans:         pbs,
	}, nil
}

func (loanAPI *loanAPIServer) GetLoan(
	ctx context.Context, req *loan.GetLoanRequest,
) (*loan.Loan, error) {
	// Authorization
	_, err := loanAPI.Auth.AuthorizeGroup(ctx, loanAPI.AllowedGroups...)
	if err != nil {
		return nil, err
	}

	// Validation
	switch {
	case req == nil:
		return nil, errs.MissingField("request body")
	case req.LoanId == "":
		return nil, errs.MissingField("loan id")
	}

	db := &models.Loan{}

	err = loanAPI.SQLDB.First(db, "id = ?", req.LoanId).Error
	switch {
	case err == nil:
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, errs.DoesNotExist("loan", req.LoanId)
	default:
		return nil, errs.FailedToFind("loan", err)
	}

	return models.LoanProto(db)
}

func (loanAPI *loanAPIServer) ApproveLoan(
	ctx context.Context, req *loan.ApproveLoanRequest,
) (*emptypb.Empty, error) {
	// Authorization
	actor, err := loanAPI.Auth.AuthorizeGroup(ctx, loanAPI.AllowedGroups...)
	if err != nil {
		return nil, err
	}

	// Validation
	switch {
	case req == nil:
		return nil, errs.MissingField("request body")
	case req.LoanId == "":
		return nil, errs.MissingField("loan id")
	}

	// Get loan
	loanPB, err := loanAPI.GetLoan(ctx, &loan.GetLoanRequest{LoanId: req.LoanId})
	if err != nil {
		return nil, err
	}

	if loanPB.Status == loan.LoanStatus_FUNDS_TRANSFERED {
		return nil, errs.WrapMessage(codes.AlreadyExists, "loan funds have been disbursed")
	}

	// Update loan
	err = loanAPI.SQLDB.Model(&models.Loan{}).Updates(map[string]interface{}{
		"approved": true,
		"status":   loan.LoanStatus_APPROVED.String(),
	}).Error
	if err != nil {
		return nil, errs.FailedToUpdate("loan", err)
	}

	ctxExt := mdutil.AddFromCtx(ctx)

	// Get account
	accountPB, err := loanAPI.MoneyAccountAPI.GetChamaAccount(ctxExt, &transaction.GetChamaAccountRequest{
		OwnerId:     loanPB.ChamaId,
		AccountName: req.AccountName,
	})
	if err != nil {
		return nil, err
	}

	// Withdraw from account
	_, err = loanAPI.TransactionAPI.Withdraw(ctxExt, &transaction.WithdrawRequest{
		ActorId:     actor.ID,
		AccountId:   accountPB.AccountId,
		Description: fmt.Sprintf("Loan approval for %s", loanPB.LoaneeNames),
		Amount:      loanPB.LoanAmount,
	})
	if err != nil {
		return nil, err
	}

	// Update loan
	err = loanAPI.SQLDB.Model(&models.Loan{}).Updates(map[string]interface{}{
		"approved": true,
		"status":   loan.LoanStatus_FUNDS_WITHDRAWN_ACCOUNT.String(),
	}).Error
	if err != nil {
		return nil, errs.FailedToUpdate("loan", err)
	}

	// B2C Transfer

	// Update loan
	// err = loanAPI.SQLDB.Model(&models.Loan{}).Updates(map[string]interface{}{
	// 	"approved": true,
	// 	"status":   loan.LoanStatus_WAITING_FUNDS_TRANSFER.String(),
	// }).Error
	// if err != nil {
	// 	return nil, errs.FailedToUpdate("loan", err)
	// }

	return &emptypb.Empty{}, nil
}
