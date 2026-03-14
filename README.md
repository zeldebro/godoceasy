<p align="center">
  <img src="https://img.shields.io/badge/📘_godoceasy-Learn_Go_Packages_The_Easy_Way-326CE5?style=for-the-badge&logoColor=white&labelColor=0D1117" alt="godoceasy" height="60">
</p>

<h1 align="center">godoceasy</h1>
<h3 align="center">🎓 Learn Any Go Package — Without Reading Source Code</h3>

<p align="center">
  <em>One command → Beautiful, searchable, beginner-friendly documentation in your browser.</em><br>
  <em>Search inside any package — find structs, functions, interfaces, types instantly.</em>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go Version">
  <img src="https://img.shields.io/badge/License-MIT-green?style=flat-square" alt="License">
  <img src="https://img.shields.io/badge/Dependencies-Zero-brightgreen?style=flat-square" alt="Zero Dependencies">
  <img src="https://img.shields.io/badge/Platform-macOS%20%7C%20Linux%20%7C%20Windows-blue?style=flat-square" alt="Platform">
  <img src="https://img.shields.io/badge/PRs-Welcome-ff69b4?style=flat-square" alt="PRs Welcome">
</p>

<p align="center">
  <a href="#-what-is-godoceasy"><strong>What Is It</strong></a> · 
  <a href="#-why-does-this-exist--what-problem-does-it-solve"><strong>Why Use It</strong></a> · 
  <a href="#-when-should-you-use-godoceasy"><strong>When To Use</strong></a> · 
  <a href="#-quick-start"><strong>Quick Start</strong></a> · 
  <a href="#-search-inside-any-package"><strong>Search</strong></a> · 
  <a href="#-real-examples--how-useful-is-this"><strong>Real Examples</strong></a> · 
  <a href="#-kubernetes-libraries"><strong>Kubernetes</strong></a> · 
  <a href="#-how-it-works"><strong>How It Works</strong></a>
</p>

---

## 📘 What Is godoceasy?

**godoceasy** is a **learning tool** that turns any Go package into a **searchable, visual, interactive documentation website** — running locally in your browser.

It does **3 things** that no other Go doc tool does:

| What It Does | Why It Matters |
|:---|:---|
| 🔍 **Search inside any package** — find any struct, function, interface, type by name | You don't need to scroll through 1000 lines. Type `Server` → see it instantly with explanation + example |
| 📖 **Explains everything in plain words** — "What It Does", "When To Use", "In Simple Words" | You understand what `HandleFunc` actually does without reading the Go source code |
| 💡 **Copy-paste examples for everything** — every struct, every function, every interface | You can start using the package in 30 seconds — just copy the example into your code |

> **Think of it like this:**
> - `pkg.go.dev` = Reference manual (for experts who already know Go)
> - `godoceasy` = **Learning guide** (for anyone — explains WHY, WHEN, HOW with real examples)

---

## ❓ Why Does This Exist — What Problem Does It Solve?

### The Problem: Go Documentation Is Hard to Learn From

Go's official documentation looks like this:

```go
func HandleFunc(pattern string, handler func(ResponseWriter, *Request))
```

**If you're learning, this tells you NOTHING:**
- ❌ What does this function actually do in real life?
- ❌ When would I use this instead of something else?
- ❌ How do I search for the right struct/function inside a huge package like `net/http` (58+ functions)?
- ❌ Where's a working example I can copy-paste right now?

### The Solution: godoceasy Answers All of This

```
┌────────────────────────────────────────────────────────────────────┐
│  ⚡ Function: HandleFunc                                           │
│  ──────────────────────────────────────────────────────────────── │
│                                                                    │
│  🎯 What It Does:                                                  │
│     Registers a function that runs when someone visits a URL.      │
│     You give it a URL pattern ("/hello") and a function.           │
│     When a user's browser hits that URL → your function runs.      │
│                                                                    │
│  ⏰ When To Use:                                                    │
│     Use this when you're building a web server and want to         │
│     say: "When someone visits /hello, run this code."              │
│     This is the FIRST thing you use when building a Go web app.    │
│                                                                    │
│  💬 In Simple Words:                                                │
│     "Like putting a sign on a door — when someone knocks on        │
│      /hello, your function answers the door."                      │
│                                                                    │
│  💡 Copy-Paste Example:                                    [📋]    │
│  ┌──────────────────────────────────────────────────────────────┐  │
│  │  package main                                                │  │
│  │                                                              │  │
│  │  import "net/http"                                           │  │
│  │  import "fmt"                                                │  │
│  │                                                              │  │
│  │  func main() {                                               │  │
│  │      // STEP 1: Tell the server what to do for /hello        │  │
│  │      http.HandleFunc("/hello", func(w http.ResponseWriter,   │  │
│  │          r *http.Request) {                                  │  │
│  │          fmt.Fprintf(w, "Hello, World!")                     │  │
│  │      })                                                      │  │
│  │                                                              │  │
│  │      // STEP 2: Start the server on port 8080                │  │
│  │      http.ListenAndServe(":8080", nil)                       │  │
│  │  }                                                           │  │
│  └──────────────────────────────────────────────────────────────┘  │
└────────────────────────────────────────────────────────────────────┘
```

**Now you know:** what it does, when to use it, and you have working code. That's the difference.

---

## ⏰ When Should You Use godoceasy?

| Situation | How godoceasy Helps |
|:---|:---|
| 🆕 **"I'm new to Go and don't understand the standard library"** | Every function/struct explained with "In Simple Words" analogies — like a teacher sitting next to you |
| 🔍 **"I need to find a specific struct or function inside a huge package"** | **Search bar** lets you type any name → instantly shows matching structs, functions, interfaces, types with direct links |
| 📦 **"I found a library on GitHub but the docs are just code signatures"** | godoceasy generates "What It Does" + "When To Use" + copy-paste examples for every exported item |
| ☸️ **"I'm building a Kubernetes operator and `client-go` docs are overwhelming"** | Built-in knowledge base explains Kubernetes libraries with practical tasks: "How to list pods", "How to create a deployment" |
| 🏗️ **"I want to understand how a package is structured"** | Left sidebar shows **all structs, interfaces, functions, types** organized — click any to jump to its explanation |
| 🧪 **"I just want a working example to start with"** | Every single item has a **copy-paste ready** example with step-by-step comments |
| 🤝 **"I'm reviewing a PR that uses a library I don't know"** | One command → full documentation in your browser. Understand the library in 5 minutes |

---

## 🚀 Quick Start

### What You Need

Only **one thing**: [Go 1.22+](https://go.dev/dl/) installed.

```bash
# Check if Go is installed
go version
# Output: go version go1.22.x ...  ← You're good!
```

### Install & Run (3 commands)

```bash
# 1. Clone
git clone https://github.com/zeldebro/godoceasy.git
cd godoceasy

# 2. Build (creates a single binary — no other dependencies)
go build -o godoceasy ./cmd/godoceasy

# 3. Run with any Go package
./godoceasy fmt
```

Your browser opens automatically at **http://localhost:8080** 🎉

```
🔍 godoceasy v1.0.0
📦 Package: fmt

⬇️  Step 1: Fetching package source...            ✅
📖 Step 2: Parsing Go documentation...            ✅ Found 1 package
✏️  Step 3: Generating beginner-friendly docs...   ✅
🌐 Step 4: Starting web server at http://localhost:8080
   Press Ctrl+C to stop
```

### Try More Packages

```bash
# Go standard library
./godoceasy net/http          # Web servers & HTTP clients
./godoceasy encoding/json     # JSON encode/decode
./godoceasy os                # File system & environment
./godoceasy context           # Cancellation & timeouts
./godoceasy sync              # Concurrency primitives

# Popular open-source libraries (downloaded automatically)
./godoceasy github.com/gin-gonic/gin          # Web framework
./godoceasy github.com/spf13/cobra            # CLI framework
./godoceasy github.com/sirupsen/logrus        # Structured logging

# Kubernetes libraries
./godoceasy k8s.io/client-go                  # K8s API client
./godoceasy sigs.k8s.io/controller-runtime    # Operator framework
```

### Change Port

```bash
GODOCEASY_PORT=9090 ./godoceasy fmt
# → Opens http://localhost:9090
```

---

## 🔍 Search Inside Any Package

This is the **most useful feature**. When you open a package with 50+ functions and 20+ structs, you don't want to scroll. You want to **search**.

### How Search Works

1. **Top search bar** — type any name: `Server`, `Handler`, `Marshal`, `Client`
2. **Results show instantly** — with type badge (Struct / Function / Interface / Type / Method)
3. **Click any result** — jumps directly to the full explanation with example

### What You Can Search For

| Search For | What You Find | Example |
|:---|:---|:---|
| **Struct name** | Full struct with fields, methods, example | Search `Server` → finds `http.Server` with all 13 fields explained |
| **Function name** | Signature + what it does + when to use + example | Search `HandleFunc` → finds the function with copy-paste code |
| **Interface name** | Contract definition + how to implement + example | Search `Handler` → finds `http.Handler` interface with implementation guide |
| **Type name** | Type alias explanation + why it exists + example | Search `HandlerFunc` → explains this type adapter pattern |
| **Method name** | Struct method with receiver + parameters + example | Search `ListenAndServe` → finds `Server.ListenAndServe()` |
| **Partial match** | Any name containing your search text | Search `read` → finds `ReadAll`, `ReadCloser`, `ReadRequest`, etc. |

### Example: Searching Inside `net/http`

```
Search: "server"

Results:
  📐 Struct   │ Server          │ HTTP server with configurable settings
  ⚡ Function │ ListenAndServe  │ Start an HTTP server on a port
  ⚡ Function │ ListenAndServeTLS │ Start an HTTPS server
  ⚡ Function │ Serve           │ Serve connections on a listener
  ⚡ Method   │ Server.Close    │ Immediately close all connections
  ⚡ Method   │ Server.Shutdown │ Gracefully shutdown the server
```

> 💡 **Why this matters:** In `net/http`, there are 58+ functions and 18+ structs. Without search, you'd spend 10 minutes scrolling. With search, you find what you need in 2 seconds.

---

## 🖥️ What You See In The Browser

### Left Sidebar — Navigate Everything

Every struct, interface, function, and type is listed in the sidebar. Click any item to jump directly to it.

```
┌─────────────────┬─────────────────────────────────────────────┐
│ 📦 Packages     │                                             │
│  ├ http         │   📦 Package: net/http                      │
│  ├ url          │   ────────────────────────────               │
│  └ cookiejar    │   📐 18 Structs │ 🔌 3 Interfaces │ ⚡ 58 Fn│
│                 │                                             │
│ 📑 This Package │   🎯 What It Does:                          │
│  ├ Overview     │   Go's built-in web toolkit. Build servers, │
│  ├ 📐 Structs   │   make HTTP requests, handle routes.        │
│  │  ├ Client    │                                             │
│  │  ├ Server    │   ⏰ When To Use:                            │
│  │  ├ Request   │   • Building a web server or REST API       │
│  │  ├ Response  │   • Making HTTP requests to other services  │
│  │  └ ...       │   • Handling URL routing                    │
│  ├ 🔌 Interfaces│                                             │
│  │  ├ Handler   │   💬 In Simple Words:                        │
│  │  └ ...       │   "Everything you need to build a website   │
│  ├ ⚡ Functions  │    or talk to web APIs — all built into Go."│
│  │  ├ Get       │                                             │
│  │  ├ Post      │                                             │
│  │  ├ Handle    │                                             │
│  │  └ ...       │                                             │
│  ├ 🏷️ Types     │                                             │
│  ├ 📌 Constants │                                             │
│  └ 📎 Variables │                                             │
└─────────────────┴─────────────────────────────────────────────┘
```

### Every Item Has 4 Sections

No matter what you click — struct, function, interface, or type — you always see:

```
┌────────────────────────────────────────────────────────────────┐
│  📐 Struct: Client                                             │
│  ────────────────────────────────────────────────────────────  │
│                                                                │
│  🎯 WHAT IT DOES                                               │
│  • Makes HTTP requests (GET, POST, PUT, DELETE)                │
│  • Manages cookies, redirects, and timeouts automatically      │
│  • Has 4 fields: Transport, CheckRedirect, Jar, Timeout       │
│  • Has 8 methods: .Do(), .Get(), .Post(), .Head()...           │
│                                                                │
│  ⏰ WHEN TO USE                                                 │
│  • Use http.Client when you need to make requests to other     │
│    services (HTTP calls, API requests, etc.)                   │
│  • Create it once → reuse it for many requests                 │
│  • An empty http.Client{} works fine — customize only when     │
│    you need timeouts, proxies, or TLS settings                 │
│                                                                │
│  💬 IN SIMPLE WORDS                                             │
│  "📱 Think of Client like a phone. You set it up once          │
│   (phone number, settings), then use it to make calls          │
│   (send requests). The fields are phone settings, the          │
│   methods are the actions."                                    │
│                                                                │
│  💡 COPY-PASTE EXAMPLE                                  [📋]   │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  package main                                            │  │
│  │                                                          │  │
│  │  import "net/http"                                       │  │
│  │  import "fmt"                                            │  │
│  │                                                          │  │
│  │  func main() {                                           │  │
│  │      // STEP 1: Create a Client (like setting up a phone)│  │
│  │      myClient := http.Client{                            │  │
│  │          Timeout: 30 * time.Second, // how long to wait  │  │
│  │      }                                                   │  │
│  │                                                          │  │
│  │      // STEP 2: Use it! Make a GET request               │  │
│  │      resp, err := myClient.Get("https://example.com")    │  │
│  │      if err != nil {                                     │  │
│  │          fmt.Println("Oops:", err)                       │  │
│  │      }                                                   │  │
│  │      fmt.Println("Status:", resp.StatusCode)             │  │
│  │  }                                                       │  │
│  └──────────────────────────────────────────────────────────┘  │
│                                                                │
│  📋 FIELDS                                                     │
│  ┌──────────────────┬──────────────┬─────────────────────────┐ │
│  │ Field            │ Type         │ What It Means           │ │
│  ├──────────────────┼──────────────┼─────────────────────────┤ │
│  │ Transport        │ RoundTripper │ controls how requests   │ │
│  │                  │              │ are sent (proxy, TLS)   │ │
│  │ Timeout          │ Duration     │ how long to wait        │ │
│  │ CheckRedirect    │ func         │ what to do on redirect  │ │
│  │ Jar              │ CookieJar   │ stores cookies           │ │
│  └──────────────────┴──────────────┴─────────────────────────┘ │
│                                                                │
│  ⚡ METHODS                                                     │
│  ├ .Do(req)        — Send any custom request                   │
│  ├ .Get(url)       — Make a GET request                        │
│  ├ .Post(url,...)  — Make a POST request                       │
│  └ .Head(url)      — Make a HEAD request                       │
└────────────────────────────────────────────────────────────────┘
```

### Builder Pattern Detection

When a struct uses the builder pattern (chain `.With()` / `.Set()` calls), godoceasy detects and explains it:

```
🔗 Builder Pattern Detected!

  obj.WithTimeout(5*time.Second).WithRetries(3).Build()
  ↑ Each method returns the same object — you chain them together
    like filling in a form field by field.
```

### 🗺️ Package Structure Diagram

Every package page includes an **interactive structure map** — a visual diagram showing the entire package at a glance:

```
┌──────────────────────────────────────────────────────────────────┐
│  🗺️ Package Structure Map                                        │
│  Interactive diagram — click any item to jump to its docs        │
│                                                                  │
│  📦 net/http                                                     │
│  📐 18 Structs │ 🔌 5 Interfaces │ ⚡ 58 Functions │ 🏷️ 9 Types  │
│                                                                  │
│  ── 📐 Structs ──────────────────────────────────────────────    │
│  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐             │
│  │ 📐 Server    │ │ 📐 Client    │ │ 📐 Request   │             │
│  │ ──────────── │ │ ──────────── │ │ ──────────── │             │
│  │ Addr  string │ │ Transport    │ │ Method string│             │
│  │ Handler      │ │ Timeout      │ │ URL    *URL  │             │
│  │ TLSConfig    │ │ Jar          │ │ Header       │             │
│  │ +10 more     │ │ CheckRedirect│ │ Body         │             │
│  │ ──────────── │ │ ──────────── │ │ +12 more     │             │
│  │ .Serve()     │ │ .Do()        │ │ ──────────── │             │
│  │ .Shutdown()  │ │ .Get()       │ │ .Cookie()    │             │
│  │ .Close()     │ │ .Post()      │ │ .AddCookie() │             │
│  └──────────────┘ └──────────────┘ └──────────────┘             │
│                                                                  │
│  ── 🔌 Interfaces ──────────────────────────────────────────     │
│  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐             │
│  │ 🔌 Handler   │ │ 🔌 Flusher   │ │ 🔌 Hijacker  │             │
│  │ .ServeHTTP() │ │ .Flush()     │ │ .Hijack()    │             │
│  └──────────────┘ └──────────────┘ └──────────────┘             │
│                                                                  │
│  ── ⚡ Functions ────────────────────────────────────────────     │
│  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐             │
│  │ ⚡ Get       │ │ ⚡ Post      │ │ ⚡ HandleFunc │             │
│  │ 📥 url      │ │ 📥 url,type  │ │ 📥 pattern   │             │
│  │ 📤 *Response│ │  body        │ │    handler    │             │
│  └──────────────┘ └──────────────┘ └──────────────┘             │
│                                                                  │
│  ── 🏷️ Types ───────────────────────────────────────────────     │
│  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐             │
│  │ 🏷️ HandlerFn │ │ 🏷️ ConnState│ │ 🏷️ Header    │             │
│  └──────────────┘ └──────────────┘ └──────────────┘             │
└──────────────────────────────────────────────────────────────────┘
```

**What the diagram shows:**
- **Structs** — with their top fields and methods listed inside the card
- **Interfaces** — with required methods shown
- **Functions** — with parameters and return types
- **Types** — named types and aliases
- **Click any card** → jumps directly to the full documentation for that item

---

## 📚 Real Examples — How Useful Is This?

### Example 1: "I need to read a JSON file in Go"

```bash
./godoceasy encoding/json
```

Search `Unmarshal` → you instantly see:

```
🎯 What It Does:  Converts JSON bytes into a Go struct
⏰ When To Use:   When you receive JSON data and need to use it in your code
💡 Example:
    var user User
    json.Unmarshal(jsonBytes, &user)
    fmt.Println(user.Name)
```

**Without godoceasy:** You'd read `func Unmarshal(data []byte, v any) error` and guess what `v any` means.

### Example 2: "I need to build a CLI tool like kubectl"

```bash
./godoceasy github.com/spf13/cobra
```

Search `Command` → you see the full struct with:
- Every field explained ("Use" = the command name, "Short" = one-line description)
- How to create a root command
- How to add sub-commands
- How to add flags
- Working example you can copy into `main.go`

**Without godoceasy:** You'd read cobra's GitHub README (which is good but doesn't explain every field).

### Example 3: "I'm debugging a Kubernetes operator and need to understand client-go"

```bash
./godoceasy k8s.io/client-go
```

Search `Clientset` → you see:

```
🎯 What It Does:
   The main entry point for ALL Kubernetes API calls.
   It holds authenticated connections to every API group.

⏰ When To Use:
   Create ONE Clientset at startup, then use it everywhere
   to list pods, create deployments, delete services.

💡 Task List:
   🎯 List all pods        → clientset.CoreV1().Pods("default").List(...)
   🎯 Create a deployment  → clientset.AppsV1().Deployments(ns).Create(...)
   🎯 Watch for changes    → clientset.CoreV1().Pods(ns).Watch(...)
   🎯 Delete a pod         → clientset.CoreV1().Pods(ns).Delete(...)
```

**Without godoceasy:** You'd spend 30 minutes reading `client-go` source code trying to figure out where to start.

### Example 4: "What's the difference between Handler and HandlerFunc in net/http?"

```bash
./godoceasy net/http
```

Search `Handler` → see both results side by side:

```
🔌 Interface: Handler
   "A contract — any type with a ServeHTTP() method is a Handler"

🏷️ Type: HandlerFunc  
   "A shortcut — converts a plain function into a Handler
    so you don't need to create a struct"
```

Now you understand the difference in 10 seconds.

---

## ☸️ Kubernetes Libraries

godoceasy has **built-in practical knowledge** for Kubernetes libraries. Instead of just showing function signatures, it shows **real tasks** you can accomplish:

```bash
./godoceasy k8s.io/client-go          # API client
./godoceasy k8s.io/api                # Pod, Deployment, Service types
./godoceasy k8s.io/apimachinery       # ListOptions, ObjectMeta
./godoceasy sigs.k8s.io/controller-runtime   # Operator framework
```

### What Makes Kubernetes Docs Special Here

| Regular Docs Show | godoceasy Shows |
|:---|:---|
| `func NewForConfig(c *rest.Config) (*Clientset, error)` | **Task:** "Create a clientset to talk to the API" → `clientset, err := kubernetes.NewForConfig(config)` |
| `type ListOptions struct { ... 8 fields ... }` | **Task:** "Filter pods by label" → `ListOptions{LabelSelector: "app=nginx"}` |
| `type ObjectMeta struct { ... 15 fields ... }` | **In Simple Words:** "The identity card of every K8s resource — name, namespace, labels, annotations" |

### Pre-Mapped Vanity Imports

These Kubernetes import paths are automatically resolved:

| You Type | godoceasy Resolves To |
|:---|:---|
| `k8s.io/client-go` | `github.com/kubernetes/client-go` |
| `k8s.io/api` | `github.com/kubernetes/api` |
| `k8s.io/apimachinery` | `github.com/kubernetes/apimachinery` |
| `k8s.io/kubectl` | `github.com/kubernetes/kubectl` |
| `k8s.io/utils` | `github.com/kubernetes/utils` |
| `sigs.k8s.io/controller-runtime` | `github.com/kubernetes-sigs/controller-runtime` |
| `sigs.k8s.io/yaml` | `github.com/kubernetes-sigs/yaml` |
| `go.uber.org/zap` | `github.com/uber-go/zap` |
| `google.golang.org/grpc` | `github.com/grpc/grpc-go` |

---

## 🧠 How It Works

```
  You type:  ./godoceasy net/http
              │
              ▼
  ┌─────────────────────────────────────────────────────┐
  │  📥 STEP 1: FETCH                                   │
  │                                                     │
  │  Finds the package source code:                     │
  │  • Standard lib (fmt, net/http) → uses your GOROOT  │
  │  • Kubernetes (k8s.io/*) → resolves vanity URL →    │
  │    git clone from GitHub                            │
  │  • Any GitHub package → git clone (shallow, fast)   │
  │  • Fallback → go mod download                       │
  └───────────────────────┬─────────────────────────────┘
                          ▼
  ┌─────────────────────────────────────────────────────┐
  │  🔬 STEP 2: PARSE                                   │
  │                                                     │
  │  Reads every .go file using Go's built-in AST       │
  │  parser and extracts:                               │
  │  • Structs (with fields + methods)                  │
  │  • Functions (with params + returns)                │
  │  • Interfaces (with required methods)               │
  │  • Types, Constants, Variables                      │
  │  • Doc comments                                     │
  │  Walks ALL sub-packages — no depth limit            │
  └───────────────────────┬─────────────────────────────┘
                          ▼
  ┌─────────────────────────────────────────────────────┐
  │  ✏️  STEP 3: EXPLAIN                                 │
  │                                                     │
  │  For every item, generates:                         │
  │  • 🎯 "What It Does" — summarized from doc +       │
  │       smart detection of patterns                   │
  │  • ⏰ "When To Use" — practical guidance based on   │
  │       function name, params, return type            │
  │  • 💬 "In Simple Words" — real-world analogy        │
  │  • 💡 "Example" — copy-paste code with comments     │
  │                                                     │
  │  Special: Builder pattern, New*/Must*/Is*/With*     │
  │  detection + knowledge base for popular packages    │
  └───────────────────────┬─────────────────────────────┘
                          ▼
  ┌─────────────────────────────────────────────────────┐
  │  🌐 STEP 4: SERVE                                   │
  │                                                     │
  │  Starts a local web server with:                    │
  │  • 📑 Left sidebar — navigate structs/funcs/types   │
  │  • 🔍 Search bar — find anything by name            │
  │  • 📋 Copy buttons — on every code example          │
  │  • 🗺️ Structure diagram — visual map of the package │
  │  • 🔗 Deep links — every item has an anchor URL     │
  │  • 📦 Package grid — browse sub-packages            │
  │                                                     │
  │  Opens browser automatically → http://localhost:8080│
  └─────────────────────────────────────────────────────┘
```

---

## 📁 Project Structure

```
godoceasy/
│
├── cmd/godoceasy/
│   └── main.go                    # CLI entry point — parses args, runs all 4 steps
│
├── internal/
│   ├── fetcher/
│   │   └── github.go              # Downloads package source (git clone / go mod / GOROOT)
│   │
│   ├── parser/
│   │   └── ast_parser.go          # Reads .go files → extracts structs, funcs, interfaces
│   │
│   ├── docs/
│   │   ├── generator.go           # Creates "What It Does", "When To Use", "In Simple Words"
│   │   └── package_knowledge.go   # Built-in knowledge for K8s, gin, cobra, std lib, etc.
│   │
│   └── server/
│       ├── server.go              # HTTP server + JSON search API
│       ├── templates/
│       │   ├── index.html         # Home page with package grid
│       │   ├── package.html       # Package detail page (all structs/funcs/etc.)
│       │   └── search.html        # Search results page
│       └── static/
│           ├── css/style.css      # Dark theme UI
│           └── js/search.js       # Live search with dropdown
│
├── go.mod                         # Go module — ZERO external dependencies
├── CONTRIBUTING.md
├── LICENSE
└── README.md
```

---

## 🆚 How Is godoceasy Different?

| | [pkg.go.dev](https://pkg.go.dev) | `go doc` (CLI) | **godoceasy** |
|:---|:---:|:---:|:---:|
| Explains "What It Does" in plain words | ❌ Only shows doc comments | ❌ Raw signatures | ✅ Auto-generated explanations |
| Explains "When To Use" | ❌ | ❌ | ✅ Practical guidance for every item |
| "In Simple Words" real-world analogies | ❌ | ❌ | ✅ "Think of Server like a restaurant..." |
| **Search** structs/funcs/interfaces/types | ✅ (needs internet) | ❌ | ✅ **Local, instant, with type badges** |
| Copy-paste examples for everything | ❌ Only if author wrote them | ❌ | ✅ Auto-generated for every item |
| Builder pattern detection | ❌ | ❌ | ✅ Detects and explains chaining |
| Kubernetes task-based docs | ❌ | ❌ | ✅ "How to list pods", "How to create deployment" |
| **Package structure diagram** | ❌ | ❌ | ✅ **Visual map of structs, funcs, interfaces, types** |
| Works offline | ❌ | ✅ | ✅ (std lib fully offline) |
| Beautiful web UI with sidebar | ✅ | ❌ | ✅ |
| Zero dependencies | — | — | ✅ Single binary, nothing to install |

---

## ❓ Frequently Asked Questions

<details>
<summary><strong>Do I need to know Go to use this?</strong></summary>

**No.** You only need Go installed to build the binary. After that, godoceasy explains packages in plain English — you don't need to read any Go source code.
</details>

<details>
<summary><strong>Can I search for a specific function inside a package?</strong></summary>

**Yes!** That's one of the main features. Use the search bar at the top — type any struct name, function name, interface name, type name, or even a partial match. Results appear instantly with type badges (Struct / Function / Interface / Type / Method) and direct links.
</details>

<details>
<summary><strong>Does it work with private/internal packages?</strong></summary>

**Yes.** If you can `git clone` the repository from your machine (with proper credentials), godoceasy can fetch and document it.
</details>

<details>
<summary><strong>Does it use AI or LLMs?</strong></summary>

**No.** All explanations are generated using smart pattern matching, naming conventions, and a built-in knowledge base. No API keys, no internet needed for explanations. Everything runs locally.
</details>

<details>
<summary><strong>Why are some explanations better than others?</strong></summary>

godoceasy has a **built-in knowledge base** for popular packages (Kubernetes, gin, cobra, standard library). For these, you get task-based docs like "How to list pods." For unknown packages, it generates explanations based on function names, parameter types, and doc comments — still useful, but less specific.
</details>

<details>
<summary><strong>Can I use this to learn Go from scratch?</strong></summary>

**Absolutely.** Start with `./godoceasy fmt` (simplest package), then try `./godoceasy net/http`, then `./godoceasy encoding/json`. Each one builds on what you learn. The "In Simple Words" analogies and step-by-step examples make it like having a tutor.
</details>

---

## 🤝 Contributing

See **[CONTRIBUTING.md](CONTRIBUTING.md)** for the full guide.

| What | File | Difficulty |
|:-----|:-----|:----------|
| 🐛 Fix a parsing bug | `internal/parser/ast_parser.go` | 🟢 Easy |
| ✨ Better sample values in examples | `internal/docs/generator.go` | 🟢 Easy |
| 📦 Add more vanity import mappings | `cmd/godoceasy/main.go` | 🟢 Easy |
| 📖 Add knowledge for a new package | `internal/docs/package_knowledge.go` | 🟢 Easy |
| 🔍 Improve search ranking | `internal/server/server.go` | 🟡 Medium |
| 🎨 UI improvements | `internal/server/static/` | 🟡 Medium |
| 📤 Export docs to Markdown/PDF | New feature | 🔴 Advanced |

```bash
# Fork → clone → test
go build -o godoceasy ./cmd/godoceasy
go vet ./...
./godoceasy fmt            # Test with standard library
./godoceasy net/http       # Test with a bigger package
```

---

## 🛠️ Built With

<p>
  <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/badge/HTML5-E34F26?style=for-the-badge&logo=html5&logoColor=white" alt="HTML5">
  <img src="https://img.shields.io/badge/CSS3-1572B6?style=for-the-badge&logo=css3&logoColor=white" alt="CSS3">
  <img src="https://img.shields.io/badge/JavaScript-F7DF1E?style=for-the-badge&logo=javascript&logoColor=black" alt="JavaScript">
</p>

- **100% Go standard library** — `go/parser`, `go/ast`, `net/http`, `html/template`, `embed`
- **Zero external dependencies** — no npm, no frameworks, no build tools
- **Single binary** — all HTML/CSS/JS embedded with `go:embed`

---

## 📜 License

MIT — see **[LICENSE](LICENSE)** for details.

---

<p align="center">
  <strong>⭐ Star this repo if godoceasy helped you understand a Go package faster!</strong>
</p>

<p align="center">
  <em>Built for learners, by learners 💙</em>
</p>

