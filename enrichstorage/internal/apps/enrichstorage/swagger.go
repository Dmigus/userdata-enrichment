package enrichstorage

import (
	"github.com/gin-gonic/gin"
)

// @title          Enricher
// @version         1.0
// @description     Сервис, который будет получать ФИО, из открытых апи обогащать ответ наиболее вероятными возрастом, полом и национальностью.

// @host      localhost:8081
// @BasePath  /api/v1
// @securitydefinitions.oauth2.password OAuth2
// @tokenUrl http://localhost:8085/realms/enricherrealm/protocol/openid-connect/token
// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

// GetRecords godoc
//
// @Summary		Получение данных
// @Description	Получение данных с различными фильтрами и пагинацией
// @Tags			records
// @Security OAuth2
// @Accept			json
// @Produce		json
// @Param			name	query		string	false	"name ="
// @Param			surname	query		string	false	"surname ="
// @Param			patronymic	query		string	false	"patronymic ="
// @Param 			age[gte] query 		int		false	"age >="
// @Param 			age[lte] query 		int		false	"age <="
// @Param 			sex query 		string		false	"sex ="
// @Param 			nationality query 		string		false	"nationality ="
// @Param 			limit query 		int		false	"limit ="
// @Param 			after query 		string		false	"after ="
// @Param 			before query 		string		false	"before ="
// @Success		200	{object}	get.response "resp"
//
// @Router       /records/get [get]
func fakeGet(g *gin.Context) {

}

// CreateRecord godoc
//
// @Summary		Новая запись
// @Description	Создание новой записи на обогащение
// @Tags			records
// @Security OAuth2
// @Accept			json
// @Produce		json
// @Param			body	body		create.request	true	"ФИО"
// @Success		200		 "Запись успешно создана"
//
//	@Router       /records/create [post]
func fakeCreate(g *gin.Context) {

}

// DeleteRecord godoc
//
// @Summary		Удалить запись
// @Description	Удаление записи о ФИО
// @Tags			records
// @Security OAuth2
// @Accept			json
// @Produce		json
// @Param			body	body		delete.request	true	"ФИО"
// @Success		200		 "Запись успешно удалена"
//
//	@Router       /records/delete [post]
func fakeDelete(g *gin.Context) {

}

// UpdateRecord godoc
//
// @Summary		Обновить запись
// @Description	Обновление записи для ФИО
// @Tags			records
// @Security OAuth2
// @Accept			json
// @Produce		json
// @Param			body	body		update.request	true	"ФИО"
// @Success		200		 "Запись успешно изменена"
//
//	@Router       /records/update [post]
func fakeUpdate(g *gin.Context) {

}
