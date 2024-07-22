disable_mlock = true
ui = true
skip_setcap = true

storage "postgresql" {
  ha_enabled = "true"           
}

listener "tcp" {
  address = "[::]:8200"
  cluster_address = "[::]:8201"
  tls_disable = 1
}

# Advertise the non-loopback interface
api_addr = "http://vault.vault:8200"
