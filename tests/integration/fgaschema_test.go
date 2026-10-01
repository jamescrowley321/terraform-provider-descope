//go:build integration || fork

package integration

import (
	"context"
	"testing"

	"github.com/descope/go-sdk/descope"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFGASchemaCRUD(t *testing.T) {
	// Provision a project to own the schema for the duration of the test. Its
	// harness is created first so its cleanup runs last: t.Cleanup is LIFO, so
	// the schema resource is torn down before the project that holds it.
	host := NewHarness(t)
	hostName := GenerateName(t)
	projectAttrs := host.ApplyFixture("project/create.tf", "descope_project.test", "name="+hostName)
	projectID := StringAttr(projectAttrs, "id")
	require.NotEmpty(t, projectID, "throwaway project must have an id")

	h := NewHarnessInProject(t, projectID)
	name := GenerateName(t)
	nameVar := "name=" + name
	address := "descope_fga_schema.test"

	// Create
	attrs := h.ApplyFixture("fgaschema/create.tf", address, nameVar)

	id := StringAttr(attrs, "id")
	require.NotEmpty(t, id)

	schema := StringAttr(attrs, "schema")
	assert.Contains(t, schema, "document")
	assert.Contains(t, schema, "owner")

	// Verify via SDK, scoped to the throwaway project. Reading the shared
	// project here would assert against a schema this test never wrote.
	sdkSchema := LoadFGASchemaViaSDKInProject(t, projectID)
	assert.Contains(t, sdkSchema.Schema, "document")
	assert.Contains(t, sdkSchema.Schema, "owner")

	require.Equal(t, false, h.StateResource("data.descope_fga_check.owner")["allowed"])
	management := newSDKClientWithProject(t, projectID).Management
	relations := []*descope.FGARelation{{Resource: "document-1", ResourceType: "document", Relation: "owner", Target: "user-1", TargetType: "user"}}
	require.NoError(t, management.FGA().CreateRelations(context.Background(), relations))
	h.Apply(nameVar)
	require.Equal(t, true, h.StateResource("data.descope_fga_check.owner")["allowed"])
	require.NoError(t, management.FGA().DeleteRelations(context.Background(), relations))
	h.Apply(nameVar)
	require.Equal(t, false, h.StateResource("data.descope_fga_check.owner")["allowed"])

	// Destroy
	h.Destroy(nameVar)
	assert.False(t, h.HasState())
}
