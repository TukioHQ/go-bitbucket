package bitbucket

// User is the sub struct of Client
type User struct {
	c *Client
}

// Profile is getting the user data
func (u *User) Profile() (interface{}, error) {
	urlStr := GetApiBaseURL() + "/user/"
	return u.c.execute("GET", urlStr, "")
}

// Emails is getting user's emails
func (u *User) Emails() (interface{}, error) {
	urlStr := GetApiBaseURL() + "/user/emails"
	return u.c.execute("GET", urlStr, "")
}

// RepoPermission returns user permissions on a specific bitbucket repository
func (u *User) RepoPermission(repoFullName string) (interface{}, error) {
	urlStr := u.c.requestUrl("/user/permissions/repositories?q=repository.full_name=%q", repoFullName)
	return u.c.execute("GET", urlStr, "")
}
