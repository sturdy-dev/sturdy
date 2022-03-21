export const removeBasicAuth = function(url: string): string {
    if (!url.includes("@")) {
        return url
    }

    const slashslash = url.indexOf("://")
    const at = url.indexOf("@")

    return url.substring(0, slashslash + 3) + url.substring(at + 1)
}

export const defaultLinkRepo = function (gitURL: string): string | null {
    gitURL = removeBasicAuth(gitURL)

    // GitLab (example: https://gitlab.com/zegl/sturdy-push.git)
    if (gitURL.startsWith('https://gitlab.com/') && gitURL.endsWith('.git')) {
        return gitURL.substring(0, gitURL.length - 4)
    }

    // BitBucket
    if (gitURL.startsWith('https://bitbucket.org/') && gitURL.endsWith('.git')) {
        return gitURL.substring(0, gitURL.length - 4)
    }

    // Azure
    // Push URL: https://gustav0446@dev.azure.com/gustav0446/yolo_sturdy_test/_git/yolo_sturdy_test
    // https://dev.azure.com/gustav0446/yolo_sturdy_test
    // https://dev.azure.com/gustav0446/_git/yolo_sturdy_test
    if (gitURL.includes("dev.azure.com") && gitURL.includes("/_git/")) {
        const s = gitURL.indexOf("dev.azure.com/")
        const git = gitURL.indexOf("/_git/")
        const parts = gitURL.substring(s, git).split("/")
        if (parts.length === 3) {
            return "https://" + parts[0] + "/" + parts[1] + "/_git/" + parts[2]
        }
    }

    return null
}

export const defaultLinkBranch = function (gitURL: string): string | null {
    gitURL = removeBasicAuth(gitURL)

    // GitLab (example: https://gitlab.com/zegl/sturdy-push.git)
    if (gitURL.startsWith('https://gitlab.com/') && gitURL.endsWith('.git')) {
        // Example: https://gitlab.com/foo-org/foo-repo/-/tree/sturdy-726608fc-59c0-4475-81b0-28ae5d12d53d
        return gitURL.substring(0, gitURL.length - 4) + '/-/tree/${BRANCH_NAME}'
    }

    // BitBucket
    if (gitURL.startsWith('https://bitbucket.org/') && gitURL.endsWith('.git')) {
        return gitURL.substring(0, gitURL.length - 4) + '/branch/${BRANCH_NAME}'
    }

    // Azure
    // Push URL: https://gustav0446@dev.azure.com/gustav0446/yolo_sturdy_test/_git/yolo_sturdy_test
    // https://dev.azure.com/gustav0446/yolo_sturdy_test
    if (gitURL.includes("dev.azure.com") && gitURL.includes("/_git/")) {
        const s = gitURL.indexOf("dev.azure.com/")
        const git = gitURL.indexOf("/_git/")
        return "https://" + gitURL.substring(s, git) + "/_git/yolo_sturdy_test?version=GB${BRANCH_NAME}"
    }

    return null
}
