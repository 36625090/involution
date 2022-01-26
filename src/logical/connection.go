package logical

type Connection struct {
	RemoteAddr string `json:"remote_addr" validate:"required"`
	UserAgent  string `json:"user_agent" validate:"required"`
}
