/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore-apiserver.
 *
 * kore-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */

package users

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"testing"

	"github.com/romanyx/polluter"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

var (
	once    = &sync.Once{}
	dbClean = &sync.Once{}
	bench   = &sync.Once{}
	db      *sql.DB
)

var (
	postgresURL = "host=127.0.0.1 port=5432 user=postgres dbname=postgres password=pass sslmode=disable"
)

func getTestDBURL() string {
	if os.Getenv("TEST_USERS_DATABASE_URL") != "" {
		return os.Getenv("TEST_USERS_DATABASE_URL")
	}

	return "root:pass@tcp(127.0.0.1:3306)/kore?parseTime=true"
}

func getTestDBDriver() string {
	if os.Getenv("TEST_USERS_DATABASE_DRIVER") != "" {
		return os.Getenv("TEST_USERS_DATABASE_DRIVER")
	}

	return "mysql"
}

func makeTestConfig() Config {
	return Config{
		Driver:        getTestDBDriver(),
		EnableLogging: false,
		StoreURL:      getTestDBURL(),
	}
}

type testframework interface {
	Fatalf(string, ...interface{})
}

func makeTestStore(t testframework) Interface {
	driver := getTestDBDriver()
	url := getTestDBURL()

	dbClean.Do(func() {
		d, err := sql.Open(driver, url)
		if err != nil {
			t.Fatalf("failed to open the database connection: %s", err)
		}
		db = d

		switch driver {
		case "mysql":
			if _, err := db.Exec("drop database if exists kore"); err != nil {
				t.Fatalf("failed to drop the database: %s", err)
			}

			if _, err := db.Exec("create database if not exists kore"); err != nil {
				t.Fatalf("failed to create the database: %s", err)
			}

			if _, err := db.Exec("use kore"); err != nil {
				t.Fatalf("failed to select database: %s", err)
			}
		}
	})

	store, err := New(makeTestConfig())
	if err != nil {
		t.Fatalf("faild to create a db store: %s", err)
	}

	once.Do(func() {
		for _, x := range []string{"db.yml"} {
			content, err := ioutil.ReadFile(fmt.Sprintf("fixtures/%s", x))
			if err != nil {
				t.Fatalf("failed to open database file: %s", err)
			}
			var p *polluter.Polluter
			switch driver {
			case "mysql":
				p = polluter.New(polluter.MySQLEngine(db))
			case "postgres":
				p = polluter.New(polluter.PostgresEngine(db))
			default:
				t.Fatalf("unknown driver")
			}

			if err := p.Pollute(bytes.NewReader(content)); err != nil {
				t.Fatalf("failed to pollute database: %s", err)
			}
		}
	})

	return store
}

func TestNewBad(t *testing.T) {
	store, err := New(Config{Driver: "non"})
	assert.Error(t, err)
	assert.Nil(t, store)
}

func TestNewOK(t *testing.T) {
	store, err := New(makeTestConfig())
	require.NoError(t, err)
	require.NotNil(t, store)
	defer store.Stop()
}

func TestTeams(t *testing.T) {
	store, err := New(makeTestConfig())
	require.NoError(t, err)
	require.NotNil(t, store)
	defer store.Stop()
	assert.NotNil(t, store.Teams())
}

func TestUsers(t *testing.T) {
	store, err := New(makeTestConfig())
	require.NoError(t, err)
	require.NotNil(t, store)
	defer store.Stop()
	assert.NotNil(t, store.Users())
}

func TestAudit(t *testing.T) {
	store, err := New(makeTestConfig())
	require.NoError(t, err)
	require.NotNil(t, store)
	defer store.Stop()
	assert.NotNil(t, store.Audit())
}
