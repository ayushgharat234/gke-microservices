module "gke_cluster" {
  source = "./modules/compute"
  cluster_name = var.cluster_name
  cluster_location = var.region
  node_count = var.node_count
  machine_type = var.machine_type
}