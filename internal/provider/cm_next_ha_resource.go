package provider

import (
	"context"
	"fmt"

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
	_ resource.Resource                = &NextHAResource{}
	_ resource.ResourceWithImportState = &NextHAResource{}
)

func NewNextHAResource() resource.Resource {
	return &NextHAResource{}
}

type NextHAResource struct {
	client *bigipnextsdk.BigipNextCM
}

type NextHAResourceModel struct {
	HaName                    types.String `tfsdk:"ha_name"`
	HaIP                      types.String `tfsdk:"ha_ip"`
	ActiveNodeIp              types.String `tfsdk:"active_node_ip"`
	StandbyNodeIp             types.String `tfsdk:"standby_node_ip"`
	ControlplaneVlan          types.String `tfsdk:"control_plane_vlan"`
	ControlplaneVlantag       types.Int64  `tfsdk:"control_plane_vlan_tag"`
	DataplaneVlan             types.String `tfsdk:"data_plane_vlan"`
	DataplaneVlantag          types.Int64  `tfsdk:"data_plane_vlan_tag"`
	ActiveNodeControlplaneIp  types.String `tfsdk:"active_node_control_plane_ip"`
	StandbyNodeControlplaneIp types.String `tfsdk:"standby_node_control_plane_ip"`
	ActiveNodeDataplaneIp     types.String `tfsdk:"active_node_data_plane_ip"`
	StandbyNodeDataplaneIp    types.String `tfsdk:"standby_node_data_plane_ip"`
	Timeout                   types.Int64  `tfsdk:"timeout"`
	DeviceId                  types.String `tfsdk:"device_id"`
	Id                        types.String `tfsdk:"id"`
}

func (r *NextHAResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cm_next_ha"
}

func (r *NextHAResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Configure High Availability for NEXT instances managed by CM",
		Attributes: map[string]schema.Attribute{
			"ha_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the High Availability (HA) cluster.The name must be unique and cannot be changed after the cluster is created.",
			},
			"ha_ip": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The desired management IP of the HA cluster.",
			},
			"active_node_ip": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The designated active Next instance management IP.",
			},
			"standby_node_ip": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The designated standby Next instance management IP.",
			},
			"active_node_control_plane_ip": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The HA control plane IP address on active node.",
			},
			"standby_node_control_plane_ip": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The HA control plane IP address on standby node.",
			},
			"active_node_data_plane_ip": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The HA data plane IP address on active node.",
			},
			"standby_node_data_plane_ip": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The HA data plane IP address on standby node.",
			},
			"control_plane_vlan": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The VLAN for the HA control plane.",
			},
			"data_plane_vlan": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The VLAN for the HA data plane.",
			},
			"control_plane_vlan_tag": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The tag for the HA control plane VLAN.",
			},
			"data_plane_vlan_tag": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "The tag for the HA control plane VLAN.",
			},
			"timeout": schema.Int64Attribute{
				MarkdownDescription: "The amount of time to wait for the HA creation task to finish, in seconds.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(900),
			},
			"device_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "HA Device ID",
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

func (r *NextHAResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client, resp.Diagnostics = toBigipNextCMProvider(req.ProviderData)
}

func (r *NextHAResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var resCfg *NextHAResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &resCfg)...)
	if resp.Diagnostics.HasError() { // coverage-ignore
		return
	}
	// get activeNodeID by IP
	activeNodeID, err := r.client.GetDeviceIdByIp(resCfg.ActiveNodeIp.ValueString())
	if err != nil { // coverage-ignore
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to Read Active Device Info, got error: %s", err))
		return
	}
	standbyNodeID, err := r.client.GetDeviceIdByIp(resCfg.StandbyNodeIp.ValueString())
	if err != nil { // coverage-ignore
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to Read Standby Device Info, got error: %s", err))
		return
	}

	// resCfg.ProviderId = types.StringValue(providerID.(string))
	haDeployConfig := haConfig(ctx, *activeNodeID, *standbyNodeID, resCfg)
	tflog.Info(ctx, fmt.Sprintf("[CREATE] Deploy HA :%+v\n", haDeployConfig.ClusterName))
	respData, err := r.client.PostDeviceHA(*activeNodeID, haDeployConfig, int(resCfg.Timeout.ValueInt64()))
	if err != nil { // coverage-ignore
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to Deploy Instance, got error: %s", err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[CREATE] respData ID:%+v\n", respData))
	resCfg.DeviceId = types.StringValue(respData.(map[string]interface{})["id"].(string))
	resCfg.Id = types.StringValue(haDeployConfig.ClusterManagementIP)
	resp.Diagnostics.Append(resp.State.Set(ctx, resCfg)...)
}

func (r *NextHAResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var stateCfg *NextHAResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateCfg)...)
	if resp.Diagnostics.HasError() { // coverage-ignore
		return
	}
	id := stateCfg.Id.ValueString()
	tflog.Info(ctx, fmt.Sprintf("Reading Device info for : %+v", id))

	haNodeInfo, err := r.client.GetDeviceInfoByIp(id)
	if err != nil { // coverage-ignore
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to Read HA Device Info, got error: %s", err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Reading HA Device info for : %+v", haNodeInfo))
	// map[_links:map[self:map[href:/v1/inventory?filter=address+eq+%2710.146.168.20%27/23254958-db28-4d10-b42f-ff58bc16228d]] address:10.146.168.20 certificate_validated:2023-11-27T09:58:27.605586Z certificate_validation_error:tls: failed to verify certificate: x509: cannot validate certificate for 10.146.194.141 because it doesn't contain any IP SANs certificate_validity:false hostname:raviecosyshydha id:23254958-db28-4d10-b42f-ff58bc16228d mode:HA platform_name:VMware platform_type:VE port:5443 version:20.0.1-2.139.10+0.0.136]

	// check if mode is HA from above response map
	if haNodeInfo.(map[string]interface{})["mode"] != "HA" { // coverage-ignore
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to Read HA Device Info, got error: %s", err))
		return
	}
	stateCfg.DeviceId = types.StringValue(haNodeInfo.(map[string]interface{})["id"].(string))
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateCfg)...)
}

func (r *NextHAResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var resCfg *NextHAResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &resCfg)...)

	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &resCfg)...)
}

func (r *NextHAResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var stateCfg *NextHAResourceModel
	if resp.Diagnostics.HasError() { // coverage-ignore
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &stateCfg)...)
	id := stateCfg.Id.ValueString()
	tflog.Info(ctx, fmt.Sprintf("[DELETE] Deleting Instance from CM : %s", id))
	deviceID := stateCfg.DeviceId.ValueString()
	err := r.client.DeleteDevice(deviceID)
	if err != nil { // coverage-ignore
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to Delete Instance, got error: %s", err))
		return
	}
	stateCfg.Id = types.StringValue("")
}

func (r *NextHAResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func haConfig(ctx context.Context, activeNodeID, standbyNodeId string, data *NextHAResourceModel) *bigipnextsdk.CMReqDeviceHA {
	var deployConfig bigipnextsdk.CMReqDeviceHA
	deployConfig.AutoFailback = bool(false)
	deployConfig.ClusterManagementIP = data.HaIP.ValueString()
	deployConfig.ClusterName = data.HaName.ValueString()
	deployConfig.ControlPlaneVlan.Name = data.ControlplaneVlan.ValueString()
	deployConfig.ControlPlaneVlan.Tag = int(data.ControlplaneVlantag.ValueInt64())
	deployConfig.DataPlaneVlan.Name = data.DataplaneVlan.ValueString()
	deployConfig.DataPlaneVlan.Tag = int(data.DataplaneVlantag.ValueInt64())
	deployConfig.DataPlaneVlan.NetworkInterface = "1.3"
	var haNodes bigipnextsdk.CMReqHANode
	haNodes.Name = "active-node"
	haNodes.ControlPlaneAddress = data.ActiveNodeControlplaneIp.ValueString()
	haNodes.DataPlanePrimaryAddress = data.ActiveNodeDataplaneIp.ValueString()
	haNodes.DataPlaneSecondaryAddress = ""
	deployConfig.Nodes = append(deployConfig.Nodes, haNodes)
	haNodes.Name = "standby-node"
	haNodes.ControlPlaneAddress = data.StandbyNodeControlplaneIp.ValueString()
	haNodes.DataPlanePrimaryAddress = data.StandbyNodeDataplaneIp.ValueString()
	haNodes.DataPlaneSecondaryAddress = ""
	deployConfig.Nodes = append(deployConfig.Nodes, haNodes)
	deployConfig.StandbyInstanceID = standbyNodeId
	deployConfig.TrafficVlan = []interface{}{}
	tflog.Info(ctx, fmt.Sprintf("HA Deploy Config : %+v", deployConfig))
	return &deployConfig
}
