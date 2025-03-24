# test-sample-golang

## Validation process
### ✅ Step 1: Set Up Go in WSL
Make sure you have Go installed in your WSL environment:

```bash
go version
```
If not installed:

```bash
sudo apt update
sudo apt install golang-go
```
You can also install the latest version manually from https://go.dev/doc/install.

### ✅ Step 2: Clone Each Project
Inside WSL, clone the repos:

```bash
git clone https://github.com/go-chi/chi.git
git clone https://github.com/gohugoio/hugo.git
git clone https://github.com/gin-gonic/gin.git
git clone https://github.com/valyala/fasthttp.git
```
Then cd into each one and run tests:

### ✅ Step 3: Run Tests for Each Project
Each Go project can be tested with:

```bash
cd <project-directory>
go test ./...
```
Examples:

🔹 chi:
```bash
cd chi
go test ./...
```
🔹 hugo:
```bash
cd hugo
go test ./...
```
Hugo has more dependencies and may require:

```bash
go mod tidy
go test ./...
```
🔹 gin:
```bash
cd gin
go test ./...
```
🔹 fasthttp:
```bash
cd fasthttp
go test ./...
```
### ✅ Step 4 (Optional): Clean Up & Ensure Modules Are Downloaded
For each project, run:

```bash
go mod tidy
go mod download
```

### 🧪 Tips
Use go test -v ./... for verbose output.

You can run specific test files:
go test -v path/to/file_test.go

You can run with a race detector:
go test -race ./...