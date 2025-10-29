package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTicket() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTicketCreate,
		ReadContext:   resourceTicketRead,
		UpdateContext: resourceTicketUpdate,
		DeleteContext: resourceTicketDelete,

        Schema: map[string]*schema.Schema{
            "subject": {
                Type:        schema.TypeString,
                Required:    true,
                Description: "The subject of the ticket.",
            },
            "description": {
                Type:        schema.TypeString,
                Required:    true,
                Description: "The description of the ticket.",
            },
            "priority": {
                Type:        schema.TypeInt,
                Required:    true,
                Description: "The priority of the ticket.",
            },
            "status": {
                Type:        schema.TypeInt,
                Required:    true,
                Description: "The status of the ticket.",
            },
            "email": {
                Type:        schema.TypeString,
                Optional:    true,
                Description: "The email of the ticket requester.",
            },
            "workspace_id": {
                Type:        schema.TypeInt,
                Optional:    true,
                Description: "The workspace ID associated with the ticket.",
            },
            "group_id": {
                Type:        schema.TypeInt,
                Optional:    true,
                Description: "ID of the group to which the ticket has been assigned.",
            },
            "responder_id": {
                Type:        schema.TypeInt,
                Optional:    true,
                Description: "ID of the agent to whom the ticket has been assigned.",
            },
            "assets": {
                Type:        schema.TypeList,
                Optional:    true,
                Description: "The assets associated with the ticket.",
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "display_id": {
                            Type:        schema.TypeInt,
                            Required:    true,
                            Description: "The display ID of the asset.",
                        },
                    },
                },
            },
		},
	}
}

// TicketAsset represents an asset associated with a ticket
type TicketAsset struct {
	DisplayID int `json:"display_id"`
}

// Ticket represents a Freshservice ticket
type Ticket struct {
	ID          int           `json:"id,omitempty"`
	Subject     string        `json:"subject"`
	Description string        `json:"description"`
	Priority    int           `json:"priority"`
	Status      int           `json:"status"`
	Email       string        `json:"email,omitempty"`
	WorkspaceID int           `json:"workspace_id,omitempty"`
	GroupID     int           `json:"group_id,omitempty"`
	ResponderID int           `json:"responder_id,omitempty"`
	Assets      []TicketAsset `json:"assets,omitempty"`
}

// TicketResponse represents the API response when creating/reading a ticket
type TicketResponse struct {
	Ticket Ticket `json:"ticket"`
}

// expandAssets converts Terraform list to TicketAsset structs
func expandAssets(assets []interface{}) []TicketAsset {
	if len(assets) == 0 {
		return nil
	}

	result := make([]TicketAsset, len(assets))
	for i, asset := range assets {
		assetMap := asset.(map[string]interface{})
		result[i] = TicketAsset{
			DisplayID: assetMap["display_id"].(int),
		}
	}
	return result
}

// flattenAssets converts TicketAsset structs to Terraform list
func flattenAssets(assets []TicketAsset) []interface{} {
	if len(assets) == 0 {
		return nil
	}

	result := make([]interface{}, len(assets))
	for i, asset := range assets {
		result[i] = map[string]interface{}{
			"display_id": asset.DisplayID,
		}
	}
	return result
}

func resourceTicketCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config)
	var diags diag.Diagnostics

	// Create the ticket
	ticket := Ticket{
		Subject:     d.Get("subject").(string),
		Description: d.Get("description").(string),
		Priority:    d.Get("priority").(int),
		Status:      d.Get("status").(int),
	}

	// Add optional fields if they exist
	if email, ok := d.GetOk("email"); ok {
		ticket.Email = email.(string)
	}
	if workspaceID, ok := d.GetOk("workspace_id"); ok {
		ticket.WorkspaceID = workspaceID.(int)
	}
	if groupID, ok := d.GetOk("group_id"); ok {
		ticket.GroupID = groupID.(int)
	}
	if responderID, ok := d.GetOk("responder_id"); ok {
		ticket.ResponderID = responderID.(int)
	}
	if assets, ok := d.GetOk("assets"); ok {
		ticket.Assets = expandAssets(assets.([]interface{}))
	}

	// Make the API call to create the ticket
	resp, err := client.CreateTicket(ctx, ticket)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set the ID of the new ticket
	d.SetId(strconv.Itoa(resp.ID))

	return diags
}

func resourceTicketRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config)
	var diags diag.Diagnostics

	ticketID := d.Id()

	// Make the API call to get the ticket
	resp, err := client.GetTicket(ctx, ticketID)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set the ticket data in the resource data
	d.Set("subject", resp.Subject)
	d.Set("description", resp.Description)
	d.Set("priority", resp.Priority)
	d.Set("status", resp.Status)
	d.Set("email", resp.Email)
	d.Set("workspace_id", resp.WorkspaceID)
	d.Set("group_id", resp.GroupID)
	d.Set("responder_id", resp.ResponderID)
	d.Set("assets", flattenAssets(resp.Assets))

	return diags
}

func resourceTicketUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config)

	ticketID := d.Id()

	// Create the updated ticket
	ticket := Ticket{
		Subject:     d.Get("subject").(string),
		Description: d.Get("description").(string),
		Priority:    d.Get("priority").(int),
		Status:      d.Get("status").(int),
	}

	// Add optional fields if they exist
	if email, ok := d.GetOk("email"); ok {
		ticket.Email = email.(string)
	}
	if workspaceID, ok := d.GetOk("workspace_id"); ok {
		ticket.WorkspaceID = workspaceID.(int)
	}
	if groupID, ok := d.GetOk("group_id"); ok {
		ticket.GroupID = groupID.(int)
	}
	if responderID, ok := d.GetOk("responder_id"); ok {
		ticket.ResponderID = responderID.(int)
	}
	if assets, ok := d.GetOk("assets"); ok {
		ticket.Assets = expandAssets(assets.([]interface{}))
	}

	// Make the API call to update the ticket
	_, err := client.UpdateTicket(ctx, ticketID, ticket)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceTicketRead(ctx, d, m)
}

func resourceTicketDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Config)
	var diags diag.Diagnostics

	ticketID := d.Id()

	// Make the API call to delete the ticket
	err := client.DeleteTicket(ctx, ticketID)
	if err != nil {
		return diag.FromErr(err)
	}

	// Remove the resource from the state
	d.SetId("")

	return diags
}