# Scdbeat

Welcome to Scdbeat.

Ensure that this folder is at the following location:
`${GOPATH}/github.com/junglestory`

## Getting Started with Scdbeat

### Requirements

* [Golang](https://golang.org/dl/) 1.6
* [Glide](https://github.com/Masterminds/glide) >= 0.10.0

### Init Project
To get running with Scdbeat, run the following command:

```
make init
```

To commit the first version before you modify it, run:

```
make commit
```

It will create a clean git history for each major step. Note that you can always rewrite the history if you wish before pushing your changes.

To push Scdbeat in the git repository, run the following commands:

```
git remote set-url origin https://github.com/junglestory/scdbeat
git push origin master
```

For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).

### Build

To build the binary for Scdbeat run the command below. This will generate a binary
in the same directory with the name scdbeat.

```
make
```


### Run

To run Scdbeat with debugging output enabled, run:

```
./scdbeat -c scdbeat.yml -e -d "*"
```


### Test

To test Scdbeat, run the following command:

```
make testsuite
```

alternatively:
```
make unit-tests
make system-tests
make integration-tests
make coverage-report
```

The test coverage is reported in the folder `./build/coverage/`


### Package

To be able to package Scdbeat the requirements are as follows:

 * [Docker Environment](https://docs.docker.com/engine/installation/) >= 1.10
 * $GOPATH/bin must be part of $PATH: `export PATH=${PATH}:${GOPATH}/bin`

To cross-compile and package Scdbeat for all supported platforms, run the following commands:

```
cd dev-tools/packer
make deps
make images
make
```

### Update

Each beat has a template for the mapping in elasticsearch and a documentation for the fields
which is automatically generated based on `etc/fields.yml`.
To generate etc/scdbeat.template.json and etc/scdbeat.asciidoc

```
make update
```


### Cleanup

To clean  Scdbeat source code, run the following commands:

```
make fmt
make simplify
```

To clean up the build directory and generated artifacts, run:

```
make clean
```


### Clone

To clone Scdbeat from the git repository, run the following commands:

```
mkdir -p ${GOPATH}/github.com/junglestory
cd ${GOPATH}/github.com/junglestory
git clone https://github.com/junglestory/scdbeat
```


For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).
