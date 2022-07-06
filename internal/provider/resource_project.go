package vex

import (
	"context"
	"fmt"
	vex_go "github.com/broswen/vex-go"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"
)

type projectResourceData struct {
	ID          types.String `tfsdk:"id"`
	AccountID   types.String `tfsdk:"account_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

func expandProject(d projectResourceData) *vex_go.Project {
	p := &vex_go.Project{}
	p.ID = d.ID.Value
	p.AccountID = d.AccountID.Value
	p.Name = d.Name.Value
	p.Description = d.Description.Value
	return p
}

func flattenProject(a *vex_go.Project) *projectResourceData {
	d := &projectResourceData{}
	d.ID = types.String{Value: a.ID}
	d.AccountID = types.String{Value: a.AccountID}
	d.Name = types.String{Value: a.Name}
	d.Description = types.String{Value: a.Description}
	return d
}

type resourceProjectType struct{}

func (r resourceProjectType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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

func (r resourceProjectType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceProject{
		p: *(p.(*provider)),
	}, nil
}

type resourceProject struct {
	p provider
}

// Create a new resource
func (r resourceProject) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var plan projectResourceData
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	p := expandProject(plan)
	if resp.Diagnostics.HasError() {
		return
	}
	err := r.p.client.CreateProject(ctx, p)
	if err != nil {
		resp.Diagnostics.AddError("error updating project", err.Error())
		return
	}
	diags = resp.State.Set(ctx, flattenProject(p))
	resp.Diagnostics.Append(diags...)
}

// Read resource information
func (r resourceProject) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var d projectResourceData
	diags := req.State.Get(ctx, &d)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
	p, err := r.p.client.GetProject(ctx, d.AccountID.Value, d.ID.Value)
	if err != nil {
		resp.Diagnostics.AddError("error reading project", err.Error())
		return
	}
	diags = resp.State.Set(ctx, flattenProject(p))
	resp.Diagnostics.Append(diags...)
}

// Update resource
func (r resourceProject) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var state, plan projectResourceData
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	plan.ID = state.ID
	p := expandProject(plan)
	if resp.Diagnostics.HasError() {
		return
	}
	err := r.p.client.UpdateProject(ctx, p)
	if err != nil {
		resp.Diagnostics.AddError("error updating project", err.Error())
		return
	}
	diags = resp.State.Set(ctx, flattenProject(p))
	resp.Diagnostics.Append(diags...)
}

// Delete resource
func (r resourceProject) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {

	var d projectResourceData
	diags := req.State.Get(ctx, &d)
	resp.Diagnostics.Append(diags...)

	p := expandProject(d)

	if resp.Diagnostics.HasError() {
		return
	}
	err := r.p.client.DeleteProject(ctx, p)
	if err != nil {
		resp.Diagnostics.AddError("error deleting project", err.Error())
		return
	}
	p = &vex_go.Project{}
	diags = resp.State.Set(ctx, flattenProject(p))
	resp.Diagnostics.Append(diags...)
}

// Import resource
func (r resourceProject) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	parts := strings.Split(req.ID, "/")
	tflog.Info(ctx, fmt.Sprintf("%#v", parts))
	if len(parts) != 2 {
		resp.Diagnostics.AddError("invalid project id", "project id must be in the form <account_id>/<project_id>")
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}
	a, err := r.p.client.GetProject(ctx, parts[0], parts[1])
	if err != nil {
		resp.Diagnostics.AddError("error importing project", err.Error())
		return
	}

	diags := resp.State.Set(ctx, flattenProject(a))
	resp.Diagnostics.Append(diags...)
}
