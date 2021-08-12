package product

import (
	"encoding/json"
	"net/http"
	"product_management/utils"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

func NewProductHandler(db *gorm.DB, u ProductUsecaseInterface) *ProductHandler {
	return &ProductHandler{
		usecase: u,
	}
}

type ProductHandler struct {
	usecase ProductUsecaseInterface
}

func (p *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	offset := r.URL.Query().Get("offset")
	limit := r.URL.Query().Get("limit")

	data, msg, err := p.usecase.GetProducts(r.Context(), offset, limit)
	resp := msg
	statusCode := msg["status"].(interface{}).(int)

	if err != nil {
		utils.RespondwithJSON(w, statusCode, resp)
		log.WithFields(log.Fields{"APIName": "GetProducts", "error": err}).Error(msg["message"].(interface{}).(string))
		return
	}
	resp["data"] = data
	utils.RespondwithJSON(w, statusCode, resp)

}

func (p *ProductHandler) GetProductById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.RespondwithJSON(w, 400, "Invalid product id")
		return
	}

	data, msg, err := p.usecase.GetProductById(r.Context(), id)

	resp := msg
	statusCode := msg["status"].(interface{}).(int)
	if err != nil {
		utils.RespondwithJSON(w, statusCode, resp)
		log.WithFields(log.Fields{"APIName": "GetProductById", "error": err}).Error(msg["message"].(interface{}).(string))
		return
	}
	resp["data"] = data
	utils.RespondwithJSON(w, statusCode, resp)
}

func (p *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	product := &Product{}

	err := json.NewDecoder(r.Body).Decode(product)

	if err != nil {
		utils.RespondwithJSON(w, 400, utils.Message(400, "Error while decoding request body"))
		log.WithFields(log.Fields{"APIName": "CreateProduct", "error": err}).Error("Error while decoding request body", err)
		return
	}

	data, msg, err := p.usecase.CreateProduct(r.Context(), product)
	resp := msg
	statusCode := msg["status"].(interface{}).(int)
	if err != nil {
		utils.RespondwithJSON(w, statusCode, resp)
		log.WithFields(log.Fields{"APIName": "CreateProduct", "error": err}).Error(msg["message"].(interface{}).(string))
		return
	}
	resp["data"] = data
	utils.RespondwithJSON(w, statusCode, resp)
}

func (p *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.RespondwithJSON(w, 400, "Invalid product id")
		return
	}
	msg, err := p.usecase.DeleteProduct(r.Context(), id)
	resp := msg
	statusCode := msg["status"].(interface{}).(int)
	if err != nil {
		utils.RespondwithJSON(w, statusCode, resp)
		log.WithFields(log.Fields{"APIName": "DeleteProduct", "error": err}).Error(msg["message"].(interface{}).(string))
		return
	}
	utils.RespondwithJSON(w, statusCode, resp)

}

func (p *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.RespondwithJSON(w, 400, "Invalid product id")
		return
	}
	json.NewDecoder(r.Body).Decode(&product)
	data, msg, err := p.usecase.UpdateProduct(r.Context(), id, &product)
	resp := msg
	statusCode := msg["status"].(interface{}).(int)
	if err != nil {
		utils.RespondwithJSON(w, statusCode, resp)
		log.WithFields(log.Fields{"APIName": "UpdateProduct", "error": err}).Error(msg["message"].(interface{}).(string))
		return
	}
	resp["data"] = data
	utils.RespondwithJSON(w, statusCode, resp)

}
