resource "bigipnext_cm_deploy_vmware" "vmware" {
  vsphere_provider = {
    provider_name      = "myvsphere01"
    content_library    = "Contentlibrary"
    cluster_name       = "vSAN Cluster"
    datacenter_name    = "vpshere-7.0"
    datastore_name     = "xxxxxxx"
    resource_pool_name = "xxxxxx"
    vm_template_name   = "BIG-IP-Next-20.0.1-2.139.10-0.0.136-VM-template"
  }
  instance = {
    instance_hostname     = "testecosyshydvm06"
    mgmt_address          = "10.xxx.xxx.xxxx"
    mgmt_prefix           = 24
    mgmt_gateway          = "10.xxx.xxx.xxx"
    mgmt_network_name     = "VM-mgmt"
    mgmt_user             = "admintest"
    mgmt_password         = "F5Test@123"
    external_network_name = "LocalTestVLAN-115"
  }
  ntp_servers = ["0.us.pool.ntp.org"]
  dns_servers = ["8.8.8.8"]
  timeout     = 1200
}