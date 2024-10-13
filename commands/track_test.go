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
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/malamtime/cli/model"
	"github.com/malamtime/cli/model/mocks"
	"github.com/sirupsen/logrus"
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
	logrus.SetLevel(logrus.TraceLevel)
	reqCursor := make([]int64, 0)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		body, err := io.ReadAll(r.Body)
		assert.Nil(s.T(), err)
		defer r.Body.Close()

		var payload model.PostTrackArgs

		err = json.Unmarshal(body, &payload)
		assert.Nil(s.T(), err)

		assert.GreaterOrEqual(s.T(), len(payload.Data), 7)

		assert.Contains(s.T(), string(body), "fish")
		assert.EqualValues(s.T(), "CLI TOKEN001", authorizationHeader)
		w.WriteHeader(http.StatusNoContent)
		reqCursor = append(reqCursor, payload.CursorID)
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

	times := 16

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

	assert.GreaterOrEqual(s.T(), len(reqCursor), 2)

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

	// Check the number of lines in the COMMAND_POST_STORAGE_FILE
	postFile := os.ExpandEnv("$HOME/" + model.COMMAND_POST_STORAGE_FILE)
	postContent, err := os.ReadFile(postFile)
	assert.Nil(s.T(), err)

	postLines := 0
	for _, byte := range postContent {
		if byte == '\n' {
			postLines++
		}
	}
	assert.Equal(s.T(), lines, postLines)

	// Check the CURSOR_FILE
	cursorFile := os.ExpandEnv("$HOME/" + model.COMMAND_CURSOR_STORAGE_FILE)
	cursorContent, err := os.ReadFile(cursorFile)
	assert.Nil(s.T(), err)

	var cursorValues []time.Time
	for _, line := range strings.Split(string(cursorContent), "\n") {
		if line != "" {
			nanoTime, err := strconv.ParseInt(line, 10, 64)
			assert.Nil(s.T(), err)
			cursorValues = append(cursorValues, time.Unix(0, nanoTime))
		}
	}
	assert.GreaterOrEqual(s.T(), len(cursorValues), 2)

	assert.True(s.T(), cursorValues[len(cursorValues)-1].After(cursorValues[0]))

	reqCursorStr := strings.Join(strings.Fields(fmt.Sprint(reqCursor)), "\t")

	for _, value := range cursorValues {
		cursorInStr := strconv.FormatInt(value.UnixNano(), 10)
		assert.Contains(s.T(), string(postContent), cursorInStr)
		assert.Contains(s.T(), reqCursorStr, cursorInStr)
	}
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
