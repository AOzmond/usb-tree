<script>
  import { deviceLogs } from "../../lib/state.svelte"
</script>

{#if $deviceLogs != null && $deviceLogs.length === 0}
  <div>No logs</div>
{:else}
  <div style=" width: 100%; text-align: left; display: flex; flex-direction: column;">
    {#each $deviceLogs as log}
      {#if log != null}
        <span style="color: {log.State === 'removed' ? 'red' : log.State === 'added' ? 'green' : 'inherit'};">
          {log.Time}
          {log.Text}
          {log.State}
        </span>
      {/if}
    {/each}
  </div>
{/if}
