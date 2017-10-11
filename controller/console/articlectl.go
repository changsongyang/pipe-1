// Solo.go - A small and beautiful blogging platform written in golang.
// Copyright (C) 2017, b3log.org
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Package console defines console controllers.
package console

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/b3log/solo.go/model"
	"github.com/b3log/solo.go/service"
	"github.com/b3log/solo.go/util"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func AddArticleCtl(c *gin.Context) {
	result := util.NewResult()
	defer c.JSON(http.StatusOK, result)

	sessionData := util.GetSession(c)

	article := &model.Article{}
	if err := c.BindJSON(article); nil != err {
		result.Code = -1
		result.Msg = "parses add article request failed"

		return
	}

	article.BlogID = sessionData.BID
	article.AuthorID = sessionData.UID
	if err := service.Article.ConsoleAddArticle(article); nil != err {
		result.Code = -1
		result.Msg = err.Error()
	}
}

func GetArticleCtl(c *gin.Context) {
	result := util.NewResult()
	defer c.JSON(http.StatusOK, result)

	idArg := c.Param("id")
	id, err := strconv.Atoi(idArg)
	if nil != err {
		result.Code = -1

		return
	}

	data := service.Article.ConsoleGetArticle(uint(id))
	if nil == data {
		result.Code = -1

		return
	}

	result.Data = data
}

type ConsoleArticle struct {
	ID           uint          `json:"id"`
	Author       *Author       `json:"author"`
	CreatedAt    string        `json:"createdAt"`
	Title        string        `gorm:"size:128" json:"title"`
	Tags         []*ConsoleTag `json:"tags"`
	URL          string        `json:"url"`
	Topped       bool          `json:"topped"`
	ViewCount    int           `json:"viewCount"`
	CommentCount int           `json:"commentCount"`
}

type ConsoleTag struct {
	Title string `json:"title"`
	URL   string `json:"url,omitempty"`
}

type Author struct {
	Name      string `json:"name"`
	AvatarURL string `json:"avatarURL"`
}

func GetArticlesCtl(c *gin.Context) {
	result := util.NewResult()
	defer c.JSON(http.StatusOK, result)

	sessionData := util.GetSession(c)
	articleModels, pagination := service.Article.ConsoleGetArticles(c.GetInt("p"), sessionData.BID)

	articles := []*ConsoleArticle{}
	for _, articleModel := range articleModels {
		consoleTags := []*ConsoleTag{}
		tagStrs := strings.Split(articleModel.Tags, ",")
		for _, tagStr := range tagStrs {
			consoleTag := &ConsoleTag{
				Title: tagStr,
				URL:   sessionData.BPath + "/" + tagStr,
			}
			consoleTags = append(consoleTags, consoleTag)
		}

		authorModel := service.User.GetUser(articleModel.AuthorID)
		if nil == authorModel {
			log.Errorf("not found author of article [id=%d, authorID=%d]", articleModel.ID, articleModel.AuthorID)

			continue
		}

		author := &Author{
			Name:      authorModel.Name,
			AvatarURL: authorModel.AvatarURL,
		}

		article := &ConsoleArticle{
			ID:           articleModel.ID,
			Author:       author,
			CreatedAt:    articleModel.CreatedAt.Format("2006-01-02"),
			Title:        articleModel.Title,
			Tags:         consoleTags,
			URL:          sessionData.BPath + articleModel.Path,
			Topped:       articleModel.Topped,
			ViewCount:    articleModel.ViewCount,
			CommentCount: articleModel.CommentCount,
		}

		articles = append(articles, article)
	}

	data := map[string]interface{}{}
	data["articles"] = articles
	data["pagination"] = pagination
	result.Data = data
}

func RemoveArticleCtl(c *gin.Context) {
	result := util.NewResult()
	defer c.JSON(http.StatusOK, result)

	idArg := c.Param("id")
	id, err := strconv.Atoi(idArg)
	if nil != err {
		result.Code = -1
		result.Msg = err.Error()

		return
	}

	if err := service.Article.ConsoleRemoveArticle(uint(id)); nil != err {
		result.Code = -1
		result.Msg = err.Error()
	}
}

func UpdateArticleCtl(c *gin.Context) {
	result := util.NewResult()
	defer c.JSON(http.StatusOK, result)

	idArg := c.Param("id")
	id, err := strconv.Atoi(idArg)
	if nil != err {
		result.Code = -1
		result.Msg = err.Error()

		return
	}

	article := &model.Article{Model: model.Model{ID: uint(id)}}
	if err := c.BindJSON(article); nil != err {
		result.Code = -1
		result.Msg = "parses update article request failed"

		return
	}

	if err := service.Article.ConsoleUpdateArticle(article); nil != err {
		result.Code = -1
		result.Msg = err.Error()
	}
}