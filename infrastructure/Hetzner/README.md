# Prerequisites
- install ansible via pip ([instructions](https://docs.ansible.com/ansible/latest/installation_guide/intro_installation.html#installing-ansible-with-pip))
- install terraform via package manager

# How to deploy to test server on digital ocean
1. Fill in credentials at `Hetzner/export_vars`. Below is an example of what needs be include:
   ```bash
   export DO_TOKEN="ENTER YOUR OWN DO TOKEN!"
   export SSH_FINGERPRINT="ENTER YOUR OWN MD5 SSH_FINGERPRINT"
   ```
 - To generate fingerprint run the command following command: 
    - `$ ssh-keygen -E md5 -lf ~/.ssh/id_rsa.pub | awk '{sub(/^(MD5|md5):/, "", $2); print $2}'`
2. Make sure ansible is installed and `ansible-playbook` command is available
3. Run `make run-test-server-do`
4. Once complete:
   - For prometheus: run`make ssh-tunnel-test-server-do-prometheus` and browse to `localhost:9090`
   - For grafana: run`make ssh-tunnel-test-server-do-grafana` and browse to `localhost:3000`
     - The default user and password are admin and admin.
5. To clean up, `make destroy-test-server-do`

# Monitoring backup of besu data
- Backup scripts location on hetzner machine: `/opt/scripts/`.
- Ansible configures a cronjob which creates a backup of the besu data folder.
- Basic logs are forwared to the #hetzner slack group.
- For more detailed debug logs, you can check `/var/log/backup_scripts/`, which contains
the last runs' debug logs.

### Running backup script manually
- Login to server
- Look at scripts in `/opt/scripts`, decide which backup you want to run, then:
- Run the following command as root: `bash /opt/scripts/{script}`


# Storage box backups
- Backups are made to a hetzner storage box, full documenation [here](https://wiki.hetzner.de/index.php/Storage_Boxes/en#Storage_Box).
- According to hetzner customer service:
> All Storage Boxes are hard drive based and running on local storage with ZFS.
The host machines are all running an RAID and withstand multiple drive failures.
You should get at least read/write speeds of about 50MB/s.

## Initial configuration of storage box

- The hetzner server will make backups to the storage box using the rsync command. This means the storage
box will have to be configured to accept ssh connections from the hetzner server.

## Configure SSH access

### Step 1: Create ssh-key (only if adding a new server)
- if you are adding a new server, you will have to create [ssh keys pairs](./ansible/files)
- to create keys for ssh access, leave the password field blank
- `ssh-keygen -t rsa -C {server name} -f ./{server name}`
- example `ssh-keygen -t rsa -C server-1 -f ./ {server name}`
- this will create two files in the current directory: `{server name}` and `{server name}`.pub

### Step 2: Copy `{server name}`.pub contents to storage_box's `.ssh/authorized_keys` file
- Log into the hetzner website and navigate to the [storage box settings](https://robot.your-server.de/storage)
- Get a password by clicking on **reset password** (since we don't rely on passwords for access, we can reset it without
worry)
- Mount the storage box locally with the following command:
  - `sudo  mount -t cifs {Samba/CIFS share url} {local mount path} -o username={username},password={password},workgroup=workgroup,iocharset=utf8`
    - `Samba/CIFS share url, username`: this can be found on the [storage box settings](https://robot.your-server.de/storage), on the **storage box data** tab
    - `local mount path` this a local folder where the network drive will be mounted.
    - `password`: this is the password you obtained by doing a password reset.
  - example command: `sudo  mount -t cifs //u233493.your-storagebox.de/backup /tmp/storage/ -o username=u233493,password=IKVM9SOSpaJfIeiV,workgroup=workgroup,iocharset=utf8`
- Navigate to the mounted storage and run the following:
  - (instructions copied from [here](https://wiki.hetzner.de/index.php/Backup_Space_SSH_Keys/en))
  - if `./.ssh` doesn't exist: `sudo mkdir .ssh && chmod 700 .ssh`
  - if `./.ssh/authorized_keys` doesn't exit`sudo touch .ssh/authorized_keys && chmod 600 .ssh/authorized_keys`
  - add the contents of `{server name}.pub` created previously to a new line in `.ssh/authorized_keys`

### Test Rsync access

- To test the key has properly been added
- Add the ssh-key with: `ssh-add {server name}.pub`
- create a test folder to sync
- test command example: `rsync --progress -e 'ssh -p23' --recursive {local test folder to sync} u233493@u233493.your-storagebox.de:/home`
- note the `/home` path, `/home/` actually corresponds to the root of the network storage.
- once you've confirmed things work, clean up the test files which were rsynced



# Important notes
- ansible playbooks are currently only configured to run on **UBUNTU 16.04 XENIAL** hosts.
This is due to the hetzner's default ubuntu image being 16.04 xenial.
- the test server must be deleted when not used because currently we don't share the terraform state and there will be issues if multiple people try to deploy it at once.
- if the TLS certificates or DNS need to be tested on the test server, make sure to assign the floating ip 64.225.81.78 to the test server droplet in the DigitalOcean console https://cloud.digitalocean.com/networking/floating_ips?i=977411
