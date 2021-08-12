package handlers

import (
	"encoding/json"
	"main_project/models"
	"main_project/utils"
	"net/http"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
)

var CreateProduct = func(w http.ResponseWriter, r *http.Request) {
	product := &models.Product{}
	err := json.NewDecoder(r.Body).Decode(product)
	if err != nil {
		utils.Respond(w, utils.Message(400, "Error while decoding request body"), 400)
		log.WithFields(log.Fields{"APIName": "CreateProduct", "error": err}).Error("Error in decoding request body")
		return
	}

	data, msg, err := product.Create()
	resp := msg
	statusCode := msg["status"].(interface{}).(int)
	if err != nil {
		utils.Respond(w, resp, statusCode)
		log.WithFields(log.Fields{"APIName": "CreateProduct", "error": err}).Error(msg["message"].(interface{}).(string))
		return
	}
	resp["data"] = data
	utils.Respond(w, resp, statusCode)
}

var GetProducts = func(w http.ResponseWriter, r *http.Request) {
	offset := r.URL.Query().Get("offset")
	limit := r.URL.Query().Get("limit")

	data, msg, err := models.GetProducts(offset, limit)
	resp := msg
	statusCode := msg["status"].(interface{}).(int)

	if err != nil {
		utils.Respond(w, resp, statusCode)
		log.WithFields(log.Fields{"APIName": "GetProducts", "error": err}).Error(msg["message"].(interface{}).(string))
		return
	}
	resp["data"] = data
	utils.Respond(w, resp, statusCode)
}

var GetProductById = func(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	data, msg, err := models.GetProductById(id)
	resp := msg
	statusCode := msg["status"].(interface{}).(int)
	if err != nil {
		utils.Respond(w, resp, statusCode)
		log.WithFields(log.Fields{"APIName": "GetProductById", "error": err}).Error(msg["message"].(interface{}).(string))
		return
	}
	resp["data"] = data
	utils.Respond(w, resp, statusCode)
}

var DeleteProduct = func(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	msg, err := models.DeleteProduct(id)
	resp := msg
	statusCode := msg["status"].(interface{}).(int)
	if err != nil {
		utils.Respond(w, resp, statusCode)
		log.WithFields(log.Fields{"APIName": "DeleteProduct", "error": err}).Error(msg["message"].(interface{}).(string))
		return
	}
	utils.Respond(w, resp, statusCode)

}

var UpdateProduct = func(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	id := chi.URLParam(r, "id")
	json.NewDecoder(r.Body).Decode(&product)
	data, msg, err := models.UpdateProduct(id, &product)
	resp := msg
	statusCode := msg["status"].(interface{}).(int)
	if err != nil {
		utils.Respond(w, resp, statusCode)
		log.WithFields(log.Fields{"APIName": "UpdateProduct", "error": err}).Error(msg["message"].(interface{}).(string))
		return
	}
	resp["data"] = data
	utils.Respond(w, resp, statusCode)

}
