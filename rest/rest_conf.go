/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package main

import "github.com/pvillela/gfoa/examples/transpmtmock/travelsvc/internal/boot"

var config = boot.TravelConfig{
	CfgForValidateTripRequestBf:  "CfgForValidateTripRequestBf",
	CfgForGetPostedCardStateDaf:  "CfgForGetPostedCardStateDaf",
	CfgForGetUnpostedUsagesDaf:   "CfgForGetUnpostedUsagesDaf",
	CfgForUpdateCardStateBf:      "CfgForUpdateCardStateBf",
	CfgForRateTripSc:             "CfgForRateTripSc",
	CfgForPrepareUsageBf:         "CfgForPrepareUsageBf",
	CfgForWriteUsageDaf:          "CfgForWriteUsageDaf",
	CfgForPrepareDownstreamEvtBf: "CfgForPrepareDownstreamEvtBf",
	CfgForDownstreamEp:           "CfgForDownstreamEp",
	CfgForPrepareResponseBf:      "CfgForPrepareResponseBf",
}
