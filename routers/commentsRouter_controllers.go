package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/sga_evaluacion_docente_mid/controllers:Formulario_por_tipoController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_evaluacion_docente_mid/controllers:Formulario_por_tipoController"],
        beego.ControllerComments{
            Method: "GetFormularioTipo",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_evaluacion_docente_mid/controllers:Respuesta_formularioController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_evaluacion_docente_mid/controllers:Respuesta_formularioController"],
        beego.ControllerComments{
            Method: "Post",
            Router: "/",
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_evaluacion_docente_mid/controllers:Respuesta_formularioController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_evaluacion_docente_mid/controllers:Respuesta_formularioController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: "/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_evaluacion_docente_mid/controllers:Respuesta_formularioController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_evaluacion_docente_mid/controllers:Respuesta_formularioController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: "/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_evaluacion_docente_mid/controllers:Respuesta_formularioController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_evaluacion_docente_mid/controllers:Respuesta_formularioController"],
        beego.ControllerComments{
            Method: "Put",
            Router: "/:id",
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/sga_evaluacion_docente_mid/controllers:Respuesta_formularioController"] = append(beego.GlobalControllerRouter["github.com/udistrital/sga_evaluacion_docente_mid/controllers:Respuesta_formularioController"],
        beego.ControllerComments{
            Method: "Delete",
            Router: "/:id",
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
