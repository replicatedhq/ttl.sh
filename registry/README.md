# replreg registry

This is the registry server that runs in ttl.sh. It's a configured Docker registry, using the registry:2 image. There the confguration installs a hook into the registry, that will fire notifications when an image is pushed. The work to expire and delete images happens in the hook directory.

