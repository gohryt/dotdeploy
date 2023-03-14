## To run
Place ```.deploy``` file to your working folder or pass deploy-file name as argument. Deploy file should looks like below and filled with commands from Commands section:
```
{
  "folder": "update",
  "keep": true,

  "Do": [{
    "type": "copy",
    "parallel": true,

    "from": ".deploy",
    "to": "update/.deploy"
  }, {
    "type": "run",
    "parallel": true,

    "path": "echo",
    "timeout": 4,

    "Environment": ["HELLO='FROM DEPLOY'"],
    "Query": ["hello", "from", ".deploy"]
  }]
}
```
```folder``` is folder wich will be created on start, it will be deleted at processing end while setting ```keep``` is false.
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
Copy ```from``` ```to```. ```to``` key may be ignored. In this case programm will copy file to workrirectory ```folder```.
#### Move
```
{
  "type": "move",

  "from": ".service",
  "to": "update/.service"
}
```
Move ```from``` ```to```. ```to``` key may be ignored. In this case programm will copy file to workrirectory ```folder```.
#### Run
```
{
  "type": "run",

  "path": "echo",
  "timeout": 4,
  
  "Environment": ["HELLO='FROM DEPLOY'"],
  "Query": ["hello", "$HELLO", ".deploy"]
}
```
Run some ```path``` with or without timeout. You can also set Environment and Query.