resource "google_container_cluster" "gke_cluster" {
  name     = var.cluster_name
  location = var.cluster_location

  # We recommend managing node pools separately.
  # This removes the default node pool created with the cluster.
  remove_default_node_pool = true
  initial_node_count       = 1 # Required, but the pool will be removed.
}

resource "google_container_node_pool" "primary_node_pool" {
  name       = "${var.cluster_name}-primary-pool"
  cluster    = google_container_cluster.gke_cluster.name
  location   = var.cluster_location
  node_count = var.node_count

  node_config {
    machine_type = var.machine_type
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]
  }
}