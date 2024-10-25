# Upgrade BIG-IP Next VE or appliance to a new version
resource "bigipnext_cm_next_upgrade" "upgrade_ve" {
  upgrade_type       = "ve"
  next_instance_ip   = "1.2.3.4"
  image_name         = "BIG-IP-Next-20.3.0-2.713.1.tgz"
  signature_filename = "BIG-IP-Next-20.3.0-2.713.1.tgz.512.sig"
  timeout            = 3600
}

# Upgrade BIG-IP Next appliance to a new version

resource "bigipnext_cm_next_upgrade" "upgrade_appliance" {
  upgrade_type       = "appliance"
  tenant_name        = "testtenantx"
  next_instance_ip   = "12.34.56.78"
  image_name         = "BIG-IP-Next-20.3.1-2.538.0"
  partition_address  = "10.217.30.40"
  partition_port     = 8888
  partition_username = "admin"
  partition_password = "admin"
}