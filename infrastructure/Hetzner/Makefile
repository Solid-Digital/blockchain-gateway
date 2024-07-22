SHELL:=/bin/bash
SHELLFLAGS:=-e

# VARIABLES
region=ams3
tags=all
ansible_base_command=ansible-playbook ${playbook} --diff --inventory ${inventory} --limit ${host} --tags=${tags}
ssh_tunnel_base_command=ssh -L ${port}:127.0.0.1:${port} root@${server_ip} -N
ip_regex:=([0-9]{1,3}\.){3}[0-9]{1,3}

# CONFIGURATION VIA ANSIBLE
run: playbook=setup.yml
run:
	cd ansible && ${ansible_base_command}

check: playbook=setup.yml
check:
	cd ansible && ${ansible_base_command} --check

reset: playbook=reset_besu.yml
reset:
	cd ansible && ${ansible_base_command}

delete-all-lvm: playbook=delete_all_lvm.yml
delete-all-lvm:
	cd ansible && ${ansible_base_command}

run-test-server-do: host=test_server_do
run-test-server-do: inventory=test_server_do_inventory.ini
run-test-server-do: provision-test-server-do run

run-hetzner: host=hetzner
run-hetzner: inventory=hetzner_inventory.ini
run-hetzner: run

check-test-server-do: host=test_server_do
check-test-server-do: inventory=test_server_do_inventory.ini
check-test-server-do: check

check-hetzner: host=hetzner
check-hetzner: inventory=hetzner_inventory.ini
check-hetzner: check

reset-test-server-do: host=test_server_do
reset-test-server-do: inventory=test_server_do_inventory.ini
reset-test-server-do: reset

reset-hetzner: host=hetzner
reset-hetzner: inventory=hetzner_inventory.ini
reset-hetzner: reset

delete-all-lvm-test-server-do: host=test_server_do
delete-all-lvm-test-server-do: inventory=test_server_do_inventory.ini
delete-all-lvm-test-server-do: delete-all-lvm

# delete-all-lvm-hetzner: host=hetzner # uncomment to delete all lvs on hetzner
delete-all-lvm-hetzner: inventory=hetzner_inventory.ini
delete-all-lvm-hetzner: delete-all-lvm

# PROVISIONING VIA TERRAFORM
test_server_do/.terraform:
	cd test_server_do && terraform init

provision-test-server-do: test_server_do/.terraform
	source ./export_vars && cd test_server_do && terraform apply -var "ssh_fingerprint=$$SSH_FINGERPRINT" -var "do_token=$$DO_TOKEN" -var "region=${region}"
	cp test_server_do/test_server_do_inventory.ini ansible/test_server_do_inventory.ini

destroy-test-server-do:
	source ./export_vars && cd test_server_do && terraform destroy -var "ssh_fingerprint=$$SSH_FINGERPRINT" -var "do_token=$$DO_TOKEN" -var "region=${region}"

# SSH TUNNELLING

ssh-tunnel:
	${ssh_tunnel_base_command}

ssh-tunnel-test-server-do-grafana: port=3000
ssh-tunnel-test-server-do-grafana: server_ip=$(shell egrep -o '${ip_regex}' ansible/test_server_do_inventory.ini)
ssh-tunnel-test-server-do-grafana: ssh-tunnel

ssh-tunnel-hetzner-grafana: port=3000
ssh-tunnel-hetzner-grafana: server_ip=$(shell egrep -o '${ip_regex}' ansible/hetzner_inventory.ini)
ssh-tunnel-hetzner-grafana: ssh-tunnel

ssh-tunnel-test-server-do-prometheus: port=9090
ssh-tunnel-test-server-do-prometheus: server_ip=$(shell egrep -o '${ip_regex}' ansible/test_server_do_inventory.ini)
ssh-tunnel-test-server-do-prometheus: ssh-tunnel

ssh-tunnel-hetzner-prometheus: port=9090
ssh-tunnel-hetzner-prometheus: server_ip=$(shell egrep -o '${ip_regex}' ansible/hetzner_inventory.ini)
ssh-tunnel-hetzner-prometheus: ssh-tunnel
