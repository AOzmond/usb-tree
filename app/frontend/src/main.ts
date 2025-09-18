import App from "./+layout.svelte";
import { mount } from "svelte";
import "@fontsource-variable/ibm-plex-sans";
import "@fontsource-variable/jetbrains-mono";

const app = mount(App, {
  target: document.getElementById("app")!,
});
export default app;
