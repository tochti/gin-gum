GUM REST
========

Funktionen um mit Gin RESTfull APIs zu erstellen.

```golang
r := gin.New()

r.GET("/", func (c *gin.Context) {
  err := errors.New("Resource not found!")
  gumrest.ErrorResponse(404, err)
})
```

erzeugt folgenden JSON Response

```json
{
  "error": "Resource not found!"
}
```

http response type application/json, http status code 404
