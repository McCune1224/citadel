/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./web/components/**/*.templ",
    "./cmd/**/*.templ",
  ],
  theme: {
    extend: {
      colors: {
        // Valve Half-Life 2 color palette
        'valve-orange': '#FF8C00',
        'valve-dark': '#0D0D0D',
        'valve-cyan': '#00FFFF',
        'valve-green': '#00FF00',
        'dark': '#0D0D0D',
      },
      fontFamily: {
        mono: ['Courier New', 'monospace'],
      },
    },
  },
  plugins: [],
}
