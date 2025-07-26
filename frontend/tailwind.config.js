/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  theme: {
    extend: {
      colors: {
        'electric-blue': '#2563EB',
        'deep-blue': '#1E40AF', 
        'lightning-yellow': '#FBBF24',
        'light-gray': '#F8FAFC',
        'dark-slate': '#0F172A',
        'success': '#10B981',
        'warning': '#F59E0B',
        'error': '#EF4444',
      },
      fontFamily: {
        'sans': ['Inter', 'system-ui', 'sans-serif'],
      },
    },
  },
  plugins: [],
}