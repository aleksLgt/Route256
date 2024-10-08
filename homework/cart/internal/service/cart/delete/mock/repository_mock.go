// Code generated by http://github.com/gojuno/minimock (v3.3.11). DO NOT EDIT.

package mock

//go:generate minimock -i route256/cart/internal/service/cart/delete.repository -o repository_mock.go -n RepositoryMock -p mock

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
)

// RepositoryMock implements delete.repository
type RepositoryMock struct {
	t          minimock.Tester
	finishOnce sync.Once

	funcDeleteAll          func(ctx context.Context, userID int64)
	inspectFuncDeleteAll   func(ctx context.Context, userID int64)
	afterDeleteAllCounter  uint64
	beforeDeleteAllCounter uint64
	DeleteAllMock          mRepositoryMockDeleteAll
}

// NewRepositoryMock returns a mock for delete.repository
func NewRepositoryMock(t minimock.Tester) *RepositoryMock {
	m := &RepositoryMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.DeleteAllMock = mRepositoryMockDeleteAll{mock: m}
	m.DeleteAllMock.callArgs = []*RepositoryMockDeleteAllParams{}

	t.Cleanup(m.MinimockFinish)

	return m
}

type mRepositoryMockDeleteAll struct {
	optional           bool
	mock               *RepositoryMock
	defaultExpectation *RepositoryMockDeleteAllExpectation
	expectations       []*RepositoryMockDeleteAllExpectation

	callArgs []*RepositoryMockDeleteAllParams
	mutex    sync.RWMutex

	expectedInvocations uint64
}

// RepositoryMockDeleteAllExpectation specifies expectation struct of the repository.DeleteAll
type RepositoryMockDeleteAllExpectation struct {
	mock      *RepositoryMock
	params    *RepositoryMockDeleteAllParams
	paramPtrs *RepositoryMockDeleteAllParamPtrs

	Counter uint64
}

// RepositoryMockDeleteAllParams contains parameters of the repository.DeleteAll
type RepositoryMockDeleteAllParams struct {
	ctx    context.Context
	userID int64
}

// RepositoryMockDeleteAllParamPtrs contains pointers to parameters of the repository.DeleteAll
type RepositoryMockDeleteAllParamPtrs struct {
	ctx    *context.Context
	userID *int64
}

// Marks this method to be optional. The default behavior of any method with Return() is '1 or more', meaning
// the test will fail minimock's automatic final call check if the mocked method was not called at least once.
// Optional() makes method check to work in '0 or more' mode.
// It is NOT RECOMMENDED to use this option by default unless you really need it, as it helps to
// catch the problems when the expected method call is totally skipped during test run.
func (mmDeleteAll *mRepositoryMockDeleteAll) Optional() *mRepositoryMockDeleteAll {
	mmDeleteAll.optional = true
	return mmDeleteAll
}

// Expect sets up expected params for repository.DeleteAll
func (mmDeleteAll *mRepositoryMockDeleteAll) Expect(ctx context.Context, userID int64) *mRepositoryMockDeleteAll {
	if mmDeleteAll.mock.funcDeleteAll != nil {
		mmDeleteAll.mock.t.Fatalf("RepositoryMock.DeleteAll mock is already set by Set")
	}

	if mmDeleteAll.defaultExpectation == nil {
		mmDeleteAll.defaultExpectation = &RepositoryMockDeleteAllExpectation{}
	}

	if mmDeleteAll.defaultExpectation.paramPtrs != nil {
		mmDeleteAll.mock.t.Fatalf("RepositoryMock.DeleteAll mock is already set by ExpectParams functions")
	}

	mmDeleteAll.defaultExpectation.params = &RepositoryMockDeleteAllParams{ctx, userID}
	for _, e := range mmDeleteAll.expectations {
		if minimock.Equal(e.params, mmDeleteAll.defaultExpectation.params) {
			mmDeleteAll.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmDeleteAll.defaultExpectation.params)
		}
	}

	return mmDeleteAll
}

// ExpectCtxParam1 sets up expected param ctx for repository.DeleteAll
func (mmDeleteAll *mRepositoryMockDeleteAll) ExpectCtxParam1(ctx context.Context) *mRepositoryMockDeleteAll {
	if mmDeleteAll.mock.funcDeleteAll != nil {
		mmDeleteAll.mock.t.Fatalf("RepositoryMock.DeleteAll mock is already set by Set")
	}

	if mmDeleteAll.defaultExpectation == nil {
		mmDeleteAll.defaultExpectation = &RepositoryMockDeleteAllExpectation{}
	}

	if mmDeleteAll.defaultExpectation.params != nil {
		mmDeleteAll.mock.t.Fatalf("RepositoryMock.DeleteAll mock is already set by Expect")
	}

	if mmDeleteAll.defaultExpectation.paramPtrs == nil {
		mmDeleteAll.defaultExpectation.paramPtrs = &RepositoryMockDeleteAllParamPtrs{}
	}
	mmDeleteAll.defaultExpectation.paramPtrs.ctx = &ctx

	return mmDeleteAll
}

// ExpectUserIDParam2 sets up expected param userID for repository.DeleteAll
func (mmDeleteAll *mRepositoryMockDeleteAll) ExpectUserIDParam2(userID int64) *mRepositoryMockDeleteAll {
	if mmDeleteAll.mock.funcDeleteAll != nil {
		mmDeleteAll.mock.t.Fatalf("RepositoryMock.DeleteAll mock is already set by Set")
	}

	if mmDeleteAll.defaultExpectation == nil {
		mmDeleteAll.defaultExpectation = &RepositoryMockDeleteAllExpectation{}
	}

	if mmDeleteAll.defaultExpectation.params != nil {
		mmDeleteAll.mock.t.Fatalf("RepositoryMock.DeleteAll mock is already set by Expect")
	}

	if mmDeleteAll.defaultExpectation.paramPtrs == nil {
		mmDeleteAll.defaultExpectation.paramPtrs = &RepositoryMockDeleteAllParamPtrs{}
	}
	mmDeleteAll.defaultExpectation.paramPtrs.userID = &userID

	return mmDeleteAll
}

// Inspect accepts an inspector function that has same arguments as the repository.DeleteAll
func (mmDeleteAll *mRepositoryMockDeleteAll) Inspect(f func(ctx context.Context, userID int64)) *mRepositoryMockDeleteAll {
	if mmDeleteAll.mock.inspectFuncDeleteAll != nil {
		mmDeleteAll.mock.t.Fatalf("Inspect function is already set for RepositoryMock.DeleteAll")
	}

	mmDeleteAll.mock.inspectFuncDeleteAll = f

	return mmDeleteAll
}

// Return sets up results that will be returned by repository.DeleteAll
func (mmDeleteAll *mRepositoryMockDeleteAll) Return() *RepositoryMock {
	if mmDeleteAll.mock.funcDeleteAll != nil {
		mmDeleteAll.mock.t.Fatalf("RepositoryMock.DeleteAll mock is already set by Set")
	}

	if mmDeleteAll.defaultExpectation == nil {
		mmDeleteAll.defaultExpectation = &RepositoryMockDeleteAllExpectation{mock: mmDeleteAll.mock}
	}

	return mmDeleteAll.mock
}

// Set uses given function f to mock the repository.DeleteAll method
func (mmDeleteAll *mRepositoryMockDeleteAll) Set(f func(ctx context.Context, userID int64)) *RepositoryMock {
	if mmDeleteAll.defaultExpectation != nil {
		mmDeleteAll.mock.t.Fatalf("Default expectation is already set for the repository.DeleteAll method")
	}

	if len(mmDeleteAll.expectations) > 0 {
		mmDeleteAll.mock.t.Fatalf("Some expectations are already set for the repository.DeleteAll method")
	}

	mmDeleteAll.mock.funcDeleteAll = f
	return mmDeleteAll.mock
}

// Times sets number of times repository.DeleteAll should be invoked
func (mmDeleteAll *mRepositoryMockDeleteAll) Times(n uint64) *mRepositoryMockDeleteAll {
	if n == 0 {
		mmDeleteAll.mock.t.Fatalf("Times of RepositoryMock.DeleteAll mock can not be zero")
	}
	mm_atomic.StoreUint64(&mmDeleteAll.expectedInvocations, n)
	return mmDeleteAll
}

func (mmDeleteAll *mRepositoryMockDeleteAll) invocationsDone() bool {
	if len(mmDeleteAll.expectations) == 0 && mmDeleteAll.defaultExpectation == nil && mmDeleteAll.mock.funcDeleteAll == nil {
		return true
	}

	totalInvocations := mm_atomic.LoadUint64(&mmDeleteAll.mock.afterDeleteAllCounter)
	expectedInvocations := mm_atomic.LoadUint64(&mmDeleteAll.expectedInvocations)

	return totalInvocations > 0 && (expectedInvocations == 0 || expectedInvocations == totalInvocations)
}

// DeleteAll implements delete.repository
func (mmDeleteAll *RepositoryMock) DeleteAll(ctx context.Context, userID int64) {
	mm_atomic.AddUint64(&mmDeleteAll.beforeDeleteAllCounter, 1)
	defer mm_atomic.AddUint64(&mmDeleteAll.afterDeleteAllCounter, 1)

	if mmDeleteAll.inspectFuncDeleteAll != nil {
		mmDeleteAll.inspectFuncDeleteAll(ctx, userID)
	}

	mm_params := RepositoryMockDeleteAllParams{ctx, userID}

	// Record call args
	mmDeleteAll.DeleteAllMock.mutex.Lock()
	mmDeleteAll.DeleteAllMock.callArgs = append(mmDeleteAll.DeleteAllMock.callArgs, &mm_params)
	mmDeleteAll.DeleteAllMock.mutex.Unlock()

	for _, e := range mmDeleteAll.DeleteAllMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return
		}
	}

	if mmDeleteAll.DeleteAllMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmDeleteAll.DeleteAllMock.defaultExpectation.Counter, 1)
		mm_want := mmDeleteAll.DeleteAllMock.defaultExpectation.params
		mm_want_ptrs := mmDeleteAll.DeleteAllMock.defaultExpectation.paramPtrs

		mm_got := RepositoryMockDeleteAllParams{ctx, userID}

		if mm_want_ptrs != nil {

			if mm_want_ptrs.ctx != nil && !minimock.Equal(*mm_want_ptrs.ctx, mm_got.ctx) {
				mmDeleteAll.t.Errorf("RepositoryMock.DeleteAll got unexpected parameter ctx, want: %#v, got: %#v%s\n", *mm_want_ptrs.ctx, mm_got.ctx, minimock.Diff(*mm_want_ptrs.ctx, mm_got.ctx))
			}

			if mm_want_ptrs.userID != nil && !minimock.Equal(*mm_want_ptrs.userID, mm_got.userID) {
				mmDeleteAll.t.Errorf("RepositoryMock.DeleteAll got unexpected parameter userID, want: %#v, got: %#v%s\n", *mm_want_ptrs.userID, mm_got.userID, minimock.Diff(*mm_want_ptrs.userID, mm_got.userID))
			}

		} else if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmDeleteAll.t.Errorf("RepositoryMock.DeleteAll got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		return

	}
	if mmDeleteAll.funcDeleteAll != nil {
		mmDeleteAll.funcDeleteAll(ctx, userID)
		return
	}
	mmDeleteAll.t.Fatalf("Unexpected call to RepositoryMock.DeleteAll. %v %v", ctx, userID)

}

// DeleteAllAfterCounter returns a count of finished RepositoryMock.DeleteAll invocations
func (mmDeleteAll *RepositoryMock) DeleteAllAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmDeleteAll.afterDeleteAllCounter)
}

// DeleteAllBeforeCounter returns a count of RepositoryMock.DeleteAll invocations
func (mmDeleteAll *RepositoryMock) DeleteAllBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmDeleteAll.beforeDeleteAllCounter)
}

// Calls returns a list of arguments used in each call to RepositoryMock.DeleteAll.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmDeleteAll *mRepositoryMockDeleteAll) Calls() []*RepositoryMockDeleteAllParams {
	mmDeleteAll.mutex.RLock()

	argCopy := make([]*RepositoryMockDeleteAllParams, len(mmDeleteAll.callArgs))
	copy(argCopy, mmDeleteAll.callArgs)

	mmDeleteAll.mutex.RUnlock()

	return argCopy
}

// MinimockDeleteAllDone returns true if the count of the DeleteAll invocations corresponds
// the number of defined expectations
func (m *RepositoryMock) MinimockDeleteAllDone() bool {
	if m.DeleteAllMock.optional {
		// Optional methods provide '0 or more' call count restriction.
		return true
	}

	for _, e := range m.DeleteAllMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	return m.DeleteAllMock.invocationsDone()
}

// MinimockDeleteAllInspect logs each unmet expectation
func (m *RepositoryMock) MinimockDeleteAllInspect() {
	for _, e := range m.DeleteAllMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to RepositoryMock.DeleteAll with params: %#v", *e.params)
		}
	}

	afterDeleteAllCounter := mm_atomic.LoadUint64(&m.afterDeleteAllCounter)
	// if default expectation was set then invocations count should be greater than zero
	if m.DeleteAllMock.defaultExpectation != nil && afterDeleteAllCounter < 1 {
		if m.DeleteAllMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to RepositoryMock.DeleteAll")
		} else {
			m.t.Errorf("Expected call to RepositoryMock.DeleteAll with params: %#v", *m.DeleteAllMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcDeleteAll != nil && afterDeleteAllCounter < 1 {
		m.t.Error("Expected call to RepositoryMock.DeleteAll")
	}

	if !m.DeleteAllMock.invocationsDone() && afterDeleteAllCounter > 0 {
		m.t.Errorf("Expected %d calls to RepositoryMock.DeleteAll but found %d calls",
			mm_atomic.LoadUint64(&m.DeleteAllMock.expectedInvocations), afterDeleteAllCounter)
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *RepositoryMock) MinimockFinish() {
	m.finishOnce.Do(func() {
		if !m.minimockDone() {
			m.MinimockDeleteAllInspect()
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
		m.MinimockDeleteAllDone()
}
