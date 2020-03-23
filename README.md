This repository is used for testing of CDI imageio import.

Deployment example on k8s:

```yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: fakeovirt
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: fakeovirt
  template:
    metadata:
      labels:
        app: fakeovirt
    spec:
      containers:
      - name: fakeovirt
        image: machacekondra/fakeovirt
        ports:
        - containerPort: 12346
        env:
          - name: NAMESPACE
            value: default
          - name: PORT
            value: 12346
---
apiVersion: v1
kind: Service
metadata:
  name: fakeovirt
  namespace: default
spec:
  selector:
    app: fakeovirt
  type: NodePort
  ports:
  - name: fakeovirt
    port: 12346
      targetPort: 12346
```
