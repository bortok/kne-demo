# Network Emulation with SONiC

## Prerequisites for Mac

1. Multipass

```Shell
brew install --cask multipass
multipass version
````

## Build SONiC docker image

1. Use https://github.com/antongisli/sonic-builder to automate build environment. Create a build VM

```Shell
cd athena-mac
export BASEDIR=`pwd`
git clone https://github.com/antongisli/sonic-builder.git
multipass launch 20.04 -n sonic-builder -c4 -m16G -d300G --cloud-init sonic-builder/sonic-cloud-init.yaml
````

2. Build a SONiC image

```Shell
multipass shell sonic-builder
./build-vs-docker.sh
````

3. Copy the image from `sonic-builder` locally, move it to a destination host and load into a local docker registry

```Shell
multipass transfer sonic-builder:sonic-buildimage/target/docker-sonic-vs.gz docker-sonic-vs_20220111.gz
scp docker-sonic-vs_20220111.gz DESTINATION_HOST:
ssh DESTINATION_HOST
docker load < docker-sonic-vs_20220111.gz
docker tag docker-sonic-vs:latest docker-sonic-vs:20220111
````
