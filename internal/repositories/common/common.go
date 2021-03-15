package common

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func CommitTransaction(_ context.Context, tx *sqlx.Tx, err error, method string) error {
	if r := recover(); r != nil {
		err = fmt.Errorf("panic in %s", method)
	}

	if err == nil {
		commitErr := tx.Commit()
		if commitErr != nil {
			err = fmt.Errorf("commit transaction error: %w", commitErr)
		}
	}

	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			err = fmt.Errorf("rollback transaction error: %w", rollbackErr)
		}
	}

	return err
}
