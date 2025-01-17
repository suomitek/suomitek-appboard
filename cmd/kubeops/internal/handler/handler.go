package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/suomitek/kubeapps-common/response"
	"github.com/suomitek/suomitek-appboard/pkg/agent"
	"github.com/suomitek/suomitek-appboard/pkg/auth"
	chartUtils "github.com/suomitek/suomitek-appboard/pkg/chart"
	"github.com/suomitek/suomitek-appboard/pkg/chart/helm3to2"
	"github.com/suomitek/suomitek-appboard/pkg/handlerutil"
	"github.com/suomitek/suomitek-appboard/pkg/kube"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
	"helm.sh/helm/v3/pkg/action"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	authHeader     = "Authorization"
	clusterParam   = "cluster"
	namespaceParam = "namespace"
	nameParam      = "releaseName"
	authUserError  = "Unexpected error while configuring authentication"
)

const isV1SupportRequired = false

// This type represents the fact that a regular handler cannot actually be created until we have access to the request,
// because a valid action config (and hence handler config) cannot be created until then.
// If the handler config were a "this" argument instead of an explicit argument, it would be easy to create a handler with a "zero" config.
// This approach practically eliminates that risk; it is much easier to use WithHandlerConfig to create a handler guaranteed to use a valid handler config.
type dependentHandler func(cfg Config, w http.ResponseWriter, req *http.Request, params handlerutil.Params)

// Options represents options that can be created without a bearer token, i.e. once at application startup.
type Options struct {
	ListLimit          int
	Timeout            int64
	UserAgent          string
	KubeappsNamespace  string
	AdditionalClusters kube.AdditionalClustersConfig
}

// Config represents data needed by each handler to be able to create Helm 3 actions.
// It cannot be created without a bearer token, so a new one must be created upon each HTTP request.
type Config struct {
	ActionConfig *action.Configuration
	Options      Options
	ChartClient  chartUtils.Resolver
}

// WithHandlerConfig takes a dependentHandler and creates a regular (WithParams) handler that,
// for every request, will create a handler config for itself.
// Written in a curried fashion for convenient usage; see cmd/kubeops/main.go.
func WithHandlerConfig(storageForDriver agent.StorageForDriver, options Options) func(f dependentHandler) handlerutil.WithParams {
	return func(f dependentHandler) handlerutil.WithParams {
		return func(w http.ResponseWriter, req *http.Request, params handlerutil.Params) {
			// Don't assume the cluster name was in the url for backwards compatability
			// for now.
			cluster, ok := params[clusterParam]
			if !ok {
				cluster = kube.DefaultClusterName
			}
			namespace := params[namespaceParam]
			token := auth.ExtractToken(req.Header.Get(authHeader))

			inClusterConfig, err := rest.InClusterConfig()
			if err != nil {
				log.Errorf("Failed to create in-cluster config: %v", err)
				response.NewErrorResponse(http.StatusInternalServerError, authUserError).Write(w)
				return
			}

			restConfig, err := kube.NewClusterConfig(inClusterConfig, token, cluster, options.AdditionalClusters)
			if err != nil {
				log.Errorf("Failed to create in-cluster config with user token: %v", err)
				response.NewErrorResponse(http.StatusInternalServerError, authUserError).Write(w)
				return
			}
			userKubeClient, err := kubernetes.NewForConfig(restConfig)
			if err != nil {
				log.Errorf("Failed to create kube client with user config: %v", err)
				response.NewErrorResponse(http.StatusInternalServerError, authUserError).Write(w)
				return
			}
			actionConfig, err := agent.NewActionConfig(storageForDriver, restConfig, userKubeClient, namespace)
			if err != nil {
				log.Errorf("Failed to create action config with user client: %v", err)
				response.NewErrorResponse(http.StatusInternalServerError, authUserError).Write(w)
				return
			}

			kubeHandler, err := kube.NewHandler(options.KubeappsNamespace, options.AdditionalClusters)
			if err != nil {
				log.Errorf("Failed to create handler: %v", err)
				response.NewErrorResponse(http.StatusInternalServerError, authUserError).Write(w)
				return
			}

			cfg := Config{
				Options:      options,
				ActionConfig: actionConfig,
				ChartClient:  chartUtils.NewChartClient(kubeHandler, options.KubeappsNamespace, options.UserAgent),
			}
			f(cfg, w, req, params)
		}
	}
}

// AddRouteWith makes it easier to define routes in main.go and avoids code repetition.
func AddRouteWith(
	r *mux.Router,
	withHandlerConfig func(dependentHandler) handlerutil.WithParams,
) func(verb, path string, handler dependentHandler) {
	return func(verb, path string, handler dependentHandler) {
		r.Methods(verb).Path(path).Handler(negroni.New(negroni.Wrap(withHandlerConfig(handler))))
	}
}

func returnForbiddenActions(forbiddenActions []auth.Action, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	body, err := json.Marshal(forbiddenActions)
	if err != nil {
		returnErrMessage(err, w)
		return
	}
	response.NewErrorResponse(http.StatusForbidden, string(body)).Write(w)
}

func returnErrMessage(err error, w http.ResponseWriter) {
	code := handlerutil.ErrorCode(err)
	errMessage := err.Error()
	if code == http.StatusForbidden {
		forbiddenActions := auth.ParseForbiddenActions(errMessage)
		if len(forbiddenActions) > 0 {
			returnForbiddenActions(forbiddenActions, w)
		} else {
			// Unable to parse forbidden actions, return the raw message
			response.NewErrorResponse(code, errMessage).Write(w)
		}
	} else {
		response.NewErrorResponse(code, errMessage).Write(w)
	}
}

// ListReleases list existing releases.
func ListReleases(cfg Config, w http.ResponseWriter, req *http.Request, params handlerutil.Params) {
	apps, err := agent.ListReleases(cfg.ActionConfig, params[namespaceParam], cfg.Options.ListLimit, req.URL.Query().Get("statuses"))
	if err != nil {
		returnErrMessage(err, w)
		return
	}
	response.NewDataResponse(apps).Write(w)
}

// ListAllReleases list all the releases available.
func ListAllReleases(cfg Config, w http.ResponseWriter, req *http.Request, _ handlerutil.Params) {
	ListReleases(cfg, w, req, make(map[string]string))
}

// CreateRelease creates a release.
func CreateRelease(cfg Config, w http.ResponseWriter, req *http.Request, params handlerutil.Params) {
	chartDetails, chartMulti, err := handlerutil.ParseAndGetChart(req, cfg.ChartClient, isV1SupportRequired)
	if err != nil {
		returnErrMessage(err, w)
		return
	}

	ch := chartMulti.Helm3Chart
	releaseName := chartDetails.ReleaseName
	namespace := params[namespaceParam]
	valuesString := chartDetails.Values
	release, err := agent.CreateRelease(cfg.ActionConfig, releaseName, namespace, valuesString, ch, cfg.ChartClient.RegistrySecretsPerDomain())
	if err != nil {
		returnErrMessage(err, w)
		return
	}
	response.NewDataResponse(release).Write(w)
}

// OperateRelease decides which method to call depending on the "action" query param.
func OperateRelease(cfg Config, w http.ResponseWriter, req *http.Request, params handlerutil.Params) {
	switch req.FormValue("action") {
	case "upgrade":
		upgradeRelease(cfg, w, req, params)
	case "rollback":
		rollbackRelease(cfg, w, req, params)
	// TODO: Add "test" case here.
	default:
		// By default, for maintaining compatibility, we call upgrade.
		upgradeRelease(cfg, w, req, params)
	}
}

func upgradeRelease(cfg Config, w http.ResponseWriter, req *http.Request, params handlerutil.Params) {
	releaseName := params[nameParam]
	chartDetails, chartMulti, err := handlerutil.ParseAndGetChart(req, cfg.ChartClient, isV1SupportRequired)
	if err != nil {
		returnErrMessage(err, w)
		return
	}

	ch := chartMulti.Helm3Chart
	rel, err := agent.UpgradeRelease(cfg.ActionConfig, releaseName, chartDetails.Values, ch, cfg.ChartClient.RegistrySecretsPerDomain())
	if err != nil {
		returnErrMessage(err, w)
		return
	}
	compatRelease, err := helm3to2.Convert(*rel)
	if err != nil {
		returnErrMessage(err, w)
		return
	}
	response.NewDataResponse(compatRelease).Write(w)
}

func rollbackRelease(cfg Config, w http.ResponseWriter, req *http.Request, params handlerutil.Params) {
	releaseName := params[nameParam]
	revision := req.FormValue("revision")
	if revision == "" {
		response.NewErrorResponse(http.StatusUnprocessableEntity, "Missing revision to rollback in request").Write(w)
		return
	}
	revisionInt, err := strconv.ParseInt(revision, 10, 32)
	if err != nil {
		returnErrMessage(err, w)
		return
	}
	rel, err := agent.RollbackRelease(cfg.ActionConfig, releaseName, int(revisionInt))
	if err != nil {
		returnErrMessage(err, w)
		return
	}
	compatRelease, err := helm3to2.Convert(*rel)
	if err != nil {
		returnErrMessage(err, w)
		return
	}
	response.NewDataResponse(compatRelease).Write(w)
}

// GetRelease returns a release.
func GetRelease(cfg Config, w http.ResponseWriter, req *http.Request, params handlerutil.Params) {
	// Namespace is already known by the RESTClientGetter.
	releaseName := params[nameParam]
	release, err := agent.GetRelease(cfg.ActionConfig, releaseName)
	if err != nil {
		returnErrMessage(err, w)
		return
	}
	compatRelease, err := helm3to2.Convert(*release)
	if err != nil {
		returnErrMessage(err, w)
		return
	}
	response.NewDataResponse(compatRelease).Write(w)
}

// DeleteRelease deletes a release.
func DeleteRelease(cfg Config, w http.ResponseWriter, req *http.Request, params handlerutil.Params) {
	releaseName := params[nameParam]
	purge := handlerutil.QueryParamIsTruthy("purge", req)
	// Helm 3 has --purge by default; --keep-history in Helm 3 corresponds to omitting --purge in Helm 2.
	// https://stackoverflow.com/a/59210923/2135002
	keepHistory := !purge
	err := agent.DeleteRelease(cfg.ActionConfig, releaseName, keepHistory)
	if err != nil {
		returnErrMessage(err, w)
		return
	}
	w.Header().Set("Status-Code", "200")
	w.Write([]byte("OK"))
}
