version="1.1.0"

build:
	docker build -t quay.io/utilitywarehouse/sys-graylog-configurer:$(version) .

push:
	docker push quay.io/utilitywarehouse/sys-graylog-configurer:$(version)
