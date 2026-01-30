/* eslint-disable @next/next/no-before-interactive-script-outside-document */
import Script from "next/script"

const themeInitScript =
  "(function(){try{var t=localStorage.getItem('color-theme');if(t){document.documentElement.setAttribute('data-theme',t);}}catch(e){}})();"

export function ColorThemeInit() {
  return (
    <Script id="color-theme-init" strategy="beforeInteractive">
      {themeInitScript}
    </Script>
  )
}
