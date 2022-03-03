<template>
  <div class="sm:px-20 flex flex-col items-center space-y-6">
    <div class="flex flex-col items-center space-y-4 mb-10">
      <div class="flex flex-col items-center">
        <div :class="[fromColorStart, fromColor, toColor]" class="bg-gradient-to-b w-px h-20"></div>
        <div
          :class="[bgColor, glow]"
          class="text-slate-900 font-bold text-l rounded-full w-10 h-10 flex justify-center items-center"
        >
          {{ stepNumber }}
        </div>
      </div>

      <h1 :class="[textAccentColor]" class="text-3xl font-bold tracking-tight text-center">
        {{ title }}
      </h1>

      <h1 class="text-slate-50 text-5xl font-bold tracking-tight text-center">
        {{ subTitle }}
      </h1>

      <p class="md:w-2/3 text-center">
        <slot name="step-top-text"></slot>
      </p>
    </div>

    <div class="mt-24 flex flex-col">
      <div class="flex flex-col md:flex-row space-y-8 md:space-y-0">
        <div class="md:h-[30rem] md:py-0 md:-my-4 grow-0 px-20 flex flex-col items-center">
          <div :class="[fromColor, toColor]" class="bg-gradient-to-b w-px h-24 md:h-1/3"></div>
          <div
            :class="[borderColor, glow]"
            class="bg-slate-900 rounded-full border-2 w-3 h-3"
          ></div>
          <div :class="[bgColor]" class="w-px h-0 md:h-2/3"></div>
        </div>
        <div class="md:grow md:w-1/3 flex items-left justify-center md:justify-start md:flex-col">
          <div class="h-1/3"></div>
          <div class="-my-4 h-2/3">
            <slot name="1-right"></slot>
          </div>
        </div>
        <div class="md:grow md:w-1/3 md:order-first">
          <slot name="1-left"></slot>
        </div>
      </div>

      <div class="flex flex-col md:flex-row space-y-8 md:space-y-0">
        <div class="md:py-0 md:-my-4 grow-0 px-20 flex flex-col items-center">
          <div
            :class="[mdBgColor, fromColor, toColor]"
            class="w-px h-24 md:h-1/3 bg-gradient-to-b"
          ></div>
          <div
            :class="[borderColor, glow]"
            class="bg-slate-900 rounded-full border-2 w-3 h-3"
          ></div>
          <div
            :class="[fromColorEnd, toColorNext]"
            class="bg-gradient-to-b w-px h-0 md:h-2/3"
          ></div>
        </div>
        <div class="md:grow md:w-1/3 flex items-left justify-center md:justify-start md:flex-col">
          <div class="h-1/3"></div>
          <div class="-my-4 h-2/3">
            <slot name="2-right"></slot>
          </div>
        </div>
        <div class="md:grow md:w-1/3 md:order-first">
          <slot name="2-left"></slot>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
export default {
  props: {
    stepNumber: {
      type: Number,
      required: true,
    },
    title: {
      type: String,
      required: true,
    },
    subTitle: {
      type: String,
      required: true,
    },
    bgColor: {
      type: String,
      required: true,
    },
    mdBgColor: {
      type: String,
      required: true,
    },
    borderColor: {
      type: String,
      required: true,
    },
    textAccentColor: {
      type: String,
      required: true,
    },
    fromColor: {
      type: String,
      required: true,
    },
    fromColorStart: {
      type: String,
      required: true,
    },
    toColor: {
      type: String,
      required: true,
    },
    fromColorEnd: {
      type: String,
      required: true,
    },
    toColorNext: {
      type: String,
      required: true,
    },
    glow: {
      type: String,
      required: true,
    },
  },
}
</script>
