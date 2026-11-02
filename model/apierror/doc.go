// Package apierror is the passive error-identity model plus
// the two wire codes GlobalExceptionHandler renders:
//
//	private static final String ERROR_CODE = "FAT_0202";           // catch-all -> 400
//	private static final String ERROR_CODE_FORBIDDEN = "FAT_0114"; // AccessDeniedException -> 403
//
// It graduates to internal/model because the permission-miss identity is
// a widely-detectable error: a caller several layers below the HTTP
// handler (e.g. pkg/auth, wiring its required permission.Permission.Code
// at CMD Init, or a logic unit that wants to raise the exact same
// identity a connector would) needs to be able to construct/detect it
// with errors.Is/errors.As WITHOUT importing a logic layer or the
// centralized Gin error-mapping engine itself.
//
// This package lands the passive identity/model values (AccessDenied/
// Generic/NotFound/Conflict and their Err* sentinels) plus Classify — the
// pure mapping function that internal/connector/restserver/endpoint's own
// respondError helper calls for every error returned by a logic/connector
// dependency. An *APIError keeps its own status; anything else becomes a
// generic 500 so raw driver or internal messages never reach the client.
// Classify stays here, not in the connector, because it is still passive
// data transformation (error in, *APIError out) — only the actual
// Gin-facing adaptation (rendering an *APIError through  libs restserver's
// RespondError envelope, which writes onto a *gin.Context) is genuine
// behavior, and that lives in the dispatcher-connector's endpoint package
// instead (see that package's own error.go/doc.go).
package apierror
