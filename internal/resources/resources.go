package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/accesskey"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/descoper"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/inboundapp"
	"github.com/jamescrowley321/terraform-provider-descope/internal/models/managementkey"
)

func NewAccessKeyResource() resource.Resource {
	return newResource[accesskey.AccessKeyModel]("access_key", accesskey.Schema)
}

func NewDescoperResource() resource.Resource {
	return newResource[descoper.DescoperModel]("descoper", descoper.Schema)
}

func NewManagementKeyResource() resource.Resource {
	return newResource[managementkey.ManagementKeyModel]("management_key", managementkey.Schema)
}

func NewInboundAppResource() resource.Resource {
	return newResource[inboundapp.InboundAppModel]("inbound_app", inboundapp.Schema)
}
