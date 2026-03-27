package services

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/udistrital/facturacion_electronica_mid/models"
	"github.com/udistrital/utils_oas/request"
)

func ObtenerDuenoReciboTerceros(terceroDuenoId int) (models.TerceroId, error) {

	var terceroDuenoData models.TerceroId

	url := beego.AppConfig.String("TercerosService") + fmt.Sprintf("/tercero/%d", terceroDuenoId)

	err := request.GetJson(url, &terceroDuenoData)
	if err != nil {
		return models.TerceroId{}, err
	}

	return terceroDuenoData, nil
}
