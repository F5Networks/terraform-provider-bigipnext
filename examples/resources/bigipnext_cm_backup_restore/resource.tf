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