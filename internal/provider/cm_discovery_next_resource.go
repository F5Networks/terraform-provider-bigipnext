package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	bigipnextsdk "gitswarm.f5net.com/terraform-providers/bigipnext"
)

var (
	_ resource.Resource                = &CMDiscoveryNextResource{}
	_ resource.ResourceWithImportState = &CMDiscoveryNextResource{}
)

func NewCMDiscoveryNextResource() resource.Resource {
	return &CMDiscoveryNextResource{}
}

type CMDiscoveryNextResource struct {
	client *bigipnextsdk.BigipNextCM
}

type CMDiscoveryNextResourceModel struct {
	Address            types.String `tfsdk:"address"`
	Port               types.Int64  `tfsdk:"port"`
	DeviceUser         types.String `tfsdk:"device_user"`
	DevicePassword     types.String `tfsdk:"device_password"`
	ManagementUser     types.String `tfsdk:"management_user"`
	ManagementPassword types.String `tfsdk:"management_password"`
	Id                 types.String `tfsdk:"id"`
}

func (r *CMDiscoveryNextResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cm_discover_next"
}

func (r *CMDiscoveryNextResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource used for add\t(discover)\t BIG-IP Next instance to BIG-IP Next Central Manager for management",
		Attributes: map[string]schema.Attribute{
			"address": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "IP Address of the BIG-IP Next instance to be discovered",
			},
			"port": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Port number of the BIG-IP Next instance to be discovered",
			},
			"device_user": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The username that the BIG-IP Next Central Manager uses before Instance discovery for BIG-IP Next management",
			},
			"device_password": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "The password that the BIG-IP Next Central Manager uses before Instance discovery for BIG-IP Next management",
			},
			"management_user": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The username that the BIG-IP Next Central Manager uses after Instance Discovery for BIG-IP Next management",
			},
			"management_password": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "The password that the BIG-IP Next Central Manager uses after Instance Discovery for BIG-IP Next management",
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

func (r *CMDiscoveryNextResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client, resp.Diagnostics = toBigipNextCMProvider(req.ProviderData)
}

func (r *CMDiscoveryNextResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var resCfg *CMDiscoveryNextResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &resCfg)...)
	if resp.Diagnostics.HasError() { // coverage-ignore
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[CREATE] CMDiscoveryNextResource:%+v\n", resCfg.Address.ValueString()))

	providerConfig := getCMDiscoveryNextConfig(ctx, resCfg)

	tflog.Info(ctx, fmt.Sprintf("[CREATE] Device Provider config:%+v\n", providerConfig))

	respData, err := r.client.DiscoverInstance(providerConfig)
	if err != nil { // coverage-ignore
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to Create Certificate, got error: %s", err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[CREATE] respData ID:%+v\n", string(respData)))
	resCfg.Id = types.StringValue(string(respData))
	resp.Diagnostics.Append(resp.State.Set(ctx, resCfg)...)
}

func (r *CMDiscoveryNextResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var stateCfg *CMDiscoveryNextResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateCfg)...)
	if resp.Diagnostics.HasError() { // coverage-ignore
		return
	}
	id := stateCfg.Id.ValueString()
	tflog.Info(ctx, fmt.Sprintf("Reading Instance Info : %+v", id))

	deviceInfo, err := r.client.GetDeviceInfoByID(id)
	if err != nil { // coverage-ignore
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to Read Instance Info, got error: %s", err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Instance Info : %+v", deviceInfo))

	r.DiscoveryNextResourceModeltoState(ctx, deviceInfo, stateCfg)

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateCfg)...)
}

func (r *CMDiscoveryNextResource) DiscoveryNextResourceModeltoState(ctx context.Context, respData interface{}, data *CMDiscoveryNextResourceModel) {
	tflog.Debug(ctx, fmt.Sprintf("respData  %+v", respData))
	data.Address = types.StringValue(respData.(map[string]interface{})["address"].(string))
}

func (r *CMDiscoveryNextResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var resCfg *CMDiscoveryNextResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &resCfg)...)

	if resp.Diagnostics.HasError() { // coverage-ignore
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[UPDATE] Updating Device Provider: %s", resCfg.Id.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &resCfg)...)
}

func (r *CMDiscoveryNextResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var stateCfg *CMDiscoveryNextResourceModel
	if resp.Diagnostics.HasError() { // coverage-ignore
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &stateCfg)...)
	id := stateCfg.Id.ValueString()
	err := r.client.DeleteDevice(id)
	if err != nil { // coverage-ignore
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to Delete Instance, got error: %s", err))
		return
	}
	stateCfg.Id = types.StringValue("")
}

func (r *CMDiscoveryNextResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getCMDiscoveryNextConfig(ctx context.Context, data *CMDiscoveryNextResourceModel) *bigipnextsdk.DiscoverInstanceRequest {
	discoverInstanceReq := &bigipnextsdk.DiscoverInstanceRequest{}
	tflog.Info(ctx, fmt.Sprintf("discoverInstanceReq:%+v\n", discoverInstanceReq))
	discoverInstanceReq.Address = data.Address.ValueString()
	discoverInstanceReq.Port = int(data.Port.ValueInt64())
	discoverInstanceReq.DeviceUser = data.DeviceUser.ValueString()
	discoverInstanceReq.DevicePassword = data.DevicePassword.ValueString()
	discoverInstanceReq.ManagementUser = data.ManagementUser.ValueString()
	discoverInstanceReq.ManagementPassword = data.ManagementPassword.ValueString()
	return discoverInstanceReq
}
