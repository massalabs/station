/** @type {import('tailwindcss').Config} */
const { createThemes } = require('tw-colors');
const plugin = require('tailwindcss/plugin')
module.exports = {
    content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],
    theme: {
        extend: {
            fontFamily: {
              MaisonNeue: ["Maison Neue", "sans-serif"],
              TTCommons: ["TT Commons", "sans-serif"],
            },
            keyframes:{
                rotate180Smooth:{
                    "0%": { transform: "rotate(0deg)" },
                    "12%": { transform: "rotate(22deg)" },
                    "25%": { transform: "rotate(45deg)" },
                    "37%": { transform: "rotate(67deg)" },
                    "50%": { transform: "rotate(90deg)" },
                    "62%": { transform: "rotate(112deg)" },
                    "75%": { transform: "rotate(135deg)" },
                    "87%": { transform: "rotate(157deg)" },
                    "100%": { transform: "rotate(180deg)" },
                }
            },
            animation:{
                rotate180Smooth: "rotate180Smooth 0.5s forwards ",
            }
        },
    },
    plugins: [
        require("@tailwindcss/typography"),
        require("daisyui"),
        createThemes({
            light: {
                primaryBG: "#F3F3F3",
                bgCard: "#FFFFFF",
                primaryButton: "#000000",
                secondaryButton: "#FFFFFF",
                hoverbgCard: "#EBEBEB",
                hoverPrimaryButton: "#D2D2D2",
                hoverSecondaryButton: "#21252A",
                brand: "#FF0000",
                live: "#14FF00",
                font: "#000000",
                invertedfont: "#FFFFFF",
                border: "#000000",
                disabledButton: "#8F8989",
                disabledText:"#D2D2D2"
            },
            dark: {
                primaryBG: "#171A1D",
                bgCard: "#21252A",
                primaryButton: "#FFFFFF",
                secondaryButton: "#21252A",
                hoverbgCard: "#2D333A",
                hoverPrimaryButton: "#2D333A",
                hoverSecondaryButton: "#D2D2D2",
                brand: "#FF0000",
                live: "#14FF00",
                border: "#F3F3F3",
                font: "#FFFFFF",
                invertedfont: "#000000",
                disabledButton: "#8F8989",
                disabledText:"#D2D2D2"
            },
        }),
        plugin(function ({ addComponents, theme }) {
            addComponents({
              ".display": {
                fontSize: "72px",
                fontWeight: "500",
                fontFamily: theme("fontFamily.MaisonNeue"),
                lineHeight: "87px",
                fontStyle: "normal",
            },
            ".label": {
                fontSize: "32px",
                fontWeight: "600",
                fontFamily: theme("fontFamily.MaisonNeue"),
                lineHeight: "39px",
            },
            ".label2": {
                fontSize: "24px",
                fontWeight: "600",
                fontFamily: theme("fontFamily.MaisonNeue"),
                lineHeight: "29px",
            },
              ".button": {
                fontSize: "18px",
                fontWeight: "600",
                fontFamily: theme("fontFamily.MaisonNeue"),
                lineHeight: "22px",
            },
            ".buttonUnderline": {
                fontSize: "18px",
                fontWeight: "600",
                fontFamily: theme("fontFamily.MaisonNeue"),
                lineHeight: "22px",
            },
            ".text": {
                fontSize: "24px",
                fontWeight: "400",
                fontFamily: theme("fontFamily.TTCommons"),
                lineHeight: "29px",
            },
            ".text2": {
                fontSize: "18px",
                fontWeight: "400",
                fontFamily: theme("fontFamily.TTCommons"),
                lineHeight: "22px",
            },
            ".text3": {
                fontSize: "14px",
                fontWeight: "400",
                fontFamily: theme("fontFamily.TTCommons"),
                lineHeight: "17px",
            },
            ".textUnderline": {
                fontSize: "18px",
                fontWeight: "400",
                fontFamily: theme("fontFamily.TTCommons"),
                lineHeight: "22px",
            },
            ".Secondary": {
                fontSize: "32px",
                fontWeight: "400",
                fontFamily: theme("fontFamily.TTCommons"),
                lineHeight: "38px",
            },
          });
        }),
    ],
};
