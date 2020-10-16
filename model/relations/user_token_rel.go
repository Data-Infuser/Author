package relations

import (
	"net/http"

	errors "gitlab.com/promptech1/infuser-author/error"
	"gitlab.com/promptech1/infuser-author/model"
	"xorm.io/xorm"
)

// XORM 처리를 위한 별도 struct 구성
// join 쿼리 수행시 struct에 명시된 순서 주의 필요함 (https://gobook.io/read/gitea.com/xorm/manual-en-US/chapter-05/5.join.html)
type UserTokenRel struct {
	User  model.User      `xorm:"extends"`
	Token model.UserToken `xorm:"extends"`
}

func (ut *UserTokenRel) FindByUserLoginId(orm *xorm.Engine) error {
	found, err := orm.Table("user").Join(
		"LEFT OUTER", "user_token",
		"user.id = user_token.user_id",
	).Where("user.login_id = ?", ut.User.LoginId).Get(ut)

	if err != nil {
		return errors.NewWithPrefix(err, "database error")
	}

	if !found {
		return errors.NewWithCode(http.StatusNotFound, "UserToken not found")
	}

	return nil
}

func (ut *UserTokenRel) FindByRefreshToken(orm *xorm.Engine) error {
	var utr UserTokenRel
	found, err := orm.Table("user").Join(
		"INNER", "user_token",
		"user.id = user_token.user_id",
	).Where("user_token.refresh_token = ?", ut.Token.RefreshToken).Get(&utr)
	*ut = utr

	if err != nil {
		return errors.NewWithPrefix(err, "database error")
	}

	if !found {
		return errors.NewWithCode(http.StatusNotFound, "UserToken not found")
	}

	return nil
}
