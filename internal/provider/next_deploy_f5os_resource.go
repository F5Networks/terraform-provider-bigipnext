package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	bigipnextsdk "gitswarm.f5net.com/terraform-providers/bigipnext"
)

var (
	_ resource.Resource                = &NextDeployF5osResource{}
	_ resource.ResourceWithImportState = &NextDeployF5osResource{}
)

func NewNextDeployF5osResource() resource.Resource {
	return &NextDeployF5osResource{}
}

type NextDeployF5osResource struct {
	client *bigipnextsdk.BigipNextCM
}

type NextDeployF5osResourceModel struct {
	F5OSProvider types.Object `tfsdk:"f5os_provider"`
	F5OSInstance types.Object `tfsdk:"instance"`
	Timeout      types.Int64  `tfsdk:"timeout"`
	Id           types.String `tfsdk:"id"`
	ProviderId   types.String `tfsdk:"provider_id"`
}

type F5OSProviderModel struct {
	ProviderName types.String `tfsdk:"provider_name"`
	ProviderType types.String `tfsdk:"provider_type"`
}

type F5OSInstanceModel struct {
	InstanceHostname     types.String `tfsdk:"instance_hostname"`
	MgmtAddress          types.String `tfsdk:"management_address"`
	MgmtPrefix           types.Int64  `tfsdk:"management_prefix"`
	MgmtGateway          types.String `tfsdk:"management_gateway"`
	MgmtUser             types.String `tfsdk:"management_user"`
	MgmtPassword         types.String `tfsdk:"management_password"`
	VlanIDs              types.List   `tfsdk:"vlan_ids"`
	CpuCores             types.Int64  `tfsdk:"cpu_cores"`
	DiskSize             types.Int64  `tfsdk:"disk_size"`
	TenantImageName      types.String `tfsdk:"tenant_image_name"`
	TenantDeploymentFile types.String `tfsdk:"tenant_deployment_file"`
}

func (r *NextDeployF5osResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cm_deploy_f5os"
}

func (r *NextDeployF5osResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource used to Deploy New `BIG-IP-Next` Instance on F5OS platforms like rSeries/velos using Next Image available on F5OS Platform",
		Attributes: map[string]schema.Attribute{
			"f5os_provider": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"provider_name": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Name of F5OS provider to be used for deploying Instances",
					},
					"provider_type": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The Type of F5OS provider(rseries/velos)",
						Validators: []validator.String{
							stringvalidator.OneOf([]string{"rseries", "velos"}...),
						},
					},
				},
				Required: true,
			},
			"instance": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"instance_hostname": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Name of BIG-IP-Next Instance to be Deployed on F5OS(velos/rSeries),it should be `unique` string value",
					},
					"tenant_image_name": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Name of tenant image to be used to deployinstance in F5OS Provider",
					},
					"tenant_deployment_file": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Name of the tenant deployment file to be used to deploy instance in F5OS Provider",
					},
					"management_address": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Management address to be used for deployed BIG-IP Next instance in F5OS Provider",
					},
					"management_prefix": schema.Int64Attribute{
						Required:            true,
						MarkdownDescription: "Management address prefix to be used for deployed BIG-IP Next instance in F5OS Provider.",
					},
					"management_gateway": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Management gateway address to be used for deployed BIG-IP Next instance in F5OS Provider",
					},
					"management_user": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Management username of deployed BIG-IP Next instance in F5OS Provider",
					},
					"management_password": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Management password of deployed BIG-IP Next instance in F5OS Provider",
						Sensitive:           true,
					},
					"vlan_ids": schema.ListAttribute{
						MarkdownDescription: "List of integers. Specifies on which blades nodes the tenants are deployed.\nRequired for create operations.\nFor single blade platforms like rSeries only the value of 1 should be provided.",
						Optional:            true,
						Computed:            true,
						ElementType:         types.Int64Type,
					},
					"cpu_cores": schema.Int64Attribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The number of virtual processor cores to configure on the BIG-IP-Next Instance.Default is `4`.",
						Default:             int64default.StaticInt64(4),
					},
					"disk_size": schema.Int64Attribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The amount of disk size in GigBytes to configure on the BIG-IP-Next Instance.Default is `30`.",
						Default:             int64default.StaticInt64(30),
					},
				},
				Required: true,
			},
			"timeout": schema.Int64Attribute{
				MarkdownDescription: "The number of seconds to wait for instance deployment to finish.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(900),
			},
			"provider_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique Identifier for the F5OS provider",
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

func (r *NextDeployF5osResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client, resp.Diagnostics = toBigipNextCMProvider(req.ProviderData)
}

func (r *NextDeployF5osResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var resCfg *NextDeployF5osResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &resCfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var providerModel F5OSProviderModel
	diag := resCfg.F5OSProvider.As(ctx, &providerModel, basetypes.ObjectAsOptions{})
	if diag.HasError() {
		return
	}
	var instanceModel F5OSInstanceModel
	diag = resCfg.F5OSInstance.As(ctx, &instanceModel, basetypes.ObjectAsOptions{})
	if diag.HasError() {
		return
	}
	providerID, err := r.client.GetDeviceProviderIDByHostname(providerModel.ProviderName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to get provider ID:, got error: %s", err))
		return
	}
	resCfg.ProviderId = types.StringValue(providerID.(string))
	providerConfig := f5osRseriesConfig(ctx, resCfg)
	tflog.Info(ctx, fmt.Sprintf("[CREATE] Deploy Next Instance:%+v\n", providerConfig.Parameters.Hostname))
	respData, err := r.client.PostDeviceInstance(providerConfig, int(resCfg.Timeout.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to Deploy Instance, got error: %s", err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[CREATE] respData ID:%+v\n", respData))
	resCfg.Id = types.StringValue(providerConfig.Parameters.Hostname)
	resp.Diagnostics.Append(resp.State.Set(ctx, resCfg)...)
}

func (r *NextDeployF5osResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var stateCfg *NextDeployF5osResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateCfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id := stateCfg.Id.ValueString()
	tflog.Info(ctx, fmt.Sprintf("Reading Device info for : %+v", id))

	deviceDetails, err := r.client.GetDeviceIdByHostname(stateCfg.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to Read Device Info, got error: %s", err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Device Info : %+v", *deviceDetails))
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateCfg)...)
}

func (r *NextDeployF5osResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var resCfg *NextDeployF5osResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &resCfg)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "[UPDATE] Updating Next Instance deployment is not supported")

	resp.Diagnostics.Append(resp.State.Set(ctx, &resCfg)...)
}

func (r *NextDeployF5osResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var stateCfg *NextDeployF5osResourceModel
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &stateCfg)...)
	id := stateCfg.Id.ValueString()

	tflog.Info(ctx, fmt.Sprintf("[DELETE] Deleting Instance from CM : %s", id))
	deviceDetails, err := r.client.GetDeviceIdByHostname(stateCfg.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to Read Device Info, got error: %s", err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Device Info : %+v", *deviceDetails))

	err = r.client.DeleteDevice(*deviceDetails)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to Delete Instance, got error: %s", err))
		return
	}
	stateCfg.Id = types.StringValue("")
}

func (r *NextDeployF5osResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func f5osRseriesConfig(ctx context.Context, data *NextDeployF5osResourceModel) *bigipnextsdk.CMReqDeviceInstance {
	var deployConfig bigipnextsdk.CMReqDeviceInstance
	deployConfig.TemplateName = "default-standalone-rseries"
	var providerModel F5OSProviderModel
	diag := data.F5OSProvider.As(ctx, &providerModel, basetypes.ObjectAsOptions{})
	if diag.HasError() {
		tflog.Error(ctx, fmt.Sprintf("F5OSProvider diag Error: %+v", diag.Errors()))
	}
	var instanceModel F5OSInstanceModel
	diag = data.F5OSInstance.As(ctx, &instanceModel, basetypes.ObjectAsOptions{})
	if diag.HasError() {
		tflog.Error(ctx, fmt.Sprintf("F5OSInstanceModel diag Error: %+v", diag.Errors()))
	}
	var cmReqDeviceInstance bigipnextsdk.CMReqDeviceInstance
	cmReqDeviceInstance.TemplateName = "default-standalone-rseries"
	cmReqDeviceInstance.Parameters.Hostname = (instanceModel.InstanceHostname).ValueString()
	cmReqDeviceInstance.Parameters.ManagementAddress = (instanceModel.MgmtAddress).ValueString()
	cmReqDeviceInstance.Parameters.ManagementNetworkWidth = int((instanceModel.MgmtPrefix).ValueInt64())
	cmReqDeviceInstance.Parameters.DefaultGateway = (instanceModel.MgmtGateway).ValueString()
	cmReqDeviceInstance.Parameters.ManagementCredentialsUsername = (instanceModel.MgmtUser).ValueString()
	cmReqDeviceInstance.Parameters.ManagementCredentialsPassword = (instanceModel.MgmtPassword).ValueString()
	cmReqDeviceInstance.Parameters.InstanceOneTimePassword = (instanceModel.MgmtPassword).ValueString()

	var vlanIds []int
	for _, val := range instanceModel.VlanIDs.Elements() {
		var ss int
		_ = json.Unmarshal([]byte(val.String()), &ss)
		vlanIds = append(vlanIds, ss)
	}

	// var dnsServ []string
	// for _, val := range data.DnsServers.Elements() {
	// 	var ss string
	// 	_ = json.Unmarshal([]byte(val.String()), &ss)
	// 	dnsServ = append(dnsServ, ss)
	// }
	// cmReqDeviceInstance.Parameters.DnsServers = dnsServ
	// var ntpServ []string
	// for _, val := range data.NtpServers.Elements() {
	// 	var ss string
	// 	_ = json.Unmarshal([]byte(val.String()), &ss)
	// 	ntpServ = append(ntpServ, ss)
	// }
	// cmReqDeviceInstance.Parameters.NtpServers = ntpServ

	cmReqDeviceInstance.Parameters.InstantiationProvider = append(cmReqDeviceInstance.Parameters.InstantiationProvider, bigipnextsdk.CMReqInstantiationProvider{
		Id:   data.ProviderId.ValueString(),
		Name: providerModel.ProviderName.ValueString(),
		Type: "rseries",
	})
	cmReqDeviceInstance.Parameters.RseriesProperties = append(cmReqDeviceInstance.Parameters.RseriesProperties, bigipnextsdk.CMReqRseriesProperties{

		TenantImageName:      instanceModel.TenantImageName.ValueString(),
		TenantDeploymentFile: instanceModel.TenantDeploymentFile.ValueString(),
		CpuCores:             int(instanceModel.CpuCores.ValueInt64()),
		DiskSize:             int(instanceModel.DiskSize.ValueInt64()),
		VlanIds:              vlanIds,
	})
	tflog.Info(ctx, fmt.Sprintf("cmReqDeviceInstance : %+v", cmReqDeviceInstance))
	return &cmReqDeviceInstance
}
