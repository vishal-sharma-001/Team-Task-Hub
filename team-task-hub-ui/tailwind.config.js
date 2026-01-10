module.exports = {
  content: [
    "./index.html",
    "./src/**/*.{js,jsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: '#3b82f6',
        'primary-dark': '#1e40af',
        'primary-light': '#dbeafe',
      },
      fontSize: {
        xs: '0.8125rem',
        sm: '0.875rem',
        base: '0.9375rem',
        lg: '1rem',
        xl: '1.125rem',
        '2xl': '1.5rem',
        '3xl': '1.875rem',
        '4xl': '2.25rem',
      },
    },
  },
  plugins: [],
}
