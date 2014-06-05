## Github WebHook Receiver and Command 
This project is a simple web server which will receive Github WebHooks and then run a command based on the payload. 

### Config
When the server starts, it reads from `config.json` which should be located in the same directory as the binary. The config file has the following properties.

| Property | Required | Description |
|----------|----------|-------------| 
| `port` | yes | The port that the server runs on |
| `rules[][criteria][][type]` | yes | Type of Github Event |
| `rules[][criteria][][owner]` | no | Repository Owner |
| `rules[][criteria][][repository]` | no | Name of repository |
| `rules[][command]` | yes | Command to run if the criteria matches |
| `cidr_override` | no | All requests received by the server are checked to make sure the requesting IP matches the given CIDR. If `cidr_override` is not specified, Github's CIDR is used |
