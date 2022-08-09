package api

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/massalabs/thyra/api/swagger/server/models"
	"github.com/massalabs/thyra/api/swagger/server/restapi/operations"
	"github.com/massalabs/thyra/pkg/my"
)

func DomainsHandler(params operations.MyDomainsParams) middleware.Responder {
	myDomains, err := my.NewDomains()
	if err != nil {
		return operations.NewUploadWebGetInternalServerError().
			WithPayload(
				&models.Error{
					Code:    errorCodeDomainsInstatiation,
					Message: err.Error(),
				})
	}

	list := myDomains.List()

	response := []*models.Websites{}
	for i := 0; i < len(list); i++ {
		response = append(response, &models.Websites{Name: list[i].URL, Address: list[i].Address})
	}

	return operations.NewMyDomainsOK().WithPayload(response)
}
