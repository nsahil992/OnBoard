resource "aws_vpc" "main_vpc" {
    cidr_block = "10.0.0.0/16"
    enable_dns_support = true
    enable_dns_hostname = true

    tags = {
        Name = "onboard-vpc"
  }
}