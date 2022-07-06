package vex

import (
	"context"
	vex_go "github.com/broswen/vex-go"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type accountResourceData struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

func expandAccount(d accountResourceData) *vex_go.Account {
	a := &vex_go.Account{}
	a.ID = d.ID.Value
	a.Name = d.Name.Value
	a.Description = d.Description.Value
	return a
}

func flattenAccount(a *vex_go.Account) *accountResourceData {
	d := &accountResourceData{}
	d.ID = types.String{Value: a.ID}
	d.Name = types.String{Value: a.Name}
	d.Description = types.String{Value: a.Description}
	return d
}

type resourceAccountType struct{}

func (r resourceAccountType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"name": {
				Type:     types.StringType,
				Required: true,
			},
			"description": {
				Type:     types.StringType,
				Required: true,
			},
		},
	}, nil
}

func (r resourceAccountType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceAccount{
		p: *(p.(*provider)),
	}, nil
}

type resourceAccount struct {
	p provider
}

// Create a new resource
func (r resourceAccount) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	resp.Diagnostics.AddError("creating accounts is not support", "Creating vex accounts via the Terraform provider is not support.")
}

// Read resource information
func (r resourceAccount) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var d accountResourceData
	diags := req.State.Get(ctx, &d)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
	a, err := r.p.client.GetAccount(ctx, d.ID.Value)
	if err != nil {
		resp.Diagnostics.AddError("error reading account", err.Error())
		return
	}
	diags = resp.State.Set(ctx, flattenAccount(a))
	resp.Diagnostics.Append(diags...)
}

// Update resource
func (r resourceAccount) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var state, plan accountResourceData
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	plan.ID = state.ID
	a := expandAccount(plan)
	if resp.Diagnostics.HasError() {
		return
	}
	err := r.p.client.UpdateAccount(ctx, a)
	if err != nil {
		resp.Diagnostics.AddError("error updating account", err.Error())
		return
	}
	diags = resp.State.Set(ctx, flattenAccount(a))
	resp.Diagnostics.Append(diags...)
}

// Delete resource
func (r resourceAccount) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {

	var d accountResourceData
	diags := req.State.Get(ctx, &d)
	resp.Diagnostics.Append(diags...)

	a := expandAccount(d)

	if resp.Diagnostics.HasError() {
		return
	}
	err := r.p.client.DeleteAccount(ctx, a.ID)
	if err != nil {
		resp.Diagnostics.AddError("error deleting account", err.Error())
		return
	}
	a = &vex_go.Account{}
	diags = resp.State.Set(ctx, flattenAccount(a))
	resp.Diagnostics.Append(diags...)
}

// Import resource
func (r resourceAccount) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	if resp.Diagnostics.HasError() {
		return
	}
	a, err := r.p.client.GetAccount(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.AddError("error reading account", err.Error())
		return
	}

	diags := resp.State.Set(ctx, flattenAccount(a))
	resp.Diagnostics.Append(diags...)
}
