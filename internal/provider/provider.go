package vex

import (
	"context"
	vex_go "github.com/broswen/vex-go"
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
	client     *vex_go.Client
}

func (p *provider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
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
}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	var config providerData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
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
	client, err := vex_go.New(apiToken)
	if err != nil {
		resp.Diagnostics.AddError(
			err.Error(),
			err.Error(),
		)
		return
	}
	p.client = client
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
