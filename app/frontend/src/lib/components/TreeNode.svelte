<script lang="ts">
  import TreeNode from "./TreeNode.svelte"
  import type { TreeNode as TreeNodeModel } from "../models"
  import { formatSpeed } from "../utilities"
  import Added from "../../assets/svgs/added.svg?component"
  import Removed from "../../assets/svgs/removed.svg?component"
  import Normal from "../../assets/svgs/normal.svg?component"
  import DownChevron from "../../assets/svgs/downchevron.svg?component"

  interface Props {
    node: TreeNodeModel
    indent?: number
  }

  let { node, indent = 0 }: Props = $props()

  const hasChildren = $derived(() => (node.children?.length ?? 0) > 0)
  const isAdded = $derived(() => node.device.state === "added")
  const isRemoved = $derived(() => node.device.state === "removed")
  const isNormal = $derived(() => node.device.state === "normal")
</script>

<div class="{node.device.state} TreeNode" style="margin-left: {indent}rem;">
  <div class="TreeNode__info">
    {#if isAdded()}
      <Added width="24" />
    {:else if isRemoved()}
      <Removed width="24" />
    {:else if isNormal() && !hasChildren()}
      <Normal />
    {/if}
    {#if hasChildren()}
      <DownChevron />
    {/if}
    <div class="TreeNode__label">
      <span>{node.device.name}</span>
    </div>
  </div>
  <div class="TreeNode__speed">{formatSpeed(node.device.speed)}</div>
</div>
{#each node.children as child}
  <TreeNode node={child} indent={indent + 1} />
{/each}

<style lang="scss">
  .TreeNode {
    padding: 0.2rem 0;
    display: flex;
    justify-content: space-between;
    color: var(--color-text);
    align-items: flex-start;
    gap: 0.5rem;
  }
  .TreeNode__info {
    display: inline-flex;
    align-items: center;
  }
  .TreeNode__label {
    display: inline-flex;
    align-items: center;
    gap: 0.25rem;
  }
  .TreeNode__icon {
    width: 1rem;
    height: 1rem;
  }
  .TreeNode__speed {
    white-space: nowrap;
    align-self: flex-start;
  }
  .added {
    color: var(--color-added);
  }
  .removed {
    color: var(--color-removed);
  }
</style>
