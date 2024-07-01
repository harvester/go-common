package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RandomGenerator struct {
	suite.Suite
	valueInvaild bool
}

func TestRandomGenerator(t *testing.T) {
	suite.Run(t, new(RandomGenerator))
}

func (r *RandomGenerator) SetupSuite() {
	// you could do something here before all tests
	r.valueInvaild = true
}

func (r *RandomGenerator) TestGenRandNumber() {
	val, err := GenRandNumber(10)
	require.Equal(r.T(), err, nil, "Generate random number should not get error")
	// check the validation of return value
	if val >= 0 && val < 10 {
		r.valueInvaild = true
	}
	require.Equalf(r.T(), r.valueInvaild, true, "Random number should be between [0,10), but got %d", val)
}

func (r *RandomGenerator) TestRandomValueCollision() {
	val1, err := GenRandNumber(10000)
	require.Equal(r.T(), err, nil, "Generate random number should not get error")
	val2, err := GenRandNumber(10000)
	require.Equal(r.T(), err, nil, "Generate random number should not get error")
	// check the validation of return value
	require.NotEqual(r.T(), val1, val2, "Random number should not be the same in two times")
}

func (r *RandomGenerator) TearDownSuite() {
	// you could do something here after all tests
}
