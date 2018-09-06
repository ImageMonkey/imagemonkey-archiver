import signal
import subprocess
import os
import sys

def signal_handler(sig, frame):
        log.info('You pressed Ctrl+C! - Stopping container')
        if container is not None:
        	container.stop()
        sys.exit(1)
signal.signal(signal.SIGINT, signal_handler)

if __name__ == "__main__":
	print("Archiver started")

	donations_folder = "/home/imagemonkey/data/donations"
	dbdump = "/home/imagemonkey/data/db_dump.sql"
	if not os.path.isdir(donations_folder):
		print("donations folder doesn't exist. Have you mounted it?")
		sys.exit(1)

	if not os.path.isfile(dbdump):
		print("database dump doesn't exist. Have you mounted it?")
		sys.exit(1)


	print("Start PostgreSQL")
	cmd = "bash /home/imagemonkey/bin/start_postgres.sh"

	p = subprocess.Popen(cmd, stdout=subprocess.PIPE, shell=True)
	for line in iter(p.stdout.readline, b''):  
		sys.stdout.write(line.decode('utf-8'))


	output_folder = "/home/imagemonkey/data/"
	cmd = ('/home/imagemonkey/bin/archiver -donationsdir="' 
			+ donations_folder + '" -dbdump="' + dbdump 
			+ '" -output="' + output_folder + '" -verify=true -dryrun=false')
	#cmd = ('/home/imagemonkey/bin/archiver -donationsdir="' 
	#		+ donations_folder + '" -dbdump="' + dbdump 
	#		+ '" -dbpasswd="imagemonkey" -output="' + output_folder + '"')

	p = subprocess.Popen(cmd, stdout=subprocess.PIPE, shell=True)
	for line in iter(p.stdout.readline, b''):  
		sys.stdout.write(line)

	#im test folder, run:
	#docker run --mount type=bind,source="$(pwd)",target=/home/imagemonkey/data -it imagemonkey-archiver
