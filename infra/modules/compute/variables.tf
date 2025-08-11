variable "cluster_name" {
  description = "The name of the cluster"
  type        = string
}

variable "cluster_location" {
  description = "The location of the cluster"
  type        = string
}

variable "node_count" {
  description = "The number of nodes in the cluster"
  type        = number
}

variable "machine_type" {
  description = "The machine type for the cluster's nodes"
  type        = string
}