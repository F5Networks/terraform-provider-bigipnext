package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	bigipnextsdk "gitswarm.f5net.com/terraform-providers/bigipnext"
	"regexp"
)

var (
	_ resource.Resource                = &NextCMWAFReportResource{}
	_ resource.ResourceWithImportState = &NextCMWAFReportResource{}
	// mutex sync.Mutex
)

func NewNextCMWAFReportResource() resource.Resource {
	return &NextCMWAFReportResource{}
}

type NextCMWAFReportResource struct {
	client *bigipnextsdk.BigipNextCM
}

type NextCMWAFReportResourceModel struct {
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	TimeFrameInDays types.Int64  `tfsdk:"time_frame_in_days"`
	TopLevel        types.Int64  `tfsdk:"top_level"`
	RequestType     types.String `tfsdk:"request_type"`
	CreatedBy       types.String `tfsdk:"created_by"`
	Scope           Policy       `tfsdk:"scope"`
	Categories      []Category   `tfsdk:"categories"`
	UserDefined     types.Bool   `tfsdk:"user_defined"`
	Id              types.String `tfsdk:"id"`
}

type Category struct {
	Name types.String `tfsdk:"name"`
}

type Policy struct {
	Entity types.String `tfsdk:"entity"`
	All    types.Bool   `tfsdk:"all"`
	Names  types.List   `tfsdk:"names"`
}

func (r *NextCMWAFReportResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cm_waf_report"
}

func (r *NextCMWAFReportResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource used to manage(CRUD) WAF Security Report resources onto BIG-IP Next CM.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the security report.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[-a-zA-Z0-9._/:\s]+$`),
						"The name is not valid.",
					),
				},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Specifies the description of the security report. Description should be less than 255 character",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtMost(255),
				},
			},
			"time_frame_in_days": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Specifies the report time period.",
			},
			"top_level": schema.Int64Attribute{
				Required:            true,
				MarkdownDescription: "Specifies the number of the top level items of the report.",
			},
			"request_type": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Specifies the number of the top level items of the report.",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"illegal", "alerted", "blocked"}...),
				},
			},
			"user_defined": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Specifies whether the report is user defined. This is a computed value.",
			},
			"created_by": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The creator of the security report. This is a computed value.",
			},
			"scope": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "Specify the Policies/Applications for the WAF security report",
				Attributes: map[string]schema.Attribute{
					"entity": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Entity value can be policies or applications",
						Validators: []validator.String{
							stringvalidator.OneOf([]string{"policies", "applications"}...),
						},
					},
					"all": schema.BoolAttribute{
						Required:            true,
						MarkdownDescription: "Specifies whether All policies/applications are to be taken or selected ones. 'names' must not be empty if 'all' is set to false.",
					},
					"names": schema.ListAttribute{
						Optional:            true,
						ElementType:         types.StringType,
						MarkdownDescription: "Specifies the names of the scoped entities.",
					},
				},
			},
			"categories": schema.ListNestedAttribute{
				Optional:            true,
				MarkdownDescription: "List of Categories",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Optional:            true,
							MarkdownDescription: "Specifies the name of the Categories.",
							Validators: []validator.String{
								stringvalidator.OneOf([]string{"Source IPs", "Geolocations", "URLs", "Violations", "Signatures", "Attack Types", "Threat Campaigns", "Malicious IPs (IPI)", ""}...),
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

func (r *NextCMWAFReportResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client, resp.Diagnostics = toBigipNextCMProvider(req.ProviderData)
}

func (r *NextCMWAFReportResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var resCfg *NextCMWAFReportResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &resCfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[CREATE] NextCMWAFReportResource:%+v\n", resCfg.Name.ValueString()))

	reqDraft := getCMWAFReportRequestDraft(ctx, resCfg)

	tflog.Info(ctx, "[CREATE]  WAF Security Report")
	tflog.Info(ctx, fmt.Sprintf("[CREATE] :%+v\n", reqDraft))

	id, created_by, user_defined, err := r.client.PostWAFReport("POST", reqDraft)
	if err != nil {
		resp.Diagnostics.AddError("WAF Security Report Error", fmt.Sprintf("Failed to Create WAF Security Report, got error: %s", err))
		return
	}

	resCfg.Id = types.StringValue(id)
	resCfg.UserDefined = types.BoolValue(user_defined)
	resCfg.CreatedBy = types.StringValue(created_by)

	resp.Diagnostics.Append(resp.State.Set(ctx, resCfg)...)

}

func (r *NextCMWAFReportResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var stateCfg *NextCMWAFReportResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateCfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id := stateCfg.Id.ValueString()
	tflog.Info(ctx, fmt.Sprintf("[READ] Reading WAF Security Report : %s", id))
	wafData, err := r.client.GetWAFReportDetails(id)
	if err != nil {
		resp.Diagnostics.AddError("Error", fmt.Sprintf("Failed to Read WAF Security Report, got error: %s", err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[READ] WAF Security Report : %+v", wafData))
	r.WafReportModeltoState(ctx, wafData, stateCfg)
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateCfg)...)
}

func (r *NextCMWAFReportResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var resCfg *NextCMWAFReportResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &resCfg)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[UPDATE] Updating WAF Security Report: %s", resCfg.Name.ValueString()))

	reqDraft := getCMWAFReportRequestDraft(ctx, resCfg)
	reqDraft.Id = resCfg.Id.ValueString()

	tflog.Info(ctx, "[UPDATE] Updating WAF Security Report")
	tflog.Info(ctx, fmt.Sprintf("[UPDATE] :%+v\n", reqDraft))

	id, created_by, user_defined, err := r.client.PostWAFReport("PUT", reqDraft)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to Update WAF Security Report, got error: %s", err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[UPDATE] id:%+v\n", id))

	resCfg.Id = types.StringValue(id)
	resCfg.UserDefined = types.BoolValue(user_defined)
	resCfg.CreatedBy = types.StringValue(created_by)

	resp.Diagnostics.Append(resp.State.Set(ctx, &resCfg)...)
}

func (r *NextCMWAFReportResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var stateCfg *NextCMWAFReportResourceModel
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &stateCfg)...)
	id := stateCfg.Id.ValueString()

	tflog.Info(ctx, fmt.Sprintf("[DELETE] Deleting WAF Security Report : %s", id))

	err := r.client.DeleteWAFReport(id)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to Delete WAF Security Report, got error: %s", err))
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &stateCfg)...)
	stateCfg.Id = types.StringValue("")
}

func (r *NextCMWAFReportResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getCMWAFReportRequestDraft(ctx context.Context, data *NextCMWAFReportResourceModel) *bigipnextsdk.CMWAFReportRequestDraft {

	tflog.Info(ctx, "[getCMWAFReportRequestDraft] Fetching WAF Security Report Request Draft")

	var wafreportReqDraft bigipnextsdk.CMWAFReportRequestDraft
	wafreportReqDraft.Name = data.Name.ValueString()
	wafreportReqDraft.Description = data.Description.ValueString()
	wafreportReqDraft.TimeFrameInDays = int(data.TimeFrameInDays.ValueInt64())
	wafreportReqDraft.TopLevel = int(data.TopLevel.ValueInt64())
	wafreportReqDraft.RequestType = data.RequestType.ValueString()

	var categoriesList []bigipnextsdk.Category
	if len(data.Categories) != 0 {
		for _, val := range data.Categories {
			var category bigipnextsdk.Category
			category.Name = val.Name.ValueString()
			categoriesList = append(categoriesList, category)
		}
		wafreportReqDraft.Categories = categoriesList
	}

	var namesList []string
	for _, val := range data.Scope.Names.Elements() {
		var ss string
		_ = json.Unmarshal([]byte(val.String()), &ss)
		namesList = append(namesList, ss)
	}
	wafreportReqDraft.Scope.All = data.Scope.All.ValueBool()
	wafreportReqDraft.Scope.Entity = data.Scope.Entity.ValueString()
	wafreportReqDraft.Scope.Names = namesList

	tflog.Info(ctx, fmt.Sprintf("[getCMWAFReportRequestDraft] wafreportReqDraft:%+v\n", wafreportReqDraft))
	return &wafreportReqDraft
}

func (r *NextCMWAFReportResource) WafReportModeltoState(ctx context.Context, respData interface{}, data *NextCMWAFReportResourceModel) {

	tflog.Info(ctx, fmt.Sprintf("WafReportModeltoState \t name: %+v", respData.(map[string]interface{})["name"]))

	data.Name = types.StringValue(respData.(map[string]interface{})["name"].(string))
	description, ok := respData.(map[string]interface{})["description"]
	if ok {
		data.Description = types.StringValue(description.(string))
	}
	data.TimeFrameInDays = types.Int64Value(int64(respData.(map[string]interface{})["time_frame_in_days"].(float64)))
	data.TopLevel = types.Int64Value(int64(respData.(map[string]interface{})["top_level"].(float64)))
	data.RequestType = types.StringValue(respData.(map[string]interface{})["request_type"].(string))
	data.Id = types.StringValue(respData.(map[string]interface{})["id"].(string))
	data.UserDefined = types.BoolValue(respData.(map[string]interface{})["user_defined"].(bool))
	data.CreatedBy = types.StringValue(respData.(map[string]interface{})["created_by"].(string))

	_, ok = respData.(map[string]interface{})["categories"]
	if ok {
		var categoriesList []Category
		for _, category := range respData.(map[string]interface{})["categories"].([]interface{}) {
			var i Category
			i.Name = types.StringValue(category.(map[string]interface{})["name"].(string))
			tflog.Info(ctx, fmt.Sprintf("WafReportModeltoState:%+v\n", i.Name.ValueString()))
			categoriesList = append(categoriesList, i)
		}
		data.Categories = categoriesList
	}

	data.Scope.Entity = types.StringValue(respData.(map[string]interface{})["scope"].(map[string]interface{})["entity"].(string))
	data.Scope.All = types.BoolValue(respData.(map[string]interface{})["scope"].(map[string]interface{})["all"].(bool))
	data.Scope.Names, _ = types.ListValueFrom(ctx, types.StringType, respData.(map[string]interface{})["scope"].(map[string]interface{})["names"])
}
