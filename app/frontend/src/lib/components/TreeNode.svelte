<script lang="ts">
  import TreeNode from "./TreeNode.svelte"
  import { BrowserOpenURL } from "../../../wailsjs/runtime/runtime.js"
  import type { TreeNode as TreeNodeModel } from "../models"
  import { tooltipTrigger } from "../tooltip-state.svelte"
  import { formatSpeed } from "../utilities"
  import { Plus, Minus, Dot, ChevronDown } from "@lucide/svelte"

  type Props = {
    node: TreeNodeModel
    indent?: number
  }

  let { node, indent = 0 }: Props = $props()
  let isCollapsed = $state<boolean>(false)

  const defaultHref = "https://the-sz.com/products/usbid/"
  const sanitizeHex = (value: string) => value.trim().replace(/^0x/i, "").toLowerCase()
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
  let searchHref = $derived(buildSearchHref(node.device.vendorId, node.device.productId))

  //Ensures wails will open a new browser.
  const handleLinkClick = (event: MouseEvent) => {
    event.preventDefault()
    BrowserOpenURL(searchHref)
  }

  const hasChildren = $derived(() => (node.children?.length ?? 0) > 0)

  const iconByState = {
    added: Plus,
    removed: Minus,
    normal: Dot,
  } as const

  const TreeIcon = $derived(
    hasChildren() ? ChevronDown : (iconByState[node.device.state as keyof typeof iconByState] ?? Dot)
  )

  const tooltipContent = $derived(() => ({
    bus: node.device?.bus ?? null,
    vendorId: node.device?.vendorId ?? null,
    productId: node.device?.productId ?? null,
  }))

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
      class={`tree-node__chevron${isCollapsed ? " tree-node__chevron--collapsed" : ""}`}
      onclick={toggleCollapsed}
    />
    <a
      class="tree-node__label"
      href={searchHref}
      onclick={handleLinkClick}
      tabindex="0"
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
@use '@carbon/styles/scss/spacing';

  .tree-node {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    align-items: flex-start;
    padding: spacing.$spacing-02 0;
  }

  .tree-node__info {
    display: inline-flex;
    align-items: center;
    gap: spacing.$spacing-02;
  }

  .tree-node__info--button {
    background: none;
    border: none;
    padding: 0;
    cursor: pointer;
  }

  .tree-node__label {
    all: unset;
    display: inline-flex;
    align-items: center;
    gap: spacing.$spacing-02;
  }

  :global .tree-node__chevron {
    display: inline-flex;
    transition: transform 0.15s ease;
  }

  :global .tree-node__chevron--collapsed {
    transform: rotate(-90deg);
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


</style>
