package xthrift

const (
	T_STOP        = 0
	T_VOID        = 1
	T_BOOL        = 2
	T_BYTE        = 3
	T_I08         = 3
	T_DOUBLE      = 4
	T_I16         = 6
	T_I32         = 8
	T_I64         = 10
	T_STRING      = 11
	T_UTF7        = 11
	T_STRUCT      = 12
	T_MAP         = 13
	T_SET         = 14
	T_LIST        = 15
	T_UTF8        = 16
	T_UTF16  int8 = 17
)

const (
	T_CALL      = 1
	T_REPLY     = 2
	T_EXCEPTION = 3
	T_ONEWAY    = 4
)
