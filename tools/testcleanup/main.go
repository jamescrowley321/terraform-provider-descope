// Command testcleanup deletes all Descope test resources whose names start with "testacc-".
//
// It cleans up projects, access keys, management keys, and descopers.
//
// It requires DESCOPE_MANAGEMENT_KEY and DESCOPE_BASE_URL environment variables.
// Usage: source .env && go run ./tools/testcleanup
package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/descope/go-sdk/descope"
	descopeclient "github.com/descope/go-sdk/descope/client"
	"github.com/descope/go-sdk/descope/sdk"
)

const testPrefix = "testacc-"

func main() {
	managementKey := os.Getenv("DESCOPE_MANAGEMENT_KEY")
	baseURL := os.Getenv("DESCOPE_BASE_URL")

	if managementKey == "" || baseURL == "" {
		fmt.Fprintln(os.Stderr, "DESCOPE_MANAGEMENT_KEY and DESCOPE_BASE_URL must be set")
		os.Exit(1)
	}

	ctx := context.Background()

	client, err := descopeclient.NewWithConfig(&descopeclient.Config{
		ManagementKey:       managementKey,
		DescopeBaseURL:      baseURL,
		AllowEmptyProjectID: true,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create client: %v\n", err)
		os.Exit(1)
	}

	mgmt := client.Management
	var totalDeleted, totalFailed int

	d, f := cleanupAccessKeys(ctx, mgmt)
	totalDeleted += d
	totalFailed += f

	d, f = cleanupManagementKeys(ctx, mgmt)
	totalDeleted += d
	totalFailed += f

	d, f = cleanupDescopers(ctx, mgmt)
	totalDeleted += d
	totalFailed += f

	d, f = cleanupProjects(ctx, mgmt, managementKey, baseURL)
	totalDeleted += d
	totalFailed += f

	fmt.Printf("\ntotal: %d deleted, %d failed\n", totalDeleted, totalFailed)
	if totalFailed > 0 {
		os.Exit(1)
	}
}

func cleanupAccessKeys(ctx context.Context, mgmt sdk.Management) (deleted, failed int) {
	keys, err := mgmt.AccessKey().SearchAll(ctx, &descope.AccessKeysSearchOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to search access keys: %v\n", err)
		return 0, 1
	}

	for _, k := range keys {
		if !strings.HasPrefix(k.Name, testPrefix) {
			continue
		}
		fmt.Printf("deleting access key %s (%s)...\n", k.Name, k.ID)
		if err := mgmt.AccessKey().Delete(ctx, k.ID); err != nil {
			fmt.Fprintf(os.Stderr, "  failed: %v\n", err)
			failed++
			continue
		}
		deleted++
	}

	if deleted > 0 || failed > 0 {
		fmt.Printf("access keys: %d deleted, %d failed\n", deleted, failed)
	}
	return
}

func cleanupManagementKeys(ctx context.Context, mgmt sdk.Management) (deleted, failed int) {
	keys, err := mgmt.ManagementKey().Search(ctx, &descope.MgmtKeySearchOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to search management keys: %v\n", err)
		return 0, 1
	}

	for _, k := range keys {
		if !strings.HasPrefix(k.Name, testPrefix) {
			continue
		}
		fmt.Printf("deleting management key %s (%s)...\n", k.Name, k.ID)
		if _, err := mgmt.ManagementKey().Delete(ctx, []string{k.ID}); err != nil {
			fmt.Fprintf(os.Stderr, "  failed: %v\n", err)
			failed++
			continue
		}
		deleted++
	}

	if deleted > 0 || failed > 0 {
		fmt.Printf("management keys: %d deleted, %d failed\n", deleted, failed)
	}
	return
}

func cleanupDescopers(ctx context.Context, mgmt sdk.Management) (deleted, failed int) {
	descopers, _, err := mgmt.Descoper().List(ctx, &descope.DescoperLoadOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to list descopers: %v\n", err)
		return 0, 1
	}

	for _, d := range descopers {
		name := ""
		if d.Attributes != nil {
			name = d.Attributes.DisplayName
		}
		if !strings.HasPrefix(name, testPrefix) {
			continue
		}
		fmt.Printf("deleting descoper %s (%s)...\n", name, d.ID)
		if err := mgmt.Descoper().Delete(ctx, d.ID); err != nil {
			fmt.Fprintf(os.Stderr, "  failed: %v\n", err)
			failed++
			continue
		}
		deleted++
	}

	if deleted > 0 || failed > 0 {
		fmt.Printf("descopers: %d deleted, %d failed\n", deleted, failed)
	}
	return
}

func cleanupProjects(ctx context.Context, mgmt sdk.Management, managementKey, baseURL string) (deleted, failed int) {
	projects, err := mgmt.Project().ListProjects(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to list projects: %v\n", err)
		return 0, 1
	}

	for _, p := range projects {
		if !strings.HasPrefix(p.Name, testPrefix) {
			continue
		}
		fmt.Printf("deleting project %s (%s)...\n", p.Name, p.ID)

		// Project.Delete requires a project-scoped client
		projectClient, err := descopeclient.NewWithConfig(&descopeclient.Config{
			ManagementKey:  managementKey,
			DescopeBaseURL: baseURL,
			ProjectID:      p.ID,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "  failed to create project client: %v\n", err)
			failed++
			continue
		}

		if err := projectClient.Management.Project().Delete(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "  failed: %v\n", err)
			failed++
			continue
		}
		deleted++
	}

	if deleted > 0 || failed > 0 {
		fmt.Printf("projects: %d deleted, %d failed\n", deleted, failed)
	}
	return
}
