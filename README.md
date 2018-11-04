# BUILD

* clone repository with `https://github.com/bbernhard/imagemonkey-archiver.git`
* build docker image with `docker build -t imagemonkey-archiver -f env/docker/Dockerfile .`

# RUN

* `cd` into the directory where the database dump (`db_dump.sql`) and the donations folder (`donations`) resides
* run `docker run --mount type=bind,source="$(pwd)",target=/home/imagemonkey/data -it imagemonkey-archiver`