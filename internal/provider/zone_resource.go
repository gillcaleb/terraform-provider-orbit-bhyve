package provider

import (
    "context"
    
		"github.com/gillcaleb/orbit-bhyve-go-client/pkg/client"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
    _ resource.Resource              = &zoneResource{}
		_ resource.ResourceWithConfigure = &zoneResource{}
)

// NewZoneResource is a helper function to simplify the provider implementation.
func NewZoneResource() resource.Resource {
    return &zoneResource{}
}

// zoneResource is the resource implementation.
type zoneResource struct{
	client *client.Client
}

// Metadata returns the resource type name.
func (r *zoneResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_zone"
}

// Schema defines the schema for the resource.
func (r *zoneResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{}
}

// Create creates the resource and sets the initial Terraform state.
func (r *zoneResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
}

// Read refreshes the Terraform state with the latest data.
func (r *zoneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *zoneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *zoneResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
