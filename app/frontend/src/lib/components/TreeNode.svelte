<script lang="ts">
  import TreeNode from "./TreeNode.svelte"
  import type { TreeNode as TreeNodeModel } from "../models"
  import { formatSpeed } from "../utilities"
  import { Plus, Minus, Dot, ChevronDown } from '@lucide/svelte'
  import { tooltipTrigger } from "../tooltip.svelte"

  interface Props {
    node: TreeNodeModel
    indent?: number
  }

  let { node, indent = 0 }: Props = $props()
  let isCollapsed = $state<boolean>(false)

  const hasChildren = $derived(() => (node.children?.length ?? 0) > 0)
  const isAdded = $derived(() => node.device.state === "added")
  const isRemoved = $derived(() => node.device.state === "removed")
  const isNormal = $derived(() => node.device.state === "normal")
  const tooltipContent = $derived(() => ({
    bus: node.device?.bus ?? null,
    vendorId: node.device?.vendorId ?? null,
    productId: node.device?.productId ?? null,
  }))

  const toggleCollapsed = () => {
    if (!hasChildren()) {
      return
    }
    isCollapsed = !isCollapsed
  }

</script>

<div
  class="{node.device.state} TreeNode layout-row"
  style="margin-left: {indent}rem;"
  use:tooltipTrigger={{
    getContent: tooltipContent,
  }}
>
  {#if hasChildren()}
    <button
      type="button"
      class="TreeNode__info"
      aria-expanded={!isCollapsed}
      onclick={toggleCollapsed}
    >
      {#if isAdded()}
        <Plus  />
      {:else if isRemoved()}
        <Minus  />
      {/if}
      <span class="TreeNode__chevron" class:collapsed={isCollapsed}>
        <ChevronDown />
      </span>
      <div class="TreeNode__label">
        <span>{node.device.name}</span>
      </div>
    </button>
  {:else}
    <div
      class="TreeNode__info"
    >
      {#if isAdded()}
        <Plus />
      {:else if isRemoved()}
        <Minus  />
      {:else if isNormal()}
        <Dot />
      {/if}
      <div class="TreeNode__label">
        <span>{node.device.name}</span>
      </div>
    </div>
  {/if}
  <div class="TreeNode__speed">{formatSpeed(node.device.speed)}</div>
</div>
{#if hasChildren() && !isCollapsed}
  {#each node.children as child}
    <TreeNode node={child} indent={indent + 1} />
  {/each}
{/if}

<style lang="scss">
  .TreeNode {
    padding: 0.2rem 0;
  }
  .TreeNode__info {
    display: inline-flex;
    align-items: center;
    gap: 0.25rem;
  }
  button.TreeNode__info {
    background: none;
    border: none;
    padding: 0;
    cursor: pointer;
  }
  .TreeNode__label {
    display: inline-flex;
    align-items: center;
    gap: 0.25rem;
  }
  .TreeNode__chevron {
    display: inline-flex;
    transition: transform 0.15s ease;
  }
  .TreeNode__chevron.collapsed {
    transform: rotate(-90deg);
  }
  .TreeNode__speed {
    white-space: nowrap;
    align-self: flex-start;
  }
</style>
