package database

/**
Usage:

result := struct {
		database.CommonResult
		Found bool  `sql:"COUNT(1)"`
	}{}

	result.Name = ArticleAttachment{}.TableName()

	if database.SelectRowWhere(store.db, &result, nil, "type=? AND article_revision_id=?", AttachmentTypeDrawing, revision.Id) {
		return result.Found
	}
 */

type CommonResult struct {
	Name string
}

func (r CommonResult) TableName() string {
	return r.Name
}
