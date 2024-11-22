# AMS SDK

Anbox Management SDK offers a group of [Go](https://golang.org/) packages and
utilities for any external [Go](https://golang.org/) code to be able to connect
to the AMS service through the exposed REST API to manage Anbox Cloud.

the AMS SDK consist of the following Go packages:

* `client`: The main purpose of this package is to provide a REST client object
   with  relevant method to manage applications, images, nodes, addons or
   containers. Each method relates to a specific management operation and wraps
   the REST calls and listening operations

* `api`: AMS REST API objects

* `shared`: Helper methods and tools for common tasks like system tasks,
  certificates, password hashing or websocket dialing.

* `shared/rest/client`: REST client base functionality to wrap common behavior
  of the various operations of the REST protocol and websocket management for
  event listening.

* `shared/rest/api`: REST API basic objects, independent of any specific REST
  implementation

* `shared/errors`: A simple wrapper for the most commonly-used error implementation
   in the return of REST API

* `examples`: A set of examples to demonstrate how the `client` package can be used.


## Install

There are no special instructions to install the AMS SDK. You can simply add the
content of the provided SDK zip file into your projects `vendor/` directory or
your `GOPATH` and start using it.

## Examples

The SDK comes with a set of examples demonstrating the capabilities of the SDK.

## Authentication setup

Go clients using the SDK must authenticate to an AMS instance by using two way SSL
Authentication. This means that server requires a trusted client certificate to be
sent in every request and that the client must add a certificate he desires to use
to the server before sending the first request.

Assuming that you already have generated a
[client certificate](#generating-a-client-certificate) and it is stored in a file
called `client.crt`, you can easily add it to the list of server trusted certificates
using [juju](https://jujucharms.com/) on an already deployed environment:

    $ juju run-action ams/0 add-certificate cert="$(cat client.crt | grep -v ^-- | tr -d '\n')" --wait

Once the execution completes, the certificate is trusted and can be used along with
the SDK libraries and tools to connect to the AMS instance or AMS cluster.

!!!Note:
    `ams/0` is the name of the unit on which AMS has been deployed. Its number can
    be different than `0`.


## Generating a client certificate

### Self signed

The easiest way of generating a client certificate is using the certificate as
CA of itself. That is what a self-signed certificate is. This is generally intended
for testing purposes because the CA is not a well-known trusted one, though in the
case of AMS it really does not matter at all because the service trust in individual
certificates regardless of the CA signing them. So, in the end a certificate provided
this way is as valid as one signed by a CA.

You can easily create it using [OpenSSL](https://www.openssl.org/):

    $ openssl req -x509 -newkey rsa:4096 -keyout client.key -out client.crt -days 365

You will be asked for a password and certificate information interactively. Answer the
questions until completing all the process and the certificate and key will be generated:

    Generating a 4096 bit RSA private key
    ............................................++
    ...................................................................++
    writing new private key to 'client.key'
    Enter PEM pass phrase:
    Verifying - Enter PEM pass phrase:
    -----
    You are about to be asked to enter information that will be incorporated
    into your certificate request.
    What you are about to enter is what is called a Distinguished Name or a DN.
    There are quite a few fields but you can leave some blank
    For some fields there will be a default value,
    If you enter '.', the field will be left blank.
    -----
    Country Name (2 letter code) [AU]:
    State or Province Name (full name) [Some-State]:
    Locality Name (eg, city) []:
    Organization Name (eg, company) [Internet Widgits Pty Ltd]:
    Organizational Unit Name (eg, section) []:
    Common Name (e.g. server FQDN or YOUR name) []:
    Email Address []:

### Signed by a CA

If you prefer, you can provide a certificate signed by a CA. The authority singing the
certificate can be a well-known one or a customer.

#### Creating a CA

You can create a custom CA key to sign your own certificates with:

    $ openssl genrsa -out MyCompanyCA.key 2048

That will generate the CA key in a file named `MyCompanyCA.key`. You can now create
the certificate with 

    $ openssl req -x509 -new -nodes -key MyCompanyCA.key -sha256 -days 1024 -out MyCompanyCA.crt

Again, you have to fill all interactively requested data until completing the process
of certificate creation. The result is a file named `MyCompanyCA.crt`


#### Generating certificates signed by a custom CA

Once you have a custom CA, to sign a certificate you first need to create the
certificate key and a request to be signed:

    $ openssl genrsa -out client.key 2048
    $ openssl req -new -key client.key -out client.csr

The file `client.csr` needs to be signed by the CA in order to get the
final certificate:

    $ openssl x509 -req -in client.csr -CA MyCompanyCA.crt -CAkey MyCompanyCA.key -CAcreateserial -out client.crt -days 1024 -sha256

you obtain as result `client.crt` signed by MyCompanyCA certificate authority.

!!!Note:
    Check [the OpenSSL website](https://www.openssl.org/) for more information about
    the specific parameters for every OpenSSL operation.

## Custom Client Code

### Connection

Main steps for custom code to connect to the server start with the creation of
a REST client object. Such object needs a TLS configuration including the client
certificate to be sent to AMS and the server certificate the client trusts. There
are many ways of creating a TLS configuration in go. AMS-SDK provides an easy
solution involving a few lines of code:

```go
    import (
        "flag"
        "net/url"
        "os"

        "github.com/anbox-cloud/ams-sdk/client"
        "github.com/canonical/ams/pkg/network"
    )

    func main() {
        flag.Parse()
        if flag.NArg() == 0 {
            fmt.Println("Please provide AMS service URL")
            os.Exit(1)
        }

        serviceURL := flag.Arg(0)
        u, err := url.Parse(serviceURL)
        if err != nil {
            fmt.Println("Failed to parse AMS service URL")
            os.Exit(1)
        }

        serverCert, err := network.GetRemoteCertificate(serviceURL)
        if err != nil {
            fmt.Println("Failed to get remote certificates")
            os.Exit(1)
        }

        tlsConfig, err := network.GetTLSConfig(clientCert, clientKey, "", serverCert)
        if err != nil {
            fmt.Println("Failed to get TLS config")
            os.Exit(1)
        }

        ...
    }
```

!!!Note:
    Here, we take any server certificate as valid. In case you want a better compromise
    on the client side with the server certificate to trust, you simply have to 
    replace `network.GetRemoteCertificate(serviceURL)` method with code to read a server well-known certificate from a remote or local path to a x509 object and pass it to `network.GetTLSConfig()` method.

Once the TLS configuration is ready, the next step is to create the REST client object:

```go
    amsClient, err := client.New(u, tlsConfig)
    if err != nil {
        return err
    }
```

### Using the REST API

The created REST client object has a large list of methods to manage the different entities
into AMS:

* Applications:
  * Create
  * Export
  * Update
  * List
  * Show
  * Publish
  * Revoke
  * Delete
  * DeleteVersion

* Containers:
  * Launch
  * List
  * Show
  * Delete

* Nodes:
  * Add
  * List
  * Update
  * Show
  * Remove

* Images:
  * Add
  * Update
  * List
  * Delete
  * DeleteVersion

* Addons:
  * Add
  * Update
  * List
  * Delete
  * DeleteVersion

Look at the godoc generated documentation for more details of each method call.
One example of the use of the SDK client to create an application could be
either from the folder path where the application package layout is :

```go
    c.CreateApplication(".", nil)
```

or from the application package path

```go
    c.CreateApplication("application.tar.bz2", nil)
```


## Asynchronous operations

All operations modifying entities on AMS are executed asynchronously to prevent
blocking the client. This means that a call, say, to `c.CreateApplication(...)`
won't block and will return immediately, even when the operation is still not
finished in the server.

All the asynchronous operations return an
`github.com/anbox-cloud/ams-sdk/shared/rest/client/Operation` struct object

If you want your client to wait for an asynchronous operation to complete, you
can call `Operation.Wait()` method, which will block current thread until the
operation finishes or an error occurs:

```go
    ...

    operation, err := c.CreateApplication(".", nil)
    err = operation.Wait(context.Background())
    if err != nil {
        return err
    }
```

You can get in your code the operation resultant resources:

```go
    operation.Get().Resources
```
