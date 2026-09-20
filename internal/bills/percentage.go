package bills

import (
	"fmt"
	"strconv"
	"strings"
)

type Percentage int

func ParsePercentage(value string) (Percentage, error) {
	parts := strings.Split(value, ".")
	if len(parts) > 2 || len(parts[0]) == 0 || (len(parts) == 2 && len(parts[1]) == 0) {
		return 0, invalidPercentageError()
	}
	if len(parts) == 2 && len(parts[1]) > 2 {
		return 0, invalidPercentageError()
	}
	if !decimalDigits(parts[0]) || (len(parts) == 2 && !decimalDigits(parts[1])) {
		return 0, invalidPercentageError()
	}

	whole, err := strconv.ParseUint(parts[0], 10, 16)
	if err != nil {
		return 0, invalidPercentageError()
	}
	fraction := uint64(0)
	if len(parts) == 2 {
		fraction, err = strconv.ParseUint(parts[1], 10, 8)
		if err != nil {
			return 0, invalidPercentageError()
		}
		if len(parts[1]) == 1 {
			fraction *= 10
		}
	}

	basisPoints := whole*100 + fraction
	if basisPoints == 0 || basisPoints > 10_000 {
		return 0, invalidPercentageError()
	}
	return Percentage(basisPoints), nil
}

func (p Percentage) String() string {
	return fmt.Sprintf("%d.%02d", p/100, p%100)
}

func FormatCents(cents Cents) string {
	return fmt.Sprintf("%d.%02d", cents/100, cents%100)
}

func decimalDigits(value string) bool {
	for _, character := range value {
		if character < '0' || character > '9' {
			return false
		}
	}
	return true
}

func invalidPercentageError() error {
	return validationError(
		"invalid_percentage",
		"Percentage must be a decimal string between 0.01 and 100.00 with at most two decimal places",
	)
}
