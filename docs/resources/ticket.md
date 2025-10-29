---
page_title: "freshservice_ticket Resource - freshservice"
subcategory: ""
description: |-
  Manages a Freshservice ticket
---

# freshservice_ticket (Resource)

Manages a Freshservice ticket. This resource allows you to create and manage support tickets in Freshservice with associated assets, priority levels, and status tracking.

## Example Usage

```terraform
# Create a basic ticket
resource "freshservice_ticket" "example" {
  subject     = "Software Installation Request"
  description = "Need to install Microsoft Office on my laptop"
  priority    = 2
  status      = 2
  email       = "user@company.com"
}

# Create a ticket with workspace, group, responder and assets
resource "freshservice_ticket" "with_assignment" {
  subject      = "Hardware Issue - Multiple Laptops"
  description  = "Multiple laptops in the office are experiencing hardware issues"
  priority     = 3
  status       = 2
  email        = "it-support@company.com"
  workspace_id = 1
  group_id     = 5
  responder_id = 1001

  assets {
    display_id = 8
  }
  
  assets {
    display_id = 9
  }
}

# Create a high priority ticket
resource "freshservice_ticket" "urgent" {
  subject     = "Server Down - Production Environment"
  description = "The main production server is not responding"
  priority    = 4
  status      = 2
  email       = "admin@company.com"
  
  assets {
    display_id = 15
  }
}
```

## Schema

### Required

- `subject` (String) The subject of the ticket.
- `description` (String) The description of the ticket.
- `priority` (Number) The priority of the ticket. Valid values: 1 (Low), 2 (Medium), 3 (High), 4 (Urgent).
- `status` (Number) The status of the ticket. Valid values: 2 (Open), 3 (Pending), 4 (Resolved), 5 (Closed).

### Optional

- `email` (String) The email of the ticket requester.
- `workspace_id` (Number) The workspace ID associated with the ticket.
- `group_id` (Number) ID of the group to which the ticket has been assigned.
- `responder_id` (Number) ID of the agent to whom the ticket has been assigned.
- `assets` (Block List) The assets associated with the ticket. (see [below for nested schema](#nestedblock--assets))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--assets"></a>
### Nested Schema for `assets`

Required:

- `display_id` (Number) The display ID of the asset.

## Priority Levels

- `1` - Low
- `2` - Medium  
- `3` - High
- `4` - Urgent

## Status Values

- `2` - Open
- `3` - Pending
- `4` - Resolved
- `5` - Closed

## Import

Import is supported using the following syntax:

```shell
terraform import freshservice_ticket.example 12345
```

Where `12345` is the ticket ID from Freshservice.