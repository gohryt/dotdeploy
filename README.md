## To run
Place ```.deploy``` file To your working folder or pass deploy-file name as argument. Deploy file should looks like below and filled with commands from Commands section:
```
{
  "folder": "update",
  "keep": false,

  "Remote": [{
    "type": "key",
    "name": "key",

    "host": "1.1.1.1",
    "file": "/home/example/.ssh/id_ed25519",
    "username": "root",
    "password": "example"
  }, {
    "type": "password",
    "name": "password",

    "host": "1.1.1.1",
    "username": "root",
    "password": "example"
  }, {
    "type": "agent",
    "name": "agent",

    "host": "1.1.1.1",
    "username": "root"
  }],

  "Do": [{
    "type": "copy",

    "From": {
      "path": ".deploy"
    },
    "To": {
      "path": "update/.deploy"
    }
  }, {
    "type": "move",
    "parallel": true,

    "From": {
      "path": "update/.deploy"
    },
    "To": {
      "path": "update/.deploy.example"
    }
  }, {
    "type": "run",
    "parallel": true,

    "timeout": 4,

    "Path": {
      "path": "echo"
    },
    "Environment": ["HELLO='FROM DEPLOY'"],
    "Query": ["hello", "From", ".deploy"]
  }]
}
```
```folder``` is folder wich will be created on start, it will be deleted at processing end while setting ```keep``` is false.
## Remote
Each element should have ```name``` filled because of referencing (as above).
#### Key
```
{
  "type": "key",

  "host": "1.1.1.1",
  "file": "/home/example/.ssh/id_ed25519",
  "password": "example",
  "username": "root"
}
```
#### Password
```
{
  "type": "password",

  "host": "1.1.1.1",
  "username": "root",
  "password": "example"
}
```
#### Agent
```
{
  "type": "agent",

  "host": "1.1.1.1",
  "username": "root"
}
```
## Do
Each element should have ```name``` filled because of referencing (as above).
#### Parallel
```
{
  "parallel": true
}
```
Each element may be ```parallel``` which means that it will be started in goroutine.
#### Copy
```
{
  "type": "copy",

  "From": {
    "path": ".service"
  },
  "To": {
    "path": "update/.service"
  }
}
```
Copy ```From``` ```To```. ```To``` key may be ignored. In this case programm will copy file to workrirecTory ```folder```.
#### Move
```
{
  "type": "move",

  "From": {
    "path": ".service"
  },
  "To": {
    "path": "update/.service"
  }
}
```
Move ```From``` ```To```. ```To``` key may be ignored. In this case programm will copy file to workrirecTory ```folder```.
#### Run
```
{
  "type": "run",

  "timeout": 4,
  
  "Path": {
    "path": "echo"
  },
  "Environment": ["HELLO='FROM DEPLOY'"],
  "Query": ["hello", "from", ".deploy"]
}
```
Run some ```Path``` with or without timeout. You can also set Environment and Query.