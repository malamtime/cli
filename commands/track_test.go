package commands

// Basic imports
import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type trackTestSuite struct {
	suite.Suite
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (s *trackTestSuite) SetupTest() {
}

// All methods that begin with "Test" are run as tests within a
// suite.
func (s *trackTestSuite) TestExample() {
	assert.Equal(s.T(), 5, 5)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestTrackTestSuite(t *testing.T) {
	suite.Run(t, new(trackTestSuite))
}
