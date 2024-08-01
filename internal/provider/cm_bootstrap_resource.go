package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	bigipnextsdk "gitswarm.f5net.com/terraform-providers/bigipnext"
)

var (
	_ resource.Resource
	_ resource.ResourceWithImportState
)

func NewCMNextBootstrapResource() resource.Resource {
	return &CMNextBootstrapResource{}
}

type CMNextBootstrapResource struct {
	client *bigipnextsdk.BigipNextCM
}

type ExternalStorage struct {
	StorageType    types.String `tfsdk:"storage_type"`
	StorageAddress types.String `tfsdk:"storage_address"`
	StoragePath    types.String `tfsdk:"storage_path"`
	CMStorageDir   types.String `tfsdk:"cm_storage_dir"`
	Username       types.String `tfsdk:"username"`
	Password       types.String `tfsdk:"password"`
}

type CMNextBootstrapResourceModel struct {
	Id               types.String `tfsdk:"id"`
	RunSetup         types.Bool   `tfsdk:"run_setup"`
	ExternalStorage  types.Object `tfsdk:"external_storage"`
	BootstrapStatus  types.String `tfsdk:"bootstrap_status"`
	BootstrapTimeout types.Int64  `tfsdk:"bootstrap_timeout"`
}

func (r *CMNextBootstrapResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cm_bootstrap"
}

func (r *CMNextBootstrapResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource used for bootstrapping Central Manager\n\n" +
			"~> **NOTE** This resource does not support update and delete. When doing `terraform destroy` it will only remove the resource from the state",
		Attributes: map[string]schema.Attribute{
			"run_setup": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Run setup on Central Manager",
			},
			"external_storage": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "External storage configuration",
				Attributes: map[string]schema.Attribute{
					"storage_type": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Type of external storage. Supported values are NFS and SAMBA",
						Validators: []validator.String{
							stringvalidator.OneOf([]string{"NFS", "SAMBA"}...),
						},
					},
					"storage_address": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "IP Address of the external storage",
					},
					"storage_path": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Directory path that is mounted on the external storage server",
					},
					"cm_storage_dir": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Folder name created on the external storage server to store Central Manager data",
					},
					"username": schema.StringAttribute{
						MarkdownDescription: "Username to access the external storage, required if storage type is SAMBA",
						Optional:            true,
					},
					"password": schema.StringAttribute{
						MarkdownDescription: "Password to access the external storage, required if storage type is SAMBA",
						Optional:            true,
					},
				},
			},
			"bootstrap_timeout": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Timeout for the bootstrap operation",
			},
			"bootstrap_status": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Status of the bootstrap operation",
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "ID of the resource",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
		},
	}
}

func (r *CMNextBootstrapResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client, resp.Diagnostics = toBigipNextCMProvider(req.ProviderData)
}

func (r *CMNextBootstrapResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var resCfg CMNextBootstrapResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &resCfg)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !resCfg.ExternalStorage.IsNull() {
		var externalStorageModel ExternalStorage
		diag := resCfg.ExternalStorage.As(ctx, &externalStorageModel, basetypes.ObjectAsOptions{})

		if diag.HasError() {
			return
		}

		externalStorage := getExternalStorageData(&externalStorageModel)
		res, err := r.client.AddExternalStorage(externalStorage)
		if err != nil { // coverage-ignore
			resp.Diagnostics.AddError("Failed to add external storage:", err.Error())
			return
		}

		externalStoreResp := &bigipnextsdk.CMExternalStorageResp{}

		tflog.Info(ctx, "External storage setup response: "+res)
		err = json.Unmarshal([]byte(res), externalStoreResp)
		if err != nil { // coverage-ignore
			resp.Diagnostics.AddError("Failed to unmarshal external storage response:", err.Error())
			return
		}

		if externalStoreResp.Status.Setup != "SUCCESSFUL" { // coverage-ignore
			tflog.Error(ctx, "Failed to setup external storage")
			resp.Diagnostics.AddError("Failed to setup external storage", externalStoreResp.Status.FailureMessage)
			return
		}

		typeMap := map[string]attr.Type{
			"storage_type":    types.StringType,
			"storage_address": types.StringType,
			"storage_path":    types.StringType,
			"cm_storage_dir":  types.StringType,
			"username":        types.StringType,
			"password":        types.StringType,
		}
		valMap := map[string]attr.Value{
			"storage_type":    types.StringValue(externalStoreResp.Spec.StorageType),
			"storage_address": types.StringValue(externalStoreResp.Spec.StorageAddress),
			"storage_path":    types.StringValue(externalStoreResp.Spec.StorageSharePath),
		}

		if externalStoreResp.Spec.StorageType == "SAMBA" {
			valMap["username"] = types.StringValue(externalStoreResp.Spec.StorageUser.Username)
			valMap["password"] = types.StringValue(externalStoreResp.Spec.StorageUser.Password)
		} else {
			valMap["username"] = types.StringNull()
			valMap["password"] = types.StringNull()
		}
		if externalStoreResp.Spec.StorageShareDir != "" {
			valMap["cm_storage_dir"] = types.StringValue(externalStoreResp.Spec.StorageShareDir)
		} else {
			if !externalStorageModel.CMStorageDir.IsNull() {
				valMap["cm_storage_dir"] = externalStorageModel.CMStorageDir
			} else {
				valMap["cm_storage_dir"] = types.StringNull()
			}
		}

		resCfg.ExternalStorage, _ = types.ObjectValue(typeMap, valMap)
	}

	var cmBootstrapStatus string
	if resCfg.RunSetup.ValueBool() {
		var timeout int64
		if resCfg.BootstrapTimeout.IsNull() {
			timeout = 600
		} else { // coverage-ignore
			timeout = resCfg.BootstrapTimeout.ValueInt64()
		}
		res, err := r.client.BootstrapCM(timeout)

		if err != nil { // coverage-ignore
			resp.Diagnostics.AddError("Failed to bootstrap Central Manager:", err.Error())
			return
		}

		bootStrapResp := &bigipnextsdk.BootstrapCMResp{}
		err = json.Unmarshal([]byte(res), bootStrapResp)
		if err != nil { // coverage-ignore
			resp.Diagnostics.AddError("Failed to unmarshal Central Manager bootstrap response:", err.Error())
			return
		}

		cmBootstrapStatus = bootStrapResp.Status

		tflog.Info(ctx, "Central Manager bootstrap response: "+res)
	}

	id := extractIPFromUrl(r.client.Host)

	resCfg.Id = types.StringValue(fmt.Sprintf("setup-%s", id))
	resCfg.BootstrapStatus = types.StringValue(cmBootstrapStatus)
	resp.Diagnostics.Append(resp.State.Set(ctx, &resCfg)...)
}

func (r *CMNextBootstrapResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var stateCfg *CMNextBootstrapResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateCfg)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.GetCMExternalStorage()
	if err != nil { // coverage-ignore
		resp.Diagnostics.AddError("Failed to read external storage status", err.Error())
	}
	externalStorageResp := &bigipnextsdk.CMExternalStorageResp{}
	err = json.Unmarshal([]byte(res), externalStorageResp)
	if err != nil { // coverage-ignore
		resp.Diagnostics.AddError("Failed to unmarshal external storage response", err.Error())
	}

	typeMap := map[string]attr.Type{
		"storage_type":    types.StringType,
		"storage_address": types.StringType,
		"storage_path":    types.StringType,
		"cm_storage_dir":  types.StringType,
		"username":        types.StringType,
		"password":        types.StringType,
	}
	valMap := map[string]attr.Value{
		"storage_type":    types.StringValue(externalStorageResp.Spec.StorageType),
		"storage_address": types.StringValue(externalStorageResp.Spec.StorageAddress),
		"storage_path":    types.StringValue(externalStorageResp.Spec.StorageSharePath),
	}

	if externalStorageResp.Spec.StorageType == "SAMBA" {
		valMap["username"] = types.StringValue(externalStorageResp.Spec.StorageUser.Username)
		valMap["password"] = types.StringValue(externalStorageResp.Spec.StorageUser.Password)
	} else {
		valMap["username"] = types.StringNull()
		valMap["password"] = types.StringNull()
	}
	if externalStorageResp.Spec.StorageShareDir != "" {
		valMap["cm_storage_dir"] = types.StringValue(externalStorageResp.Spec.StorageShareDir)
	} else {
		valMap["cm_storage_dir"] = types.StringNull()
	}

	stateCfg.ExternalStorage, _ = types.ObjectValue(typeMap, valMap)

	bootstrap, err := r.client.GetCMBootstrap()
	if err != nil { // coverage-ignore
		resp.Diagnostics.AddError("Failed to read the bootstrap status", err.Error())
	}
	bootstrapResp := &bigipnextsdk.BootstrapCMResp{}
	err = json.Unmarshal([]byte(bootstrap), bootstrapResp)
	if err != nil { // coverage-ignore
		resp.Diagnostics.AddError("Failed to unmarshal the bootstrap response", err.Error())
	}

	id := extractIPFromUrl(r.client.Host)

	stateCfg.Id = types.StringValue(fmt.Sprintf("setup-%s", id))
	stateCfg.BootstrapStatus = types.StringValue(bootstrapResp.Status)
}

func (r *CMNextBootstrapResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "this resource does not support update operation")
}

func (r *CMNextBootstrapResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "this resource does not support delete operation, it can oly be removed from the state")
	resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(""))
}

func getExternalStorageData(plan *ExternalStorage) *bigipnextsdk.CMExternalStorage {
	externalStorage := &bigipnextsdk.CMExternalStorage{}

	externalStorage.StorageType = plan.StorageType.ValueString()
	externalStorage.StorageAddress = plan.StorageAddress.ValueString()
	externalStorage.StorageSharePath = plan.StoragePath.ValueString()

	if plan.CMStorageDir.IsNull() {
		externalStorage.StorageShareDir = ""
	} else {
		externalStorage.StorageShareDir = plan.CMStorageDir.ValueString()
	}

	if !plan.Username.IsNull() {
		externalStorage.StorageUser = &bigipnextsdk.CMExternalStorageUser{
			Username: plan.Username.ValueString(),
			Password: plan.Password.ValueString(),
		}
	}

	return externalStorage
}
