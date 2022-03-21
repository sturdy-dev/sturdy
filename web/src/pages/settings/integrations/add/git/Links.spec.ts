import {defaultLinkBranch, defaultLinkRepo, removeBasicAuth} from "./Links";

describe('Links', () => {
    it ('Remove Basic Auth', () => {
        expect(removeBasicAuth('https://foobar@gitlab.com/zegl/sturdy-push.git')).toEqual("https://gitlab.com/zegl/sturdy-push.git")
        expect(removeBasicAuth('http://foobar@gitlab.com/zegl/sturdy-push.git')).toEqual("http://gitlab.com/zegl/sturdy-push.git")
        expect(removeBasicAuth('http://example.com')).toEqual("http://example.com")
    })

    it ('GitLab', () => {
        const link = 'https://gitlab.com/zegl/sturdy-push.git'
        expect(defaultLinkRepo(link)).toEqual("https://gitlab.com/zegl/sturdy-push")
        expect(defaultLinkBranch(link)).toEqual("https://gitlab.com/zegl/sturdy-push/-/tree/${BRANCH_NAME}")
    })

    it ('Azure', () => {
        const link = 'https://gustav0446@dev.azure.com/gustav0446/yolo_sturdy_test/_git/yolo_sturdy_test'
        expect(defaultLinkRepo(link)).toEqual("https://dev.azure.com/gustav0446/_git/yolo_sturdy_test")
        expect(defaultLinkBranch(link)).toEqual("https://dev.azure.com/gustav0446/yolo_sturdy_test/_git/yolo_sturdy_test?version=GB${BRANCH_NAME}")
    })

    it ('BitBucket', () => {
        const link = 'https://zegl@bitbucket.org/zegl/taget-go.git'
        expect(defaultLinkRepo(link)).toEqual("https://bitbucket.org/zegl/taget-go")
        expect(defaultLinkBranch(link)).toEqual("https://bitbucket.org/zegl/taget-go/branch/${BRANCH_NAME}")
    })
})