package alias

import "github.com/reusee/domui"

var (
	// from https://developer.mozilla.org/en-US/docs/Web/Events

	EError        = domui.MakeEventFunc("error")
	EAbort        = domui.MakeEventFunc("abort")
	ELoad         = domui.MakeEventFunc("load")
	EBeforeUnload = domui.MakeEventFunc("beforeunload")
	EUnload       = domui.MakeEventFunc("unload")

	EOnline  = domui.MakeEventFunc("online")
	EOffline = domui.MakeEventFunc("offline")

	EFocus    = domui.MakeEventFunc("focus")
	EBlur     = domui.MakeEventFunc("blur")
	EFocusIn  = domui.MakeEventFunc("focusin")
	EFocusOut = domui.MakeEventFunc("focusout")

	EOpen    = domui.MakeEventFunc("open")
	EMessage = domui.MakeEventFunc("message")
	EClose   = domui.MakeEventFunc("close")

	EPageHide = domui.MakeEventFunc("pagehide")
	EPageShow = domui.MakeEventFunc("pageshow")
	EPopState = domui.MakeEventFunc("popstate")

	EAnimationStart     = domui.MakeEventFunc("animationstart")
	EAnimationCancel    = domui.MakeEventFunc("animationcancel")
	EAnimationEnd       = domui.MakeEventFunc("animationend")
	EAnimationIteration = domui.MakeEventFunc("animationiteration")

	ETransitionStart  = domui.MakeEventFunc("transitionstart")
	ETransitionCancel = domui.MakeEventFunc("transitioncancel")
	ETransitionEnd    = domui.MakeEventFunc("transitionend")
	ETransitionRun    = domui.MakeEventFunc("transitionrun")

	EReset  = domui.MakeEventFunc("reset")
	ESubmit = domui.MakeEventFunc("submit")

	EBeforePrint = domui.MakeEventFunc("beforeprint")
	EAfterPrint  = domui.MakeEventFunc("afterprint")

	ECompositionStart  = domui.MakeEventFunc("compositionstart")
	ECompositionUpdate = domui.MakeEventFunc("compositionupdate")
	ECompositionEnd    = domui.MakeEventFunc("compositionend")

	EFullscreenChange = domui.MakeEventFunc("fullscreenchange")
	EFullscreenError  = domui.MakeEventFunc("fullscreenerror")
	EResize           = domui.MakeEventFunc("resize")
	EScroll           = domui.MakeEventFunc("scroll")

	ECut   = domui.MakeEventFunc("cut")
	ECopy  = domui.MakeEventFunc("copy")
	EPaste = domui.MakeEventFunc("paste")

	EKeyDown  = domui.MakeEventFunc("keydown")
	EKeyPress = domui.MakeEventFunc("keypress")
	EKeyUp    = domui.MakeEventFunc("keyup")

	EAuxClick          = domui.MakeEventFunc("auxclick")
	EClick             = domui.MakeEventFunc("click")
	EContextMenu       = domui.MakeEventFunc("contextmenu")
	EDblClick          = domui.MakeEventFunc("dblclick")
	EMouseDown         = domui.MakeEventFunc("mousedown")
	EMouseEnter        = domui.MakeEventFunc("mouseenter")
	EMouseLeave        = domui.MakeEventFunc("mouseleave")
	EMouseMove         = domui.MakeEventFunc("mousemove")
	EMouseOver         = domui.MakeEventFunc("mouseover")
	EMouseOut          = domui.MakeEventFunc("mouseout")
	EMouseUp           = domui.MakeEventFunc("mouseup")
	EPointerLockChange = domui.MakeEventFunc("pointerlockchange")
	EPointerLockError  = domui.MakeEventFunc("pointerlockerror")
	ESelect            = domui.MakeEventFunc("select")
	EWheel             = domui.MakeEventFunc("wheel")

	EDrag      = domui.MakeEventFunc("drag")
	EDragEnd   = domui.MakeEventFunc("dragend")
	EDragEnter = domui.MakeEventFunc("dragenter")
	EDragStart = domui.MakeEventFunc("dragstart")
	EDragLeave = domui.MakeEventFunc("dragleave")
	EDragOver  = domui.MakeEventFunc("dragover")
	EDrop      = domui.MakeEventFunc("drop")

	EAudioProcess   = domui.MakeEventFunc("audioprocess")
	ECanPlay        = domui.MakeEventFunc("canplay")
	ECanPlayThrough = domui.MakeEventFunc("canplaythrough")
	EComplete       = domui.MakeEventFunc("complete")
	EDurationChange = domui.MakeEventFunc("durationchange")
	EEmptied        = domui.MakeEventFunc("emptied")
	EEnded          = domui.MakeEventFunc("ended")
	ELoadedData     = domui.MakeEventFunc("loadeddata")
	ELoadedMetaData = domui.MakeEventFunc("loadedmetadata")
	EPause          = domui.MakeEventFunc("pause")
	EPlay           = domui.MakeEventFunc("play")
	EPlaying        = domui.MakeEventFunc("playing")
	ERateChange     = domui.MakeEventFunc("ratechange")
	ESeeked         = domui.MakeEventFunc("seeked")
	ESeeking        = domui.MakeEventFunc("seeking")
	EStalled        = domui.MakeEventFunc("stalled")
	ESuspend        = domui.MakeEventFunc("suspend")
	ETimeUpdate     = domui.MakeEventFunc("timeupdate")
	EVolumeChange   = domui.MakeEventFunc("volumechange")
	EWaiting        = domui.MakeEventFunc("waiting")

	ELoadEnd   = domui.MakeEventFunc("loadend")
	ELoadStart = domui.MakeEventFunc("loadstart")
	EProgress  = domui.MakeEventFunc("progress")
	ETimeout   = domui.MakeEventFunc("timeout")

	EChange  = domui.MakeEventFunc("change")
	EStorage = domui.MakeEventFunc("storage")

	EChecking    = domui.MakeEventFunc("checking")
	EDownloading = domui.MakeEventFunc("downloading")
	ENoUpdate    = domui.MakeEventFunc("noupdate")
	EObsolete    = domui.MakeEventFunc("obsolete")
	EUpdateReady = domui.MakeEventFunc("updateready")

	EBroadcast           = domui.MakeEventFunc("broadcast")
	ECheckboxStateChange = domui.MakeEventFunc("checkboxstatechange")
	EHashChange          = domui.MakeEventFunc("hashchange")
	EInput               = domui.MakeEventFunc("input")
	ERadioStateChange    = domui.MakeEventFunc("radiostatechange")
	EReadyStateChange    = domui.MakeEventFunc("readystatechange")
	EValueChange         = domui.MakeEventFunc("valuechange")

	EInvalid = domui.MakeEventFunc("invalid")
	EShow    = domui.MakeEventFunc("show")
)
