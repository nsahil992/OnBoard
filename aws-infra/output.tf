# Public IP of the EC2 instance – needed for SSH and Ansible
output "ec2_public_ip" {
  description = "Public IP of the EC2 instance"
  value       = aws_instance.onboard_ec2.public_ip
}

# VPC ID – useful for debugging and reference
output "vpc_id" {
  description = "ID of the created VPC"
  value       = aws_vpc.main_vpc.id
}

# Subnet IDs – helpful if you ever need to associate more resources
output "public_subnet_id" {
  description = "ID of the public subnet"
  value       = aws_subnet.public_subnet_1.id
}

output "private_subnet_id" {
  description = "ID of the private subnet"
  value       = aws_subnet.private_subnet_1.id
}

# NAT Gateway – for private resources to reach the internet
output "nat_gateway_id" {
  description = "ID of the NAT Gateway"
  value       = aws_nat_gateway.nat.id
}

# Security Group ID – for attaching to other services later
output "ec2_security_group_id" {
  description = "Security Group ID used for EC2 instance"
  value       = aws_security_group.ec2_sg.id
}
