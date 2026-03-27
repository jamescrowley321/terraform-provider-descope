//go:build integration || fork

package integration

import (
	"context"
	"os"
	"testing"

	"github.com/descope/go-sdk/descope"
	descopeclient "github.com/descope/go-sdk/descope/client"
	"github.com/stretchr/testify/require"
)

// newProjectSDKClient creates a Descope management SDK client scoped to the
// project identified by the DESCOPE_PROJECT_ID environment variable.
func newProjectSDKClient(t *testing.T) *descopeclient.DescopeClient {
	t.Helper()
	pid := os.Getenv("DESCOPE_PROJECT_ID")
	require.NotEmpty(t, pid, "DESCOPE_PROJECT_ID must be set for SDK verification")
	return newSDKClientWithProject(t, pid)
}

// newCompanySDKClient creates a Descope management SDK client without a project
// ID, for verifying company-level resources such as descopers and management keys.
func newCompanySDKClient(t *testing.T) *descopeclient.DescopeClient {
	t.Helper()
	client, err := descopeclient.NewWithConfig(&descopeclient.Config{
		ManagementKey:       os.Getenv("DESCOPE_MANAGEMENT_KEY"),
		DescopeBaseURL:      os.Getenv("DESCOPE_BASE_URL"),
		AllowEmptyProjectID: true,
	})
	require.NoError(t, err, "creating company SDK client")
	return client
}

// newSDKClientWithProject creates a Descope management SDK client scoped to the given project.
func newSDKClientWithProject(t *testing.T, projectID string) *descopeclient.DescopeClient {
	t.Helper()
	client, err := descopeclient.NewWithConfig(&descopeclient.Config{
		ManagementKey:  os.Getenv("DESCOPE_MANAGEMENT_KEY"),
		DescopeBaseURL: os.Getenv("DESCOPE_BASE_URL"),
		ProjectID:      projectID,
	})
	require.NoError(t, err, "creating SDK client for project %s", projectID)
	return client
}

// --- Project-scoped resource loaders ---

// LoadTenantViaSDK loads a tenant directly from the Descope API.
func LoadTenantViaSDK(t *testing.T, id string) *descope.Tenant {
	t.Helper()
	client := newProjectSDKClient(t)
	tenant, err := client.Management.Tenant().Load(context.Background(), id)
	require.NoError(t, err, "loading tenant %s via SDK", id)
	return tenant
}

// LoadTenantSettingsViaSDK loads tenant settings directly from the Descope API.
func LoadTenantSettingsViaSDK(t *testing.T, tenantID string) *descope.TenantSettings {
	t.Helper()
	client := newProjectSDKClient(t)
	settings, err := client.Management.Tenant().GetSettings(context.Background(), tenantID)
	require.NoError(t, err, "loading tenant settings for %s via SDK", tenantID)
	return settings
}

// LoadSSOApplicationViaSDK loads an SSO application directly from the Descope API.
func LoadSSOApplicationViaSDK(t *testing.T, id string) *descope.SSOApplication {
	t.Helper()
	client := newProjectSDKClient(t)
	app, err := client.Management.SSOApplication().Load(context.Background(), id)
	require.NoError(t, err, "loading SSO application %s via SDK", id)
	return app
}

// LoadThirdPartyAppViaSDK loads a third-party application directly from the Descope API.
func LoadThirdPartyAppViaSDK(t *testing.T, id string) *descope.ThirdPartyApplication {
	t.Helper()
	client := newProjectSDKClient(t)
	app, err := client.Management.ThirdPartyApplication().LoadApplication(context.Background(), id)
	require.NoError(t, err, "loading third-party application %s via SDK", id)
	return app
}

// LoadFGASchemaViaSDK loads the FGA schema directly from the Descope API.
func LoadFGASchemaViaSDK(t *testing.T) *descope.FGASchema {
	t.Helper()
	client := newProjectSDKClient(t)
	schema, err := client.Management.FGA().LoadSchema(context.Background())
	require.NoError(t, err, "loading FGA schema via SDK")
	return schema
}

// LoadListViaSDK loads a list directly from the Descope API.
func LoadListViaSDK(t *testing.T, id string) *descope.List {
	t.Helper()
	client := newProjectSDKClient(t)
	list, err := client.Management.List().Load(context.Background(), id)
	require.NoError(t, err, "loading list %s via SDK", id)
	return list
}

// LoadAccessKeyViaSDK loads an access key directly from the Descope API.
func LoadAccessKeyViaSDK(t *testing.T, id string) *descope.AccessKeyResponse {
	t.Helper()
	client := newProjectSDKClient(t)
	key, err := client.Management.AccessKey().Load(context.Background(), id)
	require.NoError(t, err, "loading access key %s via SDK", id)
	return key
}

// FindPermissionViaSDK loads all permissions and returns the one matching the given name.
func FindPermissionViaSDK(t *testing.T, name string) *descope.Permission {
	t.Helper()
	client := newProjectSDKClient(t)
	perms, err := client.Management.Permission().LoadAll(context.Background())
	require.NoError(t, err, "loading permissions via SDK")
	for _, p := range perms {
		if p.Name == name {
			return p
		}
	}
	t.Fatalf("permission %q not found via SDK", name)
	return nil
}

// FindRoleViaSDK loads all roles and returns the one matching the given name.
func FindRoleViaSDK(t *testing.T, name string) *descope.Role {
	t.Helper()
	client := newProjectSDKClient(t)
	roles, err := client.Management.Role().LoadAll(context.Background())
	require.NoError(t, err, "loading roles via SDK")
	for _, r := range roles {
		if r.Name == name {
			return r
		}
	}
	t.Fatalf("role %q not found via SDK", name)
	return nil
}

// LoadSSOSettingsViaSDK loads SSO settings for a tenant directly from the Descope API.
func LoadSSOSettingsViaSDK(t *testing.T, tenantID, ssoID string) *descope.SSOTenantSettingsResponse {
	t.Helper()
	client := newProjectSDKClient(t)
	settings, err := client.Management.SSO().LoadSettings(context.Background(), tenantID, ssoID)
	require.NoError(t, err, "loading SSO settings for tenant %s via SDK", tenantID)
	return settings
}

// LoadOutboundAppViaSDK loads an outbound application directly from the Descope API.
func LoadOutboundAppViaSDK(t *testing.T, id string) *descope.OutboundApp {
	t.Helper()
	client := newProjectSDKClient(t)
	app, err := client.Management.OutboundApplication().LoadApplication(context.Background(), id)
	require.NoError(t, err, "loading outbound application %s via SDK", id)
	return app
}

// LoadPasswordSettingsViaSDK loads password settings directly from the Descope API.
// Pass an empty tenantID for project-level settings.
func LoadPasswordSettingsViaSDK(t *testing.T, tenantID string) *descope.PasswordSettings {
	t.Helper()
	client := newProjectSDKClient(t)
	settings, err := client.Management.Password().GetSettings(context.Background(), tenantID)
	require.NoError(t, err, "loading password settings via SDK")
	return settings
}

// --- Company-scoped resource loaders ---

// LoadManagementKeyViaSDK loads a management key directly from the Descope API.
func LoadManagementKeyViaSDK(t *testing.T, id string) *descope.MgmtKey {
	t.Helper()
	client := newCompanySDKClient(t)
	key, err := client.Management.ManagementKey().Get(context.Background(), id)
	require.NoError(t, err, "loading management key %s via SDK", id)
	return key
}

// LoadDescoperViaSDK loads a descoper directly from the Descope API.
func LoadDescoperViaSDK(t *testing.T, id string) *descope.Descoper {
	t.Helper()
	client := newCompanySDKClient(t)
	d, err := client.Management.Descoper().Get(context.Background(), id)
	require.NoError(t, err, "loading descoper %s via SDK", id)
	return d
}

// ProjectExistsViaSDK checks whether a project exists by listing all projects.
func ProjectExistsViaSDK(t *testing.T, projectID string) bool {
	t.Helper()
	client := newCompanySDKClient(t)
	projects, err := client.Management.Project().ListProjects(context.Background())
	require.NoError(t, err, "listing projects via SDK")
	for _, p := range projects {
		if p.ID == projectID {
			return true
		}
	}
	return false
}
