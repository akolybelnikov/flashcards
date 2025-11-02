# Go Mock Generation - Research & Implementation Plan

**Date:** November 2, 2025  
**Task:** Replace handwritten mocks with generated mocks using industry-standard tools

---

## 🔍 Major Go Mocking Libraries

### 1. **gomock** (Official Google/Uber)
- **Repository:** https://github.com/uber-go/mock (formerly golang/mock)
- **Command:** `mockgen`
- **Status:** ✅ Most popular, industry standard
- **Maintained by:** Uber (originally Google)
- **Go Generate:** ✅ Yes
- **Type-safe:** ✅ Yes

**Example:**
```go
//go:generate mockgen -source=flashcard_service.go -destination=mocks/mock_flashcard_service.go -package=mocks
```

### 2. **testify/mock**
- **Repository:** https://github.com/stretchr/testify
- **Status:** ✅ Very popular
- **Maintained by:** Stretchr
- **Go Generate:** ❌ Manual (handwritten, but with helpers)
- **Type-safe:** ⚠️ Less strict

**Note:** Testify provides mock helpers but requires manual implementation (not auto-generated).

### 3. **moq**
- **Repository:** https://github.com/matryer/moq
- **Command:** `moq`
- **Status:** ✅ Active
- **Go Generate:** ✅ Yes
- **Type-safe:** ✅ Yes
- **Philosophy:** Simpler, less feature-rich than gomock

**Example:**
```go
//go:generate moq -out mocks/flashcard_service_mock.go . FlashcardServiceInterface
```

---

## 📊 Comparison

| Feature | gomock | testify/mock | moq |
|---------|--------|--------------|-----|
| **Auto-generation** | ✅ Yes | ❌ No | ✅ Yes |
| **Type-safety** | ✅✅✅ | ⚠️ Medium | ✅✅ |
| **Industry adoption** | ✅✅✅ | ✅✅✅ | ✅ |
| **Learning curve** | Medium | Easy | Easy |
| **Expectations** | ✅ Rich | ✅ Rich | ⚠️ Basic |
| **Go Generate** | ✅ Yes | ❌ N/A | ✅ Yes |
| **Stars (GitHub)** | 9.5k+ | 23k+ | 2k+ |

---

## 🎯 Recommendation: **gomock/mockgen**

**Why:**
1. ✅ Official industry standard (Google → Uber)
2. ✅ Full auto-generation from interfaces
3. ✅ Type-safe with compile-time checks
4. ✅ Rich expectation API (EXPECT(), Times(), Return())
5. ✅ Perfect for our use case (interface-based architecture)
6. ✅ Works with `go generate`

---

## 📁 Standard Mock Organization

### Option 1: `mocks/` folder (Recommended)
```
flashcards/
├── services/
│   ├── flashcard_service.go
│   ├── flashcard_service_test.go
│   ├── llm_client.go
│   ├── translation_cache.go
│   └── mocks/                          ← Generated mocks here
│       ├── mock_flashcard_service.go
│       ├── mock_llm_client.go
│       └── mock_translation_cache.go
├── db/
│   ├── flashcard_repository.go
│   └── mocks/
│       └── mock_flashcard_repository.go
└── handlers/
    └── (uses service mocks)
```

### Option 2: `_test.go` suffix (Alternative)
```
flashcards/
├── services/
│   ├── flashcard_service.go
│   ├── flashcard_service_test.go
│   ├── mock_flashcard_service_test.go  ← Same package, test-only
│   ├── llm_client.go
│   └── mock_llm_client_test.go
```

**Recommendation:** **Option 1 (mocks/ folder)** because:
- ✅ Cleaner separation
- ✅ Mocks can be shared across packages
- ✅ Easier to gitignore if needed
- ✅ Industry standard

---

## 🔄 Migration Plan

### Current State

**Files with handwritten mocks:**
1. `services/mocks.go` - MockLLMClient, MockTranslationCache
2. `services/flashcard_service_test.go` - mockRepo (inline)
3. `handlers/flashcard_handler_test.go` - mockService (inline)

**Interfaces to mock:**
1. `services.FlashcardServiceInterface` (handlers use this)
2. `services.LLMClient`
3. `services.TranslationCache`
4. `db.FlashcardRepository`

### Step-by-Step Migration

#### Phase 1: Setup (Install & Configure)
```bash
# Install mockgen
go install go.uber.org/mock/mockgen@latest

# Verify installation
mockgen -version
```

#### Phase 2: Add go:generate directives

**File: `services/flashcard_service.go`**
```go
//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

package services
// ... rest of file
```

**File: `services/llm_client.go`**
```go
//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

package services
// ... rest of file
```

**File: `services/translation_cache.go`**
```go
//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

package services
// ... rest of file
```

**File: `db/flashcard_repository.go`**
```go
//go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE -package=mocks

package db
// ... rest of file
```

#### Phase 3: Update Makefile

**Add to Makefile:**
```makefile
# Code generation
generate:
	@echo "Generating mocks..."
	go generate ./...
	@echo "✓ Mocks generated"

# Update test target to ensure mocks are fresh
test: generate
	@echo "Running tests..."
	go test ./... -v

# Clean generated files
clean-mocks:
	@echo "Cleaning generated mocks..."
	find . -path "*/mocks/mock_*.go" -delete
	@echo "✓ Mocks cleaned"

clean: clean-mocks
	@echo "Cleaning build artifacts..."
	rm -rf $(BIN_DIR)
```

#### Phase 4: Generate Mocks

```bash
# Generate all mocks
make generate

# Or directly
go generate ./...
```

This creates:
```
services/mocks/mock_flashcard_service.go
services/mocks/mock_llm_client.go
services/mocks/mock_translation_cache.go
db/mocks/mock_flashcard_repository.go
```

#### Phase 5: Update Tests

**Before (handwritten):**
```go
// services/flashcard_service_test.go
type mockRepo struct{}

func (m *mockRepo) Create(req *models.CreateFlashcardRequest) (*models.Flashcard, error) {
	now := time.Now()
	return &models.Flashcard{ID: 1, Question: req.Question, Answer: req.Answer}, nil
}
// ... 50+ lines of mock implementations
```

**After (generated):**
```go
// services/flashcard_service_test.go
import (
	"github.com/akolybelnikov/flashcards/db/mocks"
	"go.uber.org/mock/gomock"
)

func TestCreateFlashcardValidation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mocks
	mockRepo := mocks.NewMockFlashcardRepository(ctrl)
	mockLLM := servicemocks.NewMockLLMClient(ctrl)
	mockCache := servicemocks.NewMockTranslationCache(ctrl)

	// Set expectations
	mockRepo.EXPECT().
		Create(gomock.Any()).
		Return(&models.Flashcard{
			ID:       1,
			Question: "hello",
			Answer:   "γεια σας",
		}, nil)

	// Test
	svc := NewFlashcardService(mockRepo, mockLLM, mockCache)
	fc, err := svc.CreateFlashcard(&models.CreateFlashcardRequest{
		Question: "hello",
		Answer:   "γεια σας",
	})

	// Assertions
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fc.Question != "hello" {
		t.Fatalf("expected question 'hello', got '%s'", fc.Question)
	}
}
```

#### Phase 6: Remove Handwritten Mocks

**Files to delete:**
- ❌ `services/mocks.go` (replaced by generated mocks)

**Code to remove:**
- ❌ Inline `mockRepo` in `flashcard_service_test.go`
- ❌ Inline `mockService` in `flashcard_handler_test.go`

---

## 📝 Example Generated Mock Structure

**Generated: `services/mocks/mock_llm_client.go`**
```go
// Code generated by MockGen. DO NOT EDIT.
// Source: llm_client.go

package mocks

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockLLMClient is a mock of LLMClient interface.
type MockLLMClient struct {
	ctrl     *gomock.Controller
	recorder *MockLLMClientMockRecorder
}

// MockLLMClientMockRecorder is the mock recorder for MockLLMClient.
type MockLLMClientMockRecorder struct {
	mock *MockLLMClient
}

// NewMockLLMClient creates a new mock instance.
func NewMockLLMClient(ctrl *gomock.Controller) *MockLLMClient {
	mock := &MockLLMClient{ctrl: ctrl}
	mock.recorder = &MockLLMClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLLMClient) EXPECT() *MockLLMClientMockRecorder {
	return m.recorder
}

// Translate mocks base method.
func (m *MockLLMClient) Translate(ctx context.Context, text, sourceLang, targetLang string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Translate", ctx, text, sourceLang, targetLang)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Translate indicates an expected call of Translate.
func (mr *MockLLMClientMockRecorder) Translate(ctx, text, sourceLang, targetLang interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Translate", reflect.TypeOf((*MockLLMClient)(nil).Translate), ctx, text, sourceLang, targetLang)
}
```

---

## 🎨 Visual Comparison

### Current: Handwritten Mocks

```
services/
├── flashcard_service.go          (Interface definition)
├── flashcard_service_test.go     (Tests + inline mockRepo ~60 lines)
├── mocks.go                       (MockLLMClient, MockTranslationCache ~80 lines)
└── llm_client.go

Total mock code: ~140 lines manually maintained
```

### After: Generated Mocks

```
services/
├── flashcard_service.go          //go:generate mockgen ...
├── flashcard_service_test.go     (Tests only, ~40 lines)
├── llm_client.go                  //go:generate mockgen ...
├── translation_cache.go           //go:generate mockgen ...
└── mocks/                         ← Auto-generated
    ├── mock_flashcard_service.go  (Generated, ~200 lines)
    ├── mock_llm_client.go         (Generated, ~100 lines)
    └── mock_translation_cache.go  (Generated, ~150 lines)

Total mock code: 0 lines manually maintained ✅
Total generated: ~450 lines (auto-generated, never edit)
```

---

## 📋 Implementation Checklist

### Preparation
- [ ] Install mockgen: `go install go.uber.org/mock/mockgen@latest`
- [ ] Verify installation: `mockgen -version`
- [ ] Create `.gitignore` entry for mocks (optional)

### Code Changes
- [ ] Add `//go:generate` directives to interface files
- [ ] Update Makefile with `generate` target
- [ ] Run `make generate` to create initial mocks
- [ ] Verify mocks are generated correctly

### Test Updates
- [ ] Update `services/flashcard_service_test.go`
  - Import `gomock` and generated mocks
  - Replace `mockRepo` with `mocks.NewMockFlashcardRepository`
  - Add `gomock.Controller` setup
  - Add expectations with `EXPECT()`
- [ ] Update `handlers/flashcard_handler_test.go`
  - Replace `mockService` with generated mock
  - Add expectations
- [ ] Remove `services/mocks.go` file

### Cleanup
- [ ] Run tests: `make test`
- [ ] Delete old mock code
- [ ] Update documentation

---

## 💰 Benefits

### Before (Handwritten)
- ❌ ~140 lines of manually maintained mock code
- ❌ Easy to forget updating mocks when interfaces change
- ❌ No compile-time verification of mock behavior
- ❌ Repetitive boilerplate for each mock
- ❌ Tests less expressive (can't verify call counts easily)

### After (Generated)
- ✅ 0 lines of manually maintained mock code
- ✅ Mocks auto-update when running `go generate`
- ✅ Compile-time type safety
- ✅ Rich expectations API (Times, Return, Do, etc.)
- ✅ Tests more expressive and readable
- ✅ Industry standard approach

---

## 🔒 .gitignore Consideration

**Option 1: Commit generated mocks (Recommended)**
```gitignore
# Don't ignore mocks - they're part of the build
```

**Why:**
- ✅ CI/CD doesn't need mockgen installed
- ✅ Code reviewers can see mock changes
- ✅ Faster builds

**Option 2: Ignore generated mocks**
```gitignore
# Ignore generated mocks
*/mocks/mock_*.go
```

**Why:**
- ✅ Smaller repo size
- ❌ Requires mockgen in CI/CD
- ❌ Slower CI builds

**Recommendation:** **Commit generated mocks** (Option 1)

---

## 🚀 Example Migration: One File

### Before: `services/flashcard_service_test.go`

```go
// 60 lines of mockRepo implementation
type mockRepo struct{}

func (m *mockRepo) Create(req *models.CreateFlashcardRequest) (*models.Flashcard, error) {
	now := time.Now()
	return &models.Flashcard{ID: 1, Question: req.Question, Answer: req.Answer, CreatedAt: now, UpdatedAt: now}, nil
}

func (m *mockRepo) GetAll() ([]*models.Flashcard, error) {
	now := time.Now()
	return []*models.Flashcard{{ID: 1, Question: "q", Answer: "a", CreatedAt: now, UpdatedAt: now}}, nil
}

// ... 4 more methods

func TestCreateFlashcardValidation(t *testing.T) {
	mockLLM := &MockLLMClient{}
	mockCache := NewMockTranslationCache()
	svc := NewFlashcardService(&mockRepo{}, mockLLM, mockCache)
	// ... test logic
}
```

### After: `services/flashcard_service_test.go`

```go
import (
	dbmocks "github.com/akolybelnikov/flashcards/db/mocks"
	"github.com/akolybelnikov/flashcards/services/mocks"
	"go.uber.org/mock/gomock"
)

func TestCreateFlashcardValidation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := dbmocks.NewMockFlashcardRepository(ctrl)
	mockLLM := mocks.NewMockLLMClient(ctrl)
	mockCache := mocks.NewMockTranslationCache(ctrl)

	// Set up expectations for Create call
	mockRepo.EXPECT().
		Create(gomock.Any()).
		Return(&models.Flashcard{
			ID:       1,
			Question: "hello",
			Answer:   "γεια σας",
		}, nil)

	svc := NewFlashcardService(mockRepo, mockLLM, mockCache)
	// ... test logic
}
```

**Result:** ~10 lines instead of ~70 lines per test file!

---

## 📚 Additional Resources

**gomock Documentation:**
- https://github.com/uber-go/mock
- https://pkg.go.dev/go.uber.org/mock/gomock

**Tutorial:**
- https://blog.golang.org/generate
- https://www.youtube.com/watch?v=ndmB0bj7eyw (Go Testing with Mocks)

---

## ⚡ Quick Start Commands

```bash
# 1. Install mockgen
go install go.uber.org/mock/mockgen@latest

# 2. Add go:generate directives (see Phase 2 above)

# 3. Update Makefile (see Phase 3 above)

# 4. Generate mocks
make generate

# 5. Update tests (see Phase 5 above)

# 6. Run tests
make test

# 7. Delete old mocks
rm services/mocks.go
```

---

## 🎯 Recommendation Summary

**Tool:** `gomock/mockgen` (Uber's fork of golang/mock)  
**Location:** `*/mocks/` folders  
**Command:** `go generate ./...` (via Makefile)  
**Commit mocks:** Yes  
**Effort:** ~2-4 hours for full migration  
**Benefit:** Eliminates ~140 lines of manual mock code ✅

---

**Next Steps:** Review this plan, and when ready, I can implement the migration! 🚀

