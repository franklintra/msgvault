package slack

import (
	"errors"
	"net/http"
	"strings"
)

// Credentials authenticates either an OAuth user token or a browser session.
// Cookie is the raw value of Slack's d cookie, not a Cookie header.
// Browser credentials can carry broader permissions than an OAuth token.
type Credentials struct {
	Token  string
	Cookie string
}

// Validate rejects incomplete browser sessions and unsafe header values.
// Errors deliberately exclude the credential values.
func (c Credentials) Validate() error {
	if c.Token == "" || strings.ContainsAny(c.Token, " \t\r\n") {
		return errors.New("Slack token is empty or contains whitespace")
	}
	for _, r := range c.Token {
		if r < 0x21 || r > 0x7e {
			return errors.New("Slack token contains invalid header characters")
		}
	}
	browser := strings.HasPrefix(c.Token, "xoxc-")
	if browser && c.Cookie == "" {
		return errors.New("Slack browser token requires its d cookie (use --cookie-file)")
	}
	if !browser && c.Cookie != "" {
		return errors.New("Slack d cookie requires an xoxc browser token")
	}
	if c.Cookie != "" {
		cookie := &http.Cookie{Name: "d", Value: c.Cookie}
		if cookie.Valid() != nil {
			return errors.New("Slack d cookie contains invalid cookie characters")
		}
	}
	return nil
}

// NewClientWithCredentials creates a client with a validated credential pair.
func NewClientWithCredentials(baseURL string, credentials Credentials) (*Client, error) {
	if err := credentials.Validate(); err != nil {
		return nil, err
	}
	client := NewClient(baseURL, credentials.Token)
	client.cookie = credentials.Cookie
	return client, nil
}

// canRead reports whether a read capability may be attempted. Browser sessions
// do not report OAuth grants. Their server-side session permissions remain
// authoritative; the client does not synthesize OAuth scopes for them.
func (c *Client) canRead(scope string) bool {
	if c.cookie != "" {
		return true
	}
	for _, granted := range c.scopes {
		if granted == scope {
			return true
		}
	}
	return false
}

func (c *Client) authenticate(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+c.token)
	if c.cookie != "" {
		req.AddCookie(&http.Cookie{Name: "d", Value: c.cookie})
	}
}
