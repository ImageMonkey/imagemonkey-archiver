FROM golang:1.17-buster

# when changing the base image, change also the debian version here
ARG DEBIAN_VERSION=buster

MAINTAINER bbernhard version: 0.1

RUN echo "deb http://apt.postgresql.org/pub/repos/apt/ ${DEBIAN_VERSION}-pgdg main" > /etc/apt/sources.list.d/pgdg.list \
    && wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add - \
	&& apt-get update \
	&& apt-get install -y postgresql-9.6 postgresql-server-dev-9.6 postgresql-contrib-9.6 postgresql-9.6-postgis-2.5 python3 wget git dos2unix


RUN adduser imagemonkey --disabled-password --gecos "First Last,RoomNumber,WorkPhone,HomePhone" --home /home/imagemonkey

RUN mkdir -p /tmp/imagemonkey-archiver
COPY src /tmp/imagemonkey-archiver/src/
COPY env/postgres /tmp/imagemonkey-archiver/env/postgres/

RUN /bin/su - imagemonkey -c "mkdir -p /home/imagemonkey/go" \
	&& /bin/su - imagemonkey -c "mkdir -p /home/imagemonkey/bin" \
	&& chown -R imagemonkey:imagemonkey /tmp/imagemonkey-archiver \
	&& chmod -R u+rwx /tmp/imagemonkey-archiver

ENV GOPATH=/home/imagemonkey/go
ENV GOBIN=/home/imagemonkey/bin

RUN cd /tmp/imagemonkey-archiver/src/ && /usr/local/go/bin/go install archiver.go obfuscate.go verify.go

# change listen address in postgres config file to localhost + add custom pg_hba.conf
RUN echo listen_addresses='localhost' >> /etc/postgresql/9.6/main/postgresql.conf \
 && cp /tmp/imagemonkey-archiver/env/postgres/pg_hba.conf /etc/postgresql/9.6/main/pg_hba.conf

RUN cp /tmp/imagemonkey-archiver/src/scripts/run_archiver.py /home/imagemonkey/bin/run_archiver.py \
	&& cp /tmp/imagemonkey-archiver/src/scripts/start_postgres.sh /home/imagemonkey/bin/start_postgres.sh \
	&& chown -R imagemonkey:imagemonkey /home/imagemonkey/bin \
	&& chmod -R u+rx /home/imagemonkey/bin \
	&& chmod u+rx /home/imagemonkey/bin/start_postgres.sh \
	&& dos2unix /home/imagemonkey/bin/start_postgres.sh \
	&& /bin/bash /home/imagemonkey/bin/start_postgres.sh \
	&& /bin/su - postgres -c "psql -c \"CREATE database imagemonkey;\"" \
	&& /bin/su - postgres -c "psql -d imagemonkey -c \"CREATE USER monkey WITH PASSWORD 'imagemonkey';\""

CMD ["python3", "/home/imagemonkey/bin/run_archiver.py"]
