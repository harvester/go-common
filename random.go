package gocommon

import (
	"crypto/rand"
	"math/big"
)

func GenRandNumber(boundary int64) (int64, error) {
	randNum, err := rand.Int(rand.Reader, big.NewInt(boundary))
	if err != nil {
		return 0, err
	}
	return randNum.Int64(), nil
}
