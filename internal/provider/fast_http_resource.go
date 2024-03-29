package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	bigipnextsdk "gitswarm.f5net.com/terraform-providers/bigipnext"
)

var (
	_ resource.Resource                = &NextCMFastHttpResource{}
	_ resource.ResourceWithImportState = &NextCMFastHttpResource{}
)

func NewNextCMFastHttpResource() resource.Resource {
	return &NextCMFastHttpResource{}
}

type NextCMFastHttpResource struct {
	client *bigipnextsdk.BigipNextCM
}
type NextCMFastHttpResourceModel struct {
	Name                   types.String `tfsdk:"name"`
	ApplicationDescription types.String `tfsdk:"application_description"`
	ApplicationName        types.String `tfsdk:"application_name"`
	Pools                  types.List   `tfsdk:"pools"`
	Virtuals               types.List   `tfsdk:"virtuals"`
	SetName                types.String `tfsdk:"set_name"`
	TemplateName           types.String `tfsdk:"template_name"`
	TenantName             types.String `tfsdk:"tenant_name"`
	AllowOverwrite         types.Bool   `tfsdk:"allow_overwrite"`
	Id                     types.String `tfsdk:"id"`
}

type NextCMFastHttpPoolModel struct {
	LoadBalancingMode types.String `tfsdk:"load_balancing_mode"`
	MonitorType       types.List   `tfsdk:"monitor_type"`
	PoolName          types.String `tfsdk:"pool_name"`
	ServicePort       types.Int64  `tfsdk:"service_port"`
}

// pool_name
type NextCMFastHttpVirtualModel struct {
	PoolName    types.String `tfsdk:"pool_name"`
	VirtualName types.String `tfsdk:"virtual_name"`
	VirtualPort types.Int64  `tfsdk:"virtual_port"`
}

func (r *NextCMFastHttpResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cm_fast_http"
}

func (r *NextCMFastHttpResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource used to manage(CRUD) AS3 declarations using BIG-IP Next CM onto target BIG-IP Next",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the Application",
			},
			"application_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the Application",
			},
			"tenant_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the Tenant",
			},
			"application_description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Description of the Application",
			},
			"pools": schema.ListNestedAttribute{
				Optional:            true,
				MarkdownDescription: "List of Pools",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"load_balancing_mode": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Load Balancing Mode",
						},
						"monitor_type": schema.ListAttribute{
							Required:            true,
							ElementType:         types.StringType,
							MarkdownDescription: "Monitor Type",
						},
						"pool_name": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Name of the Pool",
						},
						"service_port": schema.Int64Attribute{
							Required:            true,
							MarkdownDescription: "Service Port",
						},
					},
				},
			},
			"virtuals": schema.ListNestedAttribute{
				Optional:            true,
				MarkdownDescription: "List of Virtuals",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"pool_name": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Name of the Pool",
						},
						"virtual_name": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Name of the Virtual",
						},
						"virtual_port": schema.Int64Attribute{
							Required:            true,
							MarkdownDescription: "Virtual Port",
						},
					},
				},
			},
			"set_name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Name of the AS3 Set",
				Default:             stringdefault.StaticString("Examples"),
			},
			"template_name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Name of the AS3 Template",
				Default:             stringdefault.StaticString("http"),
			},
			"allow_overwrite": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Allow Overwrite",
				Default:             booldefault.StaticBool(false),
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

func (r *NextCMFastHttpResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client, resp.Diagnostics = toBigipNextCMProvider(req.ProviderData)
}

func (r *NextCMFastHttpResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var resCfg *NextCMFastHttpResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &resCfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[CREATE] NextCMFastHttpResource:%+v\n", resCfg.Name.ValueString()))
	reqDraft := getFastRequestDraft(ctx, resCfg)

	tflog.Info(ctx, fmt.Sprintf("[CREATE] Https:%+v\n", reqDraft))

	draftID, err := r.client.PostFastApplicationDraft(reqDraft)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to Create FAST Draft config, got error: %s", err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[CREATE] draftID:%+v\n", draftID))
	resp.Diagnostics.Append(resp.State.Set(ctx, resCfg)...)
}

func (r *NextCMFastHttpResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var stateCfg *NextCMFastHttpResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateCfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateCfg)...)
}

func (r *NextCMFastHttpResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var resCfg *NextCMFastHttpResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &resCfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &resCfg)...)
}

func (r *NextCMFastHttpResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var stateCfg *NextCMFastHttpResourceModel
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &stateCfg)...)
	stateCfg.Id = types.StringValue("")
}

func (r *NextCMFastHttpResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getFastRequestDraft(ctx context.Context, data *NextCMFastHttpResourceModel) *bigipnextsdk.FastRequestDraft {
	var fastReqDraft bigipnextsdk.FastRequestDraft
	fastReqDraft.Name = data.Name.ValueString()
	fastReqDraft.Parameters.ApplicationDescription = data.ApplicationDescription.ValueString()
	fastReqDraft.Parameters.ApplicationName = data.ApplicationName.ValueString()
	fastReqDraft.SetName = data.SetName.ValueString()
	fastReqDraft.TemplateName = data.TemplateName.ValueString()
	fastReqDraft.TenantName = data.TenantName.ValueString()
	fastReqDraft.AllowOverwrite = data.AllowOverwrite.ValueBool()
	elements := make([]types.Object, 0, len(data.Pools.Elements()))
	data.Pools.ElementsAs(ctx, &elements, false)
	for _, element := range elements {
		var fastPool bigipnextsdk.FastPool
		var objectModel NextCMFastHttpPoolModel
		element.As(ctx, &objectModel, basetypes.ObjectAsOptions{})
		fastPool.LoadBalancingMode = objectModel.LoadBalancingMode.ValueString()
		objectModel.MonitorType.ElementsAs(ctx, &fastPool.MonitorType, false)
		fastPool.PoolName = objectModel.PoolName.ValueString()
		fastPool.ServicePort = int(objectModel.ServicePort.ValueInt64())
		fastReqDraft.Parameters.Pools = append(fastReqDraft.Parameters.Pools, fastPool)
	}
	vsLists := make([]types.Object, 0, len(data.Virtuals.Elements()))
	data.Virtuals.ElementsAs(ctx, &vsLists, false)
	for _, element := range vsLists {
		var fastVS bigipnextsdk.VirtualServer
		var objectModel NextCMFastHttpVirtualModel
		element.As(ctx, &objectModel, basetypes.ObjectAsOptions{})
		fastVS.VirtualName = objectModel.VirtualName.ValueString()
		fastVS.VirtualPort = int(objectModel.VirtualPort.ValueInt64())
		fastReqDraft.Parameters.Virtuals = append(fastReqDraft.Parameters.Virtuals, fastVS)
	}
	tflog.Info(ctx, fmt.Sprintf("fastReqDraft:%+v\n", fastReqDraft))
	return &fastReqDraft
}
