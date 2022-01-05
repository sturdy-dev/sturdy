/* eslint-disable @typescript-eslint/no-var-requires */

const colors = require('tailwindcss/colors')

module.exports = {
  content: ['src/**/*.{vue,ts}', 'src/*.{vue,ts}'],
  mode: 'jit',
  theme: {
    extend: {
      colors: {
        gray: colors.neutral,
        orange: colors.orange,
        warmgray: colors.warmGray,
        green: colors.emerald,
        yellow: colors.amber,
        purple: colors.violet,
      },
    },
  },
  plugins: [
    require('@tailwindcss/typography'),
    require('@tailwindcss/forms'),
    require('@tailwindcss/line-clamp'),
    require('@tailwindcss/aspect-ratio'),
  ],
}
