const tailwindcss = require('tailwindcss');
module.exports = {
    plugins: [
		tailwindcss('./examples/app/tailwind.config.js'),
        require('autoprefixer'),
    ],
};