package service

type JWTService interface {
	GenerateToken(userID string, role string) (string, error)
}
