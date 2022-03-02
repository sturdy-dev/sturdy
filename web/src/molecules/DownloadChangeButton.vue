<template>
  <div class="flex items-center space-x-2">
    <Button v-if="fetching" size="wider" disabled>
      <Spinner class="-ml-1 mr-2 h-5 w-5 text-gray-400" aria-hidden="true" />
      <span>Preparing</span>
    </Button>
    <Button v-else-if="!fetching && data && url" size="wider">
      <DownloadIcon class="-ml-1 mr-2 h-5 w-5 text-gray-400" aria-hidden="true" />
      <a :href="url" target="_blank">Download as Zip</a>
    </Button>
    <Button v-else size="wider" @click="onDownloadChange">
      <DownloadIcon class="-ml-1 mr-2 h-5 w-5 text-gray-400" aria-hidden="true" />
      <span>Download as Zip</span>
    </Button>
  </div>
</template>

<script lang="ts">
import DownloadIcon from '@heroicons/vue/outline/DownloadIcon'
import Button from '../components/shared/Button.vue'
import Spinner from '../components/shared/Spinner.vue'

import { useQuery, gql } from '@urql/vue'
import { ref, toRefs } from 'vue'
import { DeepMaybeRef } from '@vueuse/core'

import {
  ChangeDetailsDownloadZipQuery,
  ChangeDetailsDownloadZipQueryVariables,
} from './__generated__/DownloadChangeButton'

const DOWNLOAD_CHANGE_QUERY = gql`
  query ChangeDetailsDownloadZip($changeID: ID!) {
    change(id: $changeID) {
      id
      downloadZip {
        id
        url
      }
    }
  }
`

export default {
  components: { DownloadIcon, Button, Spinner },
  props: {
    changeId: {
      type: String,
      required: true,
    },
  },
  setup(props) {
    const { changeId } = toRefs(props)
    const pauseGenerateZip = ref(true)
    const { data, fetching } = useQuery<
      ChangeDetailsDownloadZipQuery,
      DeepMaybeRef<ChangeDetailsDownloadZipQueryVariables>
    >({
      query: DOWNLOAD_CHANGE_QUERY,
      pause: pauseGenerateZip,
      variables: {
        changeID: changeId,
      },
      requestPolicy: 'network-only',
    })
    return {
      data,
      fetching,

      onDownloadChange() {
        pauseGenerateZip.value = false
      },
    }
  },
  data() {
    return {
      url: undefined as undefined | string,
    }
  },
  watch: {
    data: function (data) {
      const url = data.change.downloadZip.url
      if (url && !this.url) {
        this.url = url
        window.open(url, '_blank')?.focus()
      }
    },
  },
}
</script>
