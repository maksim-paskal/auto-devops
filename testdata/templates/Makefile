KUBECONFIG=$(HOME)/.kube/config
k8s_provider_microk8s := $(shell microk8s.kubectl version --client>/dev/null 2>&1; echo $$?)

help:
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

auto-devops: ##Loads new devops spec
	@docker pull paskalmaksim/auto-devops:dev
	@docker run -it \
	-u `id -u`:`id -u` \
	-v `pwd`:/dir \
	-w /dir \
	paskalmaksim/auto-devops:dev \
	-bootstrap=https://minio.test.com/auto-devops/auto-devops.zip

testChart: ##Test kubernetes deployment
	helm lint --strict ./chart
	helm template ./chart | kubectl apply --dry-run=client -f -

	helm lint --strict ./chart/db
	helm template ./chart/db | kubectl apply --dry-run=client -f -

	helm lint --strict ./chart/db --values=.gitlab-ci/helm/db.yaml
	helm template ./chart/db --values=.gitlab-ci/helm/db.yaml | kubectl apply --dry-run=client -f -

	helm lint --strict ./chart --values=.gitlab-ci/helm/stage.yaml
	helm template ./chart --values=.gitlab-ci/helm/stage.yaml | kubectl apply --dry-run=client -f -

	helm lint --strict ./chart --values=.gitlab-ci/helm/prod.yaml
	helm template ./chart --values=.gitlab-ci/helm/prod.yaml | kubectl apply --dry-run=client -f -
build: ##Build service
	docker-compose build --pull --parallel
ifeq ($(k8s_provider_microk8s),0)
	docker push localhost:32000/{% .GitInfo.PathFormated %}/front:local
	docker push localhost:32000/{% .GitInfo.PathFormated %}/backend:local

	microk8s.ctr image pull docker-hub-proxy.test.com/library/percona:8.0.25-15
	microk8s.ctr image pull --plain-http localhost:32000/{% .GitInfo.PathFormated %}/front:local
	microk8s.ctr image pull --plain-http localhost:32000/{% .GitInfo.PathFormated %}/backend:local
endif

run:  ##Deploy local environment
	@echo MYIP=`kubectl get no -o jsonpath='{.items[0].status.addresses[0].address}'` > .env.docker
	make stopDocker
	docker-compose --env-file .env.docker up

ip: ##Show IP of PODSs
	@echo "https://`kubectl -n paket get pod -lapp={% .Answers.name %} -o=jsonpath='{.items[0].status.podIP}'`:8001"

ipDocker: ##Show IP of containers
	@echo web=http://localhost:{% .Answers.webServerPort %}
	@echo mysql=127.0.0.1:{% .Answers.mysqlServerPort %}

stopDocker: ##Stop containers
	docker-compose down

shellDocker: ##Run shell into php-container
	docker-compose exec backend bash

shell: ##Run shell into cluster container
	kubectl --namespace paket wait pod -l=app={% .Answers.name %} --for=condition=Ready --timeout=600s
	@kubectl --namespace paket exec -it `kubectl --namespace paket get pods -l=app={% .Answers.name %}  -o jsonpath='{.items[0].metadata.name}'` -c backend bash

clear: ##Remove all containers
	helm uninstall {% .Answers.name %}-db -n paket || true
	helm uninstall {% .Answers.name %} -n paket || true
	docker-compose down --remove-orphans --volumes
	sudo rm -rf ./vendor ./var/cache
	kubectl --namespace paket wait persistentvolumeclaims/{% .Answers.name %}-db-mysql-data --for=delete --timeout=600s || true

deploy: ##Deploy to cluster
	kubectl -n paket scale deploy {% .Answers.name %} --replicas=0 || true

	helm upgrade --install {% .Answers.name %}-db ./chart/db --namespace=paket
	kubectl --namespace=paket wait --for=condition=available deployment {% .Answers.name %}-db-mysql --timeout=600s

	helm upgrade --install --timeout 600s {% .Answers.name %} ./chart --namespace=paket \
	--set devMode.enabled=true \
	--set devMode.appDir=`pwd` \
	--set devMode.homeDir=`echo ~` \
	--set hostIP=`kubectl get no -o jsonpath='{.items[0].status.addresses[0].address}'`

deploy-logs:
	kubectl --namespace paket logs -lapp={% .Answers.name %}-install --tail=1000

pods:
	kubectl --namespace paket get pods --sort-by=.metadata.name -lapp={% .Answers.name %} -o wide

test: ##Run unit tests
	kubectl --namespace paket wait pod -l=app={% .Answers.name %} --for=condition=Ready --timeout=600s
	@kubectl --namespace paket exec -it `kubectl --namespace paket get pods -l=app={% .Answers.name %} -o jsonpath='{.items[0].metadata.name}'` -c backend /app/entrypoint testOnly $(args)

checkCodeStyle: ##Check CodeStyle
	kubectl --namespace paket wait pod -l=app={% .Answers.name %} --for=condition=Ready --timeout=600s
	@kubectl --namespace paket exec -it `kubectl --namespace paket get pods -l=app={% .Answers.name %} -o jsonpath='{.items[0].metadata.name}'` -c backend /app/entrypoint checkCodeStyle $(args)

fixCodeStyle: ##Fix CodeStyle
	kubectl --namespace paket wait pod -l=app={% .Answers.name %} --for=condition=Ready --timeout=600s
	@kubectl --namespace paket exec -it `kubectl --namespace paket get pods -l=app={% .Answers.name %} -o jsonpath='{.items[0].metadata.name}'` -c backend /app/entrypoint fixCodeStyle $(args)

