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

### Example Config

```json
{
    "port": 8000,
    "rules": [
        {
            "_comments": "Update the Satis packages whenever code is pushed to a robstrong repo",
            "command": "php /opt/satis/bin/satis build /opt/satis/satis.json /srv/www -n",
            "criteria": [
                {
                    "event": "push",
                    "owner": "robstrong"
                }
            ]
        },
        {
            "_comments": "Update the issues page when an issue is updated/created",
            "command": "/opt/github-history/go-gh-history -o /srv/www/issues.html issues robstrong/hook-receiver",
            "criteria": [
                {
                    "event": "issues",
                    "owner": "robstrong",
                    "repository": "hook-receiver"
                }
            ]
        },
        {
            "_comments": "Pull down branches starting with 'release-' locally when pushed",
            "command": "cd ~/go-github-history/ && git checkout {{ .Branch }} && git pull origin {{ .Branch }}",
            "criteria": [
                {
                    "event": "push",
                    "owner": "robstrong",
                    "repository": "go-github-history",
                    "push_params": {
                        "branch": "release-*"
                    }
                }
            ]
        }
    ]
}
```

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
