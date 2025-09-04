# Data Sources to query available locations
data "hcloud_locations" "all" {}

# Select a specific location (you can change the filter criteria)
data "hcloud_location" "main" {
  name = var.location
}

# Query all available SSH keys in the account
data "hcloud_ssh_keys" "all_keys" {}

# Centralized label management
locals {
  # Common labels for all resources
  common_labels = merge({
    service     = var.service_name
    environment = var.environment
  }, var.common_labels)
  
  # Compute tier labels (merged with common)
  compute_labels = merge(local.common_labels, var.compute_tier_labels)
}

# Network Configuration
resource "hcloud_network" "main_network" {
  name     = "${var.service_name}-network"
  ip_range = var.network_cidr
}

# Private Subnet
resource "hcloud_network_subnet" "main_subnet" {
  network_id   = hcloud_network.main_network.id
  type         = "cloud"
  network_zone = data.hcloud_location.main.network_zone
  ip_range     = var.subnet_cidr
}

# Firewall Configuration
resource "hcloud_firewall" "main_firewall" {
  name = "${var.service_name}-firewall"
  
  # Allow ICMP (ping)
  rule {
    direction = "in"
    protocol  = "icmp"
    source_ips = [
      "0.0.0.0/0",
      "::/0"
    ]
  }
  
  # Allow HTTP
  rule {
    direction = "in"
    protocol  = "tcp"
    port      = "80"
    source_ips = [
      "0.0.0.0/0",
      "::/0"
    ]
  }
  
  # Allow HTTPS
  rule {
    direction = "in"
    protocol  = "tcp"
    port      = "443"
    source_ips = [
      "0.0.0.0/0",
      "::/0"
    ]
  }
  
  # Allow SSH (optional but usually needed)
  rule {
    direction = "in"
    protocol  = "tcp"
    port      = "22"
    source_ips = [
      "0.0.0.0/0",
      "::/0"
    ]
  }
}

# Compute Server Instance (can be scaled to multiple instances)
resource "hcloud_server" "instance1" {
  name         = "${var.service_name}-compute-1"  # Named to indicate it's a compute node
  image        = var.server_image
  server_type  = var.server_type
  location     = data.hcloud_location.main.name
  firewall_ids = [hcloud_firewall.main_firewall.id]
  ssh_keys     = data.hcloud_ssh_keys.all_keys.ssh_keys[*].id
  
  labels = local.compute_labels
  
  network {
    network_id = hcloud_network.main_network.id
    # IP will be auto-assigned from the subnet
  }
  
  public_net {
    ipv4_enabled = true  # Public IP for SSH management access
    ipv6_enabled = false
  }
  
  depends_on = [
    hcloud_network_subnet.main_subnet
  ]
}


# Outputs for debugging and DNS configuration
output "server_primary_ip" {
  value       = hcloud_server.instance1.ipv4_address
  description = "The server's primary IPv4 address to use for ttl.sh DNS A record"
}

output "server_ipv6" {
  value       = hcloud_server.instance1.ipv6_address
  description = "The server's IPv6 address"
}

output "selected_location" {
  value = {
    name         = data.hcloud_location.main.name
    description  = data.hcloud_location.main.description
    network_zone = data.hcloud_location.main.network_zone
  }
  description = "Currently selected location details"
}

output "ssh_keys_used" {
  value = {
    for key in data.hcloud_ssh_keys.all_keys.ssh_keys :
    key.name => key.fingerprint
  }
  description = "SSH keys that will be added to the server"
}

output "server_details" {
  value = {
    name       = hcloud_server.instance1.name
    private_ip = hcloud_server.instance1.network[*].ip
    labels     = hcloud_server.instance1.labels
  }
  description = "Details of the deployed server instance"
}