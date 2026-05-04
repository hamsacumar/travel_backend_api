package usecase

import (
	"errors"

	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	crand "crypto/rand"

	"github.com/hamsacumar/travel_backend_api/internal/domain/entity"
	"github.com/hamsacumar/travel_backend_api/internal/domain/repository"
	"github.com/hamsacumar/travel_backend_api/internal/domain/service"
	"github.com/hamsacumar/travel_backend_api/internal/http/request"
	smsotp "github.com/hamsacumar/travel_backend_api/internal/infrastructure/service/otp"
)

const usecaseLogPrefix = `travels-api.internal.usecase.auth_usecase`

type AuthUsecase struct {
	passengerRepo repository.PassengerRepository
	driverRepo    repository.DriverRepository
	travelsRepo   repository.TravelsRepository
	otpRepo       repository.OTPRepository
	jwtService    service.JWTService
}

// NewAuthUsecase creates and returns a new AuthUsecase.
func NewAuthUsecase(
	passengerRepo repository.PassengerRepository,
	driverRepo repository.DriverRepository,
	travelsRepo repository.TravelsRepository,
	otpRepo repository.OTPRepository,
	jwtService service.JWTService,
) *AuthUsecase {
	return &AuthUsecase{
		passengerRepo: passengerRepo,
		driverRepo:    driverRepo,
		travelsRepo:   travelsRepo,
		otpRepo:       otpRepo,
		jwtService:    jwtService,
	}
}

func (uc *AuthUsecase) Register(input request.SignUpInput) (res interface{}, err error) {
	log.Printf(fmt.Sprintf(`[%s] Register started for phone: %s role: %s`, usecaseLogPrefix, input.Phone, input.Role))

	if input.Role == "passenger" {
		p := entity.Passenger{
			Username: input.Username,
			Phone:    input.Phone,
			Email:    input.Email,
		} //passenger driver travel only one number
		if err := uc.passengerRepo.Create(p); err != nil {
			log.Printf(fmt.Sprintf(`[%s] Register passenger create error: %v`, usecaseLogPrefix, err))
			return nil, err
		}
		log.Printf(fmt.Sprintf(`[%s] Register passenger created: %s`, usecaseLogPrefix, p.ID))
		res = p

	} else if input.Role == "driver" {
		d := entity.Driver{
			Username:   input.Username,
			Phone:      input.Phone,
			Email:      input.Email,
			BusName:    input.BusName,
			BusNumbers: input.BusNumbers,
			BusType:    input.BusType,
			SeatCount:  input.SeatCount,
			SeatType:   input.SeatType,
		} //passenger driver travel only one number
		if err := uc.driverRepo.Create(d); err != nil {
			log.Printf(fmt.Sprintf(`[%s] Register driver create error: %v`, usecaseLogPrefix, err))
			return nil, err
		}
		log.Printf(fmt.Sprintf(`[%s] Register driver created: %s`, usecaseLogPrefix, d.ID))
		res = d

	} else if input.Role == "travel" {
		t := entity.Travels{
			Name:  input.Username,
			Phone: input.Phone,
			Email: input.Email,
		} //passenger driver travel only one number
		if err := uc.travelsRepo.Create(t); err != nil {
			log.Printf(fmt.Sprintf(`[%s] Register travel create error: %v`, usecaseLogPrefix, err))
			return nil, err
		}
		log.Printf(fmt.Sprintf(`[%s] Register travel created: %s`, usecaseLogPrefix, t.ID))
		res = t

	} else {
		return nil, errors.New("invalid role")
	}

	return res, nil
}

func (uc *AuthUsecase) Verify(phone, code, role string) (res interface{}, err error) {
	log.Printf(fmt.Sprintf(`[%s] Verify started for phone: %s`, usecaseLogPrefix, phone))

	otp, err := uc.otpRepo.Find(phone, code)
	if err != nil || otp == nil {
		log.Printf(fmt.Sprintf(`[%s] Verify invalid otp for phone: %s`, usecaseLogPrefix, phone))
		return nil, errors.New("invalid otp")
	}

	if time.Now().After(otp.ExpiresAt) {
		log.Printf(fmt.Sprintf(`[%s] Verify otp expired for phone: %s`, usecaseLogPrefix, phone))
		return nil, errors.New("otp expired")
	}

	switch role {
	case "passenger":
		p, _ := uc.passengerRepo.FindByPhone(phone)
		if p == nil {
			log.Printf(fmt.Sprintf(`[%s] Verify passenger not found for phone: %s`, usecaseLogPrefix, phone))
			return nil, errors.New("passenger not found, please register first")
		}
		if err := uc.passengerRepo.Verify(phone); err != nil {
			log.Printf(fmt.Sprintf(`[%s] Verify passenger verify error: %v`, usecaseLogPrefix, err))
			return nil, err
		}
		token, err := uc.jwtService.GenerateToken(p.ID, phone, "passenger")
		if err != nil {
			log.Printf(fmt.Sprintf(`[%s] Verify passenger token generate error: %v`, usecaseLogPrefix, err))
			return nil, err
		}
		log.Printf(fmt.Sprintf(`[%s] Verify passenger token generated for phone: %s`, usecaseLogPrefix, phone))
		return token, nil
	case "driver":
		d, _ := uc.driverRepo.FindByPhone(phone)
		if d == nil {
			log.Printf(fmt.Sprintf(`[%s] Verify driver not found for phone: %s`, usecaseLogPrefix, phone))
			return nil, errors.New("driver not found, please register first")
		}
		if err := uc.driverRepo.Verify(phone); err != nil {
			log.Printf(fmt.Sprintf(`[%s] Verify driver verify error: %v`, usecaseLogPrefix, err))
			return nil, err
		}
		token, err := uc.jwtService.GenerateToken(d.ID, phone, "driver")
		if err != nil {
			log.Printf(fmt.Sprintf(`[%s] Verify driver token generate error: %v`, usecaseLogPrefix, err))
			return nil, err
		}
		log.Printf(fmt.Sprintf(`[%s] Verify driver token generated for phone: %s`, usecaseLogPrefix, phone))
		return token, nil
	case "travel":
		tr, _ := uc.travelsRepo.FindByPhone(phone)
		if tr == nil {
			log.Printf(fmt.Sprintf(`[%s] Verify travel not found for phone: %s`, usecaseLogPrefix, phone))
			return nil, errors.New("travel not found, please register first")
		}
		if err := uc.travelsRepo.Verify(phone); err != nil {
			log.Printf(fmt.Sprintf(`[%s] Verify travel verify error: %v`, usecaseLogPrefix, err))
			return nil, err
		}
		token, err := uc.jwtService.GenerateToken(tr.ID, phone, "travel")
		if err != nil {
			log.Printf(fmt.Sprintf(`[%s] Verify travel token generate error: %v`, usecaseLogPrefix, err))
			return nil, err
		}
		log.Printf(fmt.Sprintf(`[%s] Verify travel token generated for phone: %s`, usecaseLogPrefix, phone))
		return token, nil
	case "admin":
		if phone != os.Getenv("ADMIN_PHONE") {
			log.Printf(fmt.Sprintf(`[%s] Verify admin phone mismatch for phone: %s`, usecaseLogPrefix, phone))
			return nil, errors.New("invalid admin credentials")
		}
		token, err := uc.jwtService.GenerateToken("000000", phone, "admin")
		if err != nil {
			log.Printf(fmt.Sprintf(`[%s] Verify admin token generate error: %v`, usecaseLogPrefix, err))
			return nil, err
		}
		return token, nil
	default:
		return nil, errors.New("invalid role")
	}
}

// have to think needed or not
func (uc *AuthUsecase) Login(phone, role string) (res interface{}, err error) {
	log.Printf(fmt.Sprintf(`[%s] Login started for phone: %s`, usecaseLogPrefix, phone))

	// Delegate to shared SendOTP usecase to generate, store, and send the OTP
	return uc.SendOTP(phone)
}

func (uc *AuthUsecase) SendOTP(phone string) (res interface{}, err error) {
	log.Printf(fmt.Sprintf(`[%s] Otp sended to phone: %s`, usecaseLogPrefix, phone))

	// generate random 5-digit OTP
	code, _ := random5Digit()
	otp := entity.OTP{
		Phone:     phone,
		Code:      code,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	if err := uc.otpRepo.Save(otp); err != nil {
		log.Printf(fmt.Sprintf(`[%s] Login otp save error: %v`, usecaseLogPrefix, err))
		return nil, err
	}
	log.Printf(fmt.Sprintf(`[%s] Login OTP saved successfully in db: %s`, usecaseLogPrefix, phone))
	if err := smsotp.Send(phone, code); err != nil {
		log.Printf(fmt.Sprintf(`[%s] Login OTP SMS send error: %v`, usecaseLogPrefix, err))
		return nil, err
	}
	log.Printf(fmt.Sprintf(`[%s] Login OTP send successfully for phone: %s`, usecaseLogPrefix, phone))
	return "OTP sent successfully", nil
}

// random5Digit generates a cryptographically secure 5-digit code as string (10000-99999)
func random5Digit() (string, error) {
	// generate number in range [10000, 99999]
	const min = 10000
	const max = 99999
	// generate 2 bytes and mod the range; simple and sufficient for OTP
	var b [2]byte
	if _, err := crand.Read(b[:]); err != nil {
		// fallback to time-based if crypto fails
		n := time.Now().UnixNano()%90000 + 10000
		return fmt.Sprintf("%05d", n), nil
	}
	n := int(b[0])<<8 | int(b[1])
	n = n%((max-min)+1) + min
	return fmt.Sprintf("%05d", n), nil
}

// sixDigitUserID derives a stable 6-digit user id from phone using HMAC-SHA256(secret, phone)
// and formatting modulo 1e6 as zero-padded string. This avoids storing the code while keeping it deterministic.
func generateSixDigitID() string {
	for {
		n, _ := rand.Int(rand.Reader, big.NewInt(1000000))
		id := fmt.Sprintf("%06d", n.Int64())

		return id
	}
}
