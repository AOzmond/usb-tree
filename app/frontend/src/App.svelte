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
  } from "$lib/splitLayout"

  let contentRegion: HTMLDivElement | null = null

  $: setContentRegion(contentRegion)

  onMount(() => {
    Init()
  })
</script>

<div class="app">
  <Header />
  <Content class="app__content-wrapper">
    <div class="app__content" class:app__content--resizing={$isResizing} bind:this={contentRegion}>
      <div class="app__pane app__pane--tree" style={`flex: ${$treeRatio} 1 0%;`}>
        <Tree />
      </div>
      <div
        class="app__splitter"
        role="separator"
        aria-orientation="horizontal"
        aria-label="Resize tree and log panes"
        tabindex="0"
        class:app__splitter--active={$isResizing}
        on:pointerdown={handlePointerDown}
      >
        <span class="app__splitter-grip" aria-hidden="true"></span>
      </div>
      <div class="app__pane app__pane--log" style={`flex: ${1 - $treeRatio} 1 0%;`}>
        <Log />
      </div>
    </div>
  </Content>
  <Tooltip />
</div>
<svelte:window
  on:pointermove={handleWindowPointerMove}
  on:pointerup={handleWindowPointerUp}
  on:pointercancel={handleWindowPointerUp}
/>

<style lang="scss">
  @use "./style/variables.scss" as variables;
  .app {
    display: flex;
    flex-direction: column;
    min-height: 100vh;
    color: var(--color-text);
    background: var(--color-tree-bg);
  }

  .app__content {
    flex: 1 1 auto;
    display: flex;
    flex-direction: column;
    min-height: 0;
  }

  .app__content--resizing {
    cursor: row-resize;
    user-select: none;
  }

  .app__pane {
    flex: 1 1 0%;
    min-height: 0;
    display: flex;
    flex-direction: column;
    background: var(--color-tree-bg);
  }

  .app__splitter {
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: variables.$spacing-02 0;
    cursor: row-resize;
    flex: 0 0 auto;
    touch-action: none;
  }

  .app__splitter-grip {
    width: 100%;
    max-width: 160px;
    height: 2px;
    border-radius: 999px;
    background: var(--color-divider);
    transition: background 0.2s ease;
  }

  .app__splitter:hover .app__splitter-grip,
  .app__splitter--active .app__splitter-grip,
  .app__splitter:focus-visible .app__splitter-grip {
    background: var(--color-text);
  }

  .app__splitter:focus-visible {
    outline: 2px solid var(--color-text);
    outline-offset: 2px;
  }
</style>
