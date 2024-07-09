package command

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CommandTestSuite struct {
	suite.Suite
	executor *Executor
}

func (suite *CommandTestSuite) SetupTest() {
	suite.executor = NewExecutor()
}

func (suite *CommandTestSuite) TestNewExecutor() {
	assert.NotNil(suite.T(), suite.executor)
}

func (suite *CommandTestSuite) TestSetTimeout() {
	suite.executor.SetTimeout(5 * time.Second)
	assert.Equal(suite.T(), 5*time.Second, suite.executor.cmdTimeout)
}

func (suite *CommandTestSuite) TestExecute() {
	output, err := suite.executor.Execute("echo", []string{"hello", "world"})
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "hello world\n", output)
}

func (suite *CommandTestSuite) TestExecute_Timeout() {
	suite.executor.SetTimeout(1 * time.Second)
	assert.Equal(suite.T(), 1*time.Second, suite.executor.cmdTimeout)
	output, err := suite.executor.Execute("sleep", []string{"5"})
	assert.NotNil(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "command timeout")
	assert.Empty(suite.T(), output)
}

func TestCommandTestSuite(t *testing.T) {
	suite.Run(t, new(CommandTestSuite))
}
