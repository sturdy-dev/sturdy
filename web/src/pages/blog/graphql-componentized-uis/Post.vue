<template>
  <BlogPost
    title="GraphQL & Componentized UIs"
    subtitle="Making full use of the normalized cache"
    date="December 17, 2021"
    reading-time="About 10 minutes"
    description="GraphQL has become somewhat of an industry standard for APIs that are designed to be consumed
      by component-oriented UIs. Meanwhile, GraphQL client libraries have converged on strategies
      for maintaining a normalized cache of all parts of the graph that has been queried. This
      article proposes some patterns on how to maximize the utility and maintainability of frontend
      code using this stack."
    :image="ogImageFull"
  >
    <template #introduction>
      GraphQL has become somewhat of an industry standard for APIs that are designed to be consumed
      by component-oriented UIs. Meanwhile, GraphQL client libraries have converged on strategies
      for maintaining a normalized cache of all parts of the graph that has been queried. This
      article proposes some patterns on how to maximize the utility and maintainability of frontend
      code using this stack. It assumes some knowledge of how GraphQL works. The examples use React,
      but the same strategies can be used with other component-oriented libraries. Let’s jump in!
    </template>

    <div ref="container">
      <h2>Part 1. The Graph as the Source of Truth</h2>

      <p>
        Popular GraphQL client libraries like
        <a href="https://www.apollographql.com/">Apollo</a> and
        <a href="https://formidable.com/open-source/urql/">Urql</a> (the latter of which we use
        internally at Sturdy) feature
        <a href="https://formidable.com/open-source/urql/docs/graphcache/normalized-caching/"
          >normalized caching</a
        >. For instance, let’s say you render a list of search results, where each result represents
        some entity, for which there is also a dedicated page. If the user clicks a search result,
        the app navigates to that entity’s page. Since we’ve queried the API for the search results,
        there is already some information loaded on each entity in the result. Most likely, the
        entity has some display name, which was used in the search result. We also want to show that
        name on the dedicated page for the entity. The normalized cache figures out that the entity
        from the search results list shares its type and ID with the entity shown on the page. And
        if it turns out we already have all the data we need to render the page, we can skip taking
        another round trip to the server to figure out what the name of that entity is – we already
        know from the search result.
      </p>

      <p>
        This is neat, but not necessarily revolutionary. The really cool thing is that both of these
        places – the search result and the page – are now
        <em>connected to the same cache entry</em>. This means that if the cache were to update,
        these components would rerender with the new data, without us having to manually handle that
        update.
      </p>

      <p>
        Here’s where mutations come in. When the user performs some write operation, the mutation
        declaration also specifies what fields should be returned from the API once the operation is
        completed. That response will be written to the cache, which will automatically update the
        components which currently render any data that was updated by the mutation.
      </p>

      <p>
        The best part of all this is that the code that triggers the mutation
        <em>doesn’t have to take responsibility for triggering updates in the UI</em>. Instead, the
        graph representation in the cache makes for the “source of truth” for the entire UI state.
      </p>

      <h2>Part 2. Queries Follow UI, Mutations Don't</h2>

      <p>
        Code related to dispatching mutations can be organized completely separated from UI
        components, e.g. in custom hooks. As long as the cache is automatically updated correctly
        after the mutation has been completed, the code that dispatches the mutation doesn’t have to
        do <em>anything at all</em> with the result. By contrast, queries are completely tied to the
        UI in the sense that it’s the UI components' usage of the data that dictates what fields
        need to be selected in the query.
      </p>

      <p>
        Because of this, we find that queries and mutations have diametrically opposed requirements.
        To reiterate: queries <em>select the fields that are needed to render UI</em>, while
        mutations
        <em
          >select only the fields that change in response to the called mutation being applied,
          regardless of UI</em
        >.
      </p>

      <p>
        Why do I try to make this point? I’m going to try to give a concrete example of how this
        distinction pulls the usages of mutations and queries in two different directions.
      </p>

      <p>
        When you start working with GraphQL, it’s really tempting to extract the query logic to
        something reusable. E.g. “Let's make a <code>useCurrentUser</code> hook that we can reuse
        anytime we need data from the current user.” The trouble with this is this: what fields
        should be queried inside the hook? How do you prevent it from querying too much or too
        little data for any given use of the hook?
      </p>

      <pre><code data-highlight class="language-typescript">import { useQuery, gql } from "urql";

const CURRENT_USER = gql`
  query CurrentUser {
    user {
      id
      email  # where is this used?
      name   # can we delete this?
    }
  }
`;

export function useCurrentUser() {
  return useQuery({ query: CURRENT_USER });
}</code></pre>

      <p>
        Instead, the selection of fields on a given type for use in a component
        <em>should be decided by the component itself</em>. GraphQL provides us with the fragment
        feature for exactly this kind of reason. The component that uses some field can declare that
        in an exported fragment.
      </p>

      <pre><code data-highlight class="language-typescript">import { gql } from "urql";

export const AVATAR_USER = gql`
  fragment AvatarUser on User {
    id
    avatarUrl
  }
`;

export function Avatar({ user }) {
  return &lt;img src={user?.avatarUrl ?? "/default_user.svg"} />;
}</code></pre>

      <p>
        The parent component is now given the responsibility of including the fragment in its query.
      </p>

      <pre><code data-highlight class="language-typescript">import { useQuery, gql } from "urql";
import { Avatar, AVATAR_USER } from "./Avatar";

const CURRENT_USER = gql`
  query CurrentUser {
    user {
      id
      name
      ...AvatarUser
    }
  }
  ${AVATAR_USER}
`;

export function CurrentUser() {
  const [{ data }] = useQuery({ query: CURRENT_USER });

  return (
    &lt;div>
      &lt;Avatar user={data?.user} />
      {data?.user?.name}
    &lt;/div>
  );
}</code></pre>

      <p>
        Since the fragment is now owned by the code that makes use of it, we can extend the
        behaviour of the component without having to also make changes to another file.
      </p>

      <pre><code data-highlight class="language-diff">  import { gql } from "urql";
+ import { initials } from "./initials";

  export const AVATAR_USER = gql`
    fragment AvatarUser on User {
      id
      avatarUrl
+     name # needed because it's used...
    }
  `;

  export function Avatar({ user }) {
+   if (user && !user.avatarUrl) {
+     // ... here!
+     return &lt;div>{initials(user.name)}&lt;/div>
+   }
    return &lt;img src={user?.avatarUrl ?? "/default_user.svg"} />;
  }</code></pre>

      <p>
        Additionally, this pattern is composable, since fragments can be spread into other
        fragments.
      </p>

      <pre><code data-highlight class="language-typescript">import { gql } from "urql";
import { Avatar, AVATAR_USER } from "./Avatar";

export const USER_HEADER = gql`
  fragment UserHeader on User {
    id
    name
    ...AvatarUser
  }
  ${AVATAR_USER}
`;

export function UserHeader({ user }) {
  return (
    &lt;div>
      &lt;Avatar user={user} />
      {user?.name}
    &lt;/div>
  );
}</code></pre>

      <p>
        So, again, why does this work so well? Because
        <em>the queries are intimately related to the usage of the fields inside the components</em
        >. We can now decide to extract a new component, for instance, and all selections needed in
        the extracted code that is no longer needed in the remaining code, gets extracted to the new
        component's corresponding new fragment. The queries/fragments follow the UI.
      </p>

      <pre><code data-highlight class="language-diff">  import { useQuery, gql } from "urql";
- import { Avatar, AVATAR_USER } from "./Avatar";
+ import { UserHeader, USER_HEADER } from "./UserHeader";

  const CURRENT_USER = gql`
    query CurrentUser {
      user {
        id
-       name
-       ...AvatarUser
+       ...UserHeader
      }
    }
-   ${AVATAR_USER}
+   ${USER_HEADER}
  `;

  export function CurrentUser() {
    const [{ data }] = useQuery({ query: CURRENT_USER });

    return (
      &lt;div>
-       &lt;Avatar user={data?.user} />
-       {data?.user?.name}
+       &lt;UserHeader user={data?.user} />
      &lt;/div>
    );
  }</code></pre>

      <p>
        But as we’ve established, mutations sort of go the other way. So, it turns out that
        extracting mutations to reusable hooks is a perfectly reasonable idea! Because the answer to
        “what fields should we select” isn’t dependent on the usage in the UI, but on the fact that
        the fields are updated by the mutation itself; something we want to record into the cache.
      </p>

      <pre><code data-highlight class="language-typescript">import { useMutation, gql } from "urql";

const RENAME_USER = gql`
  mutation RenameUser($newName: String!) {
    renameUser(name: $newName) {
      id
      name
    }
  }
`;

export function useRenameUser() {
  const [, execute] = useMutation(RENAME_USER);
  return (newName) => execute({ newName });
}</code></pre>

      <p>
        I’m drawing a parallel to
        <a href="https://martinfowler.com/bliki/CQRS.html"
          >CQRS (Command Query Responsibility Separation)</a
        >
        here. It doesn’t translate perfectly, but the core idea is the same – commands (mutations)
        can be issued by code that is explicitly separated from UI requirements, while queries use a
        fully read-oriented representation (the graph) that don’t at all take into account how or
        why the data would change (i.e. command requirements).
      </p>

      <h3>Subscriptions</h3>

      <p>
        Surprisingly, GraphQL Subscriptions have more in common with mutations than with queries.
        Even though subscriptions only read data and doesn’t mutate in any way, a given subscription
        event carries with it the information that the graph has somehow changed – just like the
        mutation result does.
      </p>

      <pre><code data-highlight class="language-typescript">import { useSubscription, gql } from "urql";

const USER_WAS_RENAMED = gql`
  subscription UserWasRenamed {
    userWasRenamed {
      id
      name
    }
  }
`;

export function useUserWasRenamed() {
  useSubscription(USER_WAS_RENAMED);
}</code></pre>

      <h2>Part 3. Summary</h2>

      <p>
        This article promotes a few patterns when using a component-oriented UI library together
        with a GraphQL client with a normalized cache. Here they are:
      </p>

      <ol>
        <li>
          The field selections on GraphQL queries should be declared close to where they are used,
          and should not be reused through centralized hooks.
        </li>

        <li>
          Components with corresponding fragments compose, which makes the decision of where in the
          component tree to actually dispatch the query easier to make and change.
        </li>

        <li>
          Mutations and subscriptions, on the other hand, can be centralized in hooks, and benefit
          from reuse. They should only be concerned with selecting the fields that need to be
          updated in the cache as an effect of the mutation being applied or the subscription event
          arriving.
        </li>
      </ol>

      <p>
        Another way of thinking about why these patterns make sense, is to ask the question: “what
        reasons do these field selections have to change?” If we made reusable hooks from queries,
        the answer would be “anytime any consumer of this hook changes what fields they depend on.”
        That’s a requirement that’s very hard to maintain. The same goes if you have one big query
        on the top level of a component tree. Anytime any of the descendant components change their
        data usage, we need to update that big nasty query (and make sure not to inadvertently break
        another component!) If we instead have fragments next to our components, the answer becomes
        “if I need to change this component,” which is much easier to maintain, and limits the blast
        radius to a single file.
      </p>

      <p>
        Conversely, with mutations/subscriptions, if we don’t centralize and the mutation changes
        how it affects the graph, we now have to make sure to update all the places where the
        mutation is invoked to include the new changes in each selection. If we do centralize, we
        only have to do it in one place.
      </p>
    </div>
  </BlogPost>
</template>

<script lang="ts" setup>
import BlogPost from '../BlogPost.vue'
import hljs from 'highlight.js/lib/core'
import ts from 'highlight.js/lib/languages/typescript'
import xml from 'highlight.js/lib/languages/xml'
import diff from 'highlight.js/lib/languages/diff'
import 'highlight.js/styles/atom-one-dark.css'
import ogImage from './example.png'
import { computed, onMounted, ref } from 'vue'

hljs.registerLanguage('typescript', ts)
hljs.registerLanguage('diff', diff)
hljs.registerLanguage('xml', xml)

let ogImageFull = computed(() => new URL(ogImage, import.meta.env.VITE_WEB_HOST).href)

let container = ref(null)

onMounted(() => {
  const highlightBlocks = container.value?.querySelectorAll('[data-highlight]') ?? []

  for (const block of highlightBlocks) {
    hljs.highlightElement(block)
  }
})
</script>
