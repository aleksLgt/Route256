// Code generated by http://github.com/gojuno/minimock (v3.3.11). DO NOT EDIT.

package mock

//go:generate minimock -i route256/cart/internal/service/cart/item/add.repository -o repository_mock.go -n RepositoryMock -p mock

import (
	"context"
	"route256/cart/internal/domain"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
)

// RepositoryMock implements add.repository
type RepositoryMock struct {
	t          minimock.Tester
	finishOnce sync.Once

	funcAdd          func(ctx context.Context, userID int64, item domain.Item)
	inspectFuncAdd   func(ctx context.Context, userID int64, item domain.Item)
	afterAddCounter  uint64
	beforeAddCounter uint64
	AddMock          mRepositoryMockAdd
}

// NewRepositoryMock returns a mock for add.repository
func NewRepositoryMock(t minimock.Tester) *RepositoryMock {
	m := &RepositoryMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AddMock = mRepositoryMockAdd{mock: m}
	m.AddMock.callArgs = []*RepositoryMockAddParams{}

	t.Cleanup(m.MinimockFinish)

	return m
}

type mRepositoryMockAdd struct {
	optional           bool
	mock               *RepositoryMock
	defaultExpectation *RepositoryMockAddExpectation
	expectations       []*RepositoryMockAddExpectation

	callArgs []*RepositoryMockAddParams
	mutex    sync.RWMutex

	expectedInvocations uint64
}

// RepositoryMockAddExpectation specifies expectation struct of the repository.Add
type RepositoryMockAddExpectation struct {
	mock      *RepositoryMock
	params    *RepositoryMockAddParams
	paramPtrs *RepositoryMockAddParamPtrs

	Counter uint64
}

// RepositoryMockAddParams contains parameters of the repository.Add
type RepositoryMockAddParams struct {
	ctx    context.Context
	userID int64
	item   domain.Item
}

// RepositoryMockAddParamPtrs contains pointers to parameters of the repository.Add
type RepositoryMockAddParamPtrs struct {
	ctx    *context.Context
	userID *int64
	item   *domain.Item
}

// Marks this method to be optional. The default behavior of any method with Return() is '1 or more', meaning
// the test will fail minimock's automatic final call check if the mocked method was not called at least once.
// Optional() makes method check to work in '0 or more' mode.
// It is NOT RECOMMENDED to use this option by default unless you really need it, as it helps to
// catch the problems when the expected method call is totally skipped during test run.
func (mmAdd *mRepositoryMockAdd) Optional() *mRepositoryMockAdd {
	mmAdd.optional = true
	return mmAdd
}

// Expect sets up expected params for repository.Add
func (mmAdd *mRepositoryMockAdd) Expect(ctx context.Context, userID int64, item domain.Item) *mRepositoryMockAdd {
	if mmAdd.mock.funcAdd != nil {
		mmAdd.mock.t.Fatalf("RepositoryMock.Add mock is already set by Set")
	}

	if mmAdd.defaultExpectation == nil {
		mmAdd.defaultExpectation = &RepositoryMockAddExpectation{}
	}

	if mmAdd.defaultExpectation.paramPtrs != nil {
		mmAdd.mock.t.Fatalf("RepositoryMock.Add mock is already set by ExpectParams functions")
	}

	mmAdd.defaultExpectation.params = &RepositoryMockAddParams{ctx, userID, item}
	for _, e := range mmAdd.expectations {
		if minimock.Equal(e.params, mmAdd.defaultExpectation.params) {
			mmAdd.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmAdd.defaultExpectation.params)
		}
	}

	return mmAdd
}

// ExpectCtxParam1 sets up expected param ctx for repository.Add
func (mmAdd *mRepositoryMockAdd) ExpectCtxParam1(ctx context.Context) *mRepositoryMockAdd {
	if mmAdd.mock.funcAdd != nil {
		mmAdd.mock.t.Fatalf("RepositoryMock.Add mock is already set by Set")
	}

	if mmAdd.defaultExpectation == nil {
		mmAdd.defaultExpectation = &RepositoryMockAddExpectation{}
	}

	if mmAdd.defaultExpectation.params != nil {
		mmAdd.mock.t.Fatalf("RepositoryMock.Add mock is already set by Expect")
	}

	if mmAdd.defaultExpectation.paramPtrs == nil {
		mmAdd.defaultExpectation.paramPtrs = &RepositoryMockAddParamPtrs{}
	}
	mmAdd.defaultExpectation.paramPtrs.ctx = &ctx

	return mmAdd
}

// ExpectUserIDParam2 sets up expected param userID for repository.Add
func (mmAdd *mRepositoryMockAdd) ExpectUserIDParam2(userID int64) *mRepositoryMockAdd {
	if mmAdd.mock.funcAdd != nil {
		mmAdd.mock.t.Fatalf("RepositoryMock.Add mock is already set by Set")
	}

	if mmAdd.defaultExpectation == nil {
		mmAdd.defaultExpectation = &RepositoryMockAddExpectation{}
	}

	if mmAdd.defaultExpectation.params != nil {
		mmAdd.mock.t.Fatalf("RepositoryMock.Add mock is already set by Expect")
	}

	if mmAdd.defaultExpectation.paramPtrs == nil {
		mmAdd.defaultExpectation.paramPtrs = &RepositoryMockAddParamPtrs{}
	}
	mmAdd.defaultExpectation.paramPtrs.userID = &userID

	return mmAdd
}

// ExpectItemParam3 sets up expected param item for repository.Add
func (mmAdd *mRepositoryMockAdd) ExpectItemParam3(item domain.Item) *mRepositoryMockAdd {
	if mmAdd.mock.funcAdd != nil {
		mmAdd.mock.t.Fatalf("RepositoryMock.Add mock is already set by Set")
	}

	if mmAdd.defaultExpectation == nil {
		mmAdd.defaultExpectation = &RepositoryMockAddExpectation{}
	}

	if mmAdd.defaultExpectation.params != nil {
		mmAdd.mock.t.Fatalf("RepositoryMock.Add mock is already set by Expect")
	}

	if mmAdd.defaultExpectation.paramPtrs == nil {
		mmAdd.defaultExpectation.paramPtrs = &RepositoryMockAddParamPtrs{}
	}
	mmAdd.defaultExpectation.paramPtrs.item = &item

	return mmAdd
}

// Inspect accepts an inspector function that has same arguments as the repository.Add
func (mmAdd *mRepositoryMockAdd) Inspect(f func(ctx context.Context, userID int64, item domain.Item)) *mRepositoryMockAdd {
	if mmAdd.mock.inspectFuncAdd != nil {
		mmAdd.mock.t.Fatalf("Inspect function is already set for RepositoryMock.Add")
	}

	mmAdd.mock.inspectFuncAdd = f

	return mmAdd
}

// Return sets up results that will be returned by repository.Add
func (mmAdd *mRepositoryMockAdd) Return() *RepositoryMock {
	if mmAdd.mock.funcAdd != nil {
		mmAdd.mock.t.Fatalf("RepositoryMock.Add mock is already set by Set")
	}

	if mmAdd.defaultExpectation == nil {
		mmAdd.defaultExpectation = &RepositoryMockAddExpectation{mock: mmAdd.mock}
	}

	return mmAdd.mock
}

// Set uses given function f to mock the repository.Add method
func (mmAdd *mRepositoryMockAdd) Set(f func(ctx context.Context, userID int64, item domain.Item)) *RepositoryMock {
	if mmAdd.defaultExpectation != nil {
		mmAdd.mock.t.Fatalf("Default expectation is already set for the repository.Add method")
	}

	if len(mmAdd.expectations) > 0 {
		mmAdd.mock.t.Fatalf("Some expectations are already set for the repository.Add method")
	}

	mmAdd.mock.funcAdd = f
	return mmAdd.mock
}

// Times sets number of times repository.Add should be invoked
func (mmAdd *mRepositoryMockAdd) Times(n uint64) *mRepositoryMockAdd {
	if n == 0 {
		mmAdd.mock.t.Fatalf("Times of RepositoryMock.Add mock can not be zero")
	}
	mm_atomic.StoreUint64(&mmAdd.expectedInvocations, n)
	return mmAdd
}

func (mmAdd *mRepositoryMockAdd) invocationsDone() bool {
	if len(mmAdd.expectations) == 0 && mmAdd.defaultExpectation == nil && mmAdd.mock.funcAdd == nil {
		return true
	}

	totalInvocations := mm_atomic.LoadUint64(&mmAdd.mock.afterAddCounter)
	expectedInvocations := mm_atomic.LoadUint64(&mmAdd.expectedInvocations)

	return totalInvocations > 0 && (expectedInvocations == 0 || expectedInvocations == totalInvocations)
}

// Add implements add.repository
func (mmAdd *RepositoryMock) Add(ctx context.Context, userID int64, item domain.Item) {
	mm_atomic.AddUint64(&mmAdd.beforeAddCounter, 1)
	defer mm_atomic.AddUint64(&mmAdd.afterAddCounter, 1)

	if mmAdd.inspectFuncAdd != nil {
		mmAdd.inspectFuncAdd(ctx, userID, item)
	}

	mm_params := RepositoryMockAddParams{ctx, userID, item}

	// Record call args
	mmAdd.AddMock.mutex.Lock()
	mmAdd.AddMock.callArgs = append(mmAdd.AddMock.callArgs, &mm_params)
	mmAdd.AddMock.mutex.Unlock()

	for _, e := range mmAdd.AddMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return
		}
	}

	if mmAdd.AddMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmAdd.AddMock.defaultExpectation.Counter, 1)
		mm_want := mmAdd.AddMock.defaultExpectation.params
		mm_want_ptrs := mmAdd.AddMock.defaultExpectation.paramPtrs

		mm_got := RepositoryMockAddParams{ctx, userID, item}

		if mm_want_ptrs != nil {

			if mm_want_ptrs.ctx != nil && !minimock.Equal(*mm_want_ptrs.ctx, mm_got.ctx) {
				mmAdd.t.Errorf("RepositoryMock.Add got unexpected parameter ctx, want: %#v, got: %#v%s\n", *mm_want_ptrs.ctx, mm_got.ctx, minimock.Diff(*mm_want_ptrs.ctx, mm_got.ctx))
			}

			if mm_want_ptrs.userID != nil && !minimock.Equal(*mm_want_ptrs.userID, mm_got.userID) {
				mmAdd.t.Errorf("RepositoryMock.Add got unexpected parameter userID, want: %#v, got: %#v%s\n", *mm_want_ptrs.userID, mm_got.userID, minimock.Diff(*mm_want_ptrs.userID, mm_got.userID))
			}

			if mm_want_ptrs.item != nil && !minimock.Equal(*mm_want_ptrs.item, mm_got.item) {
				mmAdd.t.Errorf("RepositoryMock.Add got unexpected parameter item, want: %#v, got: %#v%s\n", *mm_want_ptrs.item, mm_got.item, minimock.Diff(*mm_want_ptrs.item, mm_got.item))
			}

		} else if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmAdd.t.Errorf("RepositoryMock.Add got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		return

	}
	if mmAdd.funcAdd != nil {
		mmAdd.funcAdd(ctx, userID, item)
		return
	}
	mmAdd.t.Fatalf("Unexpected call to RepositoryMock.Add. %v %v %v", ctx, userID, item)

}

// AddAfterCounter returns a count of finished RepositoryMock.Add invocations
func (mmAdd *RepositoryMock) AddAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmAdd.afterAddCounter)
}

// AddBeforeCounter returns a count of RepositoryMock.Add invocations
func (mmAdd *RepositoryMock) AddBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmAdd.beforeAddCounter)
}

// Calls returns a list of arguments used in each call to RepositoryMock.Add.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmAdd *mRepositoryMockAdd) Calls() []*RepositoryMockAddParams {
	mmAdd.mutex.RLock()

	argCopy := make([]*RepositoryMockAddParams, len(mmAdd.callArgs))
	copy(argCopy, mmAdd.callArgs)

	mmAdd.mutex.RUnlock()

	return argCopy
}

// MinimockAddDone returns true if the count of the Add invocations corresponds
// the number of defined expectations
func (m *RepositoryMock) MinimockAddDone() bool {
	if m.AddMock.optional {
		// Optional methods provide '0 or more' call count restriction.
		return true
	}

	for _, e := range m.AddMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	return m.AddMock.invocationsDone()
}

// MinimockAddInspect logs each unmet expectation
func (m *RepositoryMock) MinimockAddInspect() {
	for _, e := range m.AddMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to RepositoryMock.Add with params: %#v", *e.params)
		}
	}

	afterAddCounter := mm_atomic.LoadUint64(&m.afterAddCounter)
	// if default expectation was set then invocations count should be greater than zero
	if m.AddMock.defaultExpectation != nil && afterAddCounter < 1 {
		if m.AddMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to RepositoryMock.Add")
		} else {
			m.t.Errorf("Expected call to RepositoryMock.Add with params: %#v", *m.AddMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcAdd != nil && afterAddCounter < 1 {
		m.t.Error("Expected call to RepositoryMock.Add")
	}

	if !m.AddMock.invocationsDone() && afterAddCounter > 0 {
		m.t.Errorf("Expected %d calls to RepositoryMock.Add but found %d calls",
			mm_atomic.LoadUint64(&m.AddMock.expectedInvocations), afterAddCounter)
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *RepositoryMock) MinimockFinish() {
	m.finishOnce.Do(func() {
		if !m.minimockDone() {
			m.MinimockAddInspect()
			m.t.FailNow()
		}
	})
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *RepositoryMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *RepositoryMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockAddDone()
}
