package vo

type HireStatus int

const (
	_           HireStatus = iota
	NewHire                // 试用期
	Regularized            // 转正
	Suspended              // 停职
	Regisn                 // 离职
	Retired                // 退休
	Dismissed              // 辞退
)
