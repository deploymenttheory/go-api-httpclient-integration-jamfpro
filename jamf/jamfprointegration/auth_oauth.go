package jamfprointegration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.uber.org/zap"
)

// oauth implements the authInterface for Oauth2 support
type oauth struct {
	Sugar             *zap.SugaredLogger
	baseDomain        string
	clientId          string
	clientSecret      string
	bufferPeriod      time.Duration
	hideSensitiveData bool
	expiryTime        time.Time
	token             string
	http              http.Client
}

// OAuthResponse represents the response structure when obtaining an OAuth access token from JamfPro.
type OAuthResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// TODO migrate strings

// getNewToken updates the held token and expiry information
func (a *oauth) getNewToken() error {
	data := url.Values{}
	data.Set("client_id", a.clientId)
	data.Set("client_secret", a.clientSecret)
	data.Set("grant_type", "client_credentials")

	oauthComlpeteEndpoint := a.baseDomain + oAuthTokenEndpoint
	a.Sugar.Debugf("oauth endpoint constructed: %s", oauthComlpeteEndpoint)

	req, err := http.NewRequest("POST", oauthComlpeteEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	a.Sugar.Debugf("oauth token request constructed: %+v", req)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := a.http.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("bad request getting auth token: %v", resp)
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	oauthResp := &OAuthResponse{}
	err = json.Unmarshal(bodyBytes, oauthResp)
	if err != nil {
		return fmt.Errorf("failed to decode OAuth response: %w", err)
	}

	if !a.hideSensitiveData {
		a.Sugar.Debug("token recieved: %+v", oauthResp)
	}

	if oauthResp.AccessToken == "" {
		return fmt.Errorf("empty access token received")
	}

	expiresIn := time.Duration(oauthResp.ExpiresIn) * time.Second
	a.expiryTime = time.Now().Add(expiresIn)
	a.token = oauthResp.AccessToken

	a.Sugar.Infow("Token obtained successfully", "expiry", a.expiryTime)
	return nil
}

// getTokenString returns the current token as a string
func (a *oauth) getTokenString() string {
	return a.token
}

// getExpiryTime returns the current token's expiry time as a time.Time var.
func (a *oauth) getExpiryTime() time.Time {
	return a.expiryTime
}

// tokenExpired returns a bool denoting if the current token expiry time has passed.
func (a *oauth) tokenExpired() bool {
	return a.expiryTime.Before(time.Now())
}

// tokenInBuffer returns a bool denoting if the current token's duration until expiry is within the buffer period
func (a *oauth) tokenInBuffer() bool {
	return time.Until(a.expiryTime) <= a.bufferPeriod
}

// tokenEmpty returns a bool denoting if the current token string is empty.
func (a *oauth) tokenEmpty() bool {
	return a.token == ""
}
