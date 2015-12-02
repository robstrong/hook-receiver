## Github WebHook Receiver and Command Runner
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

### Templating

The command that is run will be evaluated as a template. The payload data will be passed to the template. So you can use
data from the Github hook in your command. For example, if you received a push event, you could setup a command like this:

```
"fabric deploy_branch:{{.Branch}}
```

You can see what data is available in the *action*_event.go (push_event.go, release_event.go, etc).

#### Functions

Currently there is only one template function:

__after__ - will return the string that comes after the string that is passed in

If Branch = "release-myNewFeature"
```
{{ .Branch | after "release-" }}
```
Would output "myNewFeature"
