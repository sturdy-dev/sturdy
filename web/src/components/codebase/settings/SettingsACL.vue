<template>
  <div v-if="writeACL && writeACL.canI">
    <HorizontalDivider class="mt-4" bg="bg-white">Access Control</HorizontalDivider>

    <div class="mx-4 divide-gray-200">
      <p class="mt-1 text-sm text-gray-500">
        The JSON-based policy defines granular access control for codespace members.
        <a class="text-blue-500 underline" href="/docs/access-control" target="_blank">
          Learn more
        </a>
      </p>

      <div v-if="errors.length > 0">
        <Banner v-for="err in errors" :key="err" status="error" class="my-2" :message="err" />
      </div>

      <div class="pt-2">
        <prism-editor
          v-model="newPolicy"
          class="max-h-96 p-5 leading-normal text-base font-mono shadow-sm sm:text-sm border-gray-300 rounded-md bg-gray-50"
          :highlight="highlighter"
          line-numbers
        ></prism-editor>

        <span class="pt-2 relative z-0 inline-flex shadow-sm rounded-md">
          <Button :grouped="true" :first="true" @click="onSaveClicked">Save</Button>
          <Button :grouped="true" :last="true" @click="onResetClicked">Reset</Button>
        </span>
      </div>
    </div>
  </div>
</template>

<script>
import { PrismEditor } from 'vue-prism-editor'
import 'vue-prism-editor/dist/prismeditor.min.css'

import { highlight, languages } from 'prismjs/components/prism-core'
import 'prismjs/components/prism-json'
import 'prismjs/themes/prism-tomorrow.css'

import { gql, useMutation, useQuery } from '@urql/vue'
import { toRef, ref, watch } from 'vue'

import HorizontalDivider from '../../shared/HorizontalDivider.vue'
import Button from '../../shared/Button.vue'
import Banner from '../../shared/Banner.vue'

export default {
  name: 'SettingsACL',
  components: {
    Banner,
    Button,
    HorizontalDivider,
    PrismEditor,
  },
  props: {
    codebaseId: String,
    aclId: String,
    aclPolicy: String,
  },
  setup(props) {
    const newPolicy = ref('')
    const aclPolicy = toRef(props, 'aclPolicy')

    newPolicy.value = aclPolicy.value
    watch(aclPolicy, (updatedPolicy) => {
      if (updatedPolicy) {
        newPolicy.value = updatedPolicy
      }
    })

    const { executeMutation: updateACLPolicyResult } = useMutation(gql`
      mutation SettingsACLUpdate($codebaseID: ID!, $policy: String) {
        updateACL(input: { codebaseID: $codebaseID, policy: $policy }) {
          id
          policy
        }
      }
    `)

    const aclId = toRef(props, 'aclId')
    const codebaseId = toRef(props, 'codebaseId')

    const canIWriteACLQuery = useQuery({
      query: gql`
        query SettingsACL($codebaseID: ID!, $action: String!, $resource: String!) {
          canI(codebaseID: $codebaseID, action: $action, resource: $resource)
        }
      `,
      variables: {
        codebaseID: codebaseId.value,
        action: 'write',
        resource: `acls::${aclId.value}`,
      },
    })

    return {
      newPolicy,
      updateACLPolicyResult,
      writeACL: canIWriteACLQuery.data,
    }
  },
  data: () => ({ errors: [] }),
  methods: {
    clearErrors() {
      this.errors.splice(0)
    },

    showError(err) {
      this.errors.push(err)
    },

    onResetClicked() {
      this.newPolicy = this.aclPolicy
      this.clearErrors()
    },

    onSaveClicked() {
      this.clearErrors()

      this.updateACLPolicyResult({
        codebaseID: this.codebaseId,
        policy: this.newPolicy,
      }).then((result) => {
        if (result.error && result.error.graphQLErrors) {
          result.error.graphQLErrors.forEach((gqlError) => {
            for (const key in gqlError.extensions) {
              this.showError(`${key}: ${gqlError.extensions[key]}`)
            }
          })
        } else if (result.error) {
          throw new Error(result.error)
        } else {
          this.emitter.emit('notification', {
            title: 'Policy updated',
            message: 'ACL Policy has been updated',
          })
        }
      })
    },

    highlighter(code) {
      return highlight(code, languages.json)
    },
  },
}
</script>

<style>
.prism-editor__line-numbers {
  user-select: none;
}
</style>
