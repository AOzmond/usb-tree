<script>
  import { Refresh } from "../../../wailsjs/go/main/App"
  import { theme, toggleTheme } from "../../lib/state.svelte"

  // TODO
  let lastTimeStamp = "some timestamp"
</script>

<div class="header">
  <span>last updated: {lastTimeStamp} </span>
  <span class="icons">
    <button
      onclick={toggleTheme}
      aria-label={`Switch to ${$theme === "dark" ? "light" : "dark"} theme`}
      aria-pressed={$theme === "dark"}
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        width="24"
        height="24"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
        class="lucide lucide-toggle-left-icon lucide-toggle-left"
        ><circle cx="9" cy="12" r="3" /><rect width="20" height="14" x="2" y="5" rx="7" /></svg
      ></button
    >
    <button onclick={Refresh} aria-label="Refresh"
      ><svg
        xmlns="http://www.w3.org/2000/svg"
        width="24"
        height="24"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
        class="lucide lucide-refresh-ccw-icon lucide-refresh-ccw"
        ><path d="M21 12a9 9 0 0 0-9-9 9.75 9.75 0 0 0-6.74 2.74L3 8" /><path d="M3 3v5h5" /><path
          d="M3 12a9 9 0 0 0 9 9 9.75 9.75 0 0 0 6.74-2.74L21 16"
        /><path d="M16 16h5v5" /></svg
      ></button
    >
  </span>
</div>

<style lang="scss">
  .header {
    margin-bottom: 1px;
    box-sizing: border-box;
    padding: 0 12px;
    box-shadow: 0 1px var(--color-divider);
    height: 64px;
    background: var(--color-header-bg);
    color: var(--color-text);
    width: 100%;
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    align-items: center;
    flex-shrink: 0;
    flex-grow: 0;
  }
  .icons {
    height: 40px;
    display: flex;
    gap: 0.5rem;
  }
  button {
    all: unset;
    box-sizing: border-box;
    height: 100%;
    color: inherit;
    padding: 0 0.25rem;
    border-radius: 0.5rem;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
  }
  button:focus-visible {
    outline: 2px solid var(--color-divider);
    outline-offset: 2px;
  }
  button:hover {
    background: var(--color-divider);
  }
  svg {
    height: 100%;
    width: auto;
  }
</style>
