package sendoperation

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/shopspring/decimal"
)

const (
	MasDecimals = 9
	base        = 10
)

// NanoToMas Converts NanoMAS amount to MAS, returning it as a string.
func NanoToMas(nanoMasAmount uint64) (string, error) {
	dec, err := decimal.NewFromString(strconv.FormatUint(nanoMasAmount, 10))
	if err != nil {
		return "", fmt.Errorf("converting '%d' to decimal: %w", nanoMasAmount, err)
	}

	return dec.Shift(-MasDecimals).String(), nil
}

// MasToNano Converts MAS amount (as a string) to NanoMAS, returning it as an uint64.
func MasToNano(masAmount string) (uint64, error) {
	dec, err := decimal.NewFromString(masAmount)
	if err != nil {
		return 0, fmt.Errorf("converting '%s' to decimal: %w", masAmount, err)
	}

	nanoDec := dec.Shift(MasDecimals)
	// Use big.Int to safely handle the conversion to uint64
	nanoInt, ok := new(big.Int).SetString(nanoDec.String(), base)
	if !ok {
		return 0, fmt.Errorf("error converting '%s' to big.Int", nanoDec.String())
	}

	// Ensure the result fits in uint64
	if !nanoInt.IsUint64() {
		return 0, fmt.Errorf("value '%s' overflows uint64", nanoDec.String())
	}

	return nanoInt.Uint64(), nil
}
