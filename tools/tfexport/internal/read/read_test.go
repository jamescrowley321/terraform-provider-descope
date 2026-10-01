package read

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/descope/terraform-provider-descope/internal/infra"
	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/discover"
	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/registry"
	"github.com/descope/terraform-provider-descope/tools/tfexport/internal/warn"
	"github.com/stretchr/testify/require"
)

func TestReadFailureIsLossyEvenForExpectedAbsence(t *testing.T) {
	ctx := context.Background()
	for _, instance := range []discover.Instance{{Resource: "descope_password_settings", ID: "Pmock"}, {Resource: "descope_widget", ID: "widget", MayBeAbsent: true}} {
		t.Run(instance.Resource, func(t *testing.T) {
			for _, status := range []int{http.StatusBadRequest, http.StatusNotFound} {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					code := "E999999"
					if status == http.StatusNotFound {
						code = "E084004"
					}
					http.Error(w, `{"errorCode":"`+code+`","errorDescription":"injected read error"}`, status)
				}))
				results, warnings := All(ctx, infra.NewClient("test", "mock-management", server.URL), registry.Load(ctx), "Pmock", []discover.Instance{instance})
				server.Close()
				require.Empty(t, results)
				require.Equal(t, status != http.StatusNotFound, warn.Incomplete(warnings))
			}
		})
	}
}
