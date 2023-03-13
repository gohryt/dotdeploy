#### To run
place ```.deploy``` file to your working folder or pass deploy-file name as argument:
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