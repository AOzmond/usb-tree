<script lang="ts">
  import { Init } from "./lib/state.svelte"
  import TreeNodeView from "./lib/components/TreeNode.svelte"
  import { deviceTree, deviceLogs } from "./lib/state.svelte"
  import { Refresh } from "../wailsjs/go/main/App"
  Init()
</script>

<!-- TODO this code is just a mockup for api testing -->
<main>
  <button onclick={Refresh}>Refresh</button>
  {#if $deviceTree.length === 0}
    <div>No devices</div>
  {:else}
    <div style="width: 100%; text-align: left;">
      {#each $deviceTree as node}
        <TreeNodeView {node} indent={0} />
      {/each}
    </div>
  {/if}
  {#if $deviceLogs != null && $deviceLogs.length === 0}
    <div>No logs</div>
  {:else}
    <div style=" width: 100%; text-align: left; display: flex; flex-direction: column;">
      <br />
      <h1>LOGS</h1>
      {#each $deviceLogs as log}
        {#if log != null}
          <span
            style="color: {log.State === 'removed' ? 'red' : log.State === 'added' ? 'green' : 'inherit'};"
          >
            {log.Time}
            {log.Text}
            {log.State}
          </span>
        {/if}
      {/each}
    </div>
  {/if}
</main>

<style>
  :root {
    background: black;
    color: white;
    font-size: 14px;
  }

  main {
    height: 100vh;
    width: 100vw;
    display: flex;
    flex-direction: column;
    align-items: center;
  }
</style>
