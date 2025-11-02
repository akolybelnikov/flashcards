# Mock Generation - Quick Visual Guide

## 🎯 The Two Major Options

### 1. **gomock/mockgen** ⭐ RECOMMENDED
```
Industry Standard | Google/Uber | 9.5k+ stars
✅ Auto-generate from interfaces
✅ Type-safe
✅ Rich expectations
✅ go generate support
```

### 2. **moq**
```
Simpler Alternative | 2k+ stars
✅ Auto-generate
✅ Type-safe
⚠️ Less features
✅ go generate support
```

### 3. **testify/mock** ❌ Not recommended for our case
```
Popular but Manual | 23k+ stars
❌ NOT auto-generated
⚠️ Handwritten with helpers
```

---

## 📊 Current vs Future

### CURRENT: Handwritten Mocks (~140 lines manual code)

```
flashcards/
├── services/
│   ├── flashcard_service_test.go
│   │   └── mockRepo struct { }      ← 60 lines manually written
│   │       func (m *mockRepo) Create() { ... }
│   │       func (m *mockRepo) GetAll() { ... }
│   │       func (m *mockRepo) GetByID() { ... }
│   │       ... 3 more methods
│   │
│   └── mocks.go                      ← 80 lines manually written
│       ├── MockLLMClient
│       └── MockTranslationCache
│
└── handlers/
    └── flashcard_handler_test.go
        └── mockService struct { }    ← Inline, manual
```

**Problems:**
- ❌ Must update mocks when interfaces change
- ❌ Easy to forget
- ❌ Repetitive boilerplate
- ❌ No compile-time verification

---

### FUTURE: Generated Mocks (0 lines manual code!)

```
flashcards/
├── services/
│   ├── flashcard_service.go
│   │   //go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE
│   │
│   ├── llm_client.go
│   │   //go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE
│   │
│   ├── translation_cache.go
│   │   //go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE
│   │
│   ├── flashcard_service_test.go     ← Clean, just tests
│   │   import "github.com/.../services/mocks"
│   │   mock := mocks.NewMockLLMClient(ctrl)
│   │   mock.EXPECT().Translate(...).Return(...)
│   │
│   └── mocks/                         ← AUTO-GENERATED
│       ├── mock_flashcard_service.go  (200 lines, auto)
│       ├── mock_llm_client.go         (100 lines, auto)
│       └── mock_translation_cache.go  (150 lines, auto)
│
├── db/
│   ├── flashcard_repository.go
│   │   //go:generate mockgen -source=$GOFILE -destination=mocks/mock_$GOFILE
│   │
│   └── mocks/
│       └── mock_flashcard_repository.go
│
└── Makefile
    generate:
        go generate ./...    ← One command regenerates all!
```

**Benefits:**
- ✅ Run `make generate` to update all mocks
- ✅ Compile-time type checking
- ✅ Can't forget to update
- ✅ Rich expectations API

---

## 🔄 Migration Steps (Visual)

```
STEP 1: Install
┌─────────────────────────────────────┐
│ go install go.uber.org/mock/mockgen │
└─────────────────────────────────────┘

STEP 2: Add Directives
┌─────────────────────────────────────────────────────────────┐
│ // flashcard_service.go                                     │
│ //go:generate mockgen -source=$GOFILE -destination=...     │
│ package services                                            │
└─────────────────────────────────────────────────────────────┘

STEP 3: Update Makefile
┌─────────────────────────────────────┐
│ generate:                           │
│     go generate ./...               │
│                                     │
│ test: generate                      │
│     go test ./... -v                │
└─────────────────────────────────────┘

STEP 4: Generate
┌─────────────────────────────────────┐
│ $ make generate                     │
│ ✓ Generated services/mocks/*.go     │
│ ✓ Generated db/mocks/*.go           │
└─────────────────────────────────────┘

STEP 5: Update Tests
┌──────────────────────────────────────────────────────┐
│ // Before                                            │
│ type mockRepo struct{}                               │
│ func (m *mockRepo) Create(...) { ... }  ← Manual    │
│                                                      │
│ // After                                             │
│ ctrl := gomock.NewController(t)                     │
│ mock := mocks.NewMockFlashcardRepository(ctrl)      │
│ mock.EXPECT().Create(...).Return(...)  ← Generated │
└──────────────────────────────────────────────────────┘

STEP 6: Delete Old Mocks
┌─────────────────────────────────────┐
│ rm services/mocks.go                │
│ # Remove inline mockRepo            │
│ # Remove inline mockService         │
└─────────────────────────────────────┘

STEP 7: Test
┌─────────────────────────────────────┐
│ $ make test                         │
│ ✓ All tests pass                    │
└─────────────────────────────────────┘
```

---

## 💡 Example: Before & After

### BEFORE (Handwritten - 60+ lines)

```go
// services/flashcard_service_test.go

type mockRepo struct{}

func (m *mockRepo) Create(req *models.CreateFlashcardRequest) (*models.Flashcard, error) {
	now := time.Now()
	return &models.Flashcard{
		ID:       1,
		Question: req.Question,
		Answer:   req.Answer,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (m *mockRepo) GetAll() ([]*models.Flashcard, error) {
	// ... 15 lines
}

func (m *mockRepo) GetByID(id int) (*models.Flashcard, error) {
	// ... 10 lines
}

// ... 3 more methods, 35+ more lines

func TestCreateFlashcard(t *testing.T) {
	repo := &mockRepo{}  // Simple, but limited
	svc := NewFlashcardService(repo, nil, nil)
	// Can't verify Create was called with specific args!
}
```

### AFTER (Generated - Clean!)

```go
// services/flashcard_service_test.go

import (
	"github.com/akolybelnikov/flashcards/db/mocks"
	"go.uber.org/mock/gomock"
)

func TestCreateFlashcard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock
	repo := mocks.NewMockFlashcardRepository(ctrl)

	// Set expectations (rich API!)
	repo.EXPECT().
		Create(gomock.Any()).
		Times(1).                    // Called exactly once ✅
		Return(&models.Flashcard{
			ID:       1,
			Question: "hello",
			Answer:   "γεια σας",
		}, nil)

	// Test
	svc := NewFlashcardService(repo, nil, nil)
	fc, err := svc.CreateFlashcard(&models.CreateFlashcardRequest{
		Question: "hello",
		Answer:   "γεια σας",
	})

	// If Create wasn't called, test FAILS automatically! ✅
}
```

---

## 🎨 Project Structure Comparison

### Current
```
services/
├── flashcard_service.go
├── flashcard_service_test.go  (tests + 60 lines mock)
├── llm_client.go
├── translation_cache.go
└── mocks.go                    (80 lines manual mocks)

Total manual mock code: 140 lines ❌
```

### After gomock
```
services/
├── flashcard_service.go        //go:generate ...
├── flashcard_service_test.go   (just tests, clean!)
├── llm_client.go                //go:generate ...
├── translation_cache.go         //go:generate ...
└── mocks/                       ← New folder
    ├── mock_flashcard_service.go   (auto-generated)
    ├── mock_llm_client.go          (auto-generated)
    └── mock_translation_cache.go   (auto-generated)

Total manual mock code: 0 lines ✅
```

---

## 🚀 Quick Decision Matrix

**Use gomock if:**
- ✅ You have interface-based architecture (we do!)
- ✅ You want type-safety (we do!)
- ✅ You want rich expectations (Times, Return, Do)
- ✅ You want industry standard
- ✅ You want zero manual mock code

**Use moq if:**
- ✅ You want simpler approach
- ✅ You don't need rich expectations
- ⚠️ Less community support

**Use testify/mock if:**
- ❌ You want to write mocks manually (we don't!)

---

## 📋 Files That Change

### Files with //go:generate added:
1. ✅ `services/flashcard_service.go`
2. ✅ `services/llm_client.go`
3. ✅ `services/translation_cache.go`
4. ✅ `db/flashcard_repository.go`

### Files to update tests:
1. ✅ `services/flashcard_service_test.go`
2. ✅ `handlers/flashcard_handler_test.go`

### Files to delete:
1. ❌ `services/mocks.go`

### Files to create:
1. ✅ `Makefile` (add `generate` target)
2. ✅ `services/mocks/` (folder, auto-populated)
3. ✅ `db/mocks/` (folder, auto-populated)

---

## ⏱️ Effort Estimate

- **Setup mockgen:** 5 minutes
- **Add //go:generate:** 10 minutes
- **Update Makefile:** 5 minutes
- **Generate mocks:** 1 minute
- **Update service tests:** 30 minutes
- **Update handler tests:** 30 minutes
- **Delete old mocks:** 5 minutes
- **Testing & verification:** 15 minutes

**Total:** ~2 hours

---

## 🎯 Summary

**Recommendation:** Use **gomock/mockgen**

**Why:**
1. Industry standard (Google → Uber)
2. Zero manual mock code
3. Type-safe, compile-time checked
4. Rich expectations API
5. One command regenerates all

**Next:** Review plan, then implement! 🚀

