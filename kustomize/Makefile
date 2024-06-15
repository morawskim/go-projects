.PHONY: kustomize build kubectl-apply list-cron-jobs

build:
	CGO_ENABLED=0 go build -o kustomize/plugins/mykrmplugin krm.go

kustomize:
	KUSTOMIZE_PLUGIN_HOME=$(PWD)/kustomize/plugins kustomize build  \
		--enable-alpha-plugins --enable-exec ./kustomize/overlay/test

kubectl-apply:
	KUSTOMIZE_PLUGIN_HOME=$(PWD)/kustomize/plugins kustomize build  \
		--enable-alpha-plugins --enable-exec ./kustomize/overlay/test | kubectl apply -f -

list-cron-jobs:
	kubectl get cj -n kustomize-plugin
