package transaction

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

type transactionAPIServer struct {
	transaction.UnimplementedTransactionAPIServer
	*Options
}

func NewTransactionAPI(ctx context.Context, opt *Options) (transaction.TransactionAPIServer, error) {
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

	transactionAPI := &transactionAPIServer{
		Options: opt,
	}

	return transactionAPI, nil
}

func (transactionAPI *transactionAPIServer) Deposit(ctx context.Context, req *transaction.DepositRequest) (*emptypb.Empty, error) {
	// Authorization
	_, err := transactionAPI.Auth.AuthorizeGroup(ctx, transactionAPI.AllowedGroups...)
	if err != nil {
		return nil, err
	}

	// Validation
	switch {
	case req == nil:
		return nil, errs.MissingField("request")
	case req.AccountId == "":
		return nil, errs.MissingField("account id")
	case req.ActorId == "":
		return nil, errs.MissingField("actor id")
	case req.Description == "":
		return nil, errs.MissingField("description")
	case req.Amount == 0:
		return nil, errs.MissingField("amount")
	case req.Amount < 0:
		return nil, errs.IncorrectVal("amount")
	}

	// Confirm account exist
	err = transactionAPI.SQLDB.First(&models.ChamaAccount{}, "id = ?", req.AccountId).Error
	switch {
	case err == nil:
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, errs.DoesExist("chama account", req.AccountId)
	default:
		return nil, errs.FailedToFind("chama account", err)
	}

	tx := transactionAPI.SQLDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return nil, errs.FailedToBeginTx(err)
	}

	// Create transaction
	err = tx.Create(&models.Transaction{
		ActorID:           req.ActorId,
		AccountID:         req.AccountId,
		Description:       req.Description,
		TransactionType:   transaction.TransactionType_DEPOSIT.String(),
		TransactionAmount: req.Amount,
	}).Error
	if err != nil {
		tx.Rollback()
		return nil, errs.FailedToSave("transaction", err)
	}

	// Deposit amount
	err = transactionAPI.SQLDB.Model(&models.ChamaAccount{}).Where("id = ?", req.AccountId).
		Updates(map[string]interface{}{
			"total_deposited_amount": gorm.Expr("total_deposited_amount + ?", req.Amount),
			"available_amount":       gorm.Expr("available_amount + ?", req.Amount),
			"last_deposited_amount":  req.Amount,
		}).Error
	if err != nil {
		tx.Rollback()
		return nil, errs.FailedToUpdate("account balance", err)
	}

	// Commit transaction
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, errs.FailedToCommitTx(err)
	}

	return &emptypb.Empty{}, nil
}

func (transactionAPI *transactionAPIServer) Withdraw(ctx context.Context, req *transaction.WithdrawRequest) (*emptypb.Empty, error) {
	// Authorization
	_, err := transactionAPI.Auth.AuthorizeGroup(ctx, transactionAPI.AllowedGroups...)
	if err != nil {
		return nil, err
	}

	// Validation
	switch {
	case req == nil:
		return nil, errs.MissingField("request")
	case req.AccountId == "":
		return nil, errs.MissingField("account id")
	case req.ActorId == "":
		return nil, errs.MissingField("actor id")
	case req.Description == "":
		return nil, errs.MissingField("description")
	case req.Amount == 0:
		return nil, errs.MissingField("amount")
	case req.Amount < 0:
		return nil, errs.IncorrectVal("amount")
	}

	// Confirm account exist
	err = transactionAPI.SQLDB.First(&models.ChamaAccount{}, "id = ?", req.AccountId).Error
	switch {
	case err == nil:
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, errs.DoesExist("chama account", req.AccountId)
	default:
		return nil, errs.FailedToFind("chama account", err)
	}

	tx := transactionAPI.SQLDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return nil, errs.FailedToBeginTx(err)
	}

	// Create transaction
	err = tx.Create(&models.Transaction{
		ActorID:           req.ActorId,
		AccountID:         req.AccountId,
		Description:       req.Description,
		TransactionType:   transaction.TransactionType_WITHDRAWAL.String(),
		TransactionAmount: req.Amount,
	}).Error
	if err != nil {
		tx.Rollback()
		return nil, errs.FailedToSave("transaction", err)
	}

	// Withdraw amount
	db := transactionAPI.SQLDB.Model(&models.ChamaAccount{}).
		Where("id = ? AND available_amount >= ? AND withdrawable", req.AccountId, req.Amount, true).
		Updates(map[string]interface{}{
			"total_withdrawn_amount": gorm.Expr("total_deposited_amount + ?", req.Amount),
			"available_amount":       gorm.Expr("available_amount - ?", req.Amount),
			"last_withdrawn_amount":  req.Amount,
		})
	if db.Error != nil {
		tx.Rollback()
		return nil, errs.FailedToUpdate("account balance", db.Error)
	}

	if db.RowsAffected == 0 {
		return nil, errs.WrapMessage(codes.FailedPrecondition, "insufficient amount")
	}

	// Commit transaction
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, errs.FailedToCommitTx(err)
	}

	return &emptypb.Empty{}, nil
}

const defaultPageSize = 50

func (transactionAPI *transactionAPIServer) ListTransactions(
	ctx context.Context, req *transaction.ListTransactionsRequest,
) (*transaction.ListTransactionsResponse, error) {
	// Authorization
	actor, err := transactionAPI.Auth.AuthorizeGroup(ctx, transactionAPI.AllowedGroups...)
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
		if !transactionAPI.Auth.IsAdmin(actor.Group) {
			pageSize = defaultPageSize
		}
	}

	var ID uint
	pageToken := req.GetPageToken()
	if pageToken != "" {
		ids, err := transactionAPI.PageHasher.DecodeInt64WithError(req.GetPageToken())
		if err != nil {
			return nil, errs.WrapErrorWithCodeAndMsg(codes.InvalidArgument, err, "failed to parse page token")
		}
		ID = uint(ids[0])
	}

	db := transactionAPI.SQLDB.Unscoped().Limit(int(pageSize + 1)).Order("id DESC")
	if ID != 0 {
		db = db.Where("id<?", ID)
	}

	// Apply filters
	if req.Filter != nil {
		if len(req.Filter.TransactionIds) != 0 {
			db = db.Where("id IN (?)", req.Filter.TransactionIds)
		}
		if len(req.Filter.ActorIds) != 0 {
			db = db.Where("actor_id IN (?)", req.Filter.ActorIds)
		}
		if len(req.Filter.AccountIds) != 0 {
			db = db.Where("account_id IN (?)", req.Filter.AccountIds)
		}
		if req.Filter.TransactionType != transaction.TransactionType_TRANSACTION_TYPE_UNSPECIFIED {
			db = db.Where("transaction_type = ?", req.Filter.TransactionType.String())
		}
	}

	dbs := make([]*models.Transaction, 0, pageSize+1)
	err = db.Find(&dbs).Error
	switch {
	case err == nil:
	default:
		return nil, errs.SQLQueryFailed(err, "LIST")
	}

	pbs := make([]*transaction.Transaction, 0, len(dbs))
	for i, db := range dbs {
		if i == int(pageSize) {
			break
		}

		pb, err := models.TransactionProto(db)
		if err != nil {
			return nil, err
		}

		pbs = append(pbs, pb)

		ID = db.ID
	}

	var token string
	if len(dbs) > int(pageSize) {
		// Next page token
		token, err = transactionAPI.PageHasher.EncodeInt64([]int64{int64(ID)})
		if err != nil {
			return nil, errs.WrapErrorWithCodeAndMsg(codes.InvalidArgument, err, "failed to generate next page token")
		}
	}

	return &transaction.ListTransactionsResponse{
		NextPageToken: token,
		Transactions:  pbs,
	}, nil
}

func (transactionAPI *transactionAPIServer) GetTransaction(
	ctx context.Context, req *transaction.GetTransactionRequest,
) (*transaction.Transaction, error) {
	// Authorization
	_, err := transactionAPI.Auth.AuthorizeGroup(ctx, transactionAPI.AllowedGroups...)
	if err != nil {
		return nil, err
	}

	// Validation
	switch {
	case req == nil:
		return nil, errs.MissingField("request body")
	case req.TransactionId == "":
		return nil, errs.MissingField("account id")
	}

	db := &models.Transaction{}

	err = transactionAPI.SQLDB.First(db, "id = ?", req.TransactionId).Error
	switch {
	case err == nil:
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, errs.DoesNotExist("transaction", req.TransactionId)
	default:
		return nil, errs.FailedToFind("transaction", err)
	}

	return models.TransactionProto(db)
}
