FROM registry:2.7.1

COPY ./entrypoint.sh /heroku-entrypoint.sh
COPY ./config.yml /etc/docker/registry/config.yml

ENTRYPOINT ["/heroku-entrypoint.sh"]
CMD ["/etc/docker/registry/config.yml"]
