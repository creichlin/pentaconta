package services
// Signal table
var signals = map[string]int{
	"hangup": 1,
	"interrupt": 2,
	"quit": 3,
	"illegal instruction": 4,
	"trace/breakpoint trap": 5,
	"aborted": 6,
	"bus error": 7,
	"floating point exception": 8,
	"killed": 9,
	"user defined signal 1": 10,
	"segmentation fault": 11,
	"user defined signal 2": 12,
	"broken pipe": 13,
	"alarm clock": 14,
	"terminated": 15,
	"stack fault": 16,
	"child exited": 17,
	"continued": 18,
	"stopped (signal)": 19,
	"stopped": 20,
	"stopped (tty input)": 21,
	"stopped (tty output)": 22,
	"urgent I/O condition": 23,
	"CPU time limit exceeded": 24,
	"file size limit exceeded": 25,
	"virtual timer expired": 26,
	"profiling timer expired": 27,
	"window changed": 28,
	"I/O possible": 29,
	"power failure": 30,
	"bad system call": 31,
}
