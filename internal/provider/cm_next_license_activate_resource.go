package provider

import (
	"context"
	"fmt"
	"strings"

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
	_ resource.Resource                = &CMNextLicenseActivateResource{}
	_ resource.ResourceWithImportState = &CMNextLicenseActivateResource{}
)

func NewCMNextLicenseActivateResource() resource.Resource {
	return &CMNextLicenseActivateResource{}
}

type CMNextLicenseActivateResource struct {
	client *bigipnextsdk.BigipNextCM
}

type CMNextLicenseActivateResourceModel struct {
	Instances []InstanceActivateModel `tfsdk:"instances"`
	Id        types.String            `tfsdk:"id"`
}

type InstanceActivateModel struct {
	InstanceAddress types.String `tfsdk:"instance_address"`
	JwtId           types.String `tfsdk:"jwt_id"`
}

func (r *CMNextLicenseActivateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cm_activate_instance_license"
}

func (r *CMNextLicenseActivateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Resource used for Activate/Deactivate License for Instances on Central Manager Using JWT Token",
		Attributes: map[string]schema.Attribute{
			"instances": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"instance_address": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "IP Address of the instance to activate the license",
						},
						"jwt_id": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "JWT ID to be used to activate the license",
						},
					},
				},
				Required:            true,
				MarkdownDescription: "List of instances to activate the license",
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

func (r *CMNextLicenseActivateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client, resp.Diagnostics = toBigipNextCMProvider(req.ProviderData)
}

func (r *CMNextLicenseActivateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var resCfg *CMNextLicenseActivateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &resCfg)...)
	if resp.Diagnostics.HasError() { // coverage-ignore
		return
	}
	tflog.Info(ctx, "[CREATE] Activate License for Instances on Central Manager Using JWT Token")
	providerConfig := getCMNextLicenseActivateConfig(ctx, r.client, resCfg)
	respData, err := r.client.PostActivateLicense(providerConfig)
	if err != nil { // coverage-ignore
		resp.Diagnostics.AddError("Failed to Activate license", fmt.Sprintf("%+v", err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[CREATE] License Task IDs :%+v\n", respData))
	var deviceID []string
	for _, value := range providerConfig {
		deviceID = append(deviceID, value.DigitalAssetId)
	}
	tflog.Info(ctx, fmt.Sprintf("Device ID : %+v", deviceID))
	// join all keys into a single string
	deviceString := strings.Join(deviceID, ",")
	resCfg.Id = types.StringValue(deviceString)
	resp.Diagnostics.Append(resp.State.Set(ctx, resCfg)...)
}

func (r *CMNextLicenseActivateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var stateCfg *CMNextLicenseActivateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateCfg)...)
	if resp.Diagnostics.HasError() { // coverage-ignore
		return
	}
	id := stateCfg.Id.ValueString()
	tflog.Info(ctx, fmt.Sprintf("Instance IDs : %+v", id))
	deviceIDs := strings.Split(id, ",")
	deactivateReq := &bigipnextsdk.LicenseDeactivaeReq{}
	digitalAssetID := deviceIDs
	deactivateReq.DigitalAssetIds = digitalAssetID

	licenseInfo, err := r.client.PostLicenseInfo(deactivateReq)
	if err != nil { // coverage-ignore
		resp.Diagnostics.AddError("Failed to get License Info", fmt.Sprintf("%+v", err))
		return
	}
	// get the license info by loop over map
	var licenseStatus []string
	for _, val := range licenseInfo.(map[string]interface{}) {
		licenseStatus = append(licenseStatus, val.(map[string]interface{})["deviceLicenseStatus"].(map[string]interface{})["licenseStatus"].(string))
		tflog.Info(ctx, fmt.Sprintf("License Info : %+v", val.(map[string]interface{})["deviceLicenseStatus"].(map[string]interface{})["licenseStatus"]))
	}
	tflog.Info(ctx, fmt.Sprintf("Instance License Info : %+v", licenseStatus))
	// diags := resp.State.SetAttribute(ctx, path.Root("license_status"), licenseStatus)
	// resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateCfg)...)
}

// func (r *CMNextLicenseActivateResource) NextJwtTokenResourceModeltoState(ctx context.Context, respData interface{}, data *CMNextLicenseActivateResourceModel) {
// 	tflog.Debug(ctx, fmt.Sprintf("respData  %+v", respData))
// }

func (r *CMNextLicenseActivateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// var resCfg *CMNextLicenseActivateResourceModel
	// resp.Diagnostics.Append(req.Plan.Get(ctx, &resCfg)...)

	// if resp.Diagnostics.HasError() {
	// 	return
	// }
	tflog.Info(ctx, "[UPDATE] Update Call on Activating License Not Supported!!!!")
	// resp.Diagnostics.Append(resp.State.Set(ctx, &resCfg)...)
}

func (r *CMNextLicenseActivateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var stateCfg *CMNextLicenseActivateResourceModel
	if resp.Diagnostics.HasError() { // coverage-ignore
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &stateCfg)...)
	id := stateCfg.Id.ValueString()
	deviceIDs := strings.Split(id, ",")
	deactivateReq := &bigipnextsdk.LicenseDeactivaeReq{}
	digitalAssetID := deviceIDs
	deactivateReq.DigitalAssetIds = digitalAssetID
	res, err := r.client.PostDeactivateLicense(deactivateReq)
	if err != nil { // coverage-ignore
		resp.Diagnostics.AddError("Failed to deactivate license", fmt.Sprintf("%+v", err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Deactivate License Response : %+v", res))
	stateCfg.Id = types.StringValue("")
}

func (r *CMNextLicenseActivateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) { // coverage-ignore
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getCMNextLicenseActivateConfig(ctx context.Context, p *bigipnextsdk.BigipNextCM, data *CMNextLicenseActivateResourceModel) []*bigipnextsdk.LicenseReq {
	listLicenseRequest := []*bigipnextsdk.LicenseReq{}
	for _, val := range data.Instances {
		licenseRequest := &bigipnextsdk.LicenseReq{}
		licenseRequest.JwtId = val.JwtId.ValueString()
		deviceID, _ := p.GetDeviceIdByIp(val.InstanceAddress.ValueString())
		licenseRequest.DigitalAssetId = *deviceID
		tflog.Info(ctx, fmt.Sprintf("licenseRequest:%+v", licenseRequest))
		listLicenseRequest = append(listLicenseRequest, licenseRequest)
	}
	return listLicenseRequest
}
