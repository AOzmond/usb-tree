<script lang="ts">
  import TreeNode from "./TreeNode.svelte"
  export let node
  export let indent = 0
</script>

<div class="{node.device.state} TreeNode" style="margin-left: {indent}rem;">
  <!-- TODO Add symbol -->
  <div>{node.device.name}</div>
  <!-- TODO convert speed to number -->
  <div>{node.device.speed}</div>
</div>
{#each node.children as child}
  <TreeNode node={child} indent={indent + 1} />
{/each}

<style lang="scss">
  .TreeNode {
    padding: 0.2rem 0;
    display: flex;
    justify-content: space-between;
  }
  .added {
    color: green;
  }
  .removed {
    color: red;
  }
</style>
