resource "bigipnext_cm_device_backup_restore" "sample" {
  backup_password = "F5site02@123"
  operation       = "backup"
  file_name       = "test.tar.gz"
  device_hostname = "big-ip-next"
}