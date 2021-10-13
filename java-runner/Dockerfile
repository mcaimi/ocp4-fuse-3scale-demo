ARG IMAGE=docker.io/openliberty/open-liberty:kernel-slim-java11-openj9-ubi
FROM ${IMAGE} as staging
LABEL maintainer=mcaimi@redhat.com
LABEL runner=openjdk

ARG nexus_repo_url
ARG artifact_name
ARG nexus_auth
ARG user=runner
ARG group=runner
ARG uid=1001
ARG gid=1001

USER root
COPY scripts/nexus-fetch.py /usr/local/bin/nexus-fetch
RUN chmod +x /usr/local/bin/nexus-fetch &&\
    dnf install -y shadow-utils python39 python39-requests &&\
    mkdir -p /var/lib/deploy && /usr/local/bin/nexus-fetch -r ${nexus_repo_url} -a ${artifact_name} -d /var/lib/deploy -c ${nexus_auth} &&\
    groupadd -g ${gid} ${group} && useradd -l -c "Java Runner User" -d /home/${user} -u ${uid} -g ${gid} -m ${user}

# Copy server config so springBootUtility can be downloaded by featureUtility in the next step
COPY --chown=${uid}:0 server.xml /config/server.xml

# This script will add the requested XML snippets to enable Liberty features and grow image to be fit-for-purpose using featureUtility
RUN features.sh

COPY --chown=${uid}:0 /var/lib/deploy/runner.jar /staging/myFatApp.jar

RUN springBootUtility thin \
 --sourceAppPath=/staging/myFatApp.jar \
 --targetThinAppPath=/staging/myThinApp.jar \
 --targetLibCachePath=/staging/lib.index.cache

FROM ${IMAGE}
COPY --chown=${uid} :0 server.xml /config

# This script will add the requested XML snippets to enable Liberty features and grow image to be fit-for-purpose using featureUtility
RUN features.sh

COPY --from=staging /staging/lib.index.cache /lib.index.cache
COPY --from=staging /staging/myThinApp.jar /config/dropins/spring/myThinApp.jar

ARG VERBOSE=false
# This script will add the requested server configurations, apply any iFixes and populate caches to optimize runtime
RUN configure.sh
