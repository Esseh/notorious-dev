package FORUM

import(
	"github.com/Esseh/retrievable"
	"github.com/Esseh/notorious-dev/CONTEXT"
)

const password = "password"

func RegisterAdminAPI(ctx CONTEXT.Context) string {
	if ctx.Req.FormValue("Password") != password {
		return `{"success":false}`
	}
	ah := AdminHeader{}
	retrievable.GetEntity(ctx,nil,&ah)
	ah.UserIDS = append(ah.UserIDS,int64(ctx.User.IntID))
	retrievable.PlaceEntity(ctx,nil,&ah)
	return `{"success":true}`
}