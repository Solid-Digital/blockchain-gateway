terraform {
  required_providers {
    digitalocean = {
      source = "digitalocean/digitalocean"
      required_version = ">= 0.13"
    }
  }
}


variable "do_token" {}

variable "ssh_fingerprint" {
  default=""
}

provider "digitalocean" {
  token = var.do_token
}

variable "region" {
  type = string
}

variable "volume_size" {
  default = 35  # hetzner has 35TB disks
}
resource "digitalocean_droplet" "test-server-do" {
  image = "ubuntu-16-04-x64"
  name = "test-server-do"
  region = var.region
  size = "s-2vcpu-4gb"
  ssh_keys = [
    "${var.ssh_fingerprint}"
  ]

  connection {
    user = "root"
    type = "ssh"
    timeout = "2m"
    agent = true
  }

}

resource "digitalocean_floating_ip_assignment" "test-server-do_floating_ip_assignment" {
  ip_address = "64.225.81.78"
  droplet_id = digitalocean_droplet.test-server-do.id
}

resource "digitalocean_volume" "lvm-volume-0" {
  name = "volume-0"
  region = var.region
  size = var.volume_size
  initial_filesystem_type = "ext4"
}

resource "digitalocean_volume" "lvm-volume-1" {
  name = "volume-1"
  region = var.region
  size = var.volume_size
  initial_filesystem_type = "ext4"
}

resource "digitalocean_volume" "lvm-volume-2" {
  name = "volume-2"
  region = var.region
  size = var.volume_size
  initial_filesystem_type = "ext4"
}

resource "digitalocean_volume" "lvm-volume-3" {
  name = "volume-3"
  region = var.region
  size = var.volume_size
  initial_filesystem_type = "ext4"
}

resource "digitalocean_volume_attachment" "node_attachment-0"{
  droplet_id = digitalocean_droplet.test-server-do.id
  volume_id = digitalocean_volume.lvm-volume-0.id
}

resource "digitalocean_volume_attachment" "node_attachment-1"{
  droplet_id = digitalocean_droplet.test-server-do.id
  volume_id = digitalocean_volume.lvm-volume-1.id
}

resource "digitalocean_volume_attachment" "node_attachment-2"{
  droplet_id = digitalocean_droplet.test-server-do.id
  volume_id = digitalocean_volume.lvm-volume-2.id
}

resource "digitalocean_volume_attachment" "node_attachment-3"{
  droplet_id = digitalocean_droplet.test-server-do.id
  volume_id = digitalocean_volume.lvm-volume-3.id
}

resource "null_resource" "local_inventory"{
  provisioner "local-exec" {
    command = " echo test_server_do ansible_ssh_host=${digitalocean_droplet.test-server-do.ipv4_address} ansible_port=22 ansible_user=root > test_server_do_inventory.ini"
  }
}

output "test_server_do_ip" {
  value = digitalocean_droplet.test-server-do.ipv4_address
}
