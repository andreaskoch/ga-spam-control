package credentials

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

func Test_Serialize_EmptyToken_NoErrorIsReturned(t *testing.T) {
	// arrange
	writeBuffer := new(bytes.Buffer)
	serializer := &tokenSerializer{}
	token := &oauth2.Token{}

	// act
	err := serializer.Serialize(writeBuffer, token)

	// assert
	if err != nil {
		t.Fail()
		t.Logf("Serialize returned an error: %s", err.Error())
	}
}

func Test_Serialize_EmptyToken_JSONIsWritten(t *testing.T) {
	// arrange
	writeBuffer := new(bytes.Buffer)
	serializer := &tokenSerializer{}
	token := &oauth2.Token{}

	// act
	serializer.Serialize(writeBuffer, token)

	// assert
	if writeBuffer.String() == "" {
		t.Fail()
		t.Logf("Serialize did not write JSON")
	}
}

func Test_Serialize_InitializedToken_JSONIsWritten(t *testing.T) {
	// arrange
	writeBuffer := new(bytes.Buffer)
	serializer := &tokenSerializer{}
	token := &oauth2.Token{
		AccessToken:  "dsajdkasdljaslkjdl382109382109382109",
		TokenType:    "tokenType",
		RefreshToken: "refreshToken",
		Expiry:       time.Now(),
	}

	// act
	serializer.Serialize(writeBuffer, token)

	// assert
	if !strings.Contains(writeBuffer.String(), "dsajdkasdljaslkjdl382109382109382109") {
		t.Fail()
		t.Logf("Serialize did not write the access token")
	}
}

func Test_Deserialize_EmptyString_ErrorIsReturned(t *testing.T) {
	// arrange
	jsonReader := strings.NewReader(``)
	serializer := &tokenSerializer{}

	// act
	_, err := serializer.Deserialize(jsonReader)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("Deserialize did not return an error")
	}
}

func Test_Deserialize_InvalidJSON_ErrorIsReturned(t *testing.T) {
	// arrange
	jsonReader := strings.NewReader(`ååasdjaskldjas
    dsadöl`)
	serializer := &tokenSerializer{}

	// act
	_, err := serializer.Deserialize(jsonReader)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("Deserialize did not return an error")
	}
}

func Test_Deserialize_ValidJSON_NoErrorIsReturned(t *testing.T) {
	// arrange
	jsonReader := strings.NewReader(`{
	"access_token": "daskdlöasdksalökdlöasklökAccessToken",
	"token_type": "Bearer",
	"refresh_token": "askdlökasldökaslöRefreshToken",
	"expiry": "2016-04-16T17:52:05.265946629+02:00"
}`)
	serializer := &tokenSerializer{}

	// act
	_, err := serializer.Deserialize(jsonReader)

	// assert
	if err != nil {
		t.Fail()
		t.Logf("Deserialize returned an error: %s", err.Error())
	}
}

func Test_Deserialize_ValidJSON_TokenIsReturned(t *testing.T) {
	// arrange
	jsonReader := strings.NewReader(`{
	"access_token": "daskdlöasdksalökdlöasklökAccessToken",
	"token_type": "Bearer",
	"refresh_token": "askdlökasldökaslöRefreshToken",
	"expiry": "2016-04-16T17:52:05.265946629+02:00"
}`)
	serializer := &tokenSerializer{}

	// act
	token, _ := serializer.Deserialize(jsonReader)

	// assert
	if token.AccessToken != "daskdlöasdksalökdlöasklökAccessToken" {
		t.Fail()
		t.Logf("Deserialize the wrong token")
	}
}
