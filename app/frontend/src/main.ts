import App from "./App.svelte"
import { mount } from "svelte"

import "@fontsource-variable/ibm-plex-sans"
import "@fontsource-variable/jetbrains-mono"
import "$style/reset.scss"
import "$style/index.scss"
import "$style/variables.scss"
import "carbon-components-svelte/css/all.css"

const app = mount(App, {
  target: document.getElementById("app")!,
})
export default app
