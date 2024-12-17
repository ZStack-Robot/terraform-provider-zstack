Terraform Provider For ZStack Cloud
==================

- [Getting Started](https://registry.terraform.io/providers/ZStack-Robot/zstack/latest)
- Usage
  - Documentation
  - [Examples](https://github.com/ZStack-Robot/terraform-provider-zstack/blob/main/docs/index.md)

The ZStack provider is used to interact with the resources supported by ZStack Cloud, a powerful cloud management platform. 
This provider allows you to manage various cloud resources such as virtual machines, networks, storage, and more. 
It provides a seamless integration with Terraform, enabling you to define and manage your cloud infrastructure as code.

Supported Versions
------------------

| Terraform version | minimum provider version |maximum provider version
| ---- | ---- | ----| 
| >= 1.5.x	| 1.0.0	| latest |

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 1.5.x
-	[Go](https://golang.org/doc/install) 1.22 (to build the provider plugin)


Building The Provider
---------------------


Using the provider
----------------------
Please see [instructions](https://www.zstack.io) on how to configure the ZStack Cloud Provider.


## Contributing to the provider

The ZStack Provider for Terraform is the work of many contributors. We appreciate your help!

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.20+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.



You may also [report an issue](https://github.com/ZStack-Robot/terraform-provider-zstack/issues/new). 

    Here is an example of how to clone this repository and switch to the directory:

    ```console
    $ git clone https://github.com/ZStack-Robot/terraform-provider-zstack.git
    $ cd terraform-provider-zstack
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
