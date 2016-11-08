# The Steward Framework

This repository contains the Steward Framework for creating service catalog controllers.

The framework utilizes inversion of control, implementing all the core concerns of a service catalog controller, whilst delegating implementation-specific operations to the client library.

The framework implements all of the following concerns:

* __Service catalog publishing:__ The framework publishes `servicecatalogentries` that reflect a backing broker's offerings into the Kubernetes cluster. The query for the backing broker's offerings is delegated to the backing broker through the client library's implementation of the `framework.Cataloger` interface.

* __Event and control loops:__ The framework watches the Kubernetes event stream for changes to `ServicePlanClaims`. Events that trigger discreet provision, bind, unbind, or deprovision actions delegate those operations to the backing broker through the client library's implementation of the `framework.Lifecycler` interface.

* __API server__: At this time, the API server implements only a health-check endpoint at `/healthz`.

## Implementation guide

Utilizing the Steward Framework to implement your own service catalog controller requires the following steps:

1. Include the `github.com/deis/steward-framework` package in your project using your dependency management system of choice-- for example, [glide](https://github.com/masterminds/glide) or [godeps](https://github.com/tools/godep).

2. Implement the `framework.Cataloger` and `framework.Lifecycler interfaces`. For your convenience, those interfaces are shown below:

    ```go
    type Cataloger interface {

      List(ctx context.Context) ([]*framework.Service, error)

    }

    type Lifecycler interface {

      Provision(
        ctx context.Context,
        req *framework.ProvisionRequest,
      ) (*framework.ProvisionResponse, error)

      Bind(
        ctx context.Context,
        req *framework.BindRequest,
      ) (*framework.BindResponse, error)

      Unbind(
        ctx context.Context,
        req *framework.UnbindRequest,
      ) error

      Deprovision(
        ctx context.Context,
        req *framework.DeprovisionRequest,
      ) (*framework.DeprovisionResponse, error)

      GetOperationStatus(
        ctx context.Context,
        req *framework.OperationStatusRequest,
      ) (*framework.OperationStatusResponse, error)

    }

    ```

3. Import `github.com/deis/steward-framework/runner`

4. In your `main()` function, call the blocking `runner.Run(...)`. For convenience, the signature of that function is shown below:

    ```go
    func Run(
	  brokerName string,
	  namespaces []string,
	  cataloger framework.Cataloger,
	  lifecycler framework.Lifecycler,
	  maxAsyncDuration time.Duration,
	  apiPort int,
    ) error
    ```

## Example code

The following annotated example is adapted from [steward-cf](https://github.com/deis/steward-cf):

```go
import "github.com/deis/steward-framework/runner"

// ...

func main() {

  // A configuration object is built from environment variables.
  cfg := ...

  // Relevant implementations of the framework.Cataloger and
  // framework.Lifecycler interfaces are initialized.
  cataloger, lifecycler := ...

  // Configuration details, cataloger, and lifecycler are passed to runner.Run().
  // The Steward Framework takes it from there. The call to runner.Run() will
  // block infinitely or until the framework encounters a fatal error.
  if err = runner.Run(
    cfg.BrokerName,
    cfg.Namespaces,
    cataloger,
    lifecycler,
    cfg.getMaxAsyncDuration(),
    cfg.APIPort,
  ); err != nil {
    logger.Criticalf("error running steward-framework: %s", err)
}
```
