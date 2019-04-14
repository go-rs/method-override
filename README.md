# Middleware: Request Method Override

## How to use?
<b>Header:</b> X-HTTP-Method-Override
<br>
<b>Methods:</b> GET, POST, PUT, DELETE, OPTIONS, HEAD, PATCH

```
var api rest.API

api.Use(methodoverride.MethodOverride())

```
