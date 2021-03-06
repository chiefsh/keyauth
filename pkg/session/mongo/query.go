package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/infraboard/keyauth/pkg/session"
	"github.com/infraboard/mcube/exception"
)

func newQueryLoginLogRequest(req *session.QuerySessionRequest) (*querySessionRequest, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	return &querySessionRequest{
		QuerySessionRequest: req,
	}, nil
}

type querySessionRequest struct {
	*session.QuerySessionRequest
}

func (r *querySessionRequest) FindOptions() *options.FindOptions {
	pageSize := int64(r.PageSize)
	skip := int64(r.PageSize) * int64(r.PageNumber-1)

	opt := &options.FindOptions{
		Sort:  bson.D{{Key: "login_at", Value: -1}},
		Limit: &pageSize,
		Skip:  &skip,
	}

	return opt
}

func (r *querySessionRequest) FindFilter() bson.M {
	tk := r.GetToken()
	filter := bson.M{
		"domain": tk.Domain,
	}

	if r.Account != "" {
		filter["account"] = r.Account
	}

	if r.ApplicationID != "" {
		filter["application_id"] = r.ApplicationID
	}

	if r.LoginIP != "" {
		filter["login_ip"] = r.LoginIP
	}

	if r.LoginCity != "" {
		filter["city"] = r.LoginCity
	}

	if r.GrantType != "" {
		filter["grant_type"] = r.GrantType
	}

	loginAt := bson.A{}
	if r.StartLoginTime != nil {
		loginAt = append(loginAt, bson.M{"login_at": bson.M{"$gte": r.StartLoginTime}})
	}

	if r.EndLoginTime != nil {
		loginAt = append(loginAt, bson.M{"login_at": bson.M{"$lte": r.EndLoginTime}})
	}
	if len(loginAt) > 0 {
		filter["$and"] = loginAt
	}

	return filter
}

func newDescribeSession(req *session.DescribeSessionRequest) (*describeSessionRequest, error) {
	if err := req.Validate(); err != nil {
		return nil, exception.NewBadRequest(err.Error())
	}

	return &describeSessionRequest{req}, nil
}

type describeSessionRequest struct {
	*session.DescribeSessionRequest
}

func (r *describeSessionRequest) FindOptions() *options.FindOneOptions {
	opt := &options.FindOneOptions{
		Sort: bson.D{{Key: "login_at", Value: -1}},
	}

	return opt
}

func (r *describeSessionRequest) FindFilter() bson.M {
	filter := bson.M{}

	if r.SessionID != "" {
		filter["_id"] = r.SessionID
	}
	if r.Domain != "" {
		filter["domain"] = r.Domain
	}
	if r.Account != "" {
		filter["account"] = r.Account
	}
	if r.Login {
		filter["logout_at"] = 0
	}

	return filter
}

func newQueryUserLastSessionRequest(req *session.QueryUserLastSessionRequest) (*queryUserLastSessionRequest, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	return &queryUserLastSessionRequest{
		QueryUserLastSessionRequest: req,
		pageSize:                    1,
	}, nil
}

type queryUserLastSessionRequest struct {
	pageSize int64
	*session.QueryUserLastSessionRequest
}

func (r *queryUserLastSessionRequest) FindOptions() *options.FindOneOptions {
	opt := &options.FindOneOptions{
		Sort: bson.D{{Key: "login_at", Value: -1}},
	}

	return opt
}

func (r *queryUserLastSessionRequest) FindFilter() bson.M {
	filter := bson.M{}

	if r.Account != "" {
		filter["account"] = r.Account
	}

	return filter
}
