//go:build integration

package integration

import (
	"bytes"
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

	"github.com/hashicorp/go-uuid"
	"github.com/stretchr/testify/require"
)

var (
	buildOnce  sync.Once
	binaryPath string
	buildErr   error
)

// Harness manages a terraform workspace for integration testing.
// It builds the provider binary, configures dev_overrides, and provides
// helpers for running terraform commands and inspecting state.
type Harness struct {
	t        *testing.T
	workDir  string
	env      []string
	lastVars []string
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
	content := fmt.Sprintf("provider_installation {\n  dev_overrides {\n    \"descope/descope\" = %q\n  }\n  direct {}\n}\n", binDir)
	require.NoError(t, os.WriteFile(terraformrc, []byte(content), 0644))

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
		cmd := exec.Command("terraform", args...)
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
func (h *Harness) Apply(vars ...string) string {
	h.t.Helper()
	h.lastVars = vars
	args := []string{"apply", "-auto-approve", "-no-color", "-input=false"}
	args = append(args, varArgs(vars)...)
	return h.terraform(args...)
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
func (h *Harness) Destroy(vars ...string) string {
	h.t.Helper()
	h.lastVars = vars
	args := []string{"destroy", "-auto-approve", "-no-color", "-input=false"}
	args = append(args, varArgs(vars)...)
	return h.terraform(args...)
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
	cmd := exec.Command("terraform", args...)
	cmd.Dir = h.workDir
	cmd.Env = h.env
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	require.Error(h.t, err, "expected terraform to fail but it succeeded")
	return buf.String()
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

func (h *Harness) terraform(args ...string) string {
	h.t.Helper()
	cmd := exec.Command("terraform", args...)
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

func (h *Harness) copyTestdata(src, dst string) {
	h.t.Helper()
	data, err := os.ReadFile(testdataPath(src))
	require.NoError(h.t, err, "reading testdata/%s", src)
	require.NoError(h.t, os.WriteFile(filepath.Join(h.workDir, dst), data, 0644))
}

func requireEnvVars(t *testing.T) {
	t.Helper()
	for _, env := range []string{"DESCOPE_MANAGEMENT_KEY", "DESCOPE_BASE_URL"} {
		require.NotEmpty(t, os.Getenv(env), "required environment variable %s is not set", env)
	}
}

func requireTerraform(t *testing.T) {
	t.Helper()
	_, err := exec.LookPath("terraform")
	require.NoError(t, err, "terraform CLI must be installed and in PATH")
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
		cmd := exec.Command("go", "build", "-o", binaryPath, ".")
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
