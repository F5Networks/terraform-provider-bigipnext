package provider

import (
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	bigipnextsdk "gitswarm.f5net.com/terraform-providers/bigipnext"
)

var (
	_ resource.Resource                = &NextCMBackupRestoreResource{}
	_ resource.ResourceWithImportState = &NextCMBackupRestoreResource{}
)

func NewNextCMBackupRestoreResource() resource.Resource {
	return &NextCMBackupRestoreResource{}
}

type NextCMBackupRestoreResource struct {
	client *bigipnextsdk.BigipNextCM
}

type NextCMBackupRestoreResourceModel struct {
	FileName       types.String `tfsdk:"file_name"`
	Password       types.String `tfsdk:"backup_password"`
	DeviceIp       types.String `tfsdk:"device_ip"`
	DeviceHostname types.String `tfsdk:"device_hostname"`
	Operation      types.String `tfsdk:"operation"`
	Id             types.String `tfsdk:"instance_id"`
	Backup         types.String `tfsdk:"backup_date"`
	Restore        types.String `tfsdk:"restore_date"`
	Timeout        types.Int64  `tfsdk:"timeout"`
}

func (r *NextCMBackupRestoreResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cm_backup_restore"
}

func (r *NextCMBackupRestoreResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource used to manage(CRUD) backup and restore of BIG-IP Next instances on BIG-IP CM.",
		Attributes: map[string]schema.Attribute{
			"file_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the backup file to create, or use.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`.*\.tar\.gz$`),
						"file_name must have .tar.gz extension specified",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"device_ip": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "IP Address of the device managed by BIG-IP Next CM.\nParameter required for create operations.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("device_hostname")),
				},
			},
			"device_hostname": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Hostname of the device managed by BIG-IP Next CM.\nParameter required for create operations.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("device_ip")),
				},
			},
			"backup_password": schema.StringAttribute{
				MarkdownDescription: "User password for the backup file.",
				Sensitive:           true,
				Required:            true,
			},
			"operation": schema.StringAttribute{
				MarkdownDescription: "Type of operation to perform.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("restore", "backup"),
				},
			},
			"instance_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "UUID of the NEXT instance which config was backed up or restored.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"backup_date": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The timestamp when backup file was created. In ISO 8601 format",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"restore_date": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The timestamp when restore operation was performed. In ISO 8601 format",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"timeout": schema.Int64Attribute{
				MarkdownDescription: "The number of seconds to wait for backup or restore operation to complete.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(360),
			},
		},
	}
}

func (r *NextCMBackupRestoreResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client, resp.Diagnostics = toBigipNextCMProvider(req.ProviderData)
}

func (r *NextCMBackupRestoreResource) GetDeviceId(data *NextCMBackupRestoreResourceModel) (deviceId *string, err error) {
	if data.DeviceIp.ValueString() == "" && data.DeviceHostname.ValueString() == "" { // coverage-ignore
		return nil, fmt.Errorf("the 'device_ip' or 'device_hostname' parameter must be specified")
	}
	if data.DeviceIp.ValueString() == "" {
		device := data.DeviceHostname.ValueString()
		deviceId, err := r.client.GetDeviceIdByHostname(device)
		if err != nil { // coverage-ignore
			return nil, err
		}
		return deviceId, nil
	} else {
		device := data.DeviceIp.ValueString()
		deviceId, err = r.client.GetDeviceIdByIp(device)
		if err != nil { // coverage-ignore
			return nil, err
		}
		return deviceId, nil
	}
}

func (r *NextCMBackupRestoreResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *NextCMBackupRestoreResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() { // coverage-ignore
		return
	}
	deviceId, err := r.GetDeviceId(data)
	if err != nil { // coverage-ignore
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to obtain device id, got error: %s", err))
		return
	}
	config := getCreateBackupRestoreConfig(ctx, req, resp)

	if data.Operation.ValueString() == "backup" {
		mutex.Lock()
		respData, err := r.client.BackupTenant(deviceId, config, int(data.Timeout.ValueInt64()))
		if err != nil { // coverage-ignore
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to create config backup, got error: %s", err))
			return
		}
		tflog.Info(ctx, fmt.Sprintf("Backup operation response:%+v", respData))
		r.backupResourceModelToState(ctx, respData, data)
		mutex.Unlock()
	}
	if data.Operation.ValueString() == "restore" {
		mutex.Lock()
		respData, err := r.client.RestoreTenant(deviceId, config, int(data.Timeout.ValueInt64()))
		if err != nil { // coverage-ignore
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to restore config backup, got error: %s", err))
			return
		}
		tflog.Info(ctx, fmt.Sprintf("Restore operation response:%+v", respData))
		r.backupResourceModelToState(ctx, respData, data)
		mutex.Unlock()
	}
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NextCMBackupRestoreResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) { // coverage-ignore
	var data *NextCMBackupRestoreResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	deviceId, err := r.GetDeviceId(data)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to obtain device id, got error: %s", err))
		return
	}
	config := getUpdateBackupRestoreConfig(ctx, req, resp)

	if data.Operation.ValueString() == "backup" {
		err := r.client.DeleteBackupFile(data.FileName.ValueStringPointer())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete config backup, got error: %s", err))
			return
		}
		mutex.Lock()
		respData, err := r.client.BackupTenant(deviceId, config, int(data.Timeout.ValueInt64()))
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to create config backup, got error: %s", err))
			return
		}
		tflog.Info(ctx, fmt.Sprintf("Backup operation response:%+v", respData))
		r.backupResourceModelToState(ctx, respData, data)
		mutex.Unlock()
	}
	if data.Operation.ValueString() == "restore" {
		mutex.Lock()
		respData, err := r.client.RestoreTenant(deviceId, config, int(data.Timeout.ValueInt64()))
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to restore config backup, got error: %s", err))
			return
		}
		tflog.Info(ctx, fmt.Sprintf("Restore operation response:%+v", respData))
		r.backupResourceModelToState(ctx, respData, data)
		mutex.Unlock()
	}
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NextCMBackupRestoreResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *NextCMBackupRestoreResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() { // coverage-ignore
		return
	}

	err := r.client.DeleteBackupFile(data.FileName.ValueStringPointer())
	if err != nil { // coverage-ignore
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to Delete backup file, got error: %s", err))
		return
	}
}

func (r *NextCMBackupRestoreResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NextCMBackupRestoreResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() { // coverage-ignore
		return
	}

	//respByte, err := r.client.GetTenant(data.Name.ValueString())
	respByte, err := r.client.GetBackupFile(data.FileName.ValueStringPointer())
	if err != nil { // coverage-ignore
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to find backup file, got error: %s", err))
		return
	}
	r.backupResourceReadModelToState(ctx, respByte, data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NextCMBackupRestoreResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) { // coverage-ignore
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getCreateBackupRestoreConfig(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) *bigipnextsdk.BackupRestoreTenantRequest {
	var data *NextCMBackupRestoreResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	config := bigipnextsdk.BackupRestoreTenantRequest{}
	// due to bug in CM, during backup we need to trip extensions so that desired filename with correct extensions is created
	// once this is fixed we will remove
	if data.Operation.ValueString() == "backup" {
		name := data.FileName.ValueString()
		config.FileName = removeExtensions(name)
	} else {
		config.FileName = data.FileName.ValueString()
	}
	config.Password = data.Password.ValueString()
	return &config
}

func getUpdateBackupRestoreConfig(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) *bigipnextsdk.BackupRestoreTenantRequest { // coverage-ignore
	var data *NextCMBackupRestoreResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	config := bigipnextsdk.BackupRestoreTenantRequest{}
	name := data.FileName.ValueString()
	config.FileName = removeExtensions(name)
	config.Password = data.Password.ValueString()
	return &config
}

func (r *NextCMBackupRestoreResource) backupResourceReadModelToState(ctx context.Context, respData *bigipnextsdk.TenantBackupFile, data *NextCMBackupRestoreResourceModel) {
	tflog.Info(ctx, fmt.Sprintf("backupResourceReadModelToState:%+v", respData))
	data.FileName = types.StringValue(respData.FileName)
	data.Id = types.StringValue(respData.InstanceId)
	data.Backup = types.StringValue(respData.FileDate)
}

func (r *NextCMBackupRestoreResource) backupResourceModelToState(ctx context.Context, respData *bigipnextsdk.TenantBackupRestoreTaskStatus, data *NextCMBackupRestoreResourceModel) {
	tflog.Info(ctx, fmt.Sprintf("backupResourceModelToState:%+v", respData))
	data.Id = types.StringValue(respData.InstanceId)
	if data.Operation == types.StringValue("restore") {
		data.Restore = types.StringValue(respData.CompletionDate)
		data.Backup = types.StringValue("")
	}
	if data.Operation == types.StringValue("backup") {
		data.Backup = types.StringValue(respData.CreationDate)
		data.Restore = types.StringValue("")
	}
}

func removeExtensions(fileName string) string {
	ext := filepath.Ext(fileName)
	if ext != "" {
		return removeExtensions(strings.TrimSuffix(fileName, filepath.Ext(fileName)))
	}
	return fileName
}
