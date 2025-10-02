<script lang="ts">
  import { Refresh } from "$wailsjs/go/main/App"
  import { deviceLogs, getNextTheme, theme, toggleTheme, type CarbonTheme } from "$lib/state.svelte"
  import { formatTimestamp } from "$lib/utilities"

  import { Header, HeaderGlobalAction, HeaderUtilities } from "carbon-components-svelte"

  import { RefreshCcw, ToggleLeft } from "@lucide/svelte"

  const themeLabels: Record<CarbonTheme, string> = {
    g100: "G100",
    white: "White",
  }

  const refreshClass = $derived(isRefreshing ? " refresh-action-spinning" : "")

  let lastLog = $derived($deviceLogs?.length ? $deviceLogs[$deviceLogs.length - 1] : undefined)
  let lastUpdatedTimestamp = $derived(lastLog ? formatTimestamp(lastLog.Time) : formatTimestamp(new Date()))
  let isRefreshing = $state(false)
  let refreshResetTimer: ReturnType<typeof setTimeout> | undefined = undefined

  let nextTheme = $derived(getNextTheme($theme))
  let currentThemeLabel = $derived(themeLabels[$theme])
  let nextThemeLabel = $derived(themeLabels[nextTheme])
  let themeTone = $derived($theme === "g100" ? "dark" : "light")

  async function handleRefresh() {
    if (isRefreshing) {
      return
    }

    isRefreshing = true
    if (refreshResetTimer) {
      clearTimeout(refreshResetTimer)
      refreshResetTimer = undefined
    }

    Refresh()

    refreshResetTimer = setTimeout(() => {
      isRefreshing = false
      refreshResetTimer = undefined
    }, 450)
  }
</script>

<Header id="header" class="header" uiShellAriaLabel="USB tree status">
  <span class="label">Last updated:</span>
  <span class="timestamp">{lastUpdatedTimestamp}</span>
  <HeaderUtilities class="utilities">
    <HeaderGlobalAction
      class="theme-action"
      data-theme-tone={themeTone}
      iconDescription="Dark mode"
      kind="primary"
      icon={ToggleLeft}
      aria-label={`Switch to ${nextThemeLabel} theme (current ${currentThemeLabel})`}
      onclick={toggleTheme}
    />
    <HeaderGlobalAction
      class={`refresh-action ${refreshClass}`}
      aria-label="Refresh"
      icon={RefreshCcw}
      kind="primary"
      iconDescription="Refresh"
      tooltipAlignment="end"
      onclick={handleRefresh}
    />
  </HeaderUtilities>
</Header>

<style lang="scss">
  @use "variables.scss" as *;

  .label {
    color: var(--color-header-text);
    padding-right: $spacing-02;
    font-weight: 600;
  }

  .timestamp {
    color: var(--color-header-text);
    font-weight: 400;
  }

  :global(.theme-action :is(.bx--btn__icon, .lucide-icon) circle),
  :global(.theme-action :is(.bx--btn__icon, .lucide-icon) rect) {
    transition: 0.15s ease 0s;
    transform-box: fill-box;
    transform-origin: center;
  }

  :global(.theme-action[data-theme-tone="dark"] :is(.bx--btn__icon, .lucide-icon) circle) {
    transform: translateX($spacing-03);
  }

  :global(.refresh-action-spinning :is(.bx--btn__icon, .lucide-icon)) {
    animation-name: spin;
    animation-duration: 0.45s;
    animation-timing-function: ease-in-out;
    animation-delay: 0s;
    animation-iteration-count: 1;
    animation-direction: normal;
    animation-fill-mode: forwards;
    animation-play-state: running;
  }

  :global(#header) {
    padding: $spacing-03;
  }
</style>
