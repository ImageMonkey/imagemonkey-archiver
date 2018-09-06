FROM debian:9

MAINTAINER bbernhard version: 0.1

RUN apt-get update && apt-get install -y postgresql-9.6 postgresql-server-dev-9.6 postgresql-contrib-9.6 postgis python3 wget git dos2unix

RUN cd /tmp/ \
   && wget https://redirector.gvt1.com/edgedl/go/go1.10.3.linux-amd64.tar.gz \
   && tar -C /usr/local -xzf go1.10.3.linux-amd64.tar.gz \
   && cd /root/

RUN adduser imagemonkey --disabled-password --gecos "First Last,RoomNumber,WorkPhone,HomePhone" --home /home/imagemonkey

RUN mkdir -p /tmp/imagemonkey-archiver
COPY src /tmp/imagemonkey-archiver/src/
COPY env/postgres /tmp/imagemonkey-archiver/env/postgres/

RUN /bin/su - imagemonkey -c "mkdir -p /home/imagemonkey/go" \
 && GOPATH="/home/imagemonkey/go" \
 && /bin/su - imagemonkey -c "mkdir -p /home/imagemonkey/bin" \
 && GOBIN="/home/imagemonkey/bin" \
 && chown -R imagemonkey:imagemonkey /tmp/imagemonkey-archiver \
 && chmod -R u+rwx /tmp/imagemonkey-archiver

RUN /bin/su - imagemonkey -c "cd /tmp/imagemonkey-archiver/src/ && export GOPATH=/home/imagemonkey/go && export GOBIN=/home/imagemonkey/bin && /usr/local/go/bin/go get -d && /usr/local/go/bin/go install archiver.go obfuscate.go verify.go"

# change listen address in postgres config file to localhost + add custom pg_hba.conf
RUN echo listen_addresses='localhost' >> /etc/postgresql/9.6/main/postgresql.conf \
 && cp tmp/imagemonkey-archiver/env/postgres/pg_hba.conf /etc/postgresql/9.6/main/pg_hba.conf

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