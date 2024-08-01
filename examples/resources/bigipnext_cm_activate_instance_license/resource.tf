resource "bigipnext_cm_activate_instance_license" "tokenadd" {
  instances = [{
    instance_address = "10.xxx.xxx.xxx"
    jwt_id           = "8a3dc22e-xxxx-xxxxc-xxxx-xxxxxxxx4326"
    },
    {
      instance_address = "10.146.194.174"
      jwt_id           = "8a3dc22e-xxxx-xxxxc-xxxx-xxxxxxxx4326"
    }
  ]
}