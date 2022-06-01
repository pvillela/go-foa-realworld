/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package memdb

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/types"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"github.com/pvillela/go-foa-realworld/internal/platform/db.postgres/daf"
	"time"
)

// MDb defines an in-memory view of the database to support testing.
type MDb struct {
	users        mUsersT
	recCtxUsers  mRecCtxUsersT
	articlesPlus mArticlesPlusT
	favorites    mFavoritesT
	followings   mFollowingsT
	comments     mCommentsT
	tags         mTagsT
}

///////////////////
// Constructor

func New() MDb {
	return MDb{
		users:        mUsersT{},
		recCtxUsers:  mRecCtxUsersT{},
		articlesPlus: mArticlesPlusT{},
		favorites:    mFavoritesT{},
		followings:   mFollowingsT{},
		comments:     mCommentsT{},
		tags:         mTagsT{},
	}
}

///////////////////
// Methods

func (mdb MDb) UserGet(username string) model.User {
	return mdb.users[username]
}

func (mdb MDb) UserGet2(username string) (model.User, daf.RecCtxUser) {
	return mdb.users[username], mdb.recCtxUsers[username]
}

func (mdb MDb) UserGetAll() ([]model.User, []daf.RecCtxUser) {
	users := make([]model.User, len(mdb.users))
	recCtxs := make([]daf.RecCtxUser, len(mdb.users))
	i := 0
	for _, v := range mdb.users {
		users[i] = v
		i++
	}
	i = 0
	for _, v := range mdb.recCtxUsers {
		recCtxs[i] = v
		i++
	}
	return users, recCtxs
}

func (mdb *MDb) UserUpsert(username string, user model.User, recCtx daf.RecCtxUser) {
	mdb.users.upsert(username, user)
	mdb.recCtxUsers.upsert(username, recCtx)
}

func (mdb MDb) ArticlePlusGet(slug string) model.ArticlePlus {
	return mdb.articlesPlus[slug]
}

func (mdb MDb) ArticlePlusGetAll() []model.ArticlePlus {
	result := make([]model.ArticlePlus, len(mdb.articlesPlus))
	i := 0
	for _, v := range mdb.articlesPlus {
		result[i] = v
		i++
	}
	return result
}

func (mdb *MDb) ArticlePlusUpsert(
	article model.Article,
	favorite bool,
	user model.User,
	follows bool,
) {
	mdb.articlesPlus.upsert(article, favorite, user, follows)
}

func (mdb MDb) FavoriteGet(username string, slug string) bool {
	return mdb.favorites[username][slug]
}

func (mdb MDb) FollowingGet(followerName string, followeeName string) model.Following {
	return mdb.followings[followerName][followeeName]
}

func (mdb *MDb) FollowingUpsert(followerName string, followeeName string, followedOn time.Time) {
	follower := mdb.users[followerName]
	followee := mdb.users[followeeName]
	following := model.Following{
		FollowerID: follower.Id,
		FolloweeID: followee.Id,
		FollowedOn: followedOn,
	}
	mdb.followings[followerName][followeeName] = following
}

func (mdb MDb) CommentGet(username string, slug string, id uint) model.Comment {
	commentKey := mCommentKeyT{username: username, slug: slug}
	return mdb.comments[commentKey][id]
}

func (mdb MDb) CommentGetAllForKey(username string, slug string) []model.Comment {
	commentKey := mCommentKeyT{username: username, slug: slug}
	comments := make([]model.Comment, len(mdb.comments))
	i := 0
	for _, comment := range mdb.comments[commentKey] {
		comments[i] = comment
		i++
	}
	return comments
}

func (mdb MDb) CommentGetAllBySlug(slug string) []model.Comment {
	returnedComments := []model.Comment{}
	for k, v := range mdb.comments {
		if k.slug == slug {
			var comments []model.Comment
			for _, comment := range v {
				comments = append(comments, comment)
			}
			returnedComments = append(returnedComments, comments...)
		}
	}
	return returnedComments
}

func (mdb MDb) CommentGetAll() []model.Comment {
	allComments := []model.Comment{}
	for k, _ := range mdb.comments {
		comments := mdb.CommentGetAllForKey(k.username, k.slug)
		allComments = append(allComments, comments...)
	}
	return allComments
}

func (mdb *MDb) CommentInsert(username string, slug string, comment model.Comment) {
	commentKey := mCommentKeyT{username: username, slug: slug}
	if mdb.comments[commentKey] == nil {
		mdb.comments[commentKey] = make(map[uint]model.Comment)
	}
	mdb.comments[commentKey][comment.Id] = comment
}

func (mdb *MDb) CommentDelete(username string, slug string, id uint) {
	outerMap := mdb.comments
	outerKey := mCommentKeyT{username, slug}
	innerMap := outerMap[outerKey]
	delete(innerMap, id)
	if len(innerMap) == 0 {
		delete(outerMap, outerKey)
	}
}

func (mdb *MDb) CommentDeleteAll() {
	mdb.comments = make(mCommentsT)
}

///////////////////
// Supporting types

// key is Username
type mUsersT map[string]model.User

func (m *mUsersT) upsert(username string, user model.User) {
	(*m)[username] = user
}

// key is Username
type mRecCtxUsersT map[string]daf.RecCtxUser

func (m *mRecCtxUsersT) upsert(username string, recCtx daf.RecCtxUser) {
	(*m)[username] = recCtx
}

// key is Slug
type mArticlesPlusT map[string]model.ArticlePlus

func (m mArticlesPlusT) upsert(
	article model.Article,
	favorite bool,
	user model.User,
	follows bool,
) {
	slug := article.Slug
	m[slug] = model.ArticlePlus_FromArticle(article, favorite, model.Profile_FromUser(user, follows))
}

// key is Username, value is a map from Slug to bool
type mFavoritesT map[string]map[string]bool

func (m *mFavoritesT) upsert(username string, slug string) {
	(*m)[username][slug] = true
}

// key is follower.Username, value is a map from followee.Usesrname to model.Following
type mFollowingsT map[string]map[string]model.Following

type mCommentKeyT struct {
	username string
	slug     string
}

// key is an mCommentKeyT, value is a map from comment.Id to model.Comment
type mCommentsT map[mCommentKeyT]map[uint]model.Comment

type mTagsT map[string]types.Unit
