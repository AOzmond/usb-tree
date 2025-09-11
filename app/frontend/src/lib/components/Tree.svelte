<script>
  import TreeNode from "../../lib/components/TreeNode.svelte"
  import { deviceTree } from "../../lib/state.svelte"
</script>

{#if $deviceTree.length === 0}
  <div>No devices</div>
{:else}
  <div style="width: 100%; text-align: left;">
    {#each $deviceTree as node}
      <TreeNode {node} indent={0} />
    {/each}
  </div>
{/if}
