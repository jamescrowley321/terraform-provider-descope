package thirdpartyapp

import (
	"github.com/descope/terraform-provider-descope/internal/models/attrs/stringattr"
	"github.com/descope/terraform-provider-descope/internal/models/attrs/strsetattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

var Attributes = map[string]schema.Attribute{
	"id":                     stringattr.Identifier(),
	"name":                   stringattr.Required(),
	"description":            stringattr.Default(""),
	"logo":                   stringattr.Default(""),
	"login_page_url":         stringattr.Optional(),
	"client_id":              stringattr.Identifier(),
	"client_secret":          stringattr.SecretComputed(),
	"approved_callback_urls": strsetattr.Default(),
}

type Model struct {
	ID                   stringattr.Type `tfsdk:"id"`
	Name                 stringattr.Type `tfsdk:"name"`
	Description          stringattr.Type `tfsdk:"description"`
	Logo                 stringattr.Type `tfsdk:"logo"`
	LoginPageURL         stringattr.Type `tfsdk:"login_page_url"`
	ClientID             stringattr.Type `tfsdk:"client_id"`
	ClientSecret         stringattr.Type `tfsdk:"client_secret"`
	ApprovedCallbackUrls strsetattr.Type `tfsdk:"approved_callback_urls"`
}
