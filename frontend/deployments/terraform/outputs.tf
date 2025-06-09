# frontend/deployments/terraform/outputs.tf
output "load_balancer_dns" {
  description = "Load Balancer DNS name"
  value       = aws_lb.autosysadmin.dns_name
}

output "ecr_backend_repository_url" {
  description = "Backend ECR repository URL"
  value       = aws_ecr_repository.backend.repository_url
}

output "ecr_frontend_repository_url" {
  description = "Frontend ECR repository URL"
  value       = aws_ecr_repository.frontend.repository_url
}