terraform {
  required_version = ">= 1.5"

  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.34.1"
    }
  }
}

provider "digitalocean" {
  token = var.do_token
}

resource "digitalocean_droplet" "api_server" {
  name   = "api-locadora"
  region = "sfo2"
  size   = "s-1vcpu-512mb-10gb"
  image  = "ubuntu-22-04-x64"

  ssh_keys = [var.ssh_fingerprint]

  user_data = <<EOF
#!/bin/bash

apt update -y
apt install -y apt-transport-https ca-certificates curl software-properties-common

# Instalar Docker
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"

apt update -y
apt install -y docker-ce docker-ce-cli containerd.io

# Instalar Docker Compose
curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

# Permitir Docker sem sudo
usermod -aG docker root

EOF
}
