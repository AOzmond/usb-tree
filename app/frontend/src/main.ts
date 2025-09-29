import App from "./App.svelte"
import { mount } from "svelte"

import "@fontsource-variable/ibm-plex-sans"
import "@fontsource-variable/jetbrains-mono"
import "carbon-components-svelte/css/all.css"

import "$style/reset.scss"
import "$style/index.scss"
import "$style/variables.scss"

const app = mount(App, {
  target: document.getElementById("app")!,
})
export default app
