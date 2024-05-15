resource "bigipnext_cm_discover_next" "test" {
  address             = "10.10.10.10"
  port                = 5443
  device_user         = "admin"
  device_password     = "admin123"
  management_user     = "admin-cm"
  management_password = "admin@123"
}