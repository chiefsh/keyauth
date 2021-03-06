package mongo

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/infraboard/keyauth/pkg/endpoint"
	"github.com/infraboard/mcube/exception"
)

func newDescribeEndpointRequest(req *endpoint.DescribeEndpointRequest) (*describeEndpointRequest, error) {
	if err := req.Validate(); err != nil {
		return nil, exception.NewBadRequest(err.Error())
	}

	return &describeEndpointRequest{req}, nil
}

type describeEndpointRequest struct {
	*endpoint.DescribeEndpointRequest
}

func (req *describeEndpointRequest) String() string {
	return fmt.Sprintf("endpoint: %s", req.ID)
}

func (req *describeEndpointRequest) FindFilter() bson.M {
	filter := bson.M{}
	if req.ID != "" {
		filter["_id"] = req.ID
	}

	return filter
}

func newQueryEndpointRequest(req *endpoint.QueryEndpointRequest) *queryEndpointRequest {
	return &queryEndpointRequest{req}
}

type queryEndpointRequest struct {
	*endpoint.QueryEndpointRequest
}

func (r *queryEndpointRequest) FindOptions() *options.FindOptions {
	pageSize := int64(r.PageSize)
	skip := int64(r.PageSize) * int64(r.PageNumber-1)

	opt := &options.FindOptions{
		Sort:  bson.D{{Key: "create_at", Value: -1}},
		Limit: &pageSize,
		Skip:  &skip,
	}

	return opt
}

func (r *queryEndpointRequest) FindFilter() bson.M {
	filter := bson.M{}

	if r.ServiceID != "" {
		filter["service_id"] = r.ServiceID
	}
	if r.Method != "" {
		filter["method"] = r.Method
	}
	if r.Resource != "" {
		filter["resource"] = r.Resource
	}
	if r.Path != "" {
		filter["path"] = bson.M{"$regex": r.Path, "$options": "im"}
	}
	if r.FunctionName != "" {
		filter["function_name"] = r.FunctionName
	}

	return filter
}
