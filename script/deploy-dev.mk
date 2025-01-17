# Deploy a dev environment of Suomitek-appboard using OIDC for authentication with a
# local dex as the provider.
#
# Targets in this helper assume that kubectl is configured with a cluster
# that has been setup with OIDC support (see ./cluster-kind.mk)

deploy-dex: devel/dex.crt devel/dex.key
	kubectl create namespace dex
	kubectl -n dex create secret tls dex-web-server-tls \
		--key ./devel/dex.key \
		--cert ./devel/dex.crt
	helm install dex stable/dex --namespace dex --version 2.13.0 \
		--values ./docs/user/manifests/suomitek-appboard-local-dev-dex-values.yaml

deploy-openldap:
	kubectl create namespace ldap
	helm install ldap stable/openldap --namespace ldap \
		--values ./docs/user/manifests/suomitek-appboard-local-dev-openldap-values.yaml

devel/localhost-cert.pem:
	mkcert -key-file ./devel/localhost-key.pem -cert-file ./devel/localhost-cert.pem localhost 172.18.0.2

deploy-dependencies: deploy-dex deploy-openldap devel/localhost-cert.pem
	kubectl create namespace suomitek-appboard
	kubectl -n suomitek-appboard create secret tls localhost-tls \
		--key ./devel/localhost-key.pem \
		--cert ./devel/localhost-cert.pem

deploy-dev: deploy-dependencies
	helm install suomitek-appboard ./chart/suomitek-appboard --namespace suomitek-appboard \
		--values ./docs/user/manifests/suomitek-appboard-local-dev-values.yaml \
		--values ./docs/user/manifests/suomitek-appboard-local-dev-auth-proxy-values.yaml \
		--values ./docs/user/manifests/suomitek-appboard-local-dev-additional-kind-cluster.yaml \
		--set useHelm3=true
	@echo "\nYou can now simply open your browser at https://localhost/ to access Suomitek-appboard!"
	@echo "When logging in, you will be redirected to dex (with a self-signed cert) and can login with email as either of"
	@echo "  suomitek-appboard-operator@example.com:password"
	@echo "  suomitek-appboard-user@example.com:password"
	@echo "or with LDAP as either of"
	@echo "  suomitek-appboard-operator-ldap@example.org:password"
	@echo "  suomitek-appboard-user-ldap@example.org:password"
	@echo "to authenticate with the corresponding permissions."

reset-dev:
	helm -n suomitek-appboard delete suomitek-appboard || true
	helm -n dex delete dex || true
	helm -n ldap delete ldap || true
	kubectl delete namespace --wait dex ldap suomitek-appboard || true

.PHONY: deploy-dex deploy-dependencies deploy-dev deploy-openldap reset-dev update-apiserver-etc-hosts
