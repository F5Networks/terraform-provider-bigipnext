resource "bigipnext_cm_deploy_f5os" "rseries01" {
  f5os_provider = {
    provider_name = "myrseries"
    provider_type = "rseries"
  }
  instance = {
    instance_hostname      = "rseriesravitest04"
    management_address     = "10.xxx.xxx.xxx"
    management_prefix      = 24
    management_gateway     = "10.xxx.xxx.xxx"
    management_user        = "admin-cm"
    management_password    = "F5test@123"
    vlan_ids               = [27, 28, 29]
    tenant_deployment_file = "BIG-IP-Next-20.1.0-2.279.0+0.0.18.yaml"
    tenant_image_name      = "BIG-IP-Next-20.1.0-2.279.0+0.0.18"
  }
}
