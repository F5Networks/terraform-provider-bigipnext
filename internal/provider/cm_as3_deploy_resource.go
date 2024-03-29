package provider

import (
	"context"
	"fmt"
	"sync"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	bigipnextsdk "gitswarm.f5net.com/terraform-providers/bigipnext"
)

var (
	_     resource.Resource                = &NextCMAS3DeployResource{}
	_     resource.ResourceWithImportState = &NextCMAS3DeployResource{}
	mutex sync.Mutex
)

func NewNextCMAS3DeployResource() resource.Resource {
	return &NextCMAS3DeployResource{}
}

type NextCMAS3DeployResource struct {
	client *bigipnextsdk.BigipNextCM
}

type NextCMAS3DeployResourceModel struct {
	As3Json       types.String `tfsdk:"as3_json"`
	TargetAddress types.String `tfsdk:"target_address"`
	Timeout       types.Int64  `tfsdk:"timeout"`
	DraftId       types.String `tfsdk:"draft_id"`
	DeployId      types.String `tfsdk:"deploy_id"`
	Id            types.String `tfsdk:"id"`
}

func (r *NextCMAS3DeployResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cm_as3_deploy"
}

func (r *NextCMAS3DeployResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Deploy an AS3 declaration to a specified instance managed by BIG-IP Next Central Manager. If the deployment already exists on a different instance, the application service is removed from the existing instance before deploying to the new instance",
		Attributes: map[string]schema.Attribute{
			"as3_json": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "AS3 Json Declaration to be post onto BIG-IP Next",
			},
			"target_address": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Target Address of the Device Inventory on BIG-IP Next CM.",
			},
			"timeout": schema.Int64Attribute{
				MarkdownDescription: "The number of seconds to wait for instance deployment to finish.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(900),
			},
			"draft_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Draft ID of the AS3 declaration on BIG-IP CM Next",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"deploy_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Deploy ID of the AS3 declaration on BIG-IP CM Next",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique Identifier for the resource",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *NextCMAS3DeployResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client, resp.Diagnostics = toBigipNextCMProvider(req.ProviderData)
}

func (r *NextCMAS3DeployResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var resCfg *NextCMAS3DeployResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &resCfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	//as3Config := resCfg.As3Json.ValueString()

	tflog.Info(ctx, fmt.Sprintf("[CREATE]Posting Application service config:%+v", resCfg.As3Json.ValueString()))
	drartID, err := r.client.PostAS3DraftDocument(resCfg.As3Json.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to Create AS3 config Drart, got error: %s", err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Application Service Draft ID:%+v", drartID))
	DeployID, err := r.client.CMAS3DeployNext(drartID, resCfg.TargetAddress.ValueString(), int(resCfg.Timeout.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to Deploy AS3 config, got error: %s", err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("AS3 Deployment ID:%+v", DeployID))
	resCfg.Id = types.StringValue(drartID)
	resCfg.DraftId = types.StringValue(drartID)
	resCfg.DeployId = types.StringValue(DeployID)
	resCfg.TargetAddress = types.StringValue(resCfg.TargetAddress.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, resCfg)...)
}

func (r *NextCMAS3DeployResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var stateCfg *NextCMAS3DeployResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateCfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	draftID := stateCfg.Id.ValueString()
	deployID := stateCfg.DeployId.ValueString()
	tflog.Info(ctx, "Reading AS3 Service Deployment")
	as3Resp, err := r.client.GetAS3DeploymentTaskStatus(draftID, deployID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to READ AS3 config, got error: %s", err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("AS3 Config Response:%+v", as3Resp))
	// stateCfg.As3Json = types.StringValue(as3Resp.(string))
	// {"class":"ADC","id":"example-declaration-01","label":"Sample 1","remark":"Simple HTTP application with round robin pool","schemaVersion":"3.45.0","tenant-cm-test1":{"app-cm-test1":{"class":"Application","pool-cm-test1":{"class":"Pool","members":[{"serverAddresses":["10.62.0.1"],"servicePort":80}],"monitors":["http"]},"serviceMain":{"class":"Service_HTTP","pool":"pool-cm-test1","virtualAddresses":["10.1.0.1"]},"template":"http"},"class":"Tenant"},"tenant-cm-test2":{"app-cm-test2":{"class":"Application","pool-cm-test2":{"class":"Pool","members":[{"serverAddresses":["10.63.0.1"],"servicePort":80}],"monitors":["http"]},"serviceMain":{"class":"Service_HTTP","pool":"pool-cm-test2","virtualAddresses":["10.1.0.2"]},"template":"http"},"class":"Tenant"}}
	// stateCfg.TargetAddress = types.StringValue(trgtAddrss)
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateCfg)...)
}

func (r *NextCMAS3DeployResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var resCfg *NextCMAS3DeployResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &resCfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	as3Json := resCfg.As3Json.ValueString()
	draftID := resCfg.Id.ValueString()
	deployID := resCfg.DeployId.ValueString()

	tflog.Info(ctx, fmt.Sprintf("[UPDATE]Update AS3 application service: %s", as3Json))

	err := r.client.PutAS3DraftDocument(draftID, resCfg.As3Json.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to Update AS3 application service, got error: %s", err))
		return
	}
	tflog.Info(ctx, "Reading AS3 Service Deployment")
	// time.Sleep(5 * time.Second)
	as3Resp, err := r.client.GetAS3DeploymentTaskStatus(draftID, deployID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to READ AS3 config, got error: %s", err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("AS3 Config Response:%+v", as3Resp))
	resCfg.Id = types.StringValue(draftID)
	resCfg.DeployId = types.StringValue(deployID)
	resCfg.As3Json = types.StringValue(as3Json)
	resCfg.TargetAddress = types.StringValue(resCfg.TargetAddress.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &resCfg)...)
}

func (r *NextCMAS3DeployResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var stateCfg *NextCMAS3DeployResourceModel
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &stateCfg)...)
	draftID := stateCfg.Id.ValueString()
	tflog.Info(ctx, fmt.Sprintf("Deleting AS3 application service Draft: %s", draftID))
	err := r.client.DeleteAS3DeploymentTask(draftID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to Delete AS3 Application service, got error: %s", err))
		return
	}
	stateCfg.Id = types.StringValue("")
}

func (r *NextCMAS3DeployResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

//
//func contains(s []string, str string) bool {
//	for _, v := range s {
//		if v == str {
//			return true
//		}
//	}
//	return false
//}
