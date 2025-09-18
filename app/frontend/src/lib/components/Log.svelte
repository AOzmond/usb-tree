<script lang="ts">
  import { onMount } from "svelte"
  import { deviceLogs } from "../../lib/state.svelte"
  import LogRow from "../components/LogRow.svelte"

  let container: HTMLDivElement | null = $state(null)
  let autoScroll = $state(true)
  let visibleLogs = $derived(($deviceLogs ?? []).filter(Boolean))

  const scrollToBottom = () => {
    if (!container) {
      return
    }
    container.scrollTop = container.scrollHeight
  }

  onMount(() => {
    scrollToBottom()
  })

  const handleScroll = () => {
    if (!container) {
      return
    }
    const distanceFromBottom = container.scrollHeight - container.scrollTop - container.clientHeight
    autoScroll = distanceFromBottom <= 2
  }

  $effect(() => {
    visibleLogs
    autoScroll
    if (!autoScroll) {
      return
    }
    scrollToBottom()
  })
</script>

<div class="log-container" bind:this={container} onscroll={handleScroll}>
  {#each visibleLogs as log}
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
    min-height: 25%;
    flex: 1 1 50%;
    border-top: 1px solid var(--color-divider);
    overflow-y: auto;
    overflow-x: hidden;
    background: var(--color-log-bg);
    gap: 0.5rem;
    font-family: "JetBrains Mono", monospace;
  }
</style>
