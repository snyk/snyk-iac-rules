resource "aws_security_group" "denied" {
  name        = "allow_tcp"
  description = "Allow TCP inbound from anywhere"
  vpc_id      = "arn"

  ingress {
    from_port   = 3389
    to_port     = 3389
    protocol    = "tcp"
    cidr_blocks = ["::/0"]
  }
}
