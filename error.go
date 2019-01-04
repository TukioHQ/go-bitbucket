package bitbucket

import (
	"errors"

	"github.com/mitchellh/mapstructure"
)

// Error wrapper object for HTTP request error
type Error struct {
	Status int
	Body   struct {
		Message string `json:"message"`
	} `json:"error"`
}

// Error returns the message body of HTTP request error
func (e Error) Error() string {
	return e.Body.Message
}

// BitbucketError wrapper object
type BitbucketError struct {
	Message string
	Fields  map[string][]string
}

// DecodeError extracts errors found for request to Bitbucket
func DecodeError(e map[string]interface{}) error {
	var bitbucketError BitbucketError
	err := mapstructure.Decode(e["error"], &bitbucketError)
	if err != nil {
		return err
	}

	return errors.New(bitbucketError.Message)
}
