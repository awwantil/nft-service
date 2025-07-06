package validator

import (
	"github.com/dongri/phonenumber"

	tvoerrors "main/tools/pkg/tvo_errors"
)

// ValidPhone function validates the given phone number.
// It takes a phone number string as input and returns an error if the phone number is invalid.
func ValidPhone(phone string) error {
	// Get ISO 3166 country information based on the phone number
	isoPhone := phonenumber.GetISO3166ByNumber(phone, true)

	// If the country name is empty, it indicates an invalid phone number
	if isoPhone.CountryName == "" {
		// Return an error indicating the invalid phone number
		return tvoerrors.Wrap(phone, tvoerrors.ErrInvalidPhone)
	}

	// If the country name is not empty, the phone number is considered valid
	return nil
}
