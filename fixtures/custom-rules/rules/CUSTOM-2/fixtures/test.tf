
resource "aws_ami" "denied" {
  name                = "aws_ami_not_encrypted"
  virtualization_type = "hvm"
  root_device_name    = "/dev/xvda"

  ebs_block_device {
    device_name = "/dev/xvda"
    encrypted   = false
    volume_size = 8
  }
}