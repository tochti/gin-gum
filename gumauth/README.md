Authentication Framework
========================

Funktionen für 
  * Benutzer registration
  * Benutzer authentication
  * SQL-Datenbank Funktionen für die Benutzer verwaltung
    * CRUD
    * Funktion zum erzeugen der Datenbank Tabellen
  * SQL-Implementation des UserStore Interfaces, benötigt für SignIn und SignedIn Handler
  * SQL-Implementationd des User Interfaces, benötigt für SignIn und SignedIn Handler

Ein Möglichkeit einen Benutzer zu authentzifiern ist mittels CSRF. Das heißt man meldet sich beim Server an und bekommt bei erfolgreicher anmeldung eine Cookie zurück welcher einen Token enthält. Dieser Token muss nun bei jedem Request mitgesendet werden damit der Server einen Authetififzierung vornehmen kann.

:name und :password müssen ein base64 String sein.

```go

func handler(c *gin.Context) {
    c.String(200, "Secret!")
}

func main () {
    // user ist ein Objekt welches das UserSignIn interface implementiert
    // sessionStore ist ein Object welches das SessionStore interface implementiert
    signIn := kauth.SignIn(user, sessionStore)
    signedIn := kauth.SignedIn(sessionStore)

    srv := gin.New()
    // :user ist der String mittels dem die FindUser Funktionen den Benutzer Indetifizieren kann, Mail oder Benutzername
    srv.GET("/SignIn/:name/:password", signIn)
    srv.GET("/", signedIn(handler))
    srv.Run()
}

```
