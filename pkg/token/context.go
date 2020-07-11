package token

// NewSession todo
func NewSession() *Session {
	return &Session{}
}

// Session 请求上下文信息
type Session struct {
	tk *Token
}

// WithToken 携带token
func (s *Session) WithToken(tk *Token) {
	s.tk = tk
}

// GetToken 获取token
func (s *Session) GetToken() *Token {
	return s.tk
}

// UserID todo
func (s *Session) UserID() string {
	if s.tk == nil {
		return "Nil"
	}

	return s.tk.UserID
}