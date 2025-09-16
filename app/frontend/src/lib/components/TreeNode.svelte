<script lang="ts">
  import TreeNode from "./TreeNode.svelte"
  import type { TreeNode as TreeNodeModel } from "../models"
  import { formatSpeed } from "../utilities"

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
      <svg
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
        class="TreeNode__icon TreeNode__icon--state"
        aria-hidden="true"
      >
        <path d="M5 12h14" />
        <path d="M12 5v14" />
      </svg>
    {:else if isRemoved()}
      <svg
        xmlns="http://www.w3.org/2000/svg"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
        class="TreeNode__icon TreeNode__icon--state"
        aria-hidden="true"
      >
        <path d="M5 12h14" />
      </svg>
    {:else if isNormal() && !hasChildren()}
      <svg
        xmlns="http://www.w3.org/2000/svg"
        width="24"
        height="24"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
        class="TreeNode__icon TreeNode__icon--state"
        aria-hidden="true"
      >
        <circle cx="12.1" cy="12.1" r="1" />
      </svg>
    {/if}
    <div class="TreeNode__label">
      {#if hasChildren()}
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="24"
          height="24"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="TreeNode__icon"
          aria-hidden="true"
        >
          <path d="m6 9 6 6 6-6" />
        </svg>
      {/if}
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
    gap: 0.35rem;
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
