package domain

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/infraboard/keyauth/common/password"
	"github.com/infraboard/keyauth/pkg/user"
	"github.com/infraboard/mcube/exception"
)

// NewDefaultSecuritySetting todo
func NewDefaultSecuritySetting() *SecuritySetting {
	return &SecuritySetting{
		PasswordSecurity: NewDefaulPasswordSecurity(),
		LoginSecurity:    NewDefaultLoginSecurity(),
	}
}

// SecuritySetting 安全设置
type SecuritySetting struct {
	PasswordSecurity *PasswordSecurity `bson:"password_security" json:"password_security"` // 密码安全
	LoginSecurity    *LoginSecurity    `bson:"login_security" json:"login_security"`       // 登录安全
}

// GetPasswordRepeateLimite todo
func (s *SecuritySetting) GetPasswordRepeateLimite() uint {
	if s.PasswordSecurity == nil {
		return 0
	}
	return s.PasswordSecurity.RepeateLimite
}

// Patch todo
func (s *SecuritySetting) Patch(data *SecuritySetting) {
	patchData, _ := json.Marshal(data)
	json.Unmarshal(patchData, s)
}

// NewDefaulPasswordSecurity todo
func NewDefaulPasswordSecurity() *PasswordSecurity {
	return &PasswordSecurity{
		Length:                  8,
		IncludeNumber:           true,
		IncludeLowerLetter:      true,
		IncludeUpperLetter:      false,
		IncludeSymbols:          false,
		RepeateLimite:           1,
		PasswrodExpiredDays:     90,
		BeforeExpiredRemindDays: 10,
	}
}

// PasswordSecurity 密码安全设置
type PasswordSecurity struct {
	Length                  int  `bson:"length" json:"length" validate:"required,min=8,max=64"`                                          // 密码长度
	IncludeNumber           bool `bson:"include_number" json:"include_number"`                                                           // 包含数字
	IncludeLowerLetter      bool `bson:"include_lower_letter" json:"include_lower_letter"`                                               // 包含小写字母
	IncludeUpperLetter      bool `bson:"include_upper_letter" json:"include_upper_letter"`                                               // 包含大写字母
	IncludeSymbols          bool `bson:"include_symbols" json:"include_symbols"`                                                         // 包含特殊字符
	RepeateLimite           uint `bson:"repeate_limite" json:"repeate_limite" validate:"required,min=1,max=24"`                          // 重复限制
	PasswrodExpiredDays     uint `bson:"password_expired_days" json:"password_expired_days" validate:"required,min=0,max=365"`           // 密码过期时间, 密码过期后要求用户重置密码
	BeforeExpiredRemindDays uint `bson:"before_expired_remind_days" json:"before_expired_remind_days" validate:"required,min=0,max=365"` // 密码过期前多少天开始提醒
}

// Validate 校验对象合法性
func (p *PasswordSecurity) Validate() error {
	return validate.Struct(p)
}

// IsPasswordExpired todo
func (p *PasswordSecurity) IsPasswordExpired(pass *user.Password) error {
	if p.PasswrodExpiredDays == 0 {
		return nil
	}

	delta := p.expiredDelta(pass.UpdateAt.T())
	if delta > 0 {
		return exception.NewPasswordExired("password expired %d days", delta)
	}

	return nil
}

// SetPasswordNeedReset todo
func (p *PasswordSecurity) SetPasswordNeedReset(pass *user.Password) {
	// 密码用不过期, 不需要重置
	if p.PasswrodExpiredDays == 0 {
		return
	}

	// 计算密码是否过期
	delta := p.expiredDelta(pass.UpdateAt.T())
	if delta > 0 {
		pass.SetExpired()
		return
	}

	// 计算是否即将过期, 需要用户重置
	if -delta < int(p.BeforeExpiredRemindDays) {
		pass.SetNeedReset("密码%d天后过期, 请重置密码", -delta)
	}
}

func (p *PasswordSecurity) expiredDelta(updateAt time.Time) int {
	updateBefore := uint(time.Now().Sub(updateAt).Hours() / 24)
	return int(updateBefore) - int(p.PasswrodExpiredDays)
}

// Check todo
func (p *PasswordSecurity) Check(pass string) error {
	v := password.NewValidater(pass)

	if ok := v.LengthOK(p.Length); !ok {
		return fmt.Errorf("password length less than %d", p.Length)
	}
	if p.IncludeNumber {
		if ok := v.IncludeNumbers(); !ok {
			return fmt.Errorf("must include numbers")
		}
	}
	if p.IncludeLowerLetter {
		if ok := v.IncludeLowercaseLetters(); !ok {
			return fmt.Errorf("must include lower letter")
		}
	}
	if p.IncludeUpperLetter {
		if ok := v.IncludeUppercaseLetters(); !ok {
			return fmt.Errorf("must include upper letter")
		}
	}
	if p.IncludeSymbols {
		if ok := v.IncludeSymbols(); !ok {
			return fmt.Errorf("must include symbols")
		}
	}

	return nil
}

// NewDefaultLoginSecurity todo
func NewDefaultLoginSecurity() *LoginSecurity {
	return &LoginSecurity{
		ExceptionLock: true,
		ExceptionLockConfig: &ExceptionLockConfig{
			OtherPlaceLogin: true,
			NotLoginDays:    30,
		},
		RetryLock: true,
		RetryLockConfig: &RetryLockConig{
			RetryLimite:  5,
			LockedMinite: 30,
		},
		IPLimite: false,
		IPLimiteConfig: &IPLimiteConfig{
			IP: []string{},
		},
	}
}

// LoginSecurity 登录安全
type LoginSecurity struct {
	ExceptionLock       bool                 `bson:"exception_lock" json:"exception_lock"`               // 异常登录锁
	ExceptionLockConfig *ExceptionLockConfig `bson:"exception_lock_config" json:"exception_lock_config"` // 异常配置
	RetryLock           bool                 `bson:"retry_lock" json:"retry_lock"`                       // 重试锁
	RetryLockConfig     *RetryLockConig      `bson:"retry_lock_config" json:"retry_lock_config"`         // 重试锁配置
	IPLimite            bool                 `bson:"ip_limite" json:"ip_limite"`                         // IP限制
	IPLimiteConfig      *IPLimiteConfig      `bson:"ip_limite_config" json:"ip_limite_config"`           // IP限制配置
}

// ExceptionLockConfig todo
type ExceptionLockConfig struct {
	OtherPlaceLogin bool `bson:"other_place_login" json:"other_place_login"` // 异地登录
	NotLoginDays    uint `bson:"not_login_days" json:"not_login_days"`       // 未登录天数,
}

// IPLimiteConfig todo
type IPLimiteConfig struct {
	Type string   `bson:"type" json:"type"` // 黑名单还是白名单
	IP   []string `bson:"ip" json:"ip"`     // ip列表
}

// RetryLockConig 重试锁配置
type RetryLockConig struct {
	RetryLimite  uint `bson:"retry_limite" json:"retry_limite"`   // 重试限制
	LockedMinite uint `bson:"locked_minite" json:"locked_minite"` // 锁定时长
}

// LockedMiniteDuration todo
func (c *RetryLockConig) LockedMiniteDuration() time.Duration {
	return time.Duration(c.LockedMinite) * time.Minute
}
