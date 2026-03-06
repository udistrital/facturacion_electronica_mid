package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/udistrital/facturacion_electronica_mid/models"
	"github.com/udistrital/facturacion_electronica_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
)

type TerceroPagoController struct {
	beego.Controller
}

func (c *TerceroPagoController) URLMapping() {
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
func (c *TerceroPagoController) GetAll() {
	response := services.ObtenerRegistrosTercerosPago()
	c.Ctx.Output.SetStatus(response.Status)
	c.Data["json"] = response
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
func (c *TerceroPagoController) GetOne() {
	response := services.ObtenerTerceroPago(c.Ctx.Input.Param(":id"), c.Ctx.Input.Param(":anio"))
	c.Ctx.Output.SetStatus(response.Status)
	c.Data["json"] = response
	c.ServeJSON()
}

// Post inserta un nuevo registro llamando al servicio de terceros posible pagador
// @Title Post
// @Description Inserta un nuevo registro enviándolo al servicio JBPM
// @Param   body        body    map[string]interface{}  true        "Datos del registro a crear. Debe coincidir con la estructura esperada por el servicio externo, ej: {'_posttercero_pago': {...}}"
// @Success 201 {object} map[string]interface{} "Registro creado y confirmado por el servicio externo (respuesta del servicio externo)"
// @Success 202 {object} map[string]interface{} "Solicitud aceptada para procesamiento por el servicio externo (usualmente sin cuerpo de respuesta)"
// @Failure 400 {object} map[string]interface{} "Error en los datos de entrada"
// @Failure 500 {object} map[string]interface{} "Error interno del servidor o error inesperado del servicio externo"
// @Failure default {object} map[string]interface{} "Respuesta de error del servicio externo"
// @router / [post]
func (c *TerceroPagoController) Post() {
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

// Put actualiza un registro existente de tercero posible pagador llamando al servicio JBPM
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
func (c *TerceroPagoController) Put() {
	response := services.ActualizarDatosTerceroPago(c.Ctx.Input.Param(":id"), c.Ctx.Input.RequestBody)
	c.Ctx.Output.SetStatus(response.Status)
	c.Data["json"] = response
	c.ServeJSON()
}
