package infra

import (
	"errors"
	"strings"

	"github.com/descope/go-sdk/descope/client"
	"github.com/descope/go-sdk/descope/sdk"
)

func (c *Client) Management(projectID string) (sdk.Management, error) {
	if strings.TrimSpace(projectID) == "" {
		return nil, errors.New("project_id must not be empty")
	}
	sdkClient, err := client.NewWithConfig(&client.Config{ProjectID: projectID, ManagementKey: c.managementKey, DescopeBaseURL: c.baseURL})
	if err != nil {
		return nil, err
	}
	return sdkClient.Management, nil
}

func IsNotFoundError(err error) bool { return AsNotFoundError(err) }
