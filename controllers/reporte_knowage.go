package controllers

import (
	"github.com/astaxie/beego"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"fmt"
	"strings"
	"io"
)

// ReporteKnowageController operations for ReporteKnowage
type ReporteKnowageController struct {
	beego.Controller
}

// URLMapping ...
func (c *ReporteKnowageController) URLMapping() {
	c.Mapping("Get", c.Get)
}

// @router / [get]
// @Title Get reporte Knowage
// @Description Obtiene un reporte desde Knowage basado en parámetros proporcionados
// @Param   label	query	string	true	"Nombre del reporte"
// @Success 200 {object} interface{} "Resultado del reporte"
// @Failure 400 parámetros inválidos
// @Failure 500 error al ejecutar la consulta
func (c *ReporteKnowageController) Get() {
	// Obtener el parámetro 'label' desde la URL
    label := c.GetString("label")
    if label == "" {
        c.Ctx.Output.SetStatus(400)
        c.Ctx.Output.Body([]byte("Falta el parámetro 'label'"))
        return
    }

    // Configuración de Knowage
    knowageHost := "https://inteligenciainstitucional.portaloas.udistrital.edu.co"
    contextPath := "/knowage"
    user := "de"         // ⚠️ Reemplaza con variables de entorno seguras
    password := "de"
    role := "/spagobi/user"

    // Crear cliente con soporte de cookies
    jar, _ := cookiejar.New(nil)
    client := &http.Client{Jar: jar}

    // Paso 1: Autenticación
    loginURL := fmt.Sprintf("%s%s/servlet/AdapterHTTP", knowageHost, contextPath)
    loginData := url.Values{}
    loginData.Set("NEW_SESSION", "TRUE")
    loginData.Set("USER", user)
    loginData.Set("PASSWORD", password)
    loginData.Set("ROLE", role)

    loginReq, _ := http.NewRequest("POST", loginURL, strings.NewReader(loginData.Encode()))
    loginReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    _, err := client.Do(loginReq)
    if err != nil {
        c.Ctx.Output.SetStatus(500)
        c.Ctx.Output.Body([]byte("Error al autenticar con Knowage"))
        return
    }

    // Paso 2: Construir URL del reporte
    reporteURL := fmt.Sprintf(
        "%s%s/servlet/AdapterHTTP?ACTION_NAME=EXECUTE_DOCUMENT_ANGULAR_ACTION&OBJECT_LABEL=%s&NEW_SESSION=TRUE&USER=%s&PASSWORD=%s&ROLE=%s&TOOLBAR_VISIBLE=true",
        knowageHost, contextPath, url.QueryEscape(label), url.QueryEscape(user), url.QueryEscape(password), url.QueryEscape(role),
    )

    // Paso 3: Obtener reporte
    req, _ := http.NewRequest("GET", reporteURL, nil)
    resp, err := client.Do(req)
    if err != nil {
        c.Ctx.Output.SetStatus(500)
        c.Ctx.Output.Body([]byte("Error al obtener el reporte desde Knowage"))
        return
    }
    defer resp.Body.Close()

    // Copiar encabezados necesarios
    for k, v := range resp.Header {
        if len(v) > 0 {
            c.Ctx.Output.Header(k, v[0])
        }
    }

    // Leer el cuerpo de la respuesta (el reporte)
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        c.Ctx.Output.SetStatus(500)
        c.Ctx.Output.Body([]byte("Error al leer el reporte"))
        return
    }

    // Enviar el contenido HTML del reporte
    c.Ctx.Output.SetStatus(resp.StatusCode)
    c.Ctx.Output.Header("Content-Type", "text/html")
    c.Ctx.Output.Body(body)

	
	/*label := c.GetString("label")
    if label == "" {
        c.Ctx.Output.SetStatus(400)
        c.Ctx.Output.Body([]byte("Falta el parámetro 'label'"))
        return
    }

    // Configuración de Knowage
    knowageHost := "https://inteligenciainstitucional.portaloas.udistrital.edu.co"
    contextPath := "/knowage"
    user := "jerodrigueza@udistrital.edu.co"         // ⚠️ Reemplaza con variables de entorno seguras
    password := "jerodrigueza"
    role := "/spagobi/user"

    // Crear cliente con soporte de cookies
    jar, _ := cookiejar.New(nil)
    client := &http.Client{Jar: jar}

    // Paso 1: Autenticación
    loginURL := fmt.Sprintf("%s%s/servlet/AdapterHTTP", knowageHost, contextPath)
    loginData := url.Values{}
    loginData.Set("NEW_SESSION", "TRUE")
    loginData.Set("USER", user)
    loginData.Set("PASSWORD", password)
    loginData.Set("ROLE", role)

    loginReq, _ := http.NewRequest("POST", loginURL, strings.NewReader(loginData.Encode()))
    loginReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    _, err := client.Do(loginReq)
    if err != nil {
        c.Ctx.Output.SetStatus(500)
        c.Ctx.Output.Body([]byte("Error al autenticar con Knowage"))
        return
    }

    // Paso 2: Construir URL del reporte
    reporteURL := fmt.Sprintf(
        "%s%s/servlet/AdapterHTTP?ACTION_NAME=EXECUTE_DOCUMENT_ANGULAR_ACTION&OBJECT_LABEL=%s&NEW_SESSION=TRUE&ROLE=%s&TOOLBAR_VISIBLE=true",
        knowageHost, contextPath, url.QueryEscape(label), url.QueryEscape(role),
    )

    // Paso 3: Obtener reporte
    req, _ := http.NewRequest("GET", reporteURL, nil)
    resp, err := client.Do(req)
    if err != nil {
        c.Ctx.Output.SetStatus(500)
        c.Ctx.Output.Body([]byte("Error al obtener el reporte desde Knowage"))
        return
    }
    defer resp.Body.Close()

    // Copiar encabezados necesarios
    for k, v := range resp.Header {
        if len(v) > 0 {
            c.Ctx.Output.Header(k, v[0])
        }
    }

    c.Ctx.Output.SetStatus(resp.StatusCode)
    body, _ := io.ReadAll(resp.Body)
    c.Ctx.Output.Body(body)*/
}
