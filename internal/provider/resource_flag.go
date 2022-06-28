package vex

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Flag struct {
	ID        string `json:"id"`
	AccountID string `json:"account_id"`
	ProjectID string `json:"project_id"`
	Key       string `json:"key"`
	Type      string `json:"type"`
	Value     string `json:"Value"`
}

var FlagTypes = []string{"BOOLEAN", "STRING", "NUMBER"}

type resourceFlagType struct{}

func (r resourceFlagType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"account_id": {
				Type:     types.StringType,
				Required: true,
			},
			"project_id": {
				Type:     types.StringType,
				Required: true,
			},
			"key": {
				Type:     types.StringType,
				Required: true,
			},
			"type": {
				Type:     types.StringType,
				Required: true,
				//Validators: match to string slice of flag types
			},
			"value": {
				Type:     types.StringType,
				Required: true,
			},
		},
	}, nil
}

func (r resourceFlagType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceFlag{
		p: *(p.(*provider)),
	}, nil
}

type resourceFlag struct {
	p provider
}

// Create a new resource
func (r resourceFlag) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
}

// Read resource information
func (r resourceFlag) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
}

// Update resource
func (r resourceFlag) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
}

// Delete resource
func (r resourceFlag) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
}

// Import resource
func (r resourceFlag) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
}
