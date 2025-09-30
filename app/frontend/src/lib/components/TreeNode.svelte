<script lang="ts">
  import TreeNode from "$lib/components/TreeNode.svelte"
  import { tooltipTrigger } from "$lib/tooltip-state.svelte"
  import { formatSpeed, iconByState } from "$lib/utilities"
  import type { TreeNode as TreeNodeModel } from "$lib/models"
  import { BrowserOpenURL } from "$wailsjs/runtime/runtime.js"


  import { ChevronDown } from "@lucide/svelte"

  type Props = {
    node: TreeNodeModel
    indent?: number
  }

  let { node, indent = 0 }: Props = $props()
  let isCollapsed = $state<boolean>(false)
  let searchHref = $derived(buildSearchHref(node.device.vendorId, node.device.productId))

  const defaultHref = "https://the-sz.com/products/usbid/"

  const hasChildren = $derived(() => (node.children?.length ?? 0) > 0)
  const iconClass = $derived(hasChildren() ? "ChevronDown" : (node.device.state ?? "Dot"))

  const TreeIcon = $derived(
    hasChildren()
      ? ChevronDown
      : (iconByState[node.device.state as keyof typeof iconByState] ?? iconByState.normal)
  )

  const tooltipContent = $derived(() => ({
    bus: node.device?.bus ?? null,
    vendorId: node.device?.vendorId ?? null,
    productId: node.device?.productId ?? null,
  }))

  // Ensures wails will open a new browser.
  function handleLinkClick(event: MouseEvent) {
    event.preventDefault()
    BrowserOpenURL(searchHref)
  }

  function sanitizeHex(value: string) {
    return value.trim().replace(/^0x/i, "").toLowerCase()
  }

  function buildSearchHref(vendor: string, product: string) {
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

  function toggleCollapsed() {
    if (!hasChildren()) {
      return
    }

    isCollapsed = !isCollapsed
  }
</script>

<div
  class={`tree-node ${node.device.state}`}
  style="margin-left: {indent}rem;"
  use:tooltipTrigger={{
    getContent: tooltipContent,
  }}
>
  <div class="info info--button" aria-expanded={!isCollapsed}>
    <TreeIcon
      class={`chevron${isCollapsed ? " chevron-collapsed" : ""} ${iconClass}`}
      onclick={toggleCollapsed}
    />
    <a
      class="label"
      href={searchHref}
      onclick={handleLinkClick}
      aria-label="Open device info in browser"
    >
      <span>{node.device.name}</span>
    </a>
  </div>
  <div class="speed">{formatSpeed(node.device.speed)}</div>
</div>
{#if hasChildren() && !isCollapsed}
  {#each node.children as child}
    <TreeNode node={child} indent={indent + 1} />
  {/each}
{/if}

<style lang="scss">
  @use "variables.scss" as *;

  .tree-node {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    align-items: flex-start;
    padding: $spacing-02 0;
  }

  .info {
    display: inline-flex;
    align-items: center;
    gap: $spacing-02;
  }

  .info--button {
    background: none;
    border: none;
    padding: 0;
    cursor: pointer;
  }

  .label {
    color: var(--txt-color);
    text-decoration: none;
    display: inline-flex;
    align-items: center;
    gap: $spacing-02;
  }

  .speed {
    white-space: nowrap;
    align-self: flex-start;
  }

  .added {
    color: var(--color-added);
  }

  .removed {
    color: var(--color-removed);
  }

  :global .chevron {
    display: inline-flex;
    transition-property: transform;
    transition-duration: 0.15s;
    transition-timing-function: ease;
    transition-delay: 0s;
  }

  :global .chevron-collapsed {
    transform: rotate(-90deg);
  }

  :global .chevron.ChevronDown {
    opacity: 0.8;
  }

  :global .chevron.normal {
    opacity: 0.7;
  }
</style>
