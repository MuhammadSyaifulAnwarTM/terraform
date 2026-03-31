resource "aws_eks_node_group" "gpu_nodes" {
  cluster_name    = aws_eks_cluster.ai_cluster.name
  instance_types = ["g4dn.xlarge"]
  # Infrastructure as Code for AI-first cloud platforms