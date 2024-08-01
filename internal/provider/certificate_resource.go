package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	bigipnextsdk "gitswarm.f5net.com/terraform-providers/bigipnext"
)

var (
	_ resource.Resource                = &NextCMCertificateResource{}
	_ resource.ResourceWithImportState = &NextCMCertificateResource{}
	// mutex sync.Mutex
)

func NewNextCMCertificateResource() resource.Resource {
	return &NextCMCertificateResource{}
}

type NextCMCertificateResource struct {
	client *bigipnextsdk.BigipNextCM
}

type CertificateUpdateDraft struct {
	bigipnextsdk.CertificateRequestDraft
	ID string `json:"id"`
}

type NextCMCertificateResourceModel struct {
	Issuer                 types.String `tfsdk:"issuer"`
	Name                   types.String `tfsdk:"name"`
	CommonName             types.String `tfsdk:"common_name"`
	Division               types.List   `tfsdk:"division"`
	Organization           types.List   `tfsdk:"organization"`
	Locality               types.List   `tfsdk:"locality"`
	State                  types.List   `tfsdk:"state"`
	Country                types.List   `tfsdk:"country"`
	Email                  types.List   `tfsdk:"email"`
	SubjectAlternativeName types.String `tfsdk:"subject_alternative_name"`
	DurationInDays         types.Int64  `tfsdk:"duration_in_days"`
	KeyType                types.String `tfsdk:"key_type"`
	KeySecurityType        types.String `tfsdk:"key_security_type"`
	KeySize                types.Int64  `tfsdk:"key_size"`
	KeyCurveName           types.String `tfsdk:"key_curve_name"`
	KeyPassphrase          types.String `tfsdk:"key_passphrase"`
	AdministratorEmail     types.String `tfsdk:"administrator_email"`
	ChallengePassword      types.String `tfsdk:"challenge_password"`
	Id                     types.String `tfsdk:"id"`
}

func (r *NextCMCertificateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cm_certificate"
}

func (r *NextCMCertificateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"issuer": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "issuer details",
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"Self", "CA"}...),
				},
				Default: stringdefault.StaticString("Self"),
			},
			"common_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The fully qualified domain name of your server. The common_name of the certificate cannot be empty",
			},
			"duration_in_days": schema.Int64Attribute{
				MarkdownDescription: "duration in days",
				Required:            true,
			},
			"key_type": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Specifies Key type to be either `RSA` or `ECDSA`",
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"RSA", "ECDSA"}...),
				},
				Default: stringdefault.StaticString("RSA"),
			},
			"key_security_type": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Specifies whether key is password protected",
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"Normal"}...),
				},
				Default: stringdefault.StaticString("Normal"),
			},
			"key_size": schema.Int64Attribute{
				MarkdownDescription: "Size of key - the number of bits in a key used by a cryptographic algorithm. Supported key size for RSA - 2048, 3072, 4096",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(2048),
			},
			"division": schema.ListAttribute{
				MarkdownDescription: "The division of your organization handling the certificate. It is Array of strings",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"organization": schema.ListAttribute{
				MarkdownDescription: "The legal name of your organization. It is Array of strings",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"locality": schema.ListAttribute{
				MarkdownDescription: "The locality where your organization is located. It is Array of strings",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"state": schema.ListAttribute{
				MarkdownDescription: "The state where your organization is located. It is Array of strings",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"country": schema.ListAttribute{
				MarkdownDescription: "The country where your organization is located. An SSL certificate country code is a two-letter code that's used when you generate a CSR. It is Array of strings",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"email": schema.ListAttribute{
				MarkdownDescription: "An email address to contact your organization. It is Array of strings",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"subject_alternative_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "A SAN or subject alternative name is a structured way to indicate all of the domain names and IP addresses that are secured by the certificate",
			},
			"key_curve_name": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Supported curve names for ECDSA- secp384r1, prime256v1",
			},
			"key_passphrase": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "key passphrase, A passphrase is a word or phrase that protects private key files, It prevents unauthorized users from encrypting them. Usually it's just the secret encryption/decryption key used for Ciphers.",
			},
			"administrator_email": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "An administrator email to contact your organization",
			},
			"challenge_password": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "challenge password",
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

func (r *NextCMCertificateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client, resp.Diagnostics = toBigipNextCMProvider(req.ProviderData)
}

func (r *NextCMCertificateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var resCfg *NextCMCertificateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &resCfg)...)
	if resp.Diagnostics.HasError() { // coverage-ignore
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[CREATE] NextCMCertificateResource:%+v\n", resCfg.Name.ValueString()))

	reqDraft := getCertificateRequestDraft(ctx, resCfg)

	tflog.Info(ctx, "[CREATE] Posting Certificate")
	tflog.Info(ctx, fmt.Sprintf("[CREATE] :%+v\n", reqDraft))

	draftID, err := r.client.PostCertificateCreate(reqDraft, "CREATE")
	if err != nil { // coverage-ignore
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to Create Certificate, got error: %s", err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[CREATE] draftID:%+v\n", draftID))

	resCfg.Id = types.StringValue(draftID)
	resp.Diagnostics.Append(resp.State.Set(ctx, resCfg)...)
}

func (r *NextCMCertificateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var stateCfg *NextCMCertificateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateCfg)...)
	if resp.Diagnostics.HasError() { // coverage-ignore
		return
	}
	id := stateCfg.Id.ValueString()
	tflog.Info(ctx, fmt.Sprintf("Reading Certificate : %s", id))
	certData, err := r.client.GetNextCMCertificate(id)
	if err != nil { // coverage-ignore
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to Read Certificate, got error: %s", err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Certificate : %+v", certData))
	r.CertificateResourceModeltoState(ctx, certData, stateCfg)
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateCfg)...)
}

func (r *NextCMCertificateResource) CertificateResourceModeltoState(ctx context.Context, respData interface{}, data *NextCMCertificateResourceModel) {
	tflog.Info(ctx, fmt.Sprintf("CertificateResourceModeltoState \t key_size: %+v", int64(respData.(map[string]interface{})["key_size"].(float64))))
	data.KeySize = types.Int64Value(int64(respData.(map[string]interface{})["key_size"].(float64)))
	data.KeyType = types.StringValue(respData.(map[string]interface{})["key_type"].(string))
	data.CommonName = types.StringValue(respData.(map[string]interface{})["common_name"].(string))
	data.Issuer = types.StringValue(respData.(map[string]interface{})["issuer"].(string))
	data.Name = types.StringValue(respData.(map[string]interface{})["name"].(string))
	data.DurationInDays = types.Int64Value(int64(respData.(map[string]interface{})["duration_in_days"].(float64)))
}

func (r *NextCMCertificateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var resCfg *NextCMCertificateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &resCfg)...)

	if resp.Diagnostics.HasError() { // coverage-ignore
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[UPDATE]Posting Certificate: %s", resCfg.Name.ValueString()))

	reqDraft := getCertificateUpdateDraft(ctx, resCfg)
	reqDraft.ID = resCfg.Id.ValueString()

	tflog.Info(ctx, "[UPDATE] Posting Certificate")
	tflog.Info(ctx, fmt.Sprintf("[UPDATE] :%+v\n", reqDraft))

	draftID, err := r.client.PostCertificateCreate(reqDraft, "UPDATE")
	if err != nil { // coverage-ignore
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to Update Certificate, got error: %s", err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[CREATE] draftID:%+v\n", draftID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &resCfg)...)
}

func (r *NextCMCertificateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var stateCfg *NextCMCertificateResourceModel
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

func (r *NextCMCertificateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) { // coverage-ignore
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getCertificateRequestDraft(ctx context.Context, data *NextCMCertificateResourceModel) *bigipnextsdk.CertificateRequestDraft {
	var certificateReqDraft bigipnextsdk.CertificateRequestDraft
	certificateReqDraft.Name = data.Name.ValueString()
	certificateReqDraft.CommonName = data.CommonName.ValueString()
	certificateReqDraft.Issuer = data.Issuer.ValueString()
	divisionList := make([]string, 0, 1)
	data.Division.ElementsAs(ctx, &divisionList, false)
	certificateReqDraft.Division = divisionList
	organizationList := make([]string, 0, 1)
	data.Organization.ElementsAs(ctx, &organizationList, false)
	certificateReqDraft.Organization = organizationList
	localityList := make([]string, 0, 1)
	data.Locality.ElementsAs(ctx, &localityList, false)
	certificateReqDraft.Locality = localityList
	stateList := make([]string, 0, 1)
	data.State.ElementsAs(ctx, &stateList, false)
	certificateReqDraft.State = stateList
	countryList := make([]string, 0, 1)
	data.Country.ElementsAs(ctx, &countryList, false)
	certificateReqDraft.Country = countryList
	emailList := make([]string, 0, 1)
	data.Email.ElementsAs(ctx, &emailList, false)
	certificateReqDraft.Email = emailList
	certificateReqDraft.SubjectAlternativeName = data.SubjectAlternativeName.ValueString()
	certificateReqDraft.DurationInDays = int(data.DurationInDays.ValueInt64())
	certificateReqDraft.KeyType = data.KeyType.ValueString()
	certificateReqDraft.KeySecurityType = data.KeySecurityType.ValueString()
	certificateReqDraft.KeySize = int(data.KeySize.ValueInt64())
	certificateReqDraft.KeyCurveName = data.KeyCurveName.ValueString()
	certificateReqDraft.KeyPassphrase = data.KeyPassphrase.ValueString()
	certificateReqDraft.AdministratorEmail = data.AdministratorEmail.ValueString()
	// certificateReqDraft.ChallengePassword = data.AllowOverwrite.ValueBool()

	tflog.Info(ctx, fmt.Sprintf("certificateReqDraft:%+v\n", certificateReqDraft))
	return &certificateReqDraft
}

func getCertificateUpdateDraft(ctx context.Context, data *NextCMCertificateResourceModel) *CertificateUpdateDraft {
	var certificateReqDraft CertificateUpdateDraft
	certificateReqDraft.Name = data.Name.ValueString()
	certificateReqDraft.CommonName = data.CommonName.ValueString()
	certificateReqDraft.Issuer = data.Issuer.ValueString()
	divisionList := make([]string, 0, 1)
	data.Division.ElementsAs(ctx, &divisionList, false)
	certificateReqDraft.Division = divisionList
	organizationList := make([]string, 0, 1)
	data.Organization.ElementsAs(ctx, &organizationList, false)
	certificateReqDraft.Organization = organizationList
	localityList := make([]string, 0, 1)
	data.Locality.ElementsAs(ctx, &localityList, false)
	certificateReqDraft.Locality = localityList
	stateList := make([]string, 0, 1)
	data.State.ElementsAs(ctx, &stateList, false)
	certificateReqDraft.State = stateList
	countryList := make([]string, 0, 1)
	data.Country.ElementsAs(ctx, &countryList, false)
	certificateReqDraft.Country = countryList
	emailList := make([]string, 0, 1)
	data.Email.ElementsAs(ctx, &emailList, false)
	certificateReqDraft.Email = emailList
	certificateReqDraft.SubjectAlternativeName = data.SubjectAlternativeName.ValueString()
	certificateReqDraft.DurationInDays = int(data.DurationInDays.ValueInt64())
	certificateReqDraft.KeyType = data.KeyType.ValueString()
	certificateReqDraft.KeySecurityType = data.KeySecurityType.ValueString()
	certificateReqDraft.KeySize = int(data.KeySize.ValueInt64())
	certificateReqDraft.KeyCurveName = data.KeyCurveName.ValueString()
	certificateReqDraft.KeyPassphrase = data.KeyPassphrase.ValueString()
	certificateReqDraft.AdministratorEmail = data.AdministratorEmail.ValueString()
	// certificateReqDraft.ChallengePassword = data.AllowOverwrite.ValueBool()

	tflog.Info(ctx, fmt.Sprintf("certificateReqDraft:%+v\n", certificateReqDraft))
	return &certificateReqDraft
}
