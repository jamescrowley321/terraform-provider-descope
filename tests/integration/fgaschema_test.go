//go:build integration || fork

package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFGASchemaCRUD runs against a throwaway project of its own, never the
// shared DESCOPE_PROJECT_ID.
//
// descope_fga_schema is a PROJECT-LEVEL SINGLETON: the resource has no
// project_id attribute and always targets the provider's project, and saving a
// schema REPLACES that project's schema rather than adding to it. Applying this
// fixture against the shared project therefore deleted every relation and
// permission the project's real schema declared, leaving only the two lines the
// fixture happens to name — silently breaking anything else that authorizes
// against that project until someone restored it by hand.
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

	// Destroy
	h.Destroy(nameVar)
	assert.False(t, h.HasState())
}
