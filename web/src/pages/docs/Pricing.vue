<template>
  <DocumentationWithTableOfContents>
    <template #sidebar>
      <DocsSidebar />
    </template>
    <template #default>
      <div class="prose p-4 max-w-[800px]">
        <h1 id="pricing">Pricing</h1>
        <p>
          Sturdy comes in a few variants, with different pricing.
          <router-link :to="{ name: 'v2pricing' }">See full comparison</router-link>
          .
        </p>

        <h2 id="open-source-sturdy">Self-hosted Sturdy Open Source</h2>
        <p>
          Self-hosted, and free forever! [<a href="https://github.com/sturdy-dev/sturdy"
            >GitHub Mirror</a
          >]
        </p>

        <h2 id="sturdy-cloud">Sturdy Cloud</h2>
        <p>Sturdy-the-product, hosted by Sturdy-the-company.</p>
        <ul>
          <li>10 users per organization is free</li>
          <li>$30 / user / month (billed monthly)</li>
        </ul>

        <h2 id="self-hosted-sturdy-enterprise">Self-hosted Sturdy Enterprise</h2>
        <ul>
          <li>Up to 10 users for free</li>
          <li>$30 / user / month (billed annually, minimum 20 seats)</li>
          <li>
            Extra used seats (outside of seats included in the license) are pro-rated (billed
            monthly)
          </li>
          <li>
            Dedicated Slack-support, hosting support, and remote debugging and configuration by
            Sturdy engineers is <em>included for paying customers</em>.
          </li>
          <li>Community-support is available for all users</li>
        </ul>
      </div>
    </template>
  </DocumentationWithTableOfContents>
</template>

<script lang="ts" setup>
import DocumentationWithTableOfContents from '../../layouts/DocumentationWithTableOfContents.vue'
import { useHead } from '@vueuse/head'
import DocsSidebar from '../../organisms/docs/DocsSidebar.vue'

useHead({
  meta: [
    {
      name: 'description',
      content: 'Pricing for Sturdy open source, enterprise, and Sturdy in the cloud.',
    },
    {
      name: 'keywords',
      content: 'study learn documentation pricing seats licenses enterprise in the cloud',
    },
  ],
  title: 'Pricing | Sturdy',
})
</script>
