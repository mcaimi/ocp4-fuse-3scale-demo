ARG IMAGE=docker.io/openliberty/open-liberty
ARG LABEL=kernel-slim-java11-openj9-ubi
FROM ${IMAGE}:${LABEL} as staging
LABEL maintainer=mcaimi@redhat.com
LABEL runner=openjdk

ARG nexus_repo_url
ARG artifact_name
ARG nexus_auth

USER root
COPY scripts/nexus-fetch.py /usr/local/bin/nexus-fetch

# Copy server config so springBootUtility can be downloaded by featureUtility in the next step
COPY --chown=1001:0 server.xml /config/server.xml

RUN features.sh && chmod +x /usr/local/bin/nexus-fetch &&\
    dnf install -y shadow-utils python39-3.9.2 python39-requests && dnf clean all &&\
    mkdir -p /var/lib/deploy && /usr/local/bin/nexus-fetch -r ${nexus_repo_url} -a ${artifact_name} -d /var/lib/deploy -c ${nexus_auth} &&\
    mkdir -p /staging && cp /var/lib/deploy/runner.jar /staging/myFatApp.jar && chown -R 1001:0 /staging &&\
    springBootUtility thin \
    --sourceAppPath=/staging/myFatApp.jar \
    --targetThinAppPath=/staging/myThinApp.jar \
    --targetLibCachePath=/staging/lib.index.cache

USER 1001

FROM ${IMAGE}:${LABEL}
COPY --chown=1001:0 server.xml /config

# This script will add the requested XML snippets to enable Liberty features and grow image to be fit-for-purpose using featureUtility
RUN features.sh

COPY --from=staging /staging/lib.index.cache /lib.index.cache
COPY --from=staging /staging/myThinApp.jar /config/dropins/spring/myThinApp.jar

ARG VERBOSE=false
# This script will add the requested server configurations, apply any iFixes and populate caches to optimize runtime
RUN configure.sh
