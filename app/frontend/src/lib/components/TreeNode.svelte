<script lang="ts">
  import TreeNode from "./TreeNode.svelte"
  import { tooltipTrigger } from "../tooltip-state.svelte"
  import { formatSpeed, iconByState } from "../utilities"
  import type { TreeNode as TreeNodeModel } from "../models"
  import { BrowserOpenURL } from "../../../wailsjs/runtime/runtime.js"
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

  //Ensures wails will open a new browser.
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
  class={`tree-node tree-node--${node.device.state}`}
  style="margin-left: {indent}rem;"
  use:tooltipTrigger={{
    getContent: tooltipContent,
  }}
>
  <div class="tree-node__info tree-node__info--button" aria-expanded={!isCollapsed}>
    <TreeIcon
      class={`tree-node__chevron${isCollapsed ? " tree-node__chevron--collapsed" : ""} ${iconClass}`}
      onclick={toggleCollapsed}
    />
    <a
      class="tree-node__label"
      href={searchHref}
      onclick={handleLinkClick}
      aria-label="Open device info in browser"
    >
      <span>{node.device.name}</span>
    </a>
  </div>
  <div class="tree-node__speed">{formatSpeed(node.device.speed)}</div>
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

  .tree-node__info {
    display: inline-flex;
    align-items: center;
    gap: $spacing-02;
  }

  .tree-node__info--button {
    background: none;
    border: none;
    padding: 0;
    cursor: pointer;
  }

  .tree-node__label {
    color: var(--txt-color);
    text-decoration: none;
    display: inline-flex;
    align-items: center;
    gap: $spacing-02;
  }

  .tree-node__speed {
    white-space: nowrap;
    align-self: flex-start;
  }

  .tree-node--added {
    color: var(--color-added);
  }

  .tree-node--removed {
    color: var(--color-removed);
  }

  :global .tree-node__chevron {
    display: inline-flex;
    transition: transform 0.15s ease;
  }

  :global .tree-node__chevron--collapsed {
    transform: rotate(-90deg);
  }

  :global .tree-node__chevron.ChevronDown {
    opacity: 0.8;
  }

  :global .tree-node__chevron.normal {
    opacity: 0.7;
  }
</style>
