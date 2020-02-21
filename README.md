## Description

Check if your Docker image exists before you can process any deployment.
Currently the tool only supports `ECR`.

## Examples

```
$ go get github.com/icy/gDockerImageCheck

$ export AWS_PROFILE=my-profile

$ gDockerImageCheck /services/base/docker.go \
    1234.dkr.ecr.eu-west-1.amazonaws.com/devops:bar-v0.0.0 \
    541497436480.dkr.ecr.eu-west-1.amazonaws.com/devops:bar-v1.1.0

:: Error: 'InvalidParameterException: Invalid parameter at 'registryId' failed to satisfy constraint: 'must satisfy regular expression [0-9]{12}'
        status code: 400, request id: e89c68bb-47d8-4769-aece-43e4714c9596'.
:: Info: Found ECR image devops:bar-v1.1.0

$ echo $?
1
```

## Usage

The program accepts a list of docker images from command line,
and check if any of them exists. The tool returns zero if all images
can be found on the registry.

`ECR` docker image would be in the following format

```
<ECR_ID>.dkr.ecr.<AWS_REGION>.amazonaws.com/<REPO>:<TAG>
```

Any other argument that doesn't conform with this format will be ignored.
