package errmsg

import "robinhood/pkg/meta"

var (

	// 1000 - 1999: system error
	InternalServer   = meta.MetaErrorInternalServer.AppendMessage(1000, "The server encountered an internal error or misconfiguration and was unable to complete your request.")
	Forbidden        = meta.MetaErrorForbidden.AppendMessage(1001, "You do not have permission to access this resource.")
	MetaDataNotFound = meta.Error.AppendMessage(1002, "Metadata not found.")

	// 2000 - 2999: user error
	UserNotFound                = meta.Error.AppendMessage(2000, "User not found.")
	UserExisted                 = meta.Error.AppendMessage(2001, "User already existed.")
	UsernameOrPasswordIncorrect = meta.Error.AppendMessage(2002, "Username or Password incorrect.")
	UserRegisterFailed          = meta.Error.AppendMessage(2003, "User register failed.")
	UserLoginFailed             = meta.Error.AppendMessage(2004, "User login failed.")

	// 3000 - 3999: blog error
	BlogNotFound      = meta.Error.AppendMessage(3000, "Blog not found.")
	BlogExisted       = meta.Error.AppendMessage(3001, "Blog already existed.")
	BlogCreateFailed  = meta.Error.AppendMessage(3002, "Blog create failed.")
	BlogUpdateFailed  = meta.Error.AppendMessage(3003, "Blog update failed.")
	BlogArchiveFailed = meta.Error.AppendMessage(3004, "Blog archive failed.")
	BlogInvalidStatus = meta.Error.AppendMessage(3005, "Blog invalid status.")
	BlogGetFailed     = meta.Error.AppendMessage(3006, "Blog get failed.")
	BlogListFailed    = meta.Error.AppendMessage(3007, "Something went wrong. Cannot get blog list.")

	// 4000 - 4999: comment error
	CommentCreateFailed = meta.Error.AppendMessage(4001, "Comment create failed.")
	CommentListFailed   = meta.Error.AppendMessage(4002, "Something went wrong. Cannot get comment list.")
)

func ErrorInvalidRequest(msg string) *meta.MetaError {
	return meta.MetaErrorBadRequest.AppendMessage(1002, msg)
}

func ParseError(c int, desc string) *meta.MetaError {
	return meta.Error.AppendMessage(c, desc)
}
