// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"net/http"

	"github.com/gillcaleb/orbit-bhyve-go-client/pkg/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure bhyveProvider satisfies various provider interfaces.
var _ provider.Provider = &bhyveProvider{}
var _ provider.ProviderWithFunctions = &bhyveProvider{}

// bhyveProvider defines the provider implementation.
type bhyveProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// bhyveProviderModel maps provider schema data to a Go type.
type bhyveProviderModel struct {
	DeviceId     types.String `tfsdk:"deviceid"`
	Email        types.String `tfsdk:"email"`
	Password     types.String `tfsdk:"password"`
}


func (p *bhyveProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "bhyve"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *bhyveProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
			Attributes: map[string]schema.Attribute{
					"deviceid": schema.StringAttribute{
							Optional: false,
							Sensitive: true,
					},
					"email": schema.StringAttribute{
							Optional: false,
					},
					"password": schema.StringAttribute{
							Optional:  true,
							Sensitive: true,
					},
			},
	}
}


func (p *bhyveProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config bhyveProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
			return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.DeviceId.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
					path.Root("deviceid"),
					"Unknown DeviceId",
					"The provider cannot create the Bhyve API client as there is an unknown configuration value for bhyve Device ID. "+
							"Either target apply the source of the value first, set the value statically in the configuration, or use the BHYVE_DEVICEID environment variable.",
			)
	}

	if config.Username.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
					path.Root("email"),
					"Unknown Bhyve Email",
					"The provider cannot create the  API client as there is an unknown configuration value for the Bhyve API username. "+
							"Either target apply the source of the value first, set the value statically in the configuration, or use the BHYVE_USERNAME environment variable.",
			)
	}

	if config.Password.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
					path.Root("password"),
					"Unknown Bhyve Password",
					"The provider cannot create the Bhyve API client as there is an unknown configuration value for the Bhyve API password. "+
							"Either target apply the source of the value first, set the value statically in the configuration, or use the BHYVE_PASSWORD environment variable.",
			)
	}

	if resp.Diagnostics.HasError() {
			return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("BHYVE_HOST")
	username := os.Getenv("BHYVE_USERNAME")
	password := os.Getenv("BHYVE_PASSWORD")

	if !config.DeviceId.IsNull() {
			deviceid = config.DeviceId.ValueString()
	}

	if !config.Email.IsNull() {
			email = config.Email.ValueString()
	}

	if !config.Password.IsNull() {
			password = config.Password.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if deviceid == "" {
			resp.Diagnostics.AddAttributeError(
					path.Root("host"),
					"Missing Bhyve API Host",
					"The provider cannot create the Bhyve API client as there is a missing or empty value for the Bhyve API host. "+
							"Set the host value in the configuration or use the BHYVE_DEVICEID environment variable. "+
							"If either is already set, ensure the value is not empty.",
			)
	}

	if email == "" {
			resp.Diagnostics.AddAttributeError(
					path.Root("email"),
					"Missing Bhyve API email",
					"The provider cannot create the Bhyve API client as there is a missing or empty value for the Bhyve API username. "+
							"Set the username value in the configuration or use the BHYVE_EMAIL environment variable. "+
							"If either is already set, ensure the value is not empty.",
			)
	}

	if password == "" {
			resp.Diagnostics.AddAttributeError(
					path.Root("password"),
					"Missing Bhyve API Password",
					"The provider cannot create the Bhyve API client as there is a missing or empty value for the Bhyve API password. "+
							"Set the password value in the configuration or use the BHYVE_PASSWORD environment variable. "+
							"If either is already set, ensure the value is not empty.",
			)
	}

	if resp.Diagnostics.HasError() {
			return
	}
  
	config := client.Config{
		Endpoint: "https://api.orbitbhyve.com/v1",
		Email: email,
		Password: password,
		DeviceId: deviceid,
  }
	// Create a new Bhyve client using the configuration values
	client, err := client.NewClient(config)
	if err != nil {
			resp.Diagnostics.AddError(
					"Unable to Create Bhyve API Client",
					"An unexpected error occurred when creating the Bhyve API client. "+
							"If the error is not clear, please contact the provider developers.\n\n"+
							"Bhyve Client Error: "+err.Error(),
			)
			return
	}

	// Make the Bhyve client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}


func (p *bhyveProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewExampleResource,
	}
}

func (p *bhyveProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewExampleDataSource,
	}
}

func (p *bhyveProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		NewExampleFunction,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &bhyveProvider{
			version: version,
		}
	}
}
