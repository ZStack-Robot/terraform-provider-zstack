// Copyright (c) ZStack.io, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"
	"terraform-provider-zstack/zstack"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"zstack.io/zstack-sdk-go/pkg/client"
)

var (
	_ provider.Provider = &ZStackProvider{}
)

type ZStackProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}
type ZStackProviderModel struct {
	Host            types.String `tfsdk:"host"`
	AccountName     types.String `tfsdk:"accountname"`
	AccountPassword types.String `tfsdk:"accountpassword"`
	AccessKeyId     types.String `tfsdk:"accesskeyid"`
	AccessKeySecret types.String `tfsdk:"accesskeysecret"`
}

// Configure implements provider.Provider.
func (p *ZStackProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	tflog.Info(ctx, "Configuring ZStack client")

	//Retrieve provider data from configuration
	var config ZStackProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown ZStack Cloud API Host",
			"The provider cannt create the ZStack Cloud API client as an unknown configuration value for the ZStack Cloud API host."+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the ZSTACK_HOST environment variable.",
		)
	}

	if config.AccountName.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("accountname"),
			"Unknown ZStack API account Username",
			"Either target apply the source of the value first, set the value statically in the configuration, or use the ZSTACK_ACCOUNTNAME environment variable.",
		)
	}

	if config.AccountPassword.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("accountpassword"),
			"Unknown ZStack API account password",
			"Either target apply the source of the value first, set the value statically in the configuration, or use the ZSTACK_ACCOUNTPASSWORD environment variable.",
		)
	}

	if config.AccessKeyId.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("accesskeyid"),
			"Unknown ZStack  accessKeyId",
			"Either target apply the source of the value first, set the value statically in the configuration, or use the ZSTACK_ACCESSKEYID environment variable.",
		)
	}

	if config.AccessKeySecret.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("accesskeysecret"),
			"Unknown ZStack accessKeySecret",
			"Either target apply the source of the value first, set the value statically in the configuration, or use the ZSTACK_ACCESSKEYSECRET environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	//Defaukt value to environment vairiable, but override
	//with Terraform configuration value if set.

	host := os.Getenv("ZSTACK_HOST")
	accountname := os.Getenv("ZSTACK_ACCOUNTNAME")
	accountpassword := os.Getenv("ZSTACK_ACCOUNTPASSWORD")
	accesskeyid := os.Getenv("ZSTACK_ACCESSKEYID")
	accesskeysecret := os.Getenv("ZSTACK_ACCESSKEYSECRET")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.AccountName.IsNull() {
		accountname = config.AccountName.ValueString()
	}

	if !config.AccountPassword.IsNull() {
		accountpassword = config.AccountPassword.ValueString()
	}

	if !config.AccessKeyId.IsNull() {
		accesskeyid = config.AccessKeyId.ValueString()
	}

	if !config.AccessKeySecret.IsNull() {
		accesskeysecret = config.AccessKeySecret.ValueString()
	}

	// If any of the expected configuration are missing, return
	// errors with provider-sepecific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing ZStack API Host",
			"The provider cannot create the ZStack API client as there is a missing or empty value for the ZStack API host. "+
				"Set the host value in the configuration or use the ZSTACK_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if accountname == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("accountname"),
			"Missing ZStack API Account Username",
			"The provider cannot create the ZStack API client as there is a missing or empty value for the ZStack API Account Username. "+
				"Set the host value in the configuration or use the ZSTACK_ACCOUNTNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if accountpassword == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("accountpassword"),
			"Missing ZStack API Account Password",
			"The provider cannot create the ZStack API client as there is a missing or empty value for the ZStack API Account Password. "+
				"Set the host value in the configuration or use the ZSTACK_ACCOUNPASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}
	if accesskeyid == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("accesskeyid"),
			"Missing ZStack AccesskeyId",
			"The provider cannot create the ZStack API client as there is a missing or empty value for the ZStack API AccessKeyId. "+
				"Set the host value in the configuration or use the ZSTACK_ACCESSKEYID environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}
	if accesskeysecret == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("accesskeysecret"),
			"Missing ZStack AccesskeySecret",
			"The provider cannot create the ZStack API client as there is a missing or empty value for the ZStack API AccessKeySecret. "+
				"Set the host value in the configuration or use the ZSTACK_ACCESSKEYSECRET environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "ZStack_host", host)
	ctx = tflog.SetField(ctx, "ZStack_accountName", accountname)
	ctx = tflog.SetField(ctx, "ZStack_accountPassword", accountpassword)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "ZStack_accountPassword")

	tflog.Debug(ctx, "Creating ZStack client")

	//Create a new ZStack client using the configuration values
	client := client.NewZSClient(client.DefaultZSConfig(host).
		LoginAccount(accountname, accountpassword).ReadOnly(false).Debug(true))

	_, err := client.Login()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create ZStack API Client",
			"An unexpected error occurred when creating the ZStack API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"ZStack Client Error: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured ZStack client", map[string]any{"success": true})
}

// DataSources implements provider.Provider.
func (p *ZStackProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		zstack.ZStackImageDataSource,
		//zstack.ZStackvmsDataSource,
		//zstack.ZStackImgDataSource,
		//zstack.NewExampleDataSource,
	}

}

// Metadata implements provider.Provider.
func (p *ZStackProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "zstack"
	resp.Version = p.version
}

// Resources implements provider.Provider.
func (p *ZStackProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		zstack.ImageResource,
		zstack.ZStackvmResource,
	}

}

// Schema implements provider.Provider.
func (p *ZStackProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "ZStack Cloud MN HOST ip address. May also be provided via ZSTACK_HOST environment variable.",
				Optional:    true,
			},
			"accountname": schema.StringAttribute{
				Description: "Username for ZStack API. May also be provided via ZSTACK_ACCOUNTNAME environment variable.",
				Optional:    true,
			},
			"accountpassword": schema.StringAttribute{
				Description: "Password for ZStack API. May also be provided via ZSTACK_ACCOUNTPASSWORD environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"accesskeyid": schema.StringAttribute{
				Description: "AccessKey ID for ZStack API. Create AccessKey ID from MN,  Operational Management->Access Control->AccessKey Management. May also be provided via ZSTACK_ACCESSKEYID environment variable.",
				Optional:    true,
			},
			"accesskeysecret": schema.StringAttribute{
				Description: "AccessKey Secret for ZStack API. May also be provided via ZSTACK_ACCESSKEYSECRET environment variable.",
				Optional:    true, Sensitive: true,
			},
		},
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ZStackProvider{
			version: version,
		}
	}
}
