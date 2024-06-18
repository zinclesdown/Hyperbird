/** @type {import('tailwindcss').Config} */
export default {
  // prefix: 'tw-', // 不打算加前缀了，都一样用
  css: {
    preFetch: true,
  },
  content: ['./src/**/*.{html,js,vue,ts,jsx,tsx}'],
  theme: {
    extend: {},
  },
  plugins: [],
  corePlugins: {
    preflight: false,
  },
};
