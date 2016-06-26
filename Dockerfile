FROM ubuntu:14.04
MAINTAINER Daniel Schauenberg d@unwiredcouch.com

# General env setup
ENV POCKETCLEANER_VERSION=1:0.1.2-1

# install pocketcleaner from packagecloud
RUN apt-get update -qq && apt-get install -y curl
RUN curl https://packagecloud.io/install/repositories/mrtazz/pocketcleaner/script.deb.sh | bash
RUN apt-get install -y pocketcleaner=${POCKETCLEANER_VERSION}

ENV POCKETCLEANER_KEEP_COUNT=20

COPY  entrypoint.sh  /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
