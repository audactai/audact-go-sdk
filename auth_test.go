package audact_test

import (
	"errors"
	"testing"

	audact "github.com/audactai/audact-go-sdk"
)

func TestAuthService(t *testing.T) {
	tests := []svcCase{
		{name: "exchange key", status: 200, resp: `{"access_token":"jwt","expires_in":900,"token_type":"Bearer"}`,
			call: func(c *audact.Client) (any, error) {
				r, err := c.Auth.ExchangeKey(bg)
				if err == nil && (r.AccessToken != "jwt" || r.ExpiresIn != 900) {
					t.Errorf("resp = %+v", r)
				}
				return r, err
			},
			wantM: "POST", wantPath: "/v1/auth/exchange-key", wantAuth: "Bearer apikey"},
		{name: "refresh", status: 200, resp: `{"access_token":"jwt2"}`,
			call:  func(c *audact.Client) (any, error) { return c.Auth.Refresh(bg) },
			wantM: "POST", wantPath: "/v1/auth/refresh", wantAuth: "Bearer apikey"},
	}
	runSvcCases(t, tests)
}

func TestAuthServiceNotFound(t *testing.T) {
	c := notFoundServer(t)
	if _, err := c.Auth.ExchangeKey(bg); !errors.Is(err, audact.ErrNotFound) {
		t.Fatalf("err = %v", err)
	}
	if _, err := c.Auth.Refresh(bg); !errors.Is(err, audact.ErrNotFound) {
		t.Fatalf("err = %v", err)
	}
}
