## To run
Place ```.deploy``` file to your working folder or pass deploy-file name as argument. Deploy file should looks like below and filled with commands from Commands section:
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
## Commands
#### Copy
```
{
  "type": "copy",
  "file": ".service"
}
```
Copy file with name to working folder.