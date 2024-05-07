package pam

/*
#cgo LDFLAGS: -lpam -fPIC
#include <security/pam_appl.h>
#include <stdlib.h>

char *get_user(pam_handle_t *pamh);
int get_uid(char *user);
*/
import "C"
import "unsafe"

// https://fossies.org/dox/openpam-20190224/pam__sm__authenticate_8c.html
//
//export pam_sm_authenticate
func pam_sm_authenticate(pamh *C.pam_handle_t, flags, argc C.int, argv **C.char) C.int {
	cUsername := C.get_user(pamh)
	if cUsername == nil {
		return C.PAM_USER_UNKNOWN
	}
	defer C.free(unsafe.Pointer(cUsername))

	uid := int(C.get_uid(cUsername))
	if uid < 0 {
		return C.PAM_USER_UNKNOWN
	}

	hdl := Handle{unsafe.Pointer(pamh)}
	user, _ := hdl.Conversation(
		Message{
			Style: MessageEchoOn,
			Msg:   "Welcome to Server, Username: ",
		},
	)

	pwd, _ := hdl.Conversation(
		Message{
			Style: MessageEchoOff,
			Msg:   "Input Your Password: ",
		},
	)

	// 此处mypamAuthenticate函数，在业务层实现相关认证逻辑
	ok, err := mypamAuthenticate(uid, C.GoString(cUsername), user, pwd)
	if err != nil {
		pamLog("authenticate err: %v", err)
		return C.PAM_AUTH_ERR
	}

	if !ok {
		pamLog("authenticate failed(user:%s, pwd:%s)", idmUser, idmPwd)
		return C.PAM_AUTH_ERR
	}

	return C.PAM_SUCCESS
}

// https://fossies.org/dox/openpam-20190224/pam__sm__setcred_8c.html
//
//export pam_sm_setcred
func pam_sm_setcred(pamh *C.pam_handle_t, flags, argc C.int, argv **C.char) C.int {
	return C.PAM_SUCCESS
}

// https://fossies.org/dox/openpam-20190224/pam__sm__acct__mgmt_8c.html
//
//export pam_sm_acct_mgmt
func pam_sm_acct_mgmt(pamh *C.pam_handle_t, flags, argc C.int, argv **C.char) C.int {
	return C.PAM_SUCCESS
}

// https://fossies.org/dox/openpam-20190224/pam__sm__open__session_8c.html
//
//export pam_sm_open_session
func pam_sm_open_session(pamh *C.pam_handle_t, flags, argc C.int, argv **C.char) C.int {
	return C.PAM_SUCCESS
}

// https://fossies.org/dox/openpam-20190224/pam__sm__close__session_8c.html
//
//export pam_sm_close_session
func pam_sm_close_session(pamh *C.pam_handle_t, flags, argc C.int, argv **C.char) C.int {
	return C.PAM_SUCCESS
}

// https://fossies.org/dox/openpam-20190224/pam__get__authtok_8c.html
//
//export pam_sm_chauthtok
func pam_sm_chauthtok(pamh *C.pam_handle_t, flags, argc C.int, argv **C.char) C.int {
	return C.PAM_SUCCESS
}
