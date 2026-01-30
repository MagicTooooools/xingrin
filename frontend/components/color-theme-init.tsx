/* eslint-disable @next/next/no-before-interactive-script-outside-document */
import Script from "next/script"
import {
  COLOR_THEMES,
  COLOR_THEME_COOKIE_KEY,
  DEFAULT_COLOR_THEME_ID,
} from "@/lib/color-themes"

const themeIds = COLOR_THEMES.map((theme) => theme.id)
const darkThemeIds = COLOR_THEMES.filter((theme) => theme.isDark).map((theme) => theme.id)

const themeInitScript = `(function(){try{
  var key=${JSON.stringify(COLOR_THEME_COOKIE_KEY)};
  var theme=null;
  var match=document.cookie.match(new RegExp('(?:^|; )'+key+'=([^;]*)'));
  if(match){theme=decodeURIComponent(match[1]);}
  if(!theme){theme=localStorage.getItem(key);}
  var valid=${JSON.stringify(themeIds)};
  if(valid.indexOf(theme)===-1){theme=${JSON.stringify(DEFAULT_COLOR_THEME_ID)};}
  try{
    document.cookie=key+'='+encodeURIComponent(theme)+'; Path=/; Max-Age=${60 * 60 * 24 * 365 * 2}; SameSite=Lax';
  }catch(e){}
  var dark=${JSON.stringify(darkThemeIds)};
  var root=document.documentElement;
  root.setAttribute('data-theme',theme);
  if(dark.indexOf(theme)!==-1){root.classList.add('dark');}else{root.classList.remove('dark');}
}catch(e){}})();`

export function ColorThemeInit() {
  return (
    <Script id="color-theme-init" strategy="beforeInteractive">
      {themeInitScript}
    </Script>
  )
}
