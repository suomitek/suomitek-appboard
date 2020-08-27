#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

# Constants
ROOT_DIR="$(cd "$( dirname "${BASH_SOURCE[0]}" )/.." >/dev/null && pwd)"
RESET='\033[0m'
GREEN='\033[38;5;2m'
RED='\033[38;5;1m'
YELLOW='\033[38;5;3m'

# Load Libraries
# shellcheck disable=SC1090
. "${ROOT_DIR}/script/libtest.sh"
# shellcheck disable=SC1090
. "${ROOT_DIR}/script/liblog.sh"

# Axiliar functions
print_menu() {
    local script
    script=$(basename "${BASH_SOURCE[0]}")
    log "${RED}NAME${RESET}"
    log "    $(basename -s .sh "${BASH_SOURCE[0]}")"
    log ""
    log "${RED}SYNOPSIS${RESET}"
    log "    $script [${YELLOW}-dh${RESET}] [${YELLOW}-n ${GREEN}\"namespace\"${RESET}] [${YELLOW}--initial-repos ${GREEN}\"name\" \"url\"${RESET}]"
    log ""
    log "${RED}DESCRIPTION${RESET}"
    log "    Script to setup Suomitek-appboard on your K8s cluster."
    log ""
    log "    The options are as follow:"
    log ""
    log "      ${YELLOW}-n, --namespace ${GREEN}[namespace]${RESET}           Namespace to use for Suomitek-appboard."
    log "      ${YELLOW}--initial-repos ${GREEN}[repo_name repo_url]${RESET}   Initial repositories to configure on Suomitek-appboard. This flag can be used several times."
    log "      ${YELLOW}-h, --help${RESET}                            Print this help menu."
    log "      ${YELLOW}-u, --dry-run${RESET}                         Enable \"dry run\" mode."
    log ""
    log "${RED}EXAMPLES${RESET}"
    log "      $script --help"
    log "      $script --namespace \"suomitek-appboard\""
    log "      $script --namespace \"suomitek-appboard\" --initial-repos \"harbor-library\" \"http://harbor.harbor.svc.cluster.local/chartrepo/library\""
    log ""
}

namespace="suomitek-appboard"
initial_repos=("bitnami http://helm.yongchehang.com")
help_menu=0
dry_run=0
while [[ "$#" -gt 0 ]]; do
    case "$1" in
        -h|--help)
            help_menu=1
            ;;
        -u|--dry-run)
            dry_run=1
            ;;
        --initial-repos)
            shift; repo_name="${1:?missing repo name}"
            shift; repo_url="${1:?missing repo url}"
            initial_repos=("${initial_repos[@]}" "$repo_name $repo_url")
            ;;
        -n|--namespace)
            shift; namespace="${1:?missing namespace}"
            ;;
        *)
            error "Invalid command line flag $1" >&2
            exit 1
            ;;
    esac
    shift
done

if [[ "$help_menu" -eq 1 ]]; then
    print_menu
    exit 0
fi

# Suomitek-appboard values
values="$(cat << EOF
useHelm3: true
apprepository:
  initialRepos:
EOF
)"
for repo in "${initial_repos[@]}"; do
    values="$(cat << EOF
$values
    - name: $(echo "$repo" | awk '{print $1}')
      url: $(echo "$repo" | awk '{print $2}')
EOF
    )"
done

if [[ "$dry_run" -eq 1 ]]; then
    info "DRY RUN mode enabled!"
    info "Namespace: $namespace"
    info "Generated values.yaml:"
    printf '#####\n\n%s\n\n#####\n' "$values"
    exit 0
fi

# Install Suomitek-appboard
info "Using the values.yaml below:"
printf '#####\n\n%s\n\n#####\n' "$values"
info "Installing Suomitek-appboard in namespace '$namespace'..."
silence kubectl create ns "$namespace"
silence helm install suomitek-appboard \
    --namespace "$namespace" \
    -f <(echo "$values") \
    chartmuseum/suomitek-appboard
# Wait for Suomitek-appboard components
info "Waiting for Suomitek-appboard components to be ready..."
deployments=(
    "suomitek-appboard"
    "suomitek-appboard-internal-apprepository-controller"
    "suomitek-appboard-internal-assetsvc"
    "suomitek-appboard-internal-dashboard"
)

for dep in "${deployments[@]}"; do
    k8s_wait_for_deployment "$namespace" "$dep"
    info "Deployment ${dep} ready!"
done
echo

# Create serviceAccount
info "Creating 'example' serviceAccount and adding RBAC permissions for 'default' namespace..."
silence kubectl create serviceaccount example --namespace default
silence kubectl apply -f https://raw.githubusercontent.com/suomitek/suomitek-appboard/master/docs/user/manifests/suomitek-appboard-applications-read.yaml
silence kubectl create -n default rolebinding example-view --clusterrole=suomitek-appboard-applications-read --serviceaccount default:example
silence kubectl create -n default rolebinding example-edit --clusterrole=edit --serviceaccount default:example
silence kubectl create -n "$namespace" rolebinding example-suomitek-appboard-repositories-read --role=suomitek-appboard-repositories-read --serviceaccount default:example
silence kubectl create -n "$namespace" rolebinding example-suomitek-appboard-repositories-write --role=suomitek-appboard-repositories-write --serviceaccount default:example
echo
    
info "Use this command for port forwading to Suomitek-appboard Dashboard:"
info "kubectl port-forward --namespace $namespace svc/suomitek-appboard 8080:80 >/dev/null 2>&1 &"
info "Suomitek-appboard URL: http://127.0.0.1:8080"
info "Kubeppas API Token:"
kubectl get -n default secret "$(kubectl get serviceaccount example --namespace default -o jsonpath='{.secrets[].name}')" -o go-template='{{.data.token | base64decode}}' && echo
echo
