package testtcell

import (
	"github.com/gdamore/tcell"
	"github.com/stretchr/testify/mock"
)

//MockScreen mocks the calls to the Screen interface.
type MockScreen struct {
	mock.Mock
}

//Init mocks the operation of the same name from the Screen interface.
func (mock *MockScreen) Init() error { mock.Called(); return nil }

//Fini mocks the operation of the same name from the Screen interface.
func (mock *MockScreen) Fini() { mock.Called() }

//Clear mocks the operation of the same name from the Screen interface.
func (mock *MockScreen) Clear() { mock.Called() }

//Fill mocks the operation of the same name from the Screen interface.
func (mock *MockScreen) Fill(r rune, s tcell.Style) { mock.Called(r, s) }

//SetCell mocks the operation of the same name from the Screen interface.
func (mock *MockScreen) SetCell(x int, y int, style tcell.Style, ch ...rune) {
	mock.Called(x, y, style, ch)
}

//GetContent mocks the operation of the same name from the Screen interface.
func (mock *MockScreen) GetContent(x, y int) (mainc rune, combc []rune, style tcell.Style, width int) {
	args := mock.Called(x, y)
	return args.Get(0).(rune), args.Get(1).([]rune), args.Get(2).(tcell.Style), args.Int(3)
}

//SetContent mocks the operation of the same name from the Screen interface.
func (mock *MockScreen) SetContent(x int, y int, mainc rune, combc []rune, style tcell.Style) {
	mock.Called(x, y, mainc, combc, style)
}

//SetStyle mocks the operation of the same name from the Screen interface.
func (mock *MockScreen) SetStyle(style tcell.Style) { mock.Called(style) }

//ShowCursor mocks the operation of the same name from the Screen interface.
func (mock *MockScreen) ShowCursor(x int, y int) {}

//HideCursor mocks the operation of the same name from the Screen interface.
func (mock *MockScreen) HideCursor() {}

//Size mocks the operation of the same name from the Screen interface.
func (mock *MockScreen) Size() (int, int) { args := mock.Called(); return args.Int(0), args.Int(1) }

//PollEvent mocks the operation of the same name from the Screen interface.
func (mock *MockScreen) PollEvent() tcell.Event {
	args := mock.Called()
	return args.Get(0).(tcell.Event)
}

//PostEvent mocks the operation of the same name from the Screen interface.
func (mock *MockScreen) PostEvent(ev tcell.Event) error {
	args := mock.Called(ev)
	return args.Get(0).(error)
}

//PostEventWait mocks the operation of the same name from the Screen interface.
func (mock *MockScreen) PostEventWait(ev tcell.Event) {}

//EnableMouse mocks the operation of the same name from the Screen interface.
func (mock *MockScreen) EnableMouse() {}

//DisableMouse mocks the operation of the same name from the Screen interface.
func (mock *MockScreen) DisableMouse() {}

//HasMouse mocks the operation of the same name from the Screen interface.
func (mock *MockScreen) HasMouse() bool { args := mock.Called(); return args.Bool(0) }

//Colors mocks the operation of the same name from the Screen interface.
func (mock *MockScreen) Colors() int { args := mock.Called(); return args.Int(0) }

//Show mocks the operation of the same name from the Screen interface.
func (mock *MockScreen) Show() { mock.Called() }

//Sync mocks the operation of the same name from the Screen interface.
func (mock *MockScreen) Sync() { mock.Called() }

//CharacterSet mocks the operation of the same name from the Screen interface.
func (mock *MockScreen) CharacterSet() string { args := mock.Called(); return args.String(0) }

//RegisterRuneFallback mocks the operation of the same name from the Screen interface.
func (mock *MockScreen) RegisterRuneFallback(r rune, subst string) { mock.Called(r, subst) }

//UnregisterRuneFallback mocks the operation of the same name from the Screen interface.
func (mock *MockScreen) UnregisterRuneFallback(r rune) { mock.Called(r) }

//CanDisplay mocks the operation of the same name from the Screen interface.
func (mock *MockScreen) CanDisplay(r rune, checkFallbacks bool) bool {
	args := mock.Called(r, checkFallbacks)
	return args.Bool(0)
}

//Resize mocks the operation of the same name from the Screen interface.
func (mock *MockScreen) Resize(a, b, c, d int) { mock.Called(a, b, c, d) }

//HasKey mocks the operation of the same name from the Screen interface.
func (mock *MockScreen) HasKey(key tcell.Key) bool {
	args := mock.Called(key)
	return args.Bool(0)
}
