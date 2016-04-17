package credentials

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

func Test_GetToken_FileNotFound_ErrorIsReturned(t *testing.T) {
	// arrange
	tokenStore := TokenStore{
		filePath:   "/tmp/file-not-found",
		serializer: tokenSerializer{},
	}

	// act
	_, err := tokenStore.GetToken()

	// assert
	if err == nil {
		t.Fail()
		t.Logf("GetToken did not return an error")
	}
}

func Test_GetToken_FileEmpty_ErrorIsReturned(t *testing.T) {
	// arrange
	content := []byte("")
	tmpfile, _ := ioutil.TempFile("", ".tokenstoretest")
	filePath := tmpfile.Name()

	defer os.Remove(filePath)

	tmpfile.Write(content)
	tmpfile.Close()

	tokenStore := TokenStore{
		filePath:   filePath,
		serializer: tokenSerializer{},
	}

	// act
	_, err := tokenStore.GetToken()

	// assert
	if err == nil {
		t.Fail()
		t.Logf("GetToken did not return an error")
	}
}

func Test_GetToken_FileContainsGarbage_ErrorIsReturned(t *testing.T) {
	// arrange
	content := []byte("la di da")
	tmpfile, _ := ioutil.TempFile("", ".tokenstoretest")
	filePath := tmpfile.Name()

	defer os.Remove(filePath)

	tmpfile.Write(content)
	tmpfile.Close()

	tokenStore := TokenStore{
		filePath:   filePath,
		serializer: tokenSerializer{},
	}

	// act
	_, err := tokenStore.GetToken()

	// assert
	if err == nil {
		t.Fail()
		t.Logf("GetToken did not return an error")
	}
}

func Test_GetToken_FileExists_TokenIsReturned(t *testing.T) {
	// arrange
	content := []byte(`{
	"access_token": "daskdlöasdksalökdlöasklökAccessToken",
	"token_type": "Bearer",
	"refresh_token": "askdlökasldökaslöRefreshToken",
	"expiry": "2016-04-16T17:52:05.265946629+02:00"
}`)
	tmpfile, _ := ioutil.TempFile("", ".tokenstoretest")
	filePath := tmpfile.Name()

	defer os.Remove(filePath)

	tmpfile.Write(content)
	tmpfile.Close()

	tokenStore := TokenStore{
		filePath:   filePath,
		serializer: tokenSerializer{},
	}

	// act
	token, _ := tokenStore.GetToken()

	// assert
	if token.AccessToken != "daskdlöasdksalökdlöasklökAccessToken" {
		t.Fail()
		t.Logf("The wrong AccessToken was returned")
	}
}

func Test_SaveToken_NilGiven_FileIsCreated(t *testing.T) {
	// arrange
	tmpfile, _ := ioutil.TempFile("", ".tokenstoretest")
	filePath := tmpfile.Name()

	tmpfile.Close()
	os.Remove(filePath)

	tokenStore := TokenStore{
		filePath:   filePath,
		serializer: tokenSerializer{},
	}

	// act
	tokenStore.SaveToken(nil)

	// assert
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fail()
		t.Logf("SaveToken did not create a file")
	}
}

func Test_SaveToken_NilGiven_FileIsTruncated(t *testing.T) {
	// arrange
	content := []byte(`{
	"access_token": "daskdlöasdksalökdlöasklökAccessToken",
	"token_type": "Bearer",
	"refresh_token": "askdlökasldökaslöRefreshToken",
	"expiry": "2016-04-16T17:52:05.265946629+02:00"
}`)
	tmpfile, _ := ioutil.TempFile("", ".tokenstoretest")
	filePath := tmpfile.Name()

	defer os.Remove(filePath)

	tmpfile.Write(content)
	tmpfile.Close()

	tokenStore := TokenStore{
		filePath:   filePath,
		serializer: tokenSerializer{},
	}

	// act
	tokenStore.SaveToken(nil)

	// assert
	file, _ := os.Open(filePath)
	writtenContent, _ := ioutil.ReadAll(file)
	file.Close()

	if len(content) == len(writtenContent) {
		t.Fail()
		t.Logf("SaveToken did not override the existing file")
	}
}

func Test_SaveToken_TokenGiven_NoErrorIsReturned(t *testing.T) {
	// arrange
	tmpfile, _ := ioutil.TempFile("", ".tokenstoretest")
	filePath := tmpfile.Name()

	tmpfile.Close()

	token := &oauth2.Token{
		AccessToken:  "dsajdkasdljaslkjdl382109382109382109",
		TokenType:    "tokenType",
		RefreshToken: "refreshToken",
		Expiry:       time.Now(),
	}

	tokenStore := TokenStore{
		filePath:   filePath,
		serializer: tokenSerializer{},
	}

	// act
	err := tokenStore.SaveToken(token)

	// assert
	if err != nil {
		t.Fail()
		t.Logf("SaveToken returned an error: %s", err.Error())
	}
}

func Test_SaveToken_TokenGiven_FileIsWritten(t *testing.T) {
	// arrange
	tmpfile, _ := ioutil.TempFile("", ".tokenstoretest")
	filePath := tmpfile.Name()

	tmpfile.Close()

	token := &oauth2.Token{
		AccessToken:  "dsajdkasdljaslkjdl382109382109382109",
		TokenType:    "tokenType",
		RefreshToken: "refreshToken",
		Expiry:       time.Now(),
	}

	tokenStore := TokenStore{
		filePath:   filePath,
		serializer: tokenSerializer{},
	}

	// act
	tokenStore.SaveToken(token)

	// assert
	file, _ := os.Open(filePath)
	writtenContent, _ := ioutil.ReadAll(file)
	file.Close()

	if !strings.Contains(string(writtenContent), "dsajdkasdljaslkjdl382109382109382109") {
		t.Fail()
		t.Logf("SaveToken did not writen the access token to the file: %s", writtenContent)
	}
}
