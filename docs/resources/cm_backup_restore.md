---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "bigipnext_cm_backup_restore Resource - terraform-provider-bigipnext"
subcategory: ""
description: |-
  Resource used to manage Backup/Restore resources onto BIG-IP Next CM.
---

# bigipnext_cm_backup_restore (Resource)

Resource used to manage Backup/Restore resources onto BIG-IP Next CM.

## Example Usage

```terraform
# Resource to do CM Backup
resource "bigipnext_cm_backup_restore" "test" {
  name                = "Test"
  encryption_password = "F5site02@1234"
  backup              = true
}
# Resource to do CM Restore
resource "bigipnext_cm_backup_restore" "test" {
  name                = "Backup-20241007-075720_L_20.3.1-0.5.0_8.tgz"
  encryption_password = "F5site02@123"
  backup              = false
}

# Resource to do CM Backup with Schedule weekly
resource "bigipnext_cm_backup_restore" "test" {
  name                = "Test"
  encryption_password = "F5site02@1234"
  backup              = true
  schedule = {
    start_at = "2025-08-25T18:30:00Z"
    end_at   = "2025-09-25T18:30:00Z"
  }
  frequency               = "Weekly"
  days_of_the_week_to_run = [1, 2, 3, 4]
}

# Resource to do CM Backup with Schedule Monthly
resource "bigipnext_cm_backup_restore" "test" {
  name                = "Test"
  encryption_password = "F5site02@1234"
  backup              = true
  schedule = {
    start_at = "2025-08-25T18:30:00Z"
    end_at   = "2025-09-25T18:30:00Z"
  }
  frequency               = "Monthly"
  day_of_the_month_to_run = 10
}

# Resource to do CM Backup with Schedule Daily
resource "bigipnext_cm_backup_restore" "test" {
  name                = "Test"
  encryption_password = "F5site02@1234"
  backup              = true
  schedule = {
    start_at = "2025-08-25T18:30:00Z"
  }
  frequency = "Daily"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `backup` (Boolean) Specifies whether backup (if True) or restore (if false) is to be done on CM
- `encryption_password` (String) Encryption password for the backup to be created. Password should be minimum of 8 characters
- `name` (String) The unique name of the backup file. Actual File Name is auto-generated in the case of Instant Backup

### Optional

- `day_of_the_month_to_run` (Number) Specifies From which Day of the month backup should start.
- `days_of_the_week_to_run` (List of Number) Specifies Day of the week on backup has been scheduled. 0-Sunday, 1-Monday and so on
- `frequency` (String) Specifies what is the frequency. Example : Daily, Monthly, Weekly
- `schedule` (Attributes) Specifies whether backup is to be scheduled or not. (see [below for nested schema](#nestedatt--schedule))

### Read-Only

- `file_name` (String) Name of the backup file generate in the case of instant backup
- `id` (String) Unique Identifier for the resource
- `scheduled` (Boolean) Specifies whether backup is scheduled or not
- `type` (String) Type of the Backup

<a id="nestedatt--schedule"></a>
### Nested Schema for `schedule`

Required:

- `start_at` (String) Specifies Start time of the backup. Example: 2019-08-24T14:15:22Z

Optional:

- `end_at` (String) Specifies End time of the backup. Example: 2019-08-24T14:15:22Z
