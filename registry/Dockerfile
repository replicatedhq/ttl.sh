FROM registry:2.7.1

ADD ./entrypoint.sh /entrypoint.sh
ADD ./config.yml /etc/docker/registry/config.yml
ADD ./garbage-collect.sh /garbage-collect.sh

ENTRYPOINT ["/entrypoint.sh"]
CMD ["/etc/docker/registry/config.yml"]
