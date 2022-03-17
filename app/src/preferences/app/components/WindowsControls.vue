<template>
  <div id="window-controls">
    <div id="min-button" class="button" @click.prevent="onMinimize">
      <svg class="icon" width="11" height="1" viewBox="0 0 11 1">
        <path d="m11 0v1h-11v-1z" stroke-width=".26208" />
      </svg>
    </div>

    <div v-if="!isMaximized" id="max-button" class="button" @click.prevent="onMaximize">
      <svg class="icon" width="10" height="10" viewBox="0 0 10 10">
        <path d="m10-1.6667e-6v10h-10v-10zm-1.001 1.001h-7.998v7.998h7.998z" stroke-width=".25" />
      </svg>
    </div>

    <div v-if="isMaximized" id="restore-button" class="button" @click.prevent="onRestore">
      <svg class="icon" width="11" height="11" viewBox="0 0 11 11">
        <path
          d="m11 8.7978h-2.2021v2.2022h-8.7979v-8.7978h2.2021v-2.2022h8.7979zm-3.2979-5.5h-6.6012v6.6011h6.6012zm2.1968-2.1968h-6.6012v1.1011h5.5v5.5h1.1011z"
          stroke-width=".275"
        />
      </svg>
    </div>

    <div id="close-button" class="button" @click.prevent="onClose">
      <svg class="icon" width="12" height="12" viewBox="0 0 12 12">
        <path
          d="m6.8496 6 5.1504 5.1504-0.84961 0.84961-5.1504-5.1504-5.1504 5.1504-0.84961-0.84961 5.1504-5.1504-5.1504-5.1504 0.84961-0.84961 5.1504 5.1504 5.1504-5.1504 0.84961 0.84961z"
          stroke-width=".3"
        />
      </svg>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'

export default {
  setup() {
    const { ipc } = window

    const isMinimized = ref(false)
    ipc.isMinimized().then((is) => {
      isMinimized.value = is
    })

    const isMaximized = ref(false)
    ipc.isMaximized().then((is) => {
      isMaximized.value = is
    })

    const isNormal = ref(false)
    ipc.isNormal().then((is) => {
      isNormal.value = is
    })

    return {
      ipc,
      isMaximized,
      isMinimized,
      isNormal,
    }
  },
  methods: {
    onClose() {
      this.ipc.close()
    },
    onMinimize() {
      this.ipc.minimize().then(() => {
        this.isMinimized = true
        this.isMaximized = false
        this.isNormal = false
      })
    },
    onMaximize() {
      this.ipc.maximize().then(() => {
        this.isMinimized = false
        this.isMaximized = true
        this.isNormal = false
      })
    },
    onRestore() {
      this.ipc.unmaximize().then(() => {
        this.isMinimized = false
        this.isMaximized = false
        this.isNormal = true
      })
    },
  },
}
</script>

<style>
#window-controls {
  display: grid;
  grid-template-columns: repeat(3, 46px);
  -webkit-app-region: no-drag;
}

#window-controls .button {
  grid-row: 1 / span 1;
  display: flex;
  justify-content: center;
  align-items: center;
  width: 46px;
  height: 32px;
  user-select: none;
}
#window-controls .button:hover {
  background: rgba(0, 0, 0, 0.1);
}
#window-controls .button:active {
  background: rgba(0, 0, 0, 0.2);
}

#min-button {
  grid-column: 1;
}

#max-button,
#restore-button {
  grid-column: 2;
}

#close-button {
  grid-column: 3;
}
#close-button:hover {
  background: #e81123 !important;
}
#close-button:active {
  background: #f1707a !important;
}
#close-button:hover .icon {
  filter: invert(1);
}
#close-button:active .icon {
  filter: invert(1);
}
</style>
