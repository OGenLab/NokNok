package errcode

// Common error code
const (
	errcode_class____system_error = -10000 /////////////////////////
	SUCCESS                       = 0
	SYSTEM_ERROR                  = -10000 // unreachable code, maybe a bug.
	UNREACHABLE_CODE              = -10000 // unreachable code, maybe a bug.
	UNKNOWN_ERROR                 = -10000 // unreachable code, maybe a bug.
	DATABASE_ERROR                = -10001 // database error.
	REDIS_ERROR                   = -10002 // redis error.

	errcode_class___request_error = -20000 /////////////////////////
	PARAMETER_INVALID             = -20000 // request parameter error
	APPID_INVALID                 = -20001 // request appid invalid
	USERID_MISMATCH               = -20002 // user not match
	DATA_MISMATCH                 = -20003 // request data mismatch
	API_NOT_FOUND                 = -20404 // api not found
	INNTER_API_REQUEST            = -20405

	errcode_class__haverisk_error = 20000 //////////////////////////
	NEED_CAPTCHA                  = 20001 // need captcha
	LOGIN_FREEZEN                 = 20002 // login freezen by too many errors

	errcode_class___frontend_jump = 10000 //////////////////////////
	NOT_LOGIN                     = 10400 // user not login
	NEED_PAY_PASSWORD             = 10401 // need input pay password
	NEED_GATK                     = 10402
	GATK_ERROR                    = 10403

	// usermanager boost
	MAX_BOOST_ERR         = 50101
	UPDATE_NEXT_BOOST_ERR = 50102

	// battle hammergopher
	ROUND_QUERY_ERR      = 50111
	NOT_HAMMER_GOPHER    = 50112
	HAMMER_GOPHER_UPDATE = 50113

	// battle prizedraw
	NOT_DRAW_COUNT = 50121

	// receive award
	UNCOMPLTE_TASK_ERR  = 50131
	RECEIVED_REWARD_ERR = 50132
)

// API error code
var (
	ErrMsg = map[int]string{
		SUCCESS: "success",

		PARAMETER_INVALID: "input params error",
		DATABASE_ERROR:    "database error",
		NOT_LOGIN:         "not login",

		MAX_BOOST_ERR:         "max boost",
		UPDATE_NEXT_BOOST_ERR: "upgade fail",

		NOT_HAMMER_GOPHER: "not hit gopher",
		NOT_DRAW_COUNT:    "not draw prize count",
	}
)
