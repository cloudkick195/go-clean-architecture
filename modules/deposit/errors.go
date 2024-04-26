package deposit

import (
	"errors"
	"go_clean_architecture/commons"
)

var (
	DepositHasBeenCompleted = commons.NewCustomError(
		errors.New("DepositHasBeenCompleted"),
		"DepositHasBeenCompleted",
		"DepositHasBeenCompleted",
	)
)
