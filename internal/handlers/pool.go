package handlers

import (
	"github.com/gladiusio/gladius-application-server/internal/controller"
	"github.com/gladiusio/gladius-common/pkg/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
)

// Retrieve Pool Information
func PublicPoolInformationHandler(database *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		poolInformation, err := controller.PoolInformation(database)
		if err != nil {
			handlers.ErrorHandler(w, r, "Could retrieve Public Information", err, http.StatusBadRequest)
			return
		}

		handlers.ResponseHandler(w, r, "null", true, nil, poolInformation, nil)
	}
}

type PoolContainsWallet struct {
	ContainsWallet bool
}

func PoolContainsNode(database *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		walletAddress := vars["walletAddress"]

		containsWallet, nodeErr := controller.NodeInPool(database, walletAddress)
		poolOwnerContainsWallet, poolErr := controller.PoolWalletOwner(database, walletAddress)

		if !(containsWallet || poolOwnerContainsWallet) {
			if nodeErr != nil {
				handlers.ErrorHandler(w, r, "Could not find node record", nodeErr, http.StatusNotFound)
				return
			}

			if poolErr != nil {
				handlers.ErrorHandler(w, r, "Could not find pool record", poolErr, http.StatusNotFound)
				return
			}
		} else {
			handlers.ResponseHandler(w, r, "null", true, nil, PoolContainsWallet{ContainsWallet: (containsWallet || poolOwnerContainsWallet)}, nil)
		}
	}
}

func PoolNodes(database *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		nodes, err := controller.NodesAccepted(database)
		if err != nil {
			handlers.ErrorHandler(w, r, "Could not retrieve nodes", err, http.StatusInternalServerError)
			return
		}

		var nodeAddresses []string

		for _, node := range nodes {
			nodeAddresses = append(nodeAddresses, node.Wallet)
		}

		handlers.ResponseHandler(w, r, "null", true, nil, nodeAddresses, nil)

	}
}
