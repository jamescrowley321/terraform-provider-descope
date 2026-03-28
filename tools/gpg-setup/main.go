// tools/gpg-setup generates a GPG key pair for Terraform Registry provider signing,
// uploads the private key and passphrase to GitHub repository secrets, and
// registers the public key with the Terraform Registry.
//
// Prerequisites:
//   - gpg CLI installed
//   - gh CLI installed and authenticated
//   - HCP_API_TOKEN environment variable set (from https://portal.cloud.hashicorp.com)
//
// Usage:
//
//	go run ./tools/gpg-setup \
//	  -repo jamescrowley321/terraform-provider-descope \
//	  -namespace jamescrowley321 \
//	  -name "James Crowley" \
//	  -email jamescrowley151@gmail.com
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	repo := flag.String("repo", "jamescrowley321/terraform-provider-descope", "GitHub owner/repo")
	namespace := flag.String("namespace", "jamescrowley321", "Terraform Registry namespace")
	name := flag.String("name", "", "GPG key real name (required)")
	email := flag.String("email", "", "GPG key email (required)")
	passphrase := flag.String("passphrase", "", "GPG key passphrase (auto-generated if empty)")
	skipGitHub := flag.Bool("skip-github", false, "Skip uploading secrets to GitHub")
	skipRegistry := flag.Bool("skip-registry", false, "Skip uploading public key to Terraform Registry")
	flag.Parse()

	if *name == "" || *email == "" {
		fmt.Fprintln(os.Stderr, "error: -name and -email are required")
		flag.Usage()
		os.Exit(1)
	}

	requireCmd("gpg")
	if !*skipGitHub {
		requireCmd("gh")
	}

	// Generate passphrase if not provided
	if *passphrase == "" {
		p, err := generatePassphrase()
		if err != nil {
			fatalf("generating passphrase: %v", err)
		}
		*passphrase = p
	}

	fmt.Println("==> Generating GPG key pair...")
	fingerprint, err := generateGPGKey(*name, *email, *passphrase)
	if err != nil {
		fatalf("generating GPG key: %v", err)
	}
	fmt.Printf("    Fingerprint: %s\n", fingerprint)

	fmt.Println("==> Exporting private key...")
	privateKey, err := exportPrivateKey(fingerprint, *passphrase)
	if err != nil {
		fatalf("exporting private key: %v", err)
	}

	fmt.Println("==> Exporting public key...")
	publicKey, err := exportPublicKey(fingerprint)
	if err != nil {
		fatalf("exporting public key: %v", err)
	}

	// Upload to GitHub
	if !*skipGitHub {
		fmt.Printf("==> Uploading secrets to GitHub repo %s...\n", *repo)
		if err := setGitHubSecret(*repo, "GPG_PRIVATE_KEY", privateKey); err != nil {
			fatalf("setting GPG_PRIVATE_KEY secret: %v", err)
		}
		fmt.Println("    Set GPG_PRIVATE_KEY")

		if err := setGitHubSecret(*repo, "GPG_PASSPHRASE", *passphrase); err != nil {
			fatalf("setting GPG_PASSPHRASE secret: %v", err)
		}
		fmt.Println("    Set GPG_PASSPHRASE")
	}

	// Upload to Terraform Registry
	if !*skipRegistry {
		token := os.Getenv("HCP_API_TOKEN")
		if token == "" {
			fmt.Println("")
			fmt.Println("==> Skipping Terraform Registry upload (HCP_API_TOKEN not set)")
			fmt.Println("    To upload later, set HCP_API_TOKEN and re-run with -skip-github")
			fmt.Println("    Or manually add the public key at: https://registry.terraform.io/settings/gpg-keys")
		} else {
			fmt.Println("==> Registering public key with Terraform Registry...")
			if err := uploadToRegistry(token, *namespace, publicKey); err != nil {
				fatalf("uploading to Terraform Registry: %v", err)
			}
			fmt.Println("    Public key registered")
		}
	}

	// Write passphrase to a temp file with restricted permissions so it
	// never appears in terminal scrollback or CI logs.
	ppFile := filepath.Join(os.TempDir(), "gpg-passphrase.txt")
	if err := os.WriteFile(ppFile, []byte(*passphrase), 0600); err != nil {
		fatalf("writing passphrase file: %v", err)
	}

	fmt.Println("")
	fmt.Println("=== Setup Complete ===")
	fmt.Printf("Fingerprint:  %s\n", fingerprint)
	fmt.Printf("Passphrase:   written to %s (delete after saving)\n", ppFile)
	fmt.Println("")
	fmt.Println("Public key (for manual Terraform Registry upload if needed):")
	fmt.Println(publicKey)
}

// generateGPGKey creates a GPG key pair and returns the fingerprint.
func generateGPGKey(name, email, passphrase string) (string, error) {
	params := fmt.Sprintf(`%%no-protection
Key-Type: RSA
Key-Length: 4096
Subkey-Type: RSA
Subkey-Length: 4096
Name-Real: %s
Name-Email: %s
Expire-Date: 0
Passphrase: %s
%%commit
`, name, email, passphrase)

	cmd := exec.Command("gpg", "--batch", "--gen-key") //#nosec G204 -- gpg args are hardcoded
	cmd.Stdin = strings.NewReader(params)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%s: %w", string(out), err)
	}

	// Get the fingerprint of the most recently created key for this email
	out, err := exec.Command("gpg", "--list-keys", "--with-colons", email).Output() //#nosec G204 -- email is from CLI flag
	if err != nil {
		return "", fmt.Errorf("listing keys: %w", err)
	}
	for _, line := range strings.Split(string(out), "\n") {
		if strings.HasPrefix(line, "fpr:") {
			parts := strings.Split(line, ":")
			if len(parts) >= 10 {
				return parts[9], nil
			}
		}
	}
	return "", fmt.Errorf("could not find fingerprint for %s", email)
}

// exportPrivateKey exports the ASCII-armored private key.
func exportPrivateKey(fingerprint, passphrase string) (string, error) {
	cmd := exec.Command("gpg", "--batch", "--yes", "--pinentry-mode", "loopback", //#nosec G204 -- gpg args from internal state
		"--passphrase", passphrase, "--armor", "--export-secret-keys", fingerprint)
	out, err := cmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("%s: %w", string(ee.Stderr), err)
		}
		return "", err
	}
	return string(out), nil
}

// exportPublicKey exports the ASCII-armored public key.
func exportPublicKey(fingerprint string) (string, error) {
	out, err := exec.Command("gpg", "--armor", "--export", fingerprint).Output() //#nosec G204 -- fingerprint from internal state
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// setGitHubSecret sets a repository secret using the gh CLI.
func setGitHubSecret(repo, name, value string) error {
	cmd := exec.Command("gh", "secret", "set", name, "--repo", repo) //#nosec G204 -- gh args from CLI flags
	cmd.Stdin = strings.NewReader(value)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%s: %w", string(out), err)
	}
	return nil
}

// uploadToRegistry registers the GPG public key with the Terraform Registry.
func uploadToRegistry(token, namespace, publicKey string) error {
	payload := map[string]any{
		"data": map[string]any{
			"type": "gpg-keys",
			"attributes": map[string]any{
				"namespace":   namespace,
				"ascii-armor": publicKey,
			},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, "https://app.terraform.io/api/registry/private/v2/gpg-keys", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/vnd.api+json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}
	return nil
}

// generatePassphrase creates a random passphrase using /dev/urandom.
func generatePassphrase() (string, error) {
	out, err := exec.Command("openssl", "rand", "-base64", "32").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// requireCmd checks that a command is available on PATH.
func requireCmd(name string) {
	if _, err := exec.LookPath(name); err != nil {
		fatalf("%s is required but not found in PATH", name)
	}
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
	os.Exit(1)
}
