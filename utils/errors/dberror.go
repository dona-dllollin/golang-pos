package errorUtils

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

func MapDbError(err error) error {
	if err == nil {
		return nil
	}

	// pgx
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return mapCode(pgErr.Code)
	}

	//     // lib/pq
	//     var pqErr *pq.Error
	//     if errors.As(err, &pqErr) {
	//         return mapCode(string(pqErr.Code))
	//     }

	return ErrInternal
}

func mapCode(code string) error {
	switch code {
	case "23505": // unique_violation
		return ErrConflict
	case "23502": // not_null_violation
		return ErrBadRequest
	case "23503": // foreign_key_violation
		return ErrBadRequest
	case "23514": // check violation
		return ErrBadRequest
	case "22P02": // invalid_text_representation
		return ErrBadRequest
	}
	return ErrInternal
}
