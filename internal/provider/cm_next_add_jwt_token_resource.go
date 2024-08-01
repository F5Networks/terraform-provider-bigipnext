package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	bigipnextsdk "gitswarm.f5net.com/terraform-providers/bigipnext"
)

var (
	_ resource.Resource                = &CMNextJwtTokenResource{}
	_ resource.ResourceWithImportState = &CMNextJwtTokenResource{}
)

func NewCMNextJwtTokenResource() resource.Resource {
	return &CMNextJwtTokenResource{}
}

type CMNextJwtTokenResource struct {
	client *bigipnextsdk.BigipNextCM
}

type CMNextJwtTokenResourceModel struct {
	TokenName          types.String `tfsdk:"token_name"`
	JwtToken           types.String `tfsdk:"jwt_token"`
	OrderType          types.String `tfsdk:"order_type"`
	SubscriptionExpiry types.String `tfsdk:"subscription_expiry"`
	Id                 types.String `tfsdk:"id"`
}

func (r *CMNextJwtTokenResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cm_add_jwt_token"
}

func (r *CMNextJwtTokenResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource used for add/copy JWT Token on Central Manager",
		Attributes: map[string]schema.Attribute{
			"token_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Nickname to be used to add the JWT token on Central Manager",
			},
			"jwt_token": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				MarkdownDescription: "JWT token to be added on Central Manager",
			},
			"order_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "JWT token to be added on Central Manager",
			},
			"subscription_expiry": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "JWT token to be added on Central Manager",
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

func (r *CMNextJwtTokenResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client, resp.Diagnostics = toBigipNextCMProvider(req.ProviderData)
}

func (r *CMNextJwtTokenResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var resCfg *CMNextJwtTokenResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &resCfg)...)
	if resp.Diagnostics.HasError() { // coverage-ignore
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[CREATE] CMNextJwtTokenResource:%+v\n", resCfg.TokenName.ValueString()))

	providerConfig := getCMNextJwtTokenConfig(ctx, resCfg)
	respData, err := r.client.PostLicenseToken(providerConfig)
	if err != nil { // coverage-ignore
		resp.Diagnostics.AddError("Failed to Create Jwt token", fmt.Sprintf(", got error: %s", err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[CREATE] JWT token ID :%+v\n", string(respData)))

	resCfg.Id = types.StringValue(string(respData))

	tokenInfo, err := r.client.GetLicenseToken(string(respData))
	if err != nil { // coverage-ignore
		resp.Diagnostics.AddError("Failed to Get JWT Token Info", fmt.Sprintf(", got error: %s", err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("JWT Token Info : %+v", tokenInfo))

	r.NextJwtTokenResourceModeltoState(ctx, tokenInfo, resCfg)
	resp.Diagnostics.Append(resp.State.Set(ctx, resCfg)...)
}

func (r *CMNextJwtTokenResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var stateCfg *CMNextJwtTokenResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateCfg)...)
	if resp.Diagnostics.HasError() { // coverage-ignore
		return
	}
	id := stateCfg.Id.ValueString()
	tflog.Info(ctx, fmt.Sprintf("JWT Token ID : %+v", id))

	tokenInfo, err := r.client.GetLicenseToken(id)
	if err != nil { // coverage-ignore
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to Get JWT Token Info, got error: %s", err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("JWT Token Info : %+v", tokenInfo))

	r.NextJwtTokenResourceModeltoState(ctx, tokenInfo, stateCfg)

	resp.Diagnostics.Append(resp.State.Set(ctx, &stateCfg)...)
}

func (r *CMNextJwtTokenResource) NextJwtTokenResourceModeltoState(ctx context.Context, respData interface{}, data *CMNextJwtTokenResourceModel) {
	tflog.Debug(ctx, fmt.Sprintf("respData  %+v", respData))
	data.OrderType = types.StringValue(respData.(map[string]interface{})["orderType"].(string))
	data.SubscriptionExpiry = types.StringValue(respData.(map[string]interface{})["subscriptionExpiry"].(string))
	data.TokenName = types.StringValue(respData.(map[string]interface{})["nickName"].(string))
}

func (r *CMNextJwtTokenResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var resCfg *CMNextJwtTokenResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &resCfg)...)

	if resp.Diagnostics.HasError() { // coverage-ignore
		return
	}
	tflog.Info(ctx, "[UPDATE] Updating JWT Tokens Not Supported!!!!")
	resCfg.OrderType = types.StringValue(resCfg.OrderType.ValueString())
	resCfg.SubscriptionExpiry = types.StringValue(resCfg.SubscriptionExpiry.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &resCfg)...)
}

func (r *CMNextJwtTokenResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var stateCfg *CMNextJwtTokenResourceModel
	if resp.Diagnostics.HasError() { // coverage-ignore
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &stateCfg)...)
	id := stateCfg.Id.ValueString()

	err := r.client.DeleteLicenseToken(id)
	if err != nil { // coverage-ignore
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to Delete JWT Token, got error: %s", err))
		return
	}
	stateCfg.Id = types.StringValue("")
}

func (r *CMNextJwtTokenResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getCMNextJwtTokenConfig(ctx context.Context, data *CMNextJwtTokenResourceModel) *bigipnextsdk.JWTRequestDraft {
	jwtRequest := &bigipnextsdk.JWTRequestDraft{}
	tflog.Info(ctx, fmt.Sprintf("jwtRequest:%+v\n", jwtRequest))
	jwtRequest.NickName = data.TokenName.ValueString()
	jwtRequest.JWT = data.JwtToken.ValueString()
	return jwtRequest
}
