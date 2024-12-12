// Copyright (c) ZStack.io, Inc.

/*
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package provider

import (
	"context"
	"os"
	"strconv"

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
	Port            types.Int64  `tfsdk:"port"`
	AccountName     types.String `tfsdk:"account_name"`
	AccountPassword types.String `tfsdk:"account_password"`
	AccessKeyId     types.String `tfsdk:"access_key_id"`
	AccessKeySecret types.String `tfsdk:"access_key_secret"`
	SessionId       types.String `tfsdk:"session_id"`
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

	if config.Port.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("port"),
			"Unknown ZStack Cloud API Port",
			"The provider cannt create the ZStack Cloud API client as an unknown configuration value for the ZStack Cloud API port."+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the ZSTACK_PORT environment variable.",
		)
	}

	if config.AccountName.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("account_name"),
			"Unknown ZStack API account Username",
			"Either target apply the source of the value first, set the value statically in the configuration, or use the ZSTACK_ACCOUNT_NAME environment variable.",
		)
	}

	if config.AccountPassword.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("account_password"),
			"Unknown ZStack API account password",
			"Either target apply the source of the value first, set the value statically in the configuration, or use the ZSTACK_ACCOUNT_PASSWORD environment variable.",
		)
	}

	if config.AccessKeyId.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("access_key_id"),
			"Unknown ZStack  access_key_id",
			"Either target apply the source of the value first, set the value statically in the configuration, or use the ZSTACK_ACCESS_KEY_ID environment variable.",
		)
	}

	if config.AccessKeySecret.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("access_key_secret"),
			"Unknown ZStack accessKeySecret",
			"Either target apply the source of the value first, set the value statically in the configuration, or use the ZSTACK_ACCESS_KEY_SECRET environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	//Defaukt value to environment vairiable, but override
	//with Terraform configuration value if set.

	port := 8080
	sessionId := ""

	host := os.Getenv("ZSTACK_HOST")
	portstr := os.Getenv("ZSTACK_PORT")
	account_name := os.Getenv("ZSTACK_ACCOUNT_NAME")
	account_password := os.Getenv("ZSTACK_ACCOUNT_PASSWORD")
	access_key_id := os.Getenv("ZSTACK_ACCESS_KEY_ID")
	access_key_secret := os.Getenv("ZSTACK_ACCESS_KEY_SECRET")

	if portstr != "" {
		if portInt, err := strconv.Atoi(portstr); err == nil {
			port = portInt
		}
	}

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Port.IsNull() {
		port = int(config.Port.ValueInt64())
	}

	if !config.AccountName.IsNull() {
		account_name = config.AccountName.ValueString()
	}

	if !config.AccountPassword.IsNull() {
		account_password = config.AccountPassword.ValueString()
	}

	if !config.AccessKeyId.IsNull() {
		access_key_id = config.AccessKeyId.ValueString()
	}

	if !config.AccessKeySecret.IsNull() {
		access_key_secret = config.AccessKeySecret.ValueString()
	}

	if !config.SessionId.IsNull() {
		sessionId = config.SessionId.ValueString()
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

	//session id is just used for marketplace-server. Don't expose to user!
	if sessionId == "" && (account_name == "" || account_password == "") && (access_key_id == "" || access_key_secret == "") {
		resp.Diagnostics.AddError(
			"Missing ZStack Authorization",
			"The provider cannot create the ZStack API client as there is no zstack authorization. \n"+
				"Please set at least one authorization method: account_name + account_password OR access_key_id + access_key_secret.\n\n"+
				"account_name value can be set in the configuration or use the ZSTACK_ACCOUNT_NAME environment variable\n"+
				"account_password value in the configuration or use the ZSTACK_ACCOUNT_PASSWORD environment variable\n"+
				"access_key_id value in the configuration or use the ZSTACK_ACCESS_KEY_ID environment variable\n"+
				"access_key_secret value in the configuration or use the ZSTACK_ACCESS_KEY_SECRET environment variable\n")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	var cli *client.ZSClient

	ctx = tflog.SetField(ctx, "ZStack_host", host)
	ctx = tflog.SetField(ctx, "ZStack_port", port)

	zsConfig := client.NewZSConfig(host, port, "zstack").RetryTimes(900).ReadOnly(false).Debug(true)
	//Create a new ZStack client using the configuration values

	if sessionId != "" {
		ctx = tflog.SetField(ctx, "ZStack_sessionId", sessionId)
		ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "ZStack_sessionId")

		tflog.Debug(ctx, "Creating ZStack client with session id")
		cli = client.NewZSClient(zsConfig.Session(sessionId))
		_, err := cli.ValidateSession()
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Create ZStack API Client",
				"An unexpected error occurred when creating the ZStack API client."+
					"It might be due to an incorrect session id being set, or the session has expired.\""+
					"If the error is not clear, please contact the provider developers.\n\n"+
					"ZStack Client Error: "+err.Error(),
			)
			return
		}
	} else if account_name != "" && account_password != "" {
		ctx = tflog.SetField(ctx, "ZStack_accountName", account_name)
		ctx = tflog.SetField(ctx, "ZStack_accountPassword", account_password)
		ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "ZStack_accountPassword")

		tflog.Debug(ctx, "Creating ZStack client with account")
		cli = client.NewZSClient(client.NewZSConfig(host, port, "zstack").LoginAccount(account_name, account_password).ReadOnly(false).Debug(true))
		_, err := cli.Login()
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Create ZStack API Client",
				"An unexpected error occurred when creating the ZStack API client. "+
					"It might be due to an incorrect account name and password being set"+
					"If the error is not clear, please contact the provider developers.\n\n"+
					"ZStack Client Error: "+err.Error(),
			)
			return
		}
	} else if access_key_id != "" && access_key_secret != "" {
		ctx = tflog.SetField(ctx, "ZStack_accessKeyId", access_key_id)
		ctx = tflog.SetField(ctx, "ZStack_accessKeySecret", access_key_secret)
		ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "ZStack_accessKeySecret")

		tflog.Debug(ctx, "Creating ZStack client with access key")
		cli = client.NewZSClient(client.NewZSConfig(host, port, "zstack").AccessKey(access_key_id, access_key_secret).ReadOnly(false).Debug(true))
		// no authorization validation! this access key may be invalid！
	}
	resp.DataSourceData = cli
	resp.ResourceData = cli

	tflog.Info(ctx, "Configured ZStack client", map[string]any{"success": true})
}

/*
	if account_name != "" && account_password != "" {
		ctx = tflog.SetField(ctx, "ZStack_accountName", account_name)
		ctx = tflog.SetField(ctx, "ZStack_accountPassword", account_password)
		ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "ZStack_accountPassword")

		tflog.Debug(ctx, "Creating ZStack client with account")
		cli = client.NewZSClient(client.NewZSConfig(host, port, "zstack").LoginAccount(account_name, account_password).ReadOnly(false).Debug(true))
		_, err := cli.Login()
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Create ZStack API Client",
				"An unexpected error occurred when creating the ZStack API client. "+
					"It might be due to an incorrect account name and password being set"+
					"If the error is not clear, please contact the provider developers.\n\n"+
					"ZStack Client Error: "+err.Error(),
			)
			return
		}
	} else if access_key_id != "" && access_key_secret != "" {
		ctx = tflog.SetField(ctx, "ZStack_accessKeyId", access_key_id)
		ctx = tflog.SetField(ctx, "ZStack_accessKeySecret", access_key_secret)
		ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "ZStack_accessKeySecret")

		tflog.Debug(ctx, "Creating ZStack client with access key")
		cli = client.NewZSClient(client.NewZSConfig(host, port, "zstack").AccessKey(access_key_id, access_key_secret).ReadOnly(false).Debug(true))
		// no authorization validation! this access key may be invalid！
	}
	resp.DataSourceData = cli
	resp.ResourceData = cli

	tflog.Info(ctx, "Configured ZStack client", map[string]any{"success": true})

}
*/

// DataSources implements provider.Provider.
func (p *ZStackProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		ZStackZoneDataSource,
		ZStackImageDataSource,
		ZStackl3NetworkDataSource,
		ZStackvmsDataSource,
		ZStackClusterDataSource,
		ZStackBackupStorageDataSource,
		ZStackHostsDataSource,
		ZStackmnNodeDataSource,
		ZStackl2NetworkDataSource,
		ZStackvrouterDataSource,
		ZStackVirtualRouterImageDataSource,
		ZStackVRouterOfferingDataSource,
		ZStackVIPsDataSource,
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
		ImageResource,
		InstanceResource,
		ReservedIpResource,
		SubnetResource,
		VpcResource,
		VipResource,
		VirtualRouterImageResource,
		VirtualRouterOfferingResource,
		VirtualRouterInstanceResource,
	}
}

// Schema implements provider.Provider.
func (p *ZStackProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "ZStack Cloud MN HOST ip address. May also be provided via ZSTACK_HOST environment variable.",
				Required:    true,
			},
			"port": schema.Int64Attribute{
				Description: "ZStack Cloud MN API port. May also be provided via ZSTACK_PORT environment variable.",
				Optional:    true,
			},
			"session_id": schema.StringAttribute{
				Description: "ZStack Cloud Session id.",
				Optional:    true,
			},
			"account_name": schema.StringAttribute{
				Description: "Username for ZStack API. May also be provided via ZSTACK_ACCOUN_TNAME environment variable.",
				Optional:    true,
			},
			"account_password": schema.StringAttribute{
				Description: "Password for ZStack API. May also be provided via ZSTACK_ACCOUNT_PASSWORD environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"access_key_id": schema.StringAttribute{
				Description: "AccessKey ID for ZStack API. Create AccessKey ID from MN,  Operational Management->Access Control->AccessKey Management. May also be provided via ZSTACK_ACCESS_KEY_ID environment variable.",
				Optional:    true,
			},
			"access_key_secret": schema.StringAttribute{
				Description: "AccessKey Secret for ZStack API. May also be provided via ZSTACK_ACCESS_KEY_SECRET environment variable.",
				Optional:    true,
				Sensitive:   true,
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
