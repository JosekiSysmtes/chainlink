package adapters

import (
	"errors"
	"strconv"

	"chainlink/core/store"
	"chainlink/core/store/models"
)

// Compare adapter type takes an Operator and a Value field to
// compare to the previous task's Result.
type Compare struct {
	*NoValidationCheckAdapter
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

var (
	ErrResultNotNumber      = errors.New("The result was not a number")
	ErrValueNotNumber       = errors.New("The value was not a number")
	ErrOperatorNotSpecified = errors.New("Operator not specified")
	ErrValueNotSpecified    = errors.New("Value not specified")
)

// Perform uses the Operator to check the run's result against the
// specified Value.
func (c *Compare) Perform(input models.RunInput, _ *store.Store) models.RunOutput {
	prevResult := input.Result().String()

	if c.Value == "" {
		return models.NewRunOutputError(ErrValueNotSpecified)
	}

	switch c.Operator {
	case "eq":
		return models.NewRunOutputCompleteWithResult(c.Value == prevResult)
	case "neq":
		return models.NewRunOutputCompleteWithResult(c.Value != prevResult)
	case "gt":
		value, desired, err := getValues(prevResult, c.Value)
		if err != nil {
			return models.NewRunOutputError(err)
		}
		return models.NewRunOutputCompleteWithResult(desired < value)
	case "gte":
		value, desired, err := getValues(prevResult, c.Value)
		if err != nil {
			return models.NewRunOutputError(err)
		}
		return models.NewRunOutputCompleteWithResult(desired <= value)
	case "lt":
		value, desired, err := getValues(prevResult, c.Value)
		if err != nil {
			return models.NewRunOutputError(err)
		}
		return models.NewRunOutputCompleteWithResult(desired > value)
	case "lte":
		value, desired, err := getValues(prevResult, c.Value)
		if err != nil {
			return models.NewRunOutputError(err)
		}
		return models.NewRunOutputCompleteWithResult(desired >= value)
	default:
		return models.NewRunOutputError(ErrOperatorNotSpecified)
	}
}

func getValues(result string, d string) (float64, float64, error) {
	value, err := strconv.ParseFloat(result, 64)
	if err != nil {
		return 0, 0, ErrResultNotNumber
	}
	desired, err := strconv.ParseFloat(d, 64)
	if err != nil {
		return 0, 0, ErrValueNotNumber
	}
	return value, desired, nil
}
