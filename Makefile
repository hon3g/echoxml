NAMESPACE = echoxml

clean:
	kubectl delete namespace $(NAMESPACE) --ignore-not-found true
	kubectl create namespace $(NAMESPACE)
	kubectl config set-context --current --namespace $(NAMESPACE)

up: clean
	helm install echoxml helm/echoxml --wait
	kubectl get pods

forward-port: close-port
	nohup kubectl port-forward service/echoxml 8080:8080 > /dev/null &

close-port:
	pkill -f '^kubectl port-forward' || :
