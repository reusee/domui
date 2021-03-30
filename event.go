package domui

import (
	"sync"
	"sync/atomic"
	"syscall/js"

	"github.com/reusee/dscope"
)

var elementID int32 = 42

type EventSpec struct {
	Event string
	Func  any
}

func (_ EventSpec) IsSpec() {}

func MakeEventFunc(event string) func(fn func()) EventSpec {
	return func(fn func()) EventSpec {
		return On(event, fn)
	}
}

func On(ev string, fn any) EventSpec {
	return EventSpec{
		Event: ev,
		Func:  fn,
	}
}

var (
	// from https://developer.mozilla.org/en-US/docs/Web/Events

	OnError        = MakeEventFunc("error")
	OnAbort        = MakeEventFunc("abort")
	OnLoad         = MakeEventFunc("load")
	OnBeforeUnload = MakeEventFunc("beforeunload")
	OnUnload       = MakeEventFunc("unload")

	OnOnline  = MakeEventFunc("online")
	OnOffline = MakeEventFunc("offline")

	OnFocus    = MakeEventFunc("focus")
	OnBlur     = MakeEventFunc("blur")
	OnFocusIn  = MakeEventFunc("focusin")
	OnFocusOut = MakeEventFunc("focusout")

	OnOpen    = MakeEventFunc("open")
	OnMessage = MakeEventFunc("message")
	OnClose   = MakeEventFunc("close")

	OnPageHide = MakeEventFunc("pagehide")
	OnPageShow = MakeEventFunc("pageshow")
	OnPopState = MakeEventFunc("popstate")

	OnAnimationStart     = MakeEventFunc("animationstart")
	OnAnimationCancel    = MakeEventFunc("animationcancel")
	OnAnimationEnd       = MakeEventFunc("animationend")
	OnAnimationIteration = MakeEventFunc("animationiteration")

	OnTransitionStart  = MakeEventFunc("transitionstart")
	OnTransitionCancel = MakeEventFunc("transitioncancel")
	OnTransitionEnd    = MakeEventFunc("transitionend")
	OnTransitionRun    = MakeEventFunc("transitionrun")

	OnReset  = MakeEventFunc("reset")
	OnSubmit = MakeEventFunc("submit")

	OnBeforePrint = MakeEventFunc("beforeprint")
	OnAfterPrint  = MakeEventFunc("afterprint")

	OnCompositionStart  = MakeEventFunc("compositionstart")
	OnCompositionUpdate = MakeEventFunc("compositionupdate")
	OnCompositionEnd    = MakeEventFunc("compositionend")

	OnFullscreenChange = MakeEventFunc("fullscreenchange")
	OnFullscreenError  = MakeEventFunc("fullscreenerror")
	OnResize           = MakeEventFunc("resize")
	OnScroll           = MakeEventFunc("scroll")

	OnCut   = MakeEventFunc("cut")
	OnCopy  = MakeEventFunc("copy")
	OnPaste = MakeEventFunc("paste")

	OnKeyDown  = MakeEventFunc("keydown")
	OnKeyPress = MakeEventFunc("keypress")
	OnKeyUp    = MakeEventFunc("keyup")

	OnAuxClick          = MakeEventFunc("auxclick")
	OnClick             = MakeEventFunc("click")
	OnContextMenu       = MakeEventFunc("contextmenu")
	OnDblClick          = MakeEventFunc("dblclick")
	OnMouseDown         = MakeEventFunc("mousedown")
	OnMouseEnter        = MakeEventFunc("mouseenter")
	OnMouseLeave        = MakeEventFunc("mouseleave")
	OnMouseMove         = MakeEventFunc("mousemove")
	OnMouseOver         = MakeEventFunc("mouseover")
	OnMouseOut          = MakeEventFunc("mouseout")
	OnMouseUp           = MakeEventFunc("mouseup")
	OnPointerLockChange = MakeEventFunc("pointerlockchange")
	OnPointerLockError  = MakeEventFunc("pointerlockerror")
	OnSelect            = MakeEventFunc("select")
	OnWheel             = MakeEventFunc("wheel")

	OnDrag      = MakeEventFunc("drag")
	OnDragEnd   = MakeEventFunc("dragend")
	OnDragEnter = MakeEventFunc("dragenter")
	OnDragStart = MakeEventFunc("dragstart")
	OnDragLeave = MakeEventFunc("dragleave")
	OnDragOver  = MakeEventFunc("dragover")
	OnDrop      = MakeEventFunc("drop")

	OnAudioProcess   = MakeEventFunc("audioprocess")
	OnCanPlay        = MakeEventFunc("canplay")
	OnCanPlayThrough = MakeEventFunc("canplaythrough")
	OnComplete       = MakeEventFunc("complete")
	OnDurationChange = MakeEventFunc("durationchange")
	OnEmptied        = MakeEventFunc("emptied")
	OnEnded          = MakeEventFunc("ended")
	OnLoadedData     = MakeEventFunc("loadeddata")
	OnLoadedMetaData = MakeEventFunc("loadedmetadata")
	OnPause          = MakeEventFunc("pause")
	OnPlay           = MakeEventFunc("play")
	OnPlaying        = MakeEventFunc("playing")
	OnRateChange     = MakeEventFunc("ratechange")
	OnSeeked         = MakeEventFunc("seeked")
	OnSeeking        = MakeEventFunc("seeking")
	OnStalled        = MakeEventFunc("stalled")
	OnSuspend        = MakeEventFunc("suspend")
	OnTimeUpdate     = MakeEventFunc("timeupdate")
	OnVolumeChange   = MakeEventFunc("volumechange")
	OnWaiting        = MakeEventFunc("waiting")

	OnLoadEnd   = MakeEventFunc("loadend")
	OnLoadStart = MakeEventFunc("loadstart")
	OnProgress  = MakeEventFunc("progress")
	OnTimeout   = MakeEventFunc("timeout")

	OnChange  = MakeEventFunc("change")
	OnStorage = MakeEventFunc("storage")

	OnChecking    = MakeEventFunc("checking")
	OnDownloading = MakeEventFunc("downloading")
	OnNoUpdate    = MakeEventFunc("noupdate")
	OnObsolete    = MakeEventFunc("obsolete")
	OnUpdateReady = MakeEventFunc("updateready")

	OnBroadcast           = MakeEventFunc("broadcast")
	OnCheckboxStateChange = MakeEventFunc("checkboxstatechange")
	OnHashChange          = MakeEventFunc("hashchange")
	OnInput               = MakeEventFunc("input")
	OnRadioStateChange    = MakeEventFunc("radiostatechange")
	OnReadyStateChange    = MakeEventFunc("readystatechange")
	OnValueChange         = MakeEventFunc("valuechange")

	OnInvalid = MakeEventFunc("invalid")
	OnShow    = MakeEventFunc("show")
)

var (
	eventRegistryLock sync.RWMutex
	eventRegistry     = make(map[int32]map[string][]EventSpec)
	eventHandlerSet   = make(map[string]bool)
)

var eventHandlerScope = dscope.New()

func setEventSpecs(wrap js.Value, element js.Value, specs map[string][]EventSpec) {

	idValue := element.Get("__element_id__")
	var id int32
	if idValue.IsUndefined() {
		id = atomic.AddInt32(&elementID, 1)
		element.Set("__element_id__", id)
	} else {
		id = int32(idValue.Int())
	}

	for event := range specs {
		if eventHandlerSet[event] {
			continue
		}
		wrap.Call(
			"addEventListener",
			event,
			js.FuncOf(
				func(this js.Value, args []js.Value) any {
					go func() {
						ev := args[0]
						typ := ev.Get("type").String()
						bubbles := ev.Get("bubbles").Bool()
						for node := ev.Get("target"); !node.IsNull() && !node.IsUndefined() && !node.Equal(wrap); node = node.Get("parentNode") {
							idValue := node.Get("__element_id__")
							if idValue.IsUndefined() {
								if !bubbles {
									break
								}
								continue
							}
							id := int32(idValue.Int())
							eventRegistryLock.RLock()
							var specs []EventSpec
							if evs, ok := eventRegistry[id]; ok {
								if ss, ok := evs[typ]; ok {
									specs = append(ss[:0:0], ss...)
								}
							}
							eventRegistryLock.RUnlock()
							for _, spec := range specs {
								eventHandlerScope.Sub(
									func() js.Value {
										return node
									},
									func(node js.Value) AttrChecked {
										return AttrChecked(
											node.Get("checked").Bool(),
										)
									},
								).Call(spec.Func)
							}
							if !bubbles {
								break
							}
						}
					}()
					return nil
				},
			),
			true,
		)
		eventHandlerSet[event] = true
	}

	eventRegistryLock.Lock()
	defer eventRegistryLock.Unlock()
	eventRegistry[id] = specs

}

func unsetEventSpecs(element js.Value) {
	idValue := element.Get("__element_id__")
	var id int32
	if idValue.IsUndefined() {
		return
	} else {
		id = int32(idValue.Int())
	}
	eventRegistryLock.Lock()
	defer eventRegistryLock.Unlock()
	eventRegistry[id] = nil
	childNodes := element.Get("childNodes")
	for i := childNodes.Length() - 1; i >= 0; i-- {
		unsetEventSpecs(childNodes.Index(i))
	}
}
