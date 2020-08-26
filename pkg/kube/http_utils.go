/*
Copyright (c) 2020 Bitnami

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

package kube

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/kubeapps/kubeapps/cmd/apprepository-controller/pkg/apis/apprepository/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

const (
	defaultTimeoutSeconds = 180
)

// HTTPClient Interface to perform HTTP requests
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// clientWithDefaultHeaders implements chart.HTTPClient interface
// and includes an override of the Do method which injects our default
// headers - User-Agent and Authorization (when present)
type clientWithDefaultHeaders struct {
	client         HTTPClient
	defaultHeaders http.Header
}

// Do HTTP request
func (c *clientWithDefaultHeaders) Do(req *http.Request) (*http.Response, error) {
	for k, v := range c.defaultHeaders {
		// Only add the default header if it's not already set in the request.
		if _, ok := req.Header[k]; !ok {
			req.Header[k] = v
		}
	}
	return c.client.Do(req)
}

// InitNetClient returns an HTTP client based on the chart details loading a
// custom CA if provided (as a secret)
func InitNetClient(appRepo *v1alpha1.AppRepository, caCertSecret, authSecret *corev1.Secret, defaultHeaders http.Header) (HTTPClient, error) {
	// Require the SystemCertPool unless the env var is explicitly set.
	caCertPool, err := x509.SystemCertPool()
	if err != nil {
		if _, ok := os.LookupEnv("TILLER_PROXY_ALLOW_EMPTY_CERT_POOL"); !ok {
			return nil, err
		}
		caCertPool = x509.NewCertPool()
	}

	if caCertSecret != nil && appRepo.Spec.Auth.CustomCA != nil {
		// Append our cert to the system pool
		key := appRepo.Spec.Auth.CustomCA.SecretKeyRef.Key
		customData, ok := caCertSecret.Data[key]
		if !ok {
			customDataString, ok := caCertSecret.StringData[key]
			if !ok {
				return nil, fmt.Errorf("secret %q did not contain key %q", appRepo.Spec.Auth.CustomCA.SecretKeyRef.Name, key)
			}
			customData = []byte(customDataString)
		}
		if ok := caCertPool.AppendCertsFromPEM(customData); !ok {
			return nil, fmt.Errorf("Failed to append %s to RootCAs", appRepo.Spec.Auth.CustomCA.SecretKeyRef.Name)
		}
	}

	if defaultHeaders == nil {
		defaultHeaders = http.Header{}
	}
	if authSecret != nil && appRepo.Spec.Auth.Header != nil {
		key := appRepo.Spec.Auth.Header.SecretKeyRef.Key
		auth, ok := authSecret.StringData[key]
		if !ok {
			authBytes, ok := authSecret.Data[key]
			if !ok {
				return nil, fmt.Errorf("secret %q did not contain key %q", appRepo.Spec.Auth.Header.SecretKeyRef.Name, key)
			}
			auth = string(authBytes)
		}
		defaultHeaders.Set("Authorization", string(auth))
	}

	// Return Transport for testing purposes
	return &clientWithDefaultHeaders{
		client: &http.Client{
			Timeout: time.Second * defaultTimeoutSeconds,
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				TLSClientConfig: &tls.Config{
					RootCAs: caCertPool,
				},
			},
		},
		defaultHeaders: defaultHeaders,
	}, nil
}
