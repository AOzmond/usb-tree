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
    padding: 1rem 12px;
    display: flex;
    flex-direction: column;
    border-top: 1px solid var(--color-divider);
    overflow-y: scroll;
    overflow-x: hidden;
    background: var(--color-log-bg);
    color: var(--color-text);
    gap: 0.5rem;
  }
</style>
