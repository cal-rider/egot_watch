import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
    "./components/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        // Old Hollywood palette
        gold: {
          50: "#fffef7",
          100: "#fef9e7",
          200: "#fdf0c3",
          300: "#fce38d",
          400: "#f9d054",
          500: "#d4af37", // Classic gold
          600: "#b8960f",
          700: "#936f0a",
          800: "#7a5a0d",
          900: "#674a12",
        },
        hollywood: {
          black: "#0a0a0a",
          charcoal: "#1a1a1a",
          dark: "#2d2d2d",
          cream: "#f5f0e6",
          burgundy: "#722f37",
        },
      },
      fontFamily: {
        display: ["Playfair Display", "Georgia", "serif"],
        body: ["Inter", "system-ui", "sans-serif"],
      },
      backgroundImage: {
        "art-deco": "url('/art-deco-pattern.svg')",
        "spotlight": "radial-gradient(ellipse at center top, rgba(212,175,55,0.15) 0%, transparent 50%)",
      },
    },
  },
  plugins: [],
};
export default config;
