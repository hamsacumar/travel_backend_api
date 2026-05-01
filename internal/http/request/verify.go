package request

type VerifyInput struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
	Role  string `json:"role"`
}
