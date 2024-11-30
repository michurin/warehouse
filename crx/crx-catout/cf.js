// avito confluence
//
// Stolen from here https://github.com/guillesotelo/confluence-dark-mode.git

(() => {
  const styles = `
/*
The content of the css file will be copied
to the styles used by the script.
Its purpose is just to have a separate
file to pretify and debug stylesheet.
*/

/* --- MAIN COLOR BACKGROUND --- */
html {
    filter: invert(84%) grayscale(0%) saturate(110%) contrast(115%) brightness(110%) hue-rotate(180deg) !important;
}

html * {
    box-shadow: none !important;
}

header,
#header,
.aui-header {
    background-color: #daebff !important;
}

#header-precursor,#header-precursor * {
    background-color: #f6faff !important;
}

img,
svg,
video,
div[role=img],
*[style*="background-image"] {
    filter: invert(84%) grayscale(0%) saturate(110%) contrast(115%) brightness(110%) hue-rotate(180deg) !important;
}

.aui-header .aui-quicksearch input[type='text'], .aui-header .aui-quicksearch input[type='text'][type='text']:focus {
    background: #8e9cb3 !important;
}

/* --- SCROLLBARS --- */
::-webkit-scrollbar {
    width: 8px !important;
    height: 8px !important;
}

::-webkit-scrollbar-thumb {
    background-color: gray !important;
    border-radius: 6px !important;
}

::-webkit-scrollbar-thumb:hover {
    background-color: darkgray !important;
}

::-webkit-scrollbar-track {
    background: lightgray !important;
}
`

  const styleElement = document.createElement('style');
  styleElement.textContent = styles;
  document.head.appendChild(styleElement);
})();
