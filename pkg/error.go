package pkg

import "errors"

var (
	ErrFailedGet       = errors.New("failed to get data")
	ErrFailedSave      = errors.New("failed to save data")
	ErrFailedDelete    = errors.New("failed to delete data")
	ErrIDNotFound      = errors.New("id not found")
	ErrDateIsIncorrect = errors.New("date is incorrect")
	ErrNoForeignKey    = errors.New("Error 1452: Cannot add or update a child row: a foreign key constraint fails")
	ErrPriceNotValid   = errors.New("incorrect price entry")
	ErrIdNotValid      = errors.New("incorrect id entry")
)
