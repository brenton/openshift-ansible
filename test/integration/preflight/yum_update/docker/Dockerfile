FROM centos:7
RUN yum install -y rpm-build rpmdevtools createrepo && \
    rpmdev-setuptree && \
    mkdir -p /mnt/localrepo
ADD root /
RUN cd /root/rpmbuild/SOURCES && \
    mkdir break-yum-update-1.0 && \
    tar zfc foo.tgz break-yum-update-1.0 && \
    rpmbuild -bb /root/break-yum-update.spec  && \
    yum install -y /root/rpmbuild/RPMS/noarch/break-yum-update-1.0-1.noarch.rpm && \
    rpmbuild -bb /root/break-yum-update-2.spec  && \
    cp /root/rpmbuild/RPMS/noarch/break-yum-update-1.0-2.noarch.rpm /mnt/localrepo && \
    createrepo /mnt/localrepo
