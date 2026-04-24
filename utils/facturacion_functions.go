package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/udistrital/facturacion_electronica_mid/models"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/time_bogota"
)

func parseStringIntOrZero(value string, fieldName string) int {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0
	}

	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		beego.Error("Error convirtiendo "+fieldName+" a int:", err)
		return 0
	}

	return parsedValue
}

func GenerarDatosSofia(terceroPago models.TerceroPago, duenoRecibo models.DuenoRecibo, conceptos []models.ConceptoRecibo, terceroDuenoId int) (sofiaTerceroD models.DatosSofiaDueno, sofiaTerceroP models.DatosSofiaPagador, sofiaTerceroConceptos []models.DatosSofiaConcepto, err error) {
	defer func() {
		if r := recover(); r != nil {
			beego.Error("Error utils.GenerarDatosSofia. ", err)
			return
		}
	}()

	var terceroDuenoData models.TerceroId

	now := time_bogota.Tiempo_bogota()

	formattedDate := now.Format("02/01/2006")

	terceroDuenoData, err = obtenerNombresDueno(terceroDuenoId)
	if err != nil {
		panic("Error obteniendo datos tercero Dueno " + err.Error())
	}

	sofiaTerceroP = make(models.DatosSofiaPagador)               // llaveConexion +  33 campos
	sofiaTerceroD = make(models.DatosSofiaDueno)                 // llaveConexion + 33 campos
	sofiaTerceroConceptos = make([]models.DatosSofiaConcepto, 0) // cada elemento de  laveConexion + 19 campos

	// Llenar SofiaTerceroP con datos de posiblePagador

	sofiaTerceroP["llaveConexion"] = "SofiaDS"
	sofiaTerceroP["campo1"] = ""
	sofiaTerceroP["campo2"] = "TER"
	sofiaTerceroP["campo3"] = 0
	sofiaTerceroP["campo4"] = terceroPago.TERPA_NATURALEZA
	sofiaTerceroP["campo5"] = terceroPago.TERPA_TDO_CODVAR
	sofiaTerceroP["campo6"] = terceroPago.TERPA_NRO_DOCUMENTO
	sofiaTerceroP["campo7"] = terceroPago.TERPA_DIGITO_CHEQUEO
	sofiaTerceroP["campo8"] = terceroPago.TERPA_RAZON_SOCIAL
	sofiaTerceroP["campo9"] = strings.ToUpper(strings.TrimSpace(terceroPago.TERPA_PRIMER_APELLIDO))
	sofiaTerceroP["campo10"] = strings.ToUpper(strings.TrimSpace(terceroPago.TERPA_SEGUNDO_APELLIDO))
	sofiaTerceroP["campo11"] = strings.ToUpper(strings.TrimSpace(terceroPago.TERPA_PRIMER_NOMBRE))
	sofiaTerceroP["campo12"] = strings.ToUpper(strings.TrimSpace(terceroPago.TERPA_SEGUNDO_NOMBRE))
	sofiaTerceroP["campo13"] = terceroPago.TERPA_DIRECCION
	sofiaTerceroP["campo14"] = terceroPago.TERPA_TELEFONO
	sofiaTerceroP["campo15"] = terceroPago.TERPA_EMAIL
	sofiaTerceroP["campo16"] = "CR13"
	sofiaTerceroP["campo17"] = formattedDate
	sofiaTerceroP["campo18"] = terceroPago.TERPA_ESTADO_REGISTRO
	sofiaTerceroP["campo19"] = 11001
	sofiaTerceroP["campo20"] = 23
	sofiaTerceroP["campo21"] = ""
	sofiaTerceroP["campo22"] = ""
	sofiaTerceroP["campo23"] = terceroPago.TERPA_EMAIL
	if terceroPago.TERPA_NATURALEZA == "N" {
		sofiaTerceroP["campo24"] = "P"
	} else {
		sofiaTerceroP["campo24"] = ""
	}
	if terceroPago.TERPA_NATURALEZA == "N" {
		sofiaTerceroP["campo25"] = "N"
	} else {
		sofiaTerceroP["campo25"] = "S"
	}
	if terceroPago.TERPA_NATURALEZA == "N" {
		sofiaTerceroP["campo26"] = "N"
	} else {
		sofiaTerceroP["campo26"] = ""
	}
	if terceroPago.TERPA_NATURALEZA == "N" {
		sofiaTerceroP["campo27"] = "N"
	} else {
		sofiaTerceroP["campo27"] = ""
	}
	if terceroPago.TERPA_NATURALEZA == "N" {
		sofiaTerceroP["campo28"] = "RI06"
	} else {
		sofiaTerceroP["campo28"] = ""
	}
	if terceroPago.TERPA_NATURALEZA == "N" {
		sofiaTerceroP["campo29"] = "N/A"
	} else {
		sofiaTerceroP["campo29"] = ""
	}
	if terceroPago.TERPA_NATURALEZA == "N" {
		sofiaTerceroP["campo30"] = "N/A"
	} else {
		sofiaTerceroP["campo30"] = ""
	}
	if terceroPago.TERPA_NATURALEZA == "N" {
		sofiaTerceroP["campo31"] = "N/A"
	} else {
		sofiaTerceroP["campo31"] = ""
	}
	if terceroPago.TERPA_NATURALEZA == "N" {
		sofiaTerceroP["campo32"] = "N/A"
	} else {
		sofiaTerceroP["campo32"] = ""
	}
	if terceroPago.TERPA_NATURALEZA == "N" {
		sofiaTerceroP["campo33"] = "N/A"
	} else {
		sofiaTerceroP["campo33"] = ""
	}

	// Llenar SofiaTerceroD con datos de duenoRecibo

	sofiaTerceroD["llaveConexion"] = "SofiaDS"
	sofiaTerceroD["campo1"] = ""
	sofiaTerceroD["campo2"] = "TER"
	sofiaTerceroD["campo3"] = 0
	sofiaTerceroD["campo4"] = "N"
	sofiaTerceroD["campo5"] = duenoRecibo.TipoIdentificacion
	val6, errConv6 := strconv.Atoi(duenoRecibo.Identificacion)
	if errConv6 != nil {
		beego.Error("Error convirtiendo duenoRecibo.Identificacion a int:", errConv6)
		sofiaTerceroD["campo6"] = 0
	} else {
		sofiaTerceroD["campo6"] = val6
	}
	sofiaTerceroD["campo7"] = ""
	sofiaTerceroD["campo8"] = ""
	sofiaTerceroD["campo9"] = strings.ToUpper(strings.TrimSpace(terceroDuenoData.PrimerApellido))
	sofiaTerceroD["campo10"] = strings.ToUpper(strings.TrimSpace(terceroDuenoData.SegundoApellido))
	sofiaTerceroD["campo11"] = strings.ToUpper(strings.TrimSpace(terceroDuenoData.PrimerNombre))
	sofiaTerceroD["campo12"] = strings.ToUpper(strings.TrimSpace(terceroDuenoData.SegundoNombre))
	sofiaTerceroD["campo13"] = ""
	sofiaTerceroD["campo14"] = ""
	sofiaTerceroD["campo15"] = duenoRecibo.CorreoElectronico
	sofiaTerceroD["campo16"] = "CR13"
	sofiaTerceroD["campo17"] = formattedDate
	sofiaTerceroD["campo18"] = "A"
	sofiaTerceroD["campo19"] = 11001
	sofiaTerceroD["campo20"] = 23
	sofiaTerceroD["campo21"] = ""
	sofiaTerceroD["campo22"] = ""
	sofiaTerceroD["campo23"] = duenoRecibo.CorreoElectronico
	sofiaTerceroD["campo24"] = "P"
	sofiaTerceroD["campo25"] = "N"
	sofiaTerceroD["campo26"] = "N"
	sofiaTerceroD["campo27"] = "N"
	sofiaTerceroD["campo28"] = "RI06"
	sofiaTerceroD["campo29"] = "N/A"
	sofiaTerceroD["campo30"] = "N/A"
	sofiaTerceroD["campo31"] = "N/A"
	sofiaTerceroD["campo32"] = "N/A"
	sofiaTerceroD["campo33"] = "N/A"

	// Llenar SofiaTerceroConceptos con datos de conceptos

	// 1.) Calcular el valor total de los conceptos para asignarlo al campo correspondiente del objeto SofiaTerceroConcepto
	var valorTotalConceptos int
	for _, concepto := range conceptos {
		valor, err := strconv.ParseInt(concepto.Valor, 10, 64)
		if err != nil {
			beego.Error("Error convirtiendo concepto.Valor a int64:", err)
			continue
		}
		valorTotalConceptos += int(valor)
	}

	for i, concepto := range conceptos {
		SofiaTerceroConcepto := make(models.DatosSofiaConcepto) // Crear una nueva instancia

		// Asignar valores específicos a SofiaTerceroConcepto basados en el concepto
		// Ejemplo: SofiaTerceroConcepto["campo"] = concepto.AlgunCampo
		SofiaTerceroConcepto["llaveConexion"] = "SofiaDS"
		SofiaTerceroConcepto["campo1"] = ""
		SofiaTerceroConcepto["campo2"] = "REC"
		SofiaTerceroConcepto["campo3"] = terceroPago.TERPA_SECUENCIA
		SofiaTerceroConcepto["campo4"] = terceroPago.TERPA_ANO_PAGO
		SofiaTerceroConcepto["campo5"] = formattedDate
		SofiaTerceroConcepto["campo6"] = terceroPago.TERPA_TDO_CODVAR
		SofiaTerceroConcepto["campo7"] = terceroPago.TERPA_NRO_DOCUMENTO
		SofiaTerceroConcepto["campo8"] = valorTotalConceptos
		SofiaTerceroConcepto["campo9"] = i + 1 // Número secuencial del concepto
		SofiaTerceroConcepto["campo10"] = parseStringIntOrZero(concepto.CodConcepto, "concepto.CodConcepto")
		SofiaTerceroConcepto["campo11"] = strings.TrimSpace(concepto.Concepto)
		SofiaTerceroConcepto["campo12"] = parseStringIntOrZero(concepto.Valor, "concepto.Valor")
		SofiaTerceroConcepto["campo13"] = len(conceptos)
		SofiaTerceroConcepto["campo14"] = parseStringIntOrZero(duenoRecibo.CodCarrera, "duenoRecibo.CodCarrera")
		SofiaTerceroConcepto["campo15"] = duenoRecibo.Carrera
		SofiaTerceroConcepto["campo16"] = parseStringIntOrZero(duenoRecibo.CodEstudiante, "duenoRecibo.CodEstudiante")
		SofiaTerceroConcepto["campo17"] = duenoRecibo.CodTipoIdentificacion
		SofiaTerceroConcepto["campo18"] = parseStringIntOrZero(duenoRecibo.Identificacion, "duenoRecibo.Identificacion")
		SofiaTerceroConcepto["campo19"] = duenoRecibo.Nivel

		sofiaTerceroConceptos = append(sofiaTerceroConceptos, SofiaTerceroConcepto)

	}

	return sofiaTerceroD, sofiaTerceroP, sofiaTerceroConceptos, nil
}

func obtenerNombresDueno(terceroDuenoId int) (terceroDuenoData models.TerceroId, err error) {

	url := beego.AppConfig.String("TercerosService") + fmt.Sprintf("/tercero/%d", terceroDuenoId)

	err = request.GetJson(url, &terceroDuenoData)
	if err != nil {
		return models.TerceroId{}, err
	}

	return terceroDuenoData, nil
}
