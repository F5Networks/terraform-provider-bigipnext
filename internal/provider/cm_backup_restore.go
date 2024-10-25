package provider

import (
	"context"
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
	_ resource.Resource                = &CMBackupRestoreResource{}
	_ resource.ResourceWithImportState = &CMBackupRestoreResource{}
	// mutex sync.Mutex
)

func NewCMBackupRestoreResource() resource.Resource {
	return &CMBackupRestoreResource{}
}

type CMBackupRestoreResource struct {
	client *bigipnextsdk.BigipNextCM
}

type CMBackupRestoreResourceModel struct {
	EncryptionPassword types.String `tfsdk:"encryption_password"`
	Name               types.String `tfsdk:"name"`
	FileName           types.String `tfsdk:"file_name"`
	Schedule           types.Object `tfsdk:"schedule"`
	Backup             types.Bool   `tfsdk:"backup"`
	Frequency          types.String `tfsdk:"frequency"`
	Type               types.String `tfsdk:"type"`
	DaysOfTheWeekToRun types.List   `tfsdk:"days_of_the_week_to_run"`
	DayOfTheMonthToRun types.Int64  `tfsdk:"day_of_the_month_to_run"`
	Scheduled          types.Bool   `tfsdk:"scheduled"`
	Id                 types.String `tfsdk:"id"`
}

type ScheduleModel struct {
	StartAt types.String `tfsdk:"start_at"`
	EndAt   types.String `tfsdk:"end_at"`
}

func (r *CMBackupRestoreResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cm_backup_restore"
}

func (r *CMBackupRestoreResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource used to manage Backup/Restore resources onto BIG-IP Next CM.",
		Attributes: map[string]schema.Attribute{
			"encryption_password": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Encryption password for the backup to be created. Password should be minimum of 8 characters",
				// Computed:            true,
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The unique name of the backup file. Actual File Name is auto-generated in the case of Instant Backup",
			},
			"schedule": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Specifies whether backup is to be scheduled or not.",
				Attributes: map[string]schema.Attribute{
					"start_at": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Specifies Start time of the backup. Example: 2019-08-24T14:15:22Z",
					},
					"end_at": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Specifies End time of the backup. Example: 2019-08-24T14:15:22Z",
					},
				},
			},
			"type": schema.StringAttribute{
				//  Optional:            true,
				Computed:            true,
				MarkdownDescription: "Type of the Backup",
			},
			"file_name": schema.StringAttribute{
				//  Optional:            true,
				Computed:            true,
				MarkdownDescription: "Name of the backup file generate in the case of instant backup",
			},
			"backup": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "Specifies whether backup (if True) or restore (if false) is to be done on CM",
				// Default:             booldefault.StaticBool(true),
			},
			"scheduled": schema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Specifies whether backup is scheduled or not",
			},
			"frequency": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Specifies what is the frequency. Example : Daily, Monthly, Weekly",
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"Monthly", "Weekly", "Daily"}...),
				},
			},
			"days_of_the_week_to_run": schema.ListAttribute{
				Optional:            true,
				ElementType:         types.Int64Type,
				MarkdownDescription: "Specifies Day of the week on backup has been scheduled. 0-Sunday, 1-Monday and so on",
			},
			"day_of_the_month_to_run": schema.Int64Attribute{
				Optional:            true,
				MarkdownDescription: "Specifies From which Day of the month backup should start.",
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

func (r *CMBackupRestoreResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client, resp.Diagnostics = toBigipNextCMProvider(req.ProviderData)
}

func (r *CMBackupRestoreResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var resCfg *CMBackupRestoreResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &resCfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "[CREATE] CMBackupRestoreResource")

	if resCfg.Backup.ValueBool() {
		// backup
		reqDraft, scheduled, err := getCMBackupDraft(ctx, resCfg)
		if err != nil { // coverage-ignore
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to Create Backup, got error: %s", err))
			return
		}

		tflog.Info(ctx, fmt.Sprintf("[CREATE] :%+v\n", reqDraft))

		file_name, draftID, err := r.client.BackUpCM(reqDraft, "POST")
		if err != nil { // coverage-ignore
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to Create Backup, got error: %s", err))
			return
		}
		tflog.Info(ctx, fmt.Sprintf("[CREATE] draftID:%+v\n", draftID))

		resCfg.FileName = types.StringValue(file_name)
		resCfg.Scheduled = types.BoolValue(scheduled)
		resCfg.Id = types.StringValue(draftID)
		if !scheduled {
			resCfg.Type = types.StringValue("Light")
		} else {
			resCfg.Type = types.StringValue(reqDraft.ScheduleType)
		}
	} else {
		// restore
		file_name := resCfg.Name.ValueString()
		backupConfig, err := r.client.GetBackUpConfig(file_name, false, true)
		if err != nil { // coverage-ignore
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to Restore, got error: %s", err))
			return
		}
		cmRestoreDraft := &bigipnextsdk.CMRestoreRequestDraft{}
		cmRestoreDraft.EncryptionPassword = resCfg.EncryptionPassword.ValueString()
		cmRestoreDraft.FileId = backupConfig.(map[string]interface{})["_embedded"].(map[string]interface{})["backups"].([]interface{})[0].(map[string]interface{})["file_id"].(string)
		err = r.client.RestoreCM(cmRestoreDraft)
		if err != nil { // coverage-ignore
			resp.Diagnostics.AddError("Error", fmt.Sprintf("Failed to Restore, got error: %s", err))
			return
		}
		resCfg.Scheduled = types.BoolValue(false)
		resCfg.FileName = types.StringValue("")
		resCfg.Id = types.StringValue("")
		resCfg.Type = types.StringValue("Restore")
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, resCfg)...)
}

func (r *CMBackupRestoreResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var stateCfg *CMBackupRestoreResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateCfg)...)
	if resp.Diagnostics.HasError() { // coverage-ignore
		return
	}
	backup := stateCfg.Backup.ValueBool()

	if backup {
		id := stateCfg.Id.ValueString()
		name := stateCfg.Name.ValueString()
		file_name := stateCfg.FileName.ValueString()
		scheduled := stateCfg.Scheduled.ValueBool()

		tflog.Info(ctx, fmt.Sprintf("Reading Backup Config : %s", id))

		var backupConfig interface{}
		var err error
		if scheduled {
			backupConfig, err = r.client.GetBackUpConfig(id, scheduled, false)
		} else {
			backupConfig, err = r.client.GetBackUpConfig(file_name, scheduled, false)
		}

		if err != nil { // coverage-ignore
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to Read Backup, got error: %s", err))
			return
		}

		tflog.Info(ctx, fmt.Sprintf("[Read] : %+v", backupConfig))
		r.backUpConfigModeltoState(ctx, backupConfig, scheduled, name, stateCfg)
		resp.Diagnostics.Append(resp.State.Set(ctx, &stateCfg)...)
	}
}

func (r *CMBackupRestoreResource) backUpConfigModeltoState(ctx context.Context, respData interface{}, scheduled bool, name string, data *CMBackupRestoreResourceModel) {

	tflog.Info(ctx, fmt.Sprintf("backupResourceModelToState:%+v", respData))

	data.Name = types.StringValue(name)
	data.Scheduled = types.BoolValue(scheduled)

	if !scheduled {
		data.FileName = types.StringValue(respData.(map[string]interface{})["_embedded"].(map[string]interface{})["files"].([]interface{})[0].(map[string]interface{})["file_name"].(string))
		data.Id = types.StringValue(respData.(map[string]interface{})["_embedded"].(map[string]interface{})["files"].([]interface{})[0].(map[string]interface{})["id"].(string))
	} else {
		data.FileName = data.Name
		var scheduleModel = ScheduleModel{}
		data.Id = types.StringValue(respData.(map[string]interface{})["id"].(string))

		_, ok := respData.(map[string]interface{})["days_of_the_week"]
		if ok {
			data.DaysOfTheWeekToRun, _ = types.ListValueFrom(ctx, types.Int64Type, respData.(map[string]interface{})["days_of_the_week"].(map[string]interface{})["days_of_the_week_to_run"])
		}

		_, ok = respData.(map[string]interface{})["day_and_time_of_month"]
		if ok {
			data.DayOfTheMonthToRun = types.Int64Value(int64(respData.(map[string]interface{})["day_and_time_of_month"].(map[string]interface{})["day_of_the_month_to_run"].(float64)))
		}

		scheduleModel.StartAt = types.StringValue(respData.(map[string]interface{})["start_date"].(string))
		_, ok = respData.(map[string]interface{})["end_date"]
		if ok {
			tflog.Info(ctx, fmt.Sprintf("backupResourceModelToState:%+v", "End date is coming"))
			scheduleModel.EndAt = types.StringValue(respData.(map[string]interface{})["end_date"].(string))
		}
		scheduleAttributes := map[string]attr.Type{
			"start_date": types.StringType,
			"end_data":   types.StringType,
		}
		scheduleObjectValue, err := types.ObjectValueFrom(ctx, scheduleAttributes, scheduleModel)
		if err != nil {
			return
		}
		data.Schedule = scheduleObjectValue
	}
}

func (r *CMBackupRestoreResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var resCfg *CMBackupRestoreResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &resCfg)...)

	if resp.Diagnostics.HasError() { // coverage-ignore
		return
	}

	if !resCfg.Backup.ValueBool() { // coverage-ignore
		resp.Diagnostics.AddError("Client Error", "Restore can't be updated.")
		return
	}

	if resCfg.Frequency.IsNull() || resCfg.Frequency.IsUnknown() { // coverage-ignore
		resp.Diagnostics.AddError("Client Error", "Only Scheduled Backup can be updated.")
		return
	}

	tflog.Info(ctx, fmt.Sprintf("[UPDATE] Name: %s", resCfg.Name.ValueString()))

	reqDraft, scheduled, err := getCMBackupDraft(ctx, resCfg)
	if err != nil { // coverage-ignore
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to Update Backup, got error: %s", err))
		return
	}

	if !scheduled { // coverage-ignore
		resp.Diagnostics.AddError("Client Error", "Invalid Input for Scheduling Backup.")
		return
	}
	reqDraft.Id = resCfg.Id.ValueString()

	tflog.Info(ctx, fmt.Sprintf("[UPDATE] :%+v\n", reqDraft))

	file_name, draftID, err := r.client.BackUpCM(reqDraft, "PUT")
	if err != nil { // coverage-ignore
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to Update Backup, got error: %s", err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[UPDATE] draftID:%+v\n", draftID))
	resCfg.FileName = types.StringValue(file_name)
	resCfg.Scheduled = types.BoolValue(true)
	resCfg.Id = types.StringValue(draftID)
	resCfg.Type = types.StringValue(reqDraft.ScheduleType)

	resp.Diagnostics.Append(resp.State.Set(ctx, resCfg)...)
}

func (r *CMBackupRestoreResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var stateCfg *CMBackupRestoreResourceModel
	if resp.Diagnostics.HasError() { // coverage-ignore
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &stateCfg)...)
	id := stateCfg.Id.ValueString()
	backup := stateCfg.Backup.ValueBool()

	if backup {
		scheduled := stateCfg.Scheduled.ValueBool()

		tflog.Info(ctx, fmt.Sprintf("Deleting Backup : %s", id))

		err := r.client.DeleteBackup(id, scheduled)
		if err != nil { // coverage-ignore
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to Delete Backup, got error: %s", err))
			return
		}
	}
	stateCfg.Id = types.StringValue("")
}

func (r *CMBackupRestoreResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getCMBackupDraft(ctx context.Context, data *CMBackupRestoreResourceModel) (*bigipnextsdk.CMBackupRequestDraft, bool, error) {

	cmBackupConfig := &bigipnextsdk.CMBackupRequestDraft{}
	cmBackupConfig.EncryptionPassword = data.EncryptionPassword.ValueString()
	cmBackupConfig.Name = data.Name.ValueString()
	scheduled := false
	if !data.Schedule.IsNull() && !data.Schedule.IsUnknown() {
		// scheduling
		var providerModel ScheduleModel
		diag := data.Schedule.As(ctx, &providerModel, basetypes.ObjectAsOptions{})
		if diag.HasError() { // coverage-ignore
			tflog.Error(ctx, fmt.Sprintf("getCMBackupDraft diag Error: %+v", diag.Errors()))
		}
		scheduled = true
		cmBackupConfig.Schedule.StartAt = providerModel.StartAt.ValueString()

		if !providerModel.EndAt.IsNull() && !providerModel.EndAt.IsUnknown() {
			cmBackupConfig.Schedule.EndAt = providerModel.EndAt.ValueString()
		}

		if data.Frequency.IsNull() || data.Frequency.IsUnknown() { // coverage-ignore
			return nil, true, fmt.Errorf("Frequency should be provided in order to schedule a backup")
		} else if data.Frequency == types.StringValue("Weekly") {

			if data.DaysOfTheWeekToRun.IsNull() || data.DaysOfTheWeekToRun.IsUnknown() { // coverage-ignore
				return nil, true, fmt.Errorf("Days of the week should be provided for Weekly Backup Schedule")
			}

			cmBackupConfig.ScheduleType = "DaysOfTheWeek"
			daysoftheweekList := make([]int64, 0)
			data.DaysOfTheWeekToRun.ElementsAs(ctx, &daysoftheweekList, false)
			daysoftheweek := &bigipnextsdk.DaysOfTheWeek{}
			daysoftheweek.HourToRunOn = 10
			daysoftheweek.MinuteToRunOn = 30
			daysoftheweek.Interval = 1
			daysoftheweek.DaysOfTheWeekToRun = daysoftheweekList
			cmBackupConfig.DaysOfTheWeek = daysoftheweek

		} else if data.Frequency == types.StringValue("Monthly") {

			if data.DayOfTheMonthToRun.IsNull() || data.DayOfTheMonthToRun.IsUnknown() { // coverage-ignore
				return nil, true, fmt.Errorf("Day of the Month should be provided for Monthly Backup Schedule")
			}

			cmBackupConfig.ScheduleType = "DayAndTimeOfTheMonth"
			dayandtimeofthemonth := &bigipnextsdk.DayAndTimeOfTheMonth{}
			dayandtimeofthemonth.DayOfTheMonthToRun = int(data.DayOfTheMonthToRun.ValueInt64())
			dayandtimeofthemonth.HourToRunOn = 10
			dayandtimeofthemonth.Interval = 1
			dayandtimeofthemonth.MinuteToRunOn = 30
			cmBackupConfig.DayAndTimeOfTheMonth = dayandtimeofthemonth

		} else if data.Frequency == types.StringValue("Daily") {
			cmBackupConfig.ScheduleType = "BasicWithInterval"

			basicwithinterval := &bigipnextsdk.BasicWithInterval{}
			basicwithinterval.IntervalToRun = 24
			basicwithinterval.IntervalUnit = "HOUR"

			cmBackupConfig.BasicWithInterval = basicwithinterval
		}

	} else {
		scheduled = false
		cmBackupConfig.Type = "light"
	}

	tflog.Info(ctx, fmt.Sprintf("getCMBackupDraft:%+v\n", cmBackupConfig))
	return cmBackupConfig, scheduled, nil
}
