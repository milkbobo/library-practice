watch:
	-pkill -9 fishcmd
	-pkill -9 library
	nohup fishcmd run --watch &
	tail -f nohup.out
release:
	go install goldenstatue
	rm -rf goldenstatue
	cp ${GOPATH}/bin/goldenstatue .
	sudo supervisorctl restart goldenstatue_server
	sudo tail -f ../data/log/access.log
stop:
	sudo supervisorctl stop goldenstatue_server 
	-pkill -9 bee
	-pkill -9 goldenstatue
