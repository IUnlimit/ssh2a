package pam

/*
#cgo LDFLAGS: -lpam
#include <stdlib.h>
#include <security/pam_appl.h>

int do_conv(pam_handle_t* hdlr, int count, const struct pam_message** msgs, struct pam_response** responses) {
    int err;
    struct pam_conv* conv;

    err = pam_get_item(hdlr, PAM_CONV, (const void**)&conv);
    if(err != PAM_SUCCESS) {
        return err;
    }

    return conv->conv(count, msgs, responses, conv->appdata_ptr);
}
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// 封装实现了Conversation函数，以实现用户与ssh会话等交互消息

// MessageStyle is a style of Message
type MessageStyle int

const (
	// MessageEchoOff is for messages that shouldn't gave an echo.
	MessageEchoOff = MessageStyle(C.PAM_PROMPT_ECHO_OFF)

	// MessageEchoOn is for messages that should have an echo.
	MessageEchoOn = MessageStyle(C.PAM_PROMPT_ECHO_ON)

	// MessageErrorMsg is for messages that should be displayed as an error.
	MessageErrorMsg = MessageStyle(C.PAM_ERROR_MSG)

	// MessageTextInfo is for textual blurbs to be spat out.
	MessageTextInfo = MessageStyle(C.PAM_TEXT_INFO)
)

// Message represents something to ask / show in a Conv.Conversation call.
type Message struct {
	Style MessageStyle
	Msg   string
}

// Handle is a handle type to hang the PAM methods off of.
type Handle struct {
	Ptr unsafe.Pointer
}

func (hdl Handle) ptr() *C.pam_handle_t {
	return (*C.pam_handle_t)(hdl.Ptr)
}

// Conversation passes on the specified messages.
func (hdl Handle) Conversation(_msg Message) (string, error) {
	msgStruct := (*C.struct_pam_message)(C.malloc(C.sizeof_struct_pam_message))
	msgStruct.msg_style = C.int(_msg.Style)
	msgStruct.msg = C.CString(_msg.Msg)
	defer C.free(unsafe.Pointer(msgStruct.msg))
	defer C.free(unsafe.Pointer(msgStruct))

	respStruct := C.malloc(C.sizeof_struct_pam_response)
	defer C.free(respStruct)

	msg := (*C.struct_pam_message)(unsafe.Pointer(msgStruct))
	resp := (*C.struct_pam_response)(respStruct)

	code := C.do_conv(hdl.ptr(), C.int(1), &msg, &resp)
	if code != C.PAM_SUCCESS {
		return "", fmt.Errorf("got non-success from the function: %d", code)
	}

	ret := C.GoString(resp.resp)
	C.free(unsafe.Pointer(resp.resp))

	return ret, nil
}
