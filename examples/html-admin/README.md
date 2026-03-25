# HTML admin example

Run from the root of the repository:

```bash
$ go run ./examples/html-admin
```

Open:

- http://localhost:9094/providers
- http://localhost:9094/app-settings

This example shows:

- `x/resource.MountCollection` for `/providers`
- `x/resource.MountSingleton` for `/app-settings`
- `x/htmx` response triggers to refresh the providers list
