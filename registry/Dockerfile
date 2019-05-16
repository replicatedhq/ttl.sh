FROM registry:2.7.1

ADD ./entrypoint.sh /heroku-entrypoint.sh
ADD ./config.yml /etc/docker/registry/config.yml

ENTRYPOINT ["/heroku-entrypoint.sh"]
CMD ["/etc/docker/registry/config.yml"]
