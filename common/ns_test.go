package common

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type NSTestSuite struct {
	suite.Suite
	hostProcPath string
}

func (suite *NSTestSuite) SetupTest() {
	suite.hostProcPath = "/proc"
}

func TestNSTestSuite(t *testing.T) {
	suite.Run(t, new(NSTestSuite))
}

func (suite *NSTestSuite) TestGetPidProc() {
	pid := 66666 // should not exists
	proc, err := getPidProc(suite.hostProcPath, pid)

	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), proc)

	pid = 1
	proc, err = getPidProc(suite.hostProcPath, pid)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), proc)
	assert.Equal(suite.T(), pid, proc.PID)
}

func (suite *NSTestSuite) TestGetSelfProc() {
	proc, err := getSelfProc(suite.hostProcPath)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), proc)
	assert.Equal(suite.T(), os.Getpid(), proc.PID)
}

func (suite *NSTestSuite) TestGetHostNamespacePath() {
	expectedPath := "/proc/1/ns/"
	path := GetHostNamespacePath(suite.hostProcPath)

	assert.Equal(suite.T(), expectedPath, path)
}
