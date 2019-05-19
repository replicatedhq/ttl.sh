FROM node:10 as deps
ADD ./package.json /src/package.json
ADD ./Makefile /src/Makefile
ADD . /src
WORKDIR /src
RUN make deps  test

ENTRYPOINT ["node"]
CMD ["--no-deprecation", "build/server.js", "hooks"]
