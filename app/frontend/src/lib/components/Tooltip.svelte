<script lang="ts">
  import { BrowserOpenURL } from "../../../wailsjs/runtime/runtime.js"
  import { cancelTooltipHide, scheduleTooltipHide, tooltip } from "../tooltip.svelte"
  import type { TooltipState } from "../tooltip.svelte"

  let host: HTMLDivElement | null = $state(null)

  const formatBus = (value: number | null) => (value == null ? null : value.toString().padStart(3, "0"))
  const sanitizeHex = (value: string) => value.trim().replace(/^0x/i, "").toLowerCase()
  const defaultHref = "https://the-sz.com/products/usbid/"
  const buildSearchHref = (vendor: string, product: string) => {
    const params = new URLSearchParams()
    if (vendor) {
      params.set("v", sanitizeHex(vendor))
    }
    if (product) {
      params.set("p", sanitizeHex(product))
    }
    const query = params.toString()
    return query ? `${defaultHref}?${query}` : defaultHref
  }
  const buildIdLabel = (vendor: string, product: string) => (vendor && product ? `${vendor}:${product}` : vendor || product || "Unknown")

  let tooltipState: TooltipState = $derived($tooltip)
  let active = $derived(Boolean(tooltipState.visible && tooltipState.content && tooltipState.position))
  let position = $derived(tooltipState.position)
  let vendorLabel = $derived(tooltipState.content?.vendorId?.trim() ?? "")
  let productLabel = $derived(tooltipState.content?.productId?.trim() ?? "")
  let busLabel = $derived(formatBus(tooltipState.content?.bus ?? null))
  let idLabel = $derived(buildIdLabel(vendorLabel, productLabel))
  let searchHref = $derived(buildSearchHref(vendorLabel, productLabel))

  //Ensures wails will open a new browser.
  const handleLinkClick = (event: MouseEvent) => {
    event.preventDefault()
    BrowserOpenURL(searchHref)
  }
</script>

<div class="tooltip-host" bind:this={host}>
  {#if active && position}
    <div
      class="tooltip"
      role="tooltip"
      style={`top: ${position.y}px; left: ${position.x}px;`}
      onmouseenter={cancelTooltipHide}
      onmouseleave={() => scheduleTooltipHide(80)}
    >
      <div class="tooltip__header">
        <span class="tooltip__summary">Bus {busLabel ?? "—"}</span>
        <span class="tooltip__id">ID {idLabel}</span>
      </div>
      <a
        class="tooltip__link"
        href={searchHref}
        target="_blank"
        rel="noopener noreferrer"
        onclick={handleLinkClick}
      >Click to search on online device database</a>
    </div>
  {/if}
</div>

<style lang="scss">
  .tooltip-host {
    position: fixed;
    inset: 0;
    z-index: 1000;
    pointer-events: none;
  }

  .tooltip {
    position: fixed;
    transform: translate(-50%, calc(-100% - 12px));
    min-width: 14rem;
    max-width: 18rem;
    padding: 0.5rem 0.75rem;
    background: var(--color-tooltip-bg);
    color: var(--color-tooltip-text);
    border: 1px solid var(--color-tooltip-border);
    border-radius: 6px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
    font-size: 0.85rem;
    pointer-events: auto;
  }

  .tooltip::before,
  .tooltip::after {
    content: "";
    position: absolute;
    left: 50%;
    transform: translateX(-50%);
    width: 0;
    height: 0;
    border-style: solid;
  }

  .tooltip::before {
    top: 100%;
    border-width: 16px 20px 0 20px;
    border-color: var(--color-tooltip-border) transparent transparent transparent;
    z-index: -1;
  }

  .tooltip::after {
    top: calc(100% - 1px);
    border-width: 15px 19px 0 19px;
    border-color: var(--color-tooltip-bg) transparent transparent transparent;
  }

  .tooltip__header {
    display: flex;
    justify-content: space-between;
    gap: 1rem;
    font-weight: 600;
    margin-bottom: 0.35rem;
    color: inherit;
  }

  .tooltip__summary {
    white-space: nowrap;
  }

  .tooltip__id {
    white-space: nowrap;
  }

  .tooltip__link {
    display: block;
    font-size: 0.8rem;
    color: inherit;
    text-decoration: none;
  }

  .tooltip__link:hover,
  .tooltip__link:focus {
    text-decoration: underline;
  }
</style>
