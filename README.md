# modware-identity
<!-- ALL-CONTRIBUTORS-BADGE:START - Do not remove or modify this section -->
[![All Contributors](https://img.shields.io/badge/all_contributors-2-orange.svg?style=flat-square)](#contributors-)
<!-- ALL-CONTRIBUTORS-BADGE:END -->

[![License](https://img.shields.io/badge/License-BSD%202--Clause-blue.svg)](LICENSE)  
![Continuous integration](https://github.com/dictyBase/modware-identity/workflows/Continuous%20integration/badge.svg)
[![codecov](https://codecov.io/gh/dictyBase/modware-identity/branch/develop/graph/badge.svg)](https://codecov.io/gh/dictyBase/modware-identity)
[![Maintainability](https://api.codeclimate.com/v1/badges/21ed283a6186cfa3d003/maintainability)](https://codeclimate.com/github/dictyBase/modware-identity/maintainability)  
![Last commit](https://badgen.net/github/last-commit/dictyBase/modware-identity/develop)   
[![Funding](https://badgen.net/badge/Funding/Rex%20L%20Chisholm,dictyBase,DCR/yellow?list=|)](https://projectreporter.nih.gov/project_info_description.cfm?aid=10024726&icde=0)

[dictyBase](http://dictybase.org) **API** server to manage identity authorization
from third party providers (ORCID, Google, LinkedIn). The API server supports both gRPC
and HTTP/JSON protocol for data exchange.

## API

#### HTTP/JSON

It's [here](https://dictybase.github.io/dictybase-api), make sure you use the **identity** content from the dropdown on the top right.

#### gRPC

The protocol buffer definitions and service apis are documented
[here](https://github.com/dictyBase/dictybaseapis/tree/master/dictybase/identity).

## Usage

```
NAME:
   modware-identity - cli for modware-identity microservice

USAGE:
   main [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
     start-server          starts the modware-identity microservice with HTTP and grpc backends
     start-identity-reply  start the reply messaging(nats) backend for identity microservice
     create-identity       creates a new identity for an user
     help, h               Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --log-format value  format of the logging out, either of json or text. (default: "json")
   --log-level value   log level for the application (default: "error")
   --help, -h          show help
   --version, -v       print the version
```

## Subcommands

```
NAME:
   main start-server - starts the modware-identity microservice with HTTP and grpc backends

USAGE:
   main start-server [command options] [arguments...]

OPTIONS:
   --arangodb-pass value, --pass value    arangodb database password [$ARANGODB_PASS]
   --arangodb-database value, --db value  arangodb database name [$ARANGODB_DATABASE]
   --arangodb-user value, --user value    arangodb database user [$ARANGODB_USER]
   --arangodb-host value, --host value    arangodb database host (default: "arangodb") [$ARANGODB_SERVICE_HOST]
   --arangodb-port value                  arangodb database port (default: "8529") [$ARANGODB_SERVICE_PORT]
   --identity-api-http-host value         public hostname serving the http api, by default the default port will be appended to http://localhost [$IDENTITY_API_HTTP_HOST]
   --is-secure                            flag for secured or unsecured arangodb endpoint
   --nats-host value                      nats messaging server host [$NATS_SERVICE_HOST]
   --nats-port value                      nats messaging server port [$NATS_SERVICE_PORT]
   --port value                           tcp port at which the server will be available (default: "9560")
```

```
NAME:
   main start-identity-reply - start the reply messaging(nats) backend for identity microservice

USAGE:
   main start-identity-reply [command options] [arguments...]

OPTIONS:
   --identity-grpc-host value  grpc host address for identity service [$IDENTITY_API_SERVICE_HOST]
   --identity-grpc-port value  grpc port for identity service [$IDENTITY_API_SERVICE_PORT]
   --messaging-host value      host address for messaging server [$NATS_SERVICE_HOST]
   --messaging-port value      port for messaging server [$NATS_SERVICE_PORT]
```

```
NAME:
   main create-identity - creates a new identity for an user

USAGE:
   main create-identity [command options] [arguments...]

OPTIONS:
   --identity-grpc-host value  grpc host address for identity service [$IDENTITY_API_SERVICE_HOST]
   --identity-grpc-port value  grpc port for identity service [$IDENTITY_API_SERVICE_PORT]
   --user-grpc-host value      host address for user service [$USER_API_SERVICE_HOST]
   --user-grpc-port value      port for user service [$USER_API_SERVICE_PORT]
   --identifier value          Third party unique identifier
   --provider value            Name of provider who provides the identifier
   --email value               An user email(not the identity provider's email) that will be tied to the identifier
```
# Misc badges
![Issues](https://badgen.net/github/issues/dictyBase/modware-identity)
![Open Issues](https://badgen.net/github/open-issues/dictyBase/modware-identity)
![Closed Issues](https://badgen.net/github/closed-issues/dictyBase/modware-identity)  
![Total PRS](https://badgen.net/github/prs/dictyBase/modware-identity)
![Open PRS](https://badgen.net/github/open-prs/dictyBase/modware-identity)
![Closed PRS](https://badgen.net/github/closed-prs/dictyBase/modware-identity)
![Merged PRS](https://badgen.net/github/merged-prs/dictyBase/modware-identity)  
![Commits](https://badgen.net/github/commits/dictyBase/modware-identity/develop)
![Branches](https://badgen.net/github/branches/dictyBase/modware-identity)
![Tags](https://badgen.net/github/tags/dictyBase/modware-identity/?color=cyan)  
![GitHub repo size](https://img.shields.io/github/repo-size/dictyBase/modware-identity?style=plastic)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/dictyBase/modware-identity?style=plastic)
[![Lines of Code](https://badgen.net/codeclimate/loc/dictyBase/modware-identity)](https://codeclimate.com/github/dictyBase/modware-identity/code)  

## Contributors ✨

Thanks goes to these wonderful people ([emoji key](https://allcontributors.org/docs/en/emoji-key)):

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tr>
    <td align="center"><a href="http://www.erichartline.net/"><img src="https://avatars3.githubusercontent.com/u/13489381?v=4" width="100px;" alt=""/><br /><sub><b>Eric Hartline</b></sub></a><br /><a href="https://github.com/dictyBase/modware-identity/issues?q=author%3Awildlifehexagon" title="Bug reports">🐛</a> <a href="https://github.com/dictyBase/modware-identity/commits?author=wildlifehexagon" title="Code">💻</a> <a href="#content-wildlifehexagon" title="Content">🖋</a> <a href="https://github.com/dictyBase/modware-identity/commits?author=wildlifehexagon" title="Documentation">📖</a> <a href="#maintenance-wildlifehexagon" title="Maintenance">🚧</a></td>
    <td align="center"><a href="http://cybersiddhu.github.com/"><img src="https://avatars3.githubusercontent.com/u/48740?v=4" width="100px;" alt=""/><br /><sub><b>Siddhartha Basu</b></sub></a><br /><a href="https://github.com/dictyBase/modware-identity/issues?q=author%3Acybersiddhu" title="Bug reports">🐛</a> <a href="https://github.com/dictyBase/modware-identity/commits?author=cybersiddhu" title="Code">💻</a> <a href="#content-cybersiddhu" title="Content">🖋</a> <a href="https://github.com/dictyBase/modware-identity/commits?author=cybersiddhu" title="Documentation">📖</a> <a href="#maintenance-cybersiddhu" title="Maintenance">🚧</a></td>
  </tr>
</table>

<!-- markdownlint-enable -->
<!-- prettier-ignore-end -->
<!-- ALL-CONTRIBUTORS-LIST:END -->

This project follows the [all-contributors](https://github.com/all-contributors/all-contributors) specification. Contributions of any kind welcome!