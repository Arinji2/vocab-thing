/** @type {import('tailwindcss').Config} */
export default {
  content: ["./src/**/*.{js,jsx,ts,tsx}"],
  theme: {
    extend: {
      letterSpacing: {
        small: "0.04rem",
        large: "0.08rem",
      },
      colors: {
        brand: {
          text: "#e5f9fa",
          offwhite: {
            DEFAULT: "#6b7280",
            dark: "#374151",
            light: "#9ca3af",
          },
          destructive: {
            DEFAULT: "#ef4444",
            dark: "#b91c1c",
            light: "#fa575f",
          },
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
