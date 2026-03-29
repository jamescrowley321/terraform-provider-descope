//go:build integration || fork

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	descopeclient "github.com/descope/go-sdk/descope/client"
	"github.com/hashicorp/go-uuid"
	"github.com/stretchr/testify/require"
)

const (
	commandMaxRetries = 2
	commandRetryWait  = 5 * time.Second

	// errLastProject is the Descope API error message returned when attempting
	// to delete the only remaining project in an account ([E073008]).
	errLastProject = "Cannot delete last project"
)

var (
	buildOnce     sync.Once
	binaryPath    string
	buildErr      error
	terraformPath string
	goPath        string
)

// Harness manages a terraform workspace for integration testing.
// It builds the provider binary, configures dev_overrides, and provides
// helpers for running terraform commands and inspecting state.
type Harness struct {
	t          *testing.T
	workDir    string
	env        []string
	lastVars   []string
	projectIDs []string // project IDs created during the test
}

// NewHarness builds the provider binary (once), creates a temp workspace,
// and configures dev_overrides so terraform uses the local build.
func NewHarness(t *testing.T) *Harness {
	t.Helper()
	requireEnvVars(t)
	requireTerraform(t)
	buildProvider(t)

	workDir := t.TempDir()

	// Write .terraformrc with dev_overrides pointing to the binary directory
	binDir := filepath.Dir(binaryPath)
	terraformrc := filepath.Join(workDir, ".terraformrc")
	content := fmt.Sprintf("provider_installation {\n  dev_overrides {\n    \"jamescrowley321/descope\" = %q\n  }\n  direct {}\n}\n", binDir)
	require.NoError(t, os.WriteFile(terraformrc, []byte(content), 0600))

	env := append(os.Environ(), "TF_CLI_CONFIG_FILE="+terraformrc)

	h := &Harness{
		t:       t,
		workDir: workDir,
		env:     env,
	}

	// Copy shared provider configuration
	h.copyTestdata("provider.tf", "provider.tf")

	// Best-effort cleanup: destroy any remaining resources when the test ends
	t.Cleanup(func() {
		args := []string{"destroy", "-auto-approve", "-no-color", "-input=false"}
		args = append(args, varArgs(h.lastVars)...)
		cmd := exec.Command(terraformPath, args...)
		cmd.Dir = h.workDir
		cmd.Env = h.env
		_ = cmd.Run()
	})

	return h
}

// LoadFixture copies a fixture from testdata/<path> to main.tf in the workspace,
// replacing any previous fixture.
func (h *Harness) LoadFixture(path string) {
	h.t.Helper()
	h.copyTestdata(path, "main.tf")
}

// Apply runs terraform apply -auto-approve and returns stdout.
// It retries up to commandMaxRetries times on failure to handle transient
// API errors during resource creation.
func (h *Harness) Apply(vars ...string) string {
	h.t.Helper()
	h.lastVars = vars
	args := []string{"apply", "-auto-approve", "-no-color", "-input=false"}
	args = append(args, varArgs(vars)...)
	return h.terraformRetry(commandMaxRetries, args...)
}

// TryApply runs terraform apply and returns the combined output and any error
// instead of calling t.Fatal. This allows callers to inspect the error and
// decide how to proceed (e.g., skip on license errors).
func (h *Harness) TryApply(vars ...string) (string, error) {
	h.t.Helper()
	h.lastVars = vars
	args := []string{"apply", "-auto-approve", "-no-color", "-input=false"}
	args = append(args, varArgs(vars)...)
	cmd := exec.Command(terraformPath, args...)
	cmd.Dir = h.workDir
	cmd.Env = h.env
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String() + "\n" + stderr.String(), err
	}
	return stdout.String(), nil
}

// Plan runs terraform plan and returns stdout.
func (h *Harness) Plan(vars ...string) string {
	h.t.Helper()
	h.lastVars = vars
	args := []string{"plan", "-no-color", "-input=false"}
	args = append(args, varArgs(vars)...)
	return h.terraform(args...)
}

// Destroy runs terraform destroy -auto-approve and returns stdout.
// It retries up to commandMaxRetries times on failure to handle transient
// API errors. If the harness created any projects, it waits for Descope
// to finish deleting them asynchronously before returning.
//
// If the Descope API rejects deletion because the resource is the last
// project in the account ([E073008]), destroy logs a warning and clears
// the terraform state instead of failing the test.
func (h *Harness) Destroy(vars ...string) string {
	h.t.Helper()
	h.lastVars = vars
	args := []string{"destroy", "-auto-approve", "-no-color", "-input=false"}
	args = append(args, varArgs(vars)...)
	out, err := h.tryTerraformRetry(commandMaxRetries, args...)
	if err != nil {
		if strings.Contains(out, errLastProject) {
			h.t.Logf("warning: cannot delete last project in Descope account, clearing terraform state: %v", err)
			h.clearState()
			h.projectIDs = nil
			return out
		}
		h.t.Fatalf("terraform destroy failed:\n%s\nerror: %v", out, err)
	}
	if len(h.projectIDs) > 0 {
		h.waitForProjectDeletion()
		h.projectIDs = nil
	}
	return out
}

// Import runs terraform import for the given resource address and ID.
func (h *Harness) Import(address, id string, vars ...string) string {
	h.t.Helper()
	h.lastVars = vars
	args := []string{"import", "-no-color", "-input=false"}
	args = append(args, varArgs(vars)...)
	args = append(args, address, id)
	return h.terraform(args...)
}

// StateRM removes a resource from the terraform state without destroying it.
func (h *Harness) StateRM(address string) string {
	h.t.Helper()
	return h.terraform("state", "rm", "-no-color", address)
}

// Output returns the value of a terraform output variable.
func (h *Harness) Output(name string) string {
	h.t.Helper()
	return strings.TrimSpace(h.terraform("output", "-no-color", "-raw", name))
}

// StateResource returns a resource's attribute values from the terraform state.
func (h *Harness) StateResource(address string) map[string]any {
	h.t.Helper()
	out := h.terraform("show", "-json", "-no-color")
	var state struct {
		Values struct {
			RootModule struct {
				Resources []struct {
					Address string         `json:"address"`
					Values  map[string]any `json:"values"`
				} `json:"resources"`
			} `json:"root_module"`
		} `json:"values"`
	}
	require.NoError(h.t, json.Unmarshal([]byte(out), &state))
	for _, r := range state.Values.RootModule.Resources {
		if r.Address == address {
			return r.Values
		}
	}
	h.t.Fatalf("resource %s not found in state", address)
	return nil
}

// HasState returns true if the terraform state file contains resources.
func (h *Harness) HasState() bool {
	h.t.Helper()
	data, err := os.ReadFile(filepath.Join(h.workDir, "terraform.tfstate"))
	if err != nil {
		return false
	}
	var state struct {
		Resources []any `json:"resources"`
	}
	if err := json.Unmarshal(data, &state); err != nil {
		return false
	}
	return len(state.Resources) > 0
}

// TerraformExpectFailure runs a terraform command and expects it to fail.
// Returns the combined stdout+stderr output.
func (h *Harness) TerraformExpectFailure(args ...string) string {
	h.t.Helper()
	cmd := exec.Command(terraformPath, args...)
	cmd.Dir = h.workDir
	cmd.Env = h.env
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	require.Error(h.t, err, "expected terraform to fail but it succeeded")
	return buf.String()
}

// ApplyFixture loads a fixture and runs apply, returning the state attributes for the given resource address.
func (h *Harness) ApplyFixture(fixture, address string, vars ...string) map[string]any {
	h.t.Helper()
	h.LoadFixture(fixture)
	h.Apply(vars...)
	attrs := h.StateResource(address)
	// Delete default SSO applications from newly created projects to stay
	// within the free plan's SSO application limit, and track the project ID
	// so we can wait for async deletion to complete after destroy.
	if strings.HasPrefix(address, "descope_project.") {
		if id, ok := attrs["id"].(string); ok && id != "" {
			deleteDefaultSSOApps(h.t, id)
			h.projectIDs = append(h.projectIDs, id)
		}
	}
	return attrs
}

// ReimportResource removes a resource from state, loads a fixture, and imports it back.
func (h *Harness) ReimportResource(fixture, address, id string, vars ...string) map[string]any {
	h.t.Helper()
	h.StateRM(address)
	h.LoadFixture(fixture)
	h.Import(address, id, vars...)
	return h.StateResource(address)
}

// StringAttr returns the string representation of a state attribute value.
func StringAttr(attrs map[string]any, key string) string {
	return fmt.Sprintf("%v", attrs[key])
}

// RequireMap extracts a map attribute and fails if it is not a map.
func RequireMap(t *testing.T, attrs map[string]any, key string) map[string]any {
	t.Helper()
	m, ok := attrs[key].(map[string]any)
	require.True(t, ok, "%s should be a map", key)
	return m
}

// RequireList extracts a list attribute and fails if it is not a list.
func RequireList(t *testing.T, attrs map[string]any, key string) []any {
	t.Helper()
	l, ok := attrs[key].([]any)
	require.True(t, ok, "%s should be a list", key)
	return l
}

// RequireListLen extracts a list attribute and asserts its length.
func RequireListLen(t *testing.T, attrs map[string]any, key string, length int) []any {
	t.Helper()
	l := RequireList(t, attrs, key)
	require.Len(t, l, length, "%s should have %d elements", key, length)
	return l
}

// GenerateName creates a unique resource name for testing.
func GenerateName(t *testing.T) string {
	t.Helper()
	test := strings.TrimPrefix(t.Name(), "Test")
	ts := time.Now().Format("01021504")
	rand, err := uuid.GenerateUUID()
	require.NoError(t, err)
	return fmt.Sprintf("testacc-%s-%s-%s", test, ts, rand[len(rand)-8:])
}

// waitForProjectDeletion polls the Descope API until all projects created by
// this harness have been fully deleted. Descope processes project deletion
// asynchronously, so without this wait, sequential tests that create projects
// can fail when the previous project hasn't been fully removed yet.
func (h *Harness) waitForProjectDeletion() {
	h.t.Helper()
	client, err := descopeclient.NewWithConfig(&descopeclient.Config{
		ManagementKey:       os.Getenv("DESCOPE_MANAGEMENT_KEY"),
		DescopeBaseURL:      os.Getenv("DESCOPE_BASE_URL"),
		AllowEmptyProjectID: true,
	})
	if err != nil {
		h.t.Logf("warning: could not create client to wait for project deletion: %v", err)
		return
	}
	ctx := context.Background()
	pending := make(map[string]bool)
	for _, id := range h.projectIDs {
		pending[id] = true
	}
	for attempt := 0; attempt < 30; attempt++ {
		projects, err := client.Management.Project().ListProjects(ctx)
		if err != nil {
			h.t.Logf("warning: failed to list projects while waiting for deletion: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}
		allGone := true
		for _, p := range projects {
			if pending[p.ID] {
				allGone = false
				break
			}
		}
		if allGone {
			return
		}
		time.Sleep(2 * time.Second)
	}
	h.t.Logf("warning: timed out waiting for project deletion after 60s")
}

// deleteDefaultSSOApps removes all SSO applications from a project using the
// Descope SDK. This prevents test projects from consuming SSO application quota
// on the free plan, where each new project gets a default OIDC application.
func deleteDefaultSSOApps(t *testing.T, projectID string) {
	t.Helper()
	client, err := descopeclient.NewWithConfig(&descopeclient.Config{
		ManagementKey:  os.Getenv("DESCOPE_MANAGEMENT_KEY"),
		DescopeBaseURL: os.Getenv("DESCOPE_BASE_URL"),
		ProjectID:      projectID,
	})
	if err != nil {
		t.Logf("warning: failed to create Descope client for SSO app cleanup: %v", err)
		return
	}
	ctx := context.Background()
	apps, err := client.Management.SSOApplication().LoadAll(ctx)
	if err != nil {
		t.Logf("warning: failed to list SSO applications in project %s: %v", projectID, err)
		return
	}
	for _, app := range apps {
		if strings.HasPrefix(app.ID, "descope-default-") {
			continue
		}
		if err := client.Management.SSOApplication().Delete(ctx, app.ID); err != nil {
			t.Logf("warning: failed to delete SSO application %s (%s): %v", app.Name, app.ID, err)
		}
	}
}

func (h *Harness) terraform(args ...string) string {
	h.t.Helper()
	cmd := exec.Command(terraformPath, args...)
	cmd.Dir = h.workDir
	cmd.Env = h.env
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		h.t.Fatalf("terraform %s failed:\nstdout:\n%s\nstderr:\n%s\nerror: %v",
			strings.Join(args, " "), stdout.String(), stderr.String(), err)
	}
	return stdout.String()
}

// terraformRetry runs a terraform command with retries on failure, calling
// t.Fatal if all attempts are exhausted. It is used for operations like apply
// where transient API errors may cause the command to fail.
func (h *Harness) terraformRetry(maxRetries int, args ...string) string {
	h.t.Helper()
	out, err := h.tryTerraformRetry(maxRetries, args...)
	if err != nil {
		h.t.Fatalf("terraform %s failed:\n%s\nerror: %v", strings.Join(args, " "), out, err)
	}
	return out
}

// tryTerraformRetry runs a terraform command with retries, returning the
// combined output and any error instead of calling t.Fatal.
func (h *Harness) tryTerraformRetry(maxRetries int, args ...string) (string, error) {
	h.t.Helper()
	var lastErr error
	var lastOutput string
	for attempt := range maxRetries + 1 {
		cmd := exec.Command(terraformPath, args...)
		cmd.Dir = h.workDir
		cmd.Env = h.env
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			lastErr = err
			lastOutput = stdout.String() + "\n" + stderr.String()
			if attempt < maxRetries {
				h.t.Logf("terraform %s failed (attempt %d/%d), retrying in %v:\nstderr: %s",
					strings.Join(args, " "), attempt+1, maxRetries+1, commandRetryWait, stderr.String())
				time.Sleep(commandRetryWait)
				continue
			}
			return lastOutput, lastErr
		}
		return stdout.String(), nil
	}
	return lastOutput, lastErr
}

// clearState removes terraform state and backup files so that HasState returns false.
func (h *Harness) clearState() {
	h.t.Helper()
	for _, name := range []string{"terraform.tfstate", "terraform.tfstate.backup"} {
		p := filepath.Join(h.workDir, name)
		if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
			h.t.Logf("warning: failed to remove %s: %v", name, err)
		}
	}
}

func (h *Harness) copyTestdata(src, dst string) {
	h.t.Helper()
	data, err := os.ReadFile(testdataPath(src))
	require.NoError(h.t, err, "reading testdata/%s", src)
	require.NoError(h.t, os.WriteFile(filepath.Join(h.workDir, dst), data, 0600))
}

func requireEnvVars(t *testing.T) {
	t.Helper()
	for _, env := range []string{"DESCOPE_MANAGEMENT_KEY", "DESCOPE_BASE_URL"} {
		require.NotEmpty(t, os.Getenv(env), "required environment variable %s is not set", env)
	}
}

func requireTerraform(t *testing.T) {
	t.Helper()
	var err error
	terraformPath, err = exec.LookPath("terraform")
	require.NoError(t, err, "terraform CLI must be installed and in PATH")
	goPath, err = exec.LookPath("go")
	require.NoError(t, err, "go CLI must be installed and in PATH")
}

func buildProvider(t *testing.T) {
	t.Helper()
	buildOnce.Do(func() {
		tmpDir, err := os.MkdirTemp("", "terraform-provider-descope-*")
		if err != nil {
			buildErr = err
			return
		}
		name := "terraform-provider-descope"
		if runtime.GOOS == "windows" {
			name += ".exe"
		}
		binaryPath = filepath.Join(tmpDir, name)
		cmd := exec.Command(goPath, "build", "-o", binaryPath, ".")
		cmd.Dir = projectRoot()
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			buildErr = fmt.Errorf("build provider: %s: %w", stderr.String(), err)
		}
	})
	require.NoError(t, buildErr)
}

func projectRoot() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..")
}

func testdataPath(name string) string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "testdata", name)
}

func varArgs(vars []string) []string {
	args := make([]string, 0, len(vars)*2)
	for _, v := range vars {
		args = append(args, "-var", v)
	}
	return args
}
