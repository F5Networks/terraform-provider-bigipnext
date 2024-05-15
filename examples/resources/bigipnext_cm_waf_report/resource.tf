resource "bigipnext_cm_waf_report" "sample" {
  name               = "sample"
  description        = "WAF security report description"
  request_type       = "illegal"
  time_frame_in_days = 7
  top_level          = 5
  categories         = []
  scope = {
    entity = "applications",
    all    = false
    names  = ["test_policy", "test"]
  }
}