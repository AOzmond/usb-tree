<script>
  import TreeNode from "$lib/components/TreeNode.svelte"
  import { deviceTree } from "$lib/state.svelte"
</script>

<section class="device-tree">
  {#each $deviceTree as node}
    <TreeNode {node} />
  {/each}
</section>

<style lang="scss">
  @use "variables.scss" as *;

  .device-tree {
    padding: 0 $spacing-04;
    flex-grow: 1;
    flex-shrink: 1;
    flex-basis: auto;
    min-height: 0;
    overflow: auto;
    background-color: var(--color-tree-bg);
    border-bottom-width: 1px;
    border-bottom-style: solid;
    border-bottom-color: var(--color-divider);
  }
</style>
