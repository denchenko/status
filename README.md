# status

```
package main

import (
    "net/http"
    "github.com/gostatus/status"
)

func main() {
    healthChecker := status.NewHealthChecker().
    WithTarget("db", status.TargetImportanceHigh, func(ctx context.Context) error{
    return dbClient.Ping(ctx)
})

    statusPage := status.NewPage().
    WithHealthChecker(healthChecker).
    WithURL("Swagger", "/swagger")

    http.HandleFunc("/health", healthChecker.Handler())
    http.HandleFunc("/status", statusPage.Handler())
    http.ListenAndServe(":8080", nil)
}
```
