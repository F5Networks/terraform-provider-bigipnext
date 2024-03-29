package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	bigipnextsdk "gitswarm.f5net.com/terraform-providers/bigipnext"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ datasource.DataSource = &DeviceInventorySource{}
	// _ datasource.DataSourceWithConfigure = &DeviceInventorySource{}
)

func NewDeviceInventorySource() datasource.DataSource {
	return &DeviceInventorySource{}
}

// DeviceInventorySource defines the data source implementation.
type DeviceInventorySource struct {
	client *bigipnextsdk.BigipNextCM
}

// DeviceInventorySourceModel describes the data source data model.
type DeviceInventorySourceModel struct {
	ID              types.String `tfsdk:"id"`
	DeviceInventory types.String `tfsdk:"device_inventory"`
}

func (d *DeviceInventorySource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cm_device_inventory"
}

func (d *DeviceInventorySource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Get information about the VLANs on f5os platform.\n\n" +
			"Use this data source to get information, such as vlan",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier of this data source: hashing of the certificates in the chain.",
			},
			"device_inventory": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier of this data source: hashing of the certificates in the chain.",
			},
		},
	}
}

func (d *DeviceInventorySource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client, resp.Diagnostics = toBigipNextCMProvider(req.ProviderData)
}

func (d *DeviceInventorySource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DeviceInventorySourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	deviceInventory, err := d.client.GetDeviceInventory()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Device inventory, got error: %s", err))
		return
	}

	data.ID = types.StringValue(fmt.Sprintf("%v", deviceInventory.Embedded.Devices[0].Id))
	data.DeviceInventory = types.StringValue(fmt.Sprintf("%+v", deviceInventory.Embedded.Devices[0]))

	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
