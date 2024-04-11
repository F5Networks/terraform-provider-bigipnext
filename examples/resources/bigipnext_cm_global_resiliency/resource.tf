resource "bigipnext_cm_global_resiliency" "sample4" {
  name              = "sample4"
  dns_listener_name = "dln"
  dns_listener_port = 10
  protocols         = ["tcp"]
  instances = [
    {
      address              = "10.145.71.115"
      dns_listener_address = "2.2.2.3"
      group_sync_address   = "10.10.1.2/24"
      hostname             = "big-ip-next"
    }
  ]
}