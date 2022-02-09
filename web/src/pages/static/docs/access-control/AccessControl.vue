<template>
  <StaticPage
    title="Access Control"
    metadescription="Learn how to setup access control for your codebases"
    image="https://images.unsplash.com/photo-1582139329536-e7284fece509?ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&ixlib=rb-1.2.1&auto=format&fit=crop&w=1480&q=80"
  >
    <div class="prose prose-yellow">
      <p>
        Sturdy supports access control lists. You can use this feature to define granular access for
        each of the collaborators in your codebase.
      </p>

      <ul>
        <li>
          <a href="#introduction">Introduction</a>
          <ul>
            <li><a href="#introduction-acl-policies">ACL policies</a></li>
            <li><a href="#introduction-files">Files</a></li>
          </ul>
        </li>
        <li>
          <a href="#syntax">Syntax</a>
        </li>
        <li>
          <a href="#rules">Rules</a>
          <ul>
            <li><a href="#rules-principals">Identifiers</a></li>
            <li><a href="#rules-files">Files</a></li>
            <li><a href="#rules-action">Action</a></li>
          </ul>
        </li>
        <li>
          <a href="#groups">Groups</a>
        </li>
        <li>
          <a href="#tests">Tests</a>
        </li>
      </ul>

      <h2 id="introduction">Introduction</h2>
      <p>
        Every codebase has an Access Control Policy associated with it - a plain
        <a href="https://github.com/tailscale/hujson">HuJSON</a> file. That file declares what
        members of the codebase have access to. You can find and update the ACL policy for your
        codebase under the settings.
      </p>
      <div class="relative text-base mx-auto max-w-prose lg:max-w-none">
        <figure>
          <div>
            <img
              class="rounded-lg object-cover object-center"
              alt="ACL policy editor"
              src="./access-control.png"
            />
          </div>
          <figcaption class="mt-3 flex text-sm text-gray-500">
            <span class="ml-2">ACL policy editor</span>
          </figcaption>
        </figure>
      </div>
      <p>
        Currently there are two things you can control access to with ACLs:
        <a href="#introduction-acl-policies">ACL policies</a> itself and
        <a href="#introduction-files">Files</a>.
      </p>
      <p>
        Sturdy policies are "default deny". Meaning that a collaborator has access to a file or acl
        policy only if there is a rule that explicitly allows it.
      </p>

      <h3 id="introduction-acl-policies">ACL policies</h3>
      <p>
        Policies for ACL define who can see and update ACL policies. When you create a codebase, it
        is initialized with a default rule that allows anyone to update ACLs.
      </p>

      <h3 id="introduction-files">Files</h3>
      <p>
        Policies for files make it possible to control which files are available for collaborators.
        If a collaborator doesn't have access to a file, they will not be able to see it nor on
        their workstation, nor on the Sturdy web application. They also won't be able to share
        changes for the file.
      </p>
      <p>
        New codebases are initialized with a default rule allowing any collaborator to have a
        <code>write</code> access to any file in the codebase.
      </p>

      <h2 id="syntax">Syntax</h2>
      <p>
        A policy for each or the codebases is expressed as a
        <a href="https://github.com/tailscale/hujson">HuJSON</a> file. HuJSON is a superset of JSON
        that allows comments and commas.
      </p>
      <p>The default policy for a codebase looks like this:</p>
      <ClientOnly>
        <prism-editor
          v-model="defaultPolicy"
          class="max-h-96 p-5 leading-normal text-base font-mono shadow-sm sm:text-sm border-gray-300 rounded-md bg-white"
          :highlight="highlighter"
          readonly
          line-numbers
        ></prism-editor>
      </ClientOnly>
      <p>Every policy has several top-level directives:</p>
      <ul>
        <li><a href="#rules">rules</a>, access policies themselves</li>
        <li>
          <a href="#groups">groups</a>, collections of users or other resources for a simpler rules
          definition
        </li>
        <li>
          <a href="#tests">tests</a>, a list of assertions about the policies that are verified
          every time the policy is updated
        </li>
      </ul>

      <h2 id="rules">Rules</h2>
      <p>
        The <code>rules</code> section of the policy contains a list of access rules for your
        codebase. An access rule looks like this:
      </p>
      <ClientOnly>
        <prism-editor
          v-model="aRule"
          class="max-h-96 p-5 leading-normal text-base font-mono shadow-sm sm:text-sm border-gray-300 rounded-md bg-white"
          :highlight="highlighter"
          readonly
          line-numbers
        ></prism-editor>
      </ClientOnly>
      <p>
        You should read it like that: any principal from the <code>list-of-principals</code> can
        perform <code>action</code> on any resource from the <code>list-of-resources</code>.
      </p>

      <h3 id="rules-principals">Identifiers</h3>
      <p>
        Every principal, resource or a group are identified by it's type and id divided with "::".
        For users, type can be omitted. For example:
      </p>
      <table class="table-auto">
        <thead>
          <tr>
            <th>Identifier</th>
            <th>Type</th>
            <th>Id</th>
            <th>Description</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>
              <code>users::*</code>
            </td>
            <td><code>users</code></td>
            <td>any</td>
            <td>Any user inside a codebase, full version</td>
          </tr>
          <tr>
            <td>
              <code>*</code>
            </td>
            <td><code>users</code></td>
            <td>any</td>
            <td>Any user inside a codebase, short version</td>
          </tr>
          <tr>
            <td>
              <code>users::e849ab49-06ba-4d5a-835f-c1eac578d5a8</code>
            </td>
            <td><code>users</code></td>
            <td><code>e849ab49-06ba-4d5a-835f-c1eac578d5a8</code></td>
            <td>A user identified by id, full version</td>
          </tr>
          <tr>
            <td>
              <code>e849ab49-06ba-4d5a-835f-c1eac578d5a8</code>
            </td>
            <td><code>users</code></td>
            <td><code>e849ab49-06ba-4d5a-835f-c1eac578d5a8</code></td>
            <td>A user identified by id, short version</td>
          </tr>
          <tr>
            <td>
              <code>users::user@exmple.com</code>
            </td>
            <td><code>users</code></td>
            <td><code>user@example.com</code></td>
            <td>A user identified by email, full version</td>
          </tr>
          <tr>
            <td>
              <code>user@example.com</code>
            </td>
            <td><code>users</code></td>
            <td><code>user@example.com</code></td>
            <td>A user identified by email, short version</td>
          </tr>
          <tr>
            <td>
              <code>groups::admins</code>
            </td>
            <td><code>groups</code></td>
            <td><code>admins</code></td>
            <td>
              A group that is defined under <code>groups</code> section with <code>admins</code> id
            </td>
          </tr>
          <tr>
            <td>
              <code>groups::*</code>
            </td>
            <td><code>groups</code></td>
            <td><code>any</code></td>
            <td>Any of the defined groups</td>
          </tr>
          <tr>
            <td>
              <code>files::*</code>
            </td>
            <td><code>files</code></td>
            <td><code>any</code></td>
            <td>Any file in the codebase</td>
          </tr>
          <tr>
            <td>
              <code>files::path/to/a/file</code>
            </td>
            <td><code>files</code></td>
            <td><code>path/to/a/file</code></td>
            <td><code>path/to/a/file</code> exactly</td>
          </tr>
          <tr>
            <td>
              <code>files::path/to/a/file</code>
            </td>
            <td><code>files</code></td>
            <td><code>path/to/dir/**/*</code></td>
            <td>any file or directory under<code>path/to/dir/</code></td>
          </tr>
        </tbody>
      </table>

      <h3 id="rules-files">Files</h3>
      <p>
        File identifiers are designed to match those of .gitignore. At a base level, file
        identifiers are build using
        <a href="https://github.com/bmatcuk/doublestar">doublestar</a> package, so all of the syntax
        carries over:
      </p>
      <table class="table-auto">
        <thead>
          <tr>
            <th>Pattern</th>
            <th>Meaning</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>*</td>
            <td>matches any sequence of non-path-separators</td>
          </tr>
          <tr>
            <td>/**/</td>
            <td>matches zero or more directories</td>
          </tr>
          <tr>
            <td>[class]</td>
            <td>
              matches any single non-path-separator character against a class of characters (see
              <a href="#files-character-classes">"character classes"</a>)
            </td>
          </tr>
          <tr>
            <td>{alt1,...}</td>
            <td>
              matches a sequence of characters if one of the comma-separated alternatives matches
            </td>
          </tr>
        </tbody>
      </table>
      <h4>Character Classes</h4>
      <table id="files-character-classes" class="table-auto">
        <thead>
          <tr>
            <th>Class</th>
            <th>Meaning</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>[abc]</td>
            <td>matches any single character within the set</td>
          </tr>
          <tr>
            <td>[a-z]</td>
            <td>matches any single character in the range</td>
          </tr>
          <tr>
            <td>[^class]</td>
            <td>matches any single character which does not match the class</td>
          </tr>
          <tr>
            <td>[!class]</td>
            <td>same as ^: negates the class</td>
          </tr>
        </tbody>
      </table>

      <p>
        If a file pattern doesn't start with a <code>/</code>, then it is treated as relative to any
        subdirectory in the codebase. For example, <code>file</code> will match any of:
        <code>/file</code>, <code>/dir/file</code>, <code>/dir/subdir/file</code> and so on. While
        <code>/file</code> only matches <code>/file</code>.
      </p>
      <p>
        Prefixing a file pattern with a <code>!</code> will cause it to negate. This way you can,
        for example, give access to every file in a directory, except one:
        <code>["files::dir/", "files::dir/*", "files::!dir/file"]</code>
      </p>

      <h3 id="rules-action">Action</h3>
      <p>Currently only <code>write</code>action is allowed</p>

      <h2 id="groups">Groups</h2>
      <p>Groups is a handy way to create unions of resources to use in rules.</p>
      <p>A group is defined like so:</p>
      <ClientOnly>
        <prism-editor
          v-model="aGroup"
          class="max-h-96 p-5 leading-normal text-base font-mono shadow-sm sm:text-sm border-gray-300 rounded-md bg-white"
          :highlight="highlighter"
          readonly
          line-numbers
        ></prism-editor>
      </ClientOnly>

      <h2 id="tests">Tests</h2>
      <p>
        Last but not the least, tests let you define assertions that Sturdy validates every time
        anyone updates ACL policy. It is a great way to make sure you don't accidentally break an
        important rule. For example,
      </p>
      <ClientOnly>
        <prism-editor
          v-model="aTest"
          class="max-h-96 p-5 leading-normal text-base font-mono shadow-sm sm:text-sm border-gray-300 rounded-md bg-white"
          :highlight="highlighter"
          readonly
          line-numbers
        ></prism-editor>
      </ClientOnly>
      <p>
        This test makes sure that <code>principal</code> has <code>action</code> access to
        <code>resource</code>.
      </p>
    </div>
  </StaticPage>
</template>

<script>
import { PrismEditor } from 'vue-prism-editor'
import 'vue-prism-editor/dist/prismeditor.min.css'

import { highlight, languages } from 'prismjs/components/prism-core'
import 'prismjs/components/prism-json'
import 'prismjs/themes/prism-tomorrow.css'

import StaticPage from '../../../../layouts/StaticPage.vue'
import { ClientOnly } from 'vite-ssr/vue'

export default {
  name: 'AccessControl',
  components: { StaticPage, PrismEditor, ClientOnly },
  setup() {
    return {
      defaultPolicy: `{
  "rules": [
    {
      "id": "everyone can manage access control",
      "principals": ["groups::everyone"],
      "action": "write",
      "resources": ["acls::*"],
    },
    {
      "id": "everyone can access all files",
      "principals": ["groups::everyone"],
      "action": "write",
      "resources": ["files::*"],
    },
  ],
  "groups": [
    {
      "id": "everyone",
      "members": ["*"],
    },
  ],
  "tests": [
    {
      /*
        make sure that at least someone can manage access control lists
      */
      "id": "user@example.com can manage access control",
      "principal": "user@example.com",
      "allow": "write",
      "resource": "acls::bc3b62f8-8cdf-481f-8c04-815912fe668a",
    },
  ],
}`,
      aRule: `{
  "id": "<identifier>",
  "principals": [ "<list-of-principals>" ],
  "action": "<action>",
  "resources": [ "<list-of-resources>" ],
}`,
      aGroup: `{
  "id": "<identifier>",
  "members": [ "<list-of-principals-or-resources"> ]
}`,
      aTest: `{
  "id": "<identifier>",
  "principal": "<principal>",
  "action": "<action>",
  "resource": "<resource>",
}`,
    }
  },
  methods: {
    highlighter(code) {
      return highlight(code, languages.json)
    },
  },
}
</script>
