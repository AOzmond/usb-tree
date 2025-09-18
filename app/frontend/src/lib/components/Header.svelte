<script lang="ts">
  import { Refresh } from "../../../wailsjs/go/main/App"
  import { deviceLogs, theme, toggleTheme } from "../../lib/state.svelte"

  const formatTimestamp = (time: Date) =>
    new Date(time).toLocaleTimeString([], {
      hour: "2-digit",
      minute: "2-digit",
      second: "2-digit",
      hour12: false,
    })

  let lastLog = $derived($deviceLogs?.length ? $deviceLogs[$deviceLogs.length - 1] : null)
  let formattedTimestamp = $derived(lastLog ? formatTimestamp(lastLog.Time) : null)
  let isRefreshing = $state(false)
  let refreshCompleted = $state(false)

  const handleRefresh = async () => {
    if (isRefreshing) {
      return
    }
    isRefreshing = true
    refreshCompleted = false
    try {
      await Refresh()
    } finally {
      refreshCompleted = true
    }
  }

  const handleSpinEnd = () => {
    if (refreshCompleted) {
      isRefreshing = false
    }
  }
</script>

<div class="header">
  <span> last change: {formattedTimestamp}</span>
  <div class="icons">
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
        ><rect width="20" height="14" x="2" y="5" rx="7" /><circle cx="9" cy="12" r="6" /></svg
      ></button
    >
    <button onclick={handleRefresh} aria-label="Refresh" class:spinning={isRefreshing}
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
        class="lucide lucide-refresh-ccw-icon lucide-refresh-ccw refresh"
        onanimationend={handleSpinEnd}
        ><path d="M21 12a9 9 0 0 0-9-9 9.75 9.75 0 0 0-6.74 2.74L3 8" /><path d="M3 3v5h5" /><path
          d="M3 12a9 9 0 0 0 9 9 9.75 9.75 0 0 0 6.74-2.74L21 16"
        /><path d="M16 16h5v5" /></svg
      ></button
    >
  </div>
</div>

<style lang="scss">
  .header {
    margin-bottom: 1px;
    box-sizing: border-box;
    padding: 0 12px;
    box-shadow: 0 1px var(--color-divider);
    height: 64px;
    background: var(--color-header-bg);
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
    justify-content: center;
    align-items: center;
  }
  button {
    border-radius: 0.5rem;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    height: 40px;
    width: 40px;
  }
  button:focus-visible {
    outline: 2px solid var(--color-divider);
    outline-offset: 2px;
  }
  button:hover {
    background: var(--color-divider);
  }
  svg {
    height: 40px;
    width: auto;
    stroke: var(--color-text);
    fill: var(--color-header-bg);
  }

  button svg circle,
  button svg rect {
    transition: transform 0.25s ease;
  }
  button[aria-pressed="true"] svg circle {
    transform: translateX(6px);
  }

  .refresh {
    height: 28px;
  }

  @keyframes spin {
    from {
      transform: rotate(0deg);
    }
    to {
      transform: rotate(-360deg);
    }
  }

  button.spinning .refresh {
    animation: spin 0.9s ease-in-out forwards;
  }
</style>
