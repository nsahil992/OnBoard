resource "aws_instance" "onboard_ec2" {
    ami = "ami-0f58b397bc5c1f2e8"
    instance_type = "t2micro"
    subnet_id = "aws_subnet.public_subnet_1.id"
    associate_public_ip_address = true
    key_name = "onboard-key"
    vpc_security_group_ids = [aws_security_group.ec2_sg.id]
    tags = {
        Name = "onboard_ec2"
    }
}