```
{
  "folder": "update",

  "do": [{
    "type": "copy",
    "file": ".service"
  }, {
    "type": "run",
    "command": "CGO_ENABLED=0 go build -o main"
  }]
}
```