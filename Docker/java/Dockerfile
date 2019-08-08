FROM java:8
ARG service
RUN apt-get update && apt-get install -y whois
RUN mkdir -p /src/$service
RUN cd /src/$service
WORKDIR /src/$service