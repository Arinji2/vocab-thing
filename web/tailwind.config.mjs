/** @type {import('tailwindcss').Config} */
export default {
  content: ["./src/**/*.{js,jsx,ts,tsx}"],
  theme: {
    extend: {
      colors: {
        brand: {
          text: "#e5f9fa",
          background: "#041518",
          primary: {
            DEFAULT: "#89dfe9",
            dark: "#0f2d30",
          },
          secondary: {
            DEFAULT: "#4a1882",
            dark: "#260e40",
          },
          accent: "#df4fd6",
        },
      },
    },
  },
};
