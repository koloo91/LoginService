// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2019-08-20 22:22:08.950107 +0200 CEST m=+0.026235054

package docs

import (
	"bytes"
	"encoding/json"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/login": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Login a user",
                "operationId": "login",
                "parameters": [
                    {
                        "description": "login json",
                        "name": "loginVo",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/model.LoginVo"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/model.LoginResultVo"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/model.HttpError"
                        }
                    }
                }
            }
        },
        "/api/register": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Registers a new user",
                "operationId": "register",
                "parameters": [
                    {
                        "description": "register json",
                        "name": "registerVo",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/model.RegisterVo"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/model.UserVo"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/model.HttpError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.HttpError": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "model.LoginResultVo": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "model.LoginVo": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "model.RegisterVo": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "model.UserVo": {
            "type": "object",
            "properties": {
                "created": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "updated": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{Schemes: []string{}}

type s struct{}

func (s *s) ReadDoc() string {
	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, SwaggerInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}