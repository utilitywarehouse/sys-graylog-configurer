version="1.2.0"

build:
	docker build -t quay.io/utilitywarehouse/sys-graylog-configurer:$(version) .

push:
	docker push quay.io/utilitywarehouse/sys-graylog-configurer:$(version)
