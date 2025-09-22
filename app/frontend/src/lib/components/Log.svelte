<script lang="ts">
  import { onMount } from "svelte"
  import { deviceLogs } from "../../lib/state.svelte"
  import LogRow from "../components/LogRow.svelte"

  let container: HTMLDivElement | null = $state(null)
  let autoScroll = $state(true)

  function scrollToBottom() {
    if (!container) {
      return
    }
    container.scrollTop = container.scrollHeight
  }

  function handleScroll() {
    if (!container) {
      return
    }
    const distanceFromBottom = container.scrollHeight - container.scrollTop - container.clientHeight
    autoScroll = distanceFromBottom <= 2
  }

  $effect(() => {
    $deviceLogs
    autoScroll
    if (autoScroll) {
      scrollToBottom()
    }
  })

  onMount(() => {
    scrollToBottom()
  })
</script>

<div class="log-panel" bind:this={container} onscroll={handleScroll}>
  {#each $deviceLogs as log}
      <LogRow {log} />
  {/each}
</div>

<style lang="scss">
  .log-panel {
    padding: 12px;
    display: flex;
    flex-direction: column;
    min-height: 25%;
    flex-grow: 1;
    flex-shrink: 1;
    flex-basis: 50%;
    border-top: 1px solid var(--color-divider);
    overflow-y: auto;
    overflow-x: hidden;
    background: var(--color-log-bg);
    gap: 0.5rem;
  }
</style>
