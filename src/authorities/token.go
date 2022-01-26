package authorities

type TokenHandler interface {
	GenerateToken(auth *Authorized) (string, error)
	ParseToken(token string) (*Authorized, error)
}