ver=0.0.2
appName=test
nameSpace=annqlm
.PHONY: build
build:
	docker build -t $(nameSpace)/$(appName):$(ver) .
	docker image prune -f --filter label=stage=builder
	docker push $(nameSpace)/$(appName):$(ver)
	docker tag $(nameSpace)/$(appName):$(ver) $(nameSpace)/$(appName):latest
	docker push $(nameSpace)/$(appName):latest