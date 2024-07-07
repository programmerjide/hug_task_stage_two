package config

import (
	"errors"
	"strings"
)

type ResponseStatus struct {
	Code        string
	Description string
}

var (
	API_SUCCESS_STATUS = ResponseStatus{"success", "Request processed successfully"}
	API_FAIL_STATUS    = ResponseStatus{"fail", "Failed processing request"}
	API_ERROR_STATUS   = ResponseStatus{"error", "Error processing request"}
)

func GetResponseStatusByCode(value string) (ResponseStatus, error) {
	if value == "" {
		return ResponseStatus{}, errors.New("input value is empty")
	}

	statuses := []ResponseStatus{API_SUCCESS_STATUS, API_FAIL_STATUS, API_ERROR_STATUS}

	for _, status := range statuses {
		if strings.EqualFold(status.Code, value) {
			return status, nil
		}
	}

	return ResponseStatus{}, errors.New("invalid response code")
}

type ResponseCode struct {
	Code        string
	Description string
}

var (
	SUCCESS         = ResponseCode{"00", "Approved or completed successfully"}
	FAILED          = ResponseCode{"99", "Failed to process request"}
	UNKNOWN         = ResponseCode{"01", "Status unknown"}
	INVALID_PAYLOAD = ResponseCode{"02", "Invalid payload passed!"}
)

func GetResponseCodeByCode(value string) (ResponseCode, error) {
	if value == "" {
		return ResponseCode{}, errors.New("input value is empty")
	}

	codes := []ResponseCode{
		SUCCESS, FAILED, UNKNOWN, INVALID_PAYLOAD,
	}

	for _, code := range codes {
		if strings.EqualFold(code.Code, value) {
			return code, nil
		}
	}

	return ResponseCode{}, errors.New("invalid response code")
}
