// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2018-2019 Canonical Ltd
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

package substore

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/CanonicalLtd/serial-vault/datastore"
	"github.com/CanonicalLtd/serial-vault/service/auth"
	"github.com/CanonicalLtd/serial-vault/service/response"
)

// ListResponse is the JSON response from the API sub-stores method
type ListResponse struct {
	Success      bool                 `json:"success"`
	ErrorCode    string               `json:"error_code"`
	ErrorSubcode string               `json:"error_subcode"`
	ErrorMessage string               `json:"message"`
	Substores    []datastore.Substore `json:"substores"`
}

// listHandler is the API method to fetch the list sub-stores
func listHandler(w http.ResponseWriter, user datastore.User, apiCall bool, accountID int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err := auth.CheckUserPermissions(user, datastore.Admin, apiCall)
	if err != nil {
		response.FormatStandardResponse(false, "error-auth", "", "", w)
		return
	}

	//logs, err := datastore.Environ.DB.ListAllowedSigningLog(user)
	stores, err := datastore.Environ.DB.ListSubstores(accountID, user)
	if err != nil {
		response.FormatStandardResponse(false, "error-stores-json", "", "", w)
		return
	}

	// Return successful JSON response with the list of models
	w.WriteHeader(http.StatusOK)
	formatListResponse(true, "", "", "", stores, w)
}

func updateHandler(w http.ResponseWriter, user datastore.User, apiCall bool, storeID int, store datastore.Substore) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err := auth.CheckUserPermissions(user, datastore.Admin, apiCall)
	if err != nil {
		response.FormatStandardResponse(false, "error-auth", "", "", w)
		return
	}

	if storeID != store.ID {
		response.FormatStandardResponse(false, "error-stores-json", "", "The store IDs do not match", w)
		return
	}

	err = datastore.Environ.DB.UpdateAllowedSubstore(store, user)
	if err != nil {
		log.Println("Error updating the store:", err)
		response.FormatStandardResponse(false, "error-stores-substore", "", "Error updating the store", w)
		return
	}

	// Return successful JSON response
	w.WriteHeader(http.StatusOK)
	response.FormatStandardResponse(true, "", "", "", w)
}

func createHandler(w http.ResponseWriter, user datastore.User, apiCall bool, store datastore.Substore) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err := auth.CheckUserPermissions(user, datastore.Admin, apiCall)
	if err != nil {
		response.FormatStandardResponse(false, "error-auth", "", "", w)
		return
	}

	err = datastore.Environ.DB.CreateAllowedSubstore(store, user)
	if err != nil {
		response.FormatStandardResponse(false, "error-stores-json", "", "", w)
		return
	}

	// Return successful JSON response
	w.WriteHeader(http.StatusOK)
	response.FormatStandardResponse(true, "", "", "", w)
}

func formatListResponse(success bool, errorCode, errorSubcode, message string, stores []datastore.Substore, w http.ResponseWriter) error {
	response := ListResponse{Success: success, ErrorCode: errorCode, ErrorSubcode: errorSubcode, ErrorMessage: message, Substores: stores}

	// Encode the response as JSON
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Error forming the sub-stores response.")
		return err
	}
	return nil
}

func deleteHandler(w http.ResponseWriter, user datastore.User, apiCall bool, storeID int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err := auth.CheckUserPermissions(user, datastore.Admin, apiCall)
	if err != nil {
		response.FormatStandardResponse(false, "error-auth", "", "", w)
		return
	}

	errorSubcode, err := datastore.Environ.DB.DeleteAllowedSubstore(storeID, user)
	if err != nil {
		response.FormatStandardResponse(false, "error-deleting-store", errorSubcode, err.Error(), w)
		return
	}

	// Return successful JSON response
	w.WriteHeader(http.StatusOK)
	response.FormatStandardResponse(true, "", "", "", w)
}
