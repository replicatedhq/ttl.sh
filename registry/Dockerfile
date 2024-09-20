FROM registry:2.7.1

ADD ./entrypoint.sh /entrypoint.sh
ADD ./config.yml /etc/docker/registry/config.tmpl.yml
ADD ./garbage-collect.sh /garbage-collect.sh

# entrypoint.sh will copy config.tmpl.yml to config.yml and replace variables
RUN rm -f /etc/docker/registry/config.yml

ENTRYPOINT ["/entrypoint.sh"]
CMD ["/etc/docker/registry/config.yml"]
