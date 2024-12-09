# example to Configure NEXT instances (with DNS Servers, NTP Servers)
resource "bigipnext_cm_instance_onboard" "test" {
  dns_servers        = ["2.2.2.4"]
  ntp_servers        = ["4.pool.com"]
  management_address = "10.218.135.67"
  timeout            = 300
}

# example to Configure NEXT instances (with DNS Servers, NTP Servers, L1 Networks)
resource "bigipnext_cm_instance_onboard" "test" {
  dns_servers        = ["2.2.2.4"]
  ntp_servers        = ["4.pool.com"]
  management_address = "10.218.135.67"
  l1_networks = [{
    name = "l1network4"
    l1_link = {
      name = "1.1"
      link_type : "Interface"
    }
  }]
  timeout = 300
}

# example to Configure NEXT instances (with DNS Servers, NTP Servers, L1 Networks, VLANs)
resource "bigipnext_cm_instance_onboard" "test" {
  dns_servers        = ["2.2.2.4"]
  ntp_servers        = ["4.pool.com"]
  management_address = "10.218.135.67"
  l1_networks = [{
    name = "l1network4"
    vlans = [
      {
        tag  = 104
        name = "vlan104"
      }
    ]
    l1_link = {
      name = "1.1"
      link_type : "Interface"
    }
  }]
  timeout = 300
}

# example to Configure NEXT instances (with DNS Servers, NTP Servers, L1 Networks, VLANs, SelfIPs)
resource "bigipnext_cm_instance_onboard" "test" {
  dns_servers        = ["2.2.2.5"]
  ntp_servers        = ["5.pool.com"]
  management_address = "10.218.135.67"
  l1_networks = [{
    name = "l1network4"
    vlans = [
      {
        tag  = 104
        name = "vlan104"
        self_ips = [
          {
            address     = "20.20.20.24/24"
            device_name = "device4"
          }
        ]
      }
    ]
    l1_link = {
      name = "1.1"
      link_type : "Interface"
    }
  }]
  timeout = 300
}

# example to Configure/Update NEXT instances (with multiple L1 Networks)
resource "bigipnext_cm_instance_onboard" "sample123" {
  dns_servers        = ["2.2.2.6"]
  ntp_servers        = ["5.pool.com"]
  management_address = "10.218.135.67"
  l1_networks = [
    {
      name = "l1network-internal"
      vlans = [
        {
          tag  = 105
          name = "vlan105"
          self_ips = [
            {
              address     = "20.20.20.25/24"
              device_name = "device5"
            }
          ]
        }
      ]
      l1_link = {
        name      = "1.1"
        link_type = "Interface"
      }
    },
    {
      name = "l1network-external"
      vlans = [
        {
          tag  = 106
          name = "vlan106"
          self_ips = [
            {
              address     = "20.20.20.26/24"
              device_name = "device6"
            }
          ]
        }
      ]
      l1_link = {
        name      = "1.2"
        link_type = "Interface"
      }
    }
  ]
  timeout = 300
}