package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	bigipnextsdk "gitswarm.f5net.com/terraform-providers/bigipnext"
	// "strings"
	// "sync"
)

var (
	_ resource.Resource                = &NextCMDeviceProviderResource{}
	_ resource.ResourceWithImportState = &NextCMDeviceProviderResource{}
)

func NewNextCMDeviceProviderResource() resource.Resource {
	return &NextCMDeviceProviderResource{}
}

type NextCMDeviceProviderResource struct {
	client *bigipnextsdk.BigipNextCM
}

type NextCMDeviceProviderResourceModel struct {
	Type     types.String `tfsdk:"type"`
	Name     types.String `tfsdk:"name"`
	Address  types.String `tfsdk:"address"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	Id       types.String `tfsdk:"id"`
}

func (r *NextCMDeviceProviderResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cm_provider"
}

func (r *NextCMDeviceProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource used to manage(CRUD) providers on BIG-IP Next Central Manager",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Specifies the name of the provider on BIG-IP Next CM to create or manage",
			},
			"type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Specifies the type of provider,valid options are `RSERIES`/`VELOS`/`VSPHERE`",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"RSERIES", "VELOS", "VSPHERE"}...),
				},
			},
			"address": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The address of the provider to which Central Manager can connect to \n The address may be a hostname or an IP-address or IP-address:port \n The parameter must be specified when creating a new provider",
			},
			"username": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The username that the BIG-IP Next Central Manager uses when connecting with the specified provider",
			},
			"password": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "The password that the BIG-IP Next Central Manager uses when connecting with the specified provider",
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

func (r *NextCMDeviceProviderResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client, resp.Diagnostics = toBigipNextCMProvider(req.ProviderData)
}

func (r *NextCMDeviceProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var resCfg *NextCMDeviceProviderResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &resCfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[CREATE] NextCMDeviceProviderResource:%+v\n", resCfg.Name.ValueString()))

	providerConfig := getDeviceProvider(ctx, resCfg)

	tflog.Info(ctx, fmt.Sprintf("[CREATE] Device Provider config:%+v\n", providerConfig))
	respData, err := r.client.PostDeviceProvider(providerConfig)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to Create Certificate, got error: %s", err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[CREATE] respData ID:%+v\n", respData.Id))

	resCfg.Id = types.StringValue(respData.Id)
	resCfg.Type = types.StringValue(respData.Type)
	resCfg.Name = types.StringValue(respData.Name)
	resCfg.Address = types.StringValue(respData.Connection.Host)
	resCfg.Username = types.StringValue(respData.Connection.Authentication.Username)
	resp.Diagnostics.Append(resp.State.Set(ctx, resCfg)...)
}

func (r *NextCMDeviceProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var stateCfg *NextCMDeviceProviderResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateCfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id := stateCfg.Id.ValueString()
	tflog.Info(ctx, fmt.Sprintf("Reading Device Provider : %+v", stateCfg.Name.ValueString()))

	deviceProvider, err := r.client.GetDeviceProvider(id, stateCfg.Type.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to Read Device Provider, got error: %s", err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Device Provider : %+v", deviceProvider))
	r.DeviceProviderResourceModeltoState(ctx, deviceProvider, stateCfg)
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateCfg)...)
}

func (r *NextCMDeviceProviderResource) DeviceProviderResourceModeltoState(ctx context.Context, respData interface{}, data *NextCMDeviceProviderResourceModel) {
	tflog.Debug(ctx, fmt.Sprintf("respData  %+v", respData))
	data.Name = types.StringValue(respData.(*bigipnextsdk.DeviceProviderResponse).Name)
	data.Address = types.StringValue(respData.(*bigipnextsdk.DeviceProviderResponse).Connection.Host)
	data.Username = types.StringValue(respData.(*bigipnextsdk.DeviceProviderResponse).Connection.Authentication.Username)
	data.Type = types.StringValue(respData.(*bigipnextsdk.DeviceProviderResponse).Type)
}

func (r *NextCMDeviceProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var resCfg *NextCMDeviceProviderResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &resCfg)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[UPDATE] Updating Device Provider: %s", resCfg.Name.ValueString()))

	providerConfig := getDeviceProvider(ctx, resCfg)
	tflog.Info(ctx, fmt.Sprintf("[UPDATE] Device Provider config:%+v\n", providerConfig))

	respData, err := r.client.UpdateDeviceProvider(resCfg.Id.ValueString(), providerConfig)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to Update Device provider, got error: %s", err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[UPDATE] respData ID:%+v\n", respData.Id))
	resp.Diagnostics.Append(resp.State.Set(ctx, &resCfg)...)
}

func (r *NextCMDeviceProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var stateCfg *NextCMDeviceProviderResourceModel
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &stateCfg)...)
	id := stateCfg.Id.ValueString()

	tflog.Info(ctx, fmt.Sprintf("[DELETE] Deleting Device Provider : %s", id))
	_, err := r.client.DeleteDeviceProvider(id, stateCfg.Type.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to Delete Device Provider, got error: %s", err))
		return
	}
	stateCfg.Id = types.StringValue("")
}

func (r *NextCMDeviceProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getDeviceProvider(ctx context.Context, data *NextCMDeviceProviderResourceModel) *bigipnextsdk.DeviceProviderReq {
	var deviceProvider bigipnextsdk.DeviceProviderReq
	deviceProvider.Name = data.Name.ValueString()
	deviceProvider.Type = data.Type.ValueString()
	deviceProvider.Connection.Host = data.Address.ValueString()
	deviceProvider.Connection.Authentication.Type = "basic"
	deviceProvider.Connection.Authentication.Username = data.Username.ValueString()
	deviceProvider.Connection.Authentication.Password = data.Password.ValueString()
	tflog.Info(ctx, fmt.Sprintf("deviceProvider:%+v\n", deviceProvider))
	return &deviceProvider
}
