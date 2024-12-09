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
	_ resource.Resource                = &NextOnboardResource{}
	_ resource.ResourceWithImportState = &NextOnboardResource{}
)

func NewNextCMOnboardResource() resource.Resource {
	return &NextOnboardResource{}
}

type NextOnboardResource struct {
	client *bigipnextsdk.BigipNextCM
}

type NextOnboardResourceModel struct {
	DnsServers        types.List       `tfsdk:"dns_servers"`
	NtpServers        types.List       `tfsdk:"ntp_servers"`
	L1Networks        []L1NetworkModel `tfsdk:"l1_networks"`
	ManagementAddress types.String     `tfsdk:"management_address"`
	Timeout           types.Int64      `tfsdk:"timeout"`
	Id                types.String     `tfsdk:"id"`
}

type L1NetworkModel struct {
	Name   types.String `tfsdk:"name"`
	Vlans  []Vlan       `tfsdk:"vlans"`
	L1Link struct {
		Name     types.String `tfsdk:"name"`
		LinkType types.String `tfsdk:"link_type"`
	} `tfsdk:"l1_link"`
}

type Vlan struct {
	Tag     types.Float64 `tfsdk:"tag"`
	Name    types.String  `tfsdk:"name"`
	SelfIps []SelfIpModel `tfsdk:"self_ips"`
}

type SelfIpModel struct {
	Address    types.String `tfsdk:"address"`
	DeviceName types.String `tfsdk:"device_name"`
}

func (r *NextOnboardResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cm_instance_onboard"
}

func (r *NextOnboardResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Configure NEXT instances onboarding (DNS Servers, NTP Servers, L1 Networks along with VLANs and Self IPs).
		Note: Delete call is just for clearing the state, It makes no changes on the CM`,
		Attributes: map[string]schema.Attribute{
			"dns_servers": schema.ListAttribute{
				Required:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "DNS servers to be Added.",
			},
			"ntp_servers": schema.ListAttribute{
				Required:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "NTP servers to be Added.",
			},
			"management_address": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The desired management Address of the Deployment.",
			},
			"l1_networks": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "L1 network name.",
						},
						"vlans": schema.ListNestedAttribute{
							Optional: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Required:            true,
										MarkdownDescription: "L1 network name.",
									},
									"tag": schema.Float64Attribute{
										Required:            true,
										MarkdownDescription: "An unsigned 32-bit integer..",
									},
									"self_ips": schema.ListNestedAttribute{
										Optional: true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"address": schema.StringAttribute{
													Required:            true,
													MarkdownDescription: "An IPv4 or IPv6 prefix.",
												},
												"device_name": schema.StringAttribute{
													Required:            true,
													MarkdownDescription: "Specifies the node that this non-floating self-IP address belongs to.",
												},
											},
										},
									},
								},
							},
						},
						"l1_link": schema.SingleNestedAttribute{
							Required:            true,
							MarkdownDescription: "L1 link layer interface.",
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "A soft reference, by name, to an L1 link layer interface.",
								},
								"link_type": schema.StringAttribute{
									Required:            true,
									MarkdownDescription: "L1 Link type..",
								},
							},
						},
					},
				},
			},
			"timeout": schema.Int64Attribute{
				MarkdownDescription: "The number of seconds to wait for instance onboard to finish.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(600),
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

func (r *NextOnboardResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client, resp.Diagnostics = toBigipNextCMProvider(req.ProviderData)
}

func (r *NextOnboardResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var resCfg *NextOnboardResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &resCfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// get Instance by IP
	instanceId, err := r.client.GetDeviceIdByIp(resCfg.ManagementAddress.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to Read Instance Info, got error: %s", err))
		return
	}

	onboardInstanceConfig := onboardInstanceConfig(ctx, resCfg)
	respData, err := r.client.PatchDeviceInstance(*instanceId, onboardInstanceConfig, int(resCfg.Timeout.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to Onboard Instance, got error: %s", err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[CREATE] Instance Info :%+v\n", respData))
	resCfg.Id = types.StringValue(*instanceId)
	r.instanceModeltoState(ctx, respData, resCfg)
	resp.Diagnostics.Append(resp.State.Set(ctx, resCfg)...)
}

func (r *NextOnboardResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var stateCfg *NextOnboardResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateCfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id := stateCfg.Id.ValueString()
	tflog.Info(ctx, fmt.Sprintf("Reading Instance info for ID: %+v", id))
	instanceInfo, err := r.client.GetDeviceInfoByID(id, true)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to Read Instance Info, got error: %s", err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[READ] Instance info : %+v", instanceInfo))
	r.instanceModeltoState(ctx, instanceInfo, stateCfg)
	stateCfg.Id = types.StringValue(id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateCfg)...)
}

func (r *NextOnboardResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var resCfg *NextOnboardResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &resCfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	instanceId := resCfg.Id.ValueString()
	tflog.Info(ctx, fmt.Sprintf("[UPDATE] Updating Instance from CM : %s", instanceId))
	onboardInstanceConfig := onboardInstanceConfig(ctx, resCfg)
	respData, err := r.client.PatchDeviceInstance(instanceId, onboardInstanceConfig, int(resCfg.Timeout.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to Onboard Instance, got error: %s", err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[Update] Instance Info Updated :%+v\n", respData))
	resCfg.Id = types.StringValue(instanceId)
	r.instanceModeltoState(ctx, respData, resCfg)
	resp.Diagnostics.Append(resp.State.Set(ctx, resCfg)...)
}

func (r *NextOnboardResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var stateCfg *NextOnboardResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateCfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id := stateCfg.Id.ValueString()
	tflog.Info(ctx, fmt.Sprintf("[DELETE] Deleting Instance from CM : %s", id))
	// deviceID := stateCfg.DeviceId.ValueString()
	// err := r.client.DeleteDevice(deviceID)
	// if err != nil {
	// 	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to Delete Instance, got error: %s", err))
	// 	return
	// }
	stateCfg.Id = types.StringValue("")
}

func (r *NextOnboardResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func onboardInstanceConfig(ctx context.Context, data *NextOnboardResourceModel) *bigipnextsdk.CMReqDeviceInstance {
	var deployConfig bigipnextsdk.CMReqDeviceInstance
	deployConfig.TemplateName = "default-standalone-ve"

	dnsServersList := make([]string, 0)
	data.DnsServers.ElementsAs(ctx, &dnsServersList, false)
	deployConfig.Parameters.DnsServers = dnsServersList

	ntpServersList := make([]string, 0)
	data.NtpServers.ElementsAs(ctx, &ntpServersList, false)
	deployConfig.Parameters.NtpServers = ntpServersList

	l1Networks := []bigipnextsdk.CMReqL1Networks{}
	for _, val := range data.L1Networks {
		var l1Network bigipnextsdk.CMReqL1Networks
		l1Network.Name = val.Name.ValueString()
		l1Network.L1Link.Name = val.L1Link.Name.ValueString()
		l1Network.L1Link.LinkType = val.L1Link.LinkType.ValueString()

		l1Network.Vlans = []bigipnextsdk.CMReqVlans{}
		for _, dataVlan := range val.Vlans {
			var vlan bigipnextsdk.CMReqVlans
			vlan.Name = dataVlan.Name.ValueString()
			vlan.Tag = int(dataVlan.Tag.ValueFloat64())
			vlan.SelfIps = []bigipnextsdk.CMReqSelfIps{}
			for _, dataSelfIp := range dataVlan.SelfIps {
				var selfIp bigipnextsdk.CMReqSelfIps
				selfIp.DeviceName = dataSelfIp.DeviceName.ValueString()
				selfIp.Address = dataSelfIp.Address.ValueString()

				vlan.SelfIps = append(vlan.SelfIps, selfIp)
			}
			l1Network.Vlans = append(l1Network.Vlans, vlan)
		}
		l1Networks = append(l1Networks, l1Network)
	}
	deployConfig.Parameters.L1Networks = l1Networks
	deployConfig.Parameters.ManagementAddress = data.ManagementAddress.String()
	return &deployConfig
}

func (r *NextOnboardResource) instanceModeltoState(ctx context.Context, respData interface{}, data *NextOnboardResourceModel) {

	tflog.Info(ctx, fmt.Sprintf("[instanceModeltoState] respData : %v", respData))

	data.ManagementAddress = types.StringValue(respData.(map[string]interface{})["parameters"].(map[string]interface{})["management_address"].(string))
	data.DnsServers, _ = types.ListValueFrom(ctx, types.StringType, respData.(map[string]interface{})["parameters"].(map[string]interface{})["dns_servers"])
	data.NtpServers, _ = types.ListValueFrom(ctx, types.StringType, respData.(map[string]interface{})["parameters"].(map[string]interface{})["ntp_servers"])

	_, ok := respData.(map[string]interface{})["parameters"].(map[string]interface{})["l1Networks"]
	if ok {
		var stateL1Network []L1NetworkModel
		for _, l1Network := range respData.(map[string]interface{})["parameters"].(map[string]interface{})["l1Networks"].([]interface{}) {
			var i L1NetworkModel
			i.Name = types.StringValue(l1Network.(map[string]interface{})["name"].(string))
			i.L1Link.Name = types.StringValue(l1Network.(map[string]interface{})["l1Link"].(map[string]interface{})["name"].(string))
			i.L1Link.LinkType = types.StringValue(l1Network.(map[string]interface{})["l1Link"].(map[string]interface{})["linkType"].(string))
			_, ok := l1Network.(map[string]interface{})["vlans"]
			if ok {
				for _, vlan := range l1Network.(map[string]interface{})["vlans"].([]interface{}) {
					var j Vlan
					j.Name = types.StringValue(vlan.(map[string]interface{})["name"].(string))
					j.Tag = types.Float64Value(vlan.(map[string]interface{})["tag"].(float64))
					_, ok := vlan.(map[string]interface{})["selfIps"]
					if ok {
						for _, selfIp := range vlan.(map[string]interface{})["selfIps"].([]interface{}) {
							var k SelfIpModel
							k.Address = types.StringValue(selfIp.(map[string]interface{})["address"].(string))
							k.DeviceName = types.StringValue(selfIp.(map[string]interface{})["deviceName"].(string))

							j.SelfIps = append(j.SelfIps, k)
						}
					}

					i.Vlans = append(i.Vlans, j)
				}
			}
			stateL1Network = append(stateL1Network, i)
		}
		data.L1Networks = stateL1Network
	}
}
