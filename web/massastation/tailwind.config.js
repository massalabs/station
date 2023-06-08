/** @type {import('tailwindcss').Config} */
export default {
  content: [
    './index.html',
    './src/**/*.{js,ts,jsx,tsx}',
    './node_modules/@massalabs/react-ui-kit/src/**/*.{js,ts,jsx,tsx}',
  ],
  presets: [require('@massalabs/react-ui-kit/presets/massa-station-preset.js')],
};
