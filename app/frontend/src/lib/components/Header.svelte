<script lang="ts">
  import { Refresh } from "../../../wailsjs/go/main/App"
  import {
    deviceLogs,
    getNextTheme,
    theme,
    toggleTheme,
    type CarbonTheme,
  } from "../../lib/state.svelte"
  import { RefreshCcw, ToggleLeft } from "@lucide/svelte"

  function formatTimestamp(time: Date) {
    return new Date(time).toLocaleTimeString()
  }

  const themeLabels: Record<CarbonTheme, string> = {
    g100: "G100",
    white: "White",
  }


  let lastLog = $derived($deviceLogs?.length ? $deviceLogs[$deviceLogs.length - 1] : null)
  let lastUpdatedTimestamp = $derived(lastLog ? formatTimestamp(lastLog.Time) : null)
  let isRefreshing = $state(false)
  let refreshCompleted = $state(false)

  let nextTheme = $derived(getNextTheme($theme))
  let currentThemeLabel = $derived(themeLabels[$theme])
  let nextThemeLabel = $derived(themeLabels[nextTheme])


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
      aria-label={`Switch to ${nextThemeLabel} theme (current ${currentThemeLabel})`}
      title={`Switch to ${nextThemeLabel} theme`}
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
@use '../../style/variables.scss';

  .header {
    margin-bottom: 1px;
    box-sizing: border-box;
    padding: 0 variables.$spacing-04;
    box-shadow: 0 1px var(--color-divider);
    height: variables.$spacing-10;
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
    height: variables.$spacing-08;
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .header__theme-button,
  .header__refresh-button {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    height: variables.$spacing-08;
    width: variables.$spacing-08;
  }

  .header__theme-button:focus-visible,
  .header__refresh-button:focus-visible {
    outline: variables.$spacing-01 solid var(--color-divider);
    outline-offset: variables.$spacing-01;
  }

  .header__theme-button:hover,
  .header__refresh-button:hover {
    background: var(--color-divider);
  }

  .header :global(svg) {
    height: variables.$spacing-08;
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

  .header
    :global(.header__theme-button[data-theme-tone="dark"] .header__theme-icon circle) {
    transform: translateX(variables.$spacing-03 - variables.$spacing-01);
  }

  .header :global(.header__refresh-icon) {
    height: variables.$spacing-07 - variables.$spacing-02;
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
