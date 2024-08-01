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
	_ resource.Resource                = &NextCMCertificateResource{}
	_ resource.ResourceWithImportState = &NextCMCertificateResource{}
	// mutex sync.Mutex
)

func NewNextCMImportCertificateResource() resource.Resource {
	return &NextCMImportCertificateResource{}
}

type NextCMImportCertificateResource struct {
	client *bigipnextsdk.BigipNextCM
}

type NextCMImportCertificateResourceModel struct {
	Name           types.String `tfsdk:"name"`
	KeyPassphrase  types.String `tfsdk:"key_passphrase"`
	CertPassphrase types.String `tfsdk:"cert_passphrase"`
	KeyText        types.String `tfsdk:"key_text"`
	CertText       types.String `tfsdk:"cert_text"`
	ImportType     types.String `tfsdk:"import_type"`
	Id             types.String `tfsdk:"id"`
}

func (r *NextCMImportCertificateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cm_import_certitficate"
}

func (r *NextCMImportCertificateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource used to manage(CRUD) certificate management resources onto BIG-IP Next CM.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The unique user-given name of the certificate. Certificate names cannot contain spaces or special characters. The allowed characters are a-z, A-Z, 0-9, dot(.), dash (-) and underscore (_). Names starting with only a-z, A-Z.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"key_passphrase": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "key passphrase, A passphrase is a word or phrase that protects private key files, It prevents unauthorized users from encrypting them. Usually it's just the secret encryption/decryption key used for Ciphers.",
			},
			"cert_passphrase": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "cert passphrase, A passphrase is a word or phrase that protects files ",
			},
			"key_text": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "key content",
			},
			"cert_text": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "cert content",
			},
			"import_type": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Import Type, Value can be `PKCS12`",
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

func (r *NextCMImportCertificateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client, resp.Diagnostics = toBigipNextCMProvider(req.ProviderData)
}

func (r *NextCMImportCertificateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var resCfg *NextCMImportCertificateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &resCfg)...)
	if resp.Diagnostics.HasError() { // coverage-ignore
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[CREATE] NextCMImportCertificateResource:%+v\n", resCfg.Name.ValueString()))

	reqDraft := getImportCertificateRequestDraft(ctx, resCfg)
	reqDraft.Name = resCfg.Name.ValueString()
	// tflog.Info(ctx, fmt.Sprintf("[CREATE] :%+v\n", reqDraft))

	draftID, err := r.client.PostCertificateCreate(reqDraft, "IMPORT")
	if err != nil { // coverage-ignore
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to Import Certificate, got error: %s", err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[CREATE] draftID:%+v\n", draftID))
	resCfg.Id = types.StringValue(draftID)
	resp.Diagnostics.Append(resp.State.Set(ctx, resCfg)...)
	// tflog.Info(ctx, fmt.Sprintf("[CREATE]resonse save :%+v\n", resCfg))
}

func (r *NextCMImportCertificateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var stateCfg *NextCMImportCertificateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateCfg)...)
	if resp.Diagnostics.HasError() { // coverage-ignore
		return
	}
	id := stateCfg.Id.ValueString()
	tflog.Info(ctx, fmt.Sprintf("Reading Certificate : %s", id))
	keycertData, err := r.client.GetNextCMImportCertificateKeyData(id)
	if err != nil { // coverage-ignore
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to Read Certificate, Key Data, got error: %s", err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Certificate/Key Data : %+v", keycertData))

	certData, err := r.client.GetNextCMCertificate(id)
	if err != nil { // coverage-ignore
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to Read Certificate, got error: %s", err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Certificate : %+v", certData))

	r.CertificateResourceModeltoState(ctx, keycertData, certData, stateCfg)
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateCfg)...)
}

func (r *NextCMImportCertificateResource) CertificateResourceModeltoState(ctx context.Context, keycertData interface{}, respData interface{}, data *NextCMImportCertificateResourceModel) {
	data.Name = types.StringValue(respData.(map[string]interface{})["name"].(string))
	data.CertText = types.StringValue(keycertData.(map[string]interface{})["cert_data"].(string))
	// data.KeyText = types.StringValue(keycertData.(map[string]interface{})["key_data"].(string))
}

func (r *NextCMImportCertificateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var resCfg *NextCMImportCertificateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &resCfg)...)

	if resp.Diagnostics.HasError() { // coverage-ignore
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[UPDATE]Posting Certificate: %s", resCfg.Name.ValueString()))

	reqDraft := getImportCertificateRequestDraft(ctx, resCfg)
	reqDraft.Id = resCfg.Id.ValueString()

	tflog.Info(ctx, "[UPDATE] Posting Certificate")
	tflog.Info(ctx, fmt.Sprintf("[UPDATE] :%+v\n", reqDraft))

	draftID, err := r.client.PostCertificateCreate(reqDraft, "UPDATEIMPORT")
	if err != nil { // coverage-ignore
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to Update Certificate, got error: %s", err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[CREATE] draftID:%+v\n", draftID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &resCfg)...)
}

func (r *NextCMImportCertificateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var stateCfg *NextCMImportCertificateResourceModel
	if resp.Diagnostics.HasError() { // coverage-ignore
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &stateCfg)...)
	id := stateCfg.Id.ValueString()

	tflog.Info(ctx, fmt.Sprintf("Deleting Certificate : %s", id))

	err := r.client.DeleteNextCMCertificate(id)
	if err != nil { // coverage-ignore
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to Delete Certificate, got error: %s", err))
		return
	}
	stateCfg.Id = types.StringValue("")
}

func (r *NextCMImportCertificateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getImportCertificateRequestDraft(ctx context.Context, data *NextCMImportCertificateResourceModel) *bigipnextsdk.ImportCertificateRequestDraft {
	var certificateReqDraft bigipnextsdk.ImportCertificateRequestDraft
	certificateReqDraft.KeyPassphrase = data.KeyPassphrase.ValueString()
	certificateReqDraft.CertPassphrase = data.CertPassphrase.ValueString()
	certificateReqDraft.KeyText = data.KeyText.ValueString()
	certificateReqDraft.CertText = data.CertText.ValueString()
	certificateReqDraft.ImportType = data.ImportType.ValueString()

	tflog.Info(ctx, fmt.Sprintf("certificateReqDraft:%+v\n", certificateReqDraft))
	return &certificateReqDraft
}
