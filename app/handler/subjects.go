package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	gremcos "github.com/supplyon/gremcos"

	"github.com/imeplusplus/dont-panic-api/app/dbOperations"
	"github.com/imeplusplus/dont-panic-api/app/logger"
	apiModel "github.com/imeplusplus/dont-panic-api/app/model/api"
	storageModel "github.com/imeplusplus/dont-panic-api/app/model/storage"
)

func GetSubjects(cosmos gremcos.Cosmos, w http.ResponseWriter, _ *http.Request) {
	storageSubjects, err := dbOperations.GetSubjects(cosmos)
	if err != nil {
		msg := logger.FailedToExecuteGremlinQuery{}
		log.Error().Stack().Err(err).Msg(msg.Info())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(storageSubjects)
	if err != nil {
		msg := logger.FailedToEncodeJSON{
			Resource: "storageModel.Subject",
		}
		log.Error().Stack().Err(err).Msg(msg.Info())
		w.WriteHeader(http.StatusInternalServerError)
	}

	msg := logger.ResourceRead{
		ResourceName:    "subjects",
		ResourceContent: storageModel.PrettyPrint(storageSubjects),
	}
	log.Info().Msg(msg.Info())
}

func GetSubject(cosmos gremcos.Cosmos, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	storageSubject, err := dbOperations.GetSubjectByName(cosmos, vars["name"])
	if err != nil {
		msg := logger.FailedToExecuteGremlinQuery{}
		log.Error().Stack().Err(err).Msg(msg.Info())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(storageSubject)
	if err != nil {
		msg := logger.FailedToEncodeJSON{
			Resource: "storageModel.Subject",
		}
		log.Error().Stack().Err(err).Msg(msg.Info())
		w.WriteHeader(http.StatusInternalServerError)
	}

	msg := logger.ResourceRead{
		ResourceName:    "subjects/" + storageSubject.Name,
		ResourceContent: storageModel.PrettyPrint(storageSubject),
	}
	log.Info().Msg(msg.Info())
}

func UpdateSubject(cosmos gremcos.Cosmos, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	apiSubject := apiModel.Subject{}
	if err := json.NewDecoder(r.Body).Decode(&apiSubject); err != nil {
		msg := logger.FailedToDecodeJSON{
			Resource: "apiModel.Subject",
		}
		log.Error().Stack().Err(err).Msg(msg.Info())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	storageSubject, err := dbOperations.UpdateSubject(cosmos, storageModel.Subject(apiSubject), name)
	if err != nil {
		msg := logger.FailedToExecuteGremlinQuery{}
		log.Error().Stack().Err(err).Msg(msg.Info())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	err = json.NewEncoder(w).Encode(storageSubject)

	if err != nil {
		msg := logger.FailedToEncodeJSON{
			Resource: "apiModel.Subject",
		}
		log.Error().Stack().Err(err).Msg(msg.Info())
		w.WriteHeader(http.StatusInternalServerError)
	}

	msg := logger.ResourceUpdated{
		ResourceName:    "subjects/" + storageSubject.Name,
		ResourceContent: storageModel.PrettyPrint(storageSubject),
	}
	log.Info().Msg(msg.Info())
}

func DeleteSubject(cosmos gremcos.Cosmos, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	err := dbOperations.DeleteSubject(cosmos, name)
	if err != nil {
		msg := logger.FailedToExecuteGremlinQuery{}
		log.Error().Stack().Err(err).Msg(msg.Info())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	msg := logger.ResourceDeleted{
		ResourceName: "subjects/" + name,
	}
	log.Info().Msg(msg.Info())
}

func CreateSubject(cosmos gremcos.Cosmos, w http.ResponseWriter, r *http.Request) {
	apiSubject := apiModel.Subject{}
	if err := json.NewDecoder(r.Body).Decode(&apiSubject); err != nil {
		msg := logger.FailedToDecodeJSON{
			Resource: "apiModel.Subject",
		}
		log.Error().Stack().Err(err).Msg(msg.Info())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	storageSubject, err := dbOperations.CreateSubject(cosmos, storageModel.Subject(apiSubject))
	if err != nil {
		msg := logger.FailedToExecuteGremlinQuery{}
		log.Error().Stack().Err(err).Msg(msg.Info())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(storageSubject)
	if err != nil {
		msg := logger.FailedToEncodeJSON{
			Resource: "storageModel.Subject",
		}
		log.Error().Stack().Err(err).Msg(msg.Info())
		w.WriteHeader(http.StatusInternalServerError)
	}

	fmt.Println(storageSubject.Name)
	msg := logger.ResourceCreated{
		ResourceName:    "subjects/" + storageSubject.Name,
		ResourceContent: storageModel.PrettyPrint(storageSubject),
	}
	log.Info().Msg(msg.Info())
}
