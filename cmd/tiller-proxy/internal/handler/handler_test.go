/*
Copyright (c) 2018 Bitnami

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/proto/hapi/release"

	"github.com/kubeapps/kubeapps/pkg/auth"
	authFake "github.com/kubeapps/kubeapps/pkg/auth/fake"
	chartFake "github.com/kubeapps/kubeapps/pkg/chart/fake"
	proxyFake "github.com/kubeapps/kubeapps/pkg/proxy/fake"
)

func TestActions(t *testing.T) {
	type testScenario struct {
		// Scenario params
		Description      string
		ExistingReleases []release.Release
		ForbiddenActions []auth.Action
		// Request params
		RequestBody  string
		RequestQuery string
		Action       string
		Params       map[string]string
		// Expected result
		StatusCode        int
		RemainingReleases []release.Release
		ResponseBody      string // Optional
	}
	tests := []testScenario{
		{
			// Scenario params
			Description:      "Create a simple release with auth",
			ExistingReleases: []release.Release{},
			ForbiddenActions: []auth.Action{},
			// Request params
			RequestBody: `{"chartName": "foo", "releaseName": "foobar",	"version": "1.0.0"}`,
			RequestQuery: "",
			Action:       "create",
			Params:       map[string]string{"namespace": "default"},
			// Expected result
			StatusCode: 200,
			RemainingReleases: []release.Release{
				release.Release{Name: "foobar", Namespace: "default"},
			},
			ResponseBody: "",
		},
		{
			// Scenario params
			Description:      "Create a conflicting release",
			ExistingReleases: []release.Release{release.Release{Name: "foobar", Namespace: "default"}},
			ForbiddenActions: []auth.Action{},
			// Request params
			RequestBody: `{"chartName": "foo", "releaseName": "foobar",	"version": "1.0.0"}`,
			RequestQuery: "",
			Action:       "create",
			Params:       map[string]string{"namespace": "default"},
			// Expected result
			StatusCode: 409,
			RemainingReleases: []release.Release{
				release.Release{Name: "foobar", Namespace: "default"},
			},
			ResponseBody: "",
		},
		{
			// Scenario params
			Description:      "Create a simple release with forbidden actions",
			ExistingReleases: []release.Release{},
			ForbiddenActions: []auth.Action{
				auth.Action{APIVersion: "v1", Resource: "pods", Namespace: "default", ClusterWide: false, Verbs: []string{"create"}},
			},
			// Request params
			RequestBody: `{"chartName": "foo", "releaseName": "foobar",	"version": "1.0.0"}`,
			RequestQuery: "",
			Action:       "create",
			Params:       map[string]string{"namespace": "default"},
			// Expected result
			StatusCode:        403,
			RemainingReleases: []release.Release{},
			ResponseBody:      `{"code":403,"message":"[{\"apiGroup\":\"v1\",\"resource\":\"pods\",\"namespace\":\"default\",\"clusterWide\":false,\"verbs\":[\"create\"]}]"}`,
		},
		{
			// Scenario params
			Description:      "Upgrade a simple release",
			ExistingReleases: []release.Release{release.Release{Name: "foobar", Namespace: "default"}},
			ForbiddenActions: []auth.Action{},
			// Request params
			RequestBody: `{"chartName": "foo", "releaseName": "foobar",	"version": "1.0.0"}`,
			RequestQuery: "",
			Action:       "upgrade",
			Params:       map[string]string{"namespace": "default", "releaseName": "foobar"},
			// Expected result
			StatusCode:        200,
			RemainingReleases: []release.Release{release.Release{Name: "foobar", Namespace: "default"}},
			ResponseBody:      "",
		},
		{
			// Scenario params
			Description:      "Upgrade a missing release",
			ExistingReleases: []release.Release{},
			ForbiddenActions: []auth.Action{},
			// Request params
			RequestBody: `{"chartName": "foo", "releaseName": "foobar",	"version": "1.0.0"}`,
			RequestQuery: "",
			Action:       "upgrade",
			Params:       map[string]string{"namespace": "default", "releaseName": "foobar"},
			// Expected result
			StatusCode:        404,
			RemainingReleases: []release.Release{},
			ResponseBody:      "",
		},
		{
			// Scenario params
			Description:      "Upgrade a simple release with forbidden actions",
			ExistingReleases: []release.Release{release.Release{Name: "foobar", Namespace: "default"}},
			ForbiddenActions: []auth.Action{
				auth.Action{APIVersion: "v1", Resource: "pods", Namespace: "default", ClusterWide: false, Verbs: []string{"upgrade"}},
			},
			// Request params
			RequestBody: `{"chartName": "foo", "releaseName": "foobar",	"version": "1.0.0"}`,
			RequestQuery: "",
			Action:       "upgrade",
			Params:       map[string]string{"namespace": "default"},
			// Expected result
			StatusCode:        403,
			RemainingReleases: []release.Release{release.Release{Name: "foobar", Namespace: "default"}},
			ResponseBody:      `{"code":403,"message":"[{\"apiGroup\":\"v1\",\"resource\":\"pods\",\"namespace\":\"default\",\"clusterWide\":false,\"verbs\":[\"upgrade\"]}]"}`,
		},
		{
			// Scenario params
			Description:      "Delete a simple release",
			ExistingReleases: []release.Release{release.Release{Name: "foobar", Namespace: "default", Config: &chart.Config{Raw: ""}}},
			ForbiddenActions: []auth.Action{},
			// Request params
			RequestBody:  "",
			RequestQuery: "",
			Action:       "delete",
			Params:       map[string]string{"namespace": "default", "releaseName": "foobar"},
			// Expected result
			StatusCode:        200,
			RemainingReleases: []release.Release{release.Release{Name: "foobar", Namespace: "default", Config: &chart.Config{Raw: ""}, Info: &release.Info{Status: &release.Status{Code: release.Status_DELETED}}}},
			ResponseBody:      "",
		},
		{
			// Scenario params
			Description:      "Delete and purge a simple release",
			ExistingReleases: []release.Release{release.Release{Name: "foobar", Namespace: "default", Config: &chart.Config{Raw: ""}}},
			ForbiddenActions: []auth.Action{},
			// Request params
			RequestBody:  "",
			RequestQuery: "?purge=true",
			Action:       "delete",
			Params:       map[string]string{"namespace": "default", "releaseName": "foobar"},
			// Expected result
			StatusCode:        200,
			RemainingReleases: []release.Release{},
			ResponseBody:      "",
		},
		{
			// Scenario params
			Description:      "Delete a missing release",
			ExistingReleases: []release.Release{},
			ForbiddenActions: []auth.Action{},
			// Request params
			RequestBody: `{"chartName": "foo", "releaseName": "foobar",	"version": "1.0.0"}`,
			RequestQuery: "",
			Action:       "delete",
			Params:       map[string]string{"namespace": "default", "releaseName": "foobar"},
			// Expected result
			StatusCode:        404,
			RemainingReleases: []release.Release{},
			ResponseBody:      "",
		},
		{
			// Scenario params
			Description:      "Delete a release with forbidden actions",
			ExistingReleases: []release.Release{release.Release{Name: "foobar", Namespace: "default", Config: &chart.Config{Raw: ""}}},
			ForbiddenActions: []auth.Action{
				auth.Action{APIVersion: "v1", Resource: "pods", Namespace: "default", ClusterWide: false, Verbs: []string{"delete"}},
			},
			// Request params
			RequestBody: `{"chartName": "foo", "releaseName": "foobar",	"version": "1.0.0"}`,
			RequestQuery: "",
			Action:       "delete",
			Params:       map[string]string{"namespace": "default", "releaseName": "foobar"},
			// Expected result
			StatusCode:        403,
			RemainingReleases: []release.Release{release.Release{Name: "foobar", Namespace: "default", Config: &chart.Config{Raw: ""}}},
			ResponseBody:      `{"code":403,"message":"[{\"apiGroup\":\"v1\",\"resource\":\"pods\",\"namespace\":\"default\",\"clusterWide\":false,\"verbs\":[\"delete\"]}]"}`,
		},
		{
			// Scenario params
			Description:      "Get a simple release",
			ExistingReleases: []release.Release{release.Release{Name: "foobar", Namespace: "default", Config: &chart.Config{Raw: ""}}},
			ForbiddenActions: []auth.Action{},
			// Request params
			RequestBody:  "",
			RequestQuery: "",
			Action:       "get",
			Params:       map[string]string{"namespace": "default", "releaseName": "foobar"},
			// Expected result
			StatusCode:        200,
			RemainingReleases: []release.Release{release.Release{Name: "foobar", Namespace: "default", Config: &chart.Config{Raw: ""}}},
			ResponseBody:      `{"data":{"name":"foobar","config":{},"namespace":"default"}}`,
		},
		{
			// Scenario params
			Description:      "Get a missing release",
			ExistingReleases: []release.Release{},
			ForbiddenActions: []auth.Action{},
			// Request params
			RequestBody:  "",
			RequestQuery: "",
			Action:       "get",
			Params:       map[string]string{"namespace": "default", "releaseName": "foobar"},
			// Expected result
			StatusCode:        404,
			RemainingReleases: []release.Release{},
			ResponseBody:      "",
		},
		{
			// Scenario params
			Description:      "Get a release with forbidden actions",
			ExistingReleases: []release.Release{release.Release{Name: "foobar", Namespace: "default", Config: &chart.Config{Raw: ""}}},
			ForbiddenActions: []auth.Action{
				auth.Action{APIVersion: "v1", Resource: "pods", Namespace: "default", ClusterWide: false, Verbs: []string{"get"}},
			},
			// Request params
			RequestBody:  "",
			RequestQuery: "",
			Action:       "get",
			Params:       map[string]string{"namespace": "default", "releaseName": "foobar"},
			// Expected result
			StatusCode:        403,
			RemainingReleases: []release.Release{release.Release{Name: "foobar", Namespace: "default", Config: &chart.Config{Raw: ""}}},
			ResponseBody:      `{"code":403,"message":"[{\"apiGroup\":\"v1\",\"resource\":\"pods\",\"namespace\":\"default\",\"clusterWide\":false,\"verbs\":[\"get\"]}]"}`,
		},
		{
			// Scenario params
			Description: "List all releases",
			ExistingReleases: []release.Release{
				release.Release{Name: "foobar", Namespace: "default"},
				release.Release{Name: "foo", Namespace: "not-default"},
			},
			ForbiddenActions: []auth.Action{},
			// Request params
			RequestBody:  "",
			RequestQuery: "",
			Action:       "listall",
			Params:       map[string]string{},
			// Expected result
			StatusCode: 200,
			RemainingReleases: []release.Release{
				release.Release{Name: "foobar", Namespace: "default"},
				release.Release{Name: "foo", Namespace: "not-default"},
			},
			ResponseBody: `{"data":[{"releaseName":"foobar","version":"","namespace":"default","status":"DEPLOYED","chart":"","chartMetadata":{}},{"releaseName":"foo","version":"","namespace":"not-default","status":"DEPLOYED","chart":"","chartMetadata":{}}]}`,
		},
		{
			// Scenario params
			Description: "List releases in a namespace",
			ExistingReleases: []release.Release{
				release.Release{Name: "foobar", Namespace: "default"},
				release.Release{Name: "foo", Namespace: "not-default"},
			},
			ForbiddenActions: []auth.Action{},
			// Request params
			RequestBody:  "",
			RequestQuery: "",
			Action:       "list",
			Params:       map[string]string{"namespace": "default"},
			// Expected result
			StatusCode: 200,
			RemainingReleases: []release.Release{
				release.Release{Name: "foobar", Namespace: "default"},
				release.Release{Name: "foo", Namespace: "not-default"},
			},
			ResponseBody: `{"data":[{"releaseName":"foobar","version":"","namespace":"default","status":"DEPLOYED","chart":"","chartMetadata":{}}]}`,
		},
		{
			// Scenario params
			Description: "Filter releases based on status when listing",
			ExistingReleases: []release.Release{
				release.Release{Name: "foobar", Namespace: "default", Info: &release.Info{Status: &release.Status{Code: release.Status_DEPLOYED}}},
				release.Release{Name: "foo", Namespace: "default", Info: &release.Info{Status: &release.Status{Code: release.Status_DELETED}}},
			},
			ForbiddenActions: []auth.Action{},
			// Request params
			RequestBody:  "",
			RequestQuery: "?statuses=deployed",
			Action:       "list",
			Params:       map[string]string{"namespace": "default"},
			// Expected result
			StatusCode: 200,
			RemainingReleases: []release.Release{
				release.Release{Name: "foobar", Namespace: "default", Info: &release.Info{Status: &release.Status{Code: release.Status_DEPLOYED}}},
				release.Release{Name: "foo", Namespace: "default", Info: &release.Info{Status: &release.Status{Code: release.Status_DELETED}}},
			},
			ResponseBody: `{"data":[{"releaseName":"foobar","version":"","namespace":"default","status":"DEPLOYED","chart":"","chartMetadata":{}}]}`,
		},
		{
			// Scenario params
			Description: "Rolls back a release",
			ExistingReleases: []release.Release{
				release.Release{Name: "foo", Namespace: "default", Info: &release.Info{Status: &release.Status{Code: release.Status_DEPLOYED}}},
			},
			ForbiddenActions: []auth.Action{},
			// Request params
			RequestBody:  "",
			RequestQuery: "?revision=1",
			Action:       "rollback",
			Params:       map[string]string{"namespace": "default", "releaseName": "foo"},
			// Expected result
			StatusCode: 200,
			RemainingReleases: []release.Release{
				release.Release{Name: "foo", Namespace: "default", Info: &release.Info{Status: &release.Status{Code: release.Status_DEPLOYED}}},
			},
			ResponseBody: `{"data":{"name":"foo","info":{"status":{"code":1}},"namespace":"default"}}`,
		},
		{
			// Scenario params
			Description:      "Rollsback a missing release",
			ExistingReleases: []release.Release{},
			ForbiddenActions: []auth.Action{},
			// Request params
			RequestBody:  "",
			RequestQuery: "?revision=1",
			Action:       "rollback",
			Params:       map[string]string{"namespace": "default", "releaseName": "foobar"},
			// Expected result
			StatusCode:        404,
			RemainingReleases: []release.Release{},
			ResponseBody:      "",
		},
		{
			// Scenario params
			Description:      "Rollback without a revision",
			ExistingReleases: []release.Release{},
			ForbiddenActions: []auth.Action{},
			// Request params
			RequestBody:  "",
			RequestQuery: "",
			Action:       "rollback",
			Params:       map[string]string{"namespace": "default", "releaseName": "foobar"},
			// Expected result
			StatusCode:        422,
			RemainingReleases: []release.Release{},
			ResponseBody:      "",
		},
		{
			// Scenario params
			Description:      "Test a release successfully",
			ExistingReleases: []release.Release{release.Release{Name: "kubeapps", Namespace: "kubeapps-ns"}},
			ForbiddenActions: []auth.Action{},
			// Request params
			RequestBody:  "",
			RequestQuery: "",
			Action:       "test",
			Params:       map[string]string{"namespace": "kubeapps-ns", "releaseName": "kubeapps"},
			// Expected result
			StatusCode:        200,
			RemainingReleases: []release.Release{release.Release{Name: "kubeapps", Namespace: "kubeapps-ns"}},
			ResponseBody:      `{"data":{"UNKNOWN":["No Tests Found"]}}`,
		},
		{
			// Scenario params
			Description:      "Fail to test a release",
			ExistingReleases: []release.Release{release.Release{Name: "kubeapps", Namespace: "kubeapps-ns"}},
			ForbiddenActions: []auth.Action{},
			// Request params
			RequestBody:  "",
			RequestQuery: "",
			Action:       "test",
			Params:       map[string]string{"namespace": "default", "releaseName": "kubeapps"},
			// Expected result
			StatusCode:        404,
			RemainingReleases: []release.Release{release.Release{Name: "kubeapps", Namespace: "kubeapps-ns"}},
			ResponseBody:      `{"code":404,"message":"Unable to locate release: Release kubeapps not found"}`,
		},
	}
	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			// Prepare environment
			proxy := &proxyFake.FakeProxy{
				Releases: test.ExistingReleases,
			}
			handler := TillerProxy{
				ListLimit:   255,
				ChartClient: &chartFake.FakeChart{},
				ProxyClient: proxy,
			}
			req := httptest.NewRequest("GET", fmt.Sprintf("http://foo.bar%s", test.RequestQuery), strings.NewReader(test.RequestBody))
			handler.CheckerForRequest = func(req *http.Request) (auth.Checker, error) {
				return &authFake.FakeAuth{
					ForbiddenActions: test.ForbiddenActions,
				}, nil
			}
			response := httptest.NewRecorder()
			switch test.Action {
			case "create":
				handler.CreateRelease(response, req, test.Params)
			case "upgrade":
				handler.UpgradeRelease(response, req, test.Params)
			case "delete":
				handler.DeleteRelease(response, req, test.Params)
			case "get":
				handler.GetRelease(response, req, test.Params)
			case "rollback":
				handler.RollbackRelease(response, req, test.Params)
			case "list":
				handler.ListReleases(response, req, test.Params)
			case "listall":
				handler.ListAllReleases(response, req)
			case "test":
				handler.TestRelease(response, req, test.Params)
			default:
				t.Errorf("Unexpected action %s", test.Action)
			}
			// Check result
			if response.Code != test.StatusCode {
				t.Errorf("Expecting a StatusCode %d, received %d", test.StatusCode, response.Code)
			}
			if got, want := proxy.Releases, test.RemainingReleases; !cmp.Equal(want, got) {
				t.Errorf("mismatch (-want +got):\n%s", cmp.Diff(want, got))
			}
			if test.ResponseBody != "" {
				if test.ResponseBody != response.Body.String() {
					t.Errorf("Unexpected body response. Expecting %s, found %s", test.ResponseBody, response.Body)
				}
			}
		})
	}
}
