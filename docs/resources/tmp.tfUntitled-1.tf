resource "ecloud_firewallrule" "firewallrule-web" {
  firewall_policy_id = "fwp-abcdef12"
  sequence           = 0
  name               = "firewallrule-web"
  direction          = "IN"
  source             = "ANY"
  destination        = "ANY"
  action             = "ALLOW"
  enabled            = true

  port {
    protocol    = "TCP"
    source      = "ANY"
    destination = "80"
  }

  port {
    protocol    = "TCP"
    source      = "ANY"
    destination = "443"
  }
}
