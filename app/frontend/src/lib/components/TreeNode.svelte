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
  const iconClass = $derived(hasChildren() ? "chevron" : node.device.state)
  const collapsedClass = $derived(isCollapsed ? "collapsed" : "")

  const TreeIcon = $derived(
    hasChildren()
      ? ChevronDown
      : (iconByState[node.device.state as keyof typeof iconByState] ?? iconByState.normal)
  )

  const tooltipContent = $derived(() => ({
    bus: node.device?.bus ?? undefined,
    vendorId: node.device?.vendorId ?? undefined,
    productId: node.device?.productId ?? undefined,
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
  class={`tree-node ${iconClass}`}
  style="margin-left: {indent}rem;"
  use:tooltipTrigger={{
    getContent: tooltipContent,
  }}
>
  <div class="info" aria-expanded={!isCollapsed}>
    <TreeIcon class={`chevron ${collapsedClass}`} onclick={toggleCollapsed} />
    <a class="label" href={searchHref} onclick={handleLinkClick} aria-label="Open device info in browser">
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

    .info {
      display: inline-flex;
      align-items: center;
      gap: $spacing-02;

      .label {
        color: var(--txt-color);
        text-decoration: none;
        display: inline-flex;
        align-items: center;
        gap: $spacing-02;
      }
    }

    .speed {
      white-space: nowrap;
      align-self: flex-start;
    }

    &.added {
      color: var(--color-added);
    }

    &.removed {
      color: var(--color-removed);
    }
  }

  :global(.chevron) {
    display: inline-flex;
    transition-property: transform;
    transition-duration: 0.15s;
    transition-timing-function: ease;
    transition-delay: 0s;
  }

  :global(.collapsed) {
    transform: rotate(-90deg);
  }

  :global(.chevron svg) {
    opacity: 0.8;
    cursor: pointer;
  }

  :global(.normal svg) {
    opacity: 0.7;
  }
</style>
