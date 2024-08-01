package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	bigipnextsdk "gitswarm.f5net.com/terraform-providers/bigipnext"
)

var (
	_ resource.Resource                = &NextDeployVmwareResource{}
	_ resource.ResourceWithImportState = &NextDeployVmwareResource{}
)

func NewNextDeployVmwareResource() resource.Resource {
	return &NextDeployVmwareResource{}
}

type NextDeployVmwareResource struct {
	client *bigipnextsdk.BigipNextCM
}

type NextDeployVmwareResourceModel struct {
	VsphereProvider types.Object     `tfsdk:"vsphere_provider"`
	Instance        types.Object     `tfsdk:"instance"`
	Timeout         types.Int64      `tfsdk:"timeout"`
	DnsServers      types.List       `tfsdk:"dns_servers"`
	NtpServers      types.List       `tfsdk:"ntp_servers"`
	Id              types.String     `tfsdk:"id"`
	ProviderId      types.String     `tfsdk:"provider_id"`
	L1networks      []L1networkModel `tfsdk:"l1_networks"`
}

type VsphereProviderModel struct {
	ProviderName     types.String `tfsdk:"provider_name"`
	DatacenterName   types.String `tfsdk:"datacenter_name"`
	ClusterName      types.String `tfsdk:"cluster_name"`
	DatastoreName    types.String `tfsdk:"datastore_name"`
	ResourcepoolName types.String `tfsdk:"resource_pool_name"`
	ContentLibrary   types.String `tfsdk:"content_library"`
	TemplateName     types.String `tfsdk:"vm_template_name"`
}

type InstanceModel struct {
	InstanceHostname          types.String `tfsdk:"instance_hostname"`
	MgmtNetworkname           types.String `tfsdk:"mgmt_network_name"`
	MgmtAddress               types.String `tfsdk:"mgmt_address"`
	MgmtPrefix                types.Int64  `tfsdk:"mgmt_prefix"`
	MgmtGateway               types.String `tfsdk:"mgmt_gateway"`
	MgmtUser                  types.String `tfsdk:"mgmt_user"`
	MgmtPassword              types.String `tfsdk:"mgmt_password"`
	Cpu                       types.Int64  `tfsdk:"cpu"`
	Memory                    types.Int64  `tfsdk:"memory"`
	ExternalNetworkname       types.String `tfsdk:"external_network_name"`
	InternalNetworkname       types.String `tfsdk:"internal_network_name"`
	HacontrolplaneNetworkname types.String `tfsdk:"ha_control_plane_network_name"`
	HadataplaneNetworkname    types.String `tfsdk:"ha_data_plane_network_name"`
}

type L1networkModel struct {
	Name  types.String `tfsdk:"name"`
	Vlans []VlanModel  `tfsdk:"vlans"`
}

type VlanModel struct {
	VlanName types.String `tfsdk:"vlan_name"`
	SelfIps  types.List   `tfsdk:"self_ips"`
	VlanTag  types.Int64  `tfsdk:"vlan_tag"`
}

func (r *NextDeployVmwareResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cm_deploy_vmware"
}

func (r *NextDeployVmwareResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource used to Deploy a new VM with a BIG-IP Next image template in the Vsphere environment",
		Attributes: map[string]schema.Attribute{
			"vsphere_provider": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"provider_name": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Name of provider to be used for deploying VMs",
					},
					"datacenter_name": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The vSphere datacenter to create the VMs",
					},
					"cluster_name": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The vSphere cluster to create the VMs",
					},
					"datastore_name": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The vSphere datastore to create the VMs",
					},
					"resource_pool_name": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The vSphere resource pool to create the VMs",
					},
					"content_library": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The vSphere Content Library from where `vm_template_name` can be used to create the VMs",
					},
					"vm_template_name": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The vSphere VM template name to create the VMs",
					},
				},
				Required: true,
			},
			"instance": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"instance_hostname": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Name of VM Deployed on vSphere,it should be `unique` string value",
					},
					"mgmt_network_name": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "The management network name to be used for deployed BIG-IP Next Instance in vSphere",
					},
					"mgmt_address": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Management address to be used for deployed BIG-IP Next Instance in vsphere",
					},
					"mgmt_prefix": schema.Int64Attribute{
						Required:            true,
						MarkdownDescription: "Management address prefix to be used for deployed BIG-IP Next instance in vsphere.",
					},
					"mgmt_gateway": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Management gateway address to be used for deployed BIG-IP Next instance in vsphere",
					},
					"mgmt_user": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Management username of deployed BIG-IP Next instance in vpshere",
					},
					"mgmt_password": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Management password of deployed BIG-IP Next instance in vpshere",
						Sensitive:           true,
					},
					"cpu": schema.Int64Attribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The number of virtual processor cores to configure on the VM.Default is `8`",
						Default:             int64default.StaticInt64(8),
					},
					"memory": schema.Int64Attribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The amount of memory in MB to configure on the VM.Default is `16384`.",
						Default:             int64default.StaticInt64(16384),
					},
					"external_network_name": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Name of the vSphere network to use as the BIG-IP Next external network",
					},
					"internal_network_name": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Name of the vSphere network to use as the BIG-IP Next internal network",
					},
					"ha_control_plane_network_name": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Name of the vSphere network to use as the BIG-IP Next HA Control plane network",
					},
					"ha_data_plane_network_name": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Name of the vSphere network to use as the BIG-IP Next HA data plane network",
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
			"dns_servers": schema.ListAttribute{
				MarkdownDescription: "List of DNS servers to assign to each deployed instance",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"ntp_servers": schema.ListAttribute{
				MarkdownDescription: "List of NTP servers to assign to each deployed instance",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"l1_networks": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Name of l1Newwork to assign to deployed instance",
						},
						"vlans": schema.ListNestedAttribute{
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"vlan_name": schema.StringAttribute{
										Required:            true,
										MarkdownDescription: "Name of vlan to be mapped for l1Network",
									},
									"vlan_tag": schema.Int64Attribute{
										MarkdownDescription: "Vlan tag to be mapped for l1Network.",
										Required:            true,
									},
									"self_ips": schema.ListAttribute{
										ElementType:         types.StringType,
										Required:            true,
										MarkdownDescription: "List of self ips to be mapped for l1Network",
									},
								},
							},
							Required:            true,
							MarkdownDescription: "List of vlans to be mapped for l1Network,each vlan is a block of attributes like vlan_name,vlan_tag,self_ips",
						},
					},
				},
				Optional:            true,
				MarkdownDescription: "List of l1networks to assign to deployed instance, each l1network is a block of attributes like name, vlans",
			},
			"provider_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique Identifier for the vpshere provider",
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

func (r *NextDeployVmwareResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client, resp.Diagnostics = toBigipNextCMProvider(req.ProviderData)
}

func (r *NextDeployVmwareResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var resCfg *NextDeployVmwareResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &resCfg)...)
	if resp.Diagnostics.HasError() { // coverage-ignore
		return
	}
	var providerModel VsphereProviderModel
	diag := resCfg.VsphereProvider.As(ctx, &providerModel, basetypes.ObjectAsOptions{})
	if diag.HasError() { // coverage-ignore
		return
	}
	var instanceModel InstanceModel
	diag = resCfg.Instance.As(ctx, &instanceModel, basetypes.ObjectAsOptions{})
	if diag.HasError() { // coverage-ignore
		return
	}
	providerID, err := r.client.GetDeviceProviderIDByHostname(providerModel.ProviderName.ValueString())
	if err != nil { // coverage-ignore
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to get provider ID:, got error: %s", err))
		return
	}
	resCfg.ProviderId = types.StringValue(providerID.(string))
	providerConfig := instanceConfig(ctx, resCfg)
	tflog.Info(ctx, fmt.Sprintf("[CREATE] Deploy Next Instance:%+v\n", providerConfig.Parameters.Hostname))
	respData, err := r.client.PostDeviceInstance(providerConfig, int(resCfg.Timeout.ValueInt64()))
	if err != nil { // coverage-ignore
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to Deploy Instance, got error: %s", err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[CREATE] respData ID:%+v\n", respData))
	deviceDetails, err := r.client.GetDeviceIdByHostname(providerConfig.Parameters.Hostname)
	if err != nil { // coverage-ignore
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to Read Device Info, got error: %s", err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[CREATE] respData ID:%+v\n", deviceDetails))
	resCfg.Id = types.StringValue(providerConfig.Parameters.Hostname)
	resp.Diagnostics.Append(resp.State.Set(ctx, resCfg)...)
}

func (r *NextDeployVmwareResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var stateCfg *NextDeployVmwareResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateCfg)...)
	if resp.Diagnostics.HasError() { // coverage-ignore
		return
	}
	id := stateCfg.Id.ValueString()
	tflog.Info(ctx, fmt.Sprintf("Reading Device info for : %+v", id))

	deviceDetails, err := r.client.GetDeviceIdByHostname(stateCfg.Id.ValueString())
	if err != nil { // coverage-ignore
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to Read Device Info, got error: %s", err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Device Info : %+v", *deviceDetails))
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateCfg)...)
}

func (r *NextDeployVmwareResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var resCfg *NextDeployVmwareResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &resCfg)...)

	if resp.Diagnostics.HasError() { // coverage-ignore
		return
	}
	providerConfig := instanceConfig(ctx, resCfg)
	tflog.Info(ctx, fmt.Sprintf("[UPDATE] Device Provider config:%+v\n", providerConfig))
	resp.Diagnostics.Append(resp.State.Set(ctx, &resCfg)...)
}

func (r *NextDeployVmwareResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var stateCfg *NextDeployVmwareResourceModel
	if resp.Diagnostics.HasError() { // coverage-ignore
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &stateCfg)...)
	id := stateCfg.Id.ValueString()

	tflog.Info(ctx, fmt.Sprintf("[DELETE] Deleting Instance from CM : %s", id))
	deviceDetails, err := r.client.GetDeviceIdByHostname(stateCfg.Id.ValueString())
	if err != nil { // coverage-ignore
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to Read Device Info, got error: %s", err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Device Info : %+v", *deviceDetails))

	err = r.client.DeleteDevice(*deviceDetails)
	if err != nil { // coverage-ignore
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to Delete Instance, got error:%s", err))
		return
	}
	stateCfg.Id = types.StringValue("")
}

func (r *NextDeployVmwareResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func instanceConfig(ctx context.Context, data *NextDeployVmwareResourceModel) *bigipnextsdk.CMReqDeviceInstance {
	var deployConfig bigipnextsdk.CMReqDeviceInstance
	deployConfig.TemplateName = "default-standalone-ve"
	var providerModel VsphereProviderModel
	diag := data.VsphereProvider.As(ctx, &providerModel, basetypes.ObjectAsOptions{})
	if diag.HasError() { // coverage-ignore
		tflog.Error(ctx, fmt.Sprintf("VsphereProviderModel diag Error: %+v", diag.Errors()))
	}
	var instanceModel InstanceModel
	diag = data.Instance.As(ctx, &instanceModel, basetypes.ObjectAsOptions{})
	if diag.HasError() { // coverage-ignore
		tflog.Error(ctx, fmt.Sprintf("InstanceModel diag Error: %+v", diag.Errors()))
	}
	var cmReqDeviceInstance bigipnextsdk.CMReqDeviceInstance
	cmReqDeviceInstance.TemplateName = "default-standalone-ve"
	cmReqDeviceInstance.Parameters.Hostname = (instanceModel.InstanceHostname).ValueString()
	cmReqDeviceInstance.Parameters.ManagementAddress = (instanceModel.MgmtAddress).ValueString()
	cmReqDeviceInstance.Parameters.ManagementNetworkWidth = int((instanceModel.MgmtPrefix).ValueInt64())
	cmReqDeviceInstance.Parameters.DefaultGateway = (instanceModel.MgmtGateway).ValueString()
	cmReqDeviceInstance.Parameters.ManagementCredentialsUsername = (instanceModel.MgmtUser).ValueString()
	cmReqDeviceInstance.Parameters.ManagementCredentialsPassword = (instanceModel.MgmtPassword).ValueString()
	cmReqDeviceInstance.Parameters.InstanceOneTimePassword = (instanceModel.MgmtPassword).ValueString()

	var dnsServ []string
	for _, val := range data.DnsServers.Elements() {
		var ss string
		_ = json.Unmarshal([]byte(val.String()), &ss)
		dnsServ = append(dnsServ, ss)
	}
	cmReqDeviceInstance.Parameters.DnsServers = dnsServ
	var ntpServ []string
	for _, val := range data.NtpServers.Elements() {
		var ss string
		_ = json.Unmarshal([]byte(val.String()), &ss)
		ntpServ = append(ntpServ, ss)
	}
	cmReqDeviceInstance.Parameters.NtpServers = ntpServ
	cmReqDeviceInstance.Parameters.InstantiationProvider = append(cmReqDeviceInstance.Parameters.InstantiationProvider, bigipnextsdk.CMReqInstantiationProvider{
		Id:   data.ProviderId.ValueString(),
		Name: providerModel.ProviderName.ValueString(),
		Type: "vsphere",
	})
	var l1Networks bigipnextsdk.CMReqL1Networks
	// var l1networkModel L1networksModel
	// interfaces := []string{"1.1", "1.2", "1.3"}
	for index, val := range data.L1networks {
		tflog.Info(ctx, fmt.Sprintf("val : %+v", val))
		l1Networks.Name = val.Name.ValueString()
		l1Networks.L1Link.Name = fmt.Sprintf("1.%d", index+1)
		l1Networks.L1Link.LinkType = "Interface"
		l1Networks.Vlans = make([]bigipnextsdk.CMReqVlans, len(val.Vlans))
		for index, vlan := range val.Vlans {
			tflog.Info(ctx, fmt.Sprintf("vlan name : %+v", vlan.VlanName.ValueString()))
			l1Networks.Vlans[index].Name = vlan.VlanName.ValueString()
			l1Networks.Vlans[index].DefaultVrf = true
			l1Networks.Vlans[index].Tag = int(vlan.VlanTag.ValueInt64())
			elements := make([]types.String, 0, len(vlan.SelfIps.Elements()))
			diags := vlan.SelfIps.ElementsAs(ctx, &elements, false)
			if diags.HasError() { // coverage-ignore
				tflog.Error(ctx, fmt.Sprintf("SelfIps diag Error: %+v", diags.Errors()))
			}
			l1Networks.Vlans[index].SelfIps = make([]bigipnextsdk.CMReqSelfIps, len(elements))
			for ii, selfip := range elements {
				l1Networks.Vlans[index].SelfIps[ii].Address = selfip.ValueString()
				l1Networks.Vlans[index].SelfIps[ii].DeviceName = fmt.Sprintf("%s-%s-%d", val.Name.ValueString(), vlan.VlanName.ValueString(), ii)
			}
		}
		cmReqDeviceInstance.Parameters.L1Networks = append(cmReqDeviceInstance.Parameters.L1Networks, l1Networks)
	}

	tflog.Info(ctx, fmt.Sprintf("l1Networks : %+v", l1Networks))

	cmReqDeviceInstance.Parameters.VSphereProperties = append(cmReqDeviceInstance.Parameters.VSphereProperties, bigipnextsdk.CMReqVsphereProperties{
		NumCpus:               int(instanceModel.Cpu.ValueInt64()),
		Memory:                int(instanceModel.Memory.ValueInt64()),
		DatacenterName:        providerModel.DatacenterName.ValueString(),
		ClusterName:           providerModel.ClusterName.ValueString(),
		DatastoreName:         providerModel.DatastoreName.ValueString(),
		ResourcePoolName:      providerModel.ResourcepoolName.ValueString(),
		VsphereContentLibrary: providerModel.ContentLibrary.ValueString(),
		VmTemplateName:        providerModel.TemplateName.ValueString(),
	})
	cmReqDeviceInstance.Parameters.VsphereNetworkAdapterSettings = append(cmReqDeviceInstance.Parameters.VsphereNetworkAdapterSettings, bigipnextsdk.CMReqVsphereNetworkAdapterSettings{
		MgmtNetworkName:           instanceModel.MgmtNetworkname.ValueString(),
		InternalNetworkName:       instanceModel.InternalNetworkname.ValueString(),
		ExternalNetworkName:       instanceModel.ExternalNetworkname.ValueString(),
		HaDataPlaneNetworkName:    instanceModel.HadataplaneNetworkname.ValueString(),
		HaControlPlaneNetworkName: instanceModel.HacontrolplaneNetworkname.ValueString(),
	})
	tflog.Info(ctx, fmt.Sprintf("cmReqDeviceInstance : %+v", cmReqDeviceInstance))
	return &cmReqDeviceInstance
}
