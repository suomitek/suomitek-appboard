/*
Copyright (c) Bitnami

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

package dbutils

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"github.com/kubeapps/common/datastore"
	"github.com/kubeapps/kubeapps/pkg/chart/models"
)

func Test_NewPGManager(t *testing.T) {
	config := datastore.Config{URL: "10.11.12.13:5432", Database: "assets", Username: "postgres", Password: "123"}
	m, err := NewPGManager(config, "kubeapps")
	if err != nil {
		t.Errorf("Found error %v", err)
	}
	expectedConnStr := "host=10.11.12.13 port=5432 user=postgres password=123 dbname=assets sslmode=disable"
	if m.connStr != expectedConnStr {
		t.Errorf("Expected %s got %s", expectedConnStr, m.connStr)
	}
}

func Test_Close(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	manager := PostgresAssetManager{
		connStr: "localhost",
		DB:      db,
	}
	mock.ExpectClose()
	err = manager.Close()
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func Test_QueryOne(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	manager := PostgresAssetManager{
		connStr: "localhost",
		DB:      db,
	}
	query := "SELECT * from charts"
	rows := sqlmock.NewRows([]string{"info"}).AddRow(`{"ID": "foo"}`)
	mock.ExpectQuery("^SELECT (.+)$").WillReturnRows(rows)
	target := models.Chart{}
	err = manager.QueryOne(&target, query)
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	expectedChart := models.Chart{ID: "foo"}
	if !cmp.Equal(target, expectedChart) {
		t.Errorf("Unexpected result %v", cmp.Diff(target, expectedChart))
	}
}

func Test_QueryAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	manager := PostgresAssetManager{
		connStr: "localhost",
		DB:      db,
	}
	query := "SELECT * from charts"
	rows := sqlmock.NewRows([]string{"info"}).
		AddRow(`{"ID": "foo"}`).
		AddRow(`{"ID": "bar"}`)
	mock.ExpectQuery("^SELECT (.+)$").WillReturnRows(rows)
	charts, err := manager.QueryAllCharts(query)
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	expectedCharts := []*models.Chart{&models.Chart{ID: "foo"}, &models.Chart{ID: "bar"}}
	if !cmp.Equal(charts, expectedCharts) {
		t.Errorf("Unexpected result %v", cmp.Diff(charts, expectedCharts))
	}
}
