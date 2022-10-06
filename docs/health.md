Health checks are powered by 
[bufbuild/connect-grpchealth-go](https://github.com/bufbuild/connect-grpchealth-go),
which exposes a health checking API that is wire-compatible with the well-known
[gRPC implementation](https://github.com/grpc/grpc/blob/master/doc/health-checking.md)
so it works with 
[grpcurl](https://github.com/fullstorydev/grpcurl),
[grpc-health-probe](https://github.com/grpc-ecosystem/grpc-health-probe/), and 
[Kubernetes gRPC liveness probes](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/#define-a-grpc-liveness-probe).

You can hit the health check API on port 8081.
```bash
http -b POST http://localhost:8081/grpc.health.v1.Health/Check service="coop.drivers.foov1beta1.FooService"
```
which should return
```json
{
    "status": "SERVING_STATUS_SERVING"
}
```