package docs

import (
	"fmt"
	"strings"
)

// PackageKnowledge holds practical, task-oriented documentation for a well-known package
type PackageKnowledge struct {
	// Short description of what this package is for
	Purpose string
	// Real-world tasks you can do with this package (task → how)
	Tasks []PackageTask
	// Package-level "When To Use" points
	WhenToUse []string
	// Package-level "What It Does" points
	WhatItDoes []string
	// Struct-specific knowledge: struct name → knowledge
	Structs map[string]StructKnowledge
	// Function-specific knowledge: func name → knowledge
	Functions map[string]FuncKnowledge
}

// PackageTask is a practical task users can accomplish with the package
type PackageTask struct {
	Task    string // e.g., "Authenticate to a Kubernetes cluster"
	How     string // e.g., "Use rest.InClusterConfig() or clientcmd.BuildConfigFromFlags()"
	Example string // short code hint
}

// StructKnowledge has practical knowledge about a specific struct
type StructKnowledge struct {
	WhatItDoes string
	WhenToUse  string
	HowToUse   string // code hint
}

// FuncKnowledge has practical knowledge about a specific function
type FuncKnowledge struct {
	WhatItDoes string
	WhenToUse  string
	HowToUse   string
}

// knowledgeBase maps import path patterns → package knowledge
// Patterns are checked with strings.Contains so "client-go/kubernetes" matches "k8s.io/client-go/kubernetes"
var knowledgeBase = map[string]*PackageKnowledge{

	// ─── k8s.io/client-go ───
	"client-go/kubernetes": {
		Purpose: "Official Go client for talking to the Kubernetes API — list pods, create deployments, watch events, everything.",
		Tasks: []PackageTask{
			{Task: "Authenticate to Kubernetes cluster", How: "Use <code>rest.InClusterConfig()</code> (inside a pod) or <code>clientcmd.BuildConfigFromFlags()</code> (from your laptop)", Example: "config, err := rest.InClusterConfig()"},
			{Task: "Create a clientset to talk to the API", How: "Pass the config to <code>kubernetes.NewForConfig(config)</code>", Example: "clientset, err := kubernetes.NewForConfig(config)"},
			{Task: "List all pods in a namespace", How: "Use <code>clientset.CoreV1().Pods(\"namespace\").List(ctx, metav1.ListOptions{})</code>", Example: "pods, err := clientset.CoreV1().Pods(\"default\").List(ctx, metav1.ListOptions{})"},
			{Task: "Create a deployment", How: "Use <code>clientset.AppsV1().Deployments(ns).Create(ctx, deployment, metav1.CreateOptions{})</code>", Example: "result, err := clientset.AppsV1().Deployments(\"default\").Create(ctx, myDeploy, metav1.CreateOptions{})"},
			{Task: "Watch for changes to resources", How: "Use <code>.Watch(ctx, metav1.ListOptions{})</code> on any resource and range over the channel", Example: "watcher, _ := clientset.CoreV1().Pods(ns).Watch(ctx, metav1.ListOptions{})"},
			{Task: "Delete a pod", How: "Use <code>clientset.CoreV1().Pods(ns).Delete(ctx, name, metav1.DeleteOptions{})</code>", Example: "err := clientset.CoreV1().Pods(\"default\").Delete(ctx, \"my-pod\", metav1.DeleteOptions{})"},
		},
		Structs: map[string]StructKnowledge{
			"Clientset": {
				WhatItDoes: "The main entry point for <strong>all Kubernetes API calls</strong>. It holds authenticated connections to every API group (CoreV1, AppsV1, BatchV1, etc.).",
				WhenToUse:  "Create ONE <code>Clientset</code> at startup, then use it everywhere to list pods, create deployments, delete services, etc.",
				HowToUse:   "clientset, err := kubernetes.NewForConfig(config)\npods, _ := clientset.CoreV1().Pods(\"default\").List(ctx, metav1.ListOptions{})",
			},
		},
	},
	"client-go/rest": {
		Purpose: "Handles authentication and HTTP connection to the Kubernetes API server.",
		Tasks: []PackageTask{
			{Task: "Get config when running INSIDE a Kubernetes pod", How: "Use <code>rest.InClusterConfig()</code> — it reads the service account token automatically", Example: "config, err := rest.InClusterConfig()"},
			{Task: "Get config when running OUTSIDE the cluster (your laptop)", How: "Use <code>clientcmd.BuildConfigFromFlags(\"\", kubeconfig)</code>", Example: "config, err := clientcmd.BuildConfigFromFlags(\"\", \"/home/user/.kube/config\")"},
			{Task: "Make raw HTTP requests to the API", How: "Use <code>rest.RESTClientFor(config)</code> for low-level HTTP access", Example: "client, err := rest.RESTClientFor(config)"},
		},
		Structs: map[string]StructKnowledge{
			"Config": {
				WhatItDoes: "Holds <strong>all connection details</strong> for the Kubernetes API — server URL, auth token, TLS certs, timeout, etc.",
				WhenToUse:  "You need this to create a <code>Clientset</code>. Get it from <code>InClusterConfig()</code> or <code>BuildConfigFromFlags()</code>.",
				HowToUse:   "config, err := rest.InClusterConfig()\n// or from kubeconfig:\nconfig, err := clientcmd.BuildConfigFromFlags(\"\", kubeconfigPath)",
			},
		},
	},
	"client-go/tools/clientcmd": {
		Purpose: "Reads and parses kubeconfig files (~/.kube/config) to connect to Kubernetes from outside the cluster.",
		Tasks: []PackageTask{
			{Task: "Load kubeconfig from default path", How: "Use <code>clientcmd.BuildConfigFromFlags(\"\", kubeconfigPath)</code>", Example: "config, _ := clientcmd.BuildConfigFromFlags(\"\", filepath.Join(homedir.HomeDir(), \".kube\", \"config\"))"},
			{Task: "Switch between multiple clusters/contexts", How: "Use <code>clientcmd.NewNonInteractiveDeferredLoadingClientConfig()</code> with overrides", Example: "// set context override to switch clusters"},
		},
	},
	"client-go/informers": {
		Purpose: "Watches Kubernetes resources efficiently using a local cache — instead of hitting the API every time.",
		Tasks: []PackageTask{
			{Task: "Watch pods efficiently with a local cache", How: "Create a SharedInformerFactory, get a pod informer, add event handlers", Example: "factory := informers.NewSharedInformerFactory(clientset, 30*time.Second)\npodInformer := factory.Core().V1().Pods()"},
			{Task: "React when a pod is created/updated/deleted", How: "Add event handlers: <code>AddFunc</code>, <code>UpdateFunc</code>, <code>DeleteFunc</code>", Example: "podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{\n  AddFunc: func(obj interface{}) { ... },\n})"},
		},
	},

	// ─── k8s.io/apimachinery ───
	"apimachinery/pkg/apis/meta/v1": {
		Purpose: "Common types used in EVERY Kubernetes API call — like ListOptions, CreateOptions, ObjectMeta, labels, etc.",
		Tasks: []PackageTask{
			{Task: "Filter resources with labels", How: "Set <code>LabelSelector</code> in <code>ListOptions</code>", Example: "opts := metav1.ListOptions{LabelSelector: \"app=nginx\"}"},
			{Task: "Limit how many results to return", How: "Set <code>Limit</code> in <code>ListOptions</code>", Example: "opts := metav1.ListOptions{Limit: 100}"},
			{Task: "Set resource metadata (name, namespace, labels)", How: "Fill in <code>ObjectMeta</code> when creating resources", Example: "meta := metav1.ObjectMeta{Name: \"my-pod\", Namespace: \"default\", Labels: map[string]string{\"app\": \"web\"}}"},
		},
		Structs: map[string]StructKnowledge{
			"ListOptions": {
				WhatItDoes: "Controls <strong>what and how many resources</strong> to fetch — filters by label, field, limit, and pagination.",
				WhenToUse:  "Pass it to any <code>.List()</code> or <code>.Watch()</code> call to filter results.",
				HowToUse:   "pods, err := clientset.CoreV1().Pods(\"default\").List(ctx, metav1.ListOptions{\n  LabelSelector: \"app=nginx\",\n  Limit: 10,\n})",
			},
			"ObjectMeta": {
				WhatItDoes: "The <strong>identity card</strong> of every Kubernetes resource — holds name, namespace, labels, annotations, UID, timestamps.",
				WhenToUse:  "Fill this in when <strong>creating</strong> any Kubernetes resource (Pod, Deployment, Service, etc.).",
				HowToUse:   "pod := &corev1.Pod{\n  ObjectMeta: metav1.ObjectMeta{\n    Name: \"my-pod\",\n    Namespace: \"default\",\n    Labels: map[string]string{\"app\": \"web\"},\n  },\n}",
			},
		},
	},

	// ─── k8s.io/api ───
	"k8s.io/api/core/v1": {
		Purpose: "All Kubernetes core resource types — Pod, Service, Node, Namespace, ConfigMap, Secret, PV, PVC, etc.",
		Tasks: []PackageTask{
			{Task: "Define a Pod spec", How: "Create a <code>corev1.Pod</code> with <code>ObjectMeta</code> and <code>PodSpec</code>", Example: "pod := &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: \"web\"}, Spec: corev1.PodSpec{Containers: []corev1.Container{{Name: \"nginx\", Image: \"nginx:latest\"}}}}"},
			{Task: "Define a container", How: "Fill in <code>corev1.Container</code> with name, image, ports, env vars", Example: "c := corev1.Container{Name: \"app\", Image: \"myapp:v1\", Ports: []corev1.ContainerPort{{ContainerPort: 8080}}}"},
			{Task: "Create a ConfigMap in code", How: "Build a <code>corev1.ConfigMap</code> with data", Example: "cm := &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: \"my-config\"}, Data: map[string]string{\"key\": \"value\"}}"},
		},
	},
	"k8s.io/api/apps/v1": {
		Purpose: "Kubernetes workload types — Deployment, StatefulSet, DaemonSet, ReplicaSet.",
		Tasks: []PackageTask{
			{Task: "Define a Deployment", How: "Create <code>appsv1.Deployment</code> with template and replicas", Example: "deploy := &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Name: \"web\"}, Spec: appsv1.DeploymentSpec{Replicas: int32Ptr(3), ...}}"},
			{Task: "Scale a deployment (change replicas)", How: "Get the deployment, change <code>.Spec.Replicas</code>, then Update", Example: "deploy.Spec.Replicas = int32Ptr(5)\nclientset.AppsV1().Deployments(ns).Update(ctx, deploy, metav1.UpdateOptions{})"},
		},
	},

	// ─── controller-runtime ───
	"sigs.k8s.io/controller-runtime": {
		Purpose: "Framework for building Kubernetes operators and controllers — handles reconciliation loops, caching, leader election.",
		Tasks: []PackageTask{
			{Task: "Create a new controller/operator", How: "Use <code>ctrl.NewManager()</code> + <code>ctrl.NewControllerManagedBy(mgr)</code>", Example: "mgr, _ := ctrl.NewManager(cfg, ctrl.Options{})\nctrl.NewControllerManagedBy(mgr).For(&myv1.MyResource{}).Complete(reconciler)"},
			{Task: "Reconcile a custom resource", How: "Implement the <code>Reconciler</code> interface with a <code>Reconcile(ctx, req)</code> method", Example: "func (r *MyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) { ... }"},
			{Task: "Read a resource from the cluster", How: "Use <code>client.Get(ctx, key, obj)</code>", Example: "err := r.Client.Get(ctx, req.NamespacedName, &myPod)"},
		},
	},

	// ─── gin ───
	"gin-gonic/gin": {
		Purpose: "Fast HTTP web framework — define routes, handle requests, return JSON, serve APIs.",
		Tasks: []PackageTask{
			{Task: "Start a web server", How: "Create a router with <code>gin.Default()</code> then <code>r.Run()</code>", Example: "r := gin.Default()\nr.Run(\":8080\")"},
			{Task: "Handle GET /hello", How: "Use <code>r.GET(\"/hello\", handlerFunc)</code>", Example: "r.GET(\"/hello\", func(c *gin.Context) {\n  c.JSON(200, gin.H{\"message\": \"hello\"})\n})"},
			{Task: "Handle POST with JSON body", How: "Use <code>c.ShouldBindJSON(&myStruct)</code> to parse the body", Example: "r.POST(\"/users\", func(c *gin.Context) {\n  var user User\n  c.ShouldBindJSON(&user)\n  c.JSON(201, user)\n})"},
			{Task: "Get URL parameters", How: "Use <code>c.Param(\"id\")</code> for path params, <code>c.Query(\"q\")</code> for query strings", Example: "r.GET(\"/users/:id\", func(c *gin.Context) {\n  id := c.Param(\"id\")\n})"},
			{Task: "Add middleware (logging, auth)", How: "Use <code>r.Use(middleware)</code>", Example: "r.Use(gin.Logger(), gin.Recovery())"},
		},
	},

	// ─── cobra ───
	"spf13/cobra": {
		Purpose: "Build beautiful CLI apps with commands, flags, and help text — like kubectl, docker, gh.",
		Tasks: []PackageTask{
			{Task: "Create a root command", How: "Use <code>cobra.Command{}</code> with Use, Short, Long, Run fields", Example: "rootCmd := &cobra.Command{Use: \"myapp\", Short: \"My CLI tool\", Run: func(cmd *cobra.Command, args []string) { ... }}"},
			{Task: "Add a sub-command (like 'myapp get')", How: "Create another <code>cobra.Command</code> and use <code>rootCmd.AddCommand(subCmd)</code>", Example: "getCmd := &cobra.Command{Use: \"get\", Run: ...}\nrootCmd.AddCommand(getCmd)"},
			{Task: "Add a flag (--verbose, --output)", How: "Use <code>cmd.Flags().StringVar()</code> or <code>cmd.Flags().BoolP()</code>", Example: "rootCmd.Flags().BoolP(\"verbose\", \"v\", false, \"verbose output\")"},
			{Task: "Run the CLI", How: "Call <code>rootCmd.Execute()</code> in main()", Example: "func main() { rootCmd.Execute() }"},
		},
	},

	// ─── viper ───
	"spf13/viper": {
		Purpose: "Read configuration from files (YAML, JSON, TOML), environment variables, and command-line flags — all in one place.",
		Tasks: []PackageTask{
			{Task: "Read a YAML config file", How: "Set file name and path, then call <code>viper.ReadInConfig()</code>", Example: "viper.SetConfigName(\"config\")\nviper.AddConfigPath(\".\")\nviper.ReadInConfig()"},
			{Task: "Get a config value", How: "Use <code>viper.GetString(\"key\")</code>, <code>viper.GetInt(\"key\")</code>, etc.", Example: "port := viper.GetInt(\"server.port\")\nname := viper.GetString(\"app.name\")"},
			{Task: "Read from environment variables", How: "Use <code>viper.AutomaticEnv()</code> — it maps MY_VAR to config key my_var", Example: "viper.AutomaticEnv()\ndbHost := viper.GetString(\"DB_HOST\")"},
			{Task: "Set default values", How: "Use <code>viper.SetDefault(\"key\", value)</code>", Example: "viper.SetDefault(\"port\", 8080)"},
		},
	},

	// ─── logrus ───
	"sirupsen/logrus": {
		Purpose: "Structured logging — log messages with fields (key-value pairs) for easy searching and filtering.",
		Tasks: []PackageTask{
			{Task: "Log a simple message", How: "Use <code>logrus.Info(\"message\")</code>, <code>logrus.Error(\"message\")</code>", Example: "logrus.Info(\"Server started\")"},
			{Task: "Log with extra fields", How: "Use <code>logrus.WithFields(logrus.Fields{...}).Info(\"msg\")</code>", Example: "logrus.WithFields(logrus.Fields{\"user\": \"admin\", \"action\": \"login\"}).Info(\"User logged in\")"},
			{Task: "Set log level", How: "Use <code>logrus.SetLevel(logrus.DebugLevel)</code>", Example: "logrus.SetLevel(logrus.DebugLevel)"},
			{Task: "Output as JSON", How: "Use <code>logrus.SetFormatter(&logrus.JSONFormatter{})</code>", Example: "logrus.SetFormatter(&logrus.JSONFormatter{})"},
		},
	},

	// ─── zap ───
	"go.uber.org/zap": {
		Purpose: "Blazing-fast structured logging from Uber — 10x faster than logrus for high-throughput applications.",
		Tasks: []PackageTask{
			{Task: "Create a logger", How: "Use <code>zap.NewProduction()</code> or <code>zap.NewDevelopment()</code>", Example: "logger, _ := zap.NewProduction()\ndefer logger.Sync()"},
			{Task: "Log a message with fields", How: "Use <code>logger.Info(\"msg\", zap.String(\"key\", \"val\"))</code>", Example: "logger.Info(\"User created\", zap.String(\"name\", \"admin\"), zap.Int(\"age\", 30))"},
			{Task: "Use the sugared (easy) API", How: "Use <code>logger.Sugar()</code> for printf-style logging", Example: "sugar := logger.Sugar()\nsugar.Infof(\"Hello %s\", name)"},
		},
	},

	// ─── testify ───
	"stretchr/testify": {
		Purpose: "Testing toolkit — assertions, mocks, and test suites that make Go tests readable and easy to write.",
		Tasks: []PackageTask{
			{Task: "Assert two values are equal", How: "Use <code>assert.Equal(t, expected, actual)</code>", Example: "assert.Equal(t, 200, resp.StatusCode)"},
			{Task: "Assert no error occurred", How: "Use <code>assert.NoError(t, err)</code>", Example: "result, err := myFunc()\nassert.NoError(t, err)"},
			{Task: "Assert a value is not nil", How: "Use <code>assert.NotNil(t, obj)</code>", Example: "assert.NotNil(t, user)"},
			{Task: "Require (fail immediately if wrong)", How: "Use <code>require.Equal(t, expected, actual)</code> — stops the test on failure", Example: "require.NoError(t, err) // test stops here if err != nil"},
		},
	},

	// ─── gorilla/mux ───
	"gorilla/mux": {
		Purpose: "Powerful HTTP router — pattern matching, path variables, middleware, subrouters.",
		Tasks: []PackageTask{
			{Task: "Create a router", How: "Use <code>mux.NewRouter()</code>", Example: "r := mux.NewRouter()"},
			{Task: "Handle a route with path variables", How: "Use <code>r.HandleFunc(\"/users/{id}\", handler)</code>", Example: "r.HandleFunc(\"/users/{id}\", func(w http.ResponseWriter, r *http.Request) {\n  vars := mux.Vars(r)\n  id := vars[\"id\"]\n})"},
			{Task: "Start the server", How: "Use <code>http.ListenAndServe(\":8080\", r)</code>", Example: "http.ListenAndServe(\":8080\", r)"},
		},
	},

	// ─── chi ───
	"go-chi/chi": {
		Purpose: "Lightweight, composable HTTP router — middleware-friendly, great for REST APIs.",
		Tasks: []PackageTask{
			{Task: "Create a router", How: "Use <code>chi.NewRouter()</code>", Example: "r := chi.NewRouter()"},
			{Task: "Add middleware", How: "Use <code>r.Use(middleware.Logger)</code>", Example: "r.Use(middleware.Logger, middleware.Recoverer)"},
			{Task: "Handle a route", How: "Use <code>r.Get(\"/path\", handler)</code>", Example: "r.Get(\"/hello\", func(w http.ResponseWriter, r *http.Request) {\n  w.Write([]byte(\"hello\"))\n})"},
		},
	},

	// ─── gRPC ───
	"google.golang.org/grpc": {
		Purpose: "Build high-performance RPC services — define services in .proto files, generate Go code, call remote functions like local ones.",
		Tasks: []PackageTask{
			{Task: "Start a gRPC server", How: "Create a listener, register your service, call <code>grpc.NewServer()</code>", Example: "lis, _ := net.Listen(\"tcp\", \":50051\")\ns := grpc.NewServer()\npb.RegisterMyServiceServer(s, &server{})\ns.Serve(lis)"},
			{Task: "Connect to a gRPC server (client)", How: "Use <code>grpc.Dial(addr, opts)</code>", Example: "conn, _ := grpc.Dial(\"localhost:50051\", grpc.WithInsecure())\nclient := pb.NewMyServiceClient(conn)"},
		},
	},

	// ─── Standard library enrichments ───
	"net/http": {
		Purpose: "Go's built-in HTTP toolkit — build web servers, make HTTP requests, handle routes.",
		Tasks: []PackageTask{
			{Task: "Start a web server", How: "Use <code>http.ListenAndServe(\":8080\", nil)</code>", Example: "http.HandleFunc(\"/\", handler)\nhttp.ListenAndServe(\":8080\", nil)"},
			{Task: "Handle a URL path", How: "Use <code>http.HandleFunc(\"/path\", handler)</code>", Example: "http.HandleFunc(\"/hello\", func(w http.ResponseWriter, r *http.Request) {\n  fmt.Fprintf(w, \"Hello!\")\n})"},
			{Task: "Make a GET request", How: "Use <code>http.Get(url)</code>", Example: "resp, err := http.Get(\"https://api.example.com/data\")\ndefer resp.Body.Close()"},
			{Task: "Make a POST request with JSON", How: "Use <code>http.Post(url, contentType, body)</code>", Example: "resp, _ := http.Post(\"https://api.example.com\", \"application/json\", bytes.NewBuffer(jsonData))"},
		},
	},
	"fmt": {
		Purpose: "Print stuff to the screen, format strings, read user input — the most basic Go package.",
		Tasks: []PackageTask{
			{Task: "Print to the screen", How: "Use <code>fmt.Println(\"hello\")</code>", Example: "fmt.Println(\"Hello, World!\")"},
			{Task: "Format a string without printing", How: "Use <code>fmt.Sprintf(\"Hello %s\", name)</code>", Example: "msg := fmt.Sprintf(\"Hello %s, you are %d\", name, age)"},
			{Task: "Print with formatting", How: "Use <code>fmt.Printf(\"format\", values)</code>", Example: "fmt.Printf(\"Name: %s, Age: %d\\n\", name, age)"},
			{Task: "Read user input", How: "Use <code>fmt.Scan(&variable)</code>", Example: "var name string\nfmt.Print(\"Enter name: \")\nfmt.Scan(&name)"},
		},
	},
	"encoding/json": {
		Purpose: "Convert between Go structs and JSON — encode (struct → JSON) and decode (JSON → struct).",
		Tasks: []PackageTask{
			{Task: "Convert a struct to JSON", How: "Use <code>json.Marshal(myStruct)</code>", Example: "data, err := json.Marshal(user)\nfmt.Println(string(data))"},
			{Task: "Convert JSON to a struct", How: "Use <code>json.Unmarshal(data, &myStruct)</code>", Example: "var user User\njson.Unmarshal(jsonBytes, &user)"},
			{Task: "Pretty-print JSON", How: "Use <code>json.MarshalIndent(obj, \"\", \"  \")</code>", Example: "pretty, _ := json.MarshalIndent(user, \"\", \"  \")"},
			{Task: "Read JSON from an HTTP response", How: "Use <code>json.NewDecoder(resp.Body).Decode(&obj)</code>", Example: "var result ApiResponse\njson.NewDecoder(resp.Body).Decode(&result)"},
		},
	},
	"os": {
		Purpose: "Interact with the operating system — read files, environment variables, command-line args, exit the program.",
		Tasks: []PackageTask{
			{Task: "Read an environment variable", How: "Use <code>os.Getenv(\"KEY\")</code>", Example: "home := os.Getenv(\"HOME\")"},
			{Task: "Read a file", How: "Use <code>os.ReadFile(\"path\")</code>", Example: "data, err := os.ReadFile(\"config.yaml\")"},
			{Task: "Create/write a file", How: "Use <code>os.WriteFile(\"path\", data, 0644)</code>", Example: "os.WriteFile(\"output.txt\", []byte(\"hello\"), 0644)"},
			{Task: "Get command-line arguments", How: "Use <code>os.Args</code>", Example: "args := os.Args[1:] // skip program name"},
			{Task: "Exit the program", How: "Use <code>os.Exit(code)</code>", Example: "os.Exit(1) // exit with error"},
		},
	},
	"io": {
		Purpose: "Basic I/O interfaces — Reader, Writer, Closer. Everything in Go that reads or writes uses these.",
		Tasks: []PackageTask{
			{Task: "Read all bytes from a reader", How: "Use <code>io.ReadAll(reader)</code>", Example: "body, err := io.ReadAll(resp.Body)"},
			{Task: "Copy data from reader to writer", How: "Use <code>io.Copy(dst, src)</code>", Example: "io.Copy(os.Stdout, resp.Body) // print response to screen"},
			{Task: "Create a reader from a string", How: "Use <code>strings.NewReader(\"text\")</code>", Example: "reader := strings.NewReader(\"hello world\")"},
		},
	},
	"context": {
		Purpose: "Pass deadlines, cancellation signals, and request-scoped values through your program — stops work that's no longer needed.",
		Tasks: []PackageTask{
			{Task: "Create a context with timeout", How: "Use <code>context.WithTimeout(parent, duration)</code>", Example: "ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)\ndefer cancel()"},
			{Task: "Create a cancellable context", How: "Use <code>context.WithCancel(parent)</code>", Example: "ctx, cancel := context.WithCancel(context.Background())\n// call cancel() when done"},
			{Task: "Pass context to a function", How: "Add <code>ctx context.Context</code> as the first parameter", Example: "func myFunc(ctx context.Context) error { ... }"},
		},
	},
	"sync": {
		Purpose: "Synchronization tools for concurrent Go programs — mutexes, wait groups, once.",
		Tasks: []PackageTask{
			{Task: "Wait for multiple goroutines to finish", How: "Use <code>sync.WaitGroup</code>", Example: "var wg sync.WaitGroup\nwg.Add(1)\ngo func() { defer wg.Done(); doWork() }()\nwg.Wait()"},
			{Task: "Protect shared data from race conditions", How: "Use <code>sync.Mutex</code>", Example: "var mu sync.Mutex\nmu.Lock()\n// access shared data\nmu.Unlock()"},
			{Task: "Run something exactly once", How: "Use <code>sync.Once</code>", Example: "var once sync.Once\nonce.Do(func() { initializeDB() })"},
		},
	},

	// ─── Additional client-go packages ───
	"client-go/tools/clientcmd/api": {
		Purpose: "Data types that represent the kubeconfig file structure — Cluster, AuthInfo, Context, Config objects.",
		Tasks: []PackageTask{
			{Task: "Understand kubeconfig structure", How: "The <code>Config</code> struct mirrors your <code>~/.kube/config</code> file: clusters, users, contexts", Example: "config := api.Config{\n  Clusters: map[string]*api.Cluster{\"prod\": {Server: \"https://...\"}},\n}"},
			{Task: "Build kubeconfig programmatically", How: "Create <code>api.Config</code> with clusters, users, and contexts, then write with <code>clientcmd.WriteToFile()</code>", Example: "cfg := api.NewConfig()\ncfg.Clusters[\"my-cluster\"] = &api.Cluster{Server: url}"},
		},
	},
	"client-go/dynamic": {
		Purpose: "Work with ANY Kubernetes resource without needing typed Go structs — useful for CRDs and unknown resource types.",
		Tasks: []PackageTask{
			{Task: "Get a CRD or custom resource", How: "Use <code>dynamic.NewForConfig(config)</code> then <code>.Resource(gvr).Get()</code>", Example: "client := dynamic.NewForConfig(config)\nobj, _ := client.Resource(gvr).Namespace(\"default\").Get(ctx, \"my-cr\", metav1.GetOptions{})"},
			{Task: "List resources by GVR (Group/Version/Resource)", How: "Pass a <code>schema.GroupVersionResource</code> to <code>.Resource()</code>", Example: "gvr := schema.GroupVersionResource{Group: \"apps\", Version: \"v1\", Resource: \"deployments\"}\nlist, _ := client.Resource(gvr).List(ctx, metav1.ListOptions{})"},
		},
	},
	"client-go/tools/cache": {
		Purpose: "Local in-memory cache for Kubernetes resources — stores objects from informers so you don't hit the API every time.",
		Tasks: []PackageTask{
			{Task: "Get a cached object by key", How: "Use <code>cache.MetaNamespaceKeyFunc(obj)</code> to get the key, then <code>store.GetByKey(key)</code>", Example: "key, _ := cache.MetaNamespaceKeyFunc(pod)\nobj, exists, _ := indexer.GetByKey(key)"},
			{Task: "Add event handlers for add/update/delete", How: "Use <code>cache.ResourceEventHandlerFuncs{AddFunc, UpdateFunc, DeleteFunc}</code>", Example: "informer.AddEventHandler(cache.ResourceEventHandlerFuncs{\n  AddFunc: func(obj interface{}) { fmt.Println(\"added\") },\n})"},
			{Task: "Wait for cache to sync before processing", How: "Use <code>cache.WaitForCacheSync(stopCh, informer.HasSynced)</code>", Example: "if !cache.WaitForCacheSync(stopCh, podInformer.HasSynced) {\n  log.Fatal(\"cache never synced\")\n}"},
		},
	},
	"client-go/tools/record": {
		Purpose: "Create Kubernetes Events — the messages you see when you run <code>kubectl describe pod</code> (like 'Pulled image', 'Started container').",
		Tasks: []PackageTask{
			{Task: "Record an event on a Kubernetes resource", How: "Create an <code>EventRecorder</code> then call <code>.Event(obj, type, reason, message)</code>", Example: "recorder.Event(pod, corev1.EventTypeNormal, \"Synced\", \"Pod synced successfully\")"},
		},
	},
	"client-go/tools/leaderelection": {
		Purpose: "Ensures only ONE replica of your controller is active at a time — critical for HA (high-availability) operators.",
		Tasks: []PackageTask{
			{Task: "Set up leader election for your controller", How: "Create a <code>LeaderElector</code> with a lock (ConfigMap or Lease) and callbacks", Example: "le, _ := leaderelection.NewLeaderElector(leaderelection.LeaderElectionConfig{\n  Lock: lock, LeaseDuration: 15*time.Second,\n  Callbacks: leaderelection.LeaderCallbacks{OnStartedLeading: run},\n})"},
		},
	},
	"client-go/util/workqueue": {
		Purpose: "Rate-limited work queue for processing Kubernetes events — ensures you don't overwhelm the API server.",
		Tasks: []PackageTask{
			{Task: "Create a rate-limited queue", How: "Use <code>workqueue.NewRateLimitingQueue()</code>", Example: "queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())"},
			{Task: "Add items and process them", How: "Use <code>queue.Add(key)</code> to enqueue, <code>queue.Get()</code> to dequeue", Example: "queue.Add(\"default/my-pod\")\nitem, _ := queue.Get()\n// process item\nqueue.Done(item)"},
		},
	},
	"client-go/transport": {
		Purpose: "HTTP transport configuration for Kubernetes API calls — TLS, auth tokens, impersonation, proxy settings.",
		Tasks: []PackageTask{
			{Task: "Customize HTTP transport for API calls", How: "Use <code>transport.New(config)</code> to get a configured <code>http.RoundTripper</code>", Example: "rt, _ := transport.New(transportConfig)"},
		},
	},
	"client-go/util/retry": {
		Purpose: "Retry failed Kubernetes API calls with backoff — handles transient errors like 409 Conflict automatically.",
		Tasks: []PackageTask{
			{Task: "Retry on conflict when updating a resource", How: "Use <code>retry.RetryOnConflict(backoff, updateFunc)</code>", Example: "retry.RetryOnConflict(retry.DefaultRetry, func() error {\n  pod, _ := client.Get(ctx, name, metav1.GetOptions{})\n  pod.Labels[\"key\"] = \"new-value\"\n  _, err := client.Update(ctx, pod, metav1.UpdateOptions{})\n  return err\n})"},
		},
	},
	"client-go/discovery": {
		Purpose: "Discover what API resources exist on a Kubernetes cluster — what versions, what resources, what groups.",
		Tasks: []PackageTask{
			{Task: "List all API resources on the cluster", How: "Use <code>clientset.Discovery().ServerResources()</code>", Example: "resources, _ := clientset.Discovery().ServerResources()"},
			{Task: "Check if a CRD exists on the cluster", How: "Use <code>clientset.Discovery().ServerResourcesForGroupVersion(gv)</code>", Example: "res, _ := clientset.Discovery().ServerResourcesForGroupVersion(\"apps/v1\")"},
		},
	},
	"client-go/kubernetes/scheme": {
		Purpose: "Registers all built-in Kubernetes types (Pod, Deployment, Service, etc.) so Go knows how to serialize/deserialize them.",
		Tasks: []PackageTask{
			{Task: "Register your CRD types with the scheme", How: "Call <code>AddToScheme(scheme.Scheme)</code> from your generated code", Example: "myv1.AddToScheme(scheme.Scheme)"},
		},
	},

	// ─── controller-runtime sub-packages ───
	"controller-runtime/pkg/client": {
		Purpose: "High-level Kubernetes client used inside operators — simpler than client-go, works with any resource type.",
		Tasks: []PackageTask{
			{Task: "Get a resource from the cluster", How: "Use <code>client.Get(ctx, key, obj)</code>", Example: "var pod corev1.Pod\nerr := c.Get(ctx, types.NamespacedName{Name: \"web\", Namespace: \"default\"}, &pod)"},
			{Task: "List resources with labels", How: "Use <code>client.List(ctx, list, opts)</code>", Example: "var pods corev1.PodList\nc.List(ctx, &pods, client.InNamespace(\"default\"), client.MatchingLabels{\"app\": \"web\"})"},
			{Task: "Create/Update/Delete a resource", How: "Use <code>c.Create(ctx, obj)</code>, <code>c.Update(ctx, obj)</code>, <code>c.Delete(ctx, obj)</code>", Example: "c.Create(ctx, &myPod)"},
			{Task: "Patch a resource", How: "Use <code>c.Patch(ctx, obj, patch)</code>", Example: "c.Patch(ctx, &pod, client.MergeFrom(oldPod))"},
		},
	},
	"controller-runtime/pkg/manager": {
		Purpose: "The main entry point for a Kubernetes operator — manages controllers, caches, webhooks, leader election, health checks.",
		Tasks: []PackageTask{
			{Task: "Create an operator manager", How: "Use <code>ctrl.NewManager(config, options)</code>", Example: "mgr, _ := ctrl.NewManager(cfg, ctrl.Options{LeaderElection: true})"},
			{Task: "Start all controllers", How: "Call <code>mgr.Start(ctx)</code> — blocks until stopped", Example: "mgr.Start(ctrl.SetupSignalHandler())"},
		},
	},
	"controller-runtime/pkg/reconcile": {
		Purpose: "Defines the reconcile loop pattern — your controller receives a Request (name+namespace) and returns a Result.",
		Tasks: []PackageTask{
			{Task: "Implement the reconcile loop", How: "Implement <code>Reconciler</code> interface with <code>Reconcile(ctx, req) (Result, error)</code>", Example: "func (r *MyReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {\n  // fetch, compare, update\n  return reconcile.Result{}, nil\n}"},
			{Task: "Requeue after a delay", How: "Return <code>reconcile.Result{RequeueAfter: 30 * time.Second}</code>", Example: "return reconcile.Result{RequeueAfter: 30 * time.Second}, nil"},
		},
	},
	"controller-runtime/pkg/controller": {
		Purpose: "Creates and configures Kubernetes controllers — watches resources and triggers reconcile when they change.",
		Tasks: []PackageTask{
			{Task: "Create a controller that watches a resource", How: "Use <code>ctrl.NewControllerManagedBy(mgr).For(&MyResource{}).Complete(reconciler)</code>", Example: "ctrl.NewControllerManagedBy(mgr).\n  For(&myv1.MyApp{}).\n  Owns(&appsv1.Deployment{}).\n  Complete(&MyReconciler{})"},
		},
	},
	"controller-runtime/pkg/webhook": {
		Purpose: "Build admission webhooks for Kubernetes — validate or mutate resources before they're created/updated.",
		Tasks: []PackageTask{
			{Task: "Create a validating webhook", How: "Implement <code>webhook.Validator</code> interface", Example: "func (r *MyResource) ValidateCreate() error { ... }"},
			{Task: "Create a mutating webhook (defaulter)", How: "Implement <code>webhook.Defaulter</code> interface", Example: "func (r *MyResource) Default() { r.Spec.Replicas = 1 }"},
		},
	},
	"controller-runtime/pkg/log": {
		Purpose: "Structured logging for Kubernetes operators — integrates with controller-runtime's log pipeline.",
		Tasks: []PackageTask{
			{Task: "Get a logger in your reconciler", How: "Use <code>log.FromContext(ctx)</code>", Example: "logger := log.FromContext(ctx)\nlogger.Info(\"reconciling\", \"name\", req.Name)"},
		},
	},
	"controller-runtime/pkg/predicate": {
		Purpose: "Filter which events trigger your controller — skip updates you don't care about to reduce reconcile calls.",
		Tasks: []PackageTask{
			{Task: "Only reconcile on spec changes (ignore status)", How: "Use <code>predicate.GenerationChangedPredicate{}</code>", Example: "ctrl.NewControllerManagedBy(mgr).For(&myv1.MyApp{}, builder.WithPredicates(predicate.GenerationChangedPredicate{}))"},
		},
	},

	// ─── k8s.io/api sub-packages ───
	"k8s.io/api/batch/v1": {
		Purpose: "Kubernetes batch workloads — Job (run once) and CronJob (run on schedule).",
		Tasks: []PackageTask{
			{Task: "Define a Job", How: "Create <code>batchv1.Job</code> with a pod template", Example: "job := &batchv1.Job{Spec: batchv1.JobSpec{Template: podTemplate}}"},
			{Task: "Define a CronJob", How: "Create <code>batchv1.CronJob</code> with a schedule", Example: "cj := &batchv1.CronJob{Spec: batchv1.CronJobSpec{Schedule: \"*/5 * * * *\", JobTemplate: jobTemplate}}"},
		},
	},
	"k8s.io/api/rbac/v1": {
		Purpose: "Kubernetes RBAC (Role-Based Access Control) — Role, ClusterRole, RoleBinding, ClusterRoleBinding.",
		Tasks: []PackageTask{
			{Task: "Define a Role with permissions", How: "Create <code>rbacv1.Role</code> with rules", Example: "role := &rbacv1.Role{Rules: []rbacv1.PolicyRule{{APIGroups: []string{\"\"}, Resources: []string{\"pods\"}, Verbs: []string{\"get\", \"list\"}}}}"},
		},
	},
	"k8s.io/api/networking/v1": {
		Purpose: "Kubernetes networking types — Ingress (expose services externally) and NetworkPolicy (firewall rules between pods).",
		Tasks: []PackageTask{
			{Task: "Define an Ingress", How: "Create <code>networkingv1.Ingress</code> with rules", Example: "ingress := &networkingv1.Ingress{Spec: networkingv1.IngressSpec{Rules: []networkingv1.IngressRule{...}}}"},
		},
	},

	// ─── More standard library ───
	"strings": {
		Purpose: "Manipulate text strings — search, replace, split, join, trim, compare, convert case.",
		Tasks: []PackageTask{
			{Task: "Check if string contains a substring", How: "Use <code>strings.Contains(s, substr)</code>", Example: "if strings.Contains(name, \"admin\") { ... }"},
			{Task: "Split a string by delimiter", How: "Use <code>strings.Split(s, sep)</code>", Example: "parts := strings.Split(\"a,b,c\", \",\") // [\"a\", \"b\", \"c\"]"},
			{Task: "Replace text in a string", How: "Use <code>strings.ReplaceAll(s, old, new)</code>", Example: "result := strings.ReplaceAll(url, \"http://\", \"https://\")"},
			{Task: "Trim whitespace", How: "Use <code>strings.TrimSpace(s)</code>", Example: "clean := strings.TrimSpace(\"  hello  \") // \"hello\""},
			{Task: "Join a slice into a string", How: "Use <code>strings.Join(slice, sep)</code>", Example: "csv := strings.Join(names, \",\")"},
		},
	},
	"strconv": {
		Purpose: "Convert between strings and numbers — parse \"42\" to int, format 3.14 to string.",
		Tasks: []PackageTask{
			{Task: "Convert string to int", How: "Use <code>strconv.Atoi(s)</code>", Example: "num, err := strconv.Atoi(\"42\")"},
			{Task: "Convert int to string", How: "Use <code>strconv.Itoa(n)</code>", Example: "s := strconv.Itoa(42)"},
			{Task: "Parse a float from string", How: "Use <code>strconv.ParseFloat(s, 64)</code>", Example: "f, _ := strconv.ParseFloat(\"3.14\", 64)"},
		},
	},
	"time": {
		Purpose: "Work with dates, times, durations, and timers — schedule tasks, measure performance, format timestamps.",
		Tasks: []PackageTask{
			{Task: "Get current time", How: "Use <code>time.Now()</code>", Example: "now := time.Now()"},
			{Task: "Sleep/wait for a duration", How: "Use <code>time.Sleep(duration)</code>", Example: "time.Sleep(5 * time.Second)"},
			{Task: "Format time as string", How: "Use <code>t.Format(layout)</code> with Go's reference time", Example: "s := time.Now().Format(\"2006-01-02 15:04:05\")"},
			{Task: "Measure elapsed time", How: "Use <code>time.Since(start)</code>", Example: "start := time.Now()\n// do work\nfmt.Println(\"took\", time.Since(start))"},
		},
	},
	"errors": {
		Purpose: "Create, wrap, and inspect errors — the standard way to handle failures in Go.",
		Tasks: []PackageTask{
			{Task: "Create a new error", How: "Use <code>errors.New(\"message\")</code>", Example: "err := errors.New(\"file not found\")"},
			{Task: "Check if error is a specific type", How: "Use <code>errors.Is(err, target)</code> or <code>errors.As(err, &target)</code>", Example: "if errors.Is(err, os.ErrNotExist) { ... }"},
			{Task: "Wrap an error with context", How: "Use <code>fmt.Errorf(\"context: %w\", err)</code>", Example: "return fmt.Errorf(\"failed to read config: %w\", err)"},
		},
	},
	"regexp": {
		Purpose: "Regular expressions — match patterns in text, find/replace, validate formats.",
		Tasks: []PackageTask{
			{Task: "Check if string matches a pattern", How: "Use <code>regexp.MatchString(pattern, s)</code>", Example: "matched, _ := regexp.MatchString(`^[a-z]+$`, input)"},
			{Task: "Find and extract matches", How: "Compile with <code>regexp.MustCompile()</code> then use <code>.FindString()</code>", Example: "re := regexp.MustCompile(`\\d+`)\nnums := re.FindAllString(text, -1)"},
		},
	},
	"path/filepath": {
		Purpose: "Work with file paths — join, split, walk directories, get extensions, make relative/absolute paths.",
		Tasks: []PackageTask{
			{Task: "Join path components safely", How: "Use <code>filepath.Join(parts...)</code>", Example: "path := filepath.Join(\"/home\", \"user\", \"config.yaml\")"},
			{Task: "Walk all files in a directory tree", How: "Use <code>filepath.WalkDir(root, walkFunc)</code>", Example: "filepath.WalkDir(\".\", func(path string, d fs.DirEntry, err error) error { ... })"},
			{Task: "Get file extension", How: "Use <code>filepath.Ext(path)</code>", Example: "ext := filepath.Ext(\"config.yaml\") // \".yaml\""},
		},
	},
	"sort": {
		Purpose: "Sort slices and collections — sort numbers, strings, or custom types.",
		Tasks: []PackageTask{
			{Task: "Sort a slice of strings", How: "Use <code>sort.Strings(slice)</code>", Example: "names := []string{\"charlie\", \"alice\", \"bob\"}\nsort.Strings(names) // [\"alice\", \"bob\", \"charlie\"]"},
			{Task: "Sort a slice of ints", How: "Use <code>sort.Ints(slice)</code>", Example: "sort.Ints(numbers)"},
			{Task: "Sort by custom criteria", How: "Use <code>sort.Slice(slice, less)</code>", Example: "sort.Slice(pods, func(i, j int) bool {\n  return pods[i].Name < pods[j].Name\n})"},
		},
	},
	"bufio": {
		Purpose: "Buffered I/O — read files line by line, scan input efficiently, write with buffering for better performance.",
		Tasks: []PackageTask{
			{Task: "Read a file line by line", How: "Use <code>bufio.NewScanner(reader)</code>", Example: "scanner := bufio.NewScanner(file)\nfor scanner.Scan() {\n  line := scanner.Text()\n}"},
		},
	},
	"net": {
		Purpose: "Low-level networking — TCP/UDP connections, DNS lookups, IP addresses, listeners.",
		Tasks: []PackageTask{
			{Task: "Start a TCP server", How: "Use <code>net.Listen(\"tcp\", addr)</code>", Example: "listener, _ := net.Listen(\"tcp\", \":8080\")\nconn, _ := listener.Accept()"},
			{Task: "Connect to a TCP server", How: "Use <code>net.Dial(\"tcp\", addr)</code>", Example: "conn, _ := net.Dial(\"tcp\", \"example.com:80\")"},
			{Task: "Resolve DNS", How: "Use <code>net.LookupHost(hostname)</code>", Example: "addrs, _ := net.LookupHost(\"google.com\")"},
		},
	},
	"crypto/tls": {
		Purpose: "TLS/SSL — secure connections, certificates, HTTPS configuration.",
		Tasks: []PackageTask{
			{Task: "Create a TLS config for HTTPS server", How: "Use <code>tls.Config{}</code> with certificates", Example: "tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}"},
			{Task: "Load a certificate from files", How: "Use <code>tls.LoadX509KeyPair(certFile, keyFile)</code>", Example: "cert, _ := tls.LoadX509KeyPair(\"cert.pem\", \"key.pem\")"},
		},
	},
	"log": {
		Purpose: "Simple logging — print messages with timestamps to stderr. For structured logging, use logrus or zap instead.",
		Tasks: []PackageTask{
			{Task: "Log a message", How: "Use <code>log.Println(\"message\")</code>", Example: "log.Println(\"Server started on port 8080\")"},
			{Task: "Log and exit on fatal error", How: "Use <code>log.Fatal(err)</code>", Example: "if err != nil { log.Fatal(err) }"},
		},
	},

	// ─── echo web framework ───
	"labstack/echo": {
		Purpose: "Minimalist Go web framework — fast HTTP routing, middleware, request binding, response rendering.",
		Tasks: []PackageTask{
			{Task: "Start a web server", How: "Create with <code>echo.New()</code> then <code>e.Start(\":8080\")</code>", Example: "e := echo.New()\ne.GET(\"/hello\", handler)\ne.Start(\":8080\")"},
			{Task: "Handle a route", How: "Use <code>e.GET(\"/path\", handler)</code>", Example: "e.GET(\"/users/:id\", func(c echo.Context) error {\n  id := c.Param(\"id\")\n  return c.JSON(200, user)\n})"},
		},
	},
}

// lookupPackageKnowledge finds knowledge for a package based on its import path
func lookupPackageKnowledge(importPath string) *PackageKnowledge {
	// Try exact match first
	if k, ok := knowledgeBase[importPath]; ok {
		return k
	}
	// Try partial match (e.g., "k8s.io/client-go/kubernetes" contains "client-go/kubernetes")
	for pattern, k := range knowledgeBase {
		if strings.Contains(importPath, pattern) {
			return k
		}
	}
	return nil
}

// lookupStructKnowledge finds knowledge for a struct in a known package
func lookupStructKnowledge(importPath string, structName string) *StructKnowledge {
	k := lookupPackageKnowledge(importPath)
	if k == nil || k.Structs == nil {
		return nil
	}
	if sk, ok := k.Structs[structName]; ok {
		return &sk
	}
	return nil
}

// lookupFuncKnowledge finds knowledge for a function in a known package
func lookupFuncKnowledge(importPath string, funcName string) *FuncKnowledge {
	k := lookupPackageKnowledge(importPath)
	if k == nil || k.Functions == nil {
		return nil
	}
	if fk, ok := k.Functions[funcName]; ok {
		return &fk
	}
	return nil
}

// generateTaskList creates an HTML list of practical tasks with how-to and code hints
func generateTaskList(tasks []PackageTask) string {
	if len(tasks) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("<div class=\"task-list\">")
	for _, t := range tasks {
		sb.WriteString("<div class=\"task-item\">")
		sb.WriteString(fmt.Sprintf("<div class=\"task-title\">🎯 <strong>%s</strong></div>", t.Task))
		sb.WriteString(fmt.Sprintf("<div class=\"task-how\">📝 %s</div>", t.How))
		if t.Example != "" {
			sb.WriteString(fmt.Sprintf("<pre class=\"task-code\"><code>%s</code></pre>", t.Example))
		}
		sb.WriteString("</div>")
	}
	sb.WriteString("</div>")
	return sb.String()
}
