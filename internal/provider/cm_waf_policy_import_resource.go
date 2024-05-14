package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	bigipnextsdk "gitswarm.f5net.com/terraform-providers/bigipnext"
)

var (
	_ resource.Resource                = &NextCMWAFPolicyImportResource{}
	_ resource.ResourceWithImportState = &NextCMWAFPolicyImportResource{}
	// mutex sync.Mutex
)

func NewNextCMWAFPolicyImportResource() resource.Resource {
	return &NextCMWAFPolicyImportResource{}
}

type NextCMWAFPolicyImportResource struct {
	client *bigipnextsdk.BigipNextCM
}

type NextCMWAFPolicyImportResourceModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	FilePath    types.String `tfsdk:"file_path"`
	FileMd5     types.String `tfsdk:"file_md5"`
	Override    types.String `tfsdk:"override"`
	Id          types.String `tfsdk:"id"`
}

func (r *NextCMWAFPolicyImportResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cm_waf_policy_import"
}

func (r *NextCMWAFPolicyImportResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource used to Import WAF Policy using policy json ( available local system disk )",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The unique user-given name of the policy. Policy names cannot contain spaces or special characters. Allowed characters are a-z, A-Z, 0-9, dot, dash (-), colon (:) and underscore (_).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[-a-zA-Z0-9._/:]+$`),
						"The name is not valid.",
					),
				},
			},
			"file_path": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Specifies WAF Policy Json file path ( available on local system disk path )",
			},
			"file_md5": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "MD5 value for WAF Policy Json file ( available on local system disk path )",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Specifies the description of the policy.",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtMost(255),
				},
			},
			"override": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Specifies Confirmation to override an existing policy with the same name.",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"true", "false"}...),
				},
				Default: stringdefault.StaticString("false"),
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

func (r *NextCMWAFPolicyImportResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client, resp.Diagnostics = toBigipNextCMProvider(req.ProviderData)
}

func (r *NextCMWAFPolicyImportResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var resCfg *NextCMWAFPolicyImportResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &resCfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	reqDraft := getCMWAFPolicyImportConfig(ctx, resCfg)
	tflog.Info(ctx, fmt.Sprintf("[CREATE] CM WAF Policy Import config : %+v\n", reqDraft))

	id, err := r.client.PolicyImport(reqDraft)
	if err != nil {
		resp.Diagnostics.AddError("Error", fmt.Sprintf("Failed to Import WAF Policy, got error: %s", err))
		return
	}
	resCfg.Id = types.StringValue(id.(map[string]interface{})["policy_id"].(string))
	resp.Diagnostics.Append(resp.State.Set(ctx, resCfg)...)

}

func (r *NextCMWAFPolicyImportResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var stateCfg *NextCMWAFPolicyImportResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateCfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id := stateCfg.Id.ValueString()
	tflog.Info(ctx, fmt.Sprintf("[READ] Reading WAF Policy : %s", id))
	wafData, err := r.client.GetWAFPolicyDetails(id)
	if err != nil {
		resp.Diagnostics.AddError("Error", fmt.Sprintf("Failed to Read WAF Policy, got error: %s", err))
		return
	}
	r.WafPolicyModeltoState(ctx, wafData, stateCfg)
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateCfg)...)
}

func (r *NextCMWAFPolicyImportResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var resCfg *NextCMWAFPolicyImportResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &resCfg)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[UPDATE] Updating WAF Policy : %s", resCfg.Name.ValueString()))

	// reqDraft := getCMWAFPolicyImportConfig(ctx, resCfg)

	tflog.Info(ctx, fmt.Sprintf("[UPDATE] id:%+v\n", resCfg.Id.ValueString()))

	reqDraft := getCMWAFPolicyImportConfig(ctx, resCfg)
	tflog.Info(ctx, fmt.Sprintf("[UPDATE] CM WAF Policy Import config : %+v\n", reqDraft))
	reqDraft.Override = "true"
	id, err := r.client.PolicyImport(reqDraft)
	if err != nil {
		resp.Diagnostics.AddError("Error", fmt.Sprintf("Failed to Import WAF Policy, got error: %s", err))
		return
	}
	resCfg.Id = types.StringValue(id.(map[string]interface{})["policy_id"].(string))

	resp.Diagnostics.Append(resp.State.Set(ctx, &resCfg)...)
}

func (r *NextCMWAFPolicyImportResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var stateCfg *NextCMWAFPolicyImportResourceModel
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &stateCfg)...)
	id := stateCfg.Id.ValueString()

	tflog.Info(ctx, fmt.Sprintf("[DELETE] Deleting WAF Policy : %s", id))

	err := r.client.DeleteWAFPolicy(id)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to Delete WAF Policy, got error: %s", err))
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &stateCfg)...)
	stateCfg.Id = types.StringValue("")
}

func (r *NextCMWAFPolicyImportResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getCMWAFPolicyImportConfig(ctx context.Context, data *NextCMWAFPolicyImportResourceModel) *bigipnextsdk.PolicyimportReqObj {
	var policyImportReqObj bigipnextsdk.PolicyimportReqObj
	policyImportReqObj.PolicyName = data.Name.ValueString()
	policyImportReqObj.FilePath = data.FilePath.ValueString()
	policyImportReqObj.Description = data.Description.ValueString()
	policyImportReqObj.Override = data.Override.ValueString()
	tflog.Info(ctx, fmt.Sprintf("[getCMWAFPolicyImportConfig] policyImportReqObj:%+v\n", policyImportReqObj))
	return &policyImportReqObj
}

func (r *NextCMWAFPolicyImportResource) WafPolicyModeltoState(ctx context.Context, respData interface{}, data *NextCMWAFPolicyImportResourceModel) {
	tflog.Info(ctx, fmt.Sprintf("WafPolicyModeltoState \t name: %+v", respData.(map[string]interface{})["name"]))
	data.Name = types.StringValue(respData.(map[string]interface{})["name"].(string))
	description, ok := respData.(map[string]interface{})["description"]
	if ok {
		data.Description = types.StringValue(description.(string))
	}
}
