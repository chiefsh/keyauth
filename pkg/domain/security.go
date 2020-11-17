package domain

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

// NewDefaulPasswordSecurity todo
func NewDefaulPasswordSecurity() *PasswordSecurity {
	return &PasswordSecurity{
		Length:             8,
		IncludeNumber:      true,
		IncludeLowerLetter: true,
		IncludeUpperLetter: false,
		IncludeSymbols:     false,
		RepeateLimite:      1,
	}
}

// PasswordSecurity 密码安全设置
type PasswordSecurity struct {
	Length             int  `bson:"length" json:"length"`                             // 密码长度
	IncludeNumber      bool `bson:"include_number" json:"include_number"`             // 包含数字
	IncludeLowerLetter bool `bson:"include_lower_letter" json:"include_lower_letter"` // 包含小写字母
	IncludeUpperLetter bool `bson:"include_upper_letter" json:"include_upper_letter"` // 包含大写字母
	IncludeSymbols     bool `bson:"include_symbols" json:"include_symbols"`           // 包含特殊字符
	RepeateLimite      uint `bson:"repeate_limite" json:"repeate_limite"`             // 重复限制
}

// NewDefaultLoginSecurity todo
func NewDefaultLoginSecurity() *LoginSecurity {
	return &LoginSecurity{
		ExceptionLock: true,
		ExceptionLockConfig: &ExceptionLockConfig{
			PasswrodExpiredDays: 0,
			NotLoginDays:        30,
		},
		RetryLock: true,
		RetryLockConfig: &RetryLockConig{
			RetryLimite:  5,
			LockedMinite: 30,
		},
		IPLimite:       false,
		IPLimiteConfig: &IPLimiteConfig{},
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
	PasswrodExpiredDays uint `bson:"password_expired_days" json:"password_expired_days"` // 密码过期时间, 密码过期后要求用户重置密码
	NotLoginDays        uint `bson:"not_login_days" json:"not_login_days"`               // 未登录天数,
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
