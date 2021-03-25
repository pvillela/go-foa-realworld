package ft

import "github.com/pvillela/go-foa-realworld/internal/model"

// ArticleGetAndCheckOwnerFlT is the function type instantiated by fs.ArticleGetAndCheckOwnerFl.
type ArticleGetAndCheckOwnerFlT = func(username, slug string) (*model.Article, error)

// ArticleFavoriteFlT is the function type instantiated by fs.ArticleFavoriteFl.
type ArticleFavoriteFlT = func(username, slug string, favorite bool) (*model.User, *model.Article, error)

// ArticleListFlT is the function type instantiated by fs.ArticleListFl.
type ArticleListFlT = func(username string, limit, offset int) (*model.User, []model.Article, error)
