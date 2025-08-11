variable "project_id" {
  description = "The Project ID for the GCP Project"
  type = string
}

variable "region" {
  description = "Region for all resources"
  type = string
  default = "us-central1"
}

variable "cluster_name" {
  description = "The name of the GKE cluster."
  type        = string
  default     = "tekton-cicd-cluster"
}

variable "node_count" {
  description = "The initial number of nodes in the cluster."
  type        = number
  default     = 2
}

variable "machine_type" {
  description = "The machine type for the cluster's nodes."
  type        = string
  default     = "e2-medium"
}