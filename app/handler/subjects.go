package handler

import (
	"encoding/json"
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
	res, err := dbOperations.GetSubjects(cosmos)

	if err != nil {
		log.Error().Stack().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)

	if err != nil {
		log.Error().Stack().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
	}
	msg := logger.ResourceRead{
		Resource: storageModel.PrettyPrint(res),
	}
	log.Info().Msg(msg.Info())
}

func GetSubject(cosmos gremcos.Cosmos, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	res, err := dbOperations.GetSubjectByName(cosmos, vars["name"])

	if err != nil {
		log.Error().Stack().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)

	if err != nil {
		log.Error().Stack().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
	}

	msg := logger.ResourceRead{
		Resource: storageModel.PrettyPrint(res),
	}
	log.Info().Msg(msg.Info())
}

func UpdateSubject(cosmos gremcos.Cosmos, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	oldSubject, err := dbOperations.GetSubjectByName(cosmos, name)

	subject := apiModel.Subject{}
	if err := json.NewDecoder(r.Body).Decode(&subject); err != nil {
		log.Error().Stack().Err(err).Msg("Couldn't parse request body into apiModel.Subject")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := dbOperations.UpdateSubject(cosmos, subject, name)

	if err != nil {
		log.Error().Stack().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	err = json.NewEncoder(w).Encode(res)

	if err != nil {
		log.Error().Stack().Err(err).Msg("Couldn't create response body with created subject")
		w.WriteHeader(http.StatusInternalServerError)
	}
	msg := logger.ResourceUpdated{
		PastResource: storageModel.PrettyPrint(oldSubject),
		NewResource:  storageModel.PrettyPrint(res),
	}
	log.Info().Msg(msg.Info())
}

func DeleteSubject(cosmos gremcos.Cosmos, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	subject, err := dbOperations.GetSubjectByName(cosmos, name)

	if err != nil {
		log.Error().Stack().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = dbOperations.DeleteSubject(cosmos, name)

	if err != nil {
		log.Error().Stack().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	msg := logger.ResourceCreated{
		Resource: storageModel.PrettyPrint(subject),
	}
	log.Info().Msg(msg.Info())
}

func CreateSubject(cosmos gremcos.Cosmos, w http.ResponseWriter, r *http.Request) {
	subject := apiModel.Subject{}
	if err := json.NewDecoder(r.Body).Decode(&subject); err != nil {
		log.Error().Stack().Err(err).Msg("Couldn't parse request body into apiModel.Subject")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := dbOperations.CreateSubject(cosmos, subject)

	if err != nil {
		log.Error().Stack().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(res)

	if err != nil {
		log.Error().Stack().Err(err).Msg("Couldn't create response body with created subject")
		w.WriteHeader(http.StatusInternalServerError)
	}

	msg := logger.ResourceCreated{
		Resource: storageModel.PrettyPrint(res),
	}
	log.Info().Msg(msg.Info())
}
