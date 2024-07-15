// Copyright (c) ZStack.io, Inc.
// SPDX-License-Identifier: MPL-2.0

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
	AccountName     types.String `tfsdk:"accountname"`
	AccountPassword types.String `tfsdk:"accountpassword"`
	AccessKeyId     types.String `tfsdk:"accesskeyid"`
	AccessKeySecret types.String `tfsdk:"accesskeysecret"`
	//SessionId       types.String `tfsdk:"sessionid"`
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

	port := 8080
	//sessionId := ""
	host := os.Getenv("ZSTACK_HOST")
	portstr := os.Getenv("ZSTACK_PORT")
	accountname := os.Getenv("ZSTACK_ACCOUNTNAME")
	accountpassword := os.Getenv("ZSTACK_ACCOUNTPASSWORD")
	accesskeyid := os.Getenv("ZSTACK_ACCESSKEYID")
	accesskeysecret := os.Getenv("ZSTACK_ACCESSKEYSECRET")

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
	/*
		if !config.SessionId.IsNull() {
			sessionId = config.SessionId.ValueString()
		}
	*/

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

	if (accountname == "" || accountpassword == "") && (accesskeyid == "" || accesskeysecret == "") {
		resp.Diagnostics.AddError(
			"Missing ZStack Authorization",
			"The provider cannot create the ZStack API client as there is no zstack authorization. \n"+
				"Please set at least one authorization method: account_name + account_password OR access_key_id + access_key_secret.\n\n"+
				"account_name value can be set in the configuration or use the ZSTACK_ACCOUNTNAME environment variable\n"+
				"account_password value in the configuration or use the ZSTACK_ACCOUNTPASSWORD environment variable\n"+
				"access_key_id value in the configuration or use the ZSTACK_ACCESSKEYID environment variable\n"+
				"access_key_secret value in the configuration or use the ZSTACK_ACCESSKEYSECRET environment variable\n",
		)
	}
	/*
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
	*/
	if resp.Diagnostics.HasError() {
		return
	}

	var cli *client.ZSClient

	ctx = tflog.SetField(ctx, "ZStack_host", host)
	//ctx = tflog.SetField(ctx, "ZStack_port", port)

	//zsConfig := client.NewZSConfig(host, port, "zstack").RetryTimes(450).ReadOnly(false).Debug(true)
	//Create a new ZStack client using the configuration values

	/* if sessionId != "" {
	ctx = tflog.SetField(ctx, "ZStack_sesssionId", sessionId)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "ZStack_sessionId")

	tflog.Debug(ctx, "Creating ZStack client with session id")
	cli = client.NewZSClient(client.NewZSConfig(host, port, "zstack").Session(sessionId).ReadOnly(false).Debug(true))
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
	} else */
	if accountname != "" && accountpassword != "" {
		ctx = tflog.SetField(ctx, "ZStack_accountName", accountname)
		ctx = tflog.SetField(ctx, "ZStack_accountPassword", accountpassword)
		ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "ZStack_accountPassword")

		tflog.Debug(ctx, "Creating ZStack client with account")
		cli = client.NewZSClient(client.NewZSConfig(host, port, "zstack").LoginAccount(accountname, accountpassword).ReadOnly(false).Debug(true))
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
	} else if accesskeyid != "" && accesskeysecret != "" {
		ctx = tflog.SetField(ctx, "ZStack_accessKeyId", accesskeyid)
		ctx = tflog.SetField(ctx, "ZStack_accessKeySecret", accesskeysecret)
		ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "ZStack_accessKeySecret")

		tflog.Debug(ctx, "Creating ZStack client with access key")
		cli = client.NewZSClient(client.NewZSConfig(host, port, "zstack").AccessKey(accesskeyid, accesskeysecret).ReadOnly(false).Debug(true))
		// no authorization validation! this access key may be invalid！
	}
	resp.DataSourceData = cli
	resp.ResourceData = cli

	tflog.Info(ctx, "Configured ZStack client", map[string]any{"success": true})
}

// DataSources implements provider.Provider.
func (p *ZStackProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		ZStackImageDataSource,
		ZStackl3NetworkDataSource,
		ZStackvmsDataSource,
		NewClusterDataSource,
		ZstackBackupStorageDataSource,
		ZStackHostsDataSource,
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
		ZStackvmResource,
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
			"port": schema.Int64Attribute{
				Description: "Stack Cloud MN API port. May also be provided via ZSTACK_PORT environment variable.",
				Optional:    true,
			},
			/*
				"sessionid": schema.StringAttribute{
					Description: "Stack Cloud MN sessionid. May also be provided via ZSTACK_SESSIONID environment variable.",
					Optional:    true,
				},*/
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
