<script lang="ts">
  import { Init } from "$lib/state.svelte"
  import Header from "$lib/components/Header.svelte"
  import Tree from "$lib/components/Tree.svelte"
  import Log from "$lib/components/Log.svelte"
  import Tooltip from "$lib/components/Tooltip.svelte"
  import { Content } from "carbon-components-svelte"
  import { onMount } from "svelte"
  import {
    handlePointerDown,
    handleWindowPointerMove,
    handleWindowPointerUp,
    isResizing,
    setContentRegion,
    treeRatio,
  } from "$lib/split-layout"

  let contentRegion: HTMLDivElement | null = null

  $: setContentRegion(contentRegion)

  onMount(() => {
    Init()
  })
</script>

<div class="app">
  <Header />
  <Content id="main-content" class="content-wrapper">
    <div class="content" class:content--resizing={$isResizing} bind:this={contentRegion}>
      <div class="pane pane--tree" style={`flex: ${$treeRatio} 1 0%;`}>
        <Tree />
      </div>
      <div
        class="splitter"
        role="separator"
        aria-orientation="horizontal"
        aria-label="Resize tree and log panes"
        class:splitter--active={$isResizing}
        onpointerdown={handlePointerDown}
      >
        <span class="splitter-grip" aria-hidden="true"></span>
      </div>
      <div class="pane pane--log" style={`flex: ${1 - $treeRatio} 1 0%;`}>
        <Log />
      </div>
    </div>
  </Content>
  <Tooltip />
</div>
<svelte:window
  onpointermove={handleWindowPointerMove}
  onpointerup={handleWindowPointerUp}
  onpointercancel={handleWindowPointerUp}
/>

<style lang="scss">
  @use "variables.scss" as *;
  .app {
    display: flex;
    flex-direction: column;
    min-height: 100vh;
    color: var(--color-text);
    background: var(--color-tree-bg);
  }

  .content {
    flex: 1 1 auto;
    display: flex;
    flex-direction: column;
    min-height: 0;
  }

  .content--resizing {
    cursor: row-resize;
    user-select: none;
  }

  .pane {
    flex: 1 1 0%;
    min-height: 0;
    display: flex;
    flex-direction: column;
    background: var(--color-tree-bg);
  }

  .splitter {
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: $spacing-02 0;
    cursor: row-resize;
    flex: 0 0 auto;
    touch-action: none;
  }

  .splitter-grip {
    width: 100%;
    max-width: 160px;
    height: 2px;
    border-radius: 999px;
    background: var(--color-divider);
    transition: background 0.2s ease;
  }

  .splitter:hover .splitter-grip,
  .splitter--active .splitter-grip,
  .splitter:focus-visible .splitter-grip {
    background: var(--color-text);
  }

  .splitter:focus-visible {
    outline: 2px solid var(--color-text);
    outline-offset: 2px;
  }

  :global #main-content {
    padding: 0;
    display: flex;
    flex-direction: column;
    flex-grow: 1;
    flex-shrink: 1;
    flex-basis: auto;
    background: var(--color-tree-bg);
    height: calc(100vh - 48px);
  }

</style>
