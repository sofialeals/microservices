# Microsserviços - Instruções Rápidas

Este projeto contém três microsserviços: **payment**, **order** e **shipping**.  
Os serviços são containerizados com Docker e orquestrados no Kubernetes com Kustomize.

---

## Rodando todos os microsserviços

Na raiz do projeto (onde está o `kustomization.yaml`), execute:

```bash
kubectl apply -k .
```

## Verificando os serviços
```bash
kubectl get pods
kubectl get svc
```

## Deletando todos os microsserviços
```bash
kubectl delete -k .
```