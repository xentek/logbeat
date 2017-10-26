package logbeat

import (
	"bytes"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type LogbeatTestSuite struct {
	suite.Suite
	Hook          *LogbeatHook
	Payload       *OpbeatPayload
	Extra         *OpbeatExtra
	Client        *OpbeatClient
	Org           string
	App           string
	Token         string
	Timestamp     time.Time
	Entry         *logrus.Entry
	OpbeatHook    *LogbeatHook
	OpbeatClient  *OpbeatClient
	OpbeatExtra   OpbeatExtra
	OpbeatMachine OpbeatMachine
	OpbeatPayload *OpbeatPayload
	OpbeatJSON    *bytes.Buffer
	OpbeatRequest *http.Request
	Endpoint      string
	EndpointURL   *url.URL
	Hostname      string
}

func (suite *LogbeatTestSuite) SetupTest() {
	suite.Hook = &LogbeatHook{}
	suite.Payload = &OpbeatPayload{}
	suite.Extra = &OpbeatExtra{}
	suite.Client = &OpbeatClient{}

	suite.Org = "TEST_ORG_ID"
	suite.App = "TEST_APP_ID"
	suite.Token = "TEST_TOKEN"

	suite.Entry = &logrus.Entry{
		Level:   logrus.PanicLevel,
		Message: "Example Logbeat Log Entry",
		Data:    logrus.Fields{"example": "true"},
		Time:    time.Date(1955, time.November, 05, 9, 11, 12, 13, time.UTC),
	}

	suite.OpbeatHook = NewOpbeatHook(suite.Org, suite.App, suite.Token)
	suite.OpbeatClient = NewOpbeatClient(suite.Org, suite.App, suite.Token)
	suite.OpbeatPayload = NewOpbeatPayload(suite.Entry)

	jsonPayload, _ := suite.OpbeatPayload.JSON()
	suite.OpbeatJSON = jsonPayload

	req, _ := suite.OpbeatClient.NewOpbeatRequest(suite.OpbeatJSON)
	suite.OpbeatRequest = req

	suite.Endpoint = "https://intake.opbeat.com/api/v1/organizations/TEST_ORG_ID/apps/TEST_APP_ID/errors/"
	parsedEndpoint, _ := url.Parse(suite.Endpoint)
	suite.EndpointURL = parsedEndpoint

	hostname, _ := os.Hostname()
	suite.Hostname = hostname
	suite.OpbeatMachine = NewOpbeatMachine()
	suite.OpbeatExtra = NewOpbeatExtra(suite.Entry)
}

func (suite *LogbeatTestSuite) TestLogbeatHookType() {
	suite.IsType(&LogbeatHook{}, suite.Hook, "expects an instance of LogbeatHook")
}

func (suite *LogbeatTestSuite) TestLogbeatHookInterface() {
	suite.Implements((*logrus.Hook)(nil), new(LogbeatHook), "expects OpbeatHook to implment logrus.Hook interface")
}

func (suite *LogbeatTestSuite) TestNewOpbeatHook() {
	suite.IsType(&LogbeatHook{}, suite.OpbeatHook, "expects an instance of LogbeatHook.")
	suite.Equal(suite.OpbeatHook.AppId, "TEST_APP_ID", "expects the correct Opbeat App ID.")
	suite.IsType(&OpbeatClient{}, suite.OpbeatHook.Opbeat, "expects an instance of OpbeatClient.")
	suite.Equal(suite.OpbeatHook.OrgId, "TEST_ORG_ID", "expects the correct Opbeat Organization ID.")
	suite.Equal(suite.OpbeatHook.SecretToken, "TEST_TOKEN", "expects the correct Opbeat Secret Token.")
}

func (suite *LogbeatTestSuite) TestFire() {
	httpmock.ActivateNonDefault(suite.OpbeatHook.Opbeat.Http)
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", suite.Endpoint, httpmock.NewStringResponder(202, ""))
	suite.NoError(suite.OpbeatHook.Fire(suite.Entry), "expects LobeatHook.Fire not to return an error")
}

func (suite *LogbeatTestSuite) TestLevels() {
	suite.IsType([]logrus.Level{}, suite.Hook.Levels(), "expects an instance of OpbeatExtra")
	suite.Contains(suite.Hook.Levels(), logrus.PanicLevel, "expects Levels to contain PanicLevel")
	suite.Contains(suite.Hook.Levels(), logrus.FatalLevel, "expects Levels to contain FatalLevel")
	suite.Contains(suite.Hook.Levels(), logrus.ErrorLevel, "expects Levels to contain ErrorLevel")
}

func TestLogbeatSuite(t *testing.T) {
	suite.Run(t, new(LogbeatTestSuite))
}
