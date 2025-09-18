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
    min-height: 25%;
    height: 50%;
    flex: 0 1 auto;
    resize: vertical;
    overflow: auto;
    background: var(--color-tree-bg);
    border-bottom: 1px solid var(--color-divider);
    font-family: "IBM Plex Sans", sans-serif;
  }
</style>
