package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/udistrital/facturacion_electronica_mid/controllers:FacturacionElectronicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/facturacion_electronica_mid/controllers:FacturacionElectronicaController"],
		beego.ControllerComments{
			Method:           "GetAll",
			Router:           "/",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/facturacion_electronica_mid/controllers:FacturacionElectronicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/facturacion_electronica_mid/controllers:FacturacionElectronicaController"],
		beego.ControllerComments{
			Method:           "Post",
			Router:           "/",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/facturacion_electronica_mid/controllers:FacturacionElectronicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/facturacion_electronica_mid/controllers:FacturacionElectronicaController"],
		beego.ControllerComments{
			Method:           "Put",
			Router:           "/:id",
			AllowHTTPMethods: []string{"put"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/udistrital/facturacion_electronica_mid/controllers:FacturacionElectronicaController"] = append(beego.GlobalControllerRouter["github.com/udistrital/facturacion_electronica_mid/controllers:FacturacionElectronicaController"],
		beego.ControllerComments{
			Method:           "GetOne",
			Router:           "/:id/:anio",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
