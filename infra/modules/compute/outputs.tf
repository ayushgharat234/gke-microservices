output "cluster_name" {
  description = "The name of the created GKE cluster"
  value = google_container_cluster.gke_cluster.name
}

output "cluster_endpoint" {
    description = "The endpoint of the GKE cluster"
    value = google_container_cluster.gke_cluster.endpoint
}