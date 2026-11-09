package util

import (
	"context"
	"fmt"
	"net/http"
	"net/mail"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var validate = validator.New()

type requestIDKey string

const (
	TimeFormatYYYYMMDD   = "2006-01-02"
	TimeFormatDDMMYYYY   = "02-01-2006"
	TimeFormatWithTime   = "2006-01-02 15:04:05"
	TimeFormatYYYYMMDDTZ = "2006-01-02T15:04:05Z"

	RequestIDKey requestIDKey = "X-NSR-Request-ID"
)

func GetRequestID(ctx context.Context) string {
	id := ctx.Value(RequestIDKey)
	if str, ok := id.(string); ok {
		return str
	}
	return ""
}

func RequestIDLog(ctx context.Context) string {
	requestID := GetRequestID(ctx)
	return fmt.Sprintf("RequestID : %v : ", requestID)
}

func GetSQLFormatIDs(ids []string) string {
	var result string
	for _, id := range ids {
		if id == "" {
			continue
		}
		result += fmt.Sprintf("'%s',", id)
	}
	if result != "" {
		result = result[:len(result)-1]
	}
	return result
}

func ValidateRequest(ctx context.Context, request interface{}) (int, error) {
	var errMessage string
	if err := validate.Struct(request); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errMessage += fmt.Sprintf("invalid %s; ", err.Field())
		}
	}

	if errMessage != "" {
		return http.StatusBadRequest, fmt.Errorf("error : %v", errMessage)
	}
	return http.StatusOK, nil
}

// ValidateIntUUIDValue set value isInt from given input string.
// If true then its value is int, else its value is UUID.
// It is only valid for int from 1 to 1000.
// Returns error if value is neither int nor UUID.
func ValidateIntUUIDValue(isInt *bool, input string) (int, error) {
	err := uuid.Validate(input)
	if err == nil {
		*isInt = false
		return http.StatusOK, nil
	}
	id, _ := strconv.Atoi(input)
	if id > 0 && id <= 1000 {
		*isInt = true
		return http.StatusOK, nil
	}
	return http.StatusBadRequest, fmt.Errorf("invalid request : invalid value %v", input)
}

// ValidateUUIDValue set value isUUID from given input string.
// If true then its value is init, else its value is string
// with length more than 1 character.
func ValidateUUIDAplhanumericValue(isUUID *bool, input string) (int, error) {
	err := uuid.Validate(input)
	if err == nil {
		*isUUID = true
		return http.StatusOK, nil
	}

	if IsAlphanumeric40(input) {
		*isUUID = false
		return http.StatusOK, nil
	}

	return http.StatusBadRequest, fmt.Errorf("invalid request : invalid value %v", input)
}

// IsAlphanumeric40 check string min 1 char and max 40 char.
// Set Max so that string more than 40 char will be rejected.
func IsAlphanumeric40(input string) bool {
	var re = regexp.MustCompile(`^[A-Za-z0-9]{1,40}$`)
	return re.Match([]byte(input))
}

func IsVersionLower(userDeviceVersion string, minDeviceVersion string) (bool, error) {

	if userDeviceVersion == "" && minDeviceVersion != "" {
		return true, nil
	}

	if minDeviceVersion == "" {
		return false, nil
	}

	versionReq, err := ParseVersion(userDeviceVersion)
	if err != nil {
		return false, err
	}
	versionCurrent, err := ParseVersion(minDeviceVersion)
	if err != nil {
		return false, err
	}

	if versionReq.Major < versionCurrent.Major {
		return true, nil
	} else if versionReq.Major > versionCurrent.Major {
		return false, nil
	}

	if versionReq.Minor < versionCurrent.Minor {
		return true, nil
	} else if versionReq.Minor > versionCurrent.Minor {
		return false, nil
	}

	if versionReq.Patch < versionCurrent.Patch {
		return true, nil
	}

	return false, nil
}

type Version struct {
	Major int
	Minor int
	Patch int
}

// ParseVersion parses a version string (e.g., "1.2.3") into a Version struct.
func ParseVersion(versionStr string) (Version, error) {
	versionParts := strings.Split(versionStr, ".")
	if len(versionParts) != 3 {
		return Version{}, fmt.Errorf("Invalid version format: %s", versionStr)
	}

	major, err := strconv.Atoi(versionParts[0])
	if err != nil {
		return Version{}, err
	}

	minor, err := strconv.Atoi(versionParts[1])
	if err != nil {
		return Version{}, err
	}

	patch, err := strconv.Atoi(versionParts[2])
	if err != nil {
		return Version{}, err
	}

	return Version{
		Major: major,
		Minor: minor,
		Patch: patch,
	}, nil
}

func GetSeatsName(input string) []string {
	output := strings.Split(input, ",")
	for i := range output {
		output[i] = strings.TrimSpace(output[i])
	}
	return output
}

func isEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
