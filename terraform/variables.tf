# Service Configuration
variable "service_name" {
  description = "Name of the service"
  type        = string
  default     = "ttl-sh"
}

variable "environment" {
  description = "Environment name (staging, production)"
  type        = string
  default     = "production"
}

# Network Configuration
variable "network_cidr" {
  description = "CIDR range for the main network"
  type        = string
  default     = "10.0.0.0/16"
}

variable "subnet_cidr" {
  description = "CIDR range for the subnet"
  type        = string
  default     = "10.0.0.0/24"
}

# Compute Configuration
variable "server_image" {
  description = "Image to use for servers"
  type        = string
  default     = "ubuntu-24.04"
}

variable "server_type" {
  description = "Server type for compute instances"
  type        = string
  default     = "cpx41"
}


# Location Configuration
variable "location" {
  description = "Hetzner Cloud location"
  type        = string
  default     = "ash"
}

# Label Configuration
variable "common_labels" {
  description = "Common labels to apply to all resources"
  type        = map(string)
  default     = {}
}

variable "compute_tier_labels" {
  description = "Additional labels for compute tier resources"
  type        = map(string)
  default = {
    tier = "compute"
    role = "worker"
  }
}