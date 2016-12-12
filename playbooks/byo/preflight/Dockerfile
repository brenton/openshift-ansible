# This builds a preflight playbook image for use with Docker.
# Look at playbooks/byo/preflight/README.md for more information.
# To build, use hack/build-images.sh

FROM docker.io/aweiteka/playbook2image

MAINTAINER "Customer Success team" <dev@lists.openshift.redhat.com>
ADD . /opt/app-root/src/
ENV PLAYBOOK_FILE=playbooks/byo/preflight/check.yml \
    INVENTORY_FILE=hosts \
    ANSIBLE_HOST_KEY_CHECKING=False \
    ALLOW_ANSIBLE_CONNECTION_LOCAL=false
CMD /usr/libexec/s2i/run
USER 0
