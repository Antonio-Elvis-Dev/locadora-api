output "public_ip" {
  description = "68.183.174.239"
  value       = digitalocean_droplet.api_server.ipv4_address
}
