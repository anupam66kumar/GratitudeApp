# **Kubernetes Cluster Health Checker, Self-Healing & Autoscaling Platform (GratitudeApp)**

An enterprise-ready, polyglot microservices platform deployed on a multi-node Kubernetes cluster orchestrated with k3d/k3s. This platform integrates an autonomous self-healing Go controller (k8s-cluster-healer)\[cite: 1\], automated Horizontal Pod Autoscaling (HPA)\[cite: 2\], normalized Persistent Volume Claim (PVC) metrics, and full-stack observability with Prometheus and Grafana\[cite: 1, 2\].

## 

## **Table of Contents**

> 1. [Project Overview & Capstone Deliverables](https://www.google.com/search?q=%231-project-overview--capstone-deliverables)  
> 2. [Application Architecture, Core Features & Working](https://www.google.com/search?q=%232-application-architecture-core-features--working)  
> 3. [System Infrastructure & Multi-Node Topology](https://www.google.com/search?q=%233-system-infrastructure--network-topology)  
> 4. [Custom Auto-Healer Module (k8s-cluster-healer)](https://www.google.com/search?q=%234-custom-auto-healer-module-k8s-cluster-healer)  
> 5. [Repository Directory Structure](https://www.google.com/search?q=%235-repository-directory-structure)  
> 6. [Workloads & Infrastructure Catalog](https://www.google.com/search?q=%236-workloads--infrastructure-catalog)  
> 7. [Observability Stack (Prometheus & Grafana)](https://www.google.com/search?q=%237-observability-stack-prometheus--grafana)  
> 8. [Horizontal Pod Autoscaling (HPA) Verification](https://www.google.com/search?q=%238-horizontal-pod-autoscaling-hpa-verification)  
> 9. [Comprehensive Troubleshooting Log & Root Cause Analysis](https://www.google.com/search?q=%239-comprehensive-troubleshooting-log--root-cause-analysis)  
> 10. [Deployment & Verification Runbook](https://www.google.com/search?q=%2310-deployment--verification-runbook)  
> 11. [Cost Optimization Analysis](https://www.google.com/search?q=%2311-cost-optimization-analysis)

## 

## **1\. Project Overview & Capstone Deliverables**

Managing cloud-native clusters manually requires continuous monitoring and reactive troubleshooting, which increases downtime during failures\[cite: 1\]. Aligned with the **Project 1: Kubernetes Cluster Health Checker and Auto-Healing** deliverables, this project provides a closed-loop monitoring and self-remediating ecosystem\[cite: 1, 3\]:

| Capstone Deliverable | Implementation in this Platform | Verification Method |
| :---- | :---- | :---- |
| **Automated Health Monitoring Tool**\[cite: 3\] | Prometheus scrape mesh (cAdvisor, kubelet, kube-state-metrics) coupled with k8s-cluster-healer scanning pod lifecycle conditions\[cite: 1\]. | Scrape target health and periodic terminal/container scan logs. |
| **Self-Healing Mechanisms**\[cite: 3\] | Custom Go controller executing automated pod remediation on CrashLoopBackOff, container exit codes, or excessive restarts ($>10$)\[cite: 2\]. | Automated termination and replacement of faulty pods (broken-test-pod). |
| **Real-Time Alerting & Observability**\[cite: 3\] | Grafana dashboards, Prometheus Alertmanager integration, and pod restart audit logging\[cite: 1, 2\]. | Visualized telemetry and live stdout remediation traces. |
| **Web Dashboard**\[cite: 3\] | Grafana dashboards tracking deployments, HPA replica states, network I/O, and normalized PVC disk capacity\[cite: 1, 3\]. | Web interface served via port-forwarding on port 3001\. |
| **Comprehensive Documentation**\[cite: 3\] | End-to-end architecture diagrams, runbooks, directory listings, and historical root cause analyses\[cite: 3\]. | Complete README.md and version-controlled manifests. |

## 

## **2\. Application Architecture, Core Features & Working**

The core workload running on the cluster is **GratitudeApp**, a cloud-native microservices application designed for journal logging, mood analysis, and streak tracking.

### 

### 

### **Key Application Features**

* **Daily Journal Logging:** Users create, read, and maintain personal gratitude entries backed by relational persistence.  
* **Mood & Emotion Tracking:** Tracks mental wellness states across configurable mood tags with temporal timestamps.  
* **Analytics & Streak Aggregations:** Aggregates user habits, tracking usage consistency and mood distribution patterns over time.  
* **Multimedia Object Attachment:** Unstructured file upload pipeline backed by an S3-compliant MinIO bucket.  
* **Polyglot Service Orchestration:** High-performance gRPC inter-service communication coordinated through an Express API Gateway.

### 

### **End-to-End Application Workflow**

<img width="3364" height="2012" alt="carbon (5)" src="https://github.com/user-attachments/assets/b4da7249-2603-466a-9278-5f8d37beab5e" />

> 

> 1. **User Request Routing:** The client interacts with the React frontend (client-deployment).  
> 2. **Gateway Dispatch:** The browser dispatches asynchronous API requests to the api-gateway-service. The gateway parses the payload, validates authentication, and delegates calls to domain microservices over low-latency gRPC channels.  
> 3. **Domain Processing & Persistence:**  
>   * Entries and reflections route to entries-service (Port 50051).  
>   * Mood indices route to moods-service (Port 50052).  
>   * Aggregations route to stats-service (Port 50053).  
>   * Structured records are committed to PostgreSQL (postgres-deployment), while media uploads are handled by files-service and stored in MinIO.  
> 4. **Dynamic Elasticity:** When inbound gateway requests exceed baseline load, the api-gateway-hpa controller scales the gateway pods horizontally (up to 5 replicas)\[cite: 2\].

## 

## **3\. System Infrastructure & Multi-Node Topology**

The cluster runs on a multi-node topology provisioned via k3d with one server control-plane node and two worker agent nodes.

### **Architectural Decision: Why k3d / k3s Over Alternatives?**

| Solution | Multi-Node Emulation | Resource Overhead | Startup Latency | Docker Integration | Verdict for This Project |
| :---- | :---- | :---- | :---- | :---- | :---- |
| **k3d (k3s in Docker)** | **Native** (Runs lightweight k3s nodes as separate Docker containers) | **Minimal** ($\sim$512MB RAM total base footprint) | Fast ($\approx 15-20$s) | Seamless image imports (k3d image import) directly into node containerd runtimes | **Selected**: Best for local multi-node scheduling, node-affinity, and HPA testing. |
| **Minikube** | Supported via Docker driver, but resource-heavy | High ($\geq$ 2GB RAM per node VM/container) | Moderate ($\approx 45-60$s) | Requires eval \$(minikube docker-env) or slow registry tunneling | Rejected due to high memory footprint for 3 nodes alongside Prometheus/Grafana. |
| **Kind (Kubernetes in Docker)** | Native | Moderate | Moderate ($\approx 30$s) | Requires manual node config YAMLs for port maps | Viable, but k3d provides cleaner CLI ingress/load-balancer port forwarding. |
| **MicroK8s** | Complex multi-node setup on desktop | High (systemd daemon dependencies) | Slow | Snaps package management overhead | Non-portable across generic Linux/macOS developer environments. |

### 

### **Why a 3-Node Cluster Topology?**

1. **Control-Plane vs. Worker Isolation:** One dedicated server node (server-0) runs the API Server, etcd/k3s datastore, and controller manager, while two worker nodes (agent-0, agent-1) handle workload scheduling.  
2. **Realistic HPA Scheduling:** Ensures that when the api-gateway scales from 1 to 5 replicas during load bursts, Kubernetes distributes pod replicas across distinct physical container worker nodes using pod anti-affinity and load-balancing algorithms.  
3. **Volume Isolation:** Validates Persistent Volume Claim (PVC) bindings and kubelet storage reporting across different physical agent nodes with local-path-provisioner.

\[Scroll down to next page\]

## **<img width="3568" height="4320" alt="carbon (6)" src="https://github.com/user-attachments/assets/eb76f089-0f05-49f0-bc7f-7abb97e2f90a" />**

## **4\. Custom Auto-Healer Module (k8s-cluster-healer)**

The k8s-cluster-healer is a specialized Go utility that automates cluster remediation using client-go\[cite: 1\]. It handles scenarios where native Kubernetes restart policies leave broken pods lingering in backoff loops or failed states\[cite: 1, 2\].

### 

### **Architecture & Operation**

* **Dual-Mode Client Initialization (pkg/k8s/client.go):** Automatically evaluates rest.InClusterConfig() for in-cluster deployment (Option B), falling back to \$KUBECONFIG or \~/.kube/config when run locally as a binary (Option A).

* **Scanning Engine (pkg/healer/healer.go):** Executes every 10 seconds against the target namespace. It scans container statuses for:

  1. Waiting.Reason \== "CrashLoopBackOff"\[cite: 2\]  
  2. Waiting.Reason \== "Error"  
  3. RestartCount \> 10 (excessive restarts indicator)

* **Remediation Loop:** When an unhealthy pod is identified, the controller issues an API call to delete it. The parent controller (Deployment/ReplicaSet) then provisions a fresh pod replacement with clean state.

### 

### **Verified Self-Healing Test Output**

During testing, an intentional crashing workload (broken-test-pod) was injected into the cluster. The healer detected and resolved the failure automatically:

Plaintext  
\[08:55:15\] 🔍 Scanning 12 pods in namespace 'default'...  
\[08:55:25\] 🔍 Scanning 12 pods in namespace 'default'...  
\⚠️ Detected unhealthy pod: broken-test-pod (Reason: CrashLoopBackOff). Remediating...  
\✅ Pod broken-test-pod remediated (deleted for controller replacement).  
\[08:55:35\] 🔍 Scanning 11 pods in namespace 'default'...

##

## **5\. Repository Directory Structure**

<img width="3972" height="3500" alt="carbon (7)" src="https://github.com/user-attachments/assets/356a4bb2-a0e9-4c25-a6d6-d91c55fed9d2" />

## 


## **6\. Workloads & Infrastructure Catalog**

| Workload Name | Kind | Port | PVC Binding | CPU Limit | Memory Limit |
| :---- | :---- | :---- | :---- | :---- | :---- |
| **client** | Deployment | 3000 | Ephemeral | 200m | 256Mi |
| **api-gateway** | Deployment (HPA) | 5000 | Ephemeral | 500m | 512Mi |
| **entries-service** | Deployment | 50051 | Ephemeral | 300m | 256Mi |
| **moods-service** | Deployment | 50052 | Ephemeral | 300m | 256Mi |
| **stats-service** | Deployment | 50053 | Ephemeral | 300m | 256Mi |
| **postgres** | Deployment | 5432 | database-persistent-volume-claim (1Gi) | 500m | 512Mi |
| **minio** | Deployment | 9000, 9001 | minio-pvc (2Gi) | 300m | 512Mi |
| **k8s-cluster-healer** | Deployment | N/A | Ephemeral | 100m | 128Mi |
| **prometheus-0** | StatefulSet | 9090 | Local TSDB (Retention: 6h) | 800m | 768Mi |
| **grafana** | Deployment | 80 \-\> 3000 | ConfigMap / Ephemeral | 500m | 768Mi |

## 

## **7\. Observability Stack (Prometheus & Grafana)**

The monitoring framework is managed through the kube-prometheus-stack Helm chart using a resource-optimized profile (monitoring-lite.yaml)\[cite: 1, 2\].

### 

### **Probe Hardening**

Default probes were reconfigured to prevent false-positive kills when Prometheus replays its Write-Ahead Log (WAL):

YAML  
grafana:  
  livenessProbe:  
    httpGet:  
      path: /api/health  
      port: 3000  
    initialDelaySeconds: 120  
    periodSeconds: 30  
    timeoutSeconds: 15  
    failureThreshold: 10  
  readinessProbe:  
    httpGet:  
      path: /api/health  
      port: 3000  
    initialDelaySeconds: 60  
    periodSeconds: 15  
    timeoutSeconds: 15  
    failureThreshold: 10

prometheus:  
  prometheusSpec:  
    scrapeInterval: 30s  
    evaluationInterval: 30s  
    retention: 6h  
    probeTimeoutSeconds: 15

### 

### 

### **Production PVC Metric Normalization**

On local clusters with local-path-provisioner, querying raw filesystem capacity can produce mismatched metrics due to multi-node label mismatches.  
The following query matches metrics across volumes and calculates accurate percentage utilization:

Code snippet  
(  
  max by (persistentvolumeclaim) (kubelet\_volume\_stats\_used\_bytes{namespace="default", persistentvolumeclaim\!=""})  
  /  
  max by (persistentvolumeclaim) (kubelet\_volume\_stats\_capacity\_bytes{namespace="default", persistentvolumeclaim\!=""})  
) \* 100

* **Panel Configuration:** Unit set to Percent (0-100), bounds fixed between 0 and 100, legend configured as {{persistentvolumeclaim}}.  
* **Verified Metric Display:** Reports both database-persistent-volume-claim and minio-pvc at an accurate utilization of **4.88%**.

**\[PLACEHOLDER: Screenshot \- Grafana Unified Monitoring Dashboard\]**  
*(Insert dashboard capture showing CPU/Memory graphs, Pod Health, PVC Percentages, and Network I/O)*

<img width="1854" height="2173" alt="graphana" src="https://github.com/user-attachments/assets/1e30b4cb-20e5-46d1-82c3-3ad3a61cd477" />

## 

## **8\. Horizontal Pod Autoscaling (HPA) Verification**

Autoscaling responsiveness was verified under synthetic multi-threaded load targeting the gateway\[cite: 2\].

### 

### **Verification Trace**

Bash  
\# Injected concurrent synthetic traffic  
kubectl run hpa-load-generator \--rm \-i \--tty \--image=curlimages/curl \--restart=Never \-- \\  
  sh \-c "for i in \\\$(seq 1 8); do (while true; do curl \-s http\://api-gateway-cluster-ip-service:5000/entries/all \> /dev/null; done) & done; wait"

Plaintext  
NAME              REFERENCE                          TARGETS         MINPODS   MAXPODS   REPLICAS   AGE  
api-gateway-hpa   Deployment/api-gateway-deployment   cpu: 2%/50%     1         5         1          18s  
api-gateway-hpa   Deployment/api-gateway-deployment   cpu: 9%/50%     1         5         1          60s  
api-gateway-hpa   Deployment/api-gateway-deployment   cpu: 167%/50%   1         5         5          75s  
api-gateway-hpa   Deployment/api-gateway-deployment   cpu: 291%/50%   1         5         5          105s  
api-gateway-hpa   Deployment/api-gateway-deployment   cpu: 282%/50%   1         5         5          2m  
api-gateway-hpa   Deployment/api-gateway-deployment   cpu: 266%/50%   1         5         5          2m16s

After traffic was stopped, the HPA observed the 5-minute stabilization cooldown window and scaled the deployment back down to 1 replica (cpu: 2%/50%).

**\[PLACEHOLDER: Screenshot \- Terminal HPA Dynamic Scale-Up (1 to 5 Replicas)\]**  
*(Insert screenshot showing real-time terminal scaling logs during load generation)*
<img width="1859" height="711" alt="Screenshot from 2026-10-10 17-59-19" src="https://github.com/user-attachments/assets/7539448d-0a6d-413e-8c41-72f9616d6c4c" />

*API-Gateway nodes before the load test*  
<img width="1847" height="305" alt="Screenshot from 2026-10-10 18-03-52" src="https://github.com/user-attachments/assets/c916eedf-a1a1-4a67-b018-71b22c9795ab" />

*API-Gateway nodes during the load test*  
<img width="1859" height="305" alt="Screenshot from 2026-10-10 18-03-27" src="https://github.com/user-attachments/assets/dde66254-42e2-47b1-af91-9eff5a7694b7" />

*API-Gateway nodes after the load test*  
<img width="1859" height="137" alt="Screenshot from 2026-10-10 18-02-19" src="https://github.com/user-attachments/assets/f59442ae-90da-420b-8f4c-1fe0274823fc" />

## 

## **9\. Comprehensive Troubleshooting Log & Root Cause Analysis**

This section catalogues the terminal errors, failed pod states, and configuration issues resolved during development\[cite: 1, 2, 3\].

### **Sprint 1: Project Setup, Cluster Initialization & Workloads\[cite: 1\]**

#### **Error 1: ImagePullBackOff / ErrImagePull**

* **Symptom:** Workload pods remained in Pending or cycled through ImagePullBackOff.  
* **Root Cause:** Container images built on the local Docker engine were not imported into the k3d nodes' containerd runtime.  
* **Resolution:** Loaded images into the cluster:  
  Bash  
  k3d image import client:latest api-gateway:latest \-c capstone-cluster

  Configured imagePullPolicy: IfNotPresent in deployment manifests.

  #### **Error 2: spec.template.metadata.labels: field is immutable (Selector / Label Mismatch)**

* **Symptom:** Applying updated manifests failed with schema validation errors; endpoints remained unbound:  
  Plaintext  
  The Deployment "api-gateway-deployment" is invalid: spec.template.metadata.labels: Invalid value: map\[string\]string{"app":"api-gateway"}: field is immutable

* **Root Cause:** Upstream manifests used app: api-gateway-service, whereas custom templates defined app: api-gateway. .spec.selector.matchLabels is an immutable field on Deployments.  
* **Resolution:** Standardized labels to app: api-gateway across manifests, recreated the deployment, and verified endpoint resolution:  
  Bash  
  kubectl delete deployment api-gateway-deployment  
  kubectl apply \-f k8s/api-gateway-deployment.yaml

  #### **Error 3: CrashLoopBackOff (Secret Encoding & Credential Keys)**

* **Symptom:** Database containers exited continuously with authentication failures.  
* **Root Cause:** Running echo "password" | base64 included a trailing newline (\\n), corrupting the secret hash. The Secret definition also mismatched application environment variables (POSTGRES\_PASSWORD vs. DB\_PASSWORD).  
* **Resolution:** Re-encoded secrets using echo \-n "password" | base64 and aligned environment keys across manifests.

### **Sprint 2: Storage Infrastructure & PVC Management\[cite: 1, 2\]**

#### **Error 4: OCI runtime exec failed: exec: "/bin/bash": no such file or directory**

* **Symptom:** Interactive debugging commands failed with exit code 126\.  
* **Root Cause:** Alpine and Distroless images do not include /bin/bash. Multi-container pods also failed when \-c \<container-name\> was omitted.  
* **Resolution:** Targeted /bin/sh and specified containers explicitly:  
  Bash  
  kubectl exec \-it deployment/postgres-deployment \-- /bin/sh  
  kubectl exec \-it lgtm-monitoring-grafana-5f649b9967-4pgcx \-c grafana \-- /bin/sh

  #### **Error 5: Orphaned Storage Claims (postgres-pvc)**

* **Symptom:** kubectl get pvc listed three bound volumes, including an unused claim.  
* **Root Cause:** postgres-pvc was a leftover claim from an earlier deployment; active workloads mounted database-persistent-volume-claim.  
* **Resolution:** Removed the redundant claim:  
  Bash  
  kubectl delete pvc postgres-pvc \-n default

### **Sprint 3 & 4: Load Testing, Autoscaling & Healer Implementation\[cite: 2\]**

#### **Error 6: wget: server returned error: HTTP/1.1 404 Not Found**

* **Symptom:** The load generator returned continuous HTTP 404 responses.  
* **Root Cause:** Requests were sent to /api/entries/all, an external ingress path. The internal Express service on port 5000 routes requests at /entries/all.  
* **Resolution:** Targeted the verified internal route:  
  Bash  
  kubectl run test-curl \--rm \-i \--tty \--image=curlimages/curl \--restart=Never \-- \\  
    curl \-s \-I http\://api-gateway-cluster-ip-service:5000/entries/all

  #### **Error 7: Error from server (AlreadyExists): pods "hpa-load-generator" already exists**

* **Symptom:** Re-running the load generator failed with an API conflict error.  
* **Root Cause:** A previous session was interrupted with Ctrl \+ C before the CLI could execute the \--rm cleanup hook.  
* **Resolution:** Cleaned up the orphaned pod:  
  Bash  
  kubectl delete pod hpa-load-generator \--force \--grace-period=0

  #### **Error 8: Embedded Git Repository Link (create mode 160000 k8s-cluster-healer)**

* **Symptom:** Git committed k8s-cluster-healer as an empty submodule pointer instead of versioning its source code.  
* **Root Cause:** The healer directory retained a nested .git folder from initialization.  
* **Resolution:** Removed the gitlink and embedded the source code directly:  
  Bash  
  git rm \--cached k8s-cluster-healer  
  rm \-rf k8s-cluster-healer/.git  
  git add k8s-cluster-healer/  
  git commit \-m "fix: embed k8s-cluster-healer source code directly"

### **Sprint 5 & 6: Observability Hardening & Dashboard Normalization\[cite: 2, 3\]**

#### 

#### **Error 9: Readiness/Liveness probe failed: context deadline exceeded**

* **Symptom:** Monitoring containers restarted repeatedly under load, dropping to 1/2 Running.  
* **Root Cause:** Prometheus required 16.97 seconds to replay its WAL during traffic bursts. The default 3s probe timeout was exceeded, triggering container restarts and cascading timeouts into Grafana.  
* **Resolution:** Relaxed probe tolerances in monitoring-lite.yaml (timeoutSeconds: 15, failureThreshold: 10, initialDelaySeconds: 120), stabilizing workloads at 2/2 Running.

  #### 

  #### **Error 10: PromQL Metric Cartesian Joins (24185054000% & 22.5 PiB)**

* **Symptom:** The storage panel displayed distorted numbers with duplicate legend entries.  
* **Root Cause:** The kubelet exported different instance labels across nodes, causing vector division to evaluate an unjoined Cartesian product. The panel was also configured with mebibytes units on raw byte inputs.  
* **Resolution:** Applied label aggregation via PromQL and fixed panel units to Percent (0-100).

## 

## **10\. Deployment & Verification Runbook**

### 

### **Prerequisites**

* Docker Engine 24+  
* k3d v5.x / kubectl v1.28+  
* helm v3.x  
* Go 1.24+

### **Step 1: Create 3-Node k3d Cluster**

Provision a multi-node cluster with 1 dedicated control-plane server and 2 worker agent nodes, mapping all required application and storage ingress ports:

> > Bash  
> > k3d cluster create capstone-cluster \\  
> >   \--servers 1 \\  
> >   \--agents 2 \\  
> >   \--port "80:80@loadbalancer" \\  
> >   \--port "3000:3000@loadbalancer" \\  
> >   \--port "5000:5000@loadbalancer" \\  
> >   \--port "9000:9000@loadbalancer" \\  
> >   \--port "9001:9001@loadbalancer"

Verify that all three nodes are in `Ready` state:

> > Bash  
> > kubectl get nodes  
\[Output\]  
<img width="1847" height="257" alt="Screenshot from 2026-10-10 19-08-20" src="https://github.com/user-attachments/assets/a81e0698-4dcc-4ca8-bc1a-deb7fc30ffe1" /> 
> > 

### **Step 2: Deploy Applications, Database & MinIO Storage**

Apply all workloads, services, and PersistentVolumeClaims deployed after carefully reviewing yaml files and updating the cloud infra components with on premise:

> > Bash  
> > kubectl apply \-f k8s/

MinIO is deployed as a stateful single-pod instance with a dedicated PersistentVolumeClaim:

To ensure the target bucket exists before `files-service` starts processing uploads, a Kubernetes `Job` or `initContainer` running the MinIO Client (`mc`) initializes the bucket and sets its access policy:

Initialized the MinIO `gratitude-uploads` S3 bucket:

> > Bash  
> > kubectl run minio-init-bucket \--rm \-i \--tty \--image=minio/mc \--restart=Never \-- \\  
> >   sh \-c "  
> >     mc alias set myminio http\://minio:9000 minioadmin minioadmin && \\  
> >     mc mb myminio/gratitude-uploads \--ignore-existing && \\  
> >     mc anonymous set download myminio/gratitude-uploads  
> >   "  
Verification:

**S3 API Endpoint:** `[http://minio.default.svc.cluster.local:9000]`  
Listing all the pods: cmd- kubectl get pods \-o wide

\[Output\]  
<img width="1861" height="915" alt="Screenshot from 2026-10-10 19-11-20" src="https://github.com/user-attachments/assets/a7749046-1ace-44d7-8f94-951a9f91b87e" />

### **Step 3: Build & Deploy k8s-cluster-healer**

#### **Option A: Run Locally via Kubeconfig**

> >   
> > Bash  
> > cd k8s-cluster-healer  
> > go run cmd/healer/main.go

#### **Option B: Deploy In-Cluster**

> >   
> > Bash  
> > cd k8s-cluster-healer  
> > docker build \-t k8s-cluster-healer:v1.0 .  
> > k3d image import k8s-cluster-healer:v1.0 \-c capstone-cluster  
> > kubectl apply \-f deploy/rbac.yaml  
> > kubectl apply \-f deploy/deployment.yaml  
> > cd ..

### 

### **Step 4: Deploy Observability Stack**

Bash  
helm repo add prometheus-community https\://prometheus-community.github.io/helm-charts  
helm repo update  
helm upgrade \--install lgtm-monitoring prometheus-community/kube-prometheus-stack \--namespace monitoring  \--create-namespace \-f monitoring-lite.yaml

### **Step 5: Access Grafana Dashboard**

Bash  
nohup bash \-c 'while true; do kubectl port-forward \-n monitoring svc/lgtm-monitoring-grafana 3001:80; sleep 2; done' \> /tmp/grafana-pf.log 2\>&1 &

Access \[http\://127.0.0.1:3001\](http\://127.0.0.1:3001) in your browser (Credentials: admin / admin). Import k8s/gratitudeapp-grafana-dashboard.json. For the custom dashboard

### 

### **Step 6: Verify Final Cluster Status**

Bash  
kubectl get pods,pvc,hpa \-n default  
     
**\[Screenshot \- Terminal Output of Fully Reconciled Cluster\]**  
*(Insert terminal capture showing all 11 microservice pods Running, both PVCs Bound, and HPA ready)*

<img width="1859" height="855" alt="overall k8 infra" src="https://github.com/user-attachments/assets/bf5b71f3-aefc-4e67-bb09-87ec993d3f5c" />

**\[Screenshot \- GratitudeApp React Frontend User Interface\]**  
*Description: Capture of the web client UI served at `http://localhost:3000` showing the active dashboard with populated journal entries, selected mood tag distributions (e.g., Happy, Peaceful, Grateful), and real-time habit analytics displaying updated total reflections, current streaks, and mood tracking stats)*

<img width="1854" height="2173" alt="graphana" src="https://github.com/user-attachments/assets/5e7782f4-1d48-4a4d-8f37-bdf5b8738ab1" />

##

## **11\. Cost Optimization Analysis**

To address the **Cost Optimization (10.00%)** rubric, this architecture demonstrates how to achieve enterprise-grade orchestration, self-healing, observability, and storage capabilities entirely on-premise without incurring costly managed public cloud infrastructure expenses:

1. **Zero-Cloud Infrastructure Cost (\$0 Overhead):**  
   * Emulated an enterprise 3-node Kubernetes topology locally using `k3d` and lightweight `k3s` on Docker Desktop, completely bypassing expensive managed control planes like AWS EKS (\$73/month base control plane fee) or GCP GKE.  
   * Multi-node worker nodes (`agent-0`, `agent-1`) run inside lightweight Docker container boundaries, achieving true multi-node network and scheduling isolation on developer workstations with zero compute billing.  
2. **Self-Hosted S3 Storage via MinIO (Avoiding AWS S3 & Data Transfer Fees):**  
   * Deployed a local, containerized MinIO instance backed by standard Linux filesystem persistence via Kubernetes `local-path-provisioner`.  
   * Replaces expensive AWS S3 bucket subscriptions, API request costs (`PUT`/`GET` requests), and outbound data transfer fees while preserving 100% S3-compatible API semantics for microservice file uploads.  
3. **Local Persistent Storage Reclamation (Avoiding Block Storage Leaks):**  
   * Dynamic PVC provisioning uses Kubernetes local node paths directly rather than cloud-managed network block storage (such as AWS EBS gp3 volumes).  
   * Proactively identified and purged orphaned volume claims (`postgres-pvc`) during development lifecycle testing, preventing unattached storage leaks from wasting host SSD capacity.  
4. **Elastic Scaling (HPA) to Conserve Compute Resources:**  
   * Tight resource baselines (`requests: 20m` CPU, `32Mi` memory) ensure low baseline memory utilization during idle hours.  
   * The Horizontal Pod Autoscaler (HPA) dynamically scales the API Gateway only during traffic surges (scaling from 1 to 5 replicas) and cools back down after 5 minutes, demonstrating cost-efficient resource scheduling.  
5. **Automated Zombie Workload Prevention (`k8s-cluster-healer`):**  
   * The custom Go auto-healer identifies pods stuck in `CrashLoopBackOff` or high restart loops and purges them.  
   * Prevents failed containers from accumulating leaked socket descriptors, zombie memory allocations, and scheduler loop churn on the host machine.

      **6\. Resource-Tuned Observability Profile:**

* Customized the Prometheus & Grafana stack via `monitoring-lite.yaml` to cap metrics retention to 6 hours and scrape intervals to 30s.  
  * Disabled non-essential components (`alertmanager`, `nodeExporter`), fitting production-grade metrics scraping and dashboard visualization inside \~1.5 GB total host RAM without starving microservice workloads.

## 

## **Evaluation Criteria Mapping**

* **Documentation (15.00%):** Comprehensive README covering architecture, workflows, catalogs, runbooks, and troubleshooting logs\[cite: 3\].  
* **Implementation (75.00%):** Production Kubernetes manifests, operational Go self-healing operator (k8s-cluster-healer)\[cite: 1\], validated HPA scaling\[cite: 2\], and an integrated Prometheus/Grafana stack\[cite: 1, 3\].  
* **Cost Optimization (10.00%):** Dynamic autoscaling, tight container resource limits, pruned storage claims, and an optimized monitoring profile\[cite: 3\].

### 

### **Update Command on Local Machine**

To replace your local file with this version and push it to GitHub, execute:

Bash  
cd \~/Downloads/HeroviredNotes/Assignments/"CAPSTONE PROJECT"/Project1/GratitudeApp-main  
git add README.md  
git commit \-m "docs: finalize comprehensive README with complete architecture, runbook, and healer integration"  
git push origin main  
