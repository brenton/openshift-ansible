# Preflight checks

Here we provide an ansible playbook (currently just one) for detecting
potential roadblocks prior to install or upgrade.

Note that currently this is only useful for RPM-based Enterprise
installations, based on supported Red Hat RPMs. Containerized and
Origin installs are excluded from checks for now.

## Running directly

For a simple trial, with an installation of Ansible 2.0 or greater,
run the playbook directly against your inventory file:

    $ git clone https://github.com/rhcarvalho/openshift-ansible.git -b pre-flight-checks
    $ ansible-playbook -i <inventory file> openshift-ansible/playbooks/byo/preflight/check.yml

## Running in a container

Rather than installing Ansible 2.0+ somewhere, you may prefer to use
the containerized version of this same playbook. This image is built
on [playbook2image](https://github.com/aweiteka/playbook2image).

With the docker daemon running, you can run a container with an ssh key
and an inventory file mounted in and get a preflight check, like so:

    $ docker run -it -v ~/.ssh/id_rsa:/root/.ssh/id_rsa \
                     -v ~/hosts:/opt/app-root/src/hosts \
                     docker.io/sosiouxme/ocp_preflight_playbook

Some limitations on this method:

* This runs as root so as to avoid key file ownership issues.
* The `id_rsa` file needs to be `chmod 0600` or ssh refuses to use it.
* The default is not to do ssh strict host checking so that you will not be
  prompted for each host that is new to the container. However for more
  control over ssh config, you can mount in a `.ssh/` directory with a
  `config` and/or `known_hosts` file. Then to enable strict host checking,
  add `-e ANSIBLE_HOST_KEY_CHECKING=True` into the docker run.
* There is no good way to enter a key passphrase.

If you have a passphrase on the ssh key or other complications requiring
terminal input, you can shell into the container instead; for example:

    $ docker run -it -v ~/.ssh/id_rsa:/root/.ssh/id_rsa  \
                     -v ~/hosts:/opt/app-root/src/hosts  \
                     docker.io/sosiouxme/ocp_preflight_playbook \
                     /bin/bash
    # eval `ssh-agent`
    Agent pid 7
    # ssh-add /root/.ssh/id_rsa
    Enter passphrase for /root/.ssh/id_rsa: <invisible>
    Identity added: /root/.ssh/id_rsa (/root/.ssh/id_rsa)

And then proceed to run the playbook:

    # /usr/libexec/s2i/run

