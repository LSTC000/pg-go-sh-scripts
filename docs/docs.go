// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/bash": {
            "post": {
                "description": "Create bash script",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Bash"
                ],
                "summary": "Create",
                "parameters": [
                    {
                        "type": "file",
                        "description": "Bash script file",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Bash"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/schema.HTTPError"
                        }
                    }
                }
            }
        },
        "/bash/execute/list": {
            "post": {
                "description": "Execute list of bash scripts",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Bash"
                ],
                "summary": "Execute List",
                "parameters": [
                    {
                        "type": "boolean",
                        "description": "Execute type: if true, then in a multithreading, otherwise in a single thread",
                        "name": "isSync",
                        "in": "query",
                        "required": true
                    },
                    {
                        "description": "List of execute bash script models",
                        "name": "execute",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/dto.ExecBashDTO"
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/schema.Message"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/schema.HTTPError"
                        }
                    }
                }
            }
        },
        "/bash/list": {
            "get": {
                "description": "Get list of bash scripts",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Bash"
                ],
                "summary": "Get list",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Limit param of pagination",
                        "name": "limit",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Offset param of pagination",
                        "name": "offset",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/schema.SwagBashPaginationLimitOffsetPage"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/schema.HTTPError"
                        }
                    }
                }
            }
        },
        "/bash/log/{bashId}/list": {
            "get": {
                "description": "Get list of bash logs by bash id",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Bash Log"
                ],
                "summary": "Get list by bash id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID of bash script",
                        "name": "bashId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Limit param of pagination",
                        "name": "limit",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Offset param of pagination",
                        "name": "offset",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/schema.SwagBashLogPaginationLimitOffsetPage"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/schema.HTTPError"
                        }
                    }
                }
            }
        },
        "/bash/{id}": {
            "get": {
                "description": "Get bash script by id",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Bash"
                ],
                "summary": "Get by id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID of bash script",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Bash"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/schema.HTTPError"
                        }
                    }
                }
            },
            "delete": {
                "description": "Remove bash script by id",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Bash"
                ],
                "summary": "Remove by id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID of bash script",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Bash"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/schema.HTTPError"
                        }
                    }
                }
            }
        },
        "/bash/{id}/file": {
            "get": {
                "description": "Get bash script file by id",
                "produces": [
                    "application/x-www-form-urlencoded"
                ],
                "tags": [
                    "Bash"
                ],
                "summary": "Get file by id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID of bash script",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "file"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/schema.HTTPError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dto.ExecBashDTO": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string",
                    "example": "59628b82-356c-4745-bc81-187015cde387"
                },
                "timeoutSeconds": {
                    "type": "integer"
                }
            }
        },
        "model.Bash": {
            "type": "object",
            "properties": {
                "body": {
                    "type": "string"
                },
                "createdAt": {
                    "type": "string",
                    "example": "2024-04-14T15:50:21.907561+07:00"
                },
                "id": {
                    "type": "string",
                    "example": "59628b82-356c-4745-bc81-187015cde387"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "model.BashLog": {
            "type": "object",
            "properties": {
                "bashId": {
                    "type": "string",
                    "example": "59628b82-356c-4745-bc81-187015cde387"
                },
                "body": {
                    "type": "string"
                },
                "createdAt": {
                    "type": "string",
                    "example": "2024-04-14T15:50:21.907561+07:00"
                },
                "id": {
                    "type": "string",
                    "example": "f4f4d096-ef4a-4649-8346-a952e2ca27d3"
                },
                "isError": {
                    "type": "boolean"
                }
            }
        },
        "schema.HTTPError": {
            "type": "object",
            "properties": {
                "detail": {
                    "type": "string"
                },
                "httpCode": {
                    "type": "integer"
                },
                "serviceCode": {
                    "type": "integer"
                }
            }
        },
        "schema.Message": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "schema.SwagBashLogPaginationLimitOffsetPage": {
            "type": "object",
            "properties": {
                "items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.BashLog"
                    }
                },
                "limit": {
                    "type": "integer"
                },
                "offset": {
                    "type": "integer"
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "schema.SwagBashPaginationLimitOffsetPage": {
            "type": "object",
            "properties": {
                "items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.Bash"
                    }
                },
                "limit": {
                    "type": "integer"
                },
                "offset": {
                    "type": "integer"
                },
                "total": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0.0",
	Host:             "0.0.0.0:8000",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Bash Scripts",
	Description:      "This is an API for running bash scripts",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
