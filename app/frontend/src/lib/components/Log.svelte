<script>
  import { deviceLogs } from "../../lib/state.svelte"
  import LogRow from "../components/LogRow.svelte"
</script>

<div class="log-container">
  {#each $deviceLogs as log}
    {#if log != null}
      <LogRow {log} />
    {/if}
  {/each}
</div>

<style lang="scss">
  .log-container {
    min-height: 25%;
    max-height: 75%;
    display: flex;
    flex-direction: column;
    border-top: 1px solid #333;
    overflow-y: scroll;
    overflow-x: hidden;
  }
</style>
