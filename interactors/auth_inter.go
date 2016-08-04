package interactors

import (
	"strings"

	"github.com/solher/auth-nginx-proxy-companion/errs"
	"github.com/solher/auth-nginx-proxy-companion/models"
	"github.com/solher/zest"
)

func init() {
	zest.Injector.Register(NewAuthInter)
}

type (
	AuthInterPoliciesInter interface {
		FindByName(name string) (*models.Policy, error)
	}

	AuthInterResourcesInter interface {
		FindByHostname(hostname string) (*models.Resource, error)
	}

	AuthInterSessionsInter interface {
		FindByToken(id string) (*models.Session, error)
	}

	AuthInter struct {
		policiesInter  AuthInterPoliciesInter
		resourcesInter AuthInterResourcesInter
		sessionsInter  AuthInterSessionsInter
	}
)

func NewAuthInter(
	policiesInter AuthInterPoliciesInter,
	resourcesInter AuthInterResourcesInter,
	sessionsInter AuthInterSessionsInter,
) *AuthInter {
	return &AuthInter{
		policiesInter:  policiesInter,
		resourcesInter: resourcesInter,
		sessionsInter:  sessionsInter,
	}
}

func (i *AuthInter) GetRedirectURL(hostname string) (string, error) {
	resource, err := i.resourcesInter.FindByHostname(hostname)
	if err != nil {
		return "", err
	}

	if resource.RedirectURL == nil {
		return "", nil
	}

	return *resource.RedirectURL, nil
}

func (i *AuthInter) AuthorizeToken(hostname, path, token string) (bool, *models.Session, error) {
	// We try to find concurrently the resource and the session corresponding to the request
	resourceCh, errCh1 := i.findResource(hostname)
	sessionCh, errCh2 := i.findSession(token)

	// If we can't find a resource, we deny the access
	if err := <-errCh1; err != nil {
		return false, nil, err
	}

	resource := <-resourceCh

	// If the found resource is marked as public, we allow the access without restriction
	if resource.Public != nil && *resource.Public {
		return true, nil, nil
	}

	// If no session is found for the token, we initiate a guest session
	// Otherwise, the access is denied with an error
	if err := <-errCh2; err != nil {
		switch err.(type) {
		case errs.ErrNotFound:
			return i.authorizeGuestSession(path, *resource.Name)
		default:
			return false, nil, err
		}
	}

	// If a session is found, we try to authorize it
	return i.authorizeSession(path, *resource.Name, <-sessionCh)
}

func (i *AuthInter) authorizeSession(path, resource string, session *models.Session) (bool, *models.Session, error) {
	ch := make(chan bool, len(session.Policies))
	errCh := make(chan error, len(session.Policies))

	// We check concurrently the associated policies and the permissions associated
	for _, policyID := range session.Policies {
		go i.checkPermissions(path, resource, policyID, ch, errCh)
	}

	// We don't wait for all the policies to be checked
	// We return as soon as we find a positive result
	for range session.Policies {
		select {
		case granted := <-ch:
			if granted {
				return true, session, nil
			}
		case err := <-errCh:
			return false, nil, err
		}
	}

	return false, nil, nil
}

func (i *AuthInter) authorizeGuestSession(path, resource string) (bool, *models.Session, error) {
	ch := make(chan bool, 1)
	errCh := make(chan error, 1)

	// We check the guest permissions
	i.checkPermissions(path, resource, "guest", ch, errCh)

	// We don't wait for all the policies to be checked
	// We return as soon as we find a positive result
	select {
	case granted := <-ch:
		if granted {
			return true, nil, nil
		}
	case err := <-errCh:
		return false, nil, err
	}

	return false, nil, nil
}

func (i *AuthInter) checkPermissions(path, resource, policyName string, ch chan bool, errCh chan error) {
	// First, we find the policy corresponding to the given name in database
	policy, err := i.policiesInter.FindByName(policyName)
	if err != nil {
		errCh <- err
		return
	}

	// If the policy is disabled, we skip it
	if policy.Enabled != nil && *policy.Enabled == false {
		ch <- false
		return
	}

	// "reqPath" is the splited path of the incoming request
	// We will use it to compare it with the permissions
	reqPath := i.splitPath(path)
	reqWeight := len(reqPath)

	// "granted" is obviously the boolean value indicating if the access must be granted
	granted := false

	// "maxWeight" is used to ponderate the permissions
	// A permission with a higher weight will override others with a lower one
	// The weight is here the number of "segments" in a permission path
	//
	// Example:
	//   "/foo" -> weight 1
	//   "/foo/bar" -> weight 2
	maxWeight := 0

	// "wildcard" indicates if the current maxWeight was set by a permission with a wildcard
	// In that case, a regular permission with the same weight would override it
	wildcard := false

	// We now check each permission of the policy
	for _, permission := range policy.Permissions {
		// If the permission does not concern the requested resource, we skip it
		if *permission.Resource != resource && *permission.Resource != "*" {
			continue
		}

		// If the permission is disabled, we skip it
		if permission.Enabled != nil && *permission.Enabled == false {
			continue
		}

		// nil paths is considered as a wildcard
		if permission.Paths == nil {
			permission.Paths = []string{"*"}
		}

		for _, path := range permission.Paths {
			// We get the splited path and the weight of the permission
			permPath := i.splitPath(path)
			permWeight := len(permPath)

			// If the weight of the permission is higher than the weight of the request, we skip it
			//
			// Example:
			//    Req: "/foo"
			//    Perm: "/foo/bar" -> Does not apply here
			if permWeight > reqWeight {
				continue
			}

			// We override the granted/maxWeight/wildcard variables if the paths match and:
			//   - Current permission weight is higher than the current maxWeight
			//   or
			//   - Current permission weight is equal to the current maxWeight but was set by a wildcard
			if ok, wc := i.match(reqPath, permPath); ok && ((permWeight > maxWeight) || (wildcard && permWeight == maxWeight)) {
				if permission.Deny == nil {
					granted = true
				} else {
					granted = !*permission.Deny
				}

				maxWeight = permWeight
				wildcard = wc
			}

		}
	}

	// We return the result
	ch <- granted
}

func (i *AuthInter) findResource(hostname string) (chan *models.Resource, chan error) {
	ch := make(chan *models.Resource, 1)
	errCh := make(chan error, 1)

	go func() {
		resource, err := i.resourcesInter.FindByHostname(hostname)
		if err != nil {
			errCh <- err
			close(ch)
			return
		}

		ch <- resource
		close(errCh)
	}()

	return ch, errCh
}

func (i *AuthInter) findSession(token string) (chan *models.Session, chan error) {
	ch := make(chan *models.Session, 1)
	errCh := make(chan error, 1)

	go func() {
		resource, err := i.sessionsInter.FindByToken(token)
		if err != nil {
			errCh <- err
			close(ch)
			return
		}

		ch <- resource
		close(errCh)
	}()

	return ch, errCh
}

func (i *AuthInter) splitPath(path string) []string {
	return strings.Split(strings.TrimPrefix(strings.TrimSuffix(path, "/"), "/"), "/")
}

func (i *AuthInter) match(reqPath, permPath []string) (bool, bool) {
	for i, p := range permPath {
		switch p {
		case reqPath[i]:
			if i == len(permPath)-1 && len(permPath) < len(reqPath) {
				return false, false
			}
			continue
		case "*":
			return true, true
		default:
			return false, false
		}
	}

	return true, false
}
