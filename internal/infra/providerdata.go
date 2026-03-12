package infra

import (
	descopeclient "github.com/descope/go-sdk/descope/client"
	"github.com/descope/go-sdk/descope/sdk"
)

// ProviderData holds all clients needed by resources.
type ProviderData struct {
	// Client is the infra API client used by project, management_key, and descoper resources.
	Client *Client

	// Management is the Descope SDK management client used by resources that
	// interact with management APIs directly (e.g., access keys).
	Management sdk.Management
}

func NewProviderData(version, managementKey, baseURL, projectID string) (*ProviderData, error) {
	client := NewClient(version, managementKey, baseURL)

	descopeClient, err := descopeclient.NewWithConfig(&descopeclient.Config{
		ProjectID:           projectID,
		ManagementKey:       managementKey,
		DescopeBaseURL:      baseURL,
		AllowEmptyProjectID: true,
	})
	if err != nil {
		return nil, err
	}

	return &ProviderData{
		Client:     client,
		Management: descopeClient.Management,
	}, nil
}
