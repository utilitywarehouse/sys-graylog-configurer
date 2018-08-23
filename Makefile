version="2.4.5-1"

build:
	docker build -t quay.io/utilitywarehouse/sys-graylog:$(version) .

push:
	docker push quay.io/utilitywarehouse/sys-graylog:$(version)
