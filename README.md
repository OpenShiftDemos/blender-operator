# Blender Remote Render Operator for Kubernetes
This is an attempt to write an operator to handle remote rendering in Kubernetes.

It is closely tied to https://github.com/openshiftdemos/blender-remote which is
the container that does the render work.

The Operator is essentially just a convenient front-end (via a
CustomResourceDefinition) to ask for a render job.

## Example CatalogSource
```yaml
apiVersion: operators.coreos.com/v1alpha1
kind: CatalogSource
metadata:
  name: blender-operator
spec:
  displayName: Blender Operator
  image: quay.io/openshiftdemos/blender-operator-catalog:v0.0.1
  sourceType: grpc
```

## Example OperatorGroup
```yaml
apiVersion: operators.coreos.com/v1
kind: OperatorGroup
metadata:
  name: blender-group
spec:
```

## Example Subscription
```yaml
apiVersion: operators.coreos.com/v1alpha1
kind: Subscription
metadata:
  name: blender-operator
spec:
  channel: "alpha"
  installPlanApproval: Automatic
  name: blender-operator
  source: blender-operator
  sourceNamespace: olm
  startingCSV: blender-operator.v0.0.1
```