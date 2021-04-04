package moneyaccount

import (
	"context"
	"errors"

	"github.com/gidyon/machama-app/internal/models"
	"github.com/gidyon/machama-app/pkg/api/transaction"
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

type moneyAccountAPIServer struct {
	transaction.UnimplementedChamaAccountAPIServer
	*Options
}

func NewChamaAccountAPI(ctx context.Context, opt *Options) (transaction.ChamaAccountAPIServer, error) {
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

	moneyAccountAPI := &moneyAccountAPIServer{
		Options: opt,
	}

	return moneyAccountAPI, nil
}

func ValidateChamaAccount(pb *transaction.ChamaAccount) error {
	switch {
	case pb == nil:
		return errs.MissingField("moneyAccount")
	case pb.AccountName == "":
		return errs.MissingField("account name")
	case pb.OwnerId == "":
		return errs.MissingField("owner id")
	case pb.AccountType == transaction.AccountType_ACCOUNT_TYPE_UNSPECIFIED:
		return errs.MissingField("account type")
	}
	return nil
}

func (moneyAccountAPI *moneyAccountAPIServer) CreateChamaAccount(
	ctx context.Context, req *transaction.CreateChamaAccountRequest,
) (*emptypb.Empty, error) {
	// Authorization
	_, err := moneyAccountAPI.Auth.AuthorizeGroup(ctx, moneyAccountAPI.AllowedGroups...)
	if err != nil {
		return nil, err
	}

	// Validate
	switch {
	case req == nil:
		return nil, errs.NilObject("request body")
	default:
		err = ValidateChamaAccount(req.ChamaAccount)
		if err != nil {
			return nil, err
		}
	}

	db, err := models.ChamaAccountModel(req.ChamaAccount)
	if err != nil {
		return nil, err
	}

	err = moneyAccountAPI.SQLDB.Create(db).Error
	if err != nil {
		return nil, errs.FailedToSave("moneyAccount", err)
	}

	return &emptypb.Empty{}, nil
}

const defaultPageSize = 50

func (moneyAccountAPI *moneyAccountAPIServer) ListChamaAccounts(
	ctx context.Context, req *transaction.ListChamaAccountsRequest,
) (*transaction.ListChamaAccountsResponse, error) {
	// Authorization
	actor, err := moneyAccountAPI.Auth.AuthorizeGroup(ctx, moneyAccountAPI.AllowedGroups...)
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
		if !moneyAccountAPI.Auth.IsAdmin(actor.Group) {
			pageSize = defaultPageSize
		}
	}

	var ID uint
	pageToken := req.GetPageToken()
	if pageToken != "" {
		ids, err := moneyAccountAPI.PageHasher.DecodeInt64WithError(req.GetPageToken())
		if err != nil {
			return nil, errs.WrapErrorWithCodeAndMsg(codes.InvalidArgument, err, "failed to parse page token")
		}
		ID = uint(ids[0])
	}

	db := moneyAccountAPI.SQLDB.Unscoped().Limit(int(pageSize + 1)).Order("id DESC")
	if ID != 0 {
		db = db.Where("id<?", ID)
	}

	// Apply tv filters
	if req.Filter != nil {
		if len(req.Filter.AccountIds) != 0 {
			db = db.Where("id IN (?)", req.Filter.AccountIds)
		}
		if len(req.Filter.OwnerIds) != 0 {
			db = db.Where("owner_id IN (?)", req.Filter.OwnerIds)
		}
		if req.Filter.AccountType != transaction.AccountType_ACCOUNT_TYPE_UNSPECIFIED {
			db = db.Where("account_type IN (?)", req.Filter.AccountType.String())
		}
		switch {
		case req.Filter.NotWithdrawable && req.Filter.Withdrawable:
		case !req.Filter.NotWithdrawable && !req.Filter.Withdrawable:
		case req.Filter.Withdrawable:
			db = db.Where("withdrawable IN (?)", true)
		case req.Filter.NotWithdrawable:
			db = db.Where("withdrawable IN (?)", false)
		}
	}

	dbs := make([]*models.ChamaAccount, 0, pageSize+1)
	err = db.Find(&dbs).Error
	switch {
	case err == nil:
	default:
		return nil, errs.SQLQueryFailed(err, "LIST")
	}

	pbs := make([]*transaction.ChamaAccount, 0, len(dbs))
	for i, db := range dbs {
		if i == int(pageSize) {
			break
		}

		pb, err := models.ChamaAccountProto(db)
		if err != nil {
			return nil, err
		}

		pbs = append(pbs, pb)

		ID = db.ID
	}

	var token string
	if len(dbs) > int(pageSize) {
		// Next page token
		token, err = moneyAccountAPI.PageHasher.EncodeInt64([]int64{int64(ID)})
		if err != nil {
			return nil, errs.WrapErrorWithCodeAndMsg(codes.InvalidArgument, err, "failed to generate next page token")
		}
	}

	return &transaction.ListChamaAccountsResponse{
		NextPageToken: token,
		ChamaAccounts: pbs,
	}, nil
}

func (moneyAccountAPI *moneyAccountAPIServer) GetChamaAccount(
	ctx context.Context, req *transaction.GetChamaAccountRequest,
) (*transaction.ChamaAccount, error) {
	// Authorization
	_, err := moneyAccountAPI.Auth.AuthorizeGroup(ctx, moneyAccountAPI.AllowedGroups...)
	if err != nil {
		return nil, err
	}

	// Validation
	switch {
	case req == nil:
		return nil, errs.MissingField("request body")
	case req.AccountId == "" && (req.OwnerId == "" || req.AccountName == ""):
		return nil, errs.MissingField("account id")
	}

	db := &models.ChamaAccount{}

	if req.AccountId != "" {
		err = moneyAccountAPI.SQLDB.First(db, "id = ?", req.AccountId).Error
	} else {
		err = moneyAccountAPI.SQLDB.First(db, "owner_id = ? AND account_name = ?", req.OwnerId, req.AccountName).Error
	}
	switch {
	case err == nil:
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, errs.DoesNotExist("moneyAccount", req.AccountId)
	default:
		return nil, errs.FailedToFind("moneyAccount", err)
	}

	return models.ChamaAccountProto(db)
}
