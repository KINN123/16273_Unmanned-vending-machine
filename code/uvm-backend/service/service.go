package service

import (
	"errors"
	"fmt"
)

var ErrorNameInvalid = errors.New("wrong name")

func test() error {
	err := errors.New("test")
	if err != nil {
		return fmt.Errorf("Aadfsadf: %w", ErrorNameInvalid)
	}
	if errors.Is(err, ErrorNameInvalid) {

	}
	return err
}
