package chama

import (
	"context"
	"errors"

	"github.com/gidyon/machama-app/internal/models"
	"github.com/gidyon/machama-app/pkg/api/chama"
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

type chamaAPIServer struct {
	chama.UnimplementedChamaAPIServer
	*Options
}

func NewChamaAPI(ctx context.Context, opt *Options) (chama.ChamaAPIServer, error) {
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

	chamaAPI := &chamaAPIServer{
		Options: opt,
	}

	return chamaAPI, nil
}

func ValidateChama(pb *chama.Chama) error {
	switch {
	case pb == nil:
		return errs.MissingField("chama")
	case pb.CreatorId == "":
		return errs.MissingField("creator id")
	case pb.Name == "":
		return errs.MissingField("chama name")
	}
	return nil
}

func (chamaAPI *chamaAPIServer) CreateChama(
	ctx context.Context, req *chama.CreateChamaRequest,
) (*emptypb.Empty, error) {
	// Authorization
	_, err := chamaAPI.Auth.AuthorizeAdmin(ctx)
	if err != nil {
		return nil, err
	}

	// Validate
	switch {
	case req == nil:
		return nil, errs.MissingField("create request")
	default:
		err = ValidateChama(req.Chama)
		if err != nil {
			return nil, err
		}
	}

	db, err := models.ChamaModel(req.Chama)
	if err != nil {
		return nil, err
	}

	err = chamaAPI.SQLDB.Create(db).Error
	if err != nil {
		return nil, errs.FailedToSave("chama", err)
	}

	return &emptypb.Empty{}, nil
}

func (chamaAPI *chamaAPIServer) UpdateChama(
	ctx context.Context, req *chama.UpdateChamaRequest,
) (*emptypb.Empty, error) {
	// Authorization
	_, err := chamaAPI.Auth.AuthorizeGroup(ctx, chamaAPI.AllowedGroups...)
	if err != nil {
		return nil, err
	}

	// Validate
	switch {
	case req == nil:
		return nil, errs.MissingField("update request")
	case req.Chama == nil:
		return nil, errs.MissingField("chama")
	case req.Chama.GetChamaId() == "":
		return nil, errs.MissingField("chama id")
	}

	db, err := models.ChamaModel(req.Chama)
	if err != nil {
		return nil, err
	}

	err = chamaAPI.SQLDB.Where("id = ?", req.Chama.ChamaId).Updates(db).Error
	if err != nil {
		return nil, errs.FailedToUpdate("chama", err)
	}

	return &emptypb.Empty{}, nil
}

const defaultPageSize = 50

func (chamaAPI *chamaAPIServer) ListChamas(
	ctx context.Context, req *chama.ListChamasRequest,
) (*chama.ListChamasResponse, error) {
	// Authorization
	actor, err := chamaAPI.Auth.AuthorizeGroup(ctx, chamaAPI.AllowedGroups...)
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
		if !chamaAPI.Auth.IsAdmin(actor.Group) {
			pageSize = defaultPageSize
		}
	}

	var ID uint
	pageToken := req.GetPageToken()
	if pageToken != "" {
		ids, err := chamaAPI.PageHasher.DecodeInt64WithError(req.GetPageToken())
		if err != nil {
			return nil, errs.WrapErrorWithCodeAndMsg(codes.InvalidArgument, err, "failed to parse page token")
		}
		ID = uint(ids[0])
	}

	db := chamaAPI.SQLDB.Unscoped().Limit(int(pageSize + 1)).Order("id DESC")
	if ID != 0 {
		db = db.Where("id<?", ID)
	}

	// Apply tv filters
	if req.Filter != nil {
		if len(req.Filter.CreatorIds) != 0 {
			db = db.Where("creator_id IN (?)", req.Filter.CreatorIds)
		}
	}

	dbs := make([]*models.Chama, 0, pageSize+1)
	err = db.Find(&dbs).Error
	switch {
	case err == nil:
	default:
		return nil, errs.SQLQueryFailed(err, "LIST")
	}

	pbs := make([]*chama.Chama, 0, len(dbs))
	for i, db := range dbs {
		if i == int(pageSize) {
			break
		}

		pb, err := models.ChamaProto(db)
		if err != nil {
			return nil, err
		}

		pbs = append(pbs, pb)

		ID = db.ID
	}

	var token string
	if len(dbs) > int(pageSize) {
		// Next page token
		token, err = chamaAPI.PageHasher.EncodeInt64([]int64{int64(ID)})
		if err != nil {
			return nil, errs.WrapErrorWithCodeAndMsg(codes.InvalidArgument, err, "failed to generate next page token")
		}
	}

	return &chama.ListChamasResponse{
		NextPageToken: token,
		Chamas:        pbs,
	}, nil
}

func (chamaAPI *chamaAPIServer) GetChama(
	ctx context.Context, req *chama.GetChamaRequest,
) (*chama.Chama, error) {
	// Authorization
	_, err := chamaAPI.Auth.AuthorizeGroup(ctx, chamaAPI.AllowedGroups...)
	if err != nil {
		return nil, err
	}

	// Validation
	switch {
	case req == nil:
		return nil, errs.MissingField("request body")
	case req.ChamaId == "":
		return nil, errs.MissingField("chama id")
	}

	db := &models.Chama{}

	err = chamaAPI.SQLDB.First(db, "id = ?", req.ChamaId).Error
	switch {
	case err == nil:
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, errs.DoesNotExist("chama", req.ChamaId)
	default:
		return nil, errs.FailedToFind("chama", err)
	}

	return models.ChamaProto(db)
}
