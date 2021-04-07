sync_files:
	scp -i "m-vas.pem" -r ./* ubuntu@54.217.120.10:/home/ubuntu/gideon/machama

rsync:
	rsync -a -e "ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -i "m-vas.pem"" ./* ubuntu@54.217.120.10:/home/ubuntu/gideon/machama