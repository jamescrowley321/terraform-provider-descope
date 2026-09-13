package main

import (
	"context"
	"errors"
	"testing"

	"github.com/descope/go-sdk/descope"
	"github.com/descope/go-sdk/descope/sdk"
	mocksmgmt "github.com/descope/go-sdk/descope/tests/mocks/mgmt"
	"github.com/stretchr/testify/require"
)

// deletion records one Tenant().Delete call so a test can assert on the exact
// set that reached the API, rather than on a count alone.
type deletion struct {
	projectID string
	tenantID  string
	cascade   bool
}

// tenantFixture builds a company whose projects hold the given tenants, and
// returns the company-scoped management mock, a project-client factory over it,
// and the pointer the factory records deletions into.
func tenantFixture(t *testing.T, projects map[string][]*descope.Tenant) (sdk.Management, projectClientFn, *[]deletion) {
	t.Helper()

	listed := make([]*descope.Project, 0, len(projects))
	for id := range projects {
		listed = append(listed, &descope.Project{ID: id, Name: id})
	}
	company := &mocksmgmt.MockManagement{
		MockProject: &mocksmgmt.MockProject{ListProjectsResponse: listed},
	}

	deletions := &[]deletion{}
	factory := func(projectID string) (sdk.Management, error) {
		tenants, ok := projects[projectID]
		require.True(t, ok, "factory asked for a project that was never listed: %s", projectID)
		return &mocksmgmt.MockManagement{
			MockTenant: &mocksmgmt.MockTenant{
				LoadAllResponse: tenants,
				DeleteAssert: func(id string, cascade bool) {
					*deletions = append(*deletions, deletion{projectID: projectID, tenantID: id, cascade: cascade})
				},
			},
		}, nil
	}
	return company, factory, deletions
}

func TestCleanupTenantsDeletesOnlyPrefixedTenantsAcrossEveryProject(t *testing.T) {
	company, factory, deletions := tenantFixture(t, map[string][]*descope.Tenant{
		"P-real": {
			{ID: "T-acme", Name: "Acme Corp"},
			{ID: "T-stray", Name: "testacc-TenantCRUD-03270600-25afd1bb"},
			{ID: "T-globex", Name: "Globex Inc"},
		},
		"P-testacc": {
			{ID: "T-owned", Name: "testacc-TenantResource"},
		},
		"P-empty": {},
	})

	deleted, failed := cleanupTenants(context.Background(), company, factory, false, true)

	require.Equal(t, 2, deleted)
	require.Equal(t, 0, failed)
	require.ElementsMatch(t, []deletion{
		{projectID: "P-real", tenantID: "T-stray", cascade: true},
		{projectID: "P-testacc", tenantID: "T-owned", cascade: true},
	}, *deletions, "only testacc- tenants may be deleted, and every delete must cascade")
}

func TestCleanupTenantsDryRunDeletesNothing(t *testing.T) {
	company, factory, deletions := tenantFixture(t, map[string][]*descope.Tenant{
		"P-real": {
			{ID: "T-acme", Name: "Acme Corp"},
			{ID: "T-stray", Name: "testacc-TenantCRUD-03270600-25afd1bb"},
		},
	})

	deleted, failed := cleanupTenants(context.Background(), company, factory, true, true)

	require.Equal(t, 1, deleted, "dry run still reports what it would delete")
	require.Equal(t, 0, failed)
	require.Empty(t, *deletions, "dry run must not call Tenant().Delete")
}

func TestCleanupTenantsCountsPerProjectFailuresWithoutAbandoningTheSweep(t *testing.T) {
	company := &mocksmgmt.MockManagement{
		MockProject: &mocksmgmt.MockProject{ListProjectsResponse: []*descope.Project{
			{ID: "P-unreachable", Name: "P-unreachable"},
			{ID: "P-listfails", Name: "P-listfails"},
			{ID: "P-delfails", Name: "P-delfails"},
			{ID: "P-ok", Name: "P-ok"},
		}},
	}

	var deletions []deletion
	factory := func(projectID string) (sdk.Management, error) {
		switch projectID {
		case "P-unreachable":
			return nil, errors.New("no client for you")
		case "P-listfails":
			return &mocksmgmt.MockManagement{
				MockTenant: &mocksmgmt.MockTenant{LoadAllError: errors.New("load all failed")},
			}, nil
		case "P-delfails":
			return &mocksmgmt.MockManagement{
				MockTenant: &mocksmgmt.MockTenant{
					LoadAllResponse: []*descope.Tenant{{ID: "T-doomed", Name: "testacc-doomed"}},
					DeleteError:     errors.New("429 slow down"),
				},
			}, nil
		default:
			return &mocksmgmt.MockManagement{
				MockTenant: &mocksmgmt.MockTenant{
					LoadAllResponse: []*descope.Tenant{{ID: "T-last", Name: "testacc-last"}},
					DeleteAssert: func(id string, cascade bool) {
						deletions = append(deletions, deletion{projectID: projectID, tenantID: id, cascade: cascade})
					},
				},
			}, nil
		}
	}

	deleted, failed := cleanupTenants(context.Background(), company, factory, false, true)

	require.Equal(t, 1, deleted)
	require.Equal(t, 3, failed)
	require.Equal(t, []deletion{{projectID: "P-ok", tenantID: "T-last", cascade: true}}, deletions,
		"a project that fails must not stop the projects after it")
}

func TestCleanupTenantsReportsOneFailureWhenProjectsCannotBeListed(t *testing.T) {
	company := &mocksmgmt.MockManagement{
		MockProject: &mocksmgmt.MockProject{ListProjectsError: errors.New("unauthorized")},
	}
	factory := func(string) (sdk.Management, error) {
		t.Fatal("no project client may be built when the project list failed")
		return nil, nil
	}

	deleted, failed := cleanupTenants(context.Background(), company, factory, false, true)

	require.Equal(t, 0, deleted)
	require.Equal(t, 1, failed)
}

func TestRunCleanupDryRunSkipsTheDeleteFunction(t *testing.T) {
	listed := []resource{{name: "testacc-key", id: "K-1"}, {name: "testacc-key-2", id: "K-2"}}
	listFn := func() ([]resource, error) { return listed, nil }
	delFn := func(r resource) error {
		t.Fatalf("dry run must not delete %s (%s)", r.name, r.id)
		return nil
	}

	deleted, failed := runCleanup("access key", listFn, delFn, true)

	require.Equal(t, len(listed), deleted)
	require.Equal(t, 0, failed)
}

func TestCleanupTenantsSkipsNonTestaccProjectsUnlessAsked(t *testing.T) {
	projects := map[string][]*descope.Tenant{
		"identity-stack": {
			{ID: "T-acme", Name: "Acme Corp"},
			// A production tenant someone named carelessly. Only -all-projects
			// may reach it; by default a name collision costs nothing.
			{ID: "T-collision", Name: "testacc-demo"},
		},
		"testacc-project": {
			{ID: "T-owned", Name: "testacc-TenantResource"},
		},
	}

	company, factory, deletions := tenantFixture(t, projects)
	deleted, failed := cleanupTenants(context.Background(), company, factory, false, false)

	require.Equal(t, 1, deleted)
	require.Equal(t, 0, failed)
	require.Equal(t, []deletion{{projectID: "testacc-project", tenantID: "T-owned", cascade: true}}, *deletions,
		"the default sweep must not reach into a project that is not itself testacc-")

	company, factory, deletions = tenantFixture(t, projects)
	deleted, failed = cleanupTenants(context.Background(), company, factory, false, true)

	require.Equal(t, 2, deleted)
	require.Equal(t, 0, failed)
	require.ElementsMatch(t, []deletion{
		{projectID: "identity-stack", tenantID: "T-collision", cascade: true},
		{projectID: "testacc-project", tenantID: "T-owned", cascade: true},
	}, *deletions, "-all-projects is what reaches strays in real projects")
}
