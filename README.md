# Warmup kubernetes pod
Standart rediness probe https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/#define-readiness-probes - runs on the container during its whole lifecycle.

If your kubernetes cluster has no enabled feature gates StartupProbe - you can add sidecar container to your pod to make readiness probes while main container answers HTTP 200 - if main container answers correctly - probes will ends
```
- name: warmup
  image: paskalmaksim/warmup:v0.0.7
  imagePullPolicy: IfNotPresent
  resources:
    requests:
      cpu: 10m
      memory: 30Mi
  command:
  - warmup
  - --url=http://127.0.0.1:3000
  - --result.file=/tmp/readinessProbeStatus
  - --http.timeout=1s
  - --try.timeout=1s
  readinessProbe:
    httpGet:
      path: /ready
      port: 12380
    initialDelaySeconds: 1
    periodSeconds: 5
  livenessProbe:
    httpGet:
      path: /healthz
      port: 12380
    initialDelaySeconds: 10
    periodSeconds: 10
```