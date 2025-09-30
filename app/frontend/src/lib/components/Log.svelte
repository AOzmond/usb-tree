<script lang="ts">
  import { onMount } from "svelte"

  import { deviceLogs } from "$lib/state.svelte"
  import LogRow from "$lib/components/LogRow.svelte"

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

  // Scrolls log to bottom on Log change if autoScroll is true.
  $effect(() => {
    $deviceLogs
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
  @use "variables.scss" as *;

  .log-panel {
    padding: $spacing-04;
    display: flex;
    flex-direction: column;
    flex-grow: 1;
    flex-shrink: 1;
    flex-basis: auto;
    border-top-width: 1px;
    border-top-style: solid;
    border-top-color: var(--color-divider);
    overflow-y: auto;
    overflow-x: hidden;
    background: var(--color-log-bg);
    gap: $spacing-03;
  }
</style>
