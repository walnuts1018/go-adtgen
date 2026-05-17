# Project Rename to go-adtgen Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Rename the project from `go-product-type` to `go-adtgen`, including module path, directives, and generated filenames.

**Architecture:** This is a coordinated rename across the module system, parser/emitter logic, and test data. We follow a sequence of Module Update -> Core Logic Update -> Test Data Update -> Documentation Update.

**Tech Stack:** Go (Module system, AST parsing, templates).

---

### Task 1: Update Go Module Path

**Files:**
- Modify: `go.mod`
- Modify: `main.go`
- Modify: All files in `internal/**/*.go` (internal imports)

- [ ] **Step 1: Update go.mod**

```go
// go.mod
module github.com/walnuts1018/go-adtgen
```

- [ ] **Step 2: Update internal imports**
Use a global find-and-replace to update all internal imports from `github.com/walnuts1018/go-adtgen` to `github.com/walnuts1018/go-adtgen`.

- [ ] **Step 3: Verify module consistency**
Run: `go mod tidy`
Expected: SUCCESS

- [ ] **Step 4: Commit**
```bash
git add go.mod go.sum main.go internal/
git commit -m "chore: rename module to go-adtgen"
```

### Task 2: Update Core Logic Directives and Output Path

**Files:**
- Modify: `internal/parser/parser.go` (if directives are defined there)
- Modify: `internal/loader/loader.go` (if build tags are handled there)
- Modify: `main.go` (output path logic)

- [ ] **Step 1: Identify and update build tag constant**
Search for `adtgen_generate` and change to `adtgen_generate`.

- [ ] **Step 2: Identify and update annotation prefix**
Search for `adtgen:` and change to `adtgen:`.

- [ ] **Step 3: Update default output filename**
In `main.go` (or `outputPath` function), change `zz_generated.adtgen.go` to `zz_generated.adtgen.go`.

- [ ] **Step 4: Commit**
```bash
git add .
git commit -m "feat: update directives and output filename to adtgen"
```

### Task 3: Update Test Data and E2E Tests

**Files:**
- Modify: `internal/testdata/**/*.go`
- Modify: `internal/e2e/e2e_test.go`
- Modify: `internal/parser/parser_test.go` (and other component tests)

- [ ] **Step 1: Update testdata build tags and annotations**
Global replace `adtgen_generate` -> `adtgen_generate` and `adtgen:` -> `adtgen:` in `internal/testdata`.

- [ ] **Step 2: Update E2E test assertions**
Update `internal/e2e/e2e_test.go` to expect the new filename `zz_generated.adtgen.go`.

- [ ] **Step 3: Run all tests**
Run: `go test ./...`
Expected: ALL PASS

- [ ] **Step 4: Commit**
```bash
git add .
git commit -m "test: update testdata and assertions for rename"
```

### Task 4: Update Documentation

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Update README.md**
Update project name, description, and usage examples to use `adtgen` and `go-adtgen`.

- [ ] **Step 2: Commit**
```bash
git add README.md
git commit -m "docs: update README for go-adtgen"
```
