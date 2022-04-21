/*
 * Copyright © 2022 Paulo Villela. All rights reserved.
 * Use of this source code is governed by the Apache 2.0 license
 * that can be found in the LICENSE file.
 */

package xxx

// See https://github.com/drashland/deno-drash-realworld-example-app/tree/main/src/models

func users() string {
	insert := "INSERT INTO users (username, email, password, bio, image) VALUES ($1, $2, $3, $4, $5);"
	update := "UPDATE users SET " +
		"username = $1, password = $2, email = $3, bio = $4, image = $5 " +
		`WHERE id = $6;`
	del := "DELETE FROM users WHERE id = $1"
	return insert + update + del
}

func articles() string {
	read := "SELECT * FROM articles "
	filtersAuthor := " WHERE author_id = '${filters.author.id}' "
	filtersOffset := " OFFSET ${filters.offset} "
	insert := "INSERT INTO articles " +
		" (author_id, title, description, body, slug, created_at, updated_at, tags)" +
		" VALUES ($1, $2, $3, $4, $5, to_timestamp($6), to_timestamp($7), $8);"
	update := "UPDATE articles SET " +
		"title = $1, description = $2, body = $3, updated_at = to_timestamp($4), tags = $5 " +
		"WHERE id = '${this.id}';"
	del := "DELETE FROM articles WHERE id = $1"
	return read + filtersAuthor + filtersOffset + insert + update + del
}

func articleComments() string {
	read := "SELECT * FROM article_comments"
	filtersArticle := " WHERE article_id = '${filters.article.id}' "
	filtersOffset := " OFFSET ${filters.offset} "
	insert := "INSERT INTO article_comments (article_id, author_image, author_id, author_username, body, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, to_timestamp($6), to_timestamp($7));"
	// may not need update
	update := "UPDATE articles SET " +
		"title = ?, description = ?, body = ?, updatedAt = to_timestamp(?) " +
		"WHERE id = '${this.id}';"
	del := "DELETE FROM article_comments WHERE id = $1"
	return read + filtersArticle + filtersOffset + insert + update + del
}

func articlesFavorites() string {
	insert := "INSERT INTO articles_favorites " +
		" (article_id, user_id, value)" +
		" VALUES ($1, $2, $3);"
	update := "UPDATE articles_favorites SET value = $1 WHERE id = $2;"
	del := "DELETE FROM articles WHERE id = $1"
	return insert + update + del
}

func sessions() string {
	read := "SELECT * FROM sessions " +
		`WHERE session_one = $1 AND session_two = $2 ` +
		"LIMIT 1;"
	insert := "INSERT INTO sessions" +
		" (user_id, session_one, session_two)" +
		" VALUES ($1, $2, $3);"
	return read + insert
}
