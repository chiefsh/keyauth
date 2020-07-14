package department

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/infraboard/mcube/exception"
	"github.com/infraboard/mcube/http/request"
	"github.com/infraboard/mcube/types/ftime"
	"github.com/rs/xid"

	"github.com/infraboard/keyauth/pkg/token"
)

// use a single instance of Validate, it caches struct info
var (
	validate = validator.New()
)

// NewDepartment 新建实例
func NewDepartment(req *CreateDepartmentRequest, d Service) (*Department, error) {
	if err := req.Validate(); err != nil {
		return nil, exception.NewBadRequest(err.Error())
	}

	tk := req.GetToken()

	ins := &Department{
		ID:                      xid.New().String(),
		CreateAt:                ftime.Now(),
		UpdateAt:                ftime.Now(),
		CreaterID:               tk.UserID,
		DomainID:                tk.DomainID,
		Grade:                   1,
		CreateDepartmentRequest: req,
	}

	if req.ParentID != "" {
		pd, err := d.DescribeDepartment(NewDescriptDepartmentRequestWithID(req.ParentID))
		if err != nil {
			return nil, err
		}
		ins.ParentPath = pd.Path()
		ins.Grade = len(strings.Split(pd.Path(), "."))
	}

	if req.ManagerID == "" {
		req.ManagerID = tk.UserID
	}

	return ins, nil
}

// NewDefaultDepartment todo
func NewDefaultDepartment() *Department {
	return &Department{
		CreateDepartmentRequest: NewCreateDepartmentRequest(),
	}
}

// Department user's department
type Department struct {
	ID                       string     `bson:"_id" json:"id"`                          // 部门ID
	ParentPath               string     `bson:"parent_path" json:"parent_path"`         // 路径
	Number                   uint64     `bson:"number" json:"number,omitempty"`         // 部门编号
	CreateAt                 ftime.Time `bson:"create_at" json:"create_at,omitempty"`   // 部门创建时间
	UpdateAt                 ftime.Time `bson:"update_at" json:"update_at,omitempty"`   // 更新时间
	CreaterID                string     `bson:"creater_id" json:"creater_id,omitempty"` // 创建人
	DomainID                 string     `bson:"domain_id" json:"domain_id,omitempty"`   // 部门所属域
	Grade                    int        `bson:"grade" json:"grade,omitempty"`           // 第几级部门, 由层数决定
	*CreateDepartmentRequest `bson:",inline"`
}

// Path 具体路径
func (d *Department) Path() string {
	return fmt.Sprintf("%s.%d", d.ParentPath, d.Number)
}

// NewCreateDepartmentRequest todo
func NewCreateDepartmentRequest() *CreateDepartmentRequest {
	return &CreateDepartmentRequest{
		Session: token.NewSession(),
	}
}

// CreateDepartmentRequest 创建部门请求
type CreateDepartmentRequest struct {
	*token.Session `bson:"-" json:"-"`
	Name           string `bson:"name" json:"name,omitempty" validate:"required,lte=60"`     // 部门名称
	ParentID       string `bson:"parent_id" json:"parent_id,omitempty" validate:"lte=200"`   // 上级部门ID
	ManagerID      string `bson:"manager_id" json:"manager_id,omitempty" validate:"lte=200"` // 部门管理者ID
}

// Validate 校验参数的合法性
func (req *CreateDepartmentRequest) Validate() error {

	if req.Session == nil {
		return fmt.Errorf("session required")
	}

	tk := req.GetToken()
	if tk == nil {
		return fmt.Errorf("session token required")
	}
	if tk.DomainID == "" {
		return fmt.Errorf("user must create domain first")
	}

	return validate.Struct(req)
}

// NewDepartmentSet 实例化
func NewDepartmentSet(req *request.PageRequest) *Set {
	return &Set{
		PageRequest: req,
		Items:       []*Department{},
	}
}

// Set 集合
type Set struct {
	*request.PageRequest

	Total int64         `json:"total"`
	Items []*Department `json:"items"`
}

// Add 添加应用
func (s *Set) Add(item *Department) {
	s.Items = append(s.Items, item)
}
