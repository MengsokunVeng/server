// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"reflect"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/database/postgres/dml"
	"github.com/go-vela/types/library"

	"gorm.io/gorm"
)

func TestPostgres_Client_GetSecretList(t *testing.T) {
	// setup types
	_secretOne := testSecret()
	_secretOne.SetID(1)
	_secretOne.SetOrg("foo")
	_secretOne.SetRepo("bar")
	_secretOne.SetName("baz")
	_secretOne.SetValue("foob")
	_secretOne.SetType("repo")
	_secretOne.SetCreatedAt(1)
	_secretOne.SetCreatedBy("user")
	_secretOne.SetUpdatedAt(1)
	_secretOne.SetUpdatedBy("user2")

	_secretTwo := testSecret()
	_secretTwo.SetID(1)
	_secretTwo.SetOrg("foo")
	_secretTwo.SetRepo("bar")
	_secretTwo.SetName("foob")
	_secretTwo.SetValue("baz")
	_secretTwo.SetType("repo")
	_secretTwo.SetCreatedAt(1)
	_secretTwo.SetCreatedBy("user")
	_secretTwo.SetUpdatedAt(1)
	_secretTwo.SetUpdatedBy("user2")

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}

	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// capture the current expected SQL query
	//
	// https://gorm.io/docs/sql_builder.html#DryRun-Mode
	_query := _database.Postgres.Session(&gorm.Session{DryRun: true}).Raw(dml.ListSecrets).Statement

	// create expected return in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "type", "org", "repo", "team", "name", "value", "images", "events", "allow_command", "created_at", "created_by", "updated_at", "updated_by"},
	).AddRow(1, "repo", "foo", "bar", "", "baz", "foob", "{}", "{}", false, 1, "user", 1, "user2").
		AddRow(1, "repo", "foo", "bar", "", "foob", "baz", "{}", "{}", false, 1, "user", 1, "user2")

	// ensure the mock expects the query
	_mock.ExpectQuery(_query.SQL.String()).WillReturnRows(_rows)

	// setup tests
	tests := []struct {
		failure bool
		want    []*library.Secret
	}{
		{
			failure: false,
			want:    []*library.Secret{_secretOne, _secretTwo},
		},
	}

	// run tests
	for _, test := range tests {
		got, err := _database.GetSecretList()

		if test.failure {
			if err == nil {
				t.Errorf("GetSecretList should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetSecretList returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetSecretList is %v, want %v", got, test.want)
		}
	}
}

func TestPostgres_Client_GetTypeSecretList_Org(t *testing.T) {
	// setup types
	_secretOne := testSecret()
	_secretOne.SetID(1)
	_secretOne.SetOrg("foo")
	_secretOne.SetRepo("*")
	_secretOne.SetName("baz")
	_secretOne.SetValue("bar")
	_secretOne.SetType("org")
	_secretOne.SetCreatedAt(1)
	_secretOne.SetCreatedBy("user")
	_secretOne.SetUpdatedAt(1)
	_secretOne.SetUpdatedBy("user2")

	_secretTwo := testSecret()
	_secretTwo.SetID(1)
	_secretTwo.SetOrg("foo")
	_secretTwo.SetRepo("*")
	_secretTwo.SetName("bar")
	_secretTwo.SetValue("baz")
	_secretTwo.SetType("org")
	_secretTwo.SetCreatedAt(1)
	_secretTwo.SetCreatedBy("user")
	_secretTwo.SetUpdatedAt(1)
	_secretTwo.SetUpdatedBy("user2")

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}

	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// capture the current expected SQL query
	//
	// https://gorm.io/docs/sql_builder.html#DryRun-Mode
	_query := _database.Postgres.Session(&gorm.Session{DryRun: true}).Raw(dml.ListOrgSecrets, "foo", 1, 10).Statement

	// create expected return in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "type", "org", "repo", "team", "name", "value", "images", "events", "allow_command", "created_at", "created_by", "updated_at", "updated_by"},
	).AddRow(1, "org", "foo", "*", "", "baz", "bar", "{}", "{}", false, 1, "user", 1, "user2").
		AddRow(1, "org", "foo", "*", "", "bar", "baz", "{}", "{}", false, 1, "user", 1, "user2")

	// ensure the mock expects the query
	_mock.ExpectQuery(_query.SQL.String()).WillReturnRows(_rows)

	// setup tests
	tests := []struct {
		failure bool
		want    []*library.Secret
	}{
		{
			failure: false,
			want:    []*library.Secret{_secretOne, _secretTwo},
		},
	}

	// run tests
	for _, test := range tests {
		got, err := _database.GetTypeSecretList("org", "foo", "*", 1, 10, []string{})

		if test.failure {
			if err == nil {
				t.Errorf("GetTypeSecretList should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetTypeSecretList returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetTypeSecretList is %v, want %v", got, test.want)
		}
	}
}

func TestPostgres_Client_GetTypeSecretList_Repo(t *testing.T) {
	// setup types
	_secretOne := testSecret()
	_secretOne.SetID(1)
	_secretOne.SetOrg("foo")
	_secretOne.SetRepo("bar")
	_secretOne.SetName("baz")
	_secretOne.SetValue("foob")
	_secretOne.SetType("repo")
	_secretOne.SetCreatedAt(1)
	_secretOne.SetCreatedBy("user")
	_secretOne.SetUpdatedAt(1)
	_secretOne.SetUpdatedBy("user2")

	_secretTwo := testSecret()
	_secretTwo.SetID(1)
	_secretTwo.SetOrg("foo")
	_secretTwo.SetRepo("bar")
	_secretTwo.SetName("foob")
	_secretTwo.SetValue("baz")
	_secretTwo.SetType("repo")
	_secretTwo.SetCreatedAt(1)
	_secretTwo.SetCreatedBy("user")
	_secretTwo.SetUpdatedAt(1)
	_secretTwo.SetUpdatedBy("user2")

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}

	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// capture the current expected SQL query
	//
	// https://gorm.io/docs/sql_builder.html#DryRun-Mode
	_query := _database.Postgres.Session(&gorm.Session{DryRun: true}).Raw(dml.ListRepoSecrets, "foo", "bar", 1, 10).Statement

	// create expected return in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "type", "org", "repo", "team", "name", "value", "images", "events", "allow_command", "created_at", "created_by", "updated_at", "updated_by"},
	).AddRow(1, "repo", "foo", "bar", "", "baz", "foob", "{}", "{}", false, 1, "user", 1, "user2").
		AddRow(1, "repo", "foo", "bar", "", "foob", "baz", "{}", "{}", false, 1, "user", 1, "user2")

	// ensure the mock expects the query
	_mock.ExpectQuery(_query.SQL.String()).WillReturnRows(_rows)

	// setup tests
	tests := []struct {
		failure bool
		want    []*library.Secret
	}{
		{
			failure: false,
			want:    []*library.Secret{_secretOne, _secretTwo},
		},
	}

	// run tests
	for _, test := range tests {
		got, err := _database.GetTypeSecretList("repo", "foo", "bar", 1, 10, []string{})

		if test.failure {
			if err == nil {
				t.Errorf("GetTypeSecretList should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetTypeSecretList returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetTypeSecretList is %v, want %v", got, test.want)
		}
	}
}

func TestPostgres_Client_GetTypeSecretList_Shared(t *testing.T) {
	// setup types
	_secretOne := testSecret()
	_secretOne.SetID(1)
	_secretOne.SetOrg("foo")
	_secretOne.SetTeam("bar")
	_secretOne.SetName("baz")
	_secretOne.SetValue("foob")
	_secretOne.SetType("shared")
	_secretOne.SetCreatedAt(1)
	_secretOne.SetCreatedBy("user")
	_secretOne.SetUpdatedAt(1)
	_secretOne.SetUpdatedBy("user2")

	_secretTwo := testSecret()
	_secretTwo.SetID(1)
	_secretTwo.SetOrg("foo")
	_secretTwo.SetTeam("bar")
	_secretTwo.SetName("foob")
	_secretTwo.SetValue("baz")
	_secretTwo.SetType("shared")
	_secretTwo.SetCreatedAt(1)
	_secretTwo.SetCreatedBy("user")
	_secretTwo.SetUpdatedAt(1)
	_secretTwo.SetUpdatedBy("user2")

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}

	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// capture the current expected SQL query
	//
	// https://gorm.io/docs/sql_builder.html#DryRun-Mode
	_query := _database.Postgres.Session(&gorm.Session{DryRun: true}).Raw(dml.ListSharedSecrets, "foo", "bar", 1, 10).Statement

	// create expected return in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "type", "org", "repo", "team", "name", "value", "images", "events", "allow_command", "created_at", "created_by", "updated_at", "updated_by"},
	).AddRow(1, "shared", "foo", "", "bar", "baz", "foob", "{}", "{}", false, 1, "user", 1, "user2").
		AddRow(1, "shared", "foo", "", "bar", "foob", "baz", "{}", "{}", false, 1, "user", 1, "user2")

	// ensure the mock expects the query
	_mock.ExpectQuery(_query.SQL.String()).WillReturnRows(_rows)

	// setup tests
	tests := []struct {
		failure bool
		want    []*library.Secret
	}{
		{
			failure: false,
			want:    []*library.Secret{_secretOne, _secretTwo},
		},
	}

	// run tests
	for _, test := range tests {
		got, err := _database.GetTypeSecretList("shared", "foo", "bar", 1, 10, []string{"bar"})

		if test.failure {
			if err == nil {
				t.Errorf("GetTypeSecretList should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetTypeSecretList returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetTypeSecretList is %v, want %v", got, test.want)
		}
	}
}

func TestPostgres_Client_GetTypeSecretList_Shared_Wildcard(t *testing.T) {
	// setup types
	_secretOne := testSecret()
	_secretOne.SetID(1)
	_secretOne.SetOrg("foo")
	_secretOne.SetTeam("bar")
	_secretOne.SetName("baz")
	_secretOne.SetValue("foob")
	_secretOne.SetType("shared")
	_secretOne.SetCreatedAt(1)
	_secretOne.SetCreatedBy("user")
	_secretOne.SetUpdatedAt(1)
	_secretOne.SetUpdatedBy("user2")

	_secretTwo := testSecret()
	_secretTwo.SetID(1)
	_secretTwo.SetOrg("foo")
	_secretTwo.SetTeam("bared")
	_secretTwo.SetName("foob")
	_secretTwo.SetValue("baz")
	_secretTwo.SetType("shared")
	_secretTwo.SetCreatedAt(1)
	_secretTwo.SetCreatedBy("user")
	_secretTwo.SetUpdatedAt(1)
	_secretTwo.SetUpdatedBy("user2")

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}

	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// create expected return in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "type", "org", "repo", "team", "name", "value", "images", "events", "allow_command", "created_at", "created_by", "updated_at", "updated_by"},
	).AddRow(1, "shared", "foo", "", "bar", "baz", "foob", "{}", "{}", false, 1, "user", 1, "user2").
		AddRow(1, "shared", "foo", "", "bared", "foob", "baz", "{}", "{}", false, 1, "user", 1, "user2")

	// ensure the mock expects the query
	_mock.ExpectQuery("SELECT * FROM \"secrets\" WHERE (type = 'shared' AND org = $1) AND LOWER(team) IN ($2,$3) ORDER BY id DESC LIMIT 10").WillReturnRows(_rows)

	// setup tests
	tests := []struct {
		failure bool
		want    []*library.Secret
	}{
		{
			failure: false,
			want:    []*library.Secret{_secretOne, _secretTwo},
		},
	}

	// run tests
	for _, test := range tests {
		got, err := _database.GetTypeSecretList("shared", "foo", "*", 1, 10, []string{"bar", "bared"})

		if test.failure {
			if err == nil {
				t.Errorf("GetTypeSecretList should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetTypeSecretList returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetTypeSecretList is %v, want %v", got, test.want)
		}
	}
}
