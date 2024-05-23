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
        "/admin/sign-in": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "SingIn admin",
                "parameters": [
                    {
                        "description": "sign in",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/admin.SignInRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/admin.SignInResponse"
                        }
                    }
                }
            }
        },
        "/admin/sign-up": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "SingUp admin",
                "parameters": [
                    {
                        "description": "sign up",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/admin.SignUpRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    }
                }
            }
        },
        "/cart/add": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Добавляет товар в корзину, если товар уже в корзине, увеличивает количество",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cart"
                ],
                "summary": "Create cart item",
                "parameters": [
                    {
                        "description": "product id and count",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/cart.AddRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    }
                }
            }
        },
        "/cart/change-count": {
            "patch": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Изменяет кол-во товара в корзине на new_count",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cart"
                ],
                "summary": "Change count curt item",
                "parameters": [
                    {
                        "description": "product id and count",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/cart.ChangeCountRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    }
                }
            }
        },
        "/cart/get": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cart"
                ],
                "summary": "Get all cart items",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/cart.GetResponse"
                        }
                    }
                }
            }
        },
        "/customer/sign-in": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "SingIn customer",
                "parameters": [
                    {
                        "description": "sign in",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/customer.SignInRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/customer.SignInResponse"
                        }
                    }
                }
            }
        },
        "/customer/sign-up": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "SingUp customer",
                "parameters": [
                    {
                        "description": "sign up",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/customer.SignUpRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    }
                }
            }
        },
        "/image/download/{id}": {
            "get": {
                "description": "При удачном запросе вернет картинку в body со статусом 200, при неудачном json с ошибкой",
                "produces": [
                    "image/jpeg",
                    "image/png",
                    "application/json"
                ],
                "tags": [
                    "files"
                ],
                "summary": "DownloadImage",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "image id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "integer"
                            }
                        }
                    }
                }
            }
        },
        "/image/upload": {
            "post": {
                "description": "В боди должна быть картинка в виде массива байт",
                "consumes": [
                    "image/jpeg",
                    "image/png"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "files"
                ],
                "summary": "UploadImage",
                "parameters": [
                    {
                        "description": "byte image",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "integer"
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/image.UploadResponse"
                        }
                    }
                }
            }
        },
        "/order/get": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "order"
                ],
                "summary": "Get all orders",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/order.GetResponse"
                        }
                    }
                }
            }
        },
        "/order/make": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Создает заказ исходя из корзины покупателя",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "order"
                ],
                "summary": "Make order",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    }
                }
            }
        },
        "/product/create": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "product"
                ],
                "summary": "Create product",
                "parameters": [
                    {
                        "description": "alias new website",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/product.CreateRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    }
                }
            }
        },
        "/product/get-by-alias/{alias}": {
            "get": {
                "description": "Возвращает все товары сайта по алиасу",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "product"
                ],
                "summary": "GetByAlias",
                "parameters": [
                    {
                        "type": "string",
                        "description": "website alias",
                        "name": "alias",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/product.GetResponse"
                        }
                    }
                }
            }
        },
        "/website/aliases": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "website"
                ],
                "summary": "Get all users aliases",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/website.AliasesResponse"
                        }
                    }
                }
            }
        },
        "/website/create": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "website"
                ],
                "summary": "Create website",
                "parameters": [
                    {
                        "description": "alias new website",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/website.CreateRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    }
                }
            }
        },
        "/website/delete/{alias}": {
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Удаляет сайт по алиасу",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "website"
                ],
                "summary": "Delete website",
                "parameters": [
                    {
                        "type": "string",
                        "description": "website alias",
                        "name": "alias",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "admin.SignInRequest": {
            "type": "object",
            "properties": {
                "login": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "admin.SignInResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "token": {
                    "type": "string"
                }
            }
        },
        "admin.SignUpRequest": {
            "type": "object",
            "properties": {
                "login": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "cart.AddRequest": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "integer"
                },
                "product_id": {
                    "type": "integer"
                }
            }
        },
        "cart.ChangeCountRequest": {
            "type": "object",
            "properties": {
                "new_count": {
                    "type": "integer"
                },
                "product_id": {
                    "type": "integer"
                }
            }
        },
        "cart.GetResponse": {
            "type": "object",
            "properties": {
                "cart_items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/storage.CartItem"
                    }
                },
                "error": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "customer.SignInRequest": {
            "type": "object",
            "properties": {
                "alias": {
                    "type": "string"
                },
                "login": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "customer.SignInResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "token": {
                    "type": "string"
                }
            }
        },
        "customer.SignUpRequest": {
            "type": "object",
            "properties": {
                "alias": {
                    "type": "string"
                },
                "login": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "image.UploadResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "order.GetResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "orders": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/storage.Order"
                    }
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "product.CreateRequest": {
            "type": "object",
            "properties": {
                "alias": {
                    "type": "string"
                },
                "product_info": {
                    "$ref": "#/definitions/product.Info"
                }
            }
        },
        "product.GetResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "products": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/storage.ProductInfo"
                    }
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "product.Info": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "image_id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "price": {
                    "type": "integer"
                }
            }
        },
        "response.Response": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "storage.CartItem": {
            "type": "object",
            "properties": {
                "cart_id": {
                    "type": "integer"
                },
                "count": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "product": {
                    "$ref": "#/definitions/storage.ProductInfo"
                }
            }
        },
        "storage.Order": {
            "type": "object",
            "properties": {
                "order_items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/storage.OrderItem"
                    }
                }
            }
        },
        "storage.OrderItem": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "order_id": {
                    "type": "integer"
                },
                "product": {
                    "$ref": "#/definitions/storage.ProductInfo"
                }
            }
        },
        "storage.ProductInfo": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "image_id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "price": {
                    "type": "integer"
                },
                "website_id": {
                    "type": "integer"
                }
            }
        },
        "website.AliasesResponse": {
            "type": "object",
            "properties": {
                "aliases": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "error": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "website.CreateRequest": {
            "type": "object",
            "properties": {
                "alias": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8082",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "ForeignKey",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
