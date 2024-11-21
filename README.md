Terraform Provider For ZStack Cloud
==================

- Tutorials: [learn.hashicorp.com](https://learn.hashicorp.com/terraform?track=getting-started#getting-started)
- Documentation: 


Supported Versions
------------------

| Terraform version | minimum provider version |maximum provider version
| ---- | ---- | ----| 
| >= 1.5.x	| 1.0.0	| latest |

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 1.5.x
-	[Go](https://golang.org/doc/install) 1.20 (to build the provider plugin)


Building The Provider
---------------------


Using the provider
----------------------
Please see [instructions](https://www.zstack.io) on how to configure the ZStack Cloud Provider.


## Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.20+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-zstack
...
```

Running `make dev` or `make devlinux` or `devwin` will only build the specified developing provider which matchs the local system.
And then, it will unarchive the provider binary and then replace the local provider plugin.

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```

## Acceptance Testing
Before making a release, the resources and data sources are tested automatically with acceptance tests (the tests are located in the zstack/*_test.go files).
You can run them by entering the following instructions in a terminal:
```
cd $GOPATH/src/xxxx/zstack/terraform-provider-zstack
export ZSTACK_HOST=xxx
export ZSTACK_ACCOUNT_NAME=xxx
export ZSTACK_ACCOUNTP_ASSWORD=xxx
export ZSTACK_ACCESS_KEY_ID=xxx
export ZSTACK_ACCESS_KEY_SECRET=xxx
export outfile=gotest.out


```
