FROM registry:2.7.1

ADD ./entrypoint.sh /entrypoint.sh
ADD ./config.yml /etc/docker/registry/config.yml

ENTRYPOINT ["/entrypoint.sh"]
CMD ["/etc/docker/registry/config.yml"]
