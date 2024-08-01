//go:build !test

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	bigipnextsdk "gitswarm.f5net.com/terraform-providers/bigipnext"
)

var (
	_ resource.Resource                = &NextCMFastTemplateResource{}
	_ resource.ResourceWithImportState = &NextCMFastTemplateResource{}
)

func NewNextCMFastTemplateResource() resource.Resource {
	return &NextCMFastTemplateResource{}
}

type NextCMFastTemplateResource struct {
	client *bigipnextsdk.BigipNextCM
}
type NextCMFastTemplateResourceModel struct {
	TemplateName           types.String `tfsdk:"template_name"`
	ApplicationDescription types.String `tfsdk:"application_description"`
	ApplicationName        types.String `tfsdk:"application_name"`
	Pools                  types.List   `tfsdk:"pools"`
	Virtuals               types.List   `tfsdk:"virtuals"`
	SetName                types.String `tfsdk:"set_name"`
	TenantName             types.String `tfsdk:"tenant_name"`
	AllowOverwrite         types.Bool   `tfsdk:"allow_overwrite"`
	Id                     types.String `tfsdk:"id"`
}

type NextCMFastTemplatePoolModel struct {
	LoadBalancingMode types.String `tfsdk:"load_balancing_mode"`
	MonitorType       types.List   `tfsdk:"monitor_type"`
	PoolName          types.String `tfsdk:"pool_name"`
	ServicePort       types.Int64  `tfsdk:"service_port"`
}

// pool_name
type NextCMFastTemplateVirtualModel struct {
	PoolName    types.String `tfsdk:"pool_name"`
	VirtualName types.String `tfsdk:"virtual_name"`
	VirtualPort types.Int64  `tfsdk:"virtual_port"`
}

func (r *NextCMFastTemplateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cm_fast_template"
}

func (r *NextCMFastTemplateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource used to Manages FAST templates on Central Manager",
		Attributes: map[string]schema.Attribute{
			"template_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the FAST template to be created",
			},
			"template_set": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the FAST template set to use\nIf a set with given name does not exist a new set will be created automatically",
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

func (r *NextCMFastTemplateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client, resp.Diagnostics = toBigipNextCMProvider(req.ProviderData)
}

func (r *NextCMFastTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var resCfg *NextCMFastTemplateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &resCfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	//git tflog.Info(ctx, fmt.Sprintf("[CREATE] NextCMFastTemplateResource:%+v\n", resCfg.Name.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, resCfg)...)
}

func (r *NextCMFastTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var stateCfg *NextCMFastTemplateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateCfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateCfg)...)
}

func (r *NextCMFastTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var resCfg *NextCMFastTemplateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &resCfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &resCfg)...)
}

func (r *NextCMFastTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var stateCfg *NextCMFastTemplateResourceModel
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &stateCfg)...)
	stateCfg.Id = types.StringValue("")
}

func (r *NextCMFastTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
