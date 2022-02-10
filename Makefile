build:
	docker build -t minitwit/webserver --target development . && \
	docker build -t minitwit/tests --target test .

start:
	make build && \
	docker run --detach -p 8080:8080 --name minitwit-webserver minitwit/webserver:latest

stop:
	docker stop minitwit-webserver
	docker stop minitwit-tests

test:
	make build && \
	docker run --name minitwit-tests minitwit/tests:latest

clean:
	./clean.sh

python-init:
	python -c"from minitwit import init_db; init_db()"

python-build:
	gcc flag_tool.c -l sqlite3 -o flag_tool

python-clean:
	rm flag_tool
