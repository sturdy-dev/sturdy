<template>
  <div ref="target" class="text-center my-8 mx-auto overflow-hidden text-2xl text-slate-100">
    <div>Have an idea</div>
    <div>
      <span :class="[targetIsVisible ? 'strike' : '', 'text-gray-500 truncate font-mono']">
        git pull
      </span>
    </div>
    <div>
      <span
        :class="[targetIsVisible ? 'strikedelay-1 strike' : '', 'text-gray-500 truncate font-mono']"
      >
        git checkout -b feature
      </span>
    </div>

    <div>Write the code</div>

    <div>
      <span
        :class="[targetIsVisible ? 'strikedelay-2 strike' : '', 'text-gray-500 truncate font-mono']"
      >
        git add .
      </span>
    </div>
    <div>
      <span
        :class="[targetIsVisible ? 'strikedelay-3 strike' : '', 'text-gray-500 truncate font-mono']"
      >
        git commit -m "My change"
      </span>
    </div>

    <div>
      <span
        :class="[
          targetIsVisible ? 'strikedelay-4 strike' : '',
          'text-gray-500 truncate font-mono hidden lg:inline-block',
        ]"
      >
        git push --set-upstream origin feature
      </span>
      <span
        :class="[
          targetIsVisible ? 'strikedelay-4 strike' : '',
          'text-gray-500 truncate font-mono lg:hidden',
        ]"
      >
        git push -u origin feature
      </span>
    </div>

    <div>
      <span
        :class="[targetIsVisible ? 'strikedelay-5 strike' : '', 'text-gray-500 truncate font-mono']"
      >
        Create Pull Request
      </span>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue'
import { useElementVisibility } from '@vueuse/core'

export default defineComponent({
  setup() {
    const target = ref(null)
    const targetIsVisible = useElementVisibility(target)

    return {
      target,
      targetIsVisible,
    }
  },
})
</script>

<style scoped>
@keyframes strike {
  0% {
    width: 0;
  }
  100% {
    width: 100%;
  }
}

.strike {
  position: relative;
}

.strikedelay-1::after {
  animation-delay: 0.5s;
}
.strikedelay-2::after {
  animation-delay: 1s;
}
.strikedelay-3::after {
  animation-delay: 1.5s;
}
.strikedelay-4::after {
  animation-delay: 2s;
}
.strikedelay-5::after {
  animation-delay: 2.5s;
}
.strikedelay-6::after {
  animation-delay: 3s;
}
.strike::after {
  content: ' ';
  position: absolute;
  top: 50%;
  left: 0;
  height: 3px;
  background: gray;
  animation-name: strike;
  animation-duration: 0.5s;
  animation-timing-function: linear;
  animation-iteration-count: 1;
  animation-fill-mode: forwards;
}

@media screen and (max-width: 992px) {
  .strike::after {
    height: 2px;
  }
}
</style>
