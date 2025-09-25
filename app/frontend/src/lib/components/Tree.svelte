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
  @use "../../style/variables.scss";

  .device-tree {
    padding: 0 variables.$spacing-04;
    flex-grow: 1;
    flex-shrink: 1;
    flex-basis: auto;
    min-height: 0;
    overflow: auto;
    background: var(--color-tree-bg);
    border-bottom: 1px solid var(--color-divider);
  }
</style>
