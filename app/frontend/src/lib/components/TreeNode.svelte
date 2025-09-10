<script lang="ts">
  import TreeNode from "./TreeNode.svelte";
  export let node;
  export let indent = 0;
</script>

<!-- TODO this code is just a mockup for api testing -->
<div class={node.device.state} style="margin-left: {indent}em;">
  {node.device.name}
</div>
{#each node.children
  .slice()
  .sort((a, b) => a.device.bus - b.device.bus) as child}
  <TreeNode node={child} indent={indent + 6} />
{/each}

<style>
  .added {
    color: green;
  }
  .removed {
    color: red;
  }
</style>
