package commands

// Basic imports
import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/malamtime/cli/model"
	"github.com/malamtime/cli/model/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/urfave/cli/v2"
)

type trackTestSuite struct {
	suite.Suite
	baseTimeFolder string
}

// before each test
func (s *trackTestSuite) SetupSuite() {
	s.baseTimeFolder = strconv.Itoa(int(time.Now().Unix()))
}

func (s *trackTestSuite) TestMultipTrackWithPre() {
	cs := mocks.NewConfigService(s.T())
	mockedConfig := model.MalamTimeConfig{}
	cs.On("ReadConfigFile").Return(mockedConfig, nil)
	model.UserMalamTimeConfig = mockedConfig
	configService = cs

	baseFolder := s.baseTimeFolder + "-withPre"
	model.InitFolder(baseFolder)

	p := filepath.Join(os.Getenv("HOME"), model.COMMAND_STORAGE_FOLDER)
	err := os.MkdirAll(p, os.ModePerm)
	assert.Nil(s.T(), err)

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

func (s *trackTestSuite) TestTrackWithSendData() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		body, err := io.ReadAll(r.Body)
		assert.Nil(s.T(), err)
		defer r.Body.Close()

		var payload model.PostTrackArgs

		err = json.Unmarshal(body, &payload)
		assert.Nil(s.T(), err)

		assert.Len(s.T(), payload.Data, 7)

		assert.Contains(s.T(), string(body), "fish")
		assert.EqualValues(s.T(), "CLI TOKEN001", authorizationHeader)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()
	cs := mocks.NewConfigService(s.T())
	mockedConfig := model.MalamTimeConfig{
		Token:       "TOKEN001",
		APIEndpoint: server.URL,
		FlushCount:  7,
		GCTime:      8,
	}
	cs.On("ReadConfigFile").Return(mockedConfig, nil)
	model.UserMalamTimeConfig = mockedConfig
	configService = cs

	baseFolder := s.baseTimeFolder + "-sendData"
	model.InitFolder(baseFolder)

	err := os.MkdirAll(filepath.Join(os.Getenv("HOME"), model.COMMAND_STORAGE_FOLDER), os.ModePerm)
	assert.Nil(s.T(), err)

	app := &cli.App{
		// mtt for malamtime-testing
		Name:  "mtt2",
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
			"mtt2",
			"track",
			"-s=fish",
			fmt.Sprintf("-id=%d", sessionID),
			fmt.Sprintf("-cmd=cmd1 %d", i),
			"-p=pre",
		}
		postCommand := []string{
			"mtt2",
			"track",
			"-s=fish",
			fmt.Sprintf("-id=%d", sessionID),
			fmt.Sprintf("-cmd=cmd1 %d", i),
			"-p=post",
		}
		go func(cmd []string, pc []string) {
			err := app.Run(cmd)
			assert.Nil(s.T(), err)
			time.Sleep(time.Millisecond * 100)
			err = app.Run(pc)
			assert.Nil(s.T(), err)
			wg.Done()
		}(command, postCommand)
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
	err := os.RemoveAll(filepath.Join(os.Getenv("HOME"), ".malamtime-"+s.baseTimeFolder+"-withPre"))
	assert.Nil(s.T(), err)

	// Delete the test folder
	err = os.RemoveAll(filepath.Join(os.Getenv("HOME"), ".malamtime-"+s.baseTimeFolder+"-sendData"))
	assert.Nil(s.T(), err)
}

func TestTrackTestSuite(t *testing.T) {
	suite.Run(t, new(trackTestSuite))
}
