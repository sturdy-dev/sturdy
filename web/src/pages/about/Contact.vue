<template>
  <Documentation>
    <div>
      <div class="mx-auto py-12 px-4 max-w-7xl sm:px-6 lg:px-8 lg:py-24 space-y-24">
        <div class="space-y-12">
          <h2 class="text-3xl font-bold sm:text-4xl">Contact us</h2>
          <div class="max-w-2xl">
            <p class="text-xl text-gray-500 pb-4">Sturdy is built by developers for developers.</p>

            <p class="text-gray-500 my-2">
              Don't hesitate to get in touch with us if you have any questions
            </p>

            <div class="mt-4">
              <ul class="space-y-2">
                <li>
                  <LinkButton
                    href="https://discord.com/invite/5HnSdzMqtA"
                    target="_blank"
                    class="!no-underline !text-black !m-0"
                  >
                    <img src="../static/docs/assets/discord.svg" class="h-8 w-8 mr-2" />
                    <span>Get help on our Discord</span>
                  </LinkButton>
                </li>
                <li>
                  <LinkButton
                    href="mailto:support@getsturdy.com"
                    class="!no-underline !text-black !m-0"
                  >
                    <AtSymbolIcon class="h-8 w-8 mr-2" />
                    <span>support@getsturdy.com</span>
                  </LinkButton>
                </li>
                <li>
                  <LinkButton href="tel:+46706389084" class="!no-underline !text-black !m-0">
                    <PhoneIcon class="h-8 w-8 mr-2" />
                    <span>Call us on +46 (58) 24 58 944</span>
                  </LinkButton>
                </li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </div>
  </Documentation>
</template>

<script lang="ts" setup>
import { useHead } from '@vueuse/head'
import LinkButton from '../../components/shared/LinkButton.vue'
import { AtSymbolIcon, PhoneIcon } from '@heroicons/vue/solid'
import Documentation from '../../layouts/Documentation.vue'

useHead({
  title: 'Contact Sturdy',
})
</script>
