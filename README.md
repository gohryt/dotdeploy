## To run
Place ```.deploy``` file to your working folder or pass deploy-file name as argument. Deploy file should looks like below and filled with commands from Commands section:
```
{
  "folder": "update",
  "remove": true,

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

    "Environment": [
      "HELLO='FROM DEPLOY'"
    ],
    "Query": [
      "hello", "from", ".deploy"
    ]
  }]
}
```
```folder``` is folder wich will be created on start, it also may be deleted at processing end with setting ```remove``` to true.
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
#### Run
```
{
  "type": "run",

  "path": "echo",
  "timeout": 4,
  
  "Environment": [
    "HELLO='FROM DEPLOY'"
  ],
  "Query": [
    "hello", "$HELLO", ".deploy"
  ]
}
```
Run some ```path``` with or without timeout. You can also set Environment and Query.