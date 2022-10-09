consul {
  address = "consul:8500"

  retry {
    enabled  = true
    attempts = 12
    backoff  = "250ms"
  }
}
template {
  source      = "/scripts/nginx.conf.tmpl"
  destination = "/etc/nginx/conf.d/default.conf"
  perms       = 0600
  command     = "service nginx restart"
}
