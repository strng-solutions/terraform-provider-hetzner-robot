terraform {
  required_providers {
    hetzner-robot = {
      version = "0.1.0"
      source  = "strng-solutions/hetzner-robot"
    }
  }

}

provider "hetzner-robot" {
  username = "yourUserNameFromRobot"
  password = "yourPasswordFromRobot"
}

resource "hetzner-robot_firewall" "firewall" {
  server_ip     = "1.1.1.1"
  active        = true
  whitelist_hos = true

  rule {
    name     = "Allow ssh"
    src_ip   = "0.0.0.0/0"
    src_port = "0-65535"
    dst_ip   = "0.0.0.0/0"
    dst_port = "22"
    protocol  = "tcp"
    #example can be one of syn, fin, rst, psh, urg or more like "syn|fin"
    tcp_flags= "syn"
    action   = "accept"
  }

  rule {
    name     = "Allow inbound"
    src_ip   = "0.0.0.0/0"
    src_port = "0-65535"
    dst_ip   = "0.0.0.0/0"
    dst_port = "32768-65535"
    action   = "accept"
  }

  rule {
    name     = "Allow ICMP"
    src_ip   = "0.0.0.0/0"
    src_port = "0-65535"
    dst_ip   = "0.0.0.0/0"
    dst_port = "0-65535"
    protocol = "icmp"
    action   = "accept"
  }


  rule {
    name     = "Deny others"
    src_ip   = "0.0.0.0/0"
    src_port = "0-65535"
    dst_ip   = "0.0.0.0/0"
    dst_port = "0-65535"
    action   = "discard"
  }
}
