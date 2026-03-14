package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/godoceasy/internal/docs"
	"github.com/godoceasy/internal/fetcher"
	"github.com/godoceasy/internal/parser"
	"github.com/godoceasy/internal/server"
)

const (
	version = "1.0.0"
	usage   = `godoceasy - Beginner-friendly Go documentation viewer

Usage:
  godoceasy <go-package>       Fetch, parse, and serve documentation
  godoceasy --help              Show this help message
  godoceasy --version           Show version

Examples:
  godoceasy fmt
  godoceasy net/http
  godoceasy encoding/json

Popular Go Libraries:
  godoceasy github.com/gin-gonic/gin              # HTTP web framework
  godoceasy github.com/gorilla/mux                 # HTTP router
  godoceasy github.com/sirupsen/logrus             # Structured logger
  godoceasy github.com/spf13/cobra                 # CLI framework
  godoceasy github.com/spf13/viper                 # Configuration
  godoceasy github.com/stretchr/testify            # Testing toolkit
  godoceasy github.com/go-chi/chi                  # Lightweight router
  godoceasy github.com/labstack/echo               # Web framework

Popular Kubernetes Libraries:
  godoceasy k8s.io/client-go                       # Kubernetes Go client
  godoceasy k8s.io/apimachinery                    # API machinery
  godoceasy k8s.io/api                             # Kubernetes API types
  godoceasy sigs.k8s.io/controller-runtime         # Controller runtime
`
)

// popularPackages maps well-known vanity import paths to their actual git repos
var popularPackages = map[string]string{
	"k8s.io/client-go":                  "https://github.com/kubernetes/client-go.git",
	"k8s.io/apimachinery":               "https://github.com/kubernetes/apimachinery.git",
	"k8s.io/api":                        "https://github.com/kubernetes/api.git",
	"k8s.io/kubectl":                    "https://github.com/kubernetes/kubectl.git",
	"k8s.io/kubelet":                    "https://github.com/kubernetes/kubelet.git",
	"k8s.io/utils":                      "https://github.com/kubernetes/utils.git",
	"sigs.k8s.io/controller-runtime":    "https://github.com/kubernetes-sigs/controller-runtime.git",
	"sigs.k8s.io/yaml":                  "https://github.com/kubernetes-sigs/yaml.git",
	"sigs.k8s.io/structured-merge-diff": "https://github.com/kubernetes-sigs/structured-merge-diff.git",
	"go.uber.org/zap":                   "https://github.com/uber-go/zap.git",
	"go.uber.org/multierr":              "https://github.com/uber-go/multierr.git",
	"go.uber.org/atomic":                "https://github.com/uber-go/atomic.git",
	"golang.org/x/net":                  "https://github.com/golang/net.git",
	"golang.org/x/text":                 "https://github.com/golang/text.git",
	"golang.org/x/sync":                 "https://github.com/golang/sync.git",
	"golang.org/x/crypto":               "https://github.com/golang/crypto.git",
	"golang.org/x/tools":                "https://github.com/golang/tools.git",
	"google.golang.org/grpc":            "https://github.com/grpc/grpc-go.git",
	"google.golang.org/protobuf":        "https://github.com/protocolbuffers/protobuf-go.git",
}

func main() {
	if len(os.Args) < 2 {
		fmt.Print(usage)
		os.Exit(1)
	}

	arg := os.Args[1]

	switch arg {
	case "--help", "-h":
		fmt.Print(usage)
		return
	case "--version", "-v":
		fmt.Printf("godoceasy version %s\n", version)
		return
	}

	packagePath := arg
	fmt.Printf("🔍 godoceasy v%s\n", version)
	fmt.Printf("📦 Package: %s\n\n", packagePath)

	// Step 1: Fetch/resolve the package
	fmt.Println("⬇️  Step 1: Fetching package source...")
	srcDir, err := fetcher.Fetch(packagePath, popularPackages)
	if err != nil {
		fmt.Printf("❌ Failed to fetch package: %v\n\n", err)
		fmt.Println("💡 Try one of these popular packages instead:")
		fmt.Println("   godoceasy fmt")
		fmt.Println("   godoceasy net/http")
		fmt.Println("   godoceasy github.com/gin-gonic/gin")
		fmt.Println("   godoceasy k8s.io/client-go")
		fmt.Println("\n   Run 'godoceasy --help' for more options.")
		os.Exit(1)
	}
	fmt.Printf("   ✅ Source at: %s\n\n", srcDir)

	// Step 2: Parse documentation
	fmt.Println("📖 Step 2: Parsing Go documentation...")
	pkgDocs, err := parser.ParsePackage(srcDir, packagePath)
	if err != nil {
		log.Fatalf("❌ Failed to parse package: %v", err)
	}
	fmt.Printf("   ✅ Found %d packages (all sub-packages included, no limit)\n\n", len(pkgDocs))

	// Step 3: Generate beginner-friendly docs
	fmt.Println("✏️  Step 3: Generating beginner-friendly documentation...")
	friendlyDocs := docs.Generate(pkgDocs)
	fmt.Printf("   ✅ Generated docs for %d packages\n\n", len(friendlyDocs))

	// Step 4: Start web server
	port := "8080"
	if p := os.Getenv("GODOCEASY_PORT"); p != "" {
		port = p
	}

	addr := fmt.Sprintf(":%s", port)
	url := fmt.Sprintf("http://localhost:%s", port)

	fmt.Printf("🌐 Step 4: Starting web server at %s\n", url)
	fmt.Println("   Press Ctrl+C to stop")
	fmt.Println("")
	fmt.Println("   💡 Tips:")
	fmt.Println("   • Press / to focus the search bar")
	fmt.Println("   • Search for any struct, function, interface, or type by name")
	fmt.Println("   • Every item has: 🎯 What It Does │ ⏰ When To Use │ 💬 In Simple Words │ 💡 Example")

	// Open browser
	go openBrowser(url)

	// Start server
	srv := server.New(friendlyDocs, packagePath)
	if err := srv.Start(addr); err != nil {
		log.Fatalf("❌ Server error: %v", err)
	}
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	}
	if cmd != nil {
		_ = cmd.Start()
	}
}
