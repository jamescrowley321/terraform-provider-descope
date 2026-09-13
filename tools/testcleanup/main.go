// Command testcleanup deletes all Descope test resources whose names start with "testacc-".
//
// It cleans up projects, tenants, access keys, management keys, and descopers.
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

	// Tenants are project-scoped, so they need a client per project rather than
	// the company-scoped one the table-driven cleanups share. Run them first:
	// deleting a tenant with cascade also removes the access keys bound only to
	// it, which shrinks the work the access-key pass has to do.
	td, tf := cleanupTenants(ctx, managementKey, baseURL, mgmt)
	totalDeleted += td
	totalFailed += tf

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

// cleanupTenants deletes testacc- tenants from every project.
//
// Tenants cannot be reached through the company-scoped client the other
// cleanups use — tenant operations resolve against whichever project the
// client is bound to — so this lists projects and rebinds a client to each one
// in turn. Stray testacc- tenants accumulate inside *real* projects when an
// acceptance test dies before its own cleanup, which is precisely the case the
// prefix-matching passes above never saw.
//
// Deletion passes cascade=true, which removes users and keys associated only
// with the tenant being deleted. Anything shared with another tenant survives.
func cleanupTenants(ctx context.Context, managementKey, baseURL string, mgmt sdk.Management) (deleted, failed int) {
	projects, err := mgmt.Project().ListProjects(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to list projects for tenant cleanup: %v\n", err)
		return 0, 1
	}

	for _, p := range projects {
		projectClient, err := descopeclient.NewWithConfig(&descopeclient.Config{
			ManagementKey:  managementKey,
			DescopeBaseURL: baseURL,
			ProjectID:      p.ID,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create client for project %s: %v\n", p.Name, err)
			failed++
			continue
		}

		tenants, err := projectClient.Management.Tenant().LoadAll(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to list tenants in project %s: %v\n", p.Name, err)
			failed++
			continue
		}

		for _, t := range tenants {
			if !strings.HasPrefix(t.Name, testPrefix) {
				continue
			}
			fmt.Printf("deleting tenant %s (%s) in project %s...\n", t.Name, t.ID, p.Name)
			if err := projectClient.Management.Tenant().Delete(ctx, t.ID, true); err != nil {
				fmt.Fprintf(os.Stderr, "  failed: %v\n", err)
				failed++
				continue
			}
			deleted++
		}
	}

	if deleted > 0 || failed > 0 {
		fmt.Printf("tenants: %d deleted, %d failed\n", deleted, failed)
	}
	return
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
