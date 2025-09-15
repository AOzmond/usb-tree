<script>
  import TreeNode from "../../lib/components/TreeNode.svelte"
  import { deviceTree } from "../../lib/state.svelte"
</script>

<div>
  {#each $deviceTree as node}
    <span><TreeNode {node} indent={0} /></span>
  {/each}
</div>

<style lang="scss">
  div {
    padding: 0 12px;
    resize: vertical;
    overflow: scroll;
    background: var(--color-tree-bg);
    color: var(--color-text);
    border-bottom: 1px solid var(--color-divider);
  }
</style>
