resource "bigipnext_cm_bootstrap" "name" {
  run_setup         = true
  bootstrap_timeout = 800
  external_storage = {
    storage_type    = "NFS"
    storage_address = "10.28.14.22"
    storage_path    = "/exports/backup"
    cm_storage_dir  = "backuppqr"
  }
}