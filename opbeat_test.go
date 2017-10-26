package logbeat

import (
	"bytes"
	"net/http"

	"github.com/jarcoal/httpmock"
)

func (suite *LogbeatTestSuite) TestOpbeatClientType() {
	suite.IsType(&OpbeatClient{}, suite.Client, "expects an instance of OpbeatClient")
}

func (suite *LogbeatTestSuite) TestOpbeatPayloadType() {
	suite.IsType(&OpbeatPayload{}, suite.Payload, "expects an instance of OpbeatPayload")
}

func (suite *LogbeatTestSuite) TestOpbeatExtraType() {
	suite.IsType(&OpbeatExtra{}, suite.Extra, "expects an instance of OpbeatExtra")
}

func (suite *LogbeatTestSuite) TestNewOpbeatClient() {
	suite.IsType(&OpbeatClient{}, suite.OpbeatClient, "expects an instance of OpbeatClient")
	suite.Equal(suite.OpbeatClient.Endpoint, suite.Endpoint, "expects the correct Endpoint URI")
	suite.IsType(&http.Client{}, suite.OpbeatClient.Http, "expects an instance of http.Client")
	suite.Equal(suite.OpbeatClient.Token, "TEST_TOKEN", "expects the correct Opbeat Token")
}

func (suite *LogbeatTestSuite) TestNewOpbeatMachine() {
	suite.IsType(OpbeatMachine{}, NewOpbeatMachine(), "expects an instance of OpbeatMachine")
	suite.Equal(suite.OpbeatMachine.Hostname, suite.Hostname, "expects the correct hostname")
}

func (suite *LogbeatTestSuite) TestNewOpbeatExtra() {
	suite.IsType(OpbeatExtra{}, NewOpbeatExtra(suite.Entry), "expects an instance of OpbeatExtra")
	suite.Equal(suite.OpbeatExtra["example"], suite.Entry.Data["example"], "expects OpbeatExtra to be the same as Log Entry Data")
}

func (suite *LogbeatTestSuite) TestOpbeatLevel() {
	suite.Equal(OpbeatLevel(suite.Entry), "critical", "expects the correct OpbeatLevel")
}

func (suite *LogbeatTestSuite) TestNewOpbeatPayload() {
	suite.IsType(&OpbeatPayload{}, suite.OpbeatPayload, "expects an instance of OpbeatPayload")
	suite.IsType(OpbeatExtra{}, suite.OpbeatPayload.Extra, "expects an instance of OpbeatExtra")
	suite.Equal(suite.OpbeatPayload.Extra["example"], suite.Entry.Data["example"], "expects OpbeatPayload Extra to be the same as Log Entry Data")
	suite.Equal(suite.OpbeatPayload.Level, "critical", "expects OpbeatPayload Level to be correct")
	suite.Contains(suite.OpbeatPayload.Logger, LogbeatVersion, "expects OpbeatPayload Logger to contain Logbeat Version")
	suite.Equal(suite.OpbeatPayload.Machine, suite.OpbeatMachine, "expects OpbeatPayload Machine to be the OpbeatMachine")
	suite.Equal(suite.OpbeatPayload.Message, suite.Entry.Message, "expects OpbeatPayload Message to be the same as Log Entry")
	suite.Equal(suite.OpbeatPayload.Timestamp, suite.Entry.Time.Format(ISO8601), "expects OpbeatPayload Timestamp to be the same as Log Entry")
}

func (suite *LogbeatTestSuite) TestJSON() {
	suite.IsType(&bytes.Buffer{}, suite.OpbeatJSON, "expects an instance of bytes.Buffer")
}

func (suite *LogbeatTestSuite) TestNewOpbeatRequest() {
	suite.IsType(&http.Request{}, suite.OpbeatRequest, "expects an instance of http.Request")
	suite.Equal(suite.OpbeatRequest.Method, "POST", "expects the correct HTTP Method")
	suite.Equal(suite.OpbeatRequest.URL, suite.EndpointURL, "expects the correct HTTP URL")
	suite.Contains(suite.OpbeatRequest.Header["Authorization"], "Bearer TEST_TOKEN", "expects the correct HTTP Header for Authorization")
	suite.Contains(suite.OpbeatRequest.Header["Content-Type"], "application/json", "expects the correct HTTP Header for Content-Type")
	suite.Equal(suite.OpbeatRequest.UserAgent(), OpbeatUserAgent(), "expects the correct HTTP Header for User Agent")
}

func (suite *LogbeatTestSuite) TestNotify() {
	httpmock.ActivateNonDefault(suite.OpbeatClient.Http)
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", suite.Endpoint, httpmock.NewStringResponder(202, ""))
	resp, err := suite.OpbeatClient.Notify(suite.Entry)
	suite.Equal(resp.StatusCode, 202, "expects HTTP Response StatusCode to be '202 Accepted'")
	suite.NoError(err, "expects Notify not to return an error")
}

func (suite *LogbeatTestSuite) TestOpbeatEndpoint() {
	suite.Equal(OpbeatEndpoint("TEST_ORG_ID", "TEST_APP_ID"), suite.Endpoint, "expects generated URI to be correct")
}

func (suite *LogbeatTestSuite) TestOpbeatBearerAuth() {
	suite.Equal(OpbeatBearerAuth("TEST_TOKEN"), "Bearer TEST_TOKEN", "expects generated Authorization Header value to be correct")
}

func (suite *LogbeatTestSuite) TestOpbeatUserAgent() {
	suite.Contains(OpbeatUserAgent(), LogbeatVersion, "expects User Agent to contain LogbeatVersion")
}
