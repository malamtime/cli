package commands

// Basic imports
import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/malamtime/cli/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/urfave/cli/v2"
)

type trackTestSuite struct {
	suite.Suite
}

// before each test
func (s *trackTestSuite) SetupSuite() {
	model.InitFolder(true)
}

func (s *trackTestSuite) TestMultipTrackWithPre() {

	app := &cli.App{
		// mtt for malamtime-testing
		Name:  "mtt",
		Usage: "just for testing",
		Commands: []*cli.Command{
			TrackCommand,
		},
	}

	times := 10

	var wg sync.WaitGroup
	wg.Add(times)

	sessionID := time.Now().Unix()

	for i := 0; i < times; i++ {
		command := []string{
			"mtt",
			"track",
			"-s=fish",
			fmt.Sprintf("-id=%d", sessionID),
			fmt.Sprintf("-cmd=cmd1 %d", i),
			"-p=pre",
		}
		go func(cmd []string) {
			err := app.Run(cmd)
			assert.Nil(s.T(), err)
			wg.Done()
		}(command)
	}

	wg.Wait()

	// Check the number of lines in the COMMAND_PRE_STORAGE_FILE
	preFile := os.ExpandEnv("$HOME/" + model.COMMAND_PRE_STORAGE_FILE)
	content, err := os.ReadFile(preFile)
	assert.Nil(s.T(), err)

	lines := 0
	for _, byte := range content {
		if byte == '\n' {
			lines++
		}
	}
	assert.Equal(s.T(), times, lines)
}

func (s *trackTestSuite) TearDownSuite() {
	// Delete the test folder
	err := os.RemoveAll(filepath.Join(os.Getenv("HOME"), ".malamtime-testing"))

	if err != nil {
		s.T().Errorf("Failed to remove testing folder: %v", err)
	}
}

func TestTrackTestSuite(t *testing.T) {
	suite.Run(t, new(trackTestSuite))
}
