<img width="1470" alt="Screenshot 2025-01-01 at 11 45 15 PM" src="https://github.com/user-attachments/assets/088900db-cac9-46a2-8908-07ac122a482a" />
## 📚 Kubernetes RBAC: Team Guide

### 📑 Table of Contents
- [📚 Kubernetes RBAC: Team Guide](#-kubernetes-rbac-team-guide)
  - [📑 Table of Contents](#-table-of-contents)
  - [🎯 Overview](#-overview)
  - [🔑 1. What is Kubernetes RBAC?](#-1-what-is-kubernetes-rbac)
  - [🚂 2. Authentication Methods](#-2-authentication-methods)
  - [🎭 3. Role \& RoleBinding (Namespace Scope)](#-3-role--rolebinding-namespace-scope)
  - [🌐 4. ClusterRole \& ClusterRoleBinding (Cluster Scope)](#-4-clusterrole--clusterrolebinding-cluster-scope)
  - [🔧 5. Token-based Authentication Example](#-5-token-based-authentication-example)
    - [5.1 Declare \& Create a ServiceAccount 🆔](#51-declare--create-a-serviceaccount-)
    - [5.2 Manually Create a Secret for the SA 🔑](#52-manually-create-a-secret-for-the-sa-)
    - [5.3 Inspect the Token \& Use It 🌐](#53-inspect-the-token--use-it-)
  - [📊 6. Demonstrating Access Levels](#-6-demonstrating-access-levels)
    - [6.1 Cluster-Level Access](#61-cluster-level-access)
    - [6.2 Namespace-Level Access](#62-namespace-level-access)
  - [📄 7. RoleBinding vs ClusterRoleBinding for Built-in Roles](#-7-rolebinding-vs-clusterrolebinding-for-built-in-roles)
  - [🛡️ 8. Best Practices \& Tips](#️-8-best-practices--tips)

---

### 🎯 Overview

Kubernetes Role​‑Based Access Control (RBAC) lets you define **who** (users, groups, or service accounts) can perform **what actions** (verbs) on **which resources** (pods, secrets, nodes, etc.), and **where** (namespaces or cluster-wide). This guide is tailored for newcomers setting up RBAC for a team.

---

### 🔑 1. What is Kubernetes RBAC?
RBAC is a **gatekeeper** that:
1. **Authenticates** the actor (verifies identity: human or machine).
2. **Authorizes** the action (checks if the actor has permission to perform the requested operation).

Without RBAC, **any** authenticated actor could do **anything** against the API.

---

### 🚂 2. Authentication Methods
Kubernetes supports various ways to prove identity:

| 🔐 Mechanism             | 📖 Description                                     | ⚙️ Use Case                             |
|--------------------------|----------------------------------------------------|-----------------------------------------|
| Certificate-based (mTLS) | Client TLS certs signed by a trusted CA           | High-security clusters                  |
| Token-based (SA or JWT)  | Bearer tokens injected into pods or delivered to users | Automation, CI/CD pipelines      |
| OIDC / OpenID Connect    | SSO via external providers (Dex, Keycloak, etc.)   | Enterprise SSO                          |
| Webhook Token Validation | Delegate token checks to an external webhook       | Custom auth logic                       |

---

### 🎭 3. Role & RoleBinding (Namespace Scope)

**Role**  
A namespaced object that **groups permissions** (verbs on resources) **within** a single namespace. Think of it as a job description scoped to one team space.

**RoleBinding**  
Associates that Role to specific **subjects** (Users, Groups, ServiceAccounts) inside the same namespace. You **must** create the Role first, then bind it.

```yaml
# role-pod-reader.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: pod-reader       # Name of your Role
  namespace: dev      # Scope: dev namespace
rules:
- apiGroups: [""]      # "" means the core API group
  resources: ["pods"]
  verbs: ["get", "watch", "list"]
```
```yaml
# binding-pod-reader.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: pod-reader-binding
  namespace: dev
- kind: ServiceAccount
  name: demo-sa    # name of service account
  namespace: dev
roleRef:
  kind: Role
  name: pod-reader
  apiGroup: rbac.authorization.k8s.io
```

---

### 🌐 4. ClusterRole & ClusterRoleBinding (Cluster Scope)

**ClusterRole**  
Similar to a Role, but **cluster​‑wide** (all namespaces) **or** for non​‑namespaced resources (nodes, certificates, CRDs).

**ClusterRoleBinding**  
Binds a ClusterRole to subjects **across** the entire cluster. You still create the ClusterRole first, then bind it.

```yaml
# clusterrole-node-reader.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: node-reader      # A reusable, cluster-wide permission set
rules:
- apiGroups: [""]
  resources: ["nodes"]
  verbs: ["get", "watch", "list"]
```
```yaml
# crb-node-reader.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: infra-node-reader
subjects:
- kind: ServiceAccount
  name: demo-sa     # name of service account
  namespace: dev
roleRef:
  kind: ClusterRole
  name: node-reader
  apiGroup: rbac.authorization.k8s.io
```

---

### 🔧 5. Token-based Authentication Example
Below we use **ServiceAccounts (SA)** to obtain tokens for API access. _Note: As of Kubernetes 1.24+, you must create SA secrets manually._

#### 5.1 Declare & Create a ServiceAccount 🆔
```yaml
# demo.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: demo-sa
  namespace: dev
```
```bash
kubectl apply -f demo.yaml
```

- **ServiceAccount**: machine identity for pods or external tools. Lives in a namespace.

#### 5.2 Manually Create a Secret for the SA 🔑
```yaml
# demo-token.yaml
apiVersion: v1
kind: Secret
metadata:
  name: demo-sa-token
  namespace: dev
  annotations:
    kubernetes.io/service-account.name: demo-sa
type: kubernetes.io/service-account-token
```
```bash
kubectl apply -f demo-token.yaml
```

- **Why?** Generates a JWT token + CA bundle for the SA. Kubernetes auto​‑populates `data.token` and `data["ca.crt"]`.

#### 5.3 Inspect the Token & Use It 🌐
```bash
# Retrieve the token value
kubectl get secret demo-sa-token -n dev \
  -o jsonpath='{.data.token}' | base64 --decode
```

> ⚠️ **NOTE**: This token in going to be used in **`kube-config-file`** to access the reources

---

### 📊 6. Demonstrating Access Levels

#### 6.1 Cluster-Level Access
1. **Custom ClusterRole**: e.g. `node-reader` (see above).
2. **ClusterRoleBinding**: bind `node-reader` to `demo-sa`:
   ```yaml
   # crb-mytool-node.yaml
   apiVersion: rbac.authorization.k8s.io/v1
   kind: ClusterRoleBinding
   metadata:
     name: mytool-node-reader
   subjects:
   - kind: ServiceAccount
     name: demo-sa
     namespace: dev
   roleRef:
     kind: ClusterRole
     name: node-reader
     apiGroup: rbac.authorization.k8s.io
   ```
3. **Built-in Roles**: you can also bind `view`, `edit`, `admin`, or `cluster-admin`:
   ```yaml
   kind: ClusterRoleBinding
   metadata:
     name: mytool-sa-view-all
   subjects:
   - kind: ServiceAccount
     name: demo-sa
     namespace: dev
   roleRef:
     kind: ClusterRole
     name: view
     apiGroup: rbac.authorization.k8s.io
   ```

| Role Type         | Scope              | Example Built​‑in                |
|-------------------|--------------------|---------------------------------|
| Custom ClusterRole| Any or non-namespaced | `node-reader`, `db-writer`    |
| Built-in ClusterRole| Curated cluster perms | `view`, `edit`, `admin`, `cluster-admin` |

#### 6.2 Namespace-Level Access
1. **Custom Role** `cm-editor` in `dev` (see section 3 example).
2. **RoleBinding**: bind it to `demo-sa` in `dev`.
3. **Built-in Role**: bind `edit` within the namespace:
   ```yaml
   apiVersion: rbac.authorization.k8s.io/v1
   kind: RoleBinding
   metadata:
     name: mytool-sa-edit
     namespace: dev
   subjects:
   - kind: ServiceAccount
     name: demo-sa
     namespace: dev
   roleRef:
     kind: ClusterRole   # built-in roles are ClusterRoles
     name: edit
     apiGroup: rbac.authorization.k8s.io
   ```

---

### 📄 7. RoleBinding vs ClusterRoleBinding for Built-in Roles

| Binding Type        | Grants Built​‑in Role | Scope                | Use Case                             |
|---------------------|----------------------|----------------------|--------------------------------------|
| RoleBinding         | e.g. `edit`          | Single namespace     | Team-specific edit permissions       |
| ClusterRoleBinding  | e.g. `edit`          | All namespaces       | Globally allow edit on every namespace |

---

### 🛡️ 8. Best Practices & Tips
- 🎯 **Principle of Least Privilege**: Grant only necessary verbs and resources.
- ⬆️ **Explicit Secrets**: As of 1.24+, create SA secrets manually for clarity and rotation.
- 📝 **Version Control**: Keep RBAC manifests in Git.
- 🔍 **Audit**: `kubectl get roles,rolebindings,clusterroles,clusterrolebindings -o yaml` regularly.
- 🛠️ **Policy Enforcement**: Use OPA/Gatekeeper or Kyverno to validate & enforce RBAC standards.

---
❤️‍🔥 Happy securing your cluster! 🚀
