/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package memdb

import (
	"github.com/pvillela/go-foa-realworld/internal/arch/util"
	"github.com/pvillela/go-foa-realworld/internal/model"
	"time"
)

// MDb defines an in-memory view of the database to support testing.
type MDb struct {
	usersByName mUsersByNameT
	usersById   mUsersByIdT
	articles    mArticlesT
	favorites   mFavoritesT
	followings  mFollowingsT
	comments    mCommentsT
	tags        mTagsT
}

///////////////////
// Constructor

func New() MDb {
	return MDb{
		usersByName: mUsersByNameT{},
		usersById:   mUsersByIdT{},
		articles:    mArticlesT{},
		favorites:   mFavoritesT{},
		followings:  mFollowingsT{},
		comments:    mCommentsT{},
		tags:        mTagsT{},
	}
}

///////////////////
// Methods

func (mdb MDb) UserGetByName(username string) model.User {
	return *mdb.usersByName[username]
}

func (mdb MDb) userGetById(id uint) model.User {
	return *mdb.usersById[id]
}

func (mdb MDb) UserGetAll() []model.User {
	userPs := util.MapToSlice(mdb.usersByName)
	users := util.SliceMap(userPs, func(puser *model.User) model.User {
		return *puser
	})
	return users
}

func (mdb *MDb) UserUpsert(user model.User) {
	existingUserByName := mdb.usersByName[user.Username]
	if existingUserByName != nil && existingUserByName.Id != user.Id {
		panic("attempt to clobber existing username in mdb")
	}
	userById := mdb.usersById[user.Id]
	if userById == nil {
		userById = &model.User{}
		mdb.usersById[user.Id] = userById
	}
	currentUsername := userById.Username
	if currentUsername != user.Username {
		delete(mdb.usersByName, currentUsername)
	}
	*userById = user
	mdb.usersByName[user.Username] = userById
}

func (mdb MDb) ArticleGetBySlug(slug string) model.Article {
	article := mdb.articles[slug]
	if article == nil {
		panic("attempt to get article for invalid slug " + slug)
	}
	return *article
}

func (mdb MDb) ArticleGetAll() []model.Article {
	var articlePs = util.MapToSlice(mdb.articles)
	return util.SliceMap(articlePs, func(pa *model.Article) model.Article {
		return *pa
	})
}

func (mdb *MDb) ArticleUpsert(
	article model.Article,
) {
	mdb.articles.upsert(article)
}

func (mdb MDb) ArticlePlusGet(currUsername string, slug string) model.ArticlePlus {
	article := mdb.ArticleGetBySlug(slug)
	author := mdb.userGetById(article.AuthorId)
	favorited := mdb.Favorited(currUsername, slug)
	follows := mdb.Follows(currUsername, author.Username)
	return model.ArticlePlus_FromArticle(article, favorited, model.Profile_FromUser(author, follows))
}

func (mdb MDb) ArticlePlusGetAll(currUsername string) []model.ArticlePlus {
	result := make([]model.ArticlePlus, len(mdb.articles))
	i := 0
	for _, a := range mdb.articles {
		result[i] = mdb.ArticlePlusGet(currUsername, a.Slug)
		i++
	}
	return result
}

func (mdb MDb) CommentGet(username string, slug string, id uint) model.Comment {
	commentKey := mCommentKeyT{username: username, slug: slug}
	return mdb.comments[commentKey][id]
}

func (mdb MDb) CommentGetAllForKey(username string, slug string) []model.Comment {
	commentKey := mCommentKeyT{username: username, slug: slug}
	return util.MapToSlice(mdb.comments[commentKey])
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

func (mdb MDb) Favorited(username string, slug string) bool {
	return mdb.favorites[mFavoriteT{username, slug}]
}

func (mdb *MDb) FavoritePut(username string, slug string) {
	mdb.favorites.put(username, slug)
}

func (mdb *MDb) FavoritedDelete(username string, slug string) {
	delete(mdb.favorites, mFavoriteT{username, slug})
}

func (mdb MDb) FollowingGet(followerName string, followeeName string) model.Following {
	return mdb.followings[mFollowingT{followerName, followeeName}]
}

func (mdb MDb) Follows(followerName string, followeeName string) bool {
	return mdb.followings[mFollowingT{followerName, followeeName}] != model.Following{}
}

func (mdb MDb) FollowingDelete(followerName string, followeeName string) {
	delete(mdb.followings, mFollowingT{followerName, followeeName})
}

func (mdb *MDb) FollowingUpsert(followerName string, followeeName string, followedOn time.Time) {
	follower := mdb.usersByName[followerName]
	followee := mdb.usersByName[followeeName]
	following := model.Following{
		FollowerId: follower.Id,
		FolloweeId: followee.Id,
		FollowedOn: followedOn,
	}
	mdb.followings[mFollowingT{followerName, followeeName}] = following
}

func (mdb MDb) TagGet(name string) model.Tag {
	return mdb.tags[name]
}

func (mdb MDb) TagExists(name string) bool {
	_, ok := mdb.tags[name]
	return ok
}

func (mdb MDb) TagGetAll() []model.Tag {
	return util.MapToSlice(mdb.tags)
}

func (mdb MDb) TagGetAllNames() []string {
	return util.SliceMap(mdb.TagGetAll(), func(tag model.Tag) string {
		return tag.Name
	})
}

func (mdb *MDb) TagUpsert(name string, tag model.Tag) {
	mdb.tags[name] = tag
}

func (mdb *MDb) TagAssignToSlug(name string, slug string) {
	if !mdb.TagExists(name) {
		panic("attempt to assign an inexistent tag")
	}
	article := mdb.ArticleGetBySlug(slug)
	if article.TagList == nil {
		article.TagList = []string{}
	}
	article.TagList = append(article.TagList, name)
	mdb.ArticleUpsert(article)
}

///////////////////
// Supporting types

// key is Username
type mUsersByNameT map[string]*model.User

// key is user ID
type mUsersByIdT map[uint]*model.User

// key is Slug
type mArticlesT map[string]*model.Article

func (m mArticlesT) upsert(
	article model.Article,
) {
	slug := article.Slug
	m[slug] = &article
}

type mCommentKeyT struct {
	username string
	slug     string
}

// key is an mCommentKeyT, value is a map from comment.Id to model.Comment
type mCommentsT map[mCommentKeyT]map[uint]model.Comment

type mFavoriteT struct {
	username string
	slug     string
}

type mFavoritesT map[mFavoriteT]bool

func (m *mFavoritesT) put(username string, slug string) {
	(*m)[mFavoriteT{username, slug}] = true
}

type mFollowingT struct {
	followerName string
	followeeName string
}

type mFollowingsT map[mFollowingT]model.Following

// key is tag name.
type mTagsT map[string]model.Tag
