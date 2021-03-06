// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2016-2017 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package service

import (
	"net/http"
	"strings"

	"github.com/CanonicalLtd/serial-vault/datastore"
	"github.com/CanonicalLtd/serial-vault/service/keypair"
	"github.com/CanonicalLtd/serial-vault/service/signinglog"
	"github.com/CanonicalLtd/serial-vault/service/substore"
	"github.com/CanonicalLtd/serial-vault/service/user"
	"github.com/CanonicalLtd/serial-vault/usso"
	"github.com/gorilla/mux"
)

// SigningRouter returns the application route handler for the signing service methods
func SigningRouter() *mux.Router {

	// Start the web service router
	router := mux.NewRouter()

	// API routes
	router.Handle("/v1/version", Middleware(http.HandlerFunc(VersionHandler))).Methods("GET")
	router.Handle("/v1/serial", Middleware(ErrorHandler(SignHandler))).Methods("POST")
	router.Handle("/v1/request-id", Middleware(ErrorHandler(RequestIDHandler))).Methods("POST")
	router.Handle("/v1/model", Middleware(ErrorHandler(ModelAssertionHandler))).Methods("POST")
	router.Handle("/v1/pivot", Middleware(ErrorHandler(PivotModelHandler))).Methods("POST")
	router.Handle("/v1/pivotmodel", Middleware(ErrorHandler(PivotModelAssertionHandler))).Methods("POST")
	router.Handle("/v1/pivotserial", Middleware(ErrorHandler(PivotSerialAssertionHandler))).Methods("POST")

	return router
}

// AdminRouter returns the application route handler for administrating the application
func AdminRouter() *mux.Router {

	// Start the web service router
	router := mux.NewRouter()

	router.Handle("/v1/version", Middleware(http.HandlerFunc(VersionHandler))).Methods("GET")

	// API routes: csrf token and auth token
	router.Handle("/v1/token", MiddlewareWithCSRF(http.HandlerFunc(TokenHandler))).Methods("GET")
	router.Handle("/v1/authtoken", MiddlewareWithCSRF(http.HandlerFunc(TokenHandler))).Methods("GET")

	// API routes: models admin
	router.Handle("/v1/models", MiddlewareWithCSRF(http.HandlerFunc(ModelsHandler))).Methods("GET")
	router.Handle("/v1/models/assertion", MiddlewareWithCSRF(http.HandlerFunc(ModelAssertionHeadersHandler))).Methods("POST")
	router.Handle("/v1/models", MiddlewareWithCSRF(http.HandlerFunc(ModelCreateHandler))).Methods("POST")
	router.Handle("/v1/models/{id:[0-9]+}", MiddlewareWithCSRF(http.HandlerFunc(ModelGetHandler))).Methods("GET")
	router.Handle("/v1/models/{id:[0-9]+}", MiddlewareWithCSRF(http.HandlerFunc(ModelUpdateHandler))).Methods("PUT")
	router.Handle("/v1/models/{id:[0-9]+}", MiddlewareWithCSRF(http.HandlerFunc(ModelDeleteHandler))).Methods("DELETE")

	// API routes: signing-keys
	router.Handle("/v1/keypairs", MiddlewareWithCSRF(http.HandlerFunc(keypair.List))).Methods("GET")
	router.Handle("/v1/keypairs", MiddlewareWithCSRF(http.HandlerFunc(KeypairCreateHandler))).Methods("POST")
	router.Handle("/v1/keypairs/{id:[0-9]+}/disable", MiddlewareWithCSRF(http.HandlerFunc(KeypairDisableHandler))).Methods("POST")
	router.Handle("/v1/keypairs/{id:[0-9]+}/enable", MiddlewareWithCSRF(http.HandlerFunc(KeypairEnableHandler))).Methods("POST")
	router.Handle("/v1/keypairs/assertion", MiddlewareWithCSRF(http.HandlerFunc(KeypairAssertionHandler))).Methods("POST")

	router.Handle("/v1/keypairs/generate", MiddlewareWithCSRF(http.HandlerFunc(KeypairGenerateHandler))).Methods("POST")
	router.Handle("/v1/keypairs/status/{authorityID}/{keyName}", MiddlewareWithCSRF(http.HandlerFunc(KeypairStatusHandler))).Methods("GET")
	router.Handle("/v1/keypairs/status", MiddlewareWithCSRF(http.HandlerFunc(KeypairStatusProgressHandler))).Methods("GET")
	router.Handle("/v1/keypairs/register", MiddlewareWithCSRF(http.HandlerFunc(StoreKeyRegisterHandler))).Methods("POST")

	// API routes: signing log
	router.Handle("/v1/signinglog", MiddlewareWithCSRF(http.HandlerFunc(signinglog.List))).Methods("GET")
	router.Handle("/v1/signinglog/filters", MiddlewareWithCSRF(http.HandlerFunc(signinglog.ListFilters))).Methods("GET")

	// API routes: account assertions
	router.Handle("/v1/accounts", MiddlewareWithCSRF(http.HandlerFunc(AccountsHandler))).Methods("GET")
	router.Handle("/v1/accounts", MiddlewareWithCSRF(http.HandlerFunc(AccountCreateHandler))).Methods("POST")
	router.Handle("/v1/accounts/{id:[0-9]+}", MiddlewareWithCSRF(http.HandlerFunc(AccountUpdateHandler))).Methods("PUT")
	router.Handle("/v1/accounts/{id:[0-9]+}", MiddlewareWithCSRF(http.HandlerFunc(AccountGetHandler))).Methods("GET")
	router.Handle("/v1/accounts/upload", MiddlewareWithCSRF(http.HandlerFunc(AccountsUploadHandler))).Methods("POST")
	router.Handle("/v1/accounts/{id:[0-9]+}/stores", MiddlewareWithCSRF(http.HandlerFunc(substore.List))).Methods("GET")
	router.Handle("/v1/accounts/stores/{id:[0-9]+}", MiddlewareWithCSRF(http.HandlerFunc(substore.Update))).Methods("PUT")
	router.Handle("/v1/accounts/stores/{id:[0-9]+}", MiddlewareWithCSRF(http.HandlerFunc(substore.Delete))).Methods("DELETE")
	router.Handle("/v1/accounts/stores", MiddlewareWithCSRF(http.HandlerFunc(substore.Create))).Methods("POST")

	// API routes: system-user assertion
	router.Handle("/v1/assertions", MiddlewareWithCSRF(http.HandlerFunc(SystemUserAssertionHandler))).Methods("POST")

	// API routes: users management
	router.Handle("/v1/users", MiddlewareWithCSRF(http.HandlerFunc(user.List))).Methods("GET")
	router.Handle("/v1/users", MiddlewareWithCSRF(http.HandlerFunc(user.Create))).Methods("POST")
	router.Handle("/v1/users/{id:[0-9]+}", MiddlewareWithCSRF(http.HandlerFunc(user.Get))).Methods("GET")
	router.Handle("/v1/users/{id:[0-9]+}", MiddlewareWithCSRF(http.HandlerFunc(user.Update))).Methods("PUT")
	router.Handle("/v1/users/{id:[0-9]+}", MiddlewareWithCSRF(http.HandlerFunc(user.Delete))).Methods("DELETE")
	router.Handle("/v1/users/{id:[0-9]+}/otheraccounts", MiddlewareWithCSRF(http.HandlerFunc(user.GetOtherAccounts))).Methods("GET")

	// OpenID routes: using Ubuntu SSO
	router.Handle("/login", MiddlewareWithCSRF(http.HandlerFunc(usso.LoginHandler)))
	router.Handle("/logout", MiddlewareWithCSRF(http.HandlerFunc(usso.LogoutHandler)))

	// Web application routes
	path := []string{datastore.Environ.Config.DocRoot, "/static/"}
	fs := http.StripPrefix("/static/", http.FileServer(http.Dir(strings.Join(path, ""))))
	router.PathPrefix("/static/").Handler(fs)
	router.PathPrefix("/signing-keys").Handler(MiddlewareWithCSRF(http.HandlerFunc(IndexHandler)))
	router.PathPrefix("/models").Handler(MiddlewareWithCSRF(http.HandlerFunc(IndexHandler)))
	router.PathPrefix("/keypairs").Handler(MiddlewareWithCSRF(http.HandlerFunc(IndexHandler)))
	router.PathPrefix("/accounts").Handler(MiddlewareWithCSRF(http.HandlerFunc(IndexHandler)))
	router.PathPrefix("/signinglog").Handler(MiddlewareWithCSRF(http.HandlerFunc(IndexHandler)))
	router.PathPrefix("/systemuser").Handler(MiddlewareWithCSRF(http.HandlerFunc(IndexHandler)))
	router.PathPrefix("/users").Handler(MiddlewareWithCSRF(http.HandlerFunc(IndexHandler)))
	router.PathPrefix("/notfound").Handler(MiddlewareWithCSRF(http.HandlerFunc(IndexHandler)))
	router.Handle("/", MiddlewareWithCSRF(http.HandlerFunc(IndexHandler))).Methods("GET")

	// Admin API routes
	router.Handle("/api/signinglog", Middleware(http.HandlerFunc(signinglog.APIList))).Methods("GET")
	router.Handle("/api/keypairs", Middleware(http.HandlerFunc(keypair.APIList))).Methods("GET")
	router.Handle("/api/accounts/{id:[0-9]+}/stores", Middleware(http.HandlerFunc(substore.APIList))).Methods("GET")
	router.Handle("/api/accounts/stores/{id:[0-9]+}", Middleware(http.HandlerFunc(substore.APIUpdate))).Methods("PUT")
	router.Handle("/api/accounts/stores/{id:[0-9]+}", Middleware(http.HandlerFunc(substore.APIDelete))).Methods("DELETE")
	router.Handle("/api/accounts/stores", Middleware(http.HandlerFunc(substore.APICreate))).Methods("POST")

	return router
}

// SystemUserRouter returns the application route handler for the system-user service methods
func SystemUserRouter() *mux.Router {

	// Start the web service router
	router := mux.NewRouter()

	// API routes
	router.Handle("/v1/version", Middleware(http.HandlerFunc(VersionHandler))).Methods("GET")
	router.Handle("/v1/token", Middleware(http.HandlerFunc(TokenHandler))).Methods("GET")
	router.Handle("/v1/models", Middleware(http.HandlerFunc(ModelsHandler))).Methods("GET")
	router.Handle("/v1/assertions", Middleware(http.HandlerFunc(SystemUserAssertionHandler))).Methods("POST")

	// Web application routes
	path := []string{datastore.Environ.Config.DocRoot, "/static/"}
	fs := http.StripPrefix("/static/", http.FileServer(http.Dir(strings.Join(path, ""))))
	router.PathPrefix("/static/").Handler(fs)
	router.Handle("/", Middleware(http.HandlerFunc(UserIndexHandler))).Methods("GET")

	return router
}
