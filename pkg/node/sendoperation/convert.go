package sendoperation

import (
	"fmt"

	"github.com/shopspring/decimal"
)

const MasDecimals = 9

func NanoToMas(nanoMasAmount uint64) (string, error) {
	dec, err := decimal.NewFromString(fmt.Sprint(nanoMasAmount))
	if err != nil {
		return "", fmt.Errorf("converting '%d' to decimal: %w", nanoMasAmount, err)
	}

	return dec.Shift(-MasDecimals).String(), nil
}
