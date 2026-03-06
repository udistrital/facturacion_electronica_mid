package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/facturacion_electronica_mid/models"
	"github.com/udistrital/facturacion_electronica_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
)

type FacturacionElectronicaController struct {
	beego.Controller
}

func (c *FacturacionElectronicaController) URLMapping() {
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("Post", c.Post)
	c.Mapping("Put", c.Put)
}

// GetAll obtiene todos los registros
// @Title GetAll
// @Description Obtiene todos los registros
// @Success 200 {object} []map[string]interface{}
// @Failure 404 No se encontraron registros
// @router / [get]
func (c *FacturacionElectronicaController) GetAll() {
	var resultado map[string]interface{}

	url := "http://" + beego.AppConfig.String("FacturacionElectronicaService")

	err := request.GetJsonWSO2(url, &resultado)

	if err != nil {
		logs.Info("URL completa: " + url)
		logs.Error(err)
		c.Data["json"] = map[string]interface{}{
			"Success": false,
			"Status":  "404",
			"Message": "Error al obtener los datos: " + err.Error(),
			"Data":    nil,
		}
	} else {
		c.Data["json"] = map[string]interface{}{
			"Success": true,
			"Status":  "200",
			"Message": "Datos obtenidos correctamente",
			"Data":    resultado,
		}
	}

	c.ServeJSON()
}

// GetOne obtiene un registro específico
// @Title GetOne
// @Description Obtiene un registro por su ID y año
// @Param id path string true "ID del registro"
// @Param anio path string true "Año de consulta"
// @Success 200 {object} map[string]interface{}
// @Failure 404 No se encontró el registro
// @router /:id/:anio [get]
func (c *FacturacionElectronicaController) GetOne() {
	id := c.Ctx.Input.Param(":id")
	anio := c.Ctx.Input.Param(":anio")
	var resultado map[string]interface{}

	url := "http://" + beego.AppConfig.String("FacturacionElectronicaService") + "/" + id + "/" + anio

	err := request.GetJsonWSO2(url, &resultado)

	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"Success": false,
			"Status":  "404",
			"Message": "Error al obtener el registro: " + err.Error(),
			"Data":    nil,
		}
	} else {
		c.Data["json"] = map[string]interface{}{
			"Success": true,
			"Status":  "200",
			"Message": "Registro obtenido correctamente",
			"Data":    resultado,
		}
	}

	c.ServeJSON()
}

// Post inserta un nuevo registro llamando al servicio externo
// @Title Post
// @Description Inserta un nuevo registro enviándolo al servicio JBPM
// @Param   body        body    map[string]interface{}  true        "Datos del registro a crear. Debe coincidir con la estructura esperada por el servicio externo, ej: {'_posttercero_pago': {...}}"
// @Success 201 {object} map[string]interface{} "Registro creado y confirmado por el servicio externo (respuesta del servicio externo)"
// @Success 202 {object} map[string]interface{} "Solicitud aceptada para procesamiento por el servicio externo (usualmente sin cuerpo de respuesta)"
// @Failure 400 {object} map[string]interface{} "Error en los datos de entrada"
// @Failure 500 {object} map[string]interface{} "Error interno del servidor o error inesperado del servicio externo"
// @Failure default {object} map[string]interface{} "Respuesta de error del servicio externo"
// @router / [post]
func (c *FacturacionElectronicaController) Post() {
	defer errorhandler.HandlePanic(&c.Controller)
	var terceroPago models.TerceroPagoRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &terceroPago); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = requestresponse.APIResponseDTO(false, 400, nil, "Datos erroneos")
		c.ServeJSON()
		return
	}
	// Llamar al service
	response := services.GuardarDatosTerceroPago(terceroPago, terceroPago.TipoUsuario, terceroPago.IdTipoDocumentoDuenoRecibo, terceroPago.TerceroId)
	c.Ctx.Output.SetStatus(response.Status)
	c.Data["json"] = response
	c.ServeJSON()
}

// Put actualiza un registro existente llamando al servicio externo
// @Title Put
// @Description Actualiza un registro existente enviando los datos al servicio JBPM
// @Param   id          path    string                  true        "ID del registro a actualizar"
// @Param   body        body    map[string]interface{}  true        "Datos del registro a actualizar. Debe coincidir con la estructura esperada por el servicio externo."
// @Success 202 {object} map[string]interface{} "Solicitud de actualización aceptada para procesamiento por el servicio externo (sin cuerpo de respuesta)"
// @Success 200 {object} map[string]interface{} "Registro actualizado correctamente (respuesta opcional del servicio externo, usualmente vacía)"
// @Success 204 {object} map[string]interface{} "Registro actualizado correctamente (sin cuerpo de respuesta)"
// @Failure 400 {object} map[string]interface{} "Error en los datos de entrada o ID inválido"
// @Failure 404 {object} map[string]interface{} "Registro no encontrado (si el servicio externo lo indica)"
// @Failure 500 {object} map[string]interface{} "Error interno del servidor o error inesperado del servicio externo"
// @Failure default {object} map[string]interface{} "Respuesta de error del servicio externo"
// @router /:id [put]
func (c *FacturacionElectronicaController) Put() {
	id := c.Ctx.Input.Param(":id") // Obtener el ID de la URL
	var inputData map[string]interface{}

	// 1. Validar que el ID no esté vacío (opcional pero recomendado)
	if id == "" {
		c.Data["json"] = map[string]interface{}{
			"Success": false,
			"Status":  "400",
			"Message": "Error: Falta el ID del registro en la URL.",
			"Data":    nil,
		}
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.ServeJSON()
		return
	}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &inputData); err != nil {
		logs.Error("Error al parsear JSON de entrada para PUT:", err)
		c.Data["json"] = map[string]interface{}{
			"Success": false,
			"Status":  "400",
			"Message": "Error en el formato JSON de la solicitud: " + err.Error(),
			"Data":    nil,
		}
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.ServeJSON()
		return
	}

	serviceURL := "http://" + beego.AppConfig.String("FacturacionElectronicaService") + "/" + id

	req := httplib.Put(serviceURL)
	req.Header("Content-Type", "application/json")
	req.Header("Accept", "application/json")
	req.JSONBody(inputData)

	resp, err := req.Response()

	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"Success": false,
			"Status":  "503",
			"Message": "Error de comunicación con el servicio externo: " + err.Error(),
		}
		c.Ctx.Output.SetStatus(http.StatusServiceUnavailable)
		c.ServeJSON()
		return
	}

	defer resp.Body.Close()

	statusCode := resp.StatusCode
	statusString := resp.Status

	switch {
	case statusCode == http.StatusAccepted: // 202
		c.Data["json"] = map[string]interface{}{
			"Success": true, // ¡Éxito!
			"Status":  statusString,
			"Message": "Solicitud de actualización aceptada para procesamiento.",
		}
		c.Ctx.Output.SetStatus(statusCode) // 202

	case statusCode == http.StatusOK: // 200 (También éxito para PUT)
		c.Data["json"] = map[string]interface{}{
			"Success": true, // ¡Éxito!
			"Status":  statusString,
			"Message": "Registro actualizado correctamente por el servicio externo (200 OK).",
		}
		c.Ctx.Output.SetStatus(statusCode) // 200

	case statusCode == http.StatusNoContent:
		c.Data["json"] = map[string]interface{}{
			"Success": true, // ¡Éxito!
			"Status":  statusString,
			"Message": "Registro actualizado correctamente (sin contenido).",
		}
		c.Ctx.Output.SetStatus(statusCode) // 204

	default: // Cualquier otro código es un error
		c.Data["json"] = map[string]interface{}{
			"Success": false,
			"Status":  statusString,
			"Message": "Error reportado por el servicio externo durante la actualización.",
		}
		c.Ctx.Output.SetStatus(statusCode)
	}

	c.ServeJSON()
}
