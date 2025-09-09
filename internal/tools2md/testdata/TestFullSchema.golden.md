## Tools

### Users

- **get_me** - Get details of the authenticated GitHub user. Use this when a request include "me", "my"...
 - `reason`: Optional: reason the session was created (string, optional)

### Issues

- **get_issue** - Get details of a specific issue in a GitHub repository
 - `issue_number`: The number of the issue (number, required)
 - `owner`: The owner of the repository (string, required)
 - `repo`: The name of the repository (string, required)

- **search_issues** - Search for issues and pull requests across GitHub repositories
 - `order`: Sort order ('asc' or 'desc') (string, optional)
 - `page`: Page number for pagination (min 1) (number, optional)
 - `perPage`: Results per page for pagination (min 1, max 100) (number, optional)
 - `q`: Search query using GitHub issues search syntax (string, required)
 - `sort`: Sort field (comments, reactions, created, etc.) (string, optional)

- **list_issues** - List issues in a GitHub repository with filtering options
 - `direction`: Sort direction ('asc', 'desc') (string, optional)
 - `labels`: Filter by labels (array, optional)
 - `owner`: Repository owner (string, required)
 - `page`: Page number for pagination (min 1) (number, optional)
 - `perPage`: Results per page for pagination (min 1, max 100) (number, optional)
 - `repo`: Repository name (string, required)
 - `since`: Filter by date (ISO 8601 timestamp) (string, optional)
 - `sort`: Sort by ('created', 'updated', 'comments') (string, optional)
 - `state`: Filter by state ('open', 'closed', 'all') (string, optional)

- **get_issue_comments** - Get comments for a GitHub issue
 - `issue_number`: Issue number (number, required)
 - `owner`: Repository owner (string, required)
 - `page`: Page number (number, optional)
 - `per_page`: Number of records per page (number, optional)
 - `repo`: Repository name (string, required)

- **create_issue** - Create a new issue in a GitHub repository
 - `assignees`: Usernames to assign to this issue (array, optional)
 - `body`: Issue body content (string, optional)
 - `labels`: Labels to apply to this issue (array, optional)
 - `milestone`: Milestone number (number, optional)
 - `owner`: Repository owner (string, required)
 - `repo`: Repository name (string, required)
 - `title`: Issue title (string, required)

- **add_issue_comment** - Add a comment to an existing issue
 - `body`: Comment text (string, required)
 - `issue_number`: Issue number to comment on (number, required)
 - `owner`: Repository owner (string, required)
 - `repo`: Repository name (string, required)

- **update_issue** - Update an existing issue in a GitHub repository
 - `assignees`: New assignees (array, optional)
 - `body`: New description (string, optional)
 - `issue_number`: Issue number to update (number, required)
 - `labels`: New labels (array, optional)
 - `milestone`: New milestone number (number, optional)
 - `owner`: Repository owner (string, required)
 - `repo`: Repository name (string, required)
 - `state`: New state ('open' or 'closed') (string, optional)
 - `title`: New title (string, optional)

### Pull Requests

- **get_pull_request** - Get details of a specific pull request
 - `owner`: Repository owner (string, required)
 - `pullNumber`: Pull request number (number, required)
 - `repo`: Repository name (string, required)

- **list_pull_requests** - List and filter repository pull requests
 - `base`: Filter by base branch (string, optional)
 - `direction`: Sort direction ('asc', 'desc') (string, optional)
 - `head`: Filter by head user/org and branch (string, optional)
 - `owner`: Repository owner (string, required)
 - `page`: Page number for pagination (min 1) (number, optional)
 - `perPage`: Results per page for pagination (min 1, max 100) (number, optional)
 - `repo`: Repository name (string, required)
 - `sort`: Sort by ('created', 'updated', 'popularity', 'long-running') (string, optional)
 - `state`: Filter by state ('open', 'closed', 'all') (string, optional)

- **get_pull_request_files** - Get the list of files changed in a pull request
 - `owner`: Repository owner (string, required)
 - `pullNumber`: Pull request number (number, required)
 - `repo`: Repository name (string, required)

- **get_pull_request_status** - Get the combined status of all status checks for a pull request
 - `owner`: Repository owner (string, required)
 - `pullNumber`: Pull request number (number, required)
 - `repo`: Repository name (string, required)

- **get_pull_request_comments** - Get the review comments on a pull request
 - `owner`: Repository owner (string, required)
 - `pullNumber`: Pull request number (number, required)
 - `repo`: Repository name (string, required)

- **get_pull_request_reviews** - Get the reviews on a pull request
 - `owner`: Repository owner (string, required)
 - `pullNumber`: Pull request number (number, required)
 - `repo`: Repository name (string, required)

- **merge_pull_request** - Merge a pull request
 - `commit_message`: Extra detail for merge commit (string, optional)
 - `commit_title`: Title for merge commit (string, optional)
 - `merge_method`: Merge method ('merge', 'squash', 'rebase') (string, optional)
 - `owner`: Repository owner (string, required)
 - `pullNumber`: Pull request number (number, required)
 - `repo`: Repository name (string, required)

- **update_pull_request_branch** - Update a pull request branch with the latest changes from the base branch
 - `expectedHeadSha`: The expected SHA of the pull request's HEAD ref (string, optional)
 - `owner`: Repository owner (string, required)
 - `pullNumber`: Pull request number (number, required)
 - `repo`: Repository name (string, required)

- **create_pull_request_review** - Create a review on a pull request
 - `body`: Review comment text (string, optional)
 - `comments`: Line-specific comments array of objects to place comments on pull request changes. Requires path and body. For line comments use line or position. For multi-line comments use start_line and line with optional side parameters. (array, optional)
 - `commitId`: SHA of commit to review (string, optional)
 - `event`: Review action ('APPROVE', 'REQUEST_CHANGES', 'COMMENT') (string, required)
 - `owner`: Repository owner (string, required)
 - `pullNumber`: Pull request number (number, required)
 - `repo`: Repository name (string, required)

- **create_pull_request** - Create a new pull request in a GitHub repository
 - `base`: Branch to merge into (string, required)
 - `body`: PR description (string, optional)
 - `draft`: Create as draft PR (boolean, optional)
 - `head`: Branch containing changes (string, required)
 - `maintainer_can_modify`: Allow maintainer edits (boolean, optional)
 - `owner`: Repository owner (string, required)
 - `repo`: Repository name (string, required)
 - `title`: PR title (string, required)

- **update_pull_request** - Update an existing pull request in a GitHub repository
 - `base`: New base branch name (string, optional)
 - `body`: New description (string, optional)
 - `maintainer_can_modify`: Allow maintainer edits (boolean, optional)
 - `owner`: Repository owner (string, required)
 - `pullNumber`: Pull request number to update (number, required)
 - `repo`: Repository name (string, required)
 - `state`: New state ('open' or 'closed') (string, optional)
 - `title`: New title (string, optional)

### Repositories

- **get_file_contents** - Get the contents of a file or directory from a GitHub repository
 - `branch`: Branch to get contents from (string, optional)
 - `owner`: Repository owner (username or organization) (string, required)
 - `path`: Path to file/directory (string, required)
 - `repo`: Repository name (string, required)

- **get_commit** - Get details for a commit from a GitHub repository
 - `owner`: Repository owner (string, required)
 - `page`: Page number for pagination (min 1) (number, optional)
 - `perPage`: Results per page for pagination (min 1, max 100) (number, optional)
 - `repo`: Repository name (string, required)
 - `sha`: Commit SHA, branch name, or tag name (string, required)

- **list_commits** - Get list of commits of a branch in a GitHub repository
 - `owner`: Repository owner (string, required)
 - `page`: Page number for pagination (min 1) (number, optional)
 - `perPage`: Results per page for pagination (min 1, max 100) (number, optional)
 - `repo`: Repository name (string, required)
 - `sha`: Branch name (string, optional)

- **create_or_update_file** - Create or update a single file in a GitHub repository
 - `branch`: Branch to create/update the file in (string, required)
 - `content`: Content of the file (string, required)
 - `message`: Commit message (string, required)
 - `owner`: Repository owner (username or organization) (string, required)
 - `path`: Path where to create/update the file (string, required)
 - `repo`: Repository name (string, required)
 - `sha`: SHA of file being replaced (for updates) (string, optional)

- **create_repository** - Create a new GitHub repository in your account
 - `autoInit`: Initialize with README (boolean, optional)
 - `description`: Repository description (string, optional)
 - `name`: Repository name (string, required)
 - `private`: Whether repo should be private (boolean, optional)

- **fork_repository** - Fork a GitHub repository to your account or specified organization
 - `organization`: Organization to fork to (string, optional)
 - `owner`: Repository owner (string, required)
 - `repo`: Repository name (string, required)

- **create_branch** - Create a new branch in a GitHub repository
 - `branch`: Name for new branch (string, required)
 - `from_branch`: Source branch (defaults to repo default) (string, optional)
 - `owner`: Repository owner (string, required)
 - `repo`: Repository name (string, required)

- **push_files** - Push multiple files to a GitHub repository in a single commit
 - `branch`: Branch to push to (string, required)
 - `files`: Array of file objects to push, each object with path (string) and content (string) (array, required)
 - `message`: Commit message (string, required)
 - `owner`: Repository owner (string, required)
 - `repo`: Repository name (string, required)

### Search

- **search_repositories** - Search for GitHub repositories
 - `page`: Page number for pagination (min 1) (number, optional)
 - `perPage`: Results per page for pagination (min 1, max 100) (number, optional)
 - `query`: Search query (string, required)

- **search_code** - Search for code across GitHub repositories
 - `order`: Sort order ('asc' or 'desc') (string, optional)
 - `page`: Page number for pagination (min 1) (number, optional)
 - `perPage`: Results per page for pagination (min 1, max 100) (number, optional)
 - `q`: Search query using GitHub code search syntax (string, required)
 - `sort`: Sort field ('indexed' only) (string, optional)

- **search_users** - Search for GitHub users
 - `order`: Sort order ('asc' or 'desc') (string, optional)
 - `page`: Page number for pagination (min 1) (number, optional)
 - `perPage`: Results per page for pagination (min 1, max 100) (number, optional)
 - `q`: Search query using GitHub users search syntax (string, required)
 - `sort`: Sort field (followers, repositories, joined) (string, optional)

### Code Scanning

- **get_code_scanning_alert** - Get details of a specific code scanning alert in a GitHub repository.
 - `alertNumber`: The number of the alert. (number, required)
 - `owner`: The owner of the repository. (string, required)
 - `repo`: The name of the repository. (string, required)

- **list_code_scanning_alerts** - List code scanning alerts in a GitHub repository.
 - `owner`: The owner of the repository. (string, required)
 - `ref`: The Git reference for the results you want to list. (string, optional)
 - `repo`: The name of the repository. (string, required)
 - `severity`: Only code scanning alerts with this severity will be returned. Possible values are: critical, high, medium, low, warning, note, error. (string, optional)
 - `state`: State of the code scanning alerts to list. Set to closed to list only closed code scanning alerts. Default: open (string, optional)
