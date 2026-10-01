//go:build integration || fork

package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/descope/go-sdk/descope"
	"github.com/stretchr/testify/require"
)

func TestM2MTokenLifecycle(t *testing.T) {
	h := NewHarness(t)
	h.LoadFixture("m2m/main.tf")
	name := "name=" + GenerateName(t)
	h.Apply(name)
	projectID := StringAttr(h.StateResource("descope_project.test"), "id")
	client := newSDKClientWithProject(t, projectID)
	key := h.StateResource("descope_access_key.service")
	keyID, secret := StringAttr(key, "id"), StringAttr(key, "cleartext")
	require.NotEmpty(t, secret)
	exchange := func(accessKey string, options *descope.AccessKeyLoginOptions) *descope.Token {
		t.Helper()
		ok, token, err := client.Auth.ExchangeAccessKey(context.Background(), accessKey, options)
		require.NoError(t, err)
		require.True(t, ok)
		require.NotNil(t, token)
		valid, verified, err := client.Auth.ValidateSessionWithToken(context.Background(), token.JWT)
		require.NoError(t, err)
		require.True(t, valid)
		return verified
	}
	oauth := func(wantSuccess bool) *descope.Token {
		t.Helper()
		form := url.Values{"grant_type": {"client_credentials"}, "scope": {"openid descope.claims descope.custom_claims"}}
		request, err := http.NewRequestWithContext(context.Background(), http.MethodPost, strings.TrimRight(os.Getenv("DESCOPE_BASE_URL"), "/")+"/oauth2/v1/token", strings.NewReader(form.Encode()))
		require.NoError(t, err)
		request.SetBasicAuth(StringAttr(key, "client_id"), secret)
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		response, err := (&http.Client{Timeout: 30 * time.Second}).Do(request)
		require.NoError(t, err)
		defer func() { _ = response.Body.Close() }()
		if !wantSuccess {
			require.Contains(t, []int{http.StatusBadRequest, http.StatusUnauthorized}, response.StatusCode)
			return nil
		}
		require.Equal(t, http.StatusOK, response.StatusCode)
		var tokens struct {
			AccessToken string `json:"access_token"`
		}
		require.NoError(t, json.NewDecoder(response.Body).Decode(&tokens))
		require.NotEmpty(t, tokens.AccessToken)
		ok, verified, err := client.Auth.ValidateSessionWithToken(context.Background(), tokens.AccessToken)
		require.NoError(t, err)
		require.True(t, ok)
		return verified
	}

	stored := LoadAccessKeyViaSDK(t, projectID, keyID)
	require.Equal(t, "worker", stored.CustomAttributes["owner"])
	token := exchange(secret, nil)
	require.Equal(t, "worker", token.Claims["service"])
	require.Equal(t, "worker", oauth(true).Claims["service"])
	require.Equal(t, true, token.Claims["enabled"])
	require.EqualValues(t, 3, token.Claims["retries"])
	require.Equal(t, "terraform-m2m-test", token.Claims["purpose"])
	require.Contains(t, token.Claims["roles"], "service")
	require.Contains(t, token.Claims[descope.ClaimAuthorizedGlobalPermissions], "read:jobs")
	tenantID := StringAttr(h.StateResource("descope_tenant.test"), "id")
	tenantSecret := StringAttr(h.StateResource("descope_access_key.tenant_service"), "cleartext")
	tenantToken := exchange(tenantSecret, &descope.AccessKeyLoginOptions{SelectedTenant: tenantID})
	require.Contains(t, tenantToken.GetTenantValue(tenantID, "roles"), "tenant-service")
	require.Contains(t, tenantToken.GetTenantValue(tenantID, descope.ClaimAuthorizedGlobalPermissions), "read:jobs")

	h.Apply(name, "service=updated")
	require.Equal(t, keyID, StringAttr(h.StateResource("descope_access_key.service"), "id"))
	require.Equal(t, secret, StringAttr(h.StateResource("descope_access_key.service"), "cleartext"))
	require.Equal(t, "updated", exchange(secret, nil).Claims["service"])
	require.Equal(t, "updated", oauth(true).Claims["service"])

	_, err := client.Management.AccessKey().Update(context.Background(), keyID, "service", nil, nil, nil, map[string]any{"service": "drift"}, nil, map[string]any{"owner": "drift"})
	require.NoError(t, err)
	h.Apply(name, "service=updated")
	require.Equal(t, "updated", exchange(secret, nil).Claims["service"])
	require.Equal(t, "updated", oauth(true).Claims["service"])

	h.Apply(name, "clear=true")
	cleared := exchange(secret, nil)
	require.NotContains(t, cleared.Claims, "service")
	require.Empty(t, cleared.Claims["roles"])
	require.Empty(t, cleared.Claims[descope.ClaimAuthorizedGlobalPermissions])
	require.NotContains(t, oauth(true).Claims, "service")
	require.Equal(t, "{}", StringAttr(h.StateResource("descope_access_key.service"), "custom_claims"))
	require.Equal(t, "{}", StringAttr(h.StateResource("descope_access_key.service"), "custom_attributes"))
	require.Empty(t, LoadAccessKeyViaSDK(t, projectID, keyID).CustomAttributes)

	h.Apply(name, "clear=true", "status=inactive")
	ok, _, err := client.Auth.ExchangeAccessKey(context.Background(), secret, nil)
	require.False(t, ok)
	require.Error(t, err)
	oauth(false)
	h.Apply(name, "service=restored")
	require.Equal(t, "restored", exchange(secret, nil).Claims["service"])
	h.Destroy(name, "service=restored")
	require.False(t, h.HasState())
}
