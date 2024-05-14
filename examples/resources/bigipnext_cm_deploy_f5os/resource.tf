# example to deploy Next Instance in rSeries
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

# example to deploy Next Instance in Velos
resource "bigipnext_cm_deploy_f5os" "velos01" {
  f5os_provider = {
    provider_name = "myvelos01"
    provider_type = "velos"
  }
  instance = {
    instance_hostname      = "ravitestvelos01"
    management_address     = "10.xx.xxx.xxx"
    management_prefix      = 24
    management_gateway     = "10.xxx.xxx.xx"
    management_user        = "admin-cm"
    management_password    = "site02@123"
    vlan_ids               = [100]
    slot_ids               = [1]
    tenant_deployment_file = "BIG-IP-Next-20.1.0-2.375.1+0.0.43.yaml"
    tenant_image_name      = "BIG-IP-Next-20.1.0-2.375.1+0.0.43"

  }
}
