package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	bigipnextsdk "gitswarm.f5net.com/terraform-providers/bigipnext"
)

var (
	_ resource.Resource                = &NextGlobalResiliencyResource{}
	_ resource.ResourceWithImportState = &NextGlobalResiliencyResource{}
	// mutex sync.Mutex
)

func NewNextGlobalResiliencyResource() resource.Resource {
	return &NextGlobalResiliencyResource{}
}

type NextGlobalResiliencyResource struct {
	client *bigipnextsdk.BigipNextCM
}

type NextGlobalResiliencyResourceModel struct {
	Name            types.String `tfsdk:"name"`
	DNSListenerName types.String `tfsdk:"dns_listener_name"`
	Protocols       types.List   `tfsdk:"protocols"`
	Instances       []Instance   `tfsdk:"instances"`
	DNSListenerPort types.Int64  `tfsdk:"dns_listener_port"`
	Id              types.String `tfsdk:"id"`
}

type Instance struct {
	Hostname           types.String `tfsdk:"hostname"`
	Address            types.String `tfsdk:"address"`
	DNSListenerAddress types.String `tfsdk:"dns_listener_address"`
	GroupSyncAddress   types.String `tfsdk:"group_sync_address"`
}

func (r *NextGlobalResiliencyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cm_global_resiliency"
}

func (r *NextGlobalResiliencyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource used to manage(CRUD) Global Resiliency resources onto BIG-IP Next CM.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Global Resiliency Group Name. The group name must start with lowercase letters (a-z) and consist only of lowercase letters (a-z) and digits (0-9).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-z][a-z0-9]*$`),
						"The group name must start with lowercase letters (a-z) and consist only of lowercase letters (a-z) and digits (0-9).",
					),
				},
			},
			"dns_listener_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "DNS Listener Name. The DNS listener name must start with lowercase letters (a-z) and consist only of lowercase letters (a-z) and digits (0-9).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-z][a-z0-9]*$`),
						"The dns listener name must start with lowercase letters (a-z) and consist only of lowercase letters (a-z) and digits (0-9).",
					),
				},
			},
			"protocols": schema.ListAttribute{
				Required:            true,
				MarkdownDescription: "Protocols to be added to the Global Resiliency Group. Protocols cannot be updated once created.",
				ElementType:         types.StringType,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(stringvalidator.OneOf([]string{"udp", "tcp"}...)),
				},
			},
			"dns_listener_port": schema.Int64Attribute{
				MarkdownDescription: "DNS Listener Port. Port number must be greater than or equal to 1. Port number must not exceed 65535. Port cannot be updated once created",
				Required:            true,
				// Validators: []validator.Int64{
				// 	int64validator.Between(1, 65535),
				// },
			},
			"instances": schema.ListNestedAttribute{
				Required:            true,
				MarkdownDescription: "List of Instances",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"hostname": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Hostname of the Instance to be added",
						},
						"address": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Address of the Bip-IP Next. A valid IP Address is required",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(`^([01]?\d{1,2}|2[0-4]\d|25[0-5])\.([01]?\d{1,2}|2[0-4]\d|25[0-5])\.([01]?\d{1,2}|2[0-4]\d|25[0-5])\.([01]?\d{1,2}|2[0-4]\d|25[0-5])$`),
									"given address is not a valid IPV4 address",
								),
							},
						},
						"dns_listener_address": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "DNS Listener Address. A valid IP Address is required",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(`^([01]?\d{1,2}|2[0-4]\d|25[0-5])\.([01]?\d{1,2}|2[0-4]\d|25[0-5])\.([01]?\d{1,2}|2[0-4]\d|25[0-5])\.([01]?\d{1,2}|2[0-4]\d|25[0-5])$`),
									"given dns_listener_address is not a valid IPV4 address",
								),
							},
						},
						"group_sync_address": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "GR Group Sunc IP. A valid IP Address with mask is required",
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(`^([01]?\d{1,2}|2[0-4]\d|25[0-5])\.([01]?\d{1,2}|2[0-4]\d|25[0-5])\.([01]?\d{1,2}|2[0-4]\d|25[0-5])\.([01]?\d{1,2}|2[0-4]\d|25[0-5])/([12]?\d|3[0-2])$`),
									"given group_sync_address must be a valid IPV4 address in CIDR format",
								),
							},
						},
					},
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

func (r *NextGlobalResiliencyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client, resp.Diagnostics = toBigipNextCMProvider(req.ProviderData)
}

func (r *NextGlobalResiliencyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var resCfg *NextGlobalResiliencyResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &resCfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[CREATE] NextGlobalResiliencyResource:%+v\n", resCfg.Name.ValueString()))

	reqDraft := getGlobalResiliencyRequestDraft(ctx, resCfg)

	tflog.Info(ctx, "[CREATE]  Global Resiliency Group")
	tflog.Info(ctx, fmt.Sprintf("[CREATE] :%+v\n", reqDraft))

	id, err := r.client.PostGlobalResiliencyGroup("POST", reqDraft)
	if err != nil {
		resp.Diagnostics.AddError("Global Resiliency Error", fmt.Sprintf("Failed to Create Global Resiliency Group, got error: %s", err))
		return
	}

	tflog.Info(ctx, "[CREATE] Global Resiliency Group created successfully")

	resCfg.Id = types.StringValue(id)
	resp.Diagnostics.Append(resp.State.Set(ctx, resCfg)...)

}

func (r *NextGlobalResiliencyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var stateCfg *NextGlobalResiliencyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateCfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id := stateCfg.Id.ValueString()
	tflog.Info(ctx, fmt.Sprintf("Reading Global Resiliency Group : %s", id))
	grData, err := r.client.GetGlobalResiliencyGroupDetails(id)
	if err != nil {
		resp.Diagnostics.AddError("Error", fmt.Sprintf("Failed to Read  Global Resiliency Group, got error: %s", err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Global Resiliency Group : %+v", grData))
	r.GlobalResilienceModeltoState(ctx, grData, stateCfg)
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateCfg)...)
}

func (r *NextGlobalResiliencyResource) GlobalResilienceModeltoState(ctx context.Context, respData interface{}, data *NextGlobalResiliencyResourceModel) {
	tflog.Info(ctx, fmt.Sprintf("GlobalResilienceModeltoState \t name: %+v", respData.(map[string]interface{})["name"].(string)))
	data.Name = types.StringValue(respData.(map[string]interface{})["name"].(string))
	data.DNSListenerName = types.StringValue(respData.(map[string]interface{})["dns_listener_name"].(string))
	data.DNSListenerPort = types.Int64Value(int64(respData.(map[string]interface{})["dns_listener_port"].(float64)))
	data.Id = types.StringValue(respData.(map[string]interface{})["id"].(string))

	var instanceList []Instance
	for _, instance := range respData.(map[string]interface{})["instances"].([]interface{}) {
		var i Instance
		i.Hostname = types.StringValue(instance.(map[string]interface{})["hostname"].(string))
		i.Address = types.StringValue(instance.(map[string]interface{})["address"].(string))
		i.DNSListenerAddress = types.StringValue(instance.(map[string]interface{})["dns_listener_address"].(string))
		i.GroupSyncAddress = types.StringValue(instance.(map[string]interface{})["group_sync_address"].(string))

		instanceList = append(instanceList, i)
	}
	data.Instances = instanceList
}

func (r *NextGlobalResiliencyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var resCfg *NextGlobalResiliencyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &resCfg)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[UPDATE] Updating Global Resiliency Group: %s", resCfg.Name.ValueString()))

	reqDraft := getGlobalResiliencyRequestDraft(ctx, resCfg)
	reqDraft.Id = resCfg.Id.ValueString()

	tflog.Info(ctx, "[UPDATE] Updating Global Resiliency Group")
	tflog.Info(ctx, fmt.Sprintf("[UPDATE] :%+v\n", reqDraft))

	id, err := r.client.PostGlobalResiliencyGroup("PUT", reqDraft)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to Update Global Resiliency Group, got error: %s", err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[Update] id:%+v\n", id))
	resCfg.Id = types.StringValue(id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &resCfg)...)
}

func (r *NextGlobalResiliencyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var stateCfg *NextGlobalResiliencyResourceModel
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &stateCfg)...)
	id := stateCfg.Id.ValueString()

	tflog.Info(ctx, fmt.Sprintf("Deleting Global Resiliency Group : %s", id))

	err := r.client.DeleteGlobalResiliencyGroup(id)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to Delete Global Resiliency Group, got error: %s", err))
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &stateCfg)...)
	stateCfg.Id = types.StringValue("")
}

func (r *NextGlobalResiliencyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getGlobalResiliencyRequestDraft(ctx context.Context, data *NextGlobalResiliencyResourceModel) *bigipnextsdk.GlobalResiliencyRequestDraft {
	var globalresilienceReqDraft bigipnextsdk.GlobalResiliencyRequestDraft
	globalresilienceReqDraft.Name = data.Name.ValueString()
	protocolList := make([]string, 0, 1)
	data.Protocols.ElementsAs(ctx, &protocolList, false)
	globalresilienceReqDraft.Protocols = protocolList
	globalresilienceReqDraft.DNSListenerName = data.DNSListenerName.ValueString()
	globalresilienceReqDraft.DNSListenerPort = int(data.DNSListenerPort.ValueInt64())
	var instanceList []bigipnextsdk.Instance
	for _, val := range data.Instances {
		var instance bigipnextsdk.Instance
		instance.Address = val.Address.ValueString()
		instance.DNSListenerAddress = val.DNSListenerAddress.ValueString()
		instance.Hostname = val.Hostname.ValueString()
		instance.GroupSyncAddress = val.GroupSyncAddress.ValueString()
		instanceList = append(instanceList, instance)
	}
	globalresilienceReqDraft.Instances = instanceList
	tflog.Info(ctx, fmt.Sprintf("globalresilienceReqDraft:%+v\n", globalresilienceReqDraft))
	return &globalresilienceReqDraft
}
