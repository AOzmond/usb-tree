<script>
  import TreeNode from "../../lib/components/TreeNode.svelte"
  import { deviceTree } from "../../lib/state.svelte"
</script>

<section class="device-tree">
  {#each $deviceTree as node}
  <TreeNode {node} />
  {/each}
</section>

<style lang="scss">
  .device-tree {
    padding: 0 12px;
    min-height: 25%;
    height: 50%;
    flex-grow: 0;
    flex-shrink: 1;
    flex-basis: auto;
    resize: vertical;
    overflow: auto;
    background: var(--color-tree-bg);
    border-bottom: 1px solid var(--color-divider);
  }
</style>
