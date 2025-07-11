package validator

import (
	"net/mail"

	tvoerrors "main/tools/pkg/tvo_errors"
)

// ValidEmail function validates the given email address.
// It takes an email string as input and returns an error if the email is invalid.
func ValidEmail(email string) error {
	// Parse the email address using the mail.ParseAddress function
	_, err := mail.ParseAddress(email)
	if err != nil {
		// If parsing fails, return an error with a message indicating the invalid email address
		return tvoerrors.Wrap(email, tvoerrors.ErrInvalidEmail)
	}
	// If parsing succeeds, return nil indicating no error
	return nil
}
