# Toolbox

A simple example of how to create a reusable Go module with commonly used tools.

The included tools are:

- [ ] Read JSON
- [ ] Write JSON
- [ ] Produce a JSON encoded error response
- [X] Upload a file to a specified directory
- [X] Download a static file
- [X] Get a random string of length n
- [ ] Post JSON to a remote service 
- [X] Create a directory, including all parent directories, if it does not already exist
- [X] Create a URL safe slug from a string

## Installation

`go get -u github.com/renatofagalde/golang-toolkit`


# Middleware de Contexto: Configuração e Recuperação

Este projeto implementa um middleware que valida os cabeçalhos `X-Request-ID` e `journey`, propagando esses valores para o contexto da requisição. Este documento explica como configurar o middleware nas rotas e recuperar os valores no contexto.

---

## Configuração do Middleware

### 1. Configurando o Middleware nas Rotas

O middleware deve ser aplicado ao grupo de rotas desejado. Por exemplo, para as rotas que possuem o prefixo `/api`:

```go
apiGroup := r.Group("/api")
apiGroup.Use(h.RequestMiddleware())
```

> **Atenção**: Verifique quais rotas precisarão do middleware antes de aplicá-lo, para evitar sobrecarga desnecessária.

---

## Recuperando Valores do Contexto

Os valores dos cabeçalhos `X-Request-ID` e `journey` podem ser recuperados em qualquer parte do código com o método `h.Give()`:

```go
journey, requestID := h.Give()
```

- **`journey`**: Identificador da jornada associado à requisição.
- **`requestID`**: Identificador único da requisição.

Caso os valores não tenham sido configurados (por ausência no cabeçalho ou falha no middleware), os retornos serão strings vazias (`""`).

---

### Exemplo de Uso no Serviço

```go
package service

import (
	"bootstrap/src/helpers"
	"log"
)

func RepositoryFindXXX() {
	var logger toolkit.Logger
	journey, requestID := h.Give()
	logger.Info(fmt.Sprintf("FindUserByEmail: %s obj %+v", email,user),
		zap.String("stage", "repository"),
		zap.String("value", query),
		zap.String("journey", journey),
		zap.String("requestID", requestID))

}
```

---

## Boas Práticas

1. **Valide os valores recuperados**:
    - Sempre verifique se os valores retornados não estão vazios antes de utilizá-los.

2. **Evite Sobrecarga do Middleware**:
    - Aplique o middleware apenas nas rotas que realmente necessitam do contexto propagado.

3. **Rastreabilidade e Auditoria**:
    - Utilize os valores de `journey` e `requestID` para rastrear requisições e associá-las a fluxos específicos.

---

Com essa configuração, você garante que o contexto seja gerenciado de forma centralizada e acessível, promovendo organização e rastreabilidade no projeto. 🚀
