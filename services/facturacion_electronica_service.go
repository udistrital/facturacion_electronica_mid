package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
	"github.com/udistrital/facturacion_electronica_mid/models"
	"github.com/udistrital/facturacion_electronica_mid/utils"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
)

func ObtenerRegistrosTercerosPago() requestresponse.APIResponse {
	var tercerosPagoResponse map[string]interface{}
	url := beego.AppConfig.String("BusserviciosBasePath") + beego.AppConfig.String("TerceroPagoPath")

	if err := request.GetJsonWSO2(url, &tercerosPagoResponse); err != nil {
		logs.Info("URL completa: %s", url)
		logs.Error("Error al obtener terceros pago: %v", err)
		return requestresponse.APIResponse{
			Success: false,
			Status:  http.StatusNotFound,
			Message: "Error al obtener los datos: " + err.Error(),
			Data:    nil,
		}
	}

	return requestresponse.APIResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Datos obtenidos correctamente",
		Data:    tercerosPagoResponse,
	}
}

func ObtenerTerceroPago(id, anio string) requestresponse.APIResponse {
	var terceroPagoResponse map[string]interface{}
	url := beego.AppConfig.String("BusserviciosBasePath") + beego.AppConfig.String("TerceroPagoPath") + "/" + id + "/" + anio

	if err := request.GetJsonWSO2(url, &terceroPagoResponse); err != nil {
		logs.Error("Error al obtener tercero pago %s/%s: %v", id, anio, err)
		return requestresponse.APIResponse{
			Success: false,
			Status:  http.StatusNotFound,
			Message: "Error al obtener el registro: " + err.Error(),
			Data:    nil,
		}
	}

	return requestresponse.APIResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Registro obtenido correctamente",
		Data:    terceroPagoResponse,
	}
}

func GuardarDatosTerceroPago(terceroPago models.TerceroPagoRequest, tipoUsuario int, idTipoDocumentoDuenoRecibo int, terceroDuenoId int) requestresponse.APIResponse {
	/* Tipos usuario
	1: aspirante
	2: admitido
	*/

	var duenoRecibo models.DuenoRecibo

	// 1. Mapeo de tipo de documento con el id tipo documento de terceros
	tipoDocumento, _ := utils.ObtenerTipoDocumentoSGA(terceroPago.IdTipoDocumentoDuenoRecibo)

	// 2. Obtener datos del dueño del recibo
	duenoRecibo, err := obtenerDatosDuenoRecibo(terceroPago, tipoUsuario, tipoDocumento)

	if err != nil {
		beego.Warning("GuardarDatosTerceroPago: error al obtener datos del dueño del recibo:", err)
		return requestresponse.APIResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Error al obtener datos del dueño del recibo: " + err.Error(),
			Data:    nil,
		}
	}

	// 3. Obtener datos de los conceptos del recibo
	conceptosRecibo, err := obtenerDatosConceptosRecibo(terceroPago, tipoUsuario)

	if err != nil {
		beego.Warning("GuardarDatosTerceroPago: error al obtener conceptos del recibo:", err)
		return requestresponse.APIResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Error al obtener conceptos del recibo: " + err.Error(),
			Data:    nil,
		}
	}

	// 4. Se arma el array de los json de datos adicionales de ACTERCERO_PAGO
	datosAdicionales, err := armarDatosAdicionalesPorConcepto(duenoRecibo, conceptosRecibo)
	if err != nil {
		beego.Warning("GuardarDatosTerceroPago: error al armar datos adicionales:", err)
		return requestresponse.APIResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Error al armar datos adicionales: " + err.Error(),
			Data:    nil,
		}
	}

	// 5. Crear un array de TerceroPago, uno por cada dato adicional, y enviarlos a ACTERCERO_PAGO
	serviceActerceroUrl := beego.AppConfig.String("BusserviciosBasePath") + beego.AppConfig.String("TerceroPagoPath")
	var respuestas []interface{}
	var errores []string

	for index, datoAdicional := range datosAdicionales {
		// Crear una copia del terceroPago original
		terceroPagoCopia := terceroPago

		// Convertir el dato adicional individual a JSON string
		datoAdicionalJSON, err := json.Marshal(datoAdicional)
		if err != nil {
			beego.Warning("GuardarDatosTerceroPago: error al convertir dato adicional a JSON:", err)
			errores = append(errores, fmt.Sprintf("Concepto %d: error al convertir dato adicional a JSON: %v", index+1, err))
			continue
		}

		terceroPagoCopia.PostTerceroPago.TERPA_DATOS_ADICIONALES = string(datoAdicionalJSON)

		// Enviar a ACTERCERO_PAGO
		respuesta, err := enviarTerceroOra(terceroPagoCopia, serviceActerceroUrl)
		if err != nil {
			beego.Warning("GuardarDatosTerceroPago: error al enviar concepto %d: %v", index+1, err)
			errores = append(errores, fmt.Sprintf("Concepto %d: %v", index+1, err))
			continue
		}

		respuestas = append(respuestas, respuesta)
	}

	// Evaluar el resultado general
	totalConceptos := len(datosAdicionales)
	exitosos := len(respuestas)
	fallidos := len(errores)

	if fallidos == totalConceptos {
		// Todos fallaron
		return requestresponse.APIResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: fmt.Sprintf("Todos los conceptos fallaron al enviarse: %v", errores),
			Data:    nil,
		}
	} else if fallidos > 0 {
		// Algunos fallaron
		return requestresponse.APIResponse{
			Success: true,
			Status:  http.StatusPartialContent,
			Message: fmt.Sprintf("Se enviaron %d/%d conceptos correctamente. Errores: %v", exitosos, totalConceptos, errores),
			Data: map[string]interface{}{
				"respuestas_exitosas": respuestas,
				"errores":             errores,
			},
		}
	}

	// Se hace el envio al ERP

	respuestasSofia, err := enviarDatosSofia(terceroPago.PostTerceroPago, duenoRecibo, conceptosRecibo, terceroDuenoId)

	fmt.Println(respuestasSofia, err)

	if err != nil {
		beego.Warning("GuardarDatosERP: error al enviar datos a ERP", err)
	}

	// Guardado de las respuestas de ERP en ACTERCERO_PAGO

	terceroPagoRespuestas := terceroPago

	orderedKeys := []string{"respSofiaTerceroP", "respSofiaTerceroD"}
	for i := 1; i <= len(conceptosRecibo); i++ {
		orderedKeys = append(orderedKeys, fmt.Sprintf("respSofiaConcepto%d", i))
	}

	jsonRespuestasERP := make(models.RespuestasERP)

	for _, respuesta := range respuestasSofia {
		var key string

		message, ok := respuesta.Message.(string)
		if !ok {
			logs.Error("respuesta.Message no es de tipo string")
			continue
		}

		// Use a tagged switch to determine the key
		switch message {
		case "respSofiaTerceroP":
			key = "respSofiaTerceroP"
		case "respSofiaTerceroD":
			key = "respSofiaTerceroD"
		default:
			key = message
		}

		jsonRespuestasERP[key] = respuesta.Status
	}

	jsonData := "{"
	for i, k := range orderedKeys {
		if val, ok := jsonRespuestasERP[k]; ok {
			if i > 0 {
				jsonData += ","
			}
			jsonData += fmt.Sprintf("\"%s\":%d", k, val)
		}
	}
	jsonData += "}"

	terceroPagoRespuestas.PostTerceroPago.TERPA_DATOS_ADICIONALES = string(jsonData)

	// Enviar respuestas de ERP a ACTERCERO_PAGO
	respuestaEnvio, err := enviarTerceroOra(terceroPagoRespuestas, serviceActerceroUrl)
	if err != nil {
		beego.Warning("GuardarDatosTerceroPago: error al enviar respuestad de ERP", err)
	}

	respuestas = append(respuestas, respuestaEnvio)

	return requestresponse.APIResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Se guardaron los datos en ACTERCERO_PAGO correctamente",
		Data:    respuestas,
	}
}

func obtenerDatosDuenoRecibo(terceroPago models.TerceroPagoRequest, tipoUsuario int, tipoDocumento string) (models.DuenoRecibo, error) {
	/* Tipos usuario
	1: aspirante
	2: admitido
	*/
	// Consulta a servicio de recibos para obtener los datos del dueno del recibo

	var duenoResponse models.DuenoReciboResponse
	urlDueno := beego.AppConfig.String("BusserviciosBasePath") + beego.AppConfig.String("ConsultarReciboPath") + "datos_recibo/" + tipoDocumento + "/" + strconv.Itoa(tipoUsuario) + "/" + strconv.Itoa(terceroPago.PostTerceroPago.TERPA_ANO_PAGO) + "/" + strconv.Itoa(terceroPago.PostTerceroPago.TERPA_SECUENCIA)

	if err := request.GetJsonWSO2(urlDueno, &duenoResponse); err != nil {
		logs.Error("No se pudo obtener los datos del dueno del recibo %s / %s: %v", strconv.Itoa(terceroPago.PostTerceroPago.TERPA_SECUENCIA), strconv.Itoa(terceroPago.PostTerceroPago.TERPA_ANO_PAGO), err)
		return models.DuenoRecibo{}, err
	}

	if len(duenoResponse.ReciboCollection.Recibo) == 0 {
		logs.Error("No se encontraron datos del dueño del recibo")
		return models.DuenoRecibo{}, fmt.Errorf("no se encontraron datos del dueño del recibo")
	}

	return duenoResponse.ReciboCollection.Recibo[0], nil
}

func obtenerDatosConceptosRecibo(terceroPago models.TerceroPagoRequest, tipoUsuario int) ([]models.ConceptoRecibo, error) {
	// Consulta a servicio de recibos para obtener los conceptos de un recibo
	var conceptosResponse models.ConceptosReciboResponse
	urlConceptos := beego.AppConfig.String("BusserviciosBasePath") + beego.AppConfig.String("ConsultarReciboPath") + "datos_conceptos_recibo/" + strconv.Itoa(terceroPago.PostTerceroPago.TERPA_ANO_PAGO) + "/" + strconv.Itoa(terceroPago.PostTerceroPago.TERPA_SECUENCIA) + "/" + strconv.Itoa(tipoUsuario)

	if err := request.GetJsonWSO2(urlConceptos, &conceptosResponse); err != nil {
		logs.Error("No se pudo obtener los conceptos del recibo %s / %s: %v", strconv.Itoa(terceroPago.PostTerceroPago.TERPA_SECUENCIA), strconv.Itoa(terceroPago.PostTerceroPago.TERPA_ANO_PAGO), err)
		return []models.ConceptoRecibo{}, err
	}

	return conceptosResponse.ReciboCollection.Recibo, nil
}

func enviarDatosSofia(pagador models.TerceroPago, dueno models.DuenoRecibo, conceptosRecibo []models.ConceptoRecibo, terceroDuenoId int) ([]requestresponse.APIResponse, error) {

	var respuestasERP []requestresponse.APIResponse

	SofiaTerceroD, SofiaTerceroP, SofiaTerceroConceptos, err := utils.GenerarDatosSofia(pagador, dueno, conceptosRecibo, terceroDuenoId)
	if err != nil {
		return nil, err
	}

	var datosSofiaPost models.DatosSofiaPost

	datosSofiaPost.TerceroD = SofiaTerceroD
	datosSofiaPost.TerceroP = SofiaTerceroP
	datosSofiaPost.ConceptoList = SofiaTerceroConceptos

	// Enviar datos a ERP

	client := &http.Client{}
	urlSofia := beego.AppConfig.String("SofiaService")

	//PAGADOR
	jsonDataTerceroP, err := json.Marshal(datosSofiaPost.TerceroP)
	if err != nil {
		logs.Error("Error al convertir datosSofiaPost a JSON: %v", err)
		// return nil, err
	}

	// Configurar la solicitud HTTP Pagador
	requestPagador, err := http.NewRequest("POST", urlSofia, bytes.NewBuffer(jsonDataTerceroP))
	if err != nil {
		logs.Error("Error al crear la solicitud HTTP: %v", err)
		return nil, err
	}
	requestPagador.Header.Set("Content-Type", "application/json")

	// Enviar la solicitud HTTP pagador y guardado de la respuesta
	responsePagador, err := client.Do(requestPagador)
	if err != nil {
		logs.Error("Error al enviar la solicitud HTTP pagador a Sofia (%s): %v", urlSofia, err)
		respuestasERP = append(respuestasERP, requestresponse.APIResponse{
			Success: false,
			Status:  http.StatusServiceUnavailable,
			Message: "respSofiaTerceroP",
			Data:    "",
		})
	} else {
		respuestasERP = append(respuestasERP, requestresponse.APIResponse{
			Success: true,
			Status:  responsePagador.StatusCode,
			Message: "respSofiaTerceroP",
			Data:    responsePagador.Body,
		})
		defer responsePagador.Body.Close()
	}

	// DUENO
	jsonDataTerceroD, err := json.Marshal(datosSofiaPost.TerceroD)
	if err != nil {
		logs.Error("Error al convertir datosSofiaPost a JSON: %v", err)
		return nil, err
	}

	requestDueno, err := http.NewRequest("POST", urlSofia, bytes.NewBuffer(jsonDataTerceroD))
	if err != nil {
		logs.Error("Error al crear la solicitud HTTP: %v", err)
		return nil, err
	}
	requestDueno.Header.Set("Content-Type", "application/json")

	// Enviar la solicitud HTTP dueno recibo y guardado de la respuesta
	responseDueno, err := client.Do(requestDueno)
	if err != nil {
		logs.Error("Error al enviar la solicitud HTTP dueño a Sofia (%s): %v", urlSofia, err)
		respuestasERP = append(respuestasERP, requestresponse.APIResponse{
			Success: false,
			Status:  http.StatusServiceUnavailable,
			Message: "respSofiaTerceroD",
			Data:    "",
		})
	} else {
		respuestasERP = append(respuestasERP, requestresponse.APIResponse{
			Success: true,
			Status:  responseDueno.StatusCode,
			Message: "respSofiaTerceroD",
			Data:    responseDueno.Body,
		})
		defer responseDueno.Body.Close()
	}

	// CONCEPTOS
	for i, concepto := range datosSofiaPost.ConceptoList {
		jsonDataConcepto, err := json.Marshal(concepto)
		if err != nil {
			logs.Error("Error al convertir datosSofiaPost a JSON: %v", err)
			// return nil, err
		}

		// Configurar la solicitud HTTP Dueno recibo
		requestConcepto, err := http.NewRequest("POST", urlSofia, bytes.NewBuffer(jsonDataConcepto))
		if err != nil {
			logs.Error("Error al crear la solicitud HTTP: %v", err)
			return nil, err
		}
		requestConcepto.Header.Set("Content-Type", "application/json")

		// Enviar la solicitud HTTP dueno recibo y guardado de la respuesta
		responseConcepto, err := client.Do(requestConcepto)
		if err != nil {
			logs.Error("Error al enviar la solicitud HTTP concepto %d a Sofia (%s): %v", i+1, urlSofia, err)
			respuestasERP = append(respuestasERP, requestresponse.APIResponse{
				Success: false,
				Status:  http.StatusServiceUnavailable,
				Message: fmt.Sprintf("respSofiaConcepto%d", i+1),
				Data:    "",
			})
		} else {
			respuestasERP = append(respuestasERP, requestresponse.APIResponse{
				Success: true,
				Status:  responseConcepto.StatusCode,
				Message: fmt.Sprintf("respSofiaConcepto%d", i+1),
				Data:    responseConcepto.Body,
			})
			defer responseConcepto.Body.Close()
		}
	}

	return respuestasERP, nil

}

func enviarTerceroOra(terceroPago models.TerceroPagoRequest, serviceURL string) (interface{}, error) {
	// Se ajusta la fecha de creacion del registro a la fecha actual en zona horaria de Bogotá
	// Formato: DD/MM/YYYY HH24:MI:SS
	bogotaLocation, _ := time.LoadLocation("America/Bogota")
	terceroPago.PostTerceroPago.TERPA_FECHA_REGISTRO = time.Now().In(bogotaLocation).Format("02/01/2006 15:04:05")

	// Crear el objeto a enviar con el wrapper _posttercero_pago
	payload := map[string]interface{}{
		"_posttercero_pago": terceroPago.PostTerceroPago,
	}

	req := httplib.Post(serviceURL)
	req.Header("Content-Type", "application/json")
	req.Header("Accept", "application/json")
	req.JSONBody(payload)

	resp, err := req.Response()
	if err != nil {
		return nil, fmt.Errorf("error wso2: %v", err)
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	var serviceResponse interface{}

	switch {
	case statusCode == http.StatusCreated: // 201
		err = json.NewDecoder(resp.Body).Decode(&serviceResponse)
		if err != nil {
			return map[string]interface{}{
				"status":  statusCode,
				"message": "Registro creado, pero respuesta no pudo ser interpretada",
			}, nil
		}
		return serviceResponse, nil

	case statusCode == http.StatusAccepted: // 202
		return map[string]interface{}{
			"status":  statusCode,
			"message": "Solicitud aceptada para procesamiento",
		}, nil

	case statusCode == http.StatusOK: // 200
		err = json.NewDecoder(resp.Body).Decode(&serviceResponse)
		if err != nil {
			return map[string]interface{}{
				"status":  statusCode,
				"message": "Solicitud procesada (200 OK), pero respuesta no pudo ser interpretada",
			}, nil
		}
		return serviceResponse, nil

	default: // Cualquier otro código (>=300) es un error
		return nil, fmt.Errorf("error wso2: servicio externo retornó código de error: %d", statusCode)
	}
}

// ActualizarDatosTerceroPago envía la solicitud PUT hacia ACTERCERO_PAGO y estandariza la respuesta
func ActualizarDatosTerceroPago(id string, requestBody []byte) requestresponse.APIResponse {
	if id == "" {
		return requestresponse.APIResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Error: Falta el ID del registro en la URL.",
			Data:    nil,
		}
	}

	var inputData map[string]interface{}
	if err := json.Unmarshal(requestBody, &inputData); err != nil {
		logs.Error("Error al parsear JSON de entrada para PUT: %v", err)
		return requestresponse.APIResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Error en el formato JSON de la solicitud: " + err.Error(),
			Data:    nil,
		}
	}

	serviceURL := beego.AppConfig.String("BusserviciosBasePath") + beego.AppConfig.String("TerceroPagoPath") + "/" + id
	req := httplib.Put(serviceURL)
	req.Header("Content-Type", "application/json")
	req.Header("Accept", "application/json")
	req.JSONBody(inputData)

	resp, err := req.Response()
	if err != nil {
		return requestresponse.APIResponse{
			Success: false,
			Status:  http.StatusServiceUnavailable,
			Message: "Error de comunicación con el servicio externo: " + err.Error(),
			Data:    nil,
		}
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusAccepted:
		return requestresponse.APIResponse{
			Success: true,
			Status:  http.StatusAccepted,
			Message: "Solicitud de actualización aceptada para procesamiento.",
			Data:    nil,
		}
	case http.StatusOK:
		return requestresponse.APIResponse{
			Success: true,
			Status:  http.StatusOK,
			Message: "Registro actualizado correctamente por el servicio externo (200 OK).",
			Data:    nil,
		}
	case http.StatusNoContent:
		return requestresponse.APIResponse{
			Success: true,
			Status:  http.StatusNoContent,
			Message: "Registro actualizado correctamente (sin contenido).",
			Data:    nil,
		}
	default:
		return requestresponse.APIResponse{
			Success: false,
			Status:  resp.StatusCode,
			Message: "Error reportado por el servicio externo durante la actualización.",
			Data:    nil,
		}
	}
}

// armarDatosAdicionalesPorConcepto crea un JSON por cada concepto del recibo
// combinando los datos del dueño del recibo con cada concepto individual
func armarDatosAdicionalesPorConcepto(duenoRecibo models.DuenoRecibo, conceptosRecibo []models.ConceptoRecibo) ([]models.DatosAdicionales, error) {
	var datosAdicionales []models.DatosAdicionales

	// Calcular cantidad de conceptos y valor total
	cantidadConceptos := len(conceptosRecibo)
	var valorTotal float64

	for _, concepto := range conceptosRecibo {
		valor, err := strconv.ParseFloat(concepto.Valor, 64)
		if err != nil {
			logs.Error("Error al convertir valor del concepto %s: %v", concepto.CodConcepto, err)
			return nil, fmt.Errorf("error al convertir valor del concepto: %v", err)
		}
		valorTotal += valor
	}

	// Crear un JSON por cada concepto
	for index, concepto := range conceptosRecibo {
		// Convertir identificacion a int (0 si está vacío)
		identificacion := 0
		if duenoRecibo.Identificacion != "" {
			var err error
			identificacion, err = strconv.Atoi(duenoRecibo.Identificacion)
			if err != nil {
				logs.Error("Error al convertir identificacion: %v", err)
				return nil, fmt.Errorf("error al convertir identificacion: %v", err)
			}
		}

		// Convertir cod_estudiante a int (0 si está vacío)
		codEstudiante := 0
		if duenoRecibo.CodEstudiante != "" {
			var err error
			codEstudiante, err = strconv.Atoi(duenoRecibo.CodEstudiante)
			if err != nil {
				logs.Error("Error al convertir cod_estudiante: %v", err)
				return nil, fmt.Errorf("error al convertir cod_estudiante: %v", err)
			}
		}

		// Convertir cod_carrera a int (0 si está vacío)
		codCarrera := 0
		if duenoRecibo.CodCarrera != "" {
			var err error
			codCarrera, err = strconv.Atoi(duenoRecibo.CodCarrera)
			if err != nil {
				logs.Error("Error al convertir cod_carrera: %v", err)
				return nil, fmt.Errorf("error al convertir cod_carrera: %v", err)
			}
		}

		// Convertir cod_concepto a int (0 si está vacío)
		codConcepto := 0
		if concepto.CodConcepto != "" {
			var err error
			codConcepto, err = strconv.Atoi(concepto.CodConcepto)
			if err != nil {
				logs.Error("Error al convertir cod_concepto: %v", err)
				return nil, fmt.Errorf("error al convertir cod_concepto: %v", err)
			}
		}

		// Convertir valor del concepto a float64 (0 si está vacío)
		valorConcepto := 0.0
		if concepto.Valor != "" {
			var err error
			valorConcepto, err = strconv.ParseFloat(concepto.Valor, 64)
			if err != nil {
				logs.Error("Error al convertir valor del concepto: %v", err)
				return nil, fmt.Errorf("error al convertir valor del concepto: %v", err)
			}
		}

		// Armar el modelo DatosAdicionales para este concepto
		datosConcepto := models.DatosAdicionales{
			Identificacion:        identificacion,
			CodTipoIdentificacion: duenoRecibo.CodTipoIdentificacion,
			Nombre:                duenoRecibo.Nombre,
			CorreoElectronico:     duenoRecibo.CorreoElectronico,
			CodEstudiante:         codEstudiante,
			CodCarrera:            codCarrera,
			Carrera:               duenoRecibo.Carrera,
			CodConcepto:           codConcepto,
			Concepto:              concepto.Concepto,
			NumeroConcepto:        index + 1,
			Valor:                 valorConcepto,
			CantidadConceptos:     cantidadConceptos,
			ValorTotal:            valorTotal,
			Nivel:                 duenoRecibo.Nivel,
		}

		datosAdicionales = append(datosAdicionales, datosConcepto)
	}

	return datosAdicionales, nil
}
