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
)

const testPrefix = "testacc-"

type resource struct {
	name string
	id   string
}

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

	cleanups := []struct {
		name   string
		listFn func() ([]resource, error)
		delFn  func(resource) error
	}{
		{
			name: "access key",
			listFn: func() ([]resource, error) {
				keys, err := mgmt.AccessKey().SearchAll(ctx, &descope.AccessKeysSearchOptions{})
				if err != nil {
					return nil, err
				}
				return collectResources(keys, func(k *descope.AccessKeyResponse) (string, string) { return k.Name, k.ID }), nil
			},
			delFn: func(r resource) error { return mgmt.AccessKey().Delete(ctx, r.id) },
		},
		{
			name: "management key",
			listFn: func() ([]resource, error) {
				keys, err := mgmt.ManagementKey().Search(ctx, &descope.MgmtKeySearchOptions{})
				if err != nil {
					return nil, err
				}
				return collectResources(keys, func(k *descope.MgmtKey) (string, string) { return k.Name, k.ID }), nil
			},
			delFn: func(r resource) error { _, err := mgmt.ManagementKey().Delete(ctx, []string{r.id}); return err },
		},
		{
			name: "descoper",
			listFn: func() ([]resource, error) {
				descopers, _, err := mgmt.Descoper().List(ctx, &descope.DescoperLoadOptions{})
				if err != nil {
					return nil, err
				}
				return collectResources(descopers, func(d *descope.Descoper) (string, string) {
					name := ""
					if d.Attributes != nil {
						name = d.Attributes.DisplayName
					}
					return name, d.ID
				}), nil
			},
			delFn: func(r resource) error { return mgmt.Descoper().Delete(ctx, r.id) },
		},
		{
			name: "project",
			listFn: func() ([]resource, error) {
				projects, err := mgmt.Project().ListProjects(ctx)
				if err != nil {
					return nil, err
				}
				return collectResources(projects, func(p *descope.Project) (string, string) { return p.Name, p.ID }), nil
			},
			delFn: func(r resource) error {
				projectClient, err := descopeclient.NewWithConfig(&descopeclient.Config{
					ManagementKey:  managementKey,
					DescopeBaseURL: baseURL,
					ProjectID:      r.id,
				})
				if err != nil {
					return err
				}
				return projectClient.Management.Project().Delete(ctx)
			},
		},
	}

	for _, c := range cleanups {
		d, f := runCleanup(c.name, c.listFn, c.delFn)
		totalDeleted += d
		totalFailed += f
	}

	fmt.Printf("\ntotal: %d deleted, %d failed\n", totalDeleted, totalFailed)
	if totalFailed > 0 {
		os.Exit(1)
	}
}

func collectResources[T any](items []T, extract func(T) (name, id string)) []resource {
	var result []resource
	for _, item := range items {
		name, id := extract(item)
		if strings.HasPrefix(name, testPrefix) {
			result = append(result, resource{name: name, id: id})
		}
	}
	return result
}

func runCleanup(typeName string, listFn func() ([]resource, error), delFn func(resource) error) (deleted, failed int) {
	resources, err := listFn()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to list %ss: %v\n", typeName, err)
		return 0, 1
	}

	for _, r := range resources {
		fmt.Printf("deleting %s %s (%s)...\n", typeName, r.name, r.id)
		if err := delFn(r); err != nil {
			fmt.Fprintf(os.Stderr, "  failed: %v\n", err)
			failed++
			continue
		}
		deleted++
	}

	if deleted > 0 || failed > 0 {
		fmt.Printf("%ss: %d deleted, %d failed\n", typeName, deleted, failed)
	}
	return
}
