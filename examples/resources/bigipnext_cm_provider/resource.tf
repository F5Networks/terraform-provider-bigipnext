resource "bigipnext_cm_provider" "vpshere" {
  name     = "testvpshereprovvm01"
  address  = "xxx70-vcenter.f5.com"
  type     = "VSPHERE"
  username = "xxxxxxx"
  password = "xxxxxxx"
}
resource "bigipnext_cm_provider" "rseries" {
  name     = "testrseriesprovvm01"
  address  = "10.1.1.80:443"
  type     = "RSERIES"
  username = "xxxxxxx"
  password = "xxxxxxxxx"
}