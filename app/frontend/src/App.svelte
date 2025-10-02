<script lang="ts">
  import { onMount } from "svelte"
  import { Content } from "carbon-components-svelte"

  import { Init } from "$lib/state.svelte"
  import Header from "$lib/components/Header.svelte"
  import Tree from "$lib/components/Tree.svelte"
  import Log from "$lib/components/Log.svelte"
  import Tooltip from "$lib/components/Tooltip.svelte"
  import {
    handlePointerDown,
    handleWindowPointerMove,
    handleWindowPointerUp,
    isResizing,
    setContentRegion,
    treeRatio,
  } from "$lib/split-layout"

  let contentRegion: HTMLDivElement

  $: setContentRegion(contentRegion)

  onMount(() => {
    Init()
  })
</script>

<div class="app">
  <Header />
  <Content id="main-content" class="content-wrapper">
    <div class="content" class:content-resizing={$isResizing} bind:this={contentRegion}>
      <div class="pane pane-tree" style={`flex-grow: ${$treeRatio};`}>
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
      <div class="pane pane-log" style={`flex-grow: ${1 - $treeRatio};`}>
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
    background-color: var(--color-tree-bg);
  }

  .content {
    flex-grow: 1;
    flex-shrink: 1;
    flex-basis: auto;
    display: flex;
    flex-direction: column;
    min-height: 0;
  }

  .content-resizing {
    cursor: row-resize;
    user-select: none;
  }

  .pane {
    flex-grow: 1;
    flex-shrink: 1;
    flex-basis: 0;
    min-height: 25%;
    display: flex;
    flex-direction: column;
    background-color: var(--color-tree-bg);
  }

  .splitter {
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: $spacing-02 0;
    cursor: row-resize;
    flex-grow: 0;
    flex-shrink: 0;
    flex-basis: auto;
    touch-action: none;
  }

  .splitter-grip {
    width: 100%;
    max-width: 160px;
    height: 2px;
    border-radius: 999px;
    background-color: var(--color-divider);
    transition-duration: 0.15s;
    transition-timing-function: ease;
    transition-delay: 0s;
  }

  .splitter:hover .splitter-grip,
  .splitter-active .splitter-grip,
  .splitter:focus-visible .splitter-grip {
    background-color: var(--color-text);
  }

  .splitter:focus-visible {
    outline-width: 2px;
    outline-style: solid;
    outline-color: var(--color-text);
    outline-offset: 2px;
  }

  :global(#main-content) {
    padding: 0;
    display: flex;
    flex-direction: column;
    flex-grow: 1;
    flex-shrink: 1;
    flex-basis: auto;
    background-color: var(--color-tree-bg);
    height: calc(100vh - $headerHeight);
  }
</style>
