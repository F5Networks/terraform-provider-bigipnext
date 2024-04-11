resource "bigipnext_cm_next_ha" "test" {
  ha_name                       = "testnextha"
  ha_ip                         = "10.xxx.xxx.xx"
  active_node_ip                = bigipnext_cm_deploy_vmware.vmware01.instance.mgmt_address
  standby_node_ip               = bigipnext_cm_deploy_vmware.vmware02.instance.mgmt_address
  control_plane_vlan            = "ha-cp-vlan"
  control_plane_vlan_tag        = 101
  data_plane_vlan               = "ha-dp-vlan"
  data_plane_vlan_tag           = 102
  active_node_control_plane_ip  = "10.xxx.xxx.xxx/xx"
  standby_node_control_plane_ip = "10.xxx.xxx.xx/xx"
  active_node_data_plane_ip     = "10.xxx.xxx.xx/xx"
  standby_node_data_plane_ip    = "10.xxx.xxx.xxx/xx"
}