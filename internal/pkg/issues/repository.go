package issues

// Repository interface definition
type Repository interface {
	FindUser(ID string) (*User, error)
	FindProjects() (Projects, error)
	FindIssues(keys []string) (Issues, error)
	FindBigSellerIssues(keys []string) (Issues, error)
}
