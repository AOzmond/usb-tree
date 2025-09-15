<script lang="ts">
  import TreeNode from "./TreeNode.svelte"
  import type { TreeNode as TreeNodeModel } from "../models"

  const SPEED_MAP: Record<string, number> = {
    low: 1.5,
    full: 12,
    high: 480,
    super: 5000,
  }

  const formatSpeed = (rawSpeed?: string): string => {
    if (!rawSpeed) {
      return "unknown"
    }

    const normalized = rawSpeed.toLowerCase()
    const value = SPEED_MAP[normalized]

    if (value === undefined) {
      return rawSpeed
    }

    if (value >= 1000) {
      const gbps = value / 1000
      return `${Number.isInteger(gbps) ? gbps : gbps} Gbps`
    }

    return `${Number.isInteger(value) ? value : value} Mbps`
  }

  interface Props {
    node: TreeNodeModel
    indent?: number
  }

  let { node, indent = 0 }: Props = $props()
</script>

<div class="{node.device.state} TreeNode" style="margin-left: {indent}rem;">
  <!-- TODO Add symbol -->
  <div>{node.device.name}</div>
  <div>{formatSpeed(node.device.speed)}</div>
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
  }
  .added {
    color: var(--color-added);
  }
  .removed {
    color: var(--color-removed);
  }
</style>
