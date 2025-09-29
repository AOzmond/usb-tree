<script lang="ts">
  import { hideTooltip, tooltip } from "../tooltip-state.svelte"
  import type { TooltipState } from "../tooltip-state.svelte"

  let host: HTMLDivElement | null = $state(null)

  const formatBus = (value: number | null) => (value == null ? null : value.toString().padStart(3, "0"))

  const buildIdLabel = (vendor: string, product: string) =>
    vendor && product ? `${vendor}:${product}` : vendor || product || "Unknown"

  let tooltipState: TooltipState = $derived($tooltip)
  let active = $derived(Boolean(tooltipState.visible && tooltipState.content && tooltipState.position))
  let position = $derived(tooltipState.position)
  let placement = $derived(position?.placement ?? "top")
  let isBottomPlacement = $derived(placement === "bottom")
  let vendorLabel = $derived(tooltipState.content?.vendorId?.trim() ?? "")
  let productLabel = $derived(tooltipState.content?.productId?.trim() ?? "")
  let busLabel = $derived(formatBus(tooltipState.content?.bus ?? null))
  let idLabel = $derived(buildIdLabel(vendorLabel, productLabel))
</script>

<div class="tooltip-host" bind:this={host}>
  {#if active && position}
    <div
      class="tooltip"
      class:tooltip--bottom={isBottomPlacement}
      role="tooltip"
      style={`top: ${position.y}px; left: ${position.x}px;`}
      onpointerenter={() => hideTooltip()}
    >
      <div class="tooltip__header">
        <span class="tooltip__summary">Bus {busLabel}</span>
        <span class="tooltip__id">ID {idLabel}</span>
      </div>
      <span>Click to search on online device database</span>
    </div>
  {/if}
</div>

<style lang="scss">
  @use "variables.scss" as *;

  .tooltip-host {
    position: absolute;
    inset: 0;
    z-index: 1000;
    pointer-events: none;
  }

  .tooltip {
    position: absolute;
    min-width: 14rem;
    max-width: 18rem;
    padding: $spacing-03 $spacing-04;
    background: var(--color-tooltip-bg);
    color: var(--color-tooltip-text);
    border: 1px solid var(--color-tooltip-border);
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

  .tooltip--bottom::before {
    top: auto;
    bottom: 100%;
    border-width: 0 20px 16px 20px;
    border-color: transparent transparent var(--color-tooltip-border) transparent;
  }

  .tooltip--bottom::after {
    top: auto;
    bottom: calc(100% - 1px);
    border-width: 0 19px 15px 19px;
    border-color: transparent transparent var(--color-tooltip-bg) transparent;
  }

  .tooltip__header {
    display: flex;
    justify-content: space-between;
    gap: $spacing-05;
    font-weight: 600;
    margin-bottom: $spacing-02 + $spacing-01;
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
