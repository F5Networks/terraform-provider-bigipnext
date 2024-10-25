package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	bigipnextsdk "gitswarm.f5net.com/terraform-providers/bigipnext"
)

var (
	_ resource.Resource = &CMNextUpgradeResource{}
)

func NewCMNextUpgradeResource() resource.Resource {
	return &CMNextUpgradeResource{}
}

type CMNextUpgradeResource struct {
	client *bigipnextsdk.BigipNextCM
}

type CMNextUpgradeResourceModel struct {
	ImageName         types.String `tfsdk:"image_name"`
	SignatureFilename types.String `tfsdk:"signature_filename"`
	NextInstanceIP    types.String `tfsdk:"next_instance_ip"`
	UpgradeType       types.String `tfsdk:"upgrade_type"`
	PartitionAddress  types.String `tfsdk:"partition_address"`
	PartitionPort     types.Int64  `tfsdk:"partition_port"`
	PartitionUsername types.String `tfsdk:"partition_username"`
	PartitionPassword types.String `tfsdk:"partition_password"`
	TenantName        types.String `tfsdk:"tenant_name"`
	Timeout           types.Int64  `tfsdk:"timeout"`
	Id                types.String `tfsdk:"id"`
}

func (r *CMNextUpgradeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cm_next_upgrade"
}

func (r *CMNextUpgradeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource used to upgrade BIG-IP Next instances that are managed by the BIG-IP Central Manager.",
		Attributes: map[string]schema.Attribute{
			"image_name": schema.StringAttribute{
				MarkdownDescription: "The name of the BIG-IP Next image file that is to be used to upgrade the BIG-IP Next instance, in case the upgrade_type is 'appliance'" +
					" the value should be the name of the image file that is present on the Velos partition or rSeries with with the BIG-IP Next instance is to be upgraded.",
				Required: true,
			},
			"signature_filename": schema.StringAttribute{
				MarkdownDescription: "The name of the signature file that is to be used to verify the image file, it is required when upgrade_type is 've'.",
				Optional:            true,
				Validators: []validator.String{
					requiredIf(path.Root("upgrade_type"), "ve"),
				},
			},
			"next_instance_ip": schema.StringAttribute{
				MarkdownDescription: "The IP address of the BIG-IP Next instance that is to be upgraded.",
				Required:            true,
			},
			"upgrade_type": schema.StringAttribute{
				MarkdownDescription: "The type of upgrade that is to be performed, it can be either 've' or 'appliance'.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("ve", "appliance"),
				},
			},
			"partition_address": schema.StringAttribute{
				MarkdownDescription: "The IP address of the Velos partition or rSeries on which the BIG-IP Next instance is to be upgraded, it is required when upgrade_type is 'appliance'.",
				Optional:            true,
				Validators: []validator.String{
					requiredIf(path.Root("upgrade_type"), "appliance"),
				},
			},
			"partition_port": schema.Int64Attribute{
				MarkdownDescription: "The port number of the Velos partition or rSeries on which the BIG-IP Next instance is to be upgraded, it is required when upgrade_type is 'appliance'.",
				Optional:            true,
				Validators: []validator.Int64{
					requiredIf(path.Root("upgrade_type"), "appliance"),
				},
			},
			"partition_username": schema.StringAttribute{
				MarkdownDescription: "The username of the Velos partition or rSeries on which the BIG-IP Next instance is to be upgraded, it is required when upgrade_type is 'appliance'.",
				Optional:            true,
				Validators: []validator.String{
					requiredIf(path.Root("upgrade_type"), "appliance"),
				},
			},
			"partition_password": schema.StringAttribute{
				MarkdownDescription: "The password of the Velos partition or rSeries on which the BIG-IP Next instance is to be upgraded, it is required when upgrade_type is 'appliance'.",
				Optional:            true,
				Validators: []validator.String{
					requiredIf(path.Root("upgrade_type"), "appliance"),
				},
			},
			"tenant_name": schema.StringAttribute{
				MarkdownDescription: "The name of the BIG-IP Next tenant that is to be upgraded, it is required when upgrade_type is 'appliance'.",
				Optional:            true,
				Validators: []validator.String{
					requiredIf(path.Root("upgrade_type"), "appliance"),
				},
			},
			"timeout": schema.Int64Attribute{
				MarkdownDescription: "The time in seconds to wait for the upgrade process to complete, the default value is 300 seconds.",
				Default:             int64default.StaticInt64(300),
				Optional:            true,
				Computed:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique Identifier for the resource.",
				Computed:            true,
			},
		},
	}
}

func (r *CMNextUpgradeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client, resp.Diagnostics = toBigipNextCMProvider(req.ProviderData)
}

func (r *CMNextUpgradeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var resCfg *CMNextUpgradeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &resCfg)...)

	nextInstanceId, err := r.client.GetNextInstanceID(resCfg.NextInstanceIP.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"error while fetching the next instance id",
			err.Error(),
		)
		return
	}

	tflog.Debug(
		ctx,
		fmt.Sprintf("Instance id of %v is %v: ", resCfg.NextInstanceIP.ValueString(), nextInstanceId),
	)

	var upgradeTaskId string
	if strings.ToLower(resCfg.UpgradeType.ValueString()) == "ve" {

		image_name, signature_name, err := r.client.GetImageAndSignatureName(
			nextInstanceId,
			resCfg.ImageName.ValueString(),
			resCfg.SignatureFilename.ValueString(),
		)
		if err != nil {
			resp.Diagnostics.AddError(
				"error while fetching the image and signature name",
				err.Error(),
			)
			return
		}

		upgradeTaskId, err = r.client.UpgradeVE(nextInstanceId, image_name, signature_name)
		if err != nil {
			resp.Diagnostics.AddError(
				"error while initiating the upgrade process",
				err.Error(),
			)
			return
		}

	} else if strings.ToLower(resCfg.UpgradeType.ValueString()) == "appliance" {
		body := make(map[string]interface{})

		body["partition_address"] = resCfg.PartitionAddress.ValueString()
		body["partition_port"] = resCfg.PartitionPort.ValueInt64()
		body["partition_username"] = resCfg.PartitionUsername.ValueString()
		body["partition_password"] = resCfg.PartitionPassword.ValueString()
		body["tenant_name"] = resCfg.TenantName.ValueString()
		body["image_name"] = resCfg.ImageName.ValueString()

		upgradeTaskId, err = r.client.UpgradeNextInstanceAppliance(nextInstanceId, body)
		if err != nil {
			resp.Diagnostics.AddError(
				"error while initiating the upgrade process",
				err.Error(),
			)
			return
		}
	}

	tflog.Debug(ctx, "Upgrade task id: "+upgradeTaskId)
	status, details, err := r.client.WaitForNextInstanceUpgrade(upgradeTaskId, resCfg.Timeout.ValueInt64())
	if err != nil { // coverage-ignore
		resp.Diagnostics.AddError(
			"error while waiting for the upgrade to complete",
			err.Error(),
		)
		return
	}

	if status != "completed" { // coverage-ignore
		resp.Diagnostics.AddError(
			"upgrade process taking too long",
			"upgrade process not complete yet: "+details,
		)
		return
	}

	if status == "failed" { // coverage-ignore
		resp.Diagnostics.AddError(
			"upgrade process failed",
			"upgrade process failed with error: "+details,
		)
		return
	}

	resCfg.Id = types.StringValue(nextInstanceId)
	resp.Diagnostics.Append(resp.State.Set(ctx, resCfg)...)
}

func (r *CMNextUpgradeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var stateCfg *CMNextUpgradeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateCfg)...)

	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateCfg)...)
}

func (r *CMNextUpgradeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var resCfg CMNextUpgradeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &resCfg)...)

	nextInstanceId, err := r.client.GetNextInstanceID(resCfg.NextInstanceIP.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"error while fetching the next instance id",
			err.Error(),
		)
		return
	}

	tflog.Debug(
		ctx,
		fmt.Sprintf("Instance id of %v is %v: ", resCfg.NextInstanceIP.ValueString(), nextInstanceId),
	)

	var upgradeTaskId string

	if strings.ToLower(resCfg.UpgradeType.ValueString()) == "ve" {

		image_name, signature_name, err := r.client.GetImageAndSignatureName(
			nextInstanceId,
			resCfg.ImageName.ValueString(),
			resCfg.SignatureFilename.ValueString(),
		)
		if err != nil {
			resp.Diagnostics.AddError(
				"error while fetching the image and signature name",
				err.Error(),
			)
			return
		}

		upgradeTaskId, err = r.client.UpgradeVE(nextInstanceId, image_name, signature_name)
		if err != nil {
			resp.Diagnostics.AddError(
				"error while initiating the upgrade process",
				err.Error(),
			)
			return
		}
	} else if strings.ToLower(resCfg.UpgradeType.ValueString()) == "appliance" {
		body := make(map[string]interface{})
		body["partition_address"] = resCfg.PartitionAddress.ValueString()
		body["partition_port"] = resCfg.PartitionPort.ValueInt64()
		body["partition_username"] = resCfg.PartitionUsername.ValueString()
		body["partition_password"] = resCfg.PartitionPassword.ValueString()
		body["tenant_name"] = resCfg.TenantName.ValueString()
		body["image_name"] = resCfg.ImageName.ValueString()

		upgradeTaskId, err = r.client.UpgradeNextInstanceAppliance(nextInstanceId, body)
		if err != nil {
			resp.Diagnostics.AddError(
				"error while initiating the upgrade process",
				err.Error(),
			)
			return
		}
	}

	status, details, err := r.client.WaitForNextInstanceUpgrade(upgradeTaskId, resCfg.Timeout.ValueInt64())
	if err != nil { // coverage-ignore
		resp.Diagnostics.AddError(
			"error while waiting for the upgrade to complete",
			err.Error(),
		)
		return
	}

	if status != "completed" { // coverage-ignore
		resp.Diagnostics.AddError(
			"upgrade process failed or taking too long",
			"upgrade process not complete yet: "+details,
		)
		return
	}

	if status == "failed" { // coverage-ignore
		resp.Diagnostics.AddError(
			"upgrade process failed",
			"upgrade process failed with error: "+details,
		)
		return
	}

	resCfg.Id = types.StringValue(nextInstanceId)
	resp.Diagnostics.Append(resp.State.Set(ctx, resCfg)...)
}

func (r *CMNextUpgradeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.State.SetAttribute(ctx, path.Root("id"), types.StringValue(""))
}
