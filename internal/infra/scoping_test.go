package infra

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSDKProjectIsolation(t *testing.T) {
	var lock sync.Mutex
	seen := map[string]int{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/mgmt/tenant" {
			lock.Lock()
			seen[r.Header.Get("Authorization")]++
			lock.Unlock()
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"tenant":{"id":"Ttest","name":"test"}}`))
	}))
	defer server.Close()
	t.Setenv("DESCOPE_PROJECT_ID", "unexpected-project")
	c := NewClient("test", "mock-management", server.URL)
	_, err := c.Management("")
	require.Error(t, err)
	_, err = c.Management("   ")
	require.Error(t, err)
	for _, project := range []string{"Pfirst", "Psecond", "Pfirst"} {
		management, err := c.Management(project)
		require.NoError(t, err)
		_, err = management.Tenant().Load(context.Background(), "Ttest")
		require.NoError(t, err)
	}
	lock.Lock()
	defer lock.Unlock()
	require.Equal(t, map[string]int{"Bearer Pfirst:mock-management": 2, "Bearer Psecond:mock-management": 1}, seen)
}
