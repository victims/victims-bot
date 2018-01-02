package github

type Commit struct {
	ID        string   `json:"id"`
	TreeID    string   `json:"tree_id"`
	Distrinct bool     `json:"distinct"`
	Message   string   `json:"message"`
	Timestamp string   `json:"timestamp"`
	URL       string   `json:"url"`
	Author    User     `json:"author"`
	Committer User     `json:"comitter"`
	Added     []string `json:"added"`
	Removed   []string `json:"removed"`
	Modified  []string `json:"modified"`
}

type Owner struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type User struct {
	Owner
	Usersname string `json:"username"`
}

type Sender struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HtmlURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	Following         string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SideAdmin         bool   `json:"side_admin"`
}

type Repository struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	FullName         string `json:"full_name"`
	Owner            Owner  `json:"owner"`
	Private          bool   `json:"private"`
	HtmlURL          string `json:"html_url"`
	Description      string `json:"description"`
	Fork             bool   `json:"fork"`
	URL              string `json:"url"`
	ForksURL         string `json:"forks_url"`
	KeysURL          string `json:"keys_url"`
	CollaboratorsURL string `json:"collaborators_url"`
	TeamsURL         string `json:"teams_url"`
	HooksURL         string `json:"hooks_url"`
	IssueEventsURL   string `json:"issue_events_url"`
	EventsURL        string `json:"events_url"`
	AssigneesURL     string `json:"assignees_url"`
	BranchesURL      string `json:"branches_url"`
	TagsURL          string `json:"tags_url"`
	BlobsURL         string `json:"blobs_url"`
	GitTagsURL       string `json:"git_tags_url"`
	GitRefsURL       string `json:"git_refs_url"`
	TreesURL         string `json:"trees_url"`
	StatusesURL      string `json:"status_url"`
	LangaugesURL     string `json:"languages_url"`
	StargazersURL    string `json:"stargazers_url"`
	ContributorsURL  string `json:"contributors_url"`
	SubscribersURL   string `json:"subscribers_url"`
	SubscriptionURL  string `json:"subscription_url"`
	CommitsURL       string `json:"commits_url"`
	GitCommitsURL    string `json:"git_commits_url"`
	CommentsURL      string `json:"comments_url"`
	IssueCommentsURL string `json:"issue_comments_url"`
	ContentsURL      string `json:"contents_url"`
	CompareURL       string `json:"compare_url"`
	MergesURL        string `json:"merges_url"`
	ArchiveURL       string `json:"archive_url"`
	DownloadsURL     string `json:"downloads_url"`
	IssuesURL        string `json:"issues_url"`
	PullsURL         string `json:"pulls_url"`
	MilestoneURL     string `json:"milestone_url"`
	NotificationURL  string `json:"notification_url"`
	LabelsURL        string `json:"labels_url"`
	ReleasesURL      string `json:"releases_url"`
	CreatedAt        int    `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
	PushedAt         int    `json:"pushed_at"`
	GitURL           string `json:"git_url"`
	SshURL           string `json:"ssh_url"`
	CloneURL         string `json:"clone_url"`
	SvnURL           string `json:"svn_url"`
	HomePage         string `json:"home_page"`
	Size             int    `json:"size"`
	StargazersCount  int    `json:"stargazers_count"`
	WatchersCount    int    `json:"watchers_count"`
	Language         string `json:"language"`
	HasIssues        bool   `json:"has_issues"`
	HasDownloads     bool   `json:"has_downloads"`
	HasWiki          bool   `json:"has_wiki"`
	HasPages         bool   `json:"has_pages"`
	ForkCount        int    `json:"fork_count"`
	MirrorURL        string `json:"mirror_url"`
	OpenIssuesCount  int    `json:"open_issues_count"`
	Forks            int    `json:"forks"`
	OpenIssues       int    `json:"open_issues"`
	Watchers         int    `json:"watchers"`
	DefaultBranch    string `json:"default_branch"`
	Stargazers       int    `json:"stargzers"`
	MasterBranch     string `json:"master_branch"`
}

type PushEvent struct {
	Ref        string     `json:"ref"`
	Before     string     `json:"before"`
	After      string     `json:"after"`
	Created    bool       `json:"created"`
	Deleted    bool       `json:"deleted"`
	Forced     bool       `json:"forced"`
	BaseRef    string     `json:"base_ref"`
	Compare    string     `json:"compare"`
	Commits    []Commit   `json:"commits"`
	HeadCommit Commit     `json:"head_commit"`
	Repository Repository `json:"repository"`
	Pusher     Owner      `json:"pusher"`
	Sender     Sender     `json:"sender"`
}
