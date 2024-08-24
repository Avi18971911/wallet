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
        "/accounts/login": {
            "post": {
                "description": "Logs in a user with the provided username and password.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "accounts"
                ],
                "summary": "Login",
                "parameters": [
                    {
                        "description": "Login payload",
                        "name": "login",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.AccountLoginDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successful login",
                        "schema": {
                            "$ref": "#/definitions/dto.AccountDetailsDTO"
                        }
                    },
                    "401": {
                        "description": "Invalid credentials",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorMessage"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorMessage"
                        }
                    }
                }
            }
        },
        "/accounts/{accountId}": {
            "get": {
                "description": "Retrieves the details of a specific account by its ID.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "accounts"
                ],
                "summary": "Get account details",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Account ID",
                        "name": "accountId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successful retrieval of account details",
                        "schema": {
                            "$ref": "#/definitions/dto.AccountDetailsDTO"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorMessage"
                        }
                    }
                }
            }
        },
        "/accounts/{accountId}/transactions": {
            "get": {
                "description": "Retrieves a list of transactions for a specific account by its ID.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "transactions"
                ],
                "summary": "Get account transactions",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Account ID",
                        "name": "accountId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successful retrieval of account transactions",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/dto.AccountTransactionDTO"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorMessage"
                        }
                    }
                }
            }
        },
        "/transactions": {
            "post": {
                "description": "Adds a new transaction to the system.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "transactions"
                ],
                "summary": "Add a new transaction",
                "parameters": [
                    {
                        "description": "Transaction request",
                        "name": "transaction",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.TransactionRequest"
                        }
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid request payload",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorMessage"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/utils.ErrorMessage"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dto.AccountDetailsDTO": {
            "type": "object",
            "required": [
                "accountNumber",
                "accountType",
                "availableBalance",
                "createdAt",
                "id",
                "knownAccounts",
                "person",
                "username"
            ],
            "properties": {
                "accountNumber": {
                    "description": "The account number",
                    "type": "string"
                },
                "accountType": {
                    "description": "The type of the account",
                    "type": "string"
                },
                "availableBalance": {
                    "description": "The available balance in the account",
                    "type": "number"
                },
                "createdAt": {
                    "description": "The creation timestamp of the account",
                    "type": "string"
                },
                "id": {
                    "description": "The unique identifier of the account",
                    "type": "string"
                },
                "knownAccounts": {
                    "description": "The list of accounts known to and recognized by the account holder",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dto.KnownAccountDTO"
                    }
                },
                "person": {
                    "description": "The account holder associated with the account",
                    "allOf": [
                        {
                            "$ref": "#/definitions/dto.PersonDTO"
                        }
                    ]
                },
                "username": {
                    "description": "The username associated with the account",
                    "type": "string"
                }
            }
        },
        "dto.AccountLoginDTO": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "password": {
                    "description": "The password for the login",
                    "type": "string"
                },
                "username": {
                    "description": "The username for the login",
                    "type": "string"
                }
            }
        },
        "dto.AccountTransactionDTO": {
            "type": "object",
            "required": [
                "accountId",
                "amount",
                "createdAt",
                "id",
                "otherAccountId",
                "transactionType"
            ],
            "properties": {
                "accountId": {
                    "description": "The primary account ID associated with the transaction",
                    "type": "string"
                },
                "amount": {
                    "description": "The amount involved in the transaction",
                    "type": "number"
                },
                "createdAt": {
                    "description": "The timestamp of when the transaction was created",
                    "type": "string"
                },
                "id": {
                    "description": "The unique identifier of the transaction",
                    "type": "string"
                },
                "otherAccountId": {
                    "description": "The other account ID involved in the transaction",
                    "type": "string"
                },
                "transactionType": {
                    "description": "The type of the transaction (debit or credit)",
                    "type": "string"
                }
            }
        },
        "dto.KnownAccountDTO": {
            "type": "object",
            "required": [
                "accountHolder",
                "accountNumber",
                "accountType"
            ],
            "properties": {
                "accountHolder": {
                    "description": "The name of the account holder",
                    "type": "string"
                },
                "accountNumber": {
                    "description": "The account number of the known account",
                    "type": "string"
                },
                "accountType": {
                    "description": "The type of the account (e.g., savings, checking)",
                    "type": "string"
                }
            }
        },
        "dto.PersonDTO": {
            "type": "object",
            "required": [
                "firstName",
                "lastName"
            ],
            "properties": {
                "firstName": {
                    "description": "The first name of the person",
                    "type": "string"
                },
                "lastName": {
                    "description": "The last name of the person",
                    "type": "string"
                }
            }
        },
        "dto.TransactionRequest": {
            "type": "object",
            "required": [
                "amount",
                "fromAccount",
                "toAccount"
            ],
            "properties": {
                "amount": {
                    "description": "The amount to be transferred",
                    "type": "number"
                },
                "fromAccount": {
                    "description": "The account number of the account from which the amount is to be transferred",
                    "type": "string"
                },
                "toAccount": {
                    "description": "The account number of the account to which the amount is to be transferred",
                    "type": "string"
                }
            }
        },
        "utils.ErrorMessage": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/backendAPI",
	Schemes:          []string{},
	Title:            "Wallet API",
	Description:      "This is a simple wallet API",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
