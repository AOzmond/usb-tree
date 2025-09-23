<script lang="ts">
  import { Refresh } from "../../../wailsjs/go/main/App"
  import { deviceLogs, theme, toggleTheme } from "../../lib/state.svelte"
  import { RefreshCcw, ToggleLeft } from "@lucide/svelte"

  function formatTimestamp(time: Date) {
    return new Date(time).toLocaleTimeString()
  }

  let lastLog = $derived($deviceLogs?.length ? $deviceLogs[$deviceLogs.length - 1] : null)
  let lastUpdatedTimestamp = $derived(lastLog ? formatTimestamp(lastLog.Time) : null)
  let isRefreshing = $state(false)
  let refreshCompleted = $state(false)

  async function handleRefresh() {
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

  function handleSpinEnd() {
    if (refreshCompleted) {
      isRefreshing = false
    }
  }
</script>

<div class="header">
  <span><b>Last updated:</b> {lastUpdatedTimestamp}</span>
  <div class="header__actions">
    <button
      type="button"
      class="header__theme-button"
      onclick={toggleTheme}
      aria-label={`Switch to ${$theme === "dark" ? "light" : "dark"} theme`}
      aria-pressed={$theme === "dark"}
    >
      <ToggleLeft class="header__theme-icon" /></button
    >
    <button
      type="button"
      class="header__refresh-button"
      class:header__refresh-button--spinning={isRefreshing}
      onclick={handleRefresh}
      aria-label="Refresh"
    >
      <RefreshCcw class="header__refresh-icon" onanimationend={handleSpinEnd} /></button
    >
  </div>
</div>

<style lang="scss">
@use '@carbon/styles/scss/spacing';

  .header {
    margin-bottom: 1px;
    box-sizing: border-box;
    padding: 0 spacing.$spacing-04;
    box-shadow: 0 1px var(--color-divider);
    height: spacing.$spacing-10;
    background: var(--color-header-bg);
    width: 100%;
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    align-items: center;
    flex-shrink: 0;
    flex-grow: 0;
  }

  .header__actions {
    height: spacing.$spacing-08;
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .header__theme-button,
  .header__refresh-button {
    border-radius: spacing.$spacing-03;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    height: spacing.$spacing-08;
    width: spacing.$spacing-08;
  }

  .header__theme-button:focus-visible,
  .header__refresh-button:focus-visible {
    outline: spacing.$spacing-01 solid var(--color-divider);
    outline-offset: spacing.$spacing-01;
  }

  .header__theme-button:hover,
  .header__refresh-button:hover {
    background: var(--color-divider);
  }

  .header :global(svg) {
    height: spacing.$spacing-08;
    width: auto;
    stroke: var(--color-text);
  }

  .header :global(.header__theme-icon circle),
  .header :global(.header__theme-icon rect) {
    transition:
      transform 0.25s ease,
      fill 0.25s ease,
      stroke 0.25s ease;
    transform-box: fill-box;
    transform-origin: center;
  }

  .header :global(.header__theme-icon circle) {
    fill: var(--color-text);
  }

  .header :global(.header__theme-button[aria-pressed="true"] .header__theme-icon circle) {
    transform: translateX(spacing.$spacing-03 - spacing.$spacing-01);
  }

  .header :global(.header__refresh-icon) {
    height: spacing.$spacing-07 - spacing.$spacing-02;
  }

  @keyframes spin {
    from {
      transform: rotate(0deg);
    }
    to {
      transform: rotate(-360deg);
    }
  }

  .header__refresh-button--spinning :global(.header__refresh-icon) {
    animation: spin 0.45s ease-in-out forwards;
  }
</style>
