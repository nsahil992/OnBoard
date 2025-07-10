resource "aws_subnet" "public_subnet_1" {
    vpc_id = "aws_vpc.main_vpc.id"
    cidr_block = "10.0.1.0/24"
    availability_zone = "ap-south-1a"
    map_public_ip_on_launch = true

      tags = {
        Name = "onboard-public-subnet"
      }
}

resource "aws_subnet" "private_subnet_1" {
  vpc_id            = aws_vpc.main_vpc.id
  cidr_block        = "10.0.2.0/24"
  availability_zone = "ap-south-1a"

  tags = {
    Name = "onboard-private-subnet"
  }
}