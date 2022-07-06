package vex

import (
	"context"
	"errors"
	"fmt"
	vex_go "github.com/broswen/vex-go"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"
)

type flagResourceData struct {
	ID        types.String `tfsdk:"id"`
	AccountID types.String `tfsdk:"account_id"`
	ProjectID types.String `tfsdk:"project_id"`
	Key       types.String `tfsdk:"key"`
	Type      types.String `tfsdk:"type"`
	Value     types.String `tfsdk:"value"`
}

func flagTypeToString(f vex_go.FlagType) string {
	switch f {
	case vex_go.BOOLEAN:
		return "BOOLEAN"
	case vex_go.STRING:
		return "STRING"
	case vex_go.NUMBER:
		return "NUMBER"
	default:
		return "UNKNOWN"
	}
}

func stringToFlagType(s string) (vex_go.FlagType, error) {
	switch s {
	case "BOOLEAN":
		return vex_go.BOOLEAN, nil
	case "STRING":
		return vex_go.STRING, nil
	case "NUMBER":
		return vex_go.NUMBER, nil
	}
	return "", errors.New("invalid flag type")
}

func expandFlag(d flagResourceData) (*vex_go.Flag, error) {
	f := &vex_go.Flag{}
	f.ID = d.ID.Value
	f.AccountID = d.AccountID.Value
	f.ProjectID = d.ProjectID.Value
	f.Key = d.Key.Value
	t, err := stringToFlagType(d.Type.Value)
	if err != nil {
		return nil, err
	}
	f.Type = t
	f.Value = d.Value.Value
	return f, nil
}

func flattenFlag(f *vex_go.Flag) *flagResourceData {
	d := &flagResourceData{}
	d.ID = types.String{Value: f.ID}
	d.AccountID = types.String{Value: f.AccountID}
	d.ProjectID = types.String{Value: f.ProjectID}
	d.Key = types.String{Value: f.Key}
	d.Type = types.String{Value: flagTypeToString(f.Type)}
	d.Value = types.String{Value: f.Value}
	return d
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
	var plan flagResourceData
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	f, err := expandFlag(plan)
	if err != nil {
		resp.Diagnostics.AddError("error expanding flag", err.Error())
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}
	err = r.p.client.CreateFlag(ctx, f)
	if err != nil {
		resp.Diagnostics.AddError("error creating flag", err.Error())
		return
	}
	diags = resp.State.Set(ctx, flattenFlag(f))
	resp.Diagnostics.Append(diags...)
}

// Read resource information
func (r resourceFlag) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var d flagResourceData
	diags := req.State.Get(ctx, &d)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
	f, err := r.p.client.GetFlag(ctx, d.AccountID.Value, d.ProjectID.Value, d.ID.Value)
	if err != nil {
		resp.Diagnostics.AddError("error reading flag", err.Error())
		return
	}
	diags = resp.State.Set(ctx, flattenFlag(f))
	resp.Diagnostics.Append(diags...)
}

// Update resource
func (r resourceFlag) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var state, plan flagResourceData
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	plan.ID = state.ID
	f, err := expandFlag(plan)
	if err != nil {
		resp.Diagnostics.AddError("error expanding flag", err.Error())
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}
	err = r.p.client.UpdateFlag(ctx, f)
	if err != nil {
		resp.Diagnostics.AddError("error updating flag", err.Error())
		return
	}
	diags = resp.State.Set(ctx, flattenFlag(f))
	resp.Diagnostics.Append(diags...)
}

// Delete resource
func (r resourceFlag) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var d flagResourceData
	diags := req.State.Get(ctx, &d)
	resp.Diagnostics.Append(diags...)

	f, err := expandFlag(d)
	if err != nil {
		resp.Diagnostics.AddError("error expanding flag", err.Error())
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}
	err = r.p.client.DeleteFlag(ctx, f)

	if err != nil {
		resp.Diagnostics.AddError("error deleting project", err.Error())
		return
	}
	f = &vex_go.Flag{}
	diags = resp.State.Set(ctx, flattenFlag(f))
	resp.Diagnostics.Append(diags...)
}

// Import resource
func (r resourceFlag) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	parts := strings.Split(req.ID, "/")
	tflog.Info(ctx, fmt.Sprintf("%#v", parts))
	if len(parts) != 3 {
		resp.Diagnostics.AddError("invalid flag id", "project id must be in the form <account_id>/<project_id>/<flag_id>")
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}
	f, err := r.p.client.GetFlag(ctx, parts[0], parts[1], parts[2])
	if err != nil {
		resp.Diagnostics.AddError("error importing flag", err.Error())
		return
	}

	diags := resp.State.Set(ctx, flattenFlag(f))
	resp.Diagnostics.Append(diags...)
}
