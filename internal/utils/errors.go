package utils

import (
	"errors"

	"github.com/wasify-io/wasify-go/logging"
)
import pkg_errors "github.com/pkg/errors"

func Aggregate(message string, elements ...error) error {
	return pkg_errors.WithStack(
		pkg_errors.Wrap(errors.Join(elements...), message),
	)
}

func Log(logger logging.Logger, err error) error {
	if err != nil {
		logger.Error(err.Error())
	}

	return err
}
