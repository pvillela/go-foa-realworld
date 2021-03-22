/*
 * Copyright © 2021 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package main

import (
	"errors"

	"github.com/pvillela/gfoa/examples/transpmtmock/travelsvc/internal/boot"
	"github.com/pvillela/gfoa/examples/transpmtmock/travelsvc/pkg/model"
	"github.com/pvillela/gfoa/examples/transpmtmock/travelsvc/pkg/rpc"
	log "github.com/sirupsen/logrus"
)

add_comment_sfl.go
authenticate_sfl.go
create_article_sfl.go
delete_article_sfl.go
delete_comment_sfl.go
favorite_article_sfl.go
feed_articles.go
follow_user_sfl.go
get_article.go
get_comments_sfl.go
get_current_user_sfl.go
get_profile_sfl.go
get_tags_sfl.go
list_articles_sfl.go
register_sfl.go
unfavorite_article_sfl.go
unfollow_user_sfl.go
update_article_sfl.go
update_user_sfl.go

var tripSvc = boot.TripSvcflowboot(config)

func tripSvcflowP(pInput interface{}) (interface{}, error) {
	return tripSvc(*pInput.(*rpc.TripRequest)), nil
}

func tripSvcflowM(m map[string]string) (interface{}, error) {
	cardInfo, cardInfoOk := m["cardInfo"]
	deviceInfo, deviceInfoOk := m["deviceInfo"]

	log.Info("m[cardInfo]", m["cardInfo"])
	log.Info("m[deviceInfo]", m["deviceInfo"])

	errMsg := ""
	if !cardInfoOk {
		errMsg = "cardInfo parameter not found"
	}
	if !deviceInfoOk {
		if errMsg != "" {
			errMsg = errMsg + ", "
		}
		errMsg = errMsg + "deviceInfo parameter not found"
	}

	var err error
	if errMsg != "" {
		err = errors.New(errMsg)
		return rpc.TripResponse{}, err
	}

	input := rpc.TripRequest{
		CardInfo:   model.CardInfo(cardInfo),
		DeviceInfo: model.DeviceInfo(deviceInfo),
	}

	return tripSvc(input), err
}
