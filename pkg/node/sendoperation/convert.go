package sendoperation

import (
	"fmt"
	"strconv"

	"github.com/shopspring/decimal"
)

const MasDecimals = 9

func NanoToMas(nanoMasAmount uint64) (string, error) {
	dec, err := decimal.NewFromString(strconv.FormatUint(nanoMasAmount, 10))
	if err != nil {
		return "", fmt.Errorf("converting '%d' to decimal: %w", nanoMasAmount, err)
	}

	return dec.Shift(-MasDecimals).String(), nil
}
