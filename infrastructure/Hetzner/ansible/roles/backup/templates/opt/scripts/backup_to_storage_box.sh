#!/bin/bash
# {{ ansible_managed }}

set -e # Exit on error
set -x # debug

# Script to backup archive data to hetzner storage box.

function timestamp_echo()
{
    msg="$(date)@$(hostname): $1"
    curl -X POST --silent -H 'Content-type: application/json' --data "{\"text\":\" $msg \"}" {{ backup_slack_url }}

    echo -e $msg
}

{% for key in besu_instances %}
if ! docker ps | grep {{ besu_instances[key].name }} 1>/dev/null
then
    timestamp_echo "{{ besu_instances[key].name }} not running, exiting....\n"
    exit 1
fi

if ! curl --max-time 5 --silent {{ besu_instances[key].host}}:{{ besu_instances[key].rpc_http_port}}/liveness | grep UP 1>/dev/null
then
    timestamp_echo "besu /liveness check failed, exiting...\n"
    exit 1
fi
{% endfor %}

snapshot_name=vg1_lv_data_to_lv_storage_box_snapshot
mountpoint={{ backup_to_storage_box_mountpoint }}

{% for key in besu_instances %}
timestamp_echo "Stopping  {{ besu_instances[key].name }} container..."
docker stop {{ besu_instances[key].name }}
timestamp_echo "{{ besu_instances[key].name }} container stopped...\n"
{% endfor %}

timestamp_echo "Creating snapshot of lv_data..."
lvcreate --size {{ backup_snapshot_size_GB }}G --snapshot --name ${snapshot_name} /dev/vg1/lv_data
timestamp_echo "${snapshot_name} created created...\n"

{% for key in besu_instances %}
timestamp_echo "Starting {{ besu_instances[key].name }} container..."
docker start {{ besu_instances[key].name }}
timestamp_echo "{{ besu_instances[key].name }} container started...\n"
{% endfor %}

timestamp_echo "Mounting snapshot..."
mount  -o nouuid /dev/vg1/${snapshot_name} ${mountpoint}
timestamp_echo "Snapshot mounted...\n"


pushd ${mountpoint}

{%for box_config in backup_storage_box_configs %}
storage_box_path=rsync/lv_data
mkdir -p /tmp/folder_structure/${storage_box_path}
timestamp_echo "rsyncing {{ box_config.besu_node }} data to {{ box_config.box_url }} ..."

# We first need to run a rsync command to create the folder if it doesn't exist.
rsync --recursive --progress -e 'ssh -o StrictHostKeyChecking=no -i /root/.ssh/server_id_rsa -p23' /tmp/folder_structure/ {{ box_config.box_url }}:/home/{{ inventory_hostname }}
rm -r /tmp/folder_structure
rsync -a -z --progress --delete -e 'ssh -o StrictHostKeyChecking=no -i /root/.ssh/server_id_rsa -p23' ./{{ box_config.besu_node }}  {{ box_config.box_url }}:/home/{{ inventory_hostname }}/${storage_box_path}
timestamp_echo "rsync complete..."

{% endfor %}

popd
timestamp_echo "all rsyncs complete..."


timestamp_echo "Unmount ${snapshot_name}..."
umount ${mountpoint}
timestamp_echo "Unmounted...\n"

timestamp_echo "Delete ${snapshot_name}..."
lvremove -f vg1/${snapshot_name}
timestamp_echo "Deleted ${snapshot_name}...\n"

timestamp_echo "(space occupied on lvm) $(du -h --max-depth 1 {{ lvm_lvs.lv_data.lv_mountpoint }})"
