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
        "/auth/me": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/controllers.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/models.Restaurant"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/auth/sign-in": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "parameters": [
                    {
                        "description": "Sign in data",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.SignInData"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/controllers.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/controllers.SignInResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/restaurants": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Admin"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/controllers.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/models.Restaurant"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Restaurants"
                ],
                "parameters": [
                    {
                        "description": "Restaurant Data",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.CreateRestaurantInput"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/controllers.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/models.Restaurant"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            },
            "patch": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Restaurants"
                ],
                "parameters": [
                    {
                        "description": "Restaurant Data",
                        "name": "restaurant",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.UpdateRestaurantInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/controllers.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/models.Restaurant"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/restaurants/{id}": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Admin"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "Restaurant ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/controllers.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/models.Restaurant"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Admin"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "Restaurant ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/controllers.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "object"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controllers.Response": {
            "type": "object",
            "properties": {
                "data": {},
                "message": {
                    "type": "string"
                },
                "success": {
                    "type": "boolean"
                }
            }
        },
        "controllers.SignInData": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "example@gmail.com"
                },
                "password": {
                    "type": "string",
                    "minLength": 8,
                    "example": "12345678"
                }
            }
        },
        "controllers.SignInResponse": {
            "type": "object",
            "required": [
                "address",
                "email",
                "name",
                "phone"
            ],
            "properties": {
                "accessToken": {
                    "type": "string"
                },
                "address": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "phone": {
                    "type": "string"
                }
            }
        },
        "models.CreateRestaurantInput": {
            "type": "object",
            "required": [
                "address",
                "email",
                "name",
                "password",
                "phone"
            ],
            "properties": {
                "address": {
                    "type": "string",
                    "example": "36 Pasteur, Ben Nghe, Quan 1, Ho Chi Minh City"
                },
                "email": {
                    "type": "string",
                    "example": "example@gmail.com"
                },
                "name": {
                    "type": "string",
                    "example": "Example Restaurant"
                },
                "password": {
                    "type": "string",
                    "minLength": 8,
                    "example": "12345678"
                },
                "phone": {
                    "type": "string",
                    "example": "1234567890"
                }
            }
        },
        "models.Restaurant": {
            "type": "object",
            "required": [
                "address",
                "email",
                "name",
                "phone"
            ],
            "properties": {
                "address": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "phone": {
                    "type": "string"
                }
            }
        },
        "models.UpdateRestaurantInput": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string",
                    "minLength": 8
                },
                "phone": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:5000",
	BasePath:         "/api/v1",
	Schemes:          []string{"http", "https"},
	Title:            "Smart Serve",
	Description:      "Smart Serve API",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
