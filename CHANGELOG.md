# Release Notes

## Version v0.0.19 - 2020-10-00

Improvements:

- Updated function labels

## Version v0.0.18 - 2020-09-16

Bug Fixes:

- Fix bug in stack declaration swagger.file
- Fix bug in use sha commit
- Fixed url of functions in `safira status`

## Version v0.0.17 - 2020-09-10

Bug Fixes:

- Fix ingress backend
- Fix declaration build.useSha, deploy.prefix and deploy.suffix

## Version v0.0.16 - 2020-09-09

Features:

- Added functionality to use sha commit in image build
- Added swagger file name declaration option
- Added the possibility to use prefix and suffix in the name of the deploy

Improvements:

- The update flag in the deploy process now impacts the swagger ui

## Version v0.0.15 - 2020-09-03

Features:

- Added possibility to disable the build of a function

Improvements:

- Rewritten `safira upgrade` command
- Openfaas url updated

## Version v0.0.14 - 2020-08-31

Features:

- Ignore node_modules in build

Bug Fixes:

- Fix output errors

## Version v0.0.13 - 2020-08-28

Bug Fixes:

- Installation Vertigo iPaaS in local cluster fixed
- A bug did not allow the use of the path attribute in environment files

## Version v0.0.12 - 2020-08-27

Features:

- Added attribute path in stack.yml

Improvements:

- Added update-template flag in `function build`
- All outputs placed in English
- Standardized displays of all outputs

Bug Fixes:

- A fixed bug that did not remove the swagger configmap when executing the `safira function undeploy` command

## Version v0.0.11 - 2020-08-21

Improvements:

- Swagger spec move to repository scope

Bug Fixes:

- Fixed functions and plugins undeploy

## Version v0.0.10 - 2020-08-21

Features:

- Added a new command in CLI to itself upgrade: `safira upgrade`

Improvements:

- Flag added to remove the function folder from the `safira function remove` command
- Added new annotations in function deploy
- Added ingress class in function deploy

Bug Fixes:

- Fixed error when hostname port was not declared in stack.yaml
- Fix namespace Kong Plugin deploy

## Version v0.0.9 - 2020-07-17

Features:

- Added kong plugin declaration on stack.yml
- Enabled creation of files from different environments

Improvements:

- Removed the kong command

## Version v0.0.8 - 2020-07-10

Features:

- Added configuration of environment variables in stack.yml
- Added new command `safira function remove`

Improvements:

- Improved command description
- Added execute examples in help flag
- Added `bp` alias in the `build-push` subcommand
- Renamed the `function remove` command to` function undeploy`
- Added flag namespace in the `function undeploy` subcommand

## Version v0.0.7 - 2020-07-03

Features:

- Added new subcommand `safira function log`

Bug Fixes:

- Fixed URL of swagger ui in `safira infra status`

## Version v0.0.6 - 2020-07-02

Features:

- Added new flag `kubeconfig` in `safira function remove`
- Added new command `safira kong`
- Added new subcommand `safira kong new`

Improvements:

- The link to the swagger editor has been added to the output of `safira infra up`
- The removed suffix from deploys at output of `safira infra status`
- Added output of URLs of the `safira infra status` command

## Version v0.0.5 - 2020-06-30

Features:

- Added declaration of CPU and memory usage limit in stack.yaml
- Added okteto binary installation
- Added new command `safira okteto login`
- Added new command `safira template pull`
- Added Dockerfile
- Added namespace flag in command `safira function deploy`
- Added .gitlab-ci.yml file creation at execution `safira function new`

Bug Fixes:

- Fix `safira remove function` for remove service and ingress

## Version v0.0.4 - 2020-06-22

Features:

- Added swagger-ui deploy
- Added swagger-editor deploy

Improvements:

- Improvements function deploy

Bug Fixes:

- Fix run `safira init` in user root

## Version v0.0.3 - 2020-06-10

Features:

- Added new flag `hostname` in `safira function deploy`
- Added config scale in a stack.yaml
- Added possibility to declare custom yamls in stack.yaml

Improvements:

- Added info UP-TO-DATE/AVAILABLE in `safira infra status`

## Version v0.0.2 - 2020-06-08

Features:

- Added new flag `kubeconfig` in `safira function deploy`
- Added the new command `sapira infra status`

Improvements:

- `safira function deploy` Add flag `update`, force the deploy to take a new image
- `safira infra secrets` Add Konga secrets
- `safira function new` Yaml project renamed to stack.yaml
- Changed the search for .env information from `safira function deploy` to yaml project
- Improved various error messages
- Removed kongplugin creation on deploy
- Added new `all-functions` flag in the deploy
- Separate the build and push command, maintaining an alternative to the two actions with the name of build-push
- Atualizado o namespace do vertigo ipaas

## Version v0.0.1-beta.2 - 2020-06-01

Features:

- Add flag verbose in all commands
- Add new command `safira init`

Improvements:

- Improvements in all errors messages

## Version v0.0.1-beta - 2020-05-26

Features:
    
- Add new command `safira infra secrets` to get user and password of applications
- `safira function new` Add a file .env

Improvements:
    
- `function new` Add a folder deploy in .gitignore

## Version v0.0.1-alpha - 2020-05-25

Initial Safira release
