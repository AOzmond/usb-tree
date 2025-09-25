<script lang="ts">
  import { Refresh } from "../../../wailsjs/go/main/App"
  import { deviceLogs, getNextTheme, theme, toggleTheme, type CarbonTheme } from "../../lib/state.svelte"
  import { Header, HeaderGlobalAction, HeaderUtilities } from "carbon-components-svelte"
  import { RefreshCcw, ToggleLeft } from "@lucide/svelte"

  function formatTimestamp(time: Date) {
    return new Date(time).toLocaleTimeString()
  }

  const themeLabels: Record<CarbonTheme, string> = {
    g100: "G100",
    white: "White",
  }

  let lastLog = $derived($deviceLogs?.length ? $deviceLogs[$deviceLogs.length - 1] : null)
  let lastUpdatedTimestamp = $derived(lastLog ? formatTimestamp(lastLog.Time) : formatTimestamp(new Date()))
  let isRefreshing = $state(false)
  let refreshResetTimer: ReturnType<typeof setTimeout> | null = null

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
      refreshResetTimer = null
    }
    try {
      await Refresh()
    } finally {
      refreshResetTimer = setTimeout(() => {
        isRefreshing = false
        refreshResetTimer = null
      }, 450)
    }
  }
</script>

<Header class="header" uiShellAriaLabel="USB tree status">
  <span slot="company" class="header__label">Last updated:&nbsp;</span>
  <span slot="platform" class="header__timestamp">{lastUpdatedTimestamp}</span>
  <HeaderUtilities class="header__utilities">
    <HeaderGlobalAction
      class="header__theme-action"
      data-theme-tone={themeTone}
      iconDescription="Dark mode"
      kind="primary"
      icon={ToggleLeft}
      aria-label={`Switch to ${nextThemeLabel} theme (current ${currentThemeLabel})`}
      on:click={toggleTheme}
    />
    <HeaderGlobalAction
      class={`header__refresh-action${isRefreshing ? " header__refresh-action--spinning" : ""}`}
      iconDescription="Refresh"
      aria-label="Refresh"
      icon={RefreshCcw}
      kind="primary"
      on:click={handleRefresh}
    />
  </HeaderUtilities>
</Header>

<style lang="scss">
  @use "../../style/variables.scss";

  .header__label {
    font-weight: 600;
  }

  .header__timestamp {
    font-weight: 400;
  }

  :global svg > * {
    fill: none;
  }

  :global(.header__theme-action :is(.bx--btn__icon, .lucide-icon) circle),
  :global(.header__theme-action :is(.bx--btn__icon, .lucide-icon) rect) {
    transition:
      transform 0.25s ease,
      fill 0.25s ease,
      stroke 0.25s ease;
    transform-box: fill-box;
    transform-origin: center;
  }

  :global(.header__theme-action[data-theme-tone="dark"] :is(.bx--btn__icon, .lucide-icon) circle) {
    transform: translateX(variables.$spacing-03 - variables.$spacing-01);
  }

  :global(.header__refresh-action--spinning :is(.bx--btn__icon, .lucide-icon)) {
    animation: spin 0.45s ease-in-out forwards;
  }

  @keyframes spin {
    from {
      transform: rotate(0deg);
    }
    to {
      transform: rotate(-360deg);
    }
  }
</style>
