package vex

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"os"
)

var stderr = os.Stderr

func New(version string) func() tfsdk.Provider {
	return func() tfsdk.Provider {
		return &provider{
			version: version,
		}
	}
}

type provider struct {
	configured bool
	version    string
	apiToken   string
	host       string
}

func (p *provider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"host": {
				Type:     types.StringType,
				Required: true,
			},
			"api_token": {
				Type:     types.StringType,
				Required: true,
			},
		},
	}, nil
}

// Provider schema struct
type providerData struct {
	APIToken types.String `tfsdk:"api_token"`
	Host     types.String `tfsdk:"host"`
}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	var config providerData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var host string
	if config.Host.Unknown {
		resp.Diagnostics.AddError(
			"Unable to create client",
			"Cannot use unknown value as host")
	}
	if config.Host.Null {
		host = os.Getenv("VEX_HOST")
	} else {
		host = config.Host.Value
	}
	if host == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find host",
			"Host cannot be an empty string",
		)
		return
	}
	var apiToken string
	if config.APIToken.Unknown {
		resp.Diagnostics.AddError(
			"Unable to create client",
			"Cannot use unknown value as API token")
	}
	if config.APIToken.Null {
		apiToken = os.Getenv("VEX_API_TOKEN")
	} else {
		apiToken = config.APIToken.Value
	}
	if apiToken == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find API token",
			"API token cannot be an empty string",
		)
		return
	}
	p.host = host
	p.apiToken = apiToken
	p.configured = true
}

// GetResources - Defines provider resources
func (p *provider) GetResources(_ context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{
		"vex_account": resourceAccountType{},
		"vex_project": resourceProjectType{},
		"vex_flag":    resourceFlagType{},
	}, nil
}

// GetDataSources - Defines provider data sources
func (p *provider) GetDataSources(_ context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{}, nil
}
