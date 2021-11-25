package utils

import (
	"path"

	"github.com/jackc/pgconn"

	"github.com/pkg/errors"
)

func CheckURL(u string) error {
	if len(u) == 0 {
		return errors.New("empty url")
	}
	if u[0] != '/' {
		return errors.New(`'/' should be the first char.`)
	}
	if u[len(u)-1] == '/' {
		return errors.New(`'/' should not be the last char.`)
	}
	if u != path.Clean(u) {
		return errors.New("Invalid path.")
	}
	return nil
}

func IsUniqueError(err error, constraint string) bool {
	pgError, ok := err.(*pgconn.PgError)
	return ok && pgError.Code == "23505" && pgError.ConstraintName == constraint
}
