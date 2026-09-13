// Command testcleanup deletes all Descope test resources whose names start with "testacc-".
//
// Passes run in this order: tenants first (cascade deletion there removes keys
// bound only to the tenant, shrinking the access-key pass), then access keys,
// management keys, descopers, and projects.
//
// It requires DESCOPE_MANAGEMENT_KEY and DESCOPE_BASE_URL environment variables.
// Usage: source .env && go run ./tools/testcleanup
//
// The tenant pass reaches inside every project the management key can see,
// including production ones, so -dry-run prints exactly what would be deleted
// without deleting anything:
//
//	source .env && go run ./tools/testcleanup -dry-run
package main

import (
	"context"
	"flag"
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

// projectClientFn binds a management client to a single project. Tenant
// operations resolve against whichever project the client carries, so the
// company-scoped client the other passes share cannot reach them.
type projectClientFn func(projectID string) (sdk.Management, error)

func main() {
	dryRun := flag.Bool("dry-run", false, "list what would be deleted without deleting anything")
	flag.Parse()

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

	newProjectClient := func(projectID string) (sdk.Management, error) {
		projectClient, err := descopeclient.NewWithConfig(&descopeclient.Config{
			ManagementKey:  managementKey,
			DescopeBaseURL: baseURL,
			ProjectID:      projectID,
		})
		if err != nil {
			return nil, err
		}
		return projectClient.Management, nil
	}

	if *dryRun {
		fmt.Println("dry run: nothing will be deleted")
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
				projectMgmt, err := newProjectClient(r.id)
				if err != nil {
					return err
				}
				return projectMgmt.Project().Delete(ctx)
			},
		},
	}

	td, tf := cleanupTenants(ctx, mgmt, newProjectClient, *dryRun)
	totalDeleted += td
	totalFailed += tf

	for _, c := range cleanups {
		d, f := runCleanup(c.name, c.listFn, c.delFn, *dryRun)
		totalDeleted += d
		totalFailed += f
	}

	verb := "deleted"
	if *dryRun {
		verb = "would delete"
	}
	fmt.Printf("\ntotal: %d %s, %d failed\n", totalDeleted, verb, totalFailed)
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
// That breadth is also the risk: every project the management key can see is
// in scope, production included, and the name prefix is the only guard. The
// scanned projects are printed before any deletion, and dryRun reports what
// would go without touching anything.
//
// Deletion passes cascade=true, which removes users and keys associated only
// with the tenant being deleted. Anything shared with another tenant survives.
//
// Neither ListProjects nor Tenant().LoadAll paginates in go-sdk v1.32.0 — both
// issue a single request and return the whole set, with no cursor on the
// request or the response — so a single call sees every project and tenant.
func cleanupTenants(ctx context.Context, mgmt sdk.Management, newProjectClient projectClientFn, dryRun bool) (deleted, failed int) {
	projects, err := mgmt.Project().ListProjects(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to list projects for tenant cleanup: %v\n", err)
		return 0, 1
	}

	names := make([]string, 0, len(projects))
	for _, p := range projects {
		names = append(names, p.Name)
	}
	fmt.Printf("scanning %d project(s) for %s tenants: %s\n", len(projects), testPrefix, strings.Join(names, ", "))

	for _, p := range projects {
		projectMgmt, err := newProjectClient(p.ID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create client for project %s: %v\n", p.Name, err)
			failed++
			continue
		}

		tenants, err := projectMgmt.Tenant().LoadAll(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to list tenants in project %s: %v\n", p.Name, err)
			failed++
			continue
		}

		for _, t := range tenants {
			if !strings.HasPrefix(t.Name, testPrefix) {
				continue
			}
			if dryRun {
				fmt.Printf("would delete tenant %s (%s) in project %s\n", t.Name, t.ID, p.Name)
				deleted++
				continue
			}
			fmt.Printf("deleting tenant %s (%s) in project %s...\n", t.Name, t.ID, p.Name)
			if err := projectMgmt.Tenant().Delete(ctx, t.ID, true); err != nil {
				fmt.Fprintf(os.Stderr, "  failed: %v\n", err)
				failed++
				continue
			}
			deleted++
		}
	}

	if deleted > 0 || failed > 0 {
		fmt.Printf("tenants: %d %s, %d failed\n", deleted, deletedVerb(dryRun), failed)
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

func runCleanup(typeName string, listFn func() ([]resource, error), delFn func(resource) error, dryRun bool) (deleted, failed int) {
	resources, err := listFn()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to list %ss: %v\n", typeName, err)
		return 0, 1
	}

	for _, r := range resources {
		if dryRun {
			fmt.Printf("would delete %s %s (%s)\n", typeName, r.name, r.id)
			deleted++
			continue
		}
		fmt.Printf("deleting %s %s (%s)...\n", typeName, r.name, r.id)
		if err := delFn(r); err != nil {
			fmt.Fprintf(os.Stderr, "  failed: %v\n", err)
			failed++
			continue
		}
		deleted++
	}

	if deleted > 0 || failed > 0 {
		fmt.Printf("%ss: %d %s, %d failed\n", typeName, deleted, deletedVerb(dryRun), failed)
	}
	return
}

func deletedVerb(dryRun bool) string {
	if dryRun {
		return "would delete"
	}
	return "deleted"
}
