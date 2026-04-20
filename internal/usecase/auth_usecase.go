package usecase

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	crand "crypto/rand"

	"github.com/google/uuid"
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
			ID:       uuid.New(),
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
			ID:         uuid.New(),
			Username:   input.Username,
			Phone:      input.Phone,
			Email:      input.Email,
			BusName:    input.BusName,
			BusNumbers: input.BusNumbers,
			BusType:    input.BusType,
			SeatType:   input.SeatType,
		} //passenger driver travel only one number
		if err := uc.driverRepo.Create(d); err != nil {
			log.Printf(fmt.Sprintf(`[%s] Register driver create error: %v`, usecaseLogPrefix, err))
			return nil, err
		}
		log.Printf(fmt.Sprintf(`[%s] Register driver created: %s`, usecaseLogPrefix, d.ID))
		res = d

	} else if input.Role == "travel" {
		busesNumbers := make([]string, 0, len(input.BusesNumbers))
		for _, b := range input.BusesNumbers {
			if b != "" {
				busesNumbers = append(busesNumbers, b)
			}
		}
		t := entity.Travels{
			ID:    uuid.New(),
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

func (uc *AuthUsecase) Verify(phone, code string) (res interface{}, err error) {
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

	p, _ := uc.passengerRepo.FindByPhone(phone)
	if p != nil {
		if err := uc.passengerRepo.Verify(phone); err != nil {
			log.Printf(fmt.Sprintf(`[%s] Verify passenger verify error: %v`, usecaseLogPrefix, err))
			return nil, err
		}
		token, err := uc.jwtService.GenerateToken(p.ID.String(), "passenger")
		if err != nil {
			log.Printf(fmt.Sprintf(`[%s] Verify passenger token generate error: %v`, usecaseLogPrefix, err))
			return nil, err
		}
		log.Printf(fmt.Sprintf(`[%s] Verify passenger token generated for phone: %s`, usecaseLogPrefix, phone))
		return token, nil
	}

	d, _ := uc.driverRepo.FindByPhone(phone)
	if d != nil {
		if err := uc.driverRepo.Verify(phone); err != nil {
			log.Printf(fmt.Sprintf(`[%s] Verify driver verify error: %v`, usecaseLogPrefix, err))
			return nil, err
		}
		token, err := uc.jwtService.GenerateToken(d.ID.String(), "driver")
		if err != nil {
			log.Printf(fmt.Sprintf(`[%s] Verify driver token generate error: %v`, usecaseLogPrefix, err))
			return nil, err
		}
		log.Printf(fmt.Sprintf(`[%s] Verify driver token generated for phone: %s`, usecaseLogPrefix, phone))
		return token, nil
	}

	tr, _ := uc.travelsRepo.FindByPhone(phone)
	if tr != nil {
		if err := uc.travelsRepo.Verify(phone); err != nil {
			log.Printf(fmt.Sprintf(`[%s] Verify travel verify error: %v`, usecaseLogPrefix, err))
			return nil, err
		}
		token, err := uc.jwtService.GenerateToken(tr.ID.String(), "travel")
		if err != nil {
			log.Printf(fmt.Sprintf(`[%s] Verify travel token generate error: %v`, usecaseLogPrefix, err))
			return nil, err
		}
		log.Printf(fmt.Sprintf(`[%s] Verify travel token generated for phone: %s`, usecaseLogPrefix, phone))
		return token, nil
	}

	if phone == os.Getenv("ADMIN_PHONE") {

		token, err := uc.jwtService.GenerateToken("S001", "admin")
		if err != nil {
			log.Printf(fmt.Sprintf(`[%s] Verify admin token generate error: %v`, usecaseLogPrefix, err))
			return nil, err
		}
		return token, nil
	}

	log.Printf(fmt.Sprintf(`[%s] Verify user not found for phone: %s`, usecaseLogPrefix, phone))
	return nil, errors.New("user not found,Please register first")
}

// have to think needed or not
func (uc *AuthUsecase) Login(phone string) (res interface{}, err error) {
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
