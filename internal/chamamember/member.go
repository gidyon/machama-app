package chamamember

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

type chamaMemberAPIServer struct {
	chama.UnimplementedChamaMemberAPIServer
	*Options
}

func NewChamaMemberAPI(ctx context.Context, opt *Options) (chama.ChamaMemberAPIServer, error) {
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

	chamaMemberAPI := &chamaMemberAPIServer{
		Options: opt,
	}

	return chamaMemberAPI, nil
}

func ValidateChamaMember(pb *chama.ChamaMember) error {
	switch {
	case pb == nil:
		return errs.MissingField("chama member")
	case pb.ChamaId == "":
		return errs.MissingField("chama id")
	case pb.FirstName == "":
		return errs.MissingField("first name")
	case pb.LastName == "":
		return errs.MissingField("last name")
	case pb.Phone == "":
		return errs.MissingField("phone")
	case pb.IdNumber == "":
		return errs.MissingField("id number")
	}
	return nil
}

func (chamaMemberAPI *chamaMemberAPIServer) CreateChamaMember(
	ctx context.Context, req *chama.CreateChamaMemberRequest,
) (*emptypb.Empty, error) {
	// Authorization
	_, err := chamaMemberAPI.Auth.AuthorizeGroup(ctx, chamaMemberAPI.AllowedGroups...)
	if err != nil {
		return nil, err
	}

	// Validate
	switch {
	case req == nil:
		return nil, errs.NilObject("request body")
	default:
		err = ValidateChamaMember(req.ChamaMember)
		if err != nil {
			return nil, err
		}
	}

	db, err := models.ChamaMemberModel(req.ChamaMember)
	if err != nil {
		return nil, err
	}

	err = chamaMemberAPI.SQLDB.Create(db).Error
	if err != nil {
		return nil, errs.FailedToSave("chamaMember", err)
	}

	return &emptypb.Empty{}, nil
}

func (chamaMemberAPI *chamaMemberAPIServer) UpdateChamaMember(
	ctx context.Context, req *chama.UpdateChamaMemberRequest,
) (*emptypb.Empty, error) {
	// Authorization
	_, err := chamaMemberAPI.Auth.AuthorizeGroup(ctx, chamaMemberAPI.AllowedGroups...)
	if err != nil {
		return nil, err
	}

	// Validate
	switch {
	case req == nil:
		return nil, errs.MissingField("update request")
	case req.ChamaMember == nil:
		return nil, errs.MissingField("chamaMember")
	case req.ChamaMember.MemberId == "":
		return nil, errs.MissingField("member id")
	}

	db, err := models.ChamaMemberModel(req.ChamaMember)
	if err != nil {
		return nil, err
	}

	err = chamaMemberAPI.SQLDB.Where("id = ?", req.ChamaMember.MemberId).Updates(db).Error
	if err != nil {
		return nil, errs.FailedToUpdate("chamaMember", err)
	}

	return &emptypb.Empty{}, nil
}

const defaultPageSize = 50

func (chamaMemberAPI *chamaMemberAPIServer) ListChamaMembers(
	ctx context.Context, req *chama.ListChamaMembersRequest,
) (*chama.ListChamaMembersResponse, error) {
	// Authorization
	actor, err := chamaMemberAPI.Auth.AuthorizeGroup(ctx, chamaMemberAPI.AllowedGroups...)
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
		if !chamaMemberAPI.Auth.IsAdmin(actor.Group) {
			pageSize = defaultPageSize
		}
	}

	var ID uint
	pageToken := req.GetPageToken()
	if pageToken != "" {
		ids, err := chamaMemberAPI.PageHasher.DecodeInt64WithError(req.GetPageToken())
		if err != nil {
			return nil, errs.WrapErrorWithCodeAndMsg(codes.InvalidArgument, err, "failed to parse page token")
		}
		ID = uint(ids[0])
	}

	db := chamaMemberAPI.SQLDB.Unscoped().Limit(int(pageSize + 1)).Order("id DESC")
	if ID != 0 {
		db = db.Where("id<?", ID)
	}

	// Apply tv filters
	if req.Filter != nil {
		if len(req.Filter.ChamaIds) != 0 {
			db = db.Where("chama_id IN (?)", req.Filter.ChamaIds)
		}
	}

	dbs := make([]*models.ChamaMember, 0, pageSize+1)
	err = db.Find(&dbs).Error
	switch {
	case err == nil:
	default:
		return nil, errs.SQLQueryFailed(err, "LIST")
	}

	pbs := make([]*chama.ChamaMember, 0, len(dbs))
	for i, db := range dbs {
		if i == int(pageSize) {
			break
		}

		pb, err := models.ChamaMemberProto(db)
		if err != nil {
			return nil, err
		}

		pbs = append(pbs, pb)

		ID = db.ID
	}

	var token string
	if len(dbs) > int(pageSize) {
		// Next page token
		token, err = chamaMemberAPI.PageHasher.EncodeInt64([]int64{int64(ID)})
		if err != nil {
			return nil, errs.WrapErrorWithCodeAndMsg(codes.InvalidArgument, err, "failed to generate next page token")
		}
	}

	return &chama.ListChamaMembersResponse{
		NextPageToken: token,
		ChamaMembers:  pbs,
	}, nil
}

func (chamaMemberAPI *chamaMemberAPIServer) GetChamaMember(
	ctx context.Context, req *chama.GetChamaMemberRequest,
) (*chama.ChamaMember, error) {
	// Authorization
	_, err := chamaMemberAPI.Auth.AuthorizeGroup(ctx, chamaMemberAPI.AllowedGroups...)
	if err != nil {
		return nil, err
	}

	// Validation
	switch {
	case req == nil:
		return nil, errs.MissingField("request body")
	case req.MemberId == "":
		return nil, errs.MissingField("member id")
	}

	db := &models.ChamaMember{}

	err = chamaMemberAPI.SQLDB.First(db, "id = ?", req.MemberId).Error
	switch {
	case err == nil:
	case errors.Is(err, gorm.ErrRecordNotFound):
		return nil, errs.DoesNotExist("chamaMember", req.MemberId)
	default:
		return nil, errs.FailedToFind("chamaMember", err)
	}

	return models.ChamaMemberProto(db)
}

func (chamaMemberAPI *chamaMemberAPIServer) DeleteChamaMember(
	ctx context.Context, req *chama.DeleteChamaMemberRequest,
) (*emptypb.Empty, error) {
	// Authorization
	_, err := chamaMemberAPI.Auth.AuthorizeGroup(ctx, chamaMemberAPI.AllowedGroups...)
	if err != nil {
		return nil, err
	}

	// Validation
	switch {
	case req == nil:
		return nil, errs.MissingField("request body")
	case req.MemberId == "":
		return nil, errs.MissingField("member id")
	}

	err = chamaMemberAPI.SQLDB.Delete(&models.ChamaMember{}, "id = ?", req.MemberId).Error
	if err != nil {
		return nil, errs.FailedToDelete("chama member", err)
	}

	return &emptypb.Empty{}, nil
}
