## To run
Place ```.deploy``` file to your working folder or pass deploy-file name as argument. Deploy file should looks like below and filled with commands from Commands section:
```
{
  "folder": "update",

  "Do": [{
    "type": "copy",
    "file": "go.mod"
  }, {
    "type": "run",
    "path": "echo",
    "arguments": ["echo"]
  }, {
    "type": "run",
    "path": "go",

    "Environments": [
      "CGO_ENABLED=0"
    ],
    "Arguments": [
      "build", "-o", "update/main"
    ]
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