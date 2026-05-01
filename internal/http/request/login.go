package request

type LoginInput struct {
	Phone string `json:"phone"`
	Role  string `json:"role"`
}
