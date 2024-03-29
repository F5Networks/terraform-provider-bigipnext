package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	bigipnextsdk "gitswarm.f5net.com/terraform-providers/bigipnext"
)

// Ensure BigipNextCMProvider satisfies various provider interfaces.
var _ provider.Provider = &BigipNextCMProvider{}

// BigipNextCMProvider defines the provider implementation
type BigipNextCMProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// BigipNextCMProviderModel describes the provider data model.
type BigipNextCMProviderModel struct {
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	// PlatformType types.String `tfsdk:"platform_type"`
	Port types.Int64 `tfsdk:"port"`
}

func (p *BigipNextCMProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "bigipnext"
	resp.Version = p.version
}

func (p *BigipNextCMProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provider plugin to interact with BIG-IP Next Central Manager(CM) Using OpenAPI",
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				MarkdownDescription: "URI for BigipNext Device. May also be provided via `BIGIPNEXT_HOST` environment variable.",
				Optional:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username for BigipNext Device. May also be provided via `BIGIPNEXT_USERNAME` environment variable.",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password for BigipNext Device. May also be provided via `BIGIPNEXT_PASSWORD` environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			// "platform_type": schema.StringAttribute{
			// 	MarkdownDescription: "Provider Host Platform type. Indicates provider is `bigipnext_cm` or `bigipnext_ve`.Default is `bigipnext_cm`.",
			// 	Optional:            true,
			// 	//Default:             stringdefault.StaticString("bigipnext_cm"),
			// 	Validators: []validator.String{
			// 		stringvalidator.OneOf([]string{"bigipnext_cm", "bigipnext_ve"}...),
			// 	},
			// },
			"port": schema.Int64Attribute{
				MarkdownDescription: "Port Number to be used to make API calls to HOST",
				Optional:            true,
			},
		},
	}
}

func (p *BigipNextCMProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring BigipNextCM client")

	// Retrieve provider data from configuration
	var config BigipNextCMProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("\n----bigipnextCmConfig :%+v", config))

	host := os.Getenv("BIGIPNEXT_HOST")
	username := os.Getenv("BIGIPNEXT_USERNAME")
	password := os.Getenv("BIGIPNEXT_PASSWORD")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}
	ctx = tflog.SetField(ctx, "bigipnext_host", host)
	ctx = tflog.SetField(ctx, "bigipnext_username", username)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "bigipnext_password")

	bigipnextCmConfig := &bigipnextsdk.BigipNextCMReqConfig{
		Host:     host,
		User:     username,
		Password: password,
		Port:     443,
	}
	tflog.Debug(ctx, fmt.Sprintf("bigipnextCmConfig client:%+v", bigipnextCmConfig))
	client, err := bigipnextsdk.CmNewSession(bigipnextCmConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create bigipnext CM Client",
			"An unexpected error occurred when creating the bigipnext client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"bigipnext Client Error: "+err.Error(),
		)
		return
	}
	resp.DataSourceData = client
	resp.ResourceData = client

	// if (config.PlatformType.IsNull() && !config.PlatformType.IsUnknown()) || config.PlatformType.ValueString() == "bigipnext_cm" {
	// 	// Example client configuration for data sources and resources
	// 	bigipnextCmConfig := &bigipnextsdk.BigipNextCMReqConfig{
	// 		Host:     host,
	// 		User:     username,
	// 		Password: password,
	// 		Port:     443,
	// 	}
	// 	tflog.Debug(ctx, fmt.Sprintf("bigipnextCmConfig client:%+v", bigipnextCmConfig))
	// 	client, err := bigipnextsdk.CmNewSession(bigipnextCmConfig)
	// 	if err != nil {
	// 		resp.Diagnostics.AddError(
	// 			"Unable to Create bigipnext CM Client",
	// 			"An unexpected error occurred when creating the bigipnext client. "+
	// 				"If the error is not clear, please contact the provider developers.\n\n"+
	// 				"bigipnext Client Error: "+err.Error(),
	// 		)
	// 		return
	// 	}
	// 	resp.DataSourceData = client
	// 	resp.ResourceData = client
	// }
	// if !config.PlatformType.IsNull() && !config.PlatformType.IsUnknown() && config.PlatformType.ValueString() == "bigipnext_ve" {
	// 	// Example client configuration for data sources and resources
	// 	bigipnextConfig := &bigipnextsdk.BigipNextConfig{
	// 		Host:     host,
	// 		User:     username,
	// 		Password: password,
	// 		Port:     5443,
	// 	}

	// 	tflog.Info(ctx, fmt.Sprintf("bigipnextConfig client:%+v", bigipnextConfig))

	// 	client, err := bigipnextsdk.NewSession(bigipnextConfig)
	// 	if err != nil {
	// 		resp.Diagnostics.AddError(
	// 			"Unable to Create bigipnext Client",
	// 			"An unexpected error occurred when creating the bigipnext client. "+
	// 				"If the error is not clear, please contact the provider developers.\n\n"+
	// 				"bigipnext Client Error: "+err.Error(),
	// 		)
	// 		return
	// 	}
	// 	resp.DataSourceData = client
	// 	resp.ResourceData = client
	// }

	tflog.Info(ctx, "Configured BIGIPNEXTCM client", map[string]any{"success": true})
}

func (p *BigipNextCMProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewNextCMAS3DeployResource,
		NewNextCMFastHttpResource,
		NewNextCMBackupRestoreResource,
		NewNextCMFastTemplateResource,
		NewNextCMCertificateResource,
		NewNextCMImportCertificateResource,
		NewNextCMDeviceProviderResource,
		NewNextDeployVmwareResource,
		NewNextDeployF5osResource,
		NewNextHAResource,
		NewNextGlobalResiliencyResource,
	}
}

func (p *BigipNextCMProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDeviceInventorySource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &BigipNextCMProvider{
			version: version,
		}
	}
}

// toProvider can be used to cast a generic provider.Provider reference to this specific provider.
// This is ideally used in DataSourceType.NewDataSource and ResourceType.NewResource calls.
func toBigipNextCMProvider(in any) (*bigipnextsdk.BigipNextCM, diag.Diagnostics) {
	if in == nil {
		return nil, nil
	}

	var diags diag.Diagnostics

	p, ok := in.(*bigipnextsdk.BigipNextCM)

	if !ok {
		diags.AddError(
			"Unexpected Provider Instance Type",
			fmt.Sprintf("While creating the data source or resource, an unexpected provider type (%T) was received. "+
				"This is always a bug in the provider code and should be reported to the provider developers.", in,
			),
		)
		return nil, diags
	}

	return p, diags
}

//// hashForState computes the hexadecimal representation of the SHA1 checksum of a string.
//// This is used by most resources/data-sources here to compute their Unique Identifier (ID).
//func hashForState(value string) string {
//	if value == "" {
//		return ""
//	}
//	hash := sha1.Sum([]byte(strings.TrimSpace(value)))
//	return hex.EncodeToString(hash[:])
//}
