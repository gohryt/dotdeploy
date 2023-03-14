## To run
Place ```.deploy``` file to your working folder or pass deploy-file name as argument. Deploy file should looks like below and filled with commands from Commands section:
```
{
  "folder": "update",

  "Do": [{
    "type": "copy",
    "parallel": true,
    "from": "go.mod",
    "to": "update/go.mod"
  }, {
    "type": "copy",
    "parallel": true,
    "from": "go.mod"
  }, {
    "type": "run",
    "parallel": true,
    "path": "echo",
    
    "ArgumentList": [
      "echo"
    ]
  }, {
    "type": "run",
    "parallel": true,
    "path": "go",

    "Environment": [
      "CGO_ENABLED=0"
    ],
    "ArgumentList": [
      "build", "-o", "update/main"
    ]
  }]
}
```
## Commands
#### Parallel
```
{
  "parallel": true
}
```
Each command may be ```parallel``` which means that it will be started in goroutine.
#### Copy
```
{
  "type": "copy",
  "from": ".service",
  "to": "update/.service"
}
```
Copy file with name to working folder. ```to``` key may be ignored. In this case programm will copy file to workrirectory ```folder```.