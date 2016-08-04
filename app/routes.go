package app

import (
	"git.wid.la/co-net/auth-server/controllers"
	"github.com/go-zoo/bone"
	"github.com/solher/zest"
)

func SetRoutes(z *zest.Zest) error {
	d := &struct {
		Router        *bone.Mux
		AuthCtrl      *controllers.AuthCtrl
		SessionsCtrl  *controllers.SessionsCtrl
		ResourcesCtrl *controllers.ResourcesCtrl
		PoliciesCtrl  *controllers.PoliciesCtrl
	}{}

	if err := z.Injector.Get(d); err != nil {
		return err
	}

	d.Router.GetFunc("/auth", d.AuthCtrl.AuthorizeToken)
	d.Router.GetFunc("/redirect", d.AuthCtrl.Redirect)

	d.Router.GetFunc("/sessions", d.SessionsCtrl.Find)
	d.Router.GetFunc("/sessions/:token", d.SessionsCtrl.FindByToken)
	d.Router.PostFunc("/sessions", d.SessionsCtrl.Create)
	d.Router.DeleteFunc("/sessions", d.SessionsCtrl.DeleteByOwnerToken)
	d.Router.DeleteFunc("/sessions/:token", d.SessionsCtrl.DeleteByToken)

	d.Router.GetFunc("/resources", d.ResourcesCtrl.Find)
	d.Router.GetFunc("/resources/:hostname", d.ResourcesCtrl.FindByHostname)
	d.Router.PostFunc("/resources", d.ResourcesCtrl.Create)
	d.Router.DeleteFunc("/resources/:hostname", d.ResourcesCtrl.DeleteByHostname)
	d.Router.PutFunc("/resources/:hostname", d.ResourcesCtrl.UpdateByHostname)

	d.Router.GetFunc("/policies", d.PoliciesCtrl.Find)
	d.Router.GetFunc("/policies/:name", d.PoliciesCtrl.FindByName)
	d.Router.PostFunc("/policies", d.PoliciesCtrl.Create)
	d.Router.DeleteFunc("/policies/:name", d.PoliciesCtrl.DeleteByName)
	d.Router.PutFunc("/policies/:name", d.PoliciesCtrl.UpdateByName)

	return nil
}
