.PHONY: all create start sysd-start stop restart attach

PROJECTNAME := pingme
RUNARGS := -p 1025:1025

default: create

create:
	docker build -t $(PROJECTNAME) .

start:
	docker run -it --name $(PROJECTNAME) $(RUNARGS) -d $(PROJECTNAME)

stop:
	docker kill $(PROJECTNAME); \
	docker rm $(PROJECTNAME)

restart: stop start

attach:
	docker exec -i -t $(PROJECTNAME) bash

