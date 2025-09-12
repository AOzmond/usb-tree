<script>
  import TreeNode from "../../lib/components/TreeNode.svelte"
  import { deviceTree } from "../../lib/state.svelte"
</script>

<div >
  {#each $deviceTree as node}
    <span><TreeNode {node} indent={0} /></span>
  {/each}
</div>

<style lang="scss">
  div {
    padding: 0 12px;
    min-height: 25%;
    resize: vertical;
    overflow: scroll;
  }
</style>
