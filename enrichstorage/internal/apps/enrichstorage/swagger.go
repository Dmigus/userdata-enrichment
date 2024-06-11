package enrichstorage

import "github.com/gin-gonic/gin"

// @title          Enricher
// @version         1.0
// @description     Сервис, который будет получать ФИО, из открытых апи обогащать ответ наиболее вероятными возрастом, полом и национальностью.

// @host      localhost:8080
// @BasePath  /api/v1

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

// GetRecords godoc
//
//		@Summary		Получение данных
//		@Description	Получение данных с различными фильтрами и пагинацией
//	 	@Tags			records
//		@Accept			json
//		@Produce		json
//		@Param			name	query		string	false	"name ="
//		@Param 			age[gte] query 		int		false	"age >="
//		@Success		200	{object}	get.response "resp"
//
//		@Router       /records/get [get]
func fakeGet(g *gin.Context) {

}

// CreateRecord godoc
//
//		@Summary		Новая запись
//		@Description	Создание новой записи на обогащение
//	 	@Tags			records
//		@Accept			json
//		@Produce		json
//		@Param			body	body		create.request	true	"ФИО"
//		@Success		200		 "Запись успешно создана"
//
//		@Router       /records/create [post]
func fakeCreate(g *gin.Context) {

}
