module.exports = {
    darkMode: 'class',
    content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
    plugins: [require('daisyui')],
    theme: {
        extend: {},
    },
    daisyui: {
        themes: ['light', 'dark'],
    },
};
