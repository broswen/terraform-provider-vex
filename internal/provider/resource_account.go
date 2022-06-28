package vex

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"log"
	"net/http"
)

type Account struct {
	ID          string `json:"id" tfsdk:"id"`
	Name        string `json:"name" tfsdk:"name"`
	Description string `json:"description" tfsdk:"description"`
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
	var a Account
	diags := req.State.Get(ctx, &a)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
	//	GET vex.broswen.com/accounts/{accountId}
	vexReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/accounts/%s", r.p.host, a.ID), nil)
	if err != nil {
		resp.Diagnostics.AddError("error reading account", err.Error())
		return
	}
	vexReq.Header.Set("Authorization", "Bearer "+r.p.apiToken)
	vexResp, err := http.DefaultClient.Do(vexReq)
	if vexResp.StatusCode != 200 {
		resp.Diagnostics.AddError("error reading account", vexResp.Status)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("error reading account", err.Error())
		return
	}
	err = json.NewDecoder(vexResp.Body).Decode(&a)
	if err != nil {
		resp.Diagnostics.AddError("error parsing account json", err.Error())
		return
	}

	diags = resp.State.Set(ctx, &a)
	resp.Diagnostics.Append(diags...)
}

// Update resource
func (r resourceAccount) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var a Account
	diags := req.State.Get(ctx, &a)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
	body, err := json.Marshal(a)
	if err != nil {
		resp.Diagnostics.AddError("error reading account", err.Error())
		return
	}
	//	PUT vex.broswen.com/accounts/{accountId}
	vexReq, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/accounts/%s", r.p.host, a.ID), bytes.NewReader(body))
	if err != nil {
		resp.Diagnostics.AddError("error updating account", err.Error())
		return
	}
	vexReq.Header.Set("Authorization", "Bearer "+r.p.apiToken)
	vexResp, err := http.DefaultClient.Do(vexReq)
	if vexResp.StatusCode != 200 {
		resp.Diagnostics.AddError("error updating account", vexResp.Status)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("error updating account", err.Error())
		return
	}
	err = json.NewDecoder(vexResp.Body).Decode(&a)
	if err != nil {
		resp.Diagnostics.AddError("error updating account json", err.Error())
		return
	}

	diags = resp.State.Set(ctx, &a)
	resp.Diagnostics.Append(diags...)
}

// Delete resource
func (r resourceAccount) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var a Account
	diags := req.State.Get(ctx, &a)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
	//	DELETE vex.broswen.com/accounts/{accountId}
	vexReq, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/accounts/%s", r.p.host, a.ID), nil)
	if err != nil {
		resp.Diagnostics.AddError("error reading account", err.Error())
		return
	}
	vexReq.Header.Set("Authorization", "Bearer "+r.p.apiToken)
	vexResp, err := http.DefaultClient.Do(vexReq)
	if vexResp.StatusCode != 200 {
		resp.Diagnostics.AddError("error reading account", vexResp.Status)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("error reading account", err.Error())
		return
	}
	err = json.NewDecoder(vexResp.Body).Decode(&a)
	if err != nil {
		resp.Diagnostics.AddError("error parsing account json", err.Error())
		return
	}

	diags = resp.State.Set(ctx, &a)
	resp.Diagnostics.Append(diags...)
}

// Import resource
func (r resourceAccount) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	var a Account
	if resp.Diagnostics.HasError() {
		return
	}
	//	GET vex.broswen.com/accounts/{accountId}
	vexReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/accounts/%s", r.p.host, req.ID), nil)
	log.Println(vexReq.URL)
	if err != nil {
		resp.Diagnostics.AddError("error reading account", err.Error())
		return
	}
	vexReq.Header.Set("Authorization", "Bearer "+r.p.apiToken)
	vexResp, err := http.DefaultClient.Do(vexReq)
	if vexResp.StatusCode != 200 {
		resp.Diagnostics.AddError("error reading account", vexResp.Status)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("error reading account", err.Error())
		return
	}
	err = json.NewDecoder(vexResp.Body).Decode(&a)
	if err != nil {
		resp.Diagnostics.AddError("error parsing account json", err.Error())
		return
	}

	diags := resp.State.Set(ctx, &a)
	resp.Diagnostics.Append(diags...)
}
