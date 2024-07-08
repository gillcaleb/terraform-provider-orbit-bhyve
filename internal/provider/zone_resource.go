package provider

import (
    "context"
		"fmt"
		"time"
		"strconv"
    
		"github.com/gillcaleb/orbit-bhyve-go-client/pkg/client"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
		"github.com/hashicorp/terraform-plugin-framework/types"
		"github.com/hashicorp/terraform-plugin-log/tflog"
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

type zoneResourceModel struct {
	ID          types.String     `tfsdk:"id"`
	LastUpdated types.String     `tfsdk:"last_updated"`
	Minutes     types.String     `tfsdk:"minutes"`
}

// Metadata returns the resource type name.
func (r *zoneResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_zone"
}

// Configure adds the provider configured client to the resource.
func (r *zoneResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
  // Add a nil check when handling ProviderData because Terraform
  // sets that data after it calls the ConfigureProvider RPC.
  if req.ProviderData == nil {
        return
    }

    client, ok := req.ProviderData.(*client.Client)

    if !ok {
        resp.Diagnostics.AddError(
            "Unexpected Data Source Configure Type",
            fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
        )

        return
    }

    r.client = client
}

// Schema defines the schema for the resource.
func (r *zoneResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
			Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
							Required: true,
					},
					"last_updated": schema.StringAttribute{
							Computed: true,
					},
					"minutes": schema.StringAttribute{
						  Required: true,
					},					
			},
	}
}


// Create creates the resource and sets the initial Terraform state.
func (r *zoneResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan zoneResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
			return
	}

	id, err := strconv.Atoi(plan.ID.ValueString())
  if err != nil {
		resp.Diagnostics.AddError(
			"Error converting ID",
			"Could not convert plan.ID to int: "+err.Error(),
		)
		return
	}

	minutes, err := strconv.Atoi(plan.Minutes.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting minutes",
			"Could not convert plan.Minutes to int: "+err.Error(),
		)
		return
	}

	// Create new zone run
	r.client.Sync()
	r.client.StartZone(id, minutes)
	tflog.Info(ctx, "Checking the status of the StartZone command")

	// Map response body to schema and populate Computed attribute values
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
			return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *zoneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state zoneResourceModel
  diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
			return
	}
  
	idi, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting ID",
			"Could not convert plan.ID to int: "+err.Error(),
		)
		return
	}
	// Get zone information
	id := r.client.ReadZone(idi)

	// Map response body to schema and populate attribute values
	state.Minutes = types.StringValue(id)

	 // Set refreshed state
	 diags = resp.State.Set(ctx, &state)
	 resp.Diagnostics.Append(diags...)
	 if resp.Diagnostics.HasError() {
			 return
	 }
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *zoneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *zoneResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
