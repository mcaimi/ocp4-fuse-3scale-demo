FROM docker.io/library/openjdk:11-jre-slim
LABEL maintainer=mcaimi@redhat.com
LABEL runner=openjdk

ARG nexus_repo_url
ARG artifact_name
ARG nexus_auth
ARG user=runner
ARG group=runner
ARG uid=1000
ARG gid=1000


USER root
COPY scripts/nexus-fetch.py /usr/local/bin/nexus-fetch
RUN chmod +x /usr/local/bin/nexus-fetch &&\
    apt-get update && apt-get install -y python3=3.7.3-1 --no-install-recommends && rm -Rf /var/lib/apt/lists/* &&\
    mkdir -p /var/lib/deploy && /usr/local/bin/nexus-fetch -r ${nexus_repo_url} -a ${artifact_name} -d /var/lib/deploy -c ${nexus_auth} &&\
    groupadd -g ${gid} ${group} && useradd -l -c "Java Runner User" -d /home/${user} -u ${uid} -g ${gid} -m ${user}

EXPOSE 8080

USER ${user}
CMD [ "java", "-jar", "/var/lib/deploy/runner.jar" ]
