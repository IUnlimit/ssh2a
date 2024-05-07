package pam

/*
#include <errno.h>
#include <pwd.h>
#include <security/pam_appl.h>
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <unistd.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <netdb.h>
#include<arpa/inet.h>
#include <sys/prctl.h>

#define RECVBUF 1024

// create a http request with tcp
// hostname support ip or domain
char *make_request(char *hostname, int port, char *message) {
    int sockfd;
      char reply[RECVBUF];
    struct sockaddr_in server;
    struct in_addr addr;
    struct hostent *h;

    memset(reply, '\0', sizeof(reply));

    sockfd = socket(AF_INET, SOCK_STREAM, 0);
    if (sockfd == -1) {
        return NULL;
    }

    if(! inet_pton(AF_INET, hostname, &addr)){
        if((h = gethostbyname(hostname)) == NULL) {
             return NULL;
        }
        hostname = inet_ntoa(*((struct in_addr *)h->h_addr));
    }

    server.sin_addr.s_addr = inet_addr(hostname);
    server.sin_family = AF_INET;
    server.sin_port = htons(port);

    if (connect(sockfd, (struct sockaddr *)&server, sizeof(server)) < 0) {
        return NULL;
    }

    // TODO: set send timeout
    if (send(sockfd, message, strlen(message), 0) < 0) {
        return NULL;
    }

    // TODO: set recv timeout
    if (recv(sockfd, reply, sizeof(reply), 0) < 0) {
        return NULL;
    }

    close(sockfd);
    return strdup(reply);
}
*/
import "C"

import (
	"unsafe"
)

func doNetSend(host string, port int, msg string) string {
	cHostName := C.CString(host)
	defer C.free(unsafe.Pointer(cHostName))

	cMsg := C.CString(msg)
	defer C.free(unsafe.Pointer(cMsg))

	resp := C.make_request(cHostName, C.int(port), cMsg)
	defer C.free(unsafe.Pointer(resp))

	return C.GoString(resp)
}
