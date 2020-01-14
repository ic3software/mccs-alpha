package util

import (
	"github.com/ic3network/mccs-alpha/global/constant"
)

// IsAcceptedStatus checks if the business status is accpeted.
func IsAcceptedStatus(status string) bool {
	if status == constant.Business.Accepted ||
		status == constant.Trading.Pending ||
		status == constant.Trading.Accepted ||
		status == constant.Trading.Rejected {
		return true
	}
	return false
}
