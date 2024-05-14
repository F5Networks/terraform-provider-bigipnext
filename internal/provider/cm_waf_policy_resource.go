package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	bigipnextsdk "gitswarm.f5net.com/terraform-providers/bigipnext"
)

var (
	_ resource.Resource                = &NextCMWAFPolicyResource{}
	_ resource.ResourceWithImportState = &NextCMWAFPolicyResource{}
	// mutex sync.Mutex
)

func NewNextCMWAFPolicyResource() resource.Resource {
	return &NextCMWAFPolicyResource{}
}

type NextCMWAFPolicyResource struct {
	client *bigipnextsdk.BigipNextCM
}

type NextCMWAFPolicyResourceModel struct {
	Name                types.String `tfsdk:"name"`
	Description         types.String `tfsdk:"description"`
	Tags                types.List   `tfsdk:"tags"`
	EnforecementMode    types.String `tfsdk:"enforcement_mode"`
	ApplicationLanguage types.String `tfsdk:"application_language"`
	TemplateName        types.String `tfsdk:"template_name"`
	BotDefense          types.Object `tfsdk:"bot_defense"`
	IpIntelligence      types.Object `tfsdk:"ip_intelligence"`
	DosProtection       types.Object `tfsdk:"dos_protection"`
	BlockingSettings    types.Object `tfsdk:"blocking_settings"`
	Id                  types.String `tfsdk:"id"`
}

type BotDefenseModel struct {
	Enabled types.Bool `tfsdk:"enabled"`
}

type IpIntelligenceModel struct {
	Enabled types.Bool `tfsdk:"enabled"`
}

type DosProtectionModel struct {
	Enabled types.Bool `tfsdk:"enabled"`
}

type BlockingSettingsModel struct {
	Enabled types.Bool `tfsdk:"enabled"`
}

type Violation struct {
	Alarm       types.Bool   `tfsdk:"alarm"`
	Block       types.Bool   `tfsdk:"block"`
	Description types.String `tfsdk:"description"`
	Name        types.String `tfsdk:"name"`
}

func (r *NextCMWAFPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cm_waf_policy"
}

func (r *NextCMWAFPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource used to manage(CRUD) WAF Policy resources onto BIG-IP Next CM.",
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
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Specifies the description of the policy.",
				Validators: []validator.String{
					stringvalidator.UTF8LengthAtMost(255),
				},
			},
			"tags": schema.ListAttribute{
				Optional:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Specifies the Tags for marking policies.",
			},
			"enforcement_mode": schema.StringAttribute{
				Required: true,
				MarkdownDescription: "Specifies How BIG-IP MA processes a request that triggers a security policy violation. \n*Blocking: When the enforcement mode is set to blocking, any triggered violation is blocked (configured for blocking)." +
					"\n*Transparent: When the enforcement mode is set to transparent, traffic is not blocked even if a violation is triggered.",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"blocking", "transparent"}...),
				},
			},
			"application_language": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The character encoding for the web application. Character encoding determines how the policy processes the character sets. ",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"iso-8859-6", "iso-8859-7", "utf-8", "windows-1256", "iso-8859-13", "iso-8859-4", "windows-1257", "windows-1250", "windows-1251", "windows-1253", "windows-1255", "windows-874", "windows-1252", "iso-8859-2", "big5", "gb18030", "gb2312", "gbk", "iso-8859-5", "koi8-r", "iso-8859-8", "euc-jp", "shift_jis", "euc-kr", "iso-8859-10", "iso-8859-16", "iso-8859-3", "iso-8859-9", "iso-8859-1", "iso-8859-15"}...),
				},
			},
			"template_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the template used to create the WAF policy. Template cannot be updated",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"Fundamental-Template", "Rating-Based-Template", "Comprehensive-Template", "RAPID-Template"}...),
				},
			},
			// "policy_import_json": schema.StringAttribute{
			// 	Optional:            true,
			// 	MarkdownDescription: "The name of the template used to create the WAF policy. Template cannot be updated",
			// },
			"bot_defense": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Specifies whether the bot defense for Policy is to be enabled or not. The default value of bot_defense is True.",
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(true),
					},
				},
			},
			"ip_intelligence": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Specifies whether the bot ip_intelligence for Policy is to be enabled or not. The default value of ip_intelligence is True.",
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(true),
					},
				},
			},
			"dos_protection": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Specifies whether the dos protection for Policy is to be enabled or not. The default value of dos_protection is False.",
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
				},
			},
			"blocking_settings": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Specifies whether the blocking setting is to be enabled or not. The default value of blocking_settings is True.",
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(true),
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

func (r *NextCMWAFPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client, resp.Diagnostics = toBigipNextCMProvider(req.ProviderData)
}

func (r *NextCMWAFPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var resCfg *NextCMWAFPolicyResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &resCfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[CREATE] NextCMWAFPolicyResource:%+v\n", resCfg))
	reqDraft := getCMWAFPolicyRequestDraft(ctx, resCfg)
	tflog.Info(ctx, fmt.Sprintf("[CREATE] WAF Policy config :%+v\n", reqDraft))

	id, err := r.client.PostWAFPolicy("POST", reqDraft)
	if err != nil {
		resp.Diagnostics.AddError("WAF Policy Error", fmt.Sprintf("Failed to Create WAF Policy, got error: %s", err))
		return
	}
	resCfg.Id = types.StringValue(id)

	// if resCfg.BotDefense.IsNull() || resCfg.BotDefense.IsUnknown() {
	// 	var botdefenseModel BotDefenseModel
	// 	botdefenseModel.Enabled = types.BoolValue(true)
	// 	diag := resCfg.BotDefense.As(ctx, &botdefenseModel, basetypes.ObjectAsOptions{})
	// 	if diag.HasError() {
	// 		return
	// 	}
	// }

	// if resCfg.IpIntelligence.IsNull() || resCfg.IpIntelligence.IsUnknown() {
	// 	var ipintelligenceModel IpIntelligenceModel
	// 	ipintelligenceModel.Enabled = types.BoolValue(true)
	// 	diag := resCfg.IpIntelligence.As(ctx, &ipintelligenceModel, basetypes.ObjectAsOptions{})
	// 	if diag.HasError() {
	// 		return
	// 	}
	// }

	// if resCfg.DosProtection.IsNull() || resCfg.DosProtection.IsUnknown() {
	// 	var dosprotectionModel DosProtectionModel
	// 	dosprotectionModel.Enabled = types.BoolValue(false)
	// 	diag := resCfg.DosProtection.As(ctx, &dosprotectionModel, basetypes.ObjectAsOptions{})
	// 	if diag.HasError() {
	// 		return
	// 	}
	// }

	// if resCfg.BlockingSettings.IsNull() || resCfg.BotDefense.IsUnknown() {
	// 	var blockingsettingdModel BlockingSettingsModel
	// 	blockingsettingdModel.Enabled = types.BoolValue(true)
	// 	diag := resCfg.BlockingSettings.As(ctx, &blockingsettingdModel, basetypes.ObjectAsOptions{})
	// 	if diag.HasError() {
	// 		return
	// 	}
	// }

	resp.Diagnostics.Append(resp.State.Set(ctx, resCfg)...)

}

func (r *NextCMWAFPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var stateCfg *NextCMWAFPolicyResourceModel
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
	tflog.Info(ctx, fmt.Sprintf("[READ] WAF Policy : %+v", wafData))
	r.WafPolicyModeltoState(ctx, wafData, stateCfg)
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateCfg)...)
}

func (r *NextCMWAFPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var resCfg *NextCMWAFPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &resCfg)...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[UPDATE] Updating WAF Policy : %s", resCfg.Name.ValueString()))

	reqDraft := getCMWAFPolicyRequestDraft(ctx, resCfg)
	reqDraft.Id = resCfg.Id.ValueString()

	tflog.Info(ctx, "[UPDATE] Updating WAF Policy")
	tflog.Info(ctx, fmt.Sprintf("[UPDATE] :%+v\n", reqDraft))

	id, err := r.client.PostWAFPolicy("PUT", reqDraft)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to Update WAF Policy, got error: %s", err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[UPDATE] id:%+v\n", id))

	resCfg.Id = types.StringValue(id)

	resp.Diagnostics.Append(resp.State.Set(ctx, &resCfg)...)
}

func (r *NextCMWAFPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var stateCfg *NextCMWAFPolicyResourceModel
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

func (r *NextCMWAFPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getCMWAFPolicyRequestDraft(ctx context.Context, data *NextCMWAFPolicyResourceModel) *bigipnextsdk.CMWAFPolicyRequestDraft {

	tflog.Info(ctx, "[getCMWAFPolicyRequestDraft] Fetching WAF Policy Request Draft ")
	var wafpolicyReqDraft bigipnextsdk.CMWAFPolicyRequestDraft
	wafpolicyReqDraft.Name = data.Name.ValueString()
	wafpolicyReqDraft.Description = data.Description.ValueString()
	wafpolicyReqDraft.EnforecementMode = data.EnforecementMode.ValueString()
	wafpolicyReqDraft.ApplicationLanguage = data.ApplicationLanguage.ValueString()
	wafpolicyReqDraft.TemplateName = data.TemplateName.ValueString()

	var tagsList []string
	for _, val := range data.Tags.Elements() {
		var ss string
		// tflog.Info(ctx, "tag name %v")
		_ = json.Unmarshal([]byte(val.String()), &ss)
		tagsList = append(tagsList, ss)
	}
	wafpolicyReqDraft.Tags = tagsList
	var botdefenseModel BotDefenseModel
	botdefenseModel.Enabled = types.BoolValue(true)
	if !data.BotDefense.IsNull() && !data.BotDefense.IsUnknown() {
		diag := data.BotDefense.As(ctx, &botdefenseModel, basetypes.ObjectAsOptions{})
		if diag.HasError() {
			tflog.Error(ctx, fmt.Sprintf("[getCMWAFPolicyRequestDraft] diag Error: %+v", diag.Errors()))
		}
	}
	wafpolicyReqDraft.Declaration.Policy.BotDefense.Settings.IsEnabled = (botdefenseModel.Enabled).ValueBool()

	var ipintelligenceModel IpIntelligenceModel
	ipintelligenceModel.Enabled = types.BoolValue(true)
	if !data.IpIntelligence.IsNull() && !data.IpIntelligence.IsUnknown() {
		diag := data.IpIntelligence.As(ctx, &ipintelligenceModel, basetypes.ObjectAsOptions{})
		if diag.HasError() {
			tflog.Error(ctx, fmt.Sprintf("[getCMWAFPolicyRequestDraft] diag Error: %+v", diag.Errors()))
		}
	}
	wafpolicyReqDraft.Declaration.Policy.IpIntelligence.Enabled = (ipintelligenceModel.Enabled).ValueBool()

	var dosprotectionModel DosProtectionModel
	dosprotectionModel.Enabled = types.BoolValue(false)
	if !data.DosProtection.IsNull() && !data.DosProtection.IsUnknown() {
		diag := data.DosProtection.As(ctx, &dosprotectionModel, basetypes.ObjectAsOptions{})
		if diag.HasError() {
			tflog.Error(ctx, fmt.Sprintf("[getCMWAFPolicyRequestDraft] diag Error: %+v", diag.Errors()))
		}
	}
	wafpolicyReqDraft.Declaration.Policy.DosProtection.Enabled = (dosprotectionModel.Enabled).ValueBool()

	var blockingsettingdModel BlockingSettingsModel
	blockingsettingdModel.Enabled = types.BoolValue(true)
	if !data.BlockingSettings.IsNull() && !data.BlockingSettings.IsUnknown() {
		diag := data.BlockingSettings.As(ctx, &blockingsettingdModel, basetypes.ObjectAsOptions{})
		if diag.HasError() {
			tflog.Error(ctx, fmt.Sprintf("[getCMWAFPolicyRequestDraft] diag Error: %+v", diag.Errors()))
		}
	}
	wafpolicyReqDraft.Declaration.Policy.Name = data.Name.ValueString()
	wafpolicyReqDraft.Declaration.Policy.Description = data.Description.ValueString()
	wafpolicyReqDraft.Declaration.Policy.Template.Name = data.TemplateName.ValueString()

	var violationsList []bigipnextsdk.Violation
	var violation bigipnextsdk.Violation
	violation.Name = "VIOL_THREAT_CAMPAIGN"
	violation.Alarm = blockingsettingdModel.Enabled.ValueBool()
	violation.Block = blockingsettingdModel.Enabled.ValueBool()
	violation.Description = "Threat Campaign detected"
	violationsList = append(violationsList, violation)

	wafpolicyReqDraft.Declaration.Policy.BlockingSettings.Violations = violationsList

	tflog.Info(ctx, fmt.Sprintf("[getCMWAFPolicyRequestDraft] wafpolicyReqDraft:%+v\n", wafpolicyReqDraft))
	return &wafpolicyReqDraft
}

func (r *NextCMWAFPolicyResource) WafPolicyModeltoState(ctx context.Context, respData interface{}, data *NextCMWAFPolicyResourceModel) {

	tflog.Info(ctx, fmt.Sprintf("WafPolicyModeltoState \t name: %+v", respData.(map[string]interface{})["name"]))

	data.Name = types.StringValue(respData.(map[string]interface{})["name"].(string))
	description, ok := respData.(map[string]interface{})["description"]
	if ok {
		data.Description = types.StringValue(description.(string))
	}
	data.Tags, _ = types.ListValueFrom(ctx, types.StringType, respData.(map[string]interface{})["tags"])
	data.EnforecementMode = types.StringValue(respData.(map[string]interface{})["enforcement_mode"].(string))
	data.ApplicationLanguage = types.StringValue(respData.(map[string]interface{})["application_language"].(string))
	data.Id = types.StringValue(respData.(map[string]interface{})["id"].(string))

	_, ok = respData.(map[string]interface{})["declaration"]
	if ok {
		_, ok = respData.(map[string]interface{})["declaration"].(map[string]interface{})["policy"]
		// fetching the settings of bot-defense, ip-intelligence, dos-protection, blocking-settings from Policy
		if ok {
			_, ok = respData.(map[string]interface{})["declaration"].(map[string]interface{})["policy"].(map[string]interface{})["bot-defense"]
			if ok {
				var botdefenseModel BotDefenseModel
				// diag := data.BotDefense.As(ctx, &botdefenseModel, basetypes.ObjectAsOptions{})
				// if diag.HasError() {
				// 	tflog.Error(ctx, fmt.Sprintf("[WafPolicyModeltoState] diag Error: %+v", diag.Errors()))
				// }
				botdefenseModel.Enabled = types.BoolValue(respData.(map[string]interface{})["declaration"].(map[string]interface{})["policy"].(map[string]interface{})["bot-defense"].(map[string]interface{})["settings"].(map[string]interface{})["isEnabled"].(bool))
			}

			_, ok = respData.(map[string]interface{})["declaration"].(map[string]interface{})["policy"].(map[string]interface{})["ip-intelligence"]
			if ok {
				var ipintelligenceModel IpIntelligenceModel
				// diag := data.IpIntelligence.As(ctx, &ipintelligenceModel, basetypes.ObjectAsOptions{})
				// if diag.HasError() {
				// 	tflog.Error(ctx, fmt.Sprintf("[WafPolicyModeltoState] diag Error: %+v", diag.Errors()))
				// }
				ipintelligenceModel.Enabled = types.BoolValue(respData.(map[string]interface{})["declaration"].(map[string]interface{})["policy"].(map[string]interface{})["ip-intelligence"].(map[string]interface{})["enabled"].(bool))
			}

			_, ok = respData.(map[string]interface{})["declaration"].(map[string]interface{})["policy"].(map[string]interface{})["template"]
			if ok {
				data.TemplateName = types.StringValue(respData.(map[string]interface{})["declaration"].(map[string]interface{})["policy"].(map[string]interface{})["template"].(map[string]interface{})["name"].(string))
			}

			_, ok = respData.(map[string]interface{})["declaration"].(map[string]interface{})["policy"].(map[string]interface{})["dos-protection"]
			if ok {
				var dosprotectionModel DosProtectionModel
				// diag := data.DosProtection.As(ctx, &dosprotectionModel, basetypes.ObjectAsOptions{})
				// if diag.HasError() {
				// 	tflog.Error(ctx, fmt.Sprintf("[WafPolicyModeltoState] diag Error: %+v", diag.Errors()))
				// }
				dosprotectionModel.Enabled = types.BoolValue(respData.(map[string]interface{})["declaration"].(map[string]interface{})["policy"].(map[string]interface{})["dos-protection"].(map[string]interface{})["enabled"].(bool))
			}
			_, ok = respData.(map[string]interface{})["declaration"].(map[string]interface{})["policy"].(map[string]interface{})["blocking-settings"]
			if ok {
				var blockingsettingdModel BlockingSettingsModel
				// diag := data.BlockingSettings.As(ctx, &blockingsettingdModel, basetypes.ObjectAsOptions{})
				// if diag.HasError() {
				// 	tflog.Error(ctx, fmt.Sprintf("[WafPolicyModeltoState] diag Error: %+v", diag.Errors()))
				// }
				for _, violation := range respData.(map[string]interface{})["declaration"].(map[string]interface{})["policy"].(map[string]interface{})["blocking-settings"].(map[string]interface{})["violations"].([]interface{}) {
					// var i Violation
					if violation.(map[string]interface{})["Name"] == "VIOL_THREAT_CAMPAIGN" {
						_, ok = violation.(map[string]interface{})["block"]
						if ok {
							blockingsettingdModel.Enabled = types.BoolValue(violation.(map[string]interface{})["block"].(bool))
						}
					}
				}
			}
		}
	}
}
